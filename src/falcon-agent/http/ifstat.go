package http

import (
	"falcon-agent/funcs"
	"net/http"
)

func configIfStatRoutes() {
	http.HandleFunc("/page/ifstat", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, funcs.IFStatsForPage())
	})

	http.HandleFunc("/proc/ifstat", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, funcs.IFStatsForProc())
	})
}
