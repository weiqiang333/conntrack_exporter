package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/weiqiang333/conntrack_exporter/pkg/conntrack"
	"github.com/weiqiang333/go-web-init-template/pkg/utils/net_tools"
)

func init() {
	pflag.String("address", ":9110", "The address on which to expose the web interface and generated Prometheus metrics.")
	pflag.String("configfile", "./config/conntrack_exporter.yaml", "exporter config file")
}

const namespace = "conntrack"

type Exporter struct {
	ConnectInfo prometheus.GaugeVec // 连接信息
}

func NewExporter() *Exporter {
	return &Exporter{
		ConnectInfo: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "connect_info",
				Help:      "network connect info, src and dst and port info",
			}, []string{"src_address", "src_port", "dst_address", "dst_port", "tcp_state_code", "tcp_state"}),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.ConnectInfo.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.ConnectInfo.Reset()

	ct := conntrack.NewConntrackInfos()
	ct.GetConntrack()
	ctInfos := ct.ConntrackInfos
	localhostIps, err := net_tools.GetLocalhostIps()
	if err != nil {
		log.Println("Wran GetLocalhostIps error:", err.Error())
	}
	log.Println(localhostIps)
	for _, i := range ctInfos {
		e.ConnectInfo.WithLabelValues(i.SrcAddress, i.SrcPort, i.DstAddress, i.DstPort, string(i.TcpStateCode), i.TcpState).Set(1)
	}

	e.ConnectInfo.Collect(ch)
}

func reloadConfig(w http.ResponseWriter, _ *http.Request) {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Println(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	fmt.Println(fmt.Sprintf("reload config file: %s", viper.ConfigFileUsed()))
	io.WriteString(w, fmt.Sprintf("rereload config file: %s", viper.ConfigFileUsed()))
}

func main() {
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal("Fatal error BindPFlags: %w", err.Error())
	}
	fmt.Println("load config file ", viper.GetString("configfile"))
	viper.SetConfigType("yaml")
	viper.SetConfigFile(viper.GetString("configfile"))
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	prometheus.MustRegister(NewExporter())

	// http server
	listenAddress := viper.GetString("address")
	fmt.Printf("http server start, address %s/metrics\n", listenAddress)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/reload", reloadConfig)
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		log.Fatal("Fatal error http: %w", err)
	}
}
