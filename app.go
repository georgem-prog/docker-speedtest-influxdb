package main

import (
	"github.com/kylegrantlucas/speedtest"
	"quadstingray/speedtest-influxdb/model"
	"log"
	"time"
)

func main() {
	settings := model.Parser()

	for true {
		log.Println("speed test started")
		stats := runTest(settings)
		model.SaveToInfluxDb(stats, settings)
		log.Printf("Ping: %3.2f ms | Download: %3.2f Mbps | Upload: %3.2f Mbps", stats.Ping, stats.Down_Mbs, stats.Up_Mbs)
		log.Printf("sleep for %v seconds", settings.Interval)
		time.Sleep(time.Duration(settings.Interval) * time.Second)
	}
}

func runTest(settings model.Settings) model.SpeedTestStatistics {
	client, err := speedtest.NewDefaultClient()
	if err != nil {
		log.Fatalf("error creating client: %v", err)
	}

	// Pass an empty string to select the fastest server
	server, err := client.GetServer(settings.Server)
	if err != nil {
		log.Fatalf("error getting server: %v", err)
	}

	dmbps, err := client.Download(server)
	if err != nil {
		log.Fatalf("error getting download: %v", err)
	}

	umbps, err := client.Upload(server)
	if err != nil {
		log.Fatalf("error getting upload: %v", err)
	}
	return model.SpeedTestStatistics{server.ID, server.Name + ", " + server.Country, server.Latency, dmbps, umbps, server.Distance}
}