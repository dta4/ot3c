package data

import (
	"sort"
	"strconv"
	"strings"

	"github.com/huaweicloud/golangsdk"
	"github.com/sirupsen/logrus"
)

//VirtualResource is the internal representation of a resource
type VirtualResource interface {
	//GetID returns a unique ID for the VR
	GetID() string
	//GetDescription returns a description for the VR
	GetDescription() string
	//GetCostItems returns all past costs up to the current date
	GetCostItems() []CostItem
	//GetPredCostItems returns future costs. Length can be 0 if no prediction has been done yet.
	GetPredCostItems() []CostItem
	//GetTags returns the VRs Tags
	GetTags() map[string]string
	//Terminate terminates the VR
	Terminate(*golangsdk.ProviderClient) error
}

//Resources is the in-memory representation of the cloud environment
var Resources []VirtualResource

//AddVirtualResource adds a new VirtualResource
func AddVirtualResource(resource VirtualResource) {
	Resources = append(Resources, resource)
}

//RemoveResourceWithID removes a VR based on its ID. Returns true on found and removed. Caution ... destroys order
func RemoveResourceWithID(id string) bool {
	for i, vr := range Resources {
		if vr.GetID() == id {
			if len(Resources) > 1 {
				Resources[i] = Resources[len(Resources)-1]
				Resources = Resources[:len(Resources)-1]
			} else {
				Resources = make([]VirtualResource, 0)
			}

			return true
		}
	}
	return false
}

var logsort *logrus.Entry = logrus.WithField("module", "VR_Sort")

//SortResourcesByPriority rearranges the VRs in the Resource List based on Priority decending.
func SortResourcesByPriority() {
	sbp := SortByPrio(Resources)
	sort.Sort(sbp)
	Resources = sbp
}

type SortByPrio []VirtualResource

func (a SortByPrio) Len() int      { return len(a) }
func (a SortByPrio) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a SortByPrio) Less(i, j int) bool {

	r1, err := strconv.Atoi(a[i].GetTags()["ot3c_prio"])
	if err != nil {
		logsort.Error(err)
		r1 = 0
	}
	r2, err := strconv.Atoi(a[j].GetTags()["ot3c_prio"])
	if err != nil {
		logsort.Error(err)
		r2 = 0
	}
	return r1 > r2
}

//FilterNonOT3C removes all VRs that have no ot3c_prio tag
func FilterNonOT3C() {
	for _, vr := range Resources {
		_, ok := vr.GetTags()["ot3c_prio"]
		if ok == false {
			RemoveResourceWithID(vr.GetID())
		}
	}
}

//GetPrio returns a prio for a vr if available. If not then it returns -1
func GetPrio(vr VirtualResource) int {
	a, ok := vr.GetTags()["ot3c_prio"]
	if ok {
		b, _ := strconv.Atoi(a)
		return b
	} else {
		return -1
	}

}

//FindVRByID returns a VR that contains the id in question . Works also on grouped VRs.
func FindVRByID(id string) VirtualResource {
	for _, vrs := range Resources {
		if strings.Contains(vrs.GetID(), id) {
			return vrs
		}
	}
	return nil
}
