package middleware

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/legacy"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func RequestValidation() gin.HandlerFunc {
	spec, err := openapi3.NewLoader().LoadFromFile("docs/openapi.yaml")
	if err != nil {
		panic(fmt.Sprintf("failed to load openapi %v", err))
	}

	router, err := legacy.NewRouter(spec)
	if err != nil {
		panic(fmt.Sprintf("failed to load router %v", err))
	}

	return func(c *gin.Context) {
		// Find route
		route, pathParams, err := router.FindRoute(c.Request)
		if err != nil {
			if err.Error() == "no matching operation was found" {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "path not found"})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    c.Request,
			PathParams: pathParams,
			Route:      route,
			Options:    &openapi3filter.Options{MultiError: true},
		}

		err = openapi3filter.ValidateRequest(c.Request.Context(), requestValidationInput)
		if err != nil {
			validationError, ok := err.(openapi3.MultiError)
			if ok {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": strings.Split(validationError.Error(), " | ")})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}
		c.Next()
	}
}
