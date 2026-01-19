package main

import (
    "bastille-ui/api"
    "bastille-ui/web"
    "flag"
)

func main() {

    apiOnly := flag.Bool("api-only", false, "Only run the API server")
    debug := flag.Bool("debug", false, "Enable debug logging")

    flag.Parse()

	api.InitLogger(*debug)
    go api.Start()

    if !*apiOnly {
        go web.Start()
    }

    select {}
}
