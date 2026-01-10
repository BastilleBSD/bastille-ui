package main

import (
    "bastille-ui/api"
    "bastille-ui/web"
    "flag"
)

func main() {

<<<<<<< HEAD
    apiOnly := flag.Bool("api-only", false, "Run only the API server")
=======
    apiOnly := flag.Bool("api-only", false, "Only run the API server")
>>>>>>> b856431accf58e3d941c801e4d8303923f73121d
    flag.Parse()

    go api.Start()

    if !*apiOnly {
        go web.Start()
    }

    select {}
}