package middleware

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"net/http"

	"github.com/kardianos/service"
	"github.com/stvp/rollbar"
	// "gopkg.in/bluesuncorp/validator.v5"
)

var (
	ErrorInternalError = errors.New("Woops! Something went wrong :(")
)

func ValidationErrorToText(e *validator.FieldError) string {
	if e == nil {
		return ""
	}
	err := *e
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", err.Field(), err.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s", err.Field(), err.Param())
	case "email":
		return fmt.Sprintf("Invalid email format")
	case "len":
		return fmt.Sprintf("%s must be %s characters long", err.Field(), err.Param())
	}
	return fmt.Sprintf("%s is not valid", err.Field())
}

// This method collects all errors and submits them to Rollbar
func Errors(env, token string, logger service.Logger) gin.HandlerFunc {
	// rollbar.Environment = env
	// rollbar.Token = token

	return func(c *gin.Context) {
		c.Next()
		// Only run if there are some errors to handle
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// Find out what type of error it is
				switch e.Type {
				case gin.ErrorTypePublic:
					// Only output public errors if nothing has been written yet
					if !c.Writer.Written() {
						c.JSON(c.Writer.Status(), gin.H{"Error": e.Error()})
					}
				case gin.ErrorTypeBind:
					var jsErr *json.UnmarshalTypeError
					if errors.As(e.Err, &jsErr) {
						status := http.StatusBadRequest
						if c.Writer.Status() != http.StatusOK {
							status = c.Writer.Status()
						}
						c.JSON(status, gin.H{"Errors": e})
					} else {
						errs := e.Err.(validator.ValidationErrors)

						list := make([]string, len(errs))

						for field, err := range errs {
							list[field] = ValidationErrorToText(&err)
						}

						// Make sure we maintain the preset response status
						status := http.StatusBadRequest
						if c.Writer.Status() != http.StatusOK {
							status = c.Writer.Status()
						}
						fmt.Println(list)
						c.JSON(status, gin.H{"Errors": list})
					}

				default:
					// Log all other errors
					rollbar.RequestError(rollbar.ERR, c.Request, e.Err)
					if logger != nil {
						logger.Error(e.Err)
					}
				}

			}
			// If there was no public or bind error, display default 500 message
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": ErrorInternalError.Error()})
			}
		}
	}
}

func BindErrorsPars(err error) gin.H {
	errs := err.(validator.ValidationErrors)

	list := make([]string, len(errs))

	for field, err := range errs {
		list[field] = ValidationErrorToText(&err)
	}

	return gin.H{"Errors": list}
}
