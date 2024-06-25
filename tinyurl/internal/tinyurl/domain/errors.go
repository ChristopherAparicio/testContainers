package domain

import (
	tinyError "github.com/christapa/tinyurl/pkg/error"
)

func NewInvalidInputError(message string) error {
	return tinyError.New(tinyError.InvalidArgument, message)
}

func NewInternalError(message string) error {
	return tinyError.New(tinyError.Internal, message)
}

func NewNotFoundError() error {
	return tinyError.New(tinyError.NotFound, "URL not found")
}
