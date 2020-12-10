package collector

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type statusCollector struct {
	ConnectionStatus     *prometheus.Desc
	CurrentNetworkTypeEx *prometheus.Desc
	WifiStatus           *prometheus.Desc
}

func newStatusCollector() hilinkCollector {

	return &statusCollector{
		/*
			2 , 3 , 5 , 8 , 20 , 21 , 23 , 27 , 28 , 29 , 30 , 31 , 32 , 33 , 65538 , 65539 , 65567 , 65568 , 131073 , 131074 , 131076 , 131078 - erreur de connexion (profil incorrect) .
			7 , 11 , 14 , 37 , 131079 , 131080 , 131081 , 131082 , 131083 , 131084 , 131085 , 131086 , 131087 , 131088 , 131089 - une erreur de connexion (accès refusé au réseau)
			12 , 13 - not connected (no roaming)
			112 - no auto-connect
			113 - no auto-connect while roaming
			114 - no reconnect
			115 - no reconnect while roaming
			201 - connexion interrupted, quota exceeded
			900 - connecting
			901 - connected
			902 - disconnected
			903 - disconnecting
			904 - connection failed
			905 - no connexion (low signal)
			906 - connexion error
		*/
		ConnectionStatus: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "status", "connection"),
			"current connection status, see exporter code for details (901 = connected)",
			nil,
			nil,
		),
		/*
			0 - aucun service
			1 - GSM
			2 - GPRS
			3 - BORD
			4 - WCDMA
			5 - HSDPA
			6 - HSUPA
			7 - HSPA
			8 - TDSCDMA
			9 - HSPA +
			10 - EVDO rév. 0
			11 - EVDO rév. ET
			12 - EVDO rév. B
			13 - 1xRTT
			14 - UMB
			15 - 1xEVDV
			16 - 3xRTT
			17 - HSPA + 64QAM
			18 - HSPA + MIMO
			19 - LTE
			21 - IS95A
			22 - IS95B
			23 - CDMA1x
			24 - EVDO rév. 0
			25 - EVDO rév. ET
			26 - EVDO rév. B
			27 - CDMA1x hybride
			28 - EVDO hybride rév. 0
			29 - EVDO hybride rév. ET
			30 - EVDO hybride rév. B
			31 - EHRPD rév. 0
			32 - EHRPD rév. ET
			33 - EHRPD rév. B
			34 - DSEH hybride rév. 0
			35 - Hybride EHRPD rév. ET
			36 - Hybrid EHRPD rév. B
			41 - WCDMA
			42 - HSDPA
			43 - HSUPA
			44 - HSPA
			45 - HSPA +
			46 - DC HSPA +
			61 - TD SCDMA
			62 - TD HSDPA
			63 - TD HSUPA
			64 - TD HSPA
			65 - TD HSPA +
			81 - 802.16E
			101 - LTE
		*/
		CurrentNetworkTypeEx: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "status", "current_network_type_ex"),
			"current network connection type, see exporter code for details (101 | 19 = LTE, 1011 = LTE+)",
			nil,
			nil,
		),
		WifiStatus: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "status", "wifi"),
			"wlan status (1 = enabled, 0 = disabled)",
			nil,
			nil,
		),
	}
}

func (c *statusCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.ConnectionStatus
	ch <- c.CurrentNetworkTypeEx
	ch <- c.WifiStatus
}

func (c *statusCollector) collect(ctx *collectorContext) error {

	log.Debug("Collecting status info")

	response, err := ctx.client.StatusInfo()
	if err != nil {
		return err
	}

	if ConnectionStatus, ok := response["ConnectionStatus"]; ok {
		if f, err := strconv.ParseFloat(fmt.Sprintf("%s", ConnectionStatus), 64); err == nil {
			ctx.ch <- prometheus.MustNewConstMetric(c.ConnectionStatus, prometheus.GaugeValue, f)
		}
	}

	if CurrentNetworkTypeEx, ok := response["CurrentNetworkTypeEx"]; ok {
		if f, err := strconv.ParseFloat(fmt.Sprintf("%s", CurrentNetworkTypeEx), 64); err == nil {
			ctx.ch <- prometheus.MustNewConstMetric(c.CurrentNetworkTypeEx, prometheus.GaugeValue, f)
		}
	}

	if WifiStatus, ok := response["WifiStatus"]; ok {
		if f, err := strconv.ParseFloat(fmt.Sprintf("%s", WifiStatus), 64); err == nil {
			ctx.ch <- prometheus.MustNewConstMetric(c.WifiStatus, prometheus.GaugeValue, f)
		}
	}

	return nil
}
