package chat

import (
	"log"
	"net/http"
	"os"

	"github.com/alexandrevicenzi/go-sse"
	"github.com/gorilla/mux"
)

//STATIC_DIR static directory
const staticDir = "/static/"

//App struct
type App struct {
	Router    *mux.Router
	SSe       *sse.Server
	Controler *Controller
}

//create SSE server
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

//set the controller
func (a *App) setController() {
	a.Controler = newController(a.SSe)
}

//set router
func (a *App) setRouter() {
	a.Router = mux.NewRouter()
	a.Router.PathPrefix("/events/").Handler(a.SSe)
	a.Controler.WatchChannel()
	a.Router.HandleFunc("/join", a.Controler.Join).Methods("POST")
	a.Router.HandleFunc("/send", a.Controler.Send).Methods("POST")
	a.Router.
		PathPrefix(staticDir).
		Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))
}

//Init set initial process of app
func (a *App) Init() {
	a.createSSE()
	a.setController()
	a.setRouter()
}

//Run run http webserver
func (a *App) Run() {
	log.Println("Server listening on port :3000")
	log.Fatal(http.ListenAndServe(":3000", a.Router))
}

//NewApp create App instance
func NewApp() *App {
	return &App{}
}
