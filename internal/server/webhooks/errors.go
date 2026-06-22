package webhooks

import "errors"

// ErrEmptyUserID is returned when an empty user ID is passed to a function that requires one.
var ErrEmptyUserID = errors.New("userID is required")
