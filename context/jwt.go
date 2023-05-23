package context

import (
	"crypto/cipher"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

const RefreshToken = "refresh-token"

type JWTConfig struct {
	SigningKey       []byte
	SigningMethod    jwt.SigningMethod
	EncryptionKey    []byte
	EncryptionMethod cipher.AEAD
	Domain           string
	Path             string
}

type JWT struct {
	response         *Response
	request          *http.Request
	encryptionNonce  []byte
	encryptionMethod cipher.AEAD
	signingKey       []byte
	signingMethod    jwt.SigningMethod
	domain           string
	path             string
}

func (l *LuxContext) JWT() *JWT {
	if l.JWTConfig == nil {
		return nil
	}

	j := new(JWT)
	j.response = l.Response
	j.request = l.Request
	j.encryptionNonce = l.JWTConfig.EncryptionKey
	j.encryptionMethod = l.JWTConfig.EncryptionMethod
	j.signingKey = l.JWTConfig.SigningKey
	j.signingMethod = l.JWTConfig.SigningMethod
	j.domain = l.JWTConfig.Domain
	j.path = l.JWTConfig.Path

	return j
}

func (j *JWT) SetRefreshTokenWithClaims(claims jwt.Claims) error {
	value, err := j.MakeRefreshToken(claims)
	if err != nil {
		return err
	}

	ck := new(http.Cookie)
	ck.Name = RefreshToken
	ck.Domain = j.domain
	ck.Path = j.path
	ck.HttpOnly = true
	ck.Secure = true
	ck.SameSite = http.SameSiteStrictMode

	ck.Value = value

	j.response.Header().Add("Set-Cookie", ck.String())

	return nil
}

func (j *JWT) MakeRefreshToken(claims jwt.Claims) (string, error) {
	t := jwt.NewWithClaims(j.signingMethod, claims)

	signed, err := t.SignedString(j.signingKey)
	if err != nil {
		return "", err
	}

	if j.encryptionMethod == nil {
		return signed, nil
	}

	encrypted := j.encryptionMethod.Seal(nil, j.encryptionNonce, []byte(signed), nil)

	return string(encrypted), nil
}

var errInvalidToken = errors.New("invalid token error")

func (j *JWT) GetRefreshTokenFromCookie() (jwt.Claims, error) {
	ck, err := j.request.Cookie(RefreshToken)
	if err != nil {
		return nil, err
	}

	token := []byte(ck.Value)

	return j.ParseRefreshToken(token)
}

func (j *JWT) ParseRefreshToken(token []byte) (jwt.Claims, error) {
	if j.encryptionMethod != nil {
		decrypted, err := j.encryptionMethod.Open(nil, j.encryptionNonce, token, nil)
		if err != nil {
			return nil, err
		}
		token = decrypted
	}

	tk, err := jwt.Parse(string(token), func(token *jwt.Token) (interface{}, error) {
		return j.signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !tk.Valid {
		return nil, errInvalidToken
	}

	return tk.Claims, nil
}

func (j *JWT) SetAccessToken(claims jwt.Claims) error {
	value, err := j.MakeAccessToken(claims)
	if err != nil {
		return err
	}

	j.response.Header().Add("Authorization", "Bearer "+value)

	return nil
}

func (j *JWT) MakeAccessToken(claims jwt.Claims) (string, error) {
	tk := jwt.NewWithClaims(j.signingMethod, claims)

	signed, err := tk.SignedString(j.signingKey)
	if err != nil {
		return "", err
	}

	if j.encryptionMethod != nil {
		value := j.encryptionMethod.Seal(nil, j.encryptionNonce, []byte(signed), nil)
		signed = string(value)
	}

	return signed, nil
}

func (j *JWT) GetAccessToken() (jwt.Claims, error) {
	value := j.request.Header.Get("Authorization")
	value = strings.TrimPrefix(value, "Bearer")
	value = strings.TrimSpace(value)

	return j.ParseAccessToken([]byte(value))
}

func (j *JWT) ParseAccessToken(tk []byte) (jwt.Claims, error) {
	if j.encryptionMethod != nil {
		decrypted, err := j.encryptionMethod.Open(nil, j.encryptionNonce, tk, nil)
		if err != nil {
			return nil, err
		}
		tk = decrypted
	}

	t, err := jwt.Parse(string(tk), func(token *jwt.Token) (interface{}, error) {
		return j.signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, errInvalidToken
	}

	return t.Claims, nil
}
