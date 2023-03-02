package utils

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func ErrorResponse(message string, err string, data interface{}) Response {
	splittedError := strings.Split(err, "\n")
	errResponse := Response{
		Status:  false,
		Message: message,
		Errors:  splittedError,
		Data:    data,
	}
	return errResponse
}
func SuccessResponse(message string, data interface{}) Response {
	SuccResponse := Response{
		Status:  true,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
	return SuccResponse
}
func ResponseJSON(c *gin.Context, response Response) {
	c.Writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(c.Writer).Encode(response)
}
