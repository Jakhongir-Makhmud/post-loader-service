package structs

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
	ErrInternal   = errors.New("something went wrong")
	ErrNoData = errors.New("no data")
	ErrTypeCast = errors.New("can not cast type")
)
