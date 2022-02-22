package main

type Response struct {
	Action  string      `json:"action"`
	ID      string      `json:"id"`
	Status  string      `json:"status"`
	Message interface{} `json:"message,omitempty"`
}
