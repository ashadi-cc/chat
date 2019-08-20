package main

import "chat"

func main() {
	app := chat.NewApp()
	app.Init()
	app.Run()
}
