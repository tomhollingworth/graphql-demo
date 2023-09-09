package graph

import (
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var (
	bucket = "history"
	org    = "demo"
	token  = "influxdb"
	// Store the URL of your InfluxDB instance
	url = "http://influxdb:8086"
)
var client influxdb2.Client

func init() {
	client = influxdb2.NewClient(url, token)
	ok, err := client.Ping(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to ping influxdb: %s", err))
	}
	if !ok {
		panic("failed to ping influxdb: not ok")
	}
}
