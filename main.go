package main

import (
    "bastille-ui/api"
    "bastille-ui/web"
    "flag"
)

func main() {

    apiOnly := flag.Bool("api-only", false, "Run only the API server")
    flag.Parse()

    go api.Start()

    if !*apiOnly {
        go web.Start()
    }

    select {}
}