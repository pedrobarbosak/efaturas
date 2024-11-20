package errors

type Type int

const (
	FatalError Type = iota
	InputError
	NotFoundError
	ForbiddenError
	UnauthorizedError
	UnprocessableEntityError
	ConflictError
	UnsupportedMediaType
)
