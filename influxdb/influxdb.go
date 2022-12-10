package influxdb

import (
	"container/list"
	"context"
	"log"
	"strconv"
	"time"

	"github.com/dtrumpfheller/toronto-water-exporter/helpers"
	"github.com/dtrumpfheller/toronto-water-exporter/torontowater"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func Export(meterNumber string, consumptions *list.List, config helpers.Config) {

	// create client objects
	client := influxdb2.NewClient(config.InfluxDB.URL, config.InfluxDB.Token)
	queryAPI := client.QueryAPI(config.InfluxDB.Organization)
	writeAPI := client.WriteAPI(config.InfluxDB.Organization, config.InfluxDB.Bucket)

	// start & end can be determined based on list elements
	startDateTime := consumptions.Front().Value.(torontowater.WaterConsumption).Time.Add(-1 * time.Hour)
	endDateTime := consumptions.Back().Value.(torontowater.WaterConsumption).Time.Add(1 * time.Hour)

	// check if entry is already stored, only consider last 7 days
	query := `from(bucket: "` + config.InfluxDB.Bucket + `")
		|> range(start: ` + strconv.FormatInt(startDateTime.Unix(), 10) + `, stop: ` + strconv.FormatInt(endDateTime.Unix(), 10) + `)
		|> filter(fn: (r) => r["_measurement"] == "toronto_water")
		|> filter(fn: (r) => r["meter"] == "` + meterNumber + `")`
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error calling InfluxDB [%s]!\n", err.Error())
		return
	}

	// remove consumptions that have already been submitted
	for result.Next() {
		var next *list.Element
		for e := consumptions.Front(); e != nil; e = next {
			next = e.Next()
			if e.Value.(torontowater.WaterConsumption).Time.Equal(result.Record().Time()) {
				consumptions.Remove(e)
				break
			}
		}
	}

	if consumptions.Len() > 0 {
		// write remaining consumptions to influxdb
		for e := consumptions.Front(); e != nil; e = e.Next() {
			consumption := e.Value.(torontowater.WaterConsumption)

			log.Println("Inserting " + consumption.Time.Format("2006-01-02 15:04:05"))
			point := influxdb2.NewPointWithMeasurement("toronto_water").
				AddTag("meter", meterNumber).
				AddField("m3", consumption.Value).
				SetTime(consumption.Time)
			writeAPI.WritePoint(point)
		}

		// force all unwritten data to be sent
		writeAPI.Flush()

	} else {
		log.Println("No new metrics available, skip export to influx")
	}

	// ensures background processes finishes
	client.Close()
}
