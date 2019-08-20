package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/alexandrevicenzi/go-sse"
	"github.com/gosimple/slug"
)

type UserJoined map[string]string

var (
	//User list variable
	Users   = make(UserJoined)
	newUser = false
)

type Controller struct {
	Sse *sse.Server
	m   sync.Mutex
}

func newController(sse *sse.Server) *Controller {
	return &Controller{Sse: sse}
}

func (h *Controller) Join(w http.ResponseWriter, r *http.Request) {
	user := User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	//set slug
	user.Username = slug.Make(user.Name)

	//check existing username is exist
	if _, ok := Users[user.Username]; ok {
		RespondError(w, http.StatusBadRequest, fmt.Sprintf("user %s already exist", user.Name))
		return
	}

	//lock process
	h.m.Lock()
	Users[user.Username] = user.Name

	RespondJSON(w, http.StatusCreated, ResponseJoin{
		Success:  true,
		Message:  fmt.Sprintf("new user joined %s", user.Name),
		Username: user.Username,
	})
	newUser = true
	h.m.Unlock()
}

func (h *Controller) Send(w http.ResponseWriter, r *http.Request) {
	message := Message{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if _, ok := Users[message.To]; !ok {
		RespondError(w, http.StatusBadRequest, fmt.Sprintf("user not found :%s", message.To))
		return
	}

	var msg = map[string]string{"message": message.Message, "from": message.From, "to": message.To}

	sendMessage(h.Sse, fmt.Sprintf("chat/%s", message.To), msg)

	RespondJSON(w, http.StatusOK, msg)
}

//WatchChannel watch online user
// and push to list channel when someone joined / leaved
func (h *Controller) WatchChannel() {
	go func(s *sse.Server) {
		for {
			time.Sleep(1 * time.Second)
			h.m.Lock()
			if newUser {
				sendMessage(s, "list", Users)
				newUser = false
			} else {
				channels := s.Channels()
				newUser = CheckOnlineUser(channels)
			}
			h.m.Unlock()
		}
	}(h.Sse)
}

func createTmpUser(channels []string) UserJoined {
	tmpUsers := make(UserJoined)

	for _, c := range channels {
		if strings.HasPrefix(c, "/events/ping/") {
			user := strings.TrimPrefix(c, "/events/ping/")
			tmpUsers[user] = Users[user]
		}
	}

	return tmpUsers
}

//CheckOnlineUser check for user leave
func CheckOnlineUser(channels []string) bool {
	tmpUsers := createTmpUser(channels)
	isChanged := false

	for key := range Users {
		if _, ok := tmpUsers[key]; !ok {
			isChanged = true
		}
	}

	if isChanged {
		Users = tmpUsers
	}

	return isChanged
}
