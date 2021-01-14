package data

import (
	"strconv"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var logGR *logrus.Entry = logrus.WithField("module", "GroupResource")

//GroupResource is a grouped resource. IDs should start with GR_
type GroupResource struct {
	List []VirtualResource
	Tags map[string]string
}

func (g GroupResource) GetID() string {
	var id string = "GR_["
	for i, vr := range g.List {
		id = id + vr.GetID()
		if i != (len(g.List) - 1) {
			id = id + ","
		}
	}
	id = id + "]"
	return id
}
func (g GroupResource) GetDescription() string {
	bytes, _ := yaml.Marshal(g)
	return string(bytes)
}

func (g GroupResource) GetCostItems() []CostItem {
	var costMap map[int64]float64 = make(map[int64]float64)

	for _, res := range g.List {
		items := res.GetCostItems()
		for _, item := range items {

			costMap[item.Time.Unix()] = costMap[item.Time.Unix()] + item.Cost
		}
	}
	var costItems []CostItem = make([]CostItem, 0)
	for t, v := range costMap {
		costItems = append(costItems, CostItem{
			Time: time.Unix(t, 0),
			Cost: v,
		})
	}
	return costItems
}
func (g GroupResource) GetPredCostItems() []CostItem {
	var costMap map[time.Time]float64 = make(map[time.Time]float64)

	for _, res := range g.List {
		items := res.GetPredCostItems()
		for _, item := range items {
			costMap[item.Time] = costMap[item.Time] + item.Cost
		}
	}
	var costItems []CostItem = make([]CostItem, 0)
	for t, v := range costMap {
		costItems = append(costItems, CostItem{
			Time: t,
			Cost: v,
		})
	}
	return costItems
}
func (g GroupResource) GetTags() map[string]string {
	if g.Tags == nil {
		g.PatchGroup()
	}
	return g.Tags
}
func (g GroupResource) Terminate(client *golangsdk.ProviderClient) error {
	for _, vr := range g.List {
		err := vr.Terminate(client)
		if err != nil {
			return err
		}
	}
	return nil
}

//PatchGroup searches through all the VRs and patches different tags on the group
//Set highest prio. Searches through all VRs, gets the highes prio and sets it for the group.
func (g GroupResource) PatchGroup() {
	patchHighestPrio(g)
	patchRemoveDuplicates(g)
}
func patchHighestPrio(g GroupResource) {

	var prio int = -1
	for _, vr := range g.List {
		//If vr contains ot3c_prio
		if val, ok := vr.GetTags()["ot3c_prio"]; ok {
			var setPrio bool = false
			//If Prio is not set, set it.
			if prio == -1 {
				setPrio = true
			}
			//Convert prio string to int
			cPrio, err := strconv.Atoi(val)
			if err != nil {
				logGR.WithField("resource", vr.GetID()).WithError(err).Errorln("Error while parsing Tags")
				return
			}

			if cPrio < prio {
				setPrio = true
			}
			//if the prio needs to be set, set it.
			if setPrio {
				prio = cPrio
			}
		}
	}
	//Set group
	if prio != -1 {
		g.Tags["ot3c_prio"] = strconv.Itoa(prio)
	}

}

//Removes VR from root list that are in this group.
func patchRemoveDuplicates(g GroupResource) {
	for _, vr := range g.List {
		//Check if VR is GR
		grc, ok := vr.(GroupResource)
		if ok {
			//Recurcive Remove
			patchRemoveDuplicates(grc)
		}
		ok = RemoveResourceWithID(vr.GetID())
	}
}
