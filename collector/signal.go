package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type signalCollector struct {
	RSSI *prometheus.Desc
	RSRP *prometheus.Desc
	RSRQ *prometheus.Desc
	SINR *prometheus.Desc
}

func newSignalCollector() hilinkCollector {

	return &signalCollector{
		RSSI: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "rssi"),
			"Received Signal Strength Indicator (dBm)",
			nil,
			nil,
		),
		RSRP: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "rsrp"),
			"Reference Signal Receive Power (dBm)",
			nil,
			nil,
		),
		RSRQ: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "rsrq"),
			"Reference Signal Receive Quality (dB)",
			nil,
			nil,
		),
		SINR: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "signal", "sinr"),
			"Signal Noise Ratio (dB, max 30)",
			nil,
			nil,
		),
	}
}

func (c *signalCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.RSSI
	ch <- c.RSRP
	ch <- c.RSRQ
	ch <- c.SINR
}

func (c *signalCollector) collect(ctx *collectorContext) error {

	log.Debug("Collecting signal info")

	response, err := ctx.client.SignalInfo()
	if err != nil {
		return err
	}

	if rssi, ok := response["rssi"]; ok {
		if f, err := parseDbValue(fmt.Sprintf("%s", rssi)); err == nil {
			ctx.ch <- prometheus.MustNewConstMetric(c.RSSI, prometheus.GaugeValue, f)
		}
	}

	if rsrp, ok := response["rsrp"]; ok {
		if f, err := parseDbValue(fmt.Sprintf("%s", rsrp)); err == nil {
			ctx.ch <- prometheus.MustNewConstMetric(c.RSRP, prometheus.GaugeValue, f)
		}
	}

	if rsrq, ok := response["rsrq"]; ok {
		if f, err := parseDbValue(fmt.Sprintf("%s", rsrq)); err == nil {
			ctx.ch <- prometheus.MustNewConstMetric(c.RSRQ, prometheus.GaugeValue, f)
		}
	}

	if sinr, ok := response["sinr"]; ok {
		if f, err := parseDbValue(fmt.Sprintf("%s", sinr)); err == nil {
			ctx.ch <- prometheus.MustNewConstMetric(c.SINR, prometheus.GaugeValue, f)
		}
	}

	return nil
}
