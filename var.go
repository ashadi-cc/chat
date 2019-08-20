package chat

//User reprenent of User model
type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

//ResponseJoin represent of joined message
type ResponseJoin struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Username string `json:"username"`
}

//Message represent of message data
type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}
