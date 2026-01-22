# BastilleBSD API + UI

## Commands

The API handles all bastille commands, and is consistent with
the syntax of the CLI. Any parameter passed via the CLI is named
the same in the API, with some exceptions. Any command that supports
both a jail or a release, will only accept a `target` parameter. See
the `destroy` example below.

Ir also handles all Rocinante commands, except for the `CP` hook.

## Setup

You can either build and run, or just run the package.

```shell
go build
./bastille-ui
```
or
```
go run .
```
To run only the API: `go run . --api-only`

To run in debug mode: `go run . --debug`

API config file: `api/config.json`

WebUI config file: `web/config.json`

Requests made via the API must contain an `Authorization: Bearer API_KEY` header. The
`API_KEY` can be set inside the `api/config.json` file. The `API_KEY` is set every time
the program starts.

For the WebUI, the `web/config.json` file contains a default username and password to
log in. Simply visit http://host:port to get started.

To use the console on the homepage, you need to `pkg install ttyd`.

## Dependencies

```
bastille
rocinante (optional)
go
ttyd
```

## API Usage

All requests called via GET will return the supported parameters and options. To actually
run the command, it must be a POST request.

Bastille endpoint: `/api/v1/bastille/command`

Rocinante endpoint: `/api/v1/rocinante/command`

Any parameter/option string that has spaces shiukd be passed with either
a `+` or `%20` as the space character. See examples below...

## API Examples

Create a jail
```
curl "http://ip:port/api/v1/bastille/create?name=test&release=15.0-release&ip=10.0.0.12&iface=vtnet0" -H "Authorization: Bearer API_KEY"
```

Create a vnet jail with custom gateway and nameserver
```
curl "http://ip:port/api/v1/bastille/create?name=test&release=15.0-release&ip=10.0.0.12&iface=vtnet0&options=-V+-g+192.168.10.1+-n+192.168.10.1" -H "Authorization: Bearer API_KEY"
```

Destroy a jail
```
curl "http://ip:port/api/v1/bastille/destroy?target=test&options=-f+-a+-y" -H "Authorization: Bearer API_KEY"
```

Run a command inside a jail
```
curl "http://ip:port/api/v1/bastille/cmd?target=test&command=echo+hi+how%20are%20you" -H "Authorization: Bearer API_KEY"
```

Apply a templateusing rocinante
```
curl "http://ip:port/api/v1/rocinante/template?template=custom/template" -H "Authorization: Bearer API_KEY"
```
