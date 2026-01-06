package api

import (
	"log"
	"net/http"
)

// Start
func Start(addr string) {
	loadRoutes()
    log.Println("Starting BastilleBSD API server on", addr)
    log.Fatal(http.ListenAndServe(addr, nil))
}

// Handle valid routes to the API
func loadRoutes() {
	routes := map[string]http.HandlerFunc{
		"/api/v1/bastille/bootstrap": BastilleBootstrapHandler,
		"/api/v1/bastille/clone":     BastilleCloneHandler,
		"/api/v1/bastille/cmd":       BastilleCmdHandler,
		"/api/v1/bastille/config":    BastilleConfigHandler,
		"/api/v1/bastille/console":   BastilleConsoleHandler,
		"/api/v1/bastille/convert":   BastilleConvertHandler,
		"/api/v1/bastille/cp":        BastilleCpHandler,
		"/api/v1/bastille/create":    BastilleCreateHandler,
		"/api/v1/bastille/destroy":   BastilleDestroyHandler,
		"/api/v1/bastille/edit":      BastilleEditHandler,
		"/api/v1/bastille/etcupdate": BastilleEtcupdateHandler,
		"/api/v1/bastille/export":    BastilleExportHandler,
		"/api/v1/bastille/htop":      BastilleHtopHandler,
		"/api/v1/bastille/import":    BastilleImportHandler,
		"/api/v1/bastille/jcp":       BastilleJcpHandler,
		"/api/v1/bastille/limits":    BastilleLimitsHandler,
		"/api/v1/bastille/list":      BastilleListHandler,
		"/api/v1/bastille/migrate":   BastilleMigrateHandler,
		"/api/v1/bastille/monitor":   BastilleMonitorHandler,
		"/api/v1/bastille/mount":     BastilleMountHandler,
		"/api/v1/bastille/network":   BastilleNetworkHandler,
		"/api/v1/bastille/pkg":       BastillePkgHandler,
		"/api/v1/bastille/rcp":       BastilleRcpHandler,
		"/api/v1/bastille/rdr":       BastilleRdrHandler,
		"/api/v1/bastille/rename":    BastilleRenameHandler,
		"/api/v1/bastille/restart":   BastilleRestartHandler,
		"/api/v1/bastille/service":   BastilleServiceHandler,
		"/api/v1/bastille/setup":     BastilleSetupHandler,
		"/api/v1/bastille/start":     BastilleStartHandler,
		"/api/v1/bastille/stop":      BastilleStopHandler,
		"/api/v1/bastille/sysrc":     BastilleSysrcHandler,
		"/api/v1/bastille/tags":      BastilleTagsHandler,
		"/api/v1/bastille/template":  BastilleTemplateHandler,
		"/api/v1/bastille/top":       BastilleTopHandler,
		"/api/v1/bastille/umount":    BastilleUmountHandler,
		"/api/v1/bastille/update":    BastilleUpdateHandler,
		"/api/v1/bastille/upgrade":   BastilleUpgradeHandler,
		"/api/v1/bastille/verify":    BastilleVerifyHandler,
		"/api/v1/bastille/zfs":       BastilleZfsHandler,
	}

	// Log first, then auth, then actual
	for path, handler := range routes {
		cmd := loggingMiddleware(apiKeyMiddleware(http.HandlerFunc(handler)))
		http.Handle(path, cmd)
	}
}
