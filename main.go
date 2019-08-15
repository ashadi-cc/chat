package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alexandrevicenzi/go-sse"
	"github.com/gorilla/mux"
)

const STATIC_DIR = "/static/"

type App struct {
	Router *mux.Router
	SSe    *sse.Server
}

func (a *App) createSSE() {
	a.SSe = sse.NewServer(&sse.Options{
		// Increase default retry interval to 10s.
		RetryInterval: 10 * 1000,
		// CORS headers
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, OPTIONS",
			"Access-Control-Allow-Headers": "Keep-Alive,X-Requested-With,Cache-Control,Content-Type,Last-Event-ID",
		},
		// Custom channel name generator
		ChannelNameFunc: func(request *http.Request) string {
			return request.URL.Path
		},
		// Print debug info
		Logger: log.New(os.Stdout, "go-sse: ", log.Ldate|log.Ltime|log.Lshortfile),
	})
}

func (a *App) Init() {
	a.createSSE()
	a.Router = mux.NewRouter()
	a.Router.PathPrefix("/events/").Handler(a.SSe)

	hh := Controller{Sse: a.SSe}
	hh.WatchChannel()
	a.Router.HandleFunc("/join", hh.Join).Methods("POST")
	a.Router.HandleFunc("/send", hh.Send).Methods("POST")

	a.Router.
		PathPrefix(STATIC_DIR).
		Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))
}

func (a *App) Run() {
	a.Init()
	//defer a.SSe.Shutdown()
	log.Println("Server listening on port :3000")
	log.Fatal(http.ListenAndServe(":3000", a.Router))
}

func main() {
	var app = &App{}
	app.Run()
}
