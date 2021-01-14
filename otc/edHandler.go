package otc

import (
	"fmt"
	"strings"
	"time"

	"github.com/dta4/ot3c/config"
	"github.com/dta4/ot3c/data"
	ed "github.com/dta4/ot3c/enterprise"
	"github.com/sirupsen/logrus"
)

var datalog *logrus.Entry = log.WithField("module", "ed_datahandler")

//LoadGhostResources creates GhostResources from Consumptions of VRs that are no longer present and have been terminated. Make sure to load all other resources first before loading the Ghosts.
func LoadGhostResources() {
	datalog = datalog.WithField("phase", "load GhostResources")
	var fullIDMap map[string][]ed.Consumption = make(map[string][]ed.Consumption)

	startDate := config.CurrentBillingStart()
	selectDate := startDate
	currentDate := time.Now()
	for currentDate.After(selectDate) {
		cons, err := ed.GetAllOT3CConsumptionsByDay(selectDate)
		if err != nil {
			log.WithError(err).WithField("date", selectDate.String()).Error("Error while loading consumptions of a day")
			continue
		}
		for _, con := range cons {
			fullIDMap[con.ResourceID] = append(fullIDMap[con.ResourceID], con)
		}
		selectDate = selectDate.Add(time.Hour * 24)

	}
	//find ghosts
	for ids, cons := range fullIDMap {
		res := data.FindVRByID(ids)
		if res == nil {
			createGhostResource(ids, cons)
		}
	}

}

func createGhostResource(id string, consuptions []ed.Consumption) {
	resourcePrefix := strings.Split(consuptions[0].Product, "_")[1]
	ghost := data.GhostResource{
		ID: fmt.Sprintf("GH_%v_%v", resourcePrefix, id),
	}
	ghost.BuildPastCostItems(consuptions)
	data.AddVirtualResource(ghost)
}
