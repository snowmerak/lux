package session

import (
	"errors"

	"github.com/Workiva/go-datastructures/trie/ctrie"
)

// Local is a session of app scope
type Local struct {
	data ctrie.Ctrie
}

// NewLocal create new Local instance
func NewLocal() *Local {
	return &Local{
		data: *ctrie.New(nil),
	}
}

// errNotExists is Not Exists Error value
var errNotExists = errors.New("not exists")

// IsNotExists is a function to check given error is Not Exists Error
func IsNotExists(err error) bool {
	return errors.Is(err, errNotExists)
}

// errNotMatchType is Not Match Type Error value
var errNotMatchType = errors.New("not match type")

// IsNotMatchType is a function to check given error is Not Match Type Error
func IsNotMatchType(err error) bool {
	return errors.Is(err, errNotMatchType)
}

// GetLocal is getting the value of given key.
// if not exists, return Not Exists Error.
// if not match type with generic, return Not Match Type Error.
func GetLocal[T any](localSession *Local, key []byte) (*T, error) {
	data, ok := localSession.data.Lookup(key)
	if !ok {
		return nil, errNotExists
	}
	t, ok := data.(*T)
	if !ok {
		return nil, errNotMatchType
	}
	return t, nil
}

// errAlreadyExists is Already Exists Error value
var errAlreadyExists = errors.New("already exists")

// IsAlreadyExists is a function to check given error is Already Exists Error
func IsAlreadyExists(err error) bool {
	return errors.Is(err, errAlreadyExists)
}

// SetLocal is setting value of key into local session.
// if already exists key value pair in local session, return Already Exists Error.
func SetLocal[T any](localSession *Local, key []byte, value T) error {
	if _, ok := localSession.data.Lookup(key); ok {
		return errAlreadyExists
	}
	localSession.data.Insert(key, &value)
	return nil
}

// RemoveLocal is deleting key value pair in local session.
// if not exists, return Not Exists Error.
// if not match type with generic, return Not Match Type Error.
func RemoveLocal[T any](localSession *Local, key []byte) (*T, error) {
	data, ok := localSession.data.Remove(key)
	if !ok {
		return nil, errNotExists
	}
	t, ok := data.(*T)
	if !ok {
		return nil, errNotMatchType
	}
	return t, nil
}
