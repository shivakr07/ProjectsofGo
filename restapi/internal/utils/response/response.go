package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	// Status string
	// Error  string
	// in response i was getting Error : "empty body" but E is capital as defined in the struct but which seems in consistent with value as in the go strings at these places suggest to use lower case
	//so how this Error or any other should look like in the response we can change that too using struct tags
	//we can define how this should look like when it gets deserialize into json
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

// it will take the response object of func(w.http.ResponseWriter, r *http.Request)
// since we don't know what kind of data we are going to receive so we take generic type  data any or data interface {}
func WriteJson(w http.ResponseWriter, status int, data interface{}) error {

	w.Header().Set("Content-Type", "application/json")
	// this allow us to send json data[we need to add header]
	w.WriteHeader(status)

	// so as we decoded the incoming json to pass in struct, in the same we way we need to again convert the struct data to json before returning [encode]
	return json.NewEncoder(w).Encode(data)
	//it also accepts interface io.Writer and this response writer implements that interface so we can pass that
	// encode method returns error in case of error so we have kept return type
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	//errs is a slice
	var errMessages []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMessages = append(errMessages, fmt.Sprintf("field %s is required field", err.Field()))

		default:
			errMessages = append(errMessages, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMessages, ", "),
		//it will join all the string elements of that slice as we have in python
	}
}
