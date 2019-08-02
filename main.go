package main

import "github.com/teohrt/cruddyAPI/app"

func main() {
	app.Start(app.Config{
		Port: "8080",
	})
}
