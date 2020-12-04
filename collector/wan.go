package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type wanCollector struct {
	desc  *prometheus.Desc
	props []string
}

func newWanCollector() hilinkCollector {

	props := []string{"WanIPAddress", "WanIPv6Address"}

	return &wanCollector{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "wan"),
			"wan info",
			props,
			nil,
		),
		props: props,
	}
}

func (c *wanCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *wanCollector) collect(ctx *collectorContext) error {

	log.Debug("Collecting wan info")

	response, err := ctx.client.DeviceInfo()
	if err != nil {
		return err
	}

	labels := getValuesFromResponse(c.props, response)
	ctx.ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, 1, labels...)

	return nil
}
