.PHONY: all
all:
	@echo "Nothing to be done. Please use make install or make uninstall"
.PHONY: install
install:
	@echo
	@echo "Installing BastilleBSD API Server"
	@echo
	@go build -o /usr/local/bin/bastille-api main.go
	@echo
	@mkdir -p /usr/local/etc/bastille-api
	@cp -Rv api/config.json /usr/local/etc/bastille-api/config.json
	@echo
	@mkdir -p /usr/local/share/bastille-api
	@cp -Rv web /usr/local/share/bastille-api/
	@echo
	@mkdir -p /usr/local/etc/rc.d
	@cp -Rv etc/rc.d/* /usr/local/etc/rc.d/
	@echo
	@echo "This method is for testing & development."
	@echo
	@echo "Please report any issues to https://github.com/BastilleBSD/bastille-ui/issues"

.PHONY: uninstall
uninstall:
	@echo
	@echo "Removing BastilleBSD API Server"
	@echo
	@rm -vf /usr/local/bin/bastille-api
	@echo
	@rm -rvf /usr/local/share/bastille-api
	@echo
	@echo "removing configuration file"
	@rm -rvf /usr/local/etc/bastille-api/config.json.sample
	@echo
	@echo "removing startup script"
	@rm -vf /usr/local/etc/rc.d/bastille-api
	@echo
	@echo "You may need to manually remove /usr/local/etc/bastille-api/config.json if it is no longer needed."