# BastilleBSD API + UI

Commands
========

The API handles all bastille commands, and is consistent with
the syntax of the CLI. Any parameter passed via the CLI is named
the same in the API.

Setup
-----

You can either build and run, or just run the package.

```shell
go build
./bastille-ui
```
or
```
go run .
```

Request made via the API must contain an `Authorization: Bearer API_KEY` header. The
`API_KEY` can be set inside the `config.json` file. The `API_KEY` is set every time
the program starts, and can also be changed via the webui.

The `config.json` file also contains a default username and password to log in via the
webui. Simply visit http://host:port to get started.

API Examples
============

Create a jail
```
curl "http://ip:port/api/v1/bastille/create?name=test&release=15.0-release&ip=10.0.0.12&iface=vtnet0" -H "Authorization: Bearer API_KEY"
```

Destroy a jail
```
curl "http://ip:port/api/v1/bastille/destroy?target=test" -H "Authorization: Bearer API_KEY"
```

