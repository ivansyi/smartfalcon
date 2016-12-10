package http

import (
	"falcon-agent/funcs"
	"net/http"
)

func configSensorsRoutes() {
	http.HandleFunc("/proc/sensors", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, funcs.SensorsForProc())
	})
}
