# bastille-api

API Interface for Bastille (https://bastillebsd.org/)


Commands
========

The API handles all bastille commands, and is consistent with
the syntax of the CLI. Any parameter passed via the CLI is named
the same in the API.

Setup
-----

Fist clone the repo then cd into bastille-ui.  Now you uneed to initialize 
the go module.

```shell
go build
./bastille-ui
```

You should see:
```shell
BastilleBSD UI started on http://localhost:8080
```

Now you are ready to run requests.  Here are some sample requests:
```shell
Create a jail
-------------
curl "http://localhost:8080/api/v1/bastille/create?options=-V+-M+--gateway+192.168.1.1&name=testjail&release=14.2-RELEASE&ip=192.168.0.10&iface=em0"

Start jail
------------
curl "http://localhost:8080/api/v1/bastille/start?name=testjail"

Rename jail
-----------
curl "http://localhost:8080/api/v1/bastille/rename?target=testjail&new_name=myjail"

Restart jail
------------
curl "http://localhost:8080/api/v1/bastille/restart?name=myjail"

Stop jail
---------
curl "http://localhost:8080/api/v1/bastille/stop?name=myjail"

Destroy jail
------------
curl "http://localhost:8080/api/v1/bastille/destroy?name=myjail"
```
For the WebUI, visit http://localhost:8080 and play with it.

