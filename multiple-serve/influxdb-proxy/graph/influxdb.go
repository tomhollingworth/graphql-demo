package graph

import (
	"context"
	"fmt"
	"os"

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

	// check if the endpoint is set in the environment
	endpoint := os.Getenv("INFLUX_ENDPOINT")
	if endpoint != "" {
		url = endpoint
	}

	client = influxdb2.NewClient(url, token)
	ok, err := client.Ping(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to ping influxdb: %s", err))
	}
	if !ok {
		panic("failed to ping influxdb: not ok")
	}
}
