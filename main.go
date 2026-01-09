package main

import (
	"bastille-ui/api"
	"bastille-ui/web"
)

func main() {
	go api.Start()
	go web.Start()
	select {}
}
