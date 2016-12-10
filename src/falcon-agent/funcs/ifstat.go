package funcs

import (
	"falcon-agent/g"
	"falcon-common/model"
	"fmt"
	"github.com/toolkits/nux"
	"log"
	"sync"
	"time"
)

var (
	ifStatsMap = make(map[string][2]*nux.NetIf)
	ifLock     = new(sync.RWMutex)
)

func UpdateIfStats() error {
	ifList, err := nux.NetIfs(g.Config().Collector.IfacePrefix)
	if err != nil {
		return err
	}

	ifLock.Lock()
	defer ifLock.Unlock()
	for i := 0; i < len(ifList); i++ {
		device := ifList[i].Iface
		ifStatsMap[device] = [2]*nux.NetIf{ifList[i], ifStatsMap[device][0]}
	}
	return nil
}

func IFInPackages(arr [2]*nux.NetIf) int64 {
	return arr[0].InPackages - arr[1].InPackages
}

func IFInBytes(arr [2]*nux.NetIf) int64 {
	return arr[0].InBytes - arr[1].InBytes
}

func IFOutPackages(arr [2]*nux.NetIf) int64 {
	return arr[0].OutPackages - arr[1].OutPackages
}

func IFOutBytes(arr [2]*nux.NetIf) int64 {
	return arr[0].OutBytes - arr[1].OutBytes
}

func IFDelta(device string, f func([2]*nux.NetIf) int64) int64 {
	val, ok := ifStatsMap[device]
	if !ok {
		return 0
	}

	if val[1] == nil {
		return 0
	}
	return f(val)
}

func NetMetrics() []*model.MetricValue {
	return CoreNetMetrics(g.Config().Collector.IfacePrefix)
}

func CoreNetMetrics(ifacePrefix []string) []*model.MetricValue {

	netIfs, err := nux.NetIfs(ifacePrefix)
	if err != nil {
		log.Println(err)
		return []*model.MetricValue{}
	}

	cnt := len(netIfs)
	ret := make([]*model.MetricValue, cnt*23)

	for idx, netIf := range netIfs {
		iface := "iface=" + netIf.Iface
		ret[idx*23+0] = CounterValue("net.if.in.bytes", netIf.InBytes, iface)
		ret[idx*23+1] = CounterValue("net.if.in.packets", netIf.InPackages, iface)
		ret[idx*23+2] = CounterValue("net.if.in.errors", netIf.InErrors, iface)
		ret[idx*23+3] = CounterValue("net.if.in.dropped", netIf.InDropped, iface)
		ret[idx*23+4] = CounterValue("net.if.in.fifo.errs", netIf.InFifoErrs, iface)
		ret[idx*23+5] = CounterValue("net.if.in.frame.errs", netIf.InFrameErrs, iface)
		ret[idx*23+6] = CounterValue("net.if.in.compressed", netIf.InCompressed, iface)
		ret[idx*23+7] = CounterValue("net.if.in.multicast", netIf.InMulticast, iface)
		ret[idx*23+8] = CounterValue("net.if.out.bytes", netIf.OutBytes, iface)
		ret[idx*23+9] = CounterValue("net.if.out.packets", netIf.OutPackages, iface)
		ret[idx*23+10] = CounterValue("net.if.out.errors", netIf.OutErrors, iface)
		ret[idx*23+11] = CounterValue("net.if.out.dropped", netIf.OutDropped, iface)
		ret[idx*23+12] = CounterValue("net.if.out.fifo.errs", netIf.OutFifoErrs, iface)
		ret[idx*23+13] = CounterValue("net.if.out.collisions", netIf.OutCollisions, iface)
		ret[idx*23+14] = CounterValue("net.if.out.carrier.errs", netIf.OutCarrierErrs, iface)
		ret[idx*23+15] = CounterValue("net.if.out.compressed", netIf.OutCompressed, iface)
		ret[idx*23+16] = CounterValue("net.if.total.bytes", netIf.TotalBytes, iface)
		ret[idx*23+17] = CounterValue("net.if.total.packets", netIf.TotalPackages, iface)
		ret[idx*23+18] = CounterValue("net.if.total.errors", netIf.TotalErrors, iface)
		ret[idx*23+19] = CounterValue("net.if.total.dropped", netIf.TotalDropped, iface)
		ret[idx*23+20] = GaugeValue("net.if.speed.bits", netIf.SpeedBits, iface)
		ret[idx*23+21] = CounterValue("net.if.in.percent", netIf.InPercent, iface)
		ret[idx*23+22] = CounterValue("net.if.out.percent", netIf.OutPercent, iface)
	}
	return ret
}

func IFStatsForPage() (L [][]string) {
	ifLock.RLock()
	defer ifLock.RUnlock()

	for device, _ := range ifStatsMap {
		in_bytes := IFDelta(device, IFInBytes)
		out_bytes := IFDelta(device, IFOutBytes)
		in_bits := in_bytes * 8 / int64(g.COLLECT_INTERVAL/time.Second)
		out_bits := out_bytes * 8 / int64(g.COLLECT_INTERVAL/time.Second)
		bandwidth := ifStatsMap[device][0].SpeedBits
		in_rate := float64(in_bits) / float64(bandwidth) * 100.0
		out_rate := float64(out_bits) / float64(bandwidth) * 100.0
		item := []string{
			device,
			fmt.Sprintf("%d", IFDelta(device, IFInPackages)),
			fmt.Sprintf("%d", IFDelta(device, IFOutPackages)),
			fmt.Sprintf("%d", in_bits),
			fmt.Sprintf("%d", out_bits),
			fmt.Sprintf("%.2f", in_rate),
			fmt.Sprintf("%.2f", out_rate),
			fmt.Sprintf("%d", bandwidth),
			fmt.Sprintf("%d", int64(g.COLLECT_INTERVAL/time.Second)),
		}
		L = append(L, item)
	}
	return
}

func IFStatsForProc() interface{} {
	ifLock.RLock()
	defer ifLock.RUnlock()

	var L map[string]interface{} = make(map[string]interface{}, 0)

	for device, _ := range ifStatsMap {
		in_bytes := IFDelta(device, IFInBytes)
		out_bytes := IFDelta(device, IFOutBytes)
		in_bits := in_bytes * 8 / int64(g.COLLECT_INTERVAL/time.Second)
		out_bits := out_bytes * 8 / int64(g.COLLECT_INTERVAL/time.Second)
		bandwidth := ifStatsMap[device][0].SpeedBits
		in_rate := float64(in_bits) / float64(bandwidth) * 100.0
		out_rate := float64(out_bits) / float64(bandwidth) * 100.0
		item := map[string]interface{}{
			"device":        device,
			"in_packages":   fmt.Sprintf("%d", IFDelta(device, IFInPackages)),
			"out_packages":  fmt.Sprintf("%d", IFDelta(device, IFOutPackages)),
			"in_bits_rate":  fmt.Sprintf("%d", in_bits),
			"out_bits_rate": fmt.Sprintf("%d", out_bits),
			"in_percent":    fmt.Sprintf("%.2f", in_rate),
			"out_percent":   fmt.Sprintf("%.2f", out_rate),
			"bandwidth":     fmt.Sprintf("%d", bandwidth),
		}
		//L = append(L, item)
		L[device] = item
	}
	return L
}
