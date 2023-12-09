package main

import "github.com/colinjlacy/mocking-http-requests/app"

func main() {
	if err := app.App.Run(); err != nil {
		panic(err)
	}
}
