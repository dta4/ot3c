package enterprise

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dta4/ot3c/config"
	"github.com/sirupsen/logrus"
)

const consumptionEndpoint string = "https://api-enterprise-dashboard.otc-service.com/v1/consumption"

var log *logrus.Entry = logrus.WithFields(
	logrus.Fields{
		"package": "enterprise",
	},
)

//GetAllConsumptions returns all Consumtions
func GetAllConsumptions() ([]Consumption, error) {

	var list []Consumption

	//first call
	var offset int = 0
	var exit bool = false

	for !exit {

		api, err := GetConsumption(100, offset)
		if err != nil {
			log.Errorf("Error on api call: %v -> abort", err)
			return nil, err
		}
		list = append(list, api.Contents...)
		var perc float32

		perc = float32(api.Offset) / float32(api.Total)
		log.Infof("Retrieved %v%%", (perc * 100))
		offset = offset + len(api.Contents)
		if api.Total <= offset {
			exit = true
		}
	}
	log.Debugf("Retrieved %v resources", len(list))

	return list, nil
}

//GetConsumption makes a single API
func GetConsumption(limit int, offset int) (ConsumptionAPIResponds, error) {
	log := log.WithField("call", "get_consuption_page")
	client := http.DefaultClient
	u, _ := url.Parse(consumptionEndpoint)

	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), http.NoBody)
	if err != nil {
		log.WithError(err).Error("Error on request creation")
		return ConsumptionAPIResponds{}, err
	}
	//Adding API Key
	var bearer = "Bearer " + config.APIKey
	req.Header.Add("Authorization", bearer)

	resp, err := client.Do(req)

	if err != nil {
		log.WithError(err).Debug("Error on ED API call")
		return ConsumptionAPIResponds{}, err
	}
	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return ConsumptionAPIResponds{}, err
	}
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		log.WithError(err).WithField("body", string(buffer)).Debug("Error on ED API call")
		return ConsumptionAPIResponds{}, err
	}

	var result ConsumptionAPIResponds

	//log.Debug(string(buffer))
	err = json.Unmarshal(buffer, &result)

	if err != nil {
		log.Debug(err.Error())
		return ConsumptionAPIResponds{}, err
	}

	return result, nil
}

func GetConsuptionByID(id string, offset int) (ConsumptionAPIResponds, error) {
	log := log.WithField("call", "get_consuption_page")
	client := http.DefaultClient
	u, _ := url.Parse(consumptionEndpoint)

	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("ressource_id", id)
	q.Add("offset", strconv.Itoa(offset))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), http.NoBody)
	if err != nil {
		log.WithError(err).Error("Error on request creation")
		return ConsumptionAPIResponds{}, err
	}
	//Adding API Key
	var bearer = "Bearer " + config.APIKey
	req.Header.Add("Authorization", bearer)

	resp, err := client.Do(req)

	if err != nil {
		log.WithError(err).Debug("Error on ED API call")
		return ConsumptionAPIResponds{}, err
	}
	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debug(err.Error())
		return ConsumptionAPIResponds{}, err
	}
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		log.WithError(err).WithField("body", string(buffer)).Debug("Error on ED API call")
		return ConsumptionAPIResponds{}, err
	}

	var result ConsumptionAPIResponds

	log.Debug(string(buffer))
	err = json.Unmarshal(buffer, &result)

	if err != nil {
		log.Error(err.Error())
		return ConsumptionAPIResponds{}, err
	}

	return result, nil
}

//GetAllConsuptionByID returns all consumptions from an id
func GetAllConsuptionByID(id string) ([]Consumption, error) {
	var list []Consumption
	log := log.WithField("id", id)

	//first call
	var offset int = 0
	var exit bool = false

	for !exit {

		api, err := GetConsuptionByID(id, offset)
		if err != nil {
			log.Errorf("Error on api call: %v -> abort", err)
			return nil, err
		}
		list = append(list, api.Contents...)
		var perc float32

		perc = float32(api.Offset) / float32(api.Total)
		log.Debugf("Retrieved %v%%", (perc * 100))
		offset = offset + len(api.Contents)
		if api.Total <= offset {
			exit = true
		}
	}
	log.Debugf("Retrieved %v resources", len(list))

	return list, nil
}

//GetAllOT3CConsumptionsByDay returns all consupltions of a day that are relevant to ot3c
func GetAllOT3CConsumptionsByDay(date time.Time) ([]Consumption, error) {
	log := log.WithField("call", "get_consuption_page")
	client := http.DefaultClient
	u, _ := url.Parse(consumptionEndpoint)
	q, _ := url.ParseQuery(u.RawQuery)
	repeat := true
	i := 0
	var returnlist []Consumption = make([]Consumption, 0)
	for repeat {
		q.Add("offset", strconv.Itoa(i*20))
		q.Add("date", fmt.Sprintf("%v-%v-%v", date.Year(), int(date.Month()), date.Day()))
		q.Add("limit", "20")
		q.Add("tagged", "true")
		q.Add("tag_key", "ot3c_prio")
		u.RawQuery = q.Encode()

		req, err := http.NewRequest("GET", u.String(), http.NoBody)
		if err != nil {
			log.WithError(err).Error("Error on request creation")
			return returnlist, err
		}
		//Adding API Key
		var bearer = "Bearer " + config.APIKey
		req.Header.Add("Authorization", bearer)

		resp, err := client.Do(req)

		if err != nil {
			log.WithError(err).Debug("Error on ED API call")
			return returnlist, err
		}
		buffer, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Debug(err.Error())
			return returnlist, err
		}
		if resp.StatusCode != 200 {
			err = errors.New(resp.Status)
			log.WithError(err).WithField("body", string(buffer)).WithField("date", fmt.Sprintf("%v-%v-%v", date.Year(), int(date.Month()), date.Day())).Error("Error on ED API call")
			return returnlist, err
		}

		var result ConsumptionAPIResponds

		log.Debug(string(buffer))
		err = json.Unmarshal(buffer, &result)

		if err != nil {
			log.Error(err.Error())
			return returnlist, err
		}

		if result.Total != len(returnlist) {
			returnlist = append(returnlist, result.Contents...)
			i++
		} else {
			repeat = false
		}

	}

	return returnlist, nil

}
