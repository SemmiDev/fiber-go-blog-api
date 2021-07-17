package entity

import "errors"

var (
	ErrRequiredEmail    = errors.New("required email")
	ErrInvalidEmail     = errors.New("invalid email")
	ErrEmailRule        = errors.New("password should be at least 6 characters")
	ErrRequiredUsername = errors.New("required username")
	ErrRequiredPassword = errors.New("required password")
	ErrRequiredTitle    = errors.New("required title")
	ErrRequiredContent  = errors.New("required content")
	ErrRequiredAuthor   = errors.New("required author")
	ErrRequiredComment  = errors.New("required comment")
)
