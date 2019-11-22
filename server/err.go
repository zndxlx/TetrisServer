package server

import (
    "errors"
)

var (
    ErrBadRoute  = errors.New("bad route")
    ErrWrongType = errors.New("wrong type")
    ErrNotFound  = errors.New("not found")

    ErrPlayerNotFound = errors.New("player not found")
)
