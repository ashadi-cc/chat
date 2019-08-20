package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexandrevicenzi/go-sse"
)

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func RespondError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{"error": message})
}

func sendMessage(s *sse.Server, path string, m interface{}) {
	channel := fmt.Sprintf("/events/%s", path)
	payload, _ := json.Marshal(m)
	s.SendMessage(channel, sse.SimpleMessage(string(payload)))
}
