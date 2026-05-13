package apperror

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyInUse = errors.New("email already in use")

	ErrInvalidEmailOrPassword = errors.New("invalid email or password")

	ErrRefreshTokenExpiredOrNotFound = errors.New("token expired or not found")

	ErrFileNotFound        = errors.New("file not found")
	ErrFileTooLarge        = errors.New("file too large")
	ErrUnsupportedMimeType = errors.New("file type unsupported")
	ErrForbiddenFileAccess = errors.New("no access to file")
)
