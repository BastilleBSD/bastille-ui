package api

import (
	"net/http"
)

// Load valid routes and method
func loadRoutes() {

	mux := http.DefaultServeMux

	bastilleRoutes := map[string]http.HandlerFunc{
		"bootstrap": BastilleBootstrapHandler,
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
	rocinanteRoutes := map[string]http.HandlerFunc{
		"bootstrap": RocinanteBootstrapHandler,
		"cmd":       RocinanteCmdHandler,
		"limits":    RocinanteLimitsHandler,
		"list":      RocinanteListHandler,
		"pkg":       RocinantePkgHandler,
		"service":   RocinanteServiceHandler,
		"sysctl":    RocinanteSysctlHandler,
		"sysrc":     RocinanteSysrcHandler,
		"template":  RocinanteTemplateHandler,
		"update":    RocinanteUpdateHandler,
		"upgrade":   RocinanteUpgradeHandler,
		"verify":    RocinanteVerifyHandler,
		"zfs":       RocinanteZfsHandler,
		"zpool":     RocinanteZpoolHandler,
	}


	for path, handler := range bastilleRoutes {
		staticPathBastille := "/api/v1/bastille/" + path
		livePathBastille := "/api/v1/bastille/live/" + path

		mux.Handle(staticPathBastille , loggingMiddleware(apiKeyMiddleware(validateMethodMiddleware(handler, path, "bastille"))))
		mux.Handle(livePathBastille , loggingMiddleware(apiKeyMiddleware(validateMethodMiddleware(handler, path, "bastille"))))

	}
	for path, handler := range rocinanteRoutes {
		staticPathRocinante := "/api/v1/rocinante/" + path
		livePathRocinante := "/api/v1/rocinante/live/" + path

		mux.Handle(staticPathRocinante , loggingMiddleware(apiKeyMiddleware(validateMethodMiddleware(handler, path, "rocinante"))))
		mux.Handle(livePathRocinante , loggingMiddleware(apiKeyMiddleware(validateMethodMiddleware(handler, path, "rocinante"))))

	}

}

