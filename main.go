package main

import (
	"log"

	api "ghost-codes/slightly-techie-blog/api"
	config "ghost-codes/slightly-techie-blog/config"

	_ "github.com/lib/pq"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("unable to load config:", err)
	}
	// conn, err := sql.Open(config.DBDriver, config.DBSource())

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	server, err := api.NewServer(config)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	err = server.Start(config.ServerAddr)
	if err != nil {
		log.Fatal("error occured, Server could not start; Error:", err)
	}
}
