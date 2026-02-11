# BastilleBSD UI Server

## Setup

Run the following commands to get started.

```shell
make install
sysrc bastille_ui_enable=YES
service bastille-ui start
```

The default username is warden and the password is bastille.

To use the console on the homepage, you need to `pkg install ttyd` on the API server side.

## Dependencies

```
bastille
rocinante (optional)
go
ttyd
```
