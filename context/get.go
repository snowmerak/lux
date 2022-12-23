package context

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
)

func (l *LuxContext) GetFormFile(name string) (multipart.File, *multipart.FileHeader, error) {
	return l.Request.FormFile(name)
}

func (l *LuxContext) GetMultipartFile(name string, maxMemoryBytes int64) ([]*multipart.FileHeader, error) {
	if l.Request.MultipartForm == nil {
		err := l.Request.ParseMultipartForm(maxMemoryBytes)
		if err != nil {
			return nil, err
		}
	}
	return l.Request.MultipartForm.File[name], nil
}

func (l *LuxContext) SaveFile(file multipart.File, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	return err
}

func (l *LuxContext) SaveMultipartFile(headers []*multipart.FileHeader, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, header := range headers {
		file, err := header.Open()
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := io.Copy(f, file); err != nil {
			return err
		}
	}
	return nil
}

func (l *LuxContext) GetFormData(key string) string {
	return l.Request.FormValue(key)
}

func (l *LuxContext) GetURLQuery(key string) string {
	return l.Request.URL.Query().Get(key)
}

func (l *LuxContext) GetPathVariable(key string) string {
	return l.RouteParams.ByName(key)
}

func (l *LuxContext) GetBody() ([]byte, error) {
	data, err := io.ReadAll(l.Request.Body)
	if err != nil {
		return nil, err
	}
	l.Request.Body.Close()
	return data, nil
}

func (l *LuxContext) GetBodyReader() io.ReadCloser {
	return l.Request.Body
}

func (l *LuxContext) GetCookie(key string) string {
	return l.Request.Header.Get("Cookie")
}

func (l *LuxContext) GetRemoteAddress() string {
	return l.Request.RemoteAddr
}

func (l *LuxContext) GetRemoteIP() string {
	i := 0
	end := -1
	addr := l.Request.RemoteAddr
	for i < len(addr) {
		if addr[i] == ':' {
			end = i
		}
		i++
	}
	if end == -1 {
		return ""
	}
	return addr[:end]
}

func (l *LuxContext) GetRemotePort() string {
	i := 0
	end := -1
	addr := l.Request.RemoteAddr
	for i < len(addr) {
		if addr[i] == ':' {
			end = i
		}
		i++
	}
	if end == -1 {
		return ""
	}
	return addr[end+1:]
}

func (l *LuxContext) ParseJSON(v interface{}) error {
	decoder := json.NewDecoder(l.Request.Body)
	if err := decoder.Decode(v); err != nil {
		return err
	}
	if err := l.Request.Body.Close(); err != nil {
		return err
	}
	return nil
}
