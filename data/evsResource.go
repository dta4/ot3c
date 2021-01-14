package data

import (
	"fmt"
	"strings"
	"time"

	"github.com/dta4/ot3c/config"
	"github.com/dta4/ot3c/enterprise"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack"
	"github.com/huaweicloud/golangsdk/openstack/blockstorage/v3/volumes"
	hVolumes "github.com/huaweicloud/golangsdk/openstack/evs/v3/volumes"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var evslog *logrus.Entry = logrus.WithField("module", "evs_data")

//ECSResource ECSResource
type EVSResource struct {
	ID          string
	consumption []enterprise.Consumption
	OTCEVS      hVolumes.Volume
	PredCost    []CostItem
}

func (c EVSResource) GetID() string {
	return c.ID
}

func (c EVSResource) GetDescription() string {
	bytes, _ := yaml.Marshal(c)
	return string(bytes)
}
func (c EVSResource) GetCostItems() []CostItem {

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
func (c EVSResource) GetPredCostItems() []CostItem {

	return c.PredCost
}
func (c EVSResource) GetTags() map[string]string {
	return c.OTCEVS.Tags
}
func (c EVSResource) Terminate(client *golangsdk.ProviderClient) error {

	bs, err := openstack.NewBlockStorageV3(client, golangsdk.EndpointOpts{
		Region: "eu-de",
	})
	if config.Dryrun {
		fmt.Printf("MOCK TERMINATE: %v\n", c.ID)
		return nil
	}
	r := volumes.Delete(bs, c.OTCEVS.ID)

	if r.Err != nil {
		return err
	}
	evslog.WithField("time", time.Now()).Infof("Terminated %v\n", c.ID)
	return nil
}
func (c *EVSResource) BuildPastCostItems(cost []enterprise.Consumption) {
	c.consumption = cost
}

func (c *EVSResource) SetPredCostItems(cost []CostItem) {
	c.PredCost = cost
}
