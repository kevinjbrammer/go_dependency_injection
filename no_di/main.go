package main

import "spaceships/server"

func main() {
	server := server.NewServer()
	server.Run()
}
