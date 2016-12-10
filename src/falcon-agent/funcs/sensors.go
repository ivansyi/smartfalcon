package funcs

import (
	"bytes"
	"falcon-common/model"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func get_cpu_temp() (i float64, err error) {
	cmd := exec.Command("sensors", "coretemp-isa-0000")
	var out bytes.Buffer
	var line string
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		line, err = out.ReadString('\n')
		if err != nil {
			return
		}
		if strings.HasPrefix(line, "Physical id 0:") {
			last := strings.Index(line, "°C")
			rs := []rune(line)
			temp := strings.TrimSpace(string(rs[15:last]))
			//log.Println(temp, line)
			i, err = strconv.ParseFloat(string(temp), 64)
			if err != nil {
				//err = fmt.Errorf("failed to convert to int64")
				return
			}
			return
		}
	}
	err = fmt.Errorf("not found")
	return
}

func SensorsMetrics() (L []*model.MetricValue) {
	temp, _ := get_cpu_temp()
	cpu := GaugeValue("cpu.temp", temp)
	return []*model.MetricValue{cpu}
}

func SensorsForProc() interface{} {
	cpu, _ := get_cpu_temp()

	items := map[string]interface{}{
		"cpu": cpu,
	}
	return items
}
