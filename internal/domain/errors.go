package domain

import "errors"

var (
	ErrBlindNotFound      = errors.New("blind not found")
	ErrBlindAlreadyExists = errors.New("blind already exists")
	ErrInvalidCoordinates = errors.New("invalid coordinates")
	ErrInvalidRadius      = errors.New("invalid radius")
)