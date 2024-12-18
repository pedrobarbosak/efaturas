package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	notFoundMessage                 = "The requested entity was not found"
	badRequestMessage               = "API server could not understand the request with the data that was given."
	unauthorizedMessage             = "Unauthorized action."
	forbiddenMessage                = "Forbidden action."
	internalServerErrorMessage      = "Internal Server Error! A error occurred while working on the request."
	unprocessableEntityErrorMessage = "The request was correct, but it could not be processed due to semantic issues."
	conflictErrorMessage            = "The requested entity conflicts with existing one"
	unsupportedMediaErrorMessage    = "The requested media type is not supported"
)

func Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		err := ctx.Errors.Last()
		if err == nil {
			return
		}

		if validatorError(err) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Err.Error()})
			return
		}

		errorToMessage(ctx, err)
	}
}

func validatorError(err error) bool {
	var validationErrors validator.ValidationErrors
	ok := errors.As(err, &validationErrors)
	if ok {
		return true
	}

	var ginError gin.Error
	ok = errors.As(err, &ginError)
	return ok
}

func errorToMessage(ctx *gin.Context, err error) {
	var returnCode int
	var errorMessage string

	switch GetCode(err) {
	case FatalError:
		returnCode = http.StatusInternalServerError
		errorMessage = internalServerErrorMessage

	case InputError:
		returnCode = http.StatusBadRequest
		errorMessage = badRequestMessage

	case NotFoundError:
		returnCode = http.StatusNotFound
		errorMessage = notFoundMessage

	case ForbiddenError:
		returnCode = http.StatusForbidden
		errorMessage = forbiddenMessage

	case UnauthorizedError:
		returnCode = http.StatusUnauthorized
		errorMessage = unauthorizedMessage

	case UnprocessableEntityError:
		returnCode = http.StatusUnprocessableEntity
		errorMessage = unprocessableEntityErrorMessage

	case ConflictError:
		returnCode = http.StatusConflict
		errorMessage = conflictErrorMessage

	case UnsupportedMediaType:
		returnCode = http.StatusUnsupportedMediaType
		errorMessage = unsupportedMediaErrorMessage

	default:
		returnCode = http.StatusInternalServerError
		errorMessage = internalServerErrorMessage
	}

	ctx.AbortWithStatusJSON(returnCode, gin.H{"message": errorMessage})
}
