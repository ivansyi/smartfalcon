package http

import (
	"falcon-agent/g"
	"github.com/toolkits/nux"
	"github.com/toolkits/sys"
	"net/http"
)

func configKernelRoutes() {
	http.HandleFunc("/proc/kernel/hostname", func(w http.ResponseWriter, r *http.Request) {
		data, err := g.Hostname()
		AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/kernel/maxproc", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.KernelMaxProc()
		AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/kernel/maxfiles", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.KernelMaxFiles()
		AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/kernel/file-max", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.KernelMaxFiles()
		AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/kernel/file-nr", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.KernelAllocateFiles()
		AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/kernel/procs", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.AllProcs()
		AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/kernel/version", func(w http.ResponseWriter, r *http.Request) {
		data, err := sys.CmdOutNoLn("uname", "-r")
		AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/socket/tcp", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.TcpPorts()
		AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/socket/udp", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.UdpPorts()
		AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/socket/sum", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.SocketStatSummary()
		AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/netstat/tcp", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.Netstat("TcpExt")
		AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/netstat/ip", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.Netstat("IpExt")
		AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/meminfo", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.MemInfo()
		AutoRender(w, data, err)
	})

}
