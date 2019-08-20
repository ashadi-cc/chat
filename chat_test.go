package chat

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var app *App

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	return rr
}

func TestMain(m *testing.M) {
	app = &App{}
	app.Init()
	code := m.Run()
	os.Exit(code)
}

func TestJoinChat(t *testing.T) {
	user := User{Name: "Ashadi cahyadi"}
	payload, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/join", bytes.NewBuffer(payload))

	if err != nil {
		t.Fatalf("got error when hit /join endpoint : %s", err.Error())
	}

	response := executeRequest(req)
	if response.Code != http.StatusCreated {
		t.Fatalf("Cannot join new user ashadi got response code %d", response.Code)
	}

	rjoin := ResponseJoin{}

	json.Unmarshal(response.Body.Bytes(), &rjoin)

	if !strings.Contains(rjoin.Message, "Ashadi cahyadi") {
		t.Fatalf("response does not match: %s", rjoin.Message)
	}
}

func TestDuplicateUsername(t *testing.T) {
	user := User{Name: "Ashadi cahyadi"}
	payload, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/join", bytes.NewBuffer(payload))
	response := executeRequest(req)
	if response.Code == http.StatusCreated {
		t.Fatal("its should not create new user as the username is duplicate")
	}
}

func TestCreateSecondUser(t *testing.T) {
	user := User{Name: "Markonah"}
	payload, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/join", bytes.NewBuffer(payload))
	response := executeRequest(req)

	if response.Code != http.StatusCreated {
		t.Fatal("its should create new user")
	}
}
func TestSendMessage(t *testing.T) {
	message := Message{From: "ashadi-cahyadi", To: "markonah", Message: "Dodol"}
	payload, _ := json.Marshal(message)
	req, _ := http.NewRequest("POST", "/send", bytes.NewBuffer(payload))
	response := executeRequest(req)

	if response.Code != http.StatusOK {
		t.Fatalf("got error when sending message: %s", string(response.Body.Bytes()))
	}
}
