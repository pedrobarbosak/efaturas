package errors

import "github.com/pedrobarbosak/go-errors"

type Error struct {
	Err  error
	Type Type
}

func (e Error) Error() string {
	return e.Err.Error()
}

func _new(eType Type, args ...any) error {
	return &Error{
		Err:  errors.NewCustom(2, args),
		Type: eType,
	}
}

func New(args ...any) error {
	return _new(FatalError, args...)
}

func NewInput(args ...any) error {
	return _new(InputError, args...)
}

func NewNotFound(args ...any) error {
	return _new(NotFoundError, args...)
}

func NewFatal(args ...any) error {
	return _new(FatalError, args...)
}

func NewConflict(args ...any) error {
	return _new(ConflictError, args...)
}

func NewForbidden(args ...any) error {
	return _new(ForbiddenError, args...)
}

func NewUnauthorized(args ...any) error {
	return _new(UnauthorizedError, args...)
}

func NewUnsupported(args ...any) error {
	return _new(UnsupportedMediaType, args...)
}

func NewUnprocessable(args ...any) error {
	return _new(UnprocessableEntityError, args...)
}

func GetCode(err error) Type {
	if err == nil {
		return FatalError
	}

	if ourErr, ok := err.(*Error); ok {
		return ourErr.Type
	}

	return FatalError
}

func Is(err error, eType Type) bool {
	if err == nil {
		return false
	}

	if ourErr, ok := err.(*Error); ok {
		return ourErr.Type == eType
	}

	return false
}
