package api

import (
	"net/http"
)

// Load valid routes and method
func loadRoutes() {

	mux := http.DefaultServeMux

	routes := map[string]http.HandlerFunc{
		"bootstrap":     BastilleBootstrapHandler,
		"clone":     BastilleCloneHandler,
		"cmd":       BastilleCmdHandler,
		"config":    BastilleConfigHandler,
		"console":   BastilleConsoleHandler,
		"convert":   BastilleConvertHandler,
		"cp":        BastilleCpHandler,
		"create":    BastilleCreateHandler,
		"destroy":   BastilleDestroyHandler,
		"edit":      BastilleEditHandler,
		"etcupdate": BastilleEtcupdateHandler,
		"export":    BastilleExportHandler,
		"htop":      BastilleHtopHandler,
		"import":    BastilleImportHandler,
		"jcp":       BastilleJcpHandler,
		"limits":    BastilleLimitsHandler,
		"list":      BastilleListHandler,
		"migrate":   BastilleMigrateHandler,
		"monitor":   BastilleMonitorHandler,
		"mount":     BastilleMountHandler,
		"network":   BastilleNetworkHandler,
		"pkg":       BastillePkgHandler,
		"rcp":       BastilleRcpHandler,
		"rdr":       BastilleRdrHandler,
		"rename":    BastilleRenameHandler,
		"restart":   BastilleRestartHandler,
		"service":   BastilleServiceHandler,
		"setup":     BastilleSetupHandler,
		"start":     BastilleStartHandler,
		"stop":      BastilleStopHandler,
		"sysrc":     BastilleSysrcHandler,
		"tags":      BastilleTagsHandler,
		"template":  BastilleTemplateHandler,
		"top":       BastilleTopHandler,
		"umount":    BastilleUmountHandler,
		"update":    BastilleUpdateHandler,
		"upgrade":   BastilleUpgradeHandler,
		"verify":    BastilleVerifyHandler,
		"zfs":       BastilleZfsHandler,
	}

	for path, handler := range routes {
		staticPath := "/api/v1/bastille/" + path
		livePath := "/api/v1/bastille/live/" + path
		mux.Handle(staticPath, loggingMiddleware(corsMiddleware(apiKeyMiddleware(handler))))
		mux.Handle(livePath, loggingMiddleware(corsMiddleware(apiKeyMiddleware(handler))))
	}
}