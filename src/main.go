package main

import "github.com/alex-1900/wishlist/src/app"

func main() {
	app := app.NewApp()
	app.PrepareRoutes()
	app.Listen()
}
