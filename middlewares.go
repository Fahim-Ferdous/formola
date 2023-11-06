package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrorInternalError = errors.New("Woops! Something went wrong :(")

// This method collects all errors and submits them to Rollbar
func Errors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// Only run if there are some errors to handle
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				fmt.Println(e)
				// Find out what type of error it is
				switch e.Type {
				case gin.ErrorTypePublic:
					// Only output public errors if nothing has been written yet
					if !c.Writer.Written() {
						c.JSON(c.Writer.Status(), gin.H{"Error": e.Error()})
					}
				case gin.ErrorTypeBind:
					// errs := e.Err.(validator.ValidationErrors)
					// list := make(map[string]string)
					// for _, err := range errs {
					// 	list[err.Field()] = ValidationErrorToText(err)
					// }

					// Make sure we maintain the preset response status
					status := http.StatusBadRequest
					if c.Writer.Status() != http.StatusOK {
						status = c.Writer.Status()
					}
					c.JSON(status, e)

				default:
					// Log all other errors
					fmt.Println(e)
				}
			}
			// If there was no public or bind error, display default 500 message
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": ErrorInternalError.Error()})
			}
		}
	}
}
