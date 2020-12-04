package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type deviceInfoCollector struct {
	desc  *prometheus.Desc
	props []string
}

func newDeviceInfoCollector() hilinkCollector {

	props := []string{"DeviceName", "HardwareVersion", "SoftwareVersion", "MacAddress1", "MacAddress2", "Imei"}

	return &deviceInfoCollector{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "device_info"),
			"device info",
			props,
			nil,
		),
		props: props,
	}
}

func (c *deviceInfoCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *deviceInfoCollector) collect(ctx *collectorContext) error {

	log.Debug("Collecting device info")

	response, err := ctx.client.DeviceInfo()
	if err != nil {
		return err
	}

	labels := getValuesFromResponse(c.props, response)
	ctx.ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, 1, labels...)

	return nil
}
