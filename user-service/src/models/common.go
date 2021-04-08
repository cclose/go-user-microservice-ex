package models

type Message struct {
	Message string `json:"message"`
}

type ErrorMessage struct {
	Message string `json:"error"`
}
