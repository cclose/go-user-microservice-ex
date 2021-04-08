// Package controllers provides the route handlers with a convenient struct for accessing dependencies
package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/cclose/go-user-microservice-ex/user-service/src/database"
	"github.com/cclose/go-user-microservice-ex/user-service/src/models"
	"github.com/cclose/go-user-microservice-ex/user-service/src/service"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type UserControllerV1 struct {
	Service *service.UserService
}

// errorResponse Handles returning a JSON encoded error message
func errorResponse(writer http.ResponseWriter, errorCode int, errorMessage string) {
	jsonResponse(writer, errorCode, models.ErrorMessage{Message: errorMessage})
}

// jsonResponse Handlers the boilerplate of encoding the payload to JSON and setting the proper headers
func jsonResponse(writer http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		errorResponse(writer, http.StatusInternalServerError, err.Error())
	}

	writer.WriteHeader(statusCode)
	_, _ = writer.Write(response)
}

// validateRequest Handles making sure the request is a json request
func validateRequest(writer http.ResponseWriter, request *http.Request) bool {
	contentType := request.Header.Get("Content-Type")
	if contentType != "application/json" {
		errorResponse(writer, http.StatusUnsupportedMediaType,
			fmt.Sprintf("Illegal Request Content-Type. Only accepcts application/json. Received: %s", contentType))
		return false
	}

	return true
}

// CreateUser creates a User
func (c *UserControllerV1) CreateUser(writer http.ResponseWriter, request *http.Request) {
	if !validateRequest(writer, request) { // Make sure we only accept JSON
		return
	}

	var user models.UserModel
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&user); err != nil {
		errorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	if err := user.Create(c.Service.Dbh); err != nil {
		if err == database.ErrDuplicateKey {
			errorResponse(writer, http.StatusConflict, "Request violates uniqueness")
		} else {
			errorResponse(writer, http.StatusBadRequest, err.Error())
		}

	} else {
		//blank the password so we don't return it
		user.Password = ""

		//return the newly created user back to the requester
		jsonResponse(writer, http.StatusCreated, user)
	}
}

// DeleteUser removes the specified user id
func (c *UserControllerV1) DeleteUser(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		errorResponse(writer, http.StatusBadRequest, "Invalid ID "+vars["id"])
		return
	}

	user := models.UserModel{ID: id}
	err = user.Delete(c.Service.Dbh)
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse(writer, http.StatusNotFound, fmt.Sprintf("No User with ID %d found", id))
		} else {
			errorResponse(writer, http.StatusInternalServerError, err.Error())
		}
	} else {
		jsonResponse(writer, http.StatusOK, models.Message{Message: fmt.Sprintf("User ID %d deleted", id)})
	}
}

const (
	DEFAULT_LIMIT  = "100"
	DEFAULT_OFFSET = "0"
)

// GetAllUsers gets all users
func (c *UserControllerV1) GetAllUsers(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	limitVal := query.Get("limit")
	if limitVal == "" {
		limitVal = DEFAULT_LIMIT
	}
	limit, err := strconv.Atoi(limitVal)
	if err != nil {
		errorResponse(writer, http.StatusBadRequest,
			fmt.Sprintf("query \"limit\" only accepts integers: received %s", limitVal))
		return
	}

	offsetVal := query.Get("offset")
	if offsetVal == "" {
		offsetVal = DEFAULT_OFFSET
	}
	var offset int
	offset, err = strconv.Atoi(offsetVal)
	if err != nil {
		errorResponse(writer, http.StatusBadRequest,
			fmt.Sprintf("query \"offset\" only accepts integers: received %s", offsetVal))
		return
	}

	var users []models.UserModel
	users, err = models.GetUsers(c.Service.Dbh, "all", "", limit, offset)
	if err != nil {
		errorResponse(writer, http.StatusInternalServerError, err.Error())
	} else {
		// if we don't have any, make an empty slice so it serializes as "[]" instead of null
		if len(users) == 0 {
			users = make([]models.UserModel, 0)
		}
		jsonResponse(writer, http.StatusOK, users)
	}
}

// GetUserById searches for a User by the specified ID
func (c *UserControllerV1) GetUserById(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	idVal := vars["id"]
	// We simply do this to validate it's an integer
	_, err := strconv.Atoi(idVal)
	if err != nil {
		errorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Invalid ID %s", idVal))
		return
	}

	var users []models.UserModel
	users, err = models.GetUsers(c.Service.Dbh, "id", idVal, 1, 0)
	if err != nil {
		errorResponse(writer, http.StatusInternalServerError, err.Error())
	} else {
		// return the first element of the slice so it isn't serialized as an array
		jsonResponse(writer, http.StatusOK, users[0])
	}
}

// GetUserById updates the user by the specified ID
func (c *UserControllerV1) UpdateUser(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		errorResponse(writer, http.StatusBadRequest, "Invalid ID "+vars["id"])
		return
	}

	var user models.UserModel
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&user)
	if err != nil {
		errorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}

	if id != user.ID {
		errorResponse(writer, http.StatusBadRequest, "changing ID is not permitted")
		return
	}

	err = user.Update(c.Service.Dbh)
	if err == sql.ErrNoRows {
		errorResponse(writer, http.StatusNotFound, fmt.Sprintf("no user with id %d found", user.ID))
	} else if err != nil {
		errorResponse(writer, http.StatusInternalServerError, err.Error())
	} else {
		//blank the password so we don't return it
		user.Password = ""

		jsonResponse(writer, http.StatusOK, user)
	}
}

// AuthenticateUser using http basic auth, tests for valid credentials
func (c *UserControllerV1) AuthenticateUser(writer http.ResponseWriter, request *http.Request) {
	if c.checkAuthentication(request) {
		jsonResponse(writer, http.StatusOK, models.Message{Message: "Success"})
	} else {
		errorResponse(writer, http.StatusUnauthorized, "Invalid credentials")
	}
}

// checkAuthentication Given a request, checks for basic auth and validates credentials
// separated so it can be used to provide authentication for other routers
func (c *UserControllerV1) checkAuthentication(request *http.Request) bool {
	username, password, success := request.BasicAuth()
	if success {
		hashedPassword, err := models.GetUserCredentials(c.Service.Dbh, username)
		if err == nil && password != "" {
			if models.CheckPassword(hashedPassword, password) {
				return true
			}
		}
	}

	return false
}
