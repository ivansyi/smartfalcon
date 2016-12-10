package http

import (
	"fmt"
	"github.com/toolkits/core"
	"github.com/toolkits/nux"
	"net/http"
	"strings"
)

func configDfRoutes() {
	http.HandleFunc("/page/df", func(w http.ResponseWriter, r *http.Request) {
		mountPoints, err := nux.ListMountPoint()
		if err != nil {
			RenderMsgJson(w, err.Error())
			return
		}

		var ret [][]interface{} = make([][]interface{}, 0)
		for idx := range mountPoints {
			var du *nux.DeviceUsage
			du, err = nux.BuildDeviceUsage(mountPoints[idx][0], mountPoints[idx][1], mountPoints[idx][2])
			if err == nil {
				ret = append(ret,
					[]interface{}{
						du.FsSpec,
						core.ReadableSize(float64(du.BlocksAll)),
						core.ReadableSize(float64(du.BlocksUsed)),
						core.ReadableSize(float64(du.BlocksFree)),
						fmt.Sprintf("%.1f%%", du.BlocksUsedPercent),
						du.FsFile,
						core.ReadableSize(float64(du.InodesAll)),
						core.ReadableSize(float64(du.InodesUsed)),
						core.ReadableSize(float64(du.InodesFree)),
						fmt.Sprintf("%.1f%%", du.InodesUsedPercent),
						du.FsVfstype,
					})
			}
		}

		RenderDataJson(w, ret)
	})

	http.HandleFunc("/proc/df", func(w http.ResponseWriter, r *http.Request) {
		mountPoints, err := nux.ListMountPoint()
		if err != nil {
			RenderMsgJson(w, err.Error())
			return
		}

		var ret map[string]interface{} = make(map[string]interface{}, 0)
		for idx := range mountPoints {
			var du *nux.DeviceUsage
			du, err = nux.BuildDeviceUsage(mountPoints[idx][0], mountPoints[idx][1], mountPoints[idx][2])
			if err == nil {
				df_info := map[string]interface{}{
					"volumn":        du.FsSpec,
					"size_total":    core.ReadableSize(float64(du.BlocksAll)),
					"size_used":     core.ReadableSize(float64(du.BlocksUsed)),
					"size_free":     core.ReadableSize(float64(du.BlocksFree)),
					"size_percent":  fmt.Sprintf("%.1f", du.BlocksUsedPercent),
					"mount_point":   du.FsFile,
					"inode_total":   core.ReadableSize(float64(du.InodesAll)),
					"inode_used":    core.ReadableSize(float64(du.InodesUsed)),
					"inode_free":    core.ReadableSize(float64(du.InodesFree)),
					"inode_percent": fmt.Sprintf("%.1f", du.InodesUsedPercent),
					"fstype":        du.FsVfstype,
				}
				if strings.HasPrefix(du.FsSpec, "/dev/") {
					ret[du.FsSpec] = df_info
					//ret = append(ret, df_info)
				}

			}
		}

		RenderDataJson(w, ret)
	})
}
