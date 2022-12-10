package torontowater

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"

	// "github.com/PuerkitoBio/goquery"
	"github.com/dtrumpfheller/toronto-water-exporter/helpers"
)

type Validate struct {
	ValidateResponse ValidateResponse `json:"validateResponse"`
}

type ValidateResponse struct {
	Status        string `json:"status"`
	AccountNumber string `json:"accountNumber"`
	RefToken      string `json:"refToken"`
}

type Accountdetails struct {
	ResultCode  int       `json:"resultCode"`
	PremiseList []Premise `json:"premiseList"`
}

type Premise struct {
	MeterList []Meter `json:"meterList"`
}

type Consumption struct {
	ResultCode int     `json:"resultCode"`
	MeterList  []Meter `json:"meterList"`
}

type Meter struct {
	MeterNumber   string     `json:"meterNumber"`
	Miu           string     `json:"miu"`
	FirstReadDate string     `json:"firstReadDate"`
	LastReadDate  string     `json:"lastReadDate"`
	IntervalList  []Interval `json:"intervalList"`
}

type Interval struct {
	LastReadDateTime    string  `json:"lastReadDateTime"`
	LastReadValue       string  `json:"lastReadValue"`
	IntConsumptionTotal float64 `json:"intConsumptionTotal"`
}

type WaterConsumption struct {
	Time  time.Time
	Value float64
}

var client http.Client

func Login(config helpers.Config) (string, error) {

	log.Println("Logging into Toronto Water... ")

	// create cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Printf("Got error while creating cookie jar [%s]!", err.Error())
		return "", err
	}
	client = http.Client{
		Jar: jar,
	}

	// login
	url := "https://secure.toronto.ca/cc_api/svcaccount_v1/WaterAccount/validate"
	if config.TorontoWater.Mock {
		url = "http://localhost:9999/cc_api/svcaccount_v1/WaterAccount/validate"
	}

	body := "json={\"API_OP\":\"VALIDATE\",\"ACCOUNT_NUMBER\":\"" + config.TorontoWater.AccountNumber +
		"\",\"LAST_NAME\":\"" + config.TorontoWater.LastName +
		"\",\"POSTAL_CODE\":\"" + config.TorontoWater.PostalCode +
		"\",\"LAST_PAYMENT_METHOD\":\"" + config.TorontoWater.LastPaymentMethod + "\"}"
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		log.Printf("Got error %s", err.Error())
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error loggin in Toronto Water [%s]!\n", err.Error())
		return "", err
	}
	if resp.StatusCode != 200 {
		log.Printf("Logging into Toronto Water failed with status code [%d]!\n", resp.StatusCode)
		return "", errors.New("Error")
	}

	// extract token
	defer resp.Body.Close()
	var validate Validate
	err = json.NewDecoder(resp.Body).Decode(&validate)
	if err != nil {
		log.Printf("Error processing Toronto Water login response [%s]!\n", err.Error())
		return "", err
	}

	if validate.ValidateResponse.Status != "SUCCESS" {
		log.Println("Login failed!")
		return "", errors.New("Error")
	}

	return validate.ValidateResponse.RefToken, nil
}

func GetAccountDetails(token string, config helpers.Config) ([]Meter, error) {

	log.Println("Getting account details")

	// get data
	jsonData := "{\"API_OP\":\"ACCOUNTDETAILS\",\"ACCOUNT_NUMBER\":\"" + config.TorontoWater.AccountNumber + "\"}"
	urlToCall := "https://secure.toronto.ca/cc_api/svcaccount_v1/WaterAccount/accountdetails"
	if config.TorontoWater.Mock {
		urlToCall = "http://localhost:9999/cc_api/svcaccount_v1/WaterAccount/accountdetails"
	}
	urlToCall += "?refToken=" + token
	urlToCall += "&json=" + url.QueryEscape(jsonData)

	req, err := http.NewRequest("GET", urlToCall, nil)
	if err != nil {
		log.Printf("Got error %s", err.Error())
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error getting account details from Toronto Water [%s]!\n", err.Error())
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Printf("Calling Toronto Water failed with status code [%d]!\n", resp.StatusCode)
		return nil, errors.New("Error")
	}

	// extract body
	defer resp.Body.Close()
	var accountdetails Accountdetails
	err = json.NewDecoder(resp.Body).Decode(&accountdetails)
	if err != nil {
		log.Printf("Error processing Toronto Water account details response [%s]!\n", err.Error())
		return nil, err
	}
	if accountdetails.ResultCode != 200 {
		log.Println("Getting account details failed!")
		return nil, errors.New("Error")
	}

	// get meters
	meters := []Meter{}
	for _, premise := range accountdetails.PremiseList {
		meters = append(meters, premise.MeterList...)
	}

	return meters, nil
}

func GetData(meter Meter, date time.Time, token string, config helpers.Config) (*list.List, error) {

	dateString := date.Format("2006-01-02")
	log.Println("Getting consumption data for meter " + meter.MeterNumber + " and date " + dateString)

	// get data
	jsonData := "{\"API_OP\":\"CONSUMPTION\",\"ACCOUNT_NUMBER\":\"" + config.TorontoWater.AccountNumber +
		"\", \"MIU_ID\": \"" + meter.Miu +
		"\", \"START_DATE\": \"" + dateString +
		"\", \"END_DATE\": \"" + dateString +
		"\", \"INTERVAL_TYPE\" : \"Hour\"}"
	urlToCall := "https://secure.toronto.ca/cc_api/svcaccount_v1/WaterAccount/consumption"
	if config.TorontoWater.Mock {
		urlToCall = "http://localhost:9999/cc_api/svcaccount_v1/WaterAccount/consumption"
	}
	urlToCall += "?refToken=" + token
	urlToCall += "&json=" + url.QueryEscape(jsonData)

	req, err := http.NewRequest("GET", urlToCall, nil)
	if err != nil {
		log.Printf("Got error %s", err.Error())
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error getting data from Toronto Water [%s]!\n", err.Error())
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Printf("Calling Toronto Water failed with status code [%d]!\n", resp.StatusCode)
		return nil, errors.New("Error")
	}

	// extract body
	defer resp.Body.Close()
	var consumption Consumption
	err = json.NewDecoder(resp.Body).Decode(&consumption)
	if err != nil {
		log.Printf("Error processing Toronto Water consumption response [%s]!\n", err.Error())
		return nil, err
	}
	if consumption.ResultCode != 200 {
		log.Println("Getting consumption failed!")
		return nil, errors.New("Error")
	}

	consumptions := list.New()
	for _, meter := range consumption.MeterList {
		for _, interval := range meter.IntervalList {
			if interval.IntConsumptionTotal > 0.0 {
				readTime, err := time.ParseInLocation("2006-01-02 15:04:05", interval.LastReadDateTime, date.Location())
				if err != nil {
					log.Printf("Error parsing date [%s]!\n", err.Error())
					continue
				}
				value, err := strconv.ParseFloat(interval.LastReadValue, 64)
				if err != nil {
					log.Printf("Error parsing value [%s]!\n", err.Error())
					continue
				}
				value += interval.IntConsumptionTotal
				consumptions.PushBack(WaterConsumption{Time: readTime, Value: value})
			}
		}
	}

	return consumptions, nil
}
