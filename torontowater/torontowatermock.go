package torontowater

import (
	"log"
	"net/http"
	"time"
)

func Mock() {
	log.Println("Mocking Toronto Water!")

	http.HandleFunc("/cc_api/svcaccount_v1/WaterAccount/validate", validate)
	http.HandleFunc("/cc_api/svcaccount_v1/WaterAccount/accountdetails", accountdetails)
	http.HandleFunc("/cc_api/svcaccount_v1/WaterAccount/consumption", consumption)

	go func() {
		log.Fatal(http.ListenAndServe(":9999", nil))
	}()
}

func validate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"validateResponse\":{\"status\":\"SUCCESS\",\"accountNumber\":\"0123456\",\"refToken\":\"abc123\"}}"))
}

func accountdetails(w http.ResponseWriter, r *http.Request) {
	body := `{
				"premiseList": [
					{
						"meterList": [
							{
								"meterNumber": "1234",
								"miu": "4321",
								"firstReadDate": "` + time.Now().Format("2006-01-02") + `",
								"lastReadDate": "` + time.Now().Format("2006-01-02") + `"
							}
						]
					}
				],
				"resultCode": 200
			}`
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func consumption(w http.ResponseWriter, r *http.Request) {
	body := `{
				"meterList": [
					{
						"meterNumber": "1234",
						"miuId": "4321",
						"intervalList": [
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 00:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 01:00:00",
								"intConsumptionTotal": 0.0,
								"intConsumptionMin": 0.0,
								"intConsumptionMax": 0.0,
								"intervalCount": 1,
								"lastReadValue": "00832.5",
								"intConsumptionAvg": 0.0,
								"synthesisFlag": "A",
								"lastReadDateTime": "` + time.Now().Format("2006-01-02") + ` 01:30:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 02:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 03:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 04:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 05:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 06:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 07:00:00",
								"intConsumptionTotal": 0.0,
								"intConsumptionMin": 0.0,
								"intConsumptionMax": 0.0,
								"intervalCount": 1,
								"lastReadValue": "00832.5",
								"intConsumptionAvg": 0.0,
								"synthesisFlag": "A",
								"lastReadDateTime": "` + time.Now().Format("2006-01-02") + ` 07:30:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 08:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 09:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 10:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 11:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 12:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 13:00:00",
								"intConsumptionTotal": 0.5,
								"intConsumptionMin": 0.5,
								"intConsumptionMax": 0.5,
								"intervalCount": 1,
								"lastReadValue": "00832.5",
								"intConsumptionAvg": 0.5,
								"synthesisFlag": "A",
								"lastReadDateTime": "` + time.Now().Format("2006-01-02") + ` 13:30:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 14:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 15:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 16:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 17:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 18:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 19:00:00",
								"intConsumptionTotal": 0.0,
								"intConsumptionMin": 0.0,
								"intConsumptionMax": 0.0,
								"intervalCount": 1,
								"lastReadValue": "00833.0",
								"intConsumptionAvg": 0.0,
								"synthesisFlag": "A",
								"lastReadDateTime": "` + time.Now().Format("2006-01-02") + ` 19:45:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 20:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 21:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 22:00:00"
							},
							{
								"intStartDate": "` + time.Now().Format("2006-01-02") + ` 23:00:00"
							}
						]
					}
				],
				"resultCode": 200
			}`
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}
