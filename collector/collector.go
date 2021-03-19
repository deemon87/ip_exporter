package collector

import (
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/deemon87/ip_exporter/info"
	"github.com/prometheus/client_golang/prometheus"
)

//Collector type
type Collector struct {
	NetInfo *prometheus.Desc
}

// NewCollector func
func NewCollector() *Collector {
	return &Collector{
		NetInfo: prometheus.NewDesc("iface",
			"Network adapter config",
			[]string{"name", "mac", "ipv4", "ipv6", "isLoopback", "isPrivate"}, nil,
		),
	}
}

// Describe method
func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.NetInfo
}

// Collect method
func (collector *Collector) Collect(ch chan<- prometheus.Metric) {
	networks, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, net := range info.GetInterfacesInfo(networks) {
		ipv4 := strings.Join(net.Ipv4, ",")
		if len(ipv4) == 0 {
			ipv4 = "none"
		}

		ipv6 := strings.Join(net.Ipv6, ",")
		if len(ipv6) == 0 {
			ipv6 = "none"
		}

		mac := net.Mac.String()
		if len(mac) == 0 {
			mac = "none"
		}

		ch <- prometheus.MustNewConstMetric(collector.NetInfo, prometheus.GaugeValue, 1,
			net.Name, mac, ipv4, ipv6, strconv.FormatBool(net.IsLoopback), strconv.FormatBool(net.IsPrivate))
	}
}
