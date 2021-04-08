package main

import (
	"fmt"
	"github.com/cclose/go-user-microservice-ex/user-service/src/controllers"
	"github.com/cclose/go-user-microservice-ex/user-service/src/service"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Booting User Service...")

	userService := service.UserService{}
	userService.Initialize()
	defer userService.Dbh.Disconnect()
	userService.Router.HandleFunc("/test", ServeTest)

	// api v1 router
	v1 := userService.Router.PathPrefix("/api/v1").Subrouter()
	// user v1 controller
	uc := controllers.UserControllerV1{Service: &userService}
	v1.HandleFunc("/user", uc.GetAllUsers).Methods(http.MethodGet)
	v1.HandleFunc("/user", uc.CreateUser).Methods(http.MethodPost)
	v1.HandleFunc("/user/{id:[0-9]+}", uc.GetUserById).Methods(http.MethodGet)
	v1.HandleFunc("/user/{id:[0-9]+}", uc.DeleteUser).Methods(http.MethodDelete)
	v1.HandleFunc("/user/{id:[0-9]+}", uc.UpdateUser).Methods(http.MethodPut)
	v1.HandleFunc("/user/auth", uc.AuthenticateUser).Methods(http.MethodPost)

	http.Handle("/", userService.Router)
	err := http.ListenAndServe(fmt.Sprintf(":%s", userService.ServicePort), userService.Router)
	if err != nil {
		fmt.Println("[status] [fatal] Caught Error on HTTP ListenAndServer: ", err)
		os.Exit(1)
	}
	userService.Logger.Println("[status] [online] User Service online and ready to serve")
}

func ServeTest(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "./static/test.html")
}
