package repository

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrBuildingSql      = errors.New("error building sql")
	ErrExecutingSql     = errors.New("error executing sql")
	ErrRetrievingData   = errors.New("error retrieving data")
	ErrCreatingTx       = errors.New("error creating tx")
	ErrConcurrentUpdate = errors.New("error concurrent update")
	ErrConcurrentDelete = errors.New("error concurrent delete")
	ErrInvalidInput     = errors.New("invalid input")
)
