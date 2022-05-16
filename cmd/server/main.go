package main

import "github.com/MustCo/Mon_go/internal/server/apiserver"

func main() {
	config := apiserver.Config{bind_addr: ":8080"}
	server := apiserver.APIServer.New(config)
	server.Start()

}
