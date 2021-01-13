package data

import (
	"strings"
	"time"

	"github.com/dta4/ot3c/enterprise"
	"github.com/huaweicloud/golangsdk"
	"gopkg.in/yaml.v2"
)

//GhostResource are VRs that once where a real VR but are now terminated and no longer produce any future costs. They are still needet for current cost calculation.
type GhostResource struct {
	ID          string
	consumption []enterprise.Consumption
}

func (c GhostResource) GetID() string {
	return c.ID
}

func (c GhostResource) GetDescription() string {
	bytes, _ := yaml.Marshal(c)
	return string(bytes)
}
func (c GhostResource) GetCostItems() []CostItem {

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
func (c GhostResource) GetPredCostItems() []CostItem {

	return nil
}
func (c GhostResource) GetTags() map[string]string {
	var nulTags map[string]string = make(map[string]string)
	nulTags["ot3c_prio"] = "0"
	return nulTags
}
func (c GhostResource) Terminate(client *golangsdk.ProviderClient) error {
	//NOP => cant terminate a already dead resource.
	return nil
}
func (c *GhostResource) BuildPastCostItems(cost []enterprise.Consumption) {
	c.consumption = cost
}
