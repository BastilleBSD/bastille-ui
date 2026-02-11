# BastilleBSD UI Server

## Setup

Run the following commands to get started.

```shell
make install
cp /usr/local/etc/bastille-ui/config.json.sample /usr/local/etc/bastille-ui/config.json
sysrc bastille_ui_enable=YES
service bastille-ui start
```

Customize the config file to your liking.

The config file contains a default username (admin) and password (admin) to
log in. Simply visit http://host:port to get started.

To use the console on the homepage, you need to `pkg install ttyd` on the API server side.

## Dependencies

The API server is required for the UI to be able
to execute commands.

[BastilleBSD API Server](https://github.com/BastilleBSD/bastille-api)

```
bastille
rocinante (optional)
go
ttyd (optional)
```
