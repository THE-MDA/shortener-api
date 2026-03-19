package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists = errors.New("url exists")
	ErrURLNotExists = errors.New("url does not exist")
)