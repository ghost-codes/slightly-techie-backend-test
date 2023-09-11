package main

import (
	"fmt"
	"log"

	api "ghost-codes/slightly-techie-blog/api"
	config "ghost-codes/slightly-techie-blog/config"
	"ghost-codes/slightly-techie-blog/docs"

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

	docs.SwaggerInfo.Title = "Slightly Techie - Blog Test"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = config.ServerUrl
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{config.SCHEME}

	fmt.Println(docs.SwaggerInfo.Host)

	server, err := api.NewServer(config)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	err = server.Start(config.ServerAddr)
	if err != nil {
		log.Fatal("error occured, Server could not start; Error:", err)
	}
}
