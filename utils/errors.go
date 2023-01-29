package utils

import "errors"

var (
	ErrInvalidFlag        = errors.New("invalid flag arguments")
	ErrDBConnectionFailed = errors.New("db connection failed")
	ErrInvalidDateFormat  = errors.New("invalid date format of start-date or end-date")
	ErrDBQueryFail        = errors.New("execute db query failed")
	ErrInvalidCategory    = errors.New("unrecognized category flag value")
	ErrFileNotExist       = errors.New("no such file or directory")
)
