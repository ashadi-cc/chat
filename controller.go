package chat

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

//UserJoined represent of users list
type UserJoined map[string]string

//Controller wrapper
type Controller struct {
	Sse     *sse.Server
	m       sync.Mutex
	newUser bool
	Users   UserJoined
}

//create new instance of controller
func newController(sse *sse.Server) *Controller {
	return &Controller{
		Sse:     sse,
		newUser: false,
		Users:   make(UserJoined),
	}
}

//Index endpoint
func (h *Controller) Index(w http.ResponseWriter, r *http.Request) {
	b, err := loadIndexFile()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

//Join /join endpoint implementation
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
	if _, ok := h.Users[user.Username]; ok {
		RespondError(w, http.StatusBadRequest, fmt.Sprintf("user %s already exist", user.Name))
		return
	}

	//lock process
	h.m.Lock()
	h.Users[user.Username] = user.Name

	RespondJSON(w, http.StatusCreated, ResponseJoin{
		Success:  true,
		Message:  fmt.Sprintf("new user joined %s", user.Name),
		Username: user.Username,
	})

	h.newUser = true
	h.m.Unlock()
}

//Send /send endpoint implementation
func (h *Controller) Send(w http.ResponseWriter, r *http.Request) {
	message := Message{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if _, ok := h.Users[message.To]; !ok {
		RespondError(w, http.StatusBadRequest, fmt.Sprintf("user not found :%s", message.To))
		return
	}

	var msg = map[string]string{"message": message.Message, "from": message.From, "to": message.To}

	sendMessage(h.Sse, fmt.Sprintf("chat/%s", message.To), msg)

	RespondJSON(w, http.StatusOK, msg)
}

// WatchChannel watch online user
// and push to list channel when someone joined / leaved
func (h *Controller) WatchChannel() {
	go func(s *sse.Server) {
		for {
			time.Sleep(1 * time.Second)
			h.m.Lock()
			if h.newUser {
				sendMessage(s, "list", h.Users)
				h.newUser = false
			} else {
				channels := s.Channels()
				h.newUser = h.CheckOnlineUser(channels)
			}
			h.m.Unlock()
		}
	}(h.Sse)
}

func (h *Controller) createTmpUser(channels []string) UserJoined {
	tmpUsers := make(UserJoined)

	for _, c := range channels {
		if strings.HasPrefix(c, "/events/ping/") {
			user := strings.TrimPrefix(c, "/events/ping/")
			tmpUsers[user] = h.Users[user]
		}
	}

	return tmpUsers
}

//CheckOnlineUser check for user leave
func (h *Controller) CheckOnlineUser(channels []string) bool {
	tmpUsers := h.createTmpUser(channels)
	isChanged := false

	for key := range h.Users {
		if _, ok := tmpUsers[key]; !ok {
			isChanged = true
		}
	}

	if isChanged {
		h.Users = tmpUsers
	}

	return isChanged
}
