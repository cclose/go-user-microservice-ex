// Package service provides a central struct for containing the dependencies of the service
package service

import (
	"fmt"
	"github.com/cclose/go-user-microservice-ex/user-service/src/database"
	"github.com/cclose/go-user-microservice-ex/user-service/src/models"
	"github.com/gorilla/mux"
	"log"
	"os"
)

// USerService manages the dependencies and subservices for the User Service
type UserService struct {
	Dbh         *database.PostGresDB
	Router      *mux.Router
	Logger      log.Logger
	ServicePort string
}

func (s *UserService) Initialize() {
	//Get our listening port
	s.ServicePort = os.Getenv("PORT")
	if s.ServicePort == "" {
		s.ServicePort = "8080"
	}

	// Get our DB Connection settings
	host := os.Getenv("PG_HOST")
	user := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASSWORD")
	db := os.Getenv("PG_DB")
	port := os.Getenv("PG_PORT")
	//TODO validate values for these

	// Boot the DB connection
	Dbh, err := database.Connect(host, user, pass, db, port)
	if err != nil {
		fmt.Println("[status] [fatal] Unable to connect to DB: ", err)
		os.Exit(1)
	}
	s.Dbh = Dbh

	//Verify the User Table
	_, err = s.Dbh.PgDbSession.Exec(models.UserSchema)
	if err != nil {
		fmt.Println("[status] [fatal] Unable to create UserModel Table: ", err)
		os.Exit(1)
	}

	s.Router = mux.NewRouter()
}
