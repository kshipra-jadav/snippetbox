package models

import "errors"

var ErrNoRecords = errors.New("model: No Matching Records Found")

var ErrInvalidCredentials = errors.New("model: invalid credentials.")

var ErrDuplicateEmail = errors.New("model: duplicate email")
