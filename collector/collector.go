package collector

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/knq/hilink"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	namespace = "hilink"
	timeout   = 5 * time.Second
)

type Options struct {
	Url      *string
	Username *string
	Password *string

	EnableSignal  bool
	EnableTraffic bool
	EnableWan     bool
}

type collectorContext struct {
	ch     chan<- prometheus.Metric
	client *hilink.Client
}

type hilinkCollector interface {
	describe(ch chan<- *prometheus.Desc)
	collect(ctx *collectorContext) error
}

type collector struct {
	options    Options
	collectors []hilinkCollector
	timeout    time.Duration
	client     *hilink.Client
}

// Describe implements the prometheus.Collector interface.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	for _, co := range c.collectors {
		co.describe(ch)
	}
}

// Collect implements the prometheus.Collector interface.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	if c.client == nil {
		sigChannel := make(chan os.Signal)
		signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)
		clientOptions := []hilink.Option{
			hilink.URL(*c.options.Url),
		}
		if *c.options.Password != "" {
			clientOptions = append(clientOptions, hilink.Auth(*c.options.Username, *c.options.Password))
		}
		log.WithFields(log.Fields{
			"url":  *c.options.Url,
			"user": *c.options.Username,
		}).Debug("Connecting")
		client, err := hilink.NewClient(clientOptions...)
		if err != nil {
			log.Fatal("Unable to contact device: %s", err)
			return
		}
		go func() {
			<-sigChannel
			log.Info("Stopping, disconnecting from device")
			c.client.Disconnect()
			os.Exit(0)
		}()
		c.client = client
		defer c.client.Disconnect()
	}

	wg := sync.WaitGroup{}
	for _, collector := range c.collectors {
		col := collector
		wg.Add(1)
		go func() {
			ctx := &collectorContext{ch, c.client}
			err := col.collect(ctx)
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()

}

func NewCollector(opts Options) prometheus.Collector {

	collectors := []hilinkCollector{
		newDeviceCollector(),
	}

	if opts.EnableSignal {
		collectors = append(collectors, newSignalCollector())
	}

	if opts.EnableTraffic {
		collectors = append(collectors, newTrafficCollector())
	}

	if opts.EnableWan {
		collectors = append(collectors, newWanCollector())
	}

	c := collector{
		options:    opts,
		collectors: collectors,
		timeout:    timeout,
	}

	return &c
}

func getValuesFromResponse(labelKeys []string, response hilink.XMLData) []string {
	labels := make([]string, 0, len(labelKeys))
	for _, prop := range labelKeys {
		v, ok := response[prop]
		if !ok {
			v = ""
			log.WithFields(log.Fields{
				"prop": prop,
			}).Error("Requested prop not found in response")
		}
		labels = append(labels, fmt.Sprintf("%s", v))
	}
	return labels
}

var dbRegex = regexp.MustCompile(`(.*?)dBm?$`)

func parseDbValue(str string) (float64, error) {
	val := dbRegex.FindStringSubmatch(str)
	if len(val) != 2 {
		return 0, fmt.Errorf("unable to match decibel string value")
	}

	return strconv.ParseFloat(val[1], 64)
}
