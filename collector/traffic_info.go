package collector

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type trafficInfoCollector struct {
	CurrentDownloadRate  *prometheus.Desc
	CurrentUploadRate    *prometheus.Desc
	TotalUpload          *prometheus.Desc
	TotalDownload        *prometheus.Desc
	CurrentUpload        *prometheus.Desc
	CurrentDownload      *prometheus.Desc
	CurrentMonthDownload *prometheus.Desc
	CurrentMonthUpload   *prometheus.Desc
}

func newTrafficInfoCollector() hilinkCollector {

	return &trafficInfoCollector{
		CurrentDownloadRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "traffic", "current_download_rate"),
			"CurrentDownloadRate (bits/s)",
			nil,
			nil,
		),
		CurrentUploadRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "traffic", "current_upload_rate"),
			"CurrentUploadRate (bits/s)",
			nil,
			nil,
		),
		TotalUpload: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "traffic", "total_upload"),
			"TotalUpload (bits)",
			nil,
			nil,
		),
		TotalDownload: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "traffic", "total_download"),
			"TotalDownload (bits)",
			nil,
			nil,
		),
		CurrentUpload: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "traffic", "current_upload"),
			"CurrentUpload (bits)",
			nil,
			nil,
		),
		CurrentDownload: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "traffic", "current_download"),
			"CurrentDownload (bits)",
			nil,
			nil,
		),
		CurrentMonthDownload: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "traffic", "current_month_download"),
			"CurrentMonthDownload (bits)",
			[]string{"month_last_clear_time"},
			nil,
		),
		CurrentMonthUpload: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "traffic", "current_month_upload"),
			"CurrentMonthUpload (bits)",
			[]string{"month_last_clear_time"},
			nil,
		),
	}
}

func (c *trafficInfoCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.CurrentDownloadRate
	ch <- c.CurrentUploadRate
	ch <- c.TotalUpload
	ch <- c.TotalDownload
	ch <- c.CurrentUpload
	ch <- c.CurrentDownload
	ch <- c.CurrentMonthDownload
	ch <- c.CurrentMonthUpload
}

func (c *trafficInfoCollector) collect(ctx *collectorContext) error {

	log.Debug("Collecting traffic info")

	response, err := ctx.client.TrafficInfo()
	if err != nil {
		return err
	}

	if CurrentDownloadRate, ok := response["CurrentDownloadRate"]; ok {
		if f, err := strconv.ParseFloat(fmt.Sprintf("%s", CurrentDownloadRate), 64); err == nil {
			ctx.ch <- prometheus.MustNewConstMetric(c.CurrentDownloadRate, prometheus.CounterValue, f)
		}
	}

	if CurrentUploadRate, ok := response["CurrentUploadRate"]; ok {
		if f, err := strconv.ParseFloat(fmt.Sprintf("%s", CurrentUploadRate), 64); err == nil {
			ctx.ch <- prometheus.MustNewConstMetric(c.CurrentUploadRate, prometheus.CounterValue, f)
		}
	}

	if TotalUpload, ok := response["TotalUpload"]; ok {
		if f, err := strconv.ParseFloat(fmt.Sprintf("%s", TotalUpload), 64); err == nil {
			ctx.ch <- prometheus.MustNewConstMetric(c.TotalUpload, prometheus.CounterValue, f)
		}
	}

	if TotalDownload, ok := response["TotalDownload"]; ok {
		if f, err := strconv.ParseFloat(fmt.Sprintf("%s", TotalDownload), 64); err == nil {
			ctx.ch <- prometheus.MustNewConstMetric(c.TotalDownload, prometheus.CounterValue, f)
		}
	}

	if CurrentUpload, ok := response["CurrentUpload"]; ok {
		if f, err := strconv.ParseFloat(fmt.Sprintf("%s", CurrentUpload), 64); err == nil {
			ctx.ch <- prometheus.MustNewConstMetric(c.CurrentUpload, prometheus.CounterValue, f)
		}
	}

	if CurrentDownload, ok := response["CurrentDownload"]; ok {
		if f, err := strconv.ParseFloat(fmt.Sprintf("%s", CurrentDownload), 64); err == nil {
			ctx.ch <- prometheus.MustNewConstMetric(c.CurrentDownload, prometheus.CounterValue, f)
		}
	}

	response, err = ctx.client.MonthInfo()
	if err != nil {
		return err
	}

	MonthLastClearTime := ""
	if v, ok := response["MonthLastClearTime"]; ok {
		MonthLastClearTime = fmt.Sprintf("%s", v)
	}

	if CurrentMonthDownload, ok := response["CurrentMonthDownload"]; ok {
		if f, err := strconv.ParseFloat(fmt.Sprintf("%s", CurrentMonthDownload), 64); err == nil {
			ctx.ch <- prometheus.MustNewConstMetric(c.CurrentMonthDownload, prometheus.CounterValue, f, MonthLastClearTime)
		}
	}

	if CurrentMonthUpload, ok := response["CurrentMonthUpload"]; ok {
		if f, err := strconv.ParseFloat(fmt.Sprintf("%s", CurrentMonthUpload), 64); err == nil {
			ctx.ch <- prometheus.MustNewConstMetric(c.CurrentMonthUpload, prometheus.CounterValue, f, MonthLastClearTime)
		}
	}

	return nil
}
