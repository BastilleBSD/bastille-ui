.PHONY: all
all:
	@echo "Nothing to be done. Please use make install or make uninstall"
.PHONY: install
install:
	@echo
	@echo "Installing BastilleBSD UI..."
	@go build -o /usr/local/bin/bastille-ui main.go
	@mkdir -p /usr/local/etc/rc.d
	@cp -Rv etc/rc.d/* /usr/local/etc/rc.d/
	@echo
	@echo "This method is for testing & development."
	@echo
	@echo "Please report any issues to https://github.com/BastilleBSD/bastille-ui/issues"
	@echo

.PHONY: uninstall
uninstall:
	@echo
	@echo "Removing BastilleBSD UI..."
	@rm -vf /usr/local/bin/bastille-ui
	@echo
	@echo "Removing startup script..."
	@rm -vf /usr/local/etc/rc.d/bastille-ui
	@echo
