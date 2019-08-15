package main

//User model
type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

type ResponseJoin struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Username string `json:"username"`
}

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}
