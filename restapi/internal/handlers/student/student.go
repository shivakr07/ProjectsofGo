package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/shivakr07/students-api/internal/storage"
	"github.com/shivakr07/students-api/internal/types"
	"github.com/shivakr07/students-api/internal/utils/response"
)

// this is convention you can use Create also
// this func will return the handler function as you saw in the GET
// router.HandleFunc("GET /",
// func(w http.ResponseWriter, r *http.Request) { .. this func
// and at that place we just need to give reference of this func

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			//we can directly return the response using write but we want to return json response
			//we make one more package response in the utils

			// response.WriteJson(w, http.StatusBadRequest, err.Error())
			// instead of the err.Error we will call that function
			// response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))

			// or you can be more specific about the error till now we were getting EOF in the error but we can be specific
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
			//returning : which make sure no further execution after this
		}
		// this errors package matches the error which we pass so here EOF means end of file [means we have got the empty object / no data from the request body]
		//this NewDecoder accepts interface of type io.Reader and this request we are getting implements that so we can pass that

		//what if we get some other except EOF we need to catch that too
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//VALIDATE THE REQUEST [don't believe on client][0 trust policy]
		//REQUEST VALIDATION
		// we can do it manually but we will use package [validator] : golang request valiation playground
		if err := validator.New().Struct(student); err != nil {
			//since ValidationError function accepts a different type than err so we need to typecaste that
			validateErrors := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrors))
			return
		}

		//create student
		//since we are receiving it as dependency then we can use that in this way
		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("user created successfully", slog.String("userId", fmt.Sprint(lastId)))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		//we need to serialize the json data we will get from request, so that we can use that

		// response.WriteJson(w, http.StatusCreated, map[string]string{"sucess": "OK"})

		//since now we are assuming everything is ok so return proper values
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})

		//NOW WE are ready to test as our handler is ready
		// we got {id:1} in response when we sent the data

		// so till now our api is working we are getting
		// EOF error [status code 400 bad request]
		// {sucess:ok} [status code 200 ok]

		//but the error is not properly descriptive we need improvement like here also we need json

		//we make one more function in the response.go
	}
}

//any dependency for the New function will be defined here as definition is not separated to make the clean everything
//we will inject the dependency here -> DEPENDENCY INJECTION
