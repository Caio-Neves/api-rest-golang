package entities

var (
	BAD_REQUEST           = "Bad Request"
	INTERNAL_SERVER_ERROR = "Internal Server Error"
	NOT_FOUND             = "Not Found"
	CONFLICT              = "Conflict"
	UNAUTHORIZED          = "Unauthorized"
	FORBIDDEN             = "Forbidden"
	NOT_IMPLEMENTED       = "Not Implemented"
	NO_CONTENT            = "No Content"
)

type Error struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Err       error  `json:"-"`
	Operation string `json:"-"`
}

func (e *Error) Error() string {
	return e.Message
}

func newError(code string, message string, err error, operation string) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Err:       err,
		Operation: operation,
	}
}

func NewBadRequestError(err error, message string, operation string) *Error {
	return newError(BAD_REQUEST, message, err, operation)
}

func NewInternalServerErrorError(err error, operation string) *Error {
	return newError(INTERNAL_SERVER_ERROR, "Um erro inesperado aconteceu, tente novamente mais tarde", err, operation)
}

func NewNotFoundError(err error, message string, operation string) *Error {
	return newError(NOT_FOUND, message, err, operation)
}

func NewConflictError(err error, message string, operation string) *Error {
	return newError(CONFLICT, message, err, operation)
}

func NewUnauthorizedError(err error, message string, operation string) *Error {
	return newError(UNAUTHORIZED, message, err, operation)
}

func NewForbiddenError(err error, message string, operation string) *Error {
	return newError(FORBIDDEN, message, err, operation)
}

func NewNotImplementedError(err error, message string, operation string) *Error {
	return newError(NOT_IMPLEMENTED, message, err, operation)
}
