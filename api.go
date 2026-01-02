package api

import (
	"log"
	"net/http"
)

var apiKey string

func Start() {
	// apiKey = os.Getenv("BASTILLE_API_KEY")
	apiKey = "testingkey"
	if apiKey == "" {
		log.Fatal("BASTILLE_API_KEY not set")
	}
	http.HandleFunc("/api/v1/bastille/bootstrap", verifyAPIKey(BastilleBootstrapHandler))
	http.HandleFunc("/api/v1/bastille/clone", verifyAPIKey(BastilleCloneHandler))
	http.HandleFunc("/api/v1/bastille/cmd", verifyAPIKey(BastilleCmdHandler))
	http.HandleFunc("/api/v1/bastille/config", verifyAPIKey(BastilleConfigHandler))
	http.HandleFunc("/api/v1/bastille/console", verifyAPIKey(BastilleConsoleHandler))
	http.HandleFunc("/api/v1/bastille/convert", verifyAPIKey(BastilleConvertHandler))
	http.HandleFunc("/api/v1/bastille/cp", verifyAPIKey(BastilleCpHandler))
	http.HandleFunc("/api/v1/bastille/create", verifyAPIKey(BastilleCreateHandler))
	http.HandleFunc("/api/v1/bastille/destroy", verifyAPIKey(BastilleDestroyHandler))
	http.HandleFunc("/api/v1/bastille/edit", verifyAPIKey(BastilleEditHandler))
	http.HandleFunc("/api/v1/bastille/etcupdate", verifyAPIKey(BastilleEtcupdateHandler))
	http.HandleFunc("/api/v1/bastille/export", verifyAPIKey(BastilleExportHandler))
	http.HandleFunc("/api/v1/bastille/htop", verifyAPIKey(BastilleHtopHandler))
	http.HandleFunc("/api/v1/bastille/import", verifyAPIKey(BastilleImportHandler))
	http.HandleFunc("/api/v1/bastille/jcp", verifyAPIKey(BastilleJcpHandler))
	http.HandleFunc("/api/v1/bastille/limits", verifyAPIKey(BastilleLimitsHandler))
	http.HandleFunc("/api/v1/bastille/list", verifyAPIKey(BastilleListHandler))
	http.HandleFunc("/api/v1/bastille/migrate", verifyAPIKey(BastilleMigrateHandler))
	http.HandleFunc("/api/v1/bastille/monitor", verifyAPIKey(BastilleMonitorHandler))
	http.HandleFunc("/api/v1/bastille/mount", verifyAPIKey(BastilleMountHandler))
	http.HandleFunc("/api/v1/bastille/network", verifyAPIKey(BastilleNetworkHandler))
	http.HandleFunc("/api/v1/bastille/pkg", verifyAPIKey(BastillePkgHandler))
	http.HandleFunc("/api/v1/bastille/rcp", verifyAPIKey(BastilleRcpHandler))
	http.HandleFunc("/api/v1/bastille/rdr", verifyAPIKey(BastilleRdrHandler))
	http.HandleFunc("/api/v1/bastille/rename", verifyAPIKey(BastilleRenameHandler))
	http.HandleFunc("/api/v1/bastille/restart", verifyAPIKey(BastilleRestartHandler))
	http.HandleFunc("/api/v1/bastille/service", verifyAPIKey(BastilleServiceHandler))
	http.HandleFunc("/api/v1/bastille/setup", verifyAPIKey(BastilleSetupHandler))
	http.HandleFunc("/api/v1/bastille/start", verifyAPIKey(BastilleStartHandler))
	http.HandleFunc("/api/v1/bastille/stop", verifyAPIKey(BastilleStopHandler))
	http.HandleFunc("/api/v1/bastille/sysrc", verifyAPIKey(BastilleSysrcHandler))
	http.HandleFunc("/api/v1/bastille/tags", verifyAPIKey(BastilleTagsHandler))
	http.HandleFunc("/api/v1/bastille/template", verifyAPIKey(BastilleTemplateHandler))
	http.HandleFunc("/api/v1/bastille/top", verifyAPIKey(BastilleTopHandler))
	http.HandleFunc("/api/v1/bastille/umount", verifyAPIKey(BastilleUmountHandler))
	http.HandleFunc("/api/v1/bastille/update", verifyAPIKey(BastilleUpdateHandler))
	http.HandleFunc("/api/v1/bastille/upgrade", verifyAPIKey(BastilleUpgradeHandler))
	http.HandleFunc("/api/v1/bastille/verify", verifyAPIKey(BastilleVerifyHandler))
	http.HandleFunc("/api/v1/bastille/zfs", verifyAPIKey(BastilleZfsHandler))
}

func verifyAPIKey(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		authToken := "Bearer " + apiKey

		if authHeader != authToken  {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		h(w, r)
	}
}
