package main

import "log"

func main() {

	config := AppConfig{
		addr: ":8082",
	}
	app := NewApplication(config)

	mux := app.mount()

	log.Fatal(app.Run(mux))
}
