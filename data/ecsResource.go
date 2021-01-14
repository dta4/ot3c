package data

import (
	"fmt"
	"strings"
	"time"

	"github.com/dta4/ot3c/config"
	"github.com/dta4/ot3c/enterprise"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/servers"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var ecslog *logrus.Entry = logrus.WithField("module", "ecs_data")

//ECSResource ECSResource
type ECSResource struct {
	ID          string
	consumption []enterprise.Consumption
	OTCServer   servers.Server
	PredCost    []CostItem
	Tags        map[string]string
}

func (c ECSResource) GetID() string {
	return c.ID
}

func (c ECSResource) GetDescription() string {
	bytes, _ := yaml.Marshal(c)
	return string(bytes)
}
func (c ECSResource) GetCostItems() []CostItem {

	var cost []CostItem
	for _, con := range c.consumption {
		time, _ := time.Parse("2006-01-02", strings.Split(con.ConsumptionDate, " ")[0])
		cost = append(cost, CostItem{
			Time: time,
			Cost: con.ListpriceAmount,
		})
	}
	return cost
}
func (c ECSResource) GetPredCostItems() []CostItem {

	return c.PredCost
}
func (c ECSResource) GetTags() map[string]string {
	return c.Tags
}
func (c ECSResource) Terminate(client *golangsdk.ProviderClient) error {

	cs, err := openstack.NewComputeV2(client, golangsdk.EndpointOpts{
		Region: "eu-de",
	})
	if config.Dryrun {
		fmt.Printf("MOCK TERMINATE: %v\n", c.ID)
		return nil
	}

	r := servers.Delete(cs, c.OTCServer.ID)

	if r.Err != nil {
		return err
	}
	evslog.WithField("time", time.Now()).Infof("Terminated %v\n", c.ID)
	return nil
	//WIP
	fmt.Printf("MOCK TERMINATE: %v\n", c.ID)
	return nil
}

//SetPastCostItems build past cost items
func (c *ECSResource) SetPastCostItems(cost []enterprise.Consumption) {

	c.consumption = cost

}

func (c *ECSResource) SetPredCostItems(cost []CostItem) {
	c.PredCost = cost
}
