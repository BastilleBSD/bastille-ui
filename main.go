package main

import (
    "bastille-ui/web"
    "flag"
)

func main() {

	webDir := flag.String("webdir", "", "Web files location")

	flag.Parse()

	web.Start(*webDir)
}
