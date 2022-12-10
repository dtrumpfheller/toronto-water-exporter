package main

import (
	"container/list"
	"flag"
	"log"
	"time"

	"github.com/dtrumpfheller/toronto-water-exporter/helpers"
	"github.com/dtrumpfheller/toronto-water-exporter/influxdb"
	"github.com/dtrumpfheller/toronto-water-exporter/torontowater"
)

var (
	configFile = flag.String("config", "config.yml", "configuration file")
	config     helpers.Config
)

func main() {

	// load arguments into variables
	flag.Parse()

	// load config file
	config = helpers.ReadConfig(*configFile)

	// setup mock if necessary
	if config.TorontoWater.Mock {
		torontowater.Mock()
	}

	for {
		// export metrics
		exportMetrics()

		if config.SleepDuration <= 0 {
			break
		}
		time.Sleep(time.Duration(config.SleepDuration) * time.Minute)
	}
}

func exportMetrics() {
	log.Println("Getting Toronto Water consumption... ")
	start := time.Now()

	token, err := torontowater.Login(config)
	if err != nil {
		return
	}

	meters, err := torontowater.GetAccountDetails(token, config)
	if err != nil {
		return
	}

	for _, meter := range meters {
		endDate, _ := time.ParseInLocation("2006-01-02", meter.LastReadDate, start.Location())
		startDate, _ := time.ParseInLocation("2006-01-02", meter.FirstReadDate, start.Location())

		date := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location()).AddDate(0, 0, -config.LookDaysInPast)
		if date.Before(startDate) {
			date = startDate
		}

		// 1. get data
		consumptions := list.New()

		for ok := (endDate.After(date) || endDate.Equal(date)); ok; ok = (endDate.After(date) || endDate.Equal(date)) {
			data, err := torontowater.GetData(meter, date, token, config)
			if err == nil {
				consumptions.PushBackList(data)
			}
			date = date.AddDate(0, 0, 1)
		}

		// 2. export data
		if consumptions.Len() > 0 {
			influxdb.Export(meter.MeterNumber, consumptions, config)
		} else {
			log.Println("No data gathered, skipping export to influxDB")
		}
	}

	log.Printf("Finished in %s\n", time.Since(start))
}
