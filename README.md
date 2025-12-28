# bastille-api

API Interface for Bastille (https://bastillebsd.org/)


Commands
========

This API currently does the following commands:

create

destroy

start

stop

restart

rename

Setup
-----

Fist clone the repo then cd into bastille-api.  Now you uneed to initialize 
the go module.

```shell
go build
./bastille-api
```

Now run the API on your server:
```sne..
go run .
```

You should see:
```shell
✅ BastilleBSD API running on http://localhost:8080
```

Now you are ready to run requests.  Here are some sample requests:
```shell
Create a jail
-------------
curl "http://localhost:8080/bastille/create?options=-V+-M+--gateway+192.168.1.1&name=testjail&release=14.2-RELEASE&ip=192.168.0.10&iface=em0"

Start jail
------------
curl "http://localhost:8080/bastille/start?name=testjail"

Rename jail
-----------
curl "http://localhost:8080/bastille/rename?target=testjail&new_name=myjail"

Restart jail
------------
curl "http://localhost:8080/bastille/restart?name=myjail"

Stop jail
---------
curl "http://localhost:8080/bastille/stop?name=myjail"

Destroy jail
------------
curl "http://localhost:8080/bastille/destroy?name=myjail"
```


