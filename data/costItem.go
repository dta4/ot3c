package data

import (
	"fmt"
	"time"

	"github.com/dta4/ot3c/config"
)

//CostItem is an item that pairs cost to a single day of a resource
type CostItem struct {
	Cost float64
	Time time.Time
}

//GetCurrentCost calculates the current cost in this cost month
func GetCurrentCost(vr VirtualResource) float64 {
	current := CostItemBetween(vr.GetCostItems(), config.CurrentBillingStart(), config.CurrentBillingEnd())
	var cost float64
	for _, con := range current {
		cost = cost + con.Cost
	}
	return cost
}

//GetCostTillBillingEnd calculates the rest cost for this cost month
func GetCostTillBillingEnd(vr VirtualResource) float64 {
	y, m, d := time.Now().Date()
	fmt.Printf("ID: %v\n", vr.GetID())
	now := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	current := CostItemBetween(vr.GetPredCostItems(), now, config.CurrentBillingEnd())
	var cost float64
	for _, con := range current {
		cost = cost + con.Cost
	}
	return cost
}

//CostItemBetween returns all cost items between t1 and t2
func CostItemBetween(c []CostItem, t1 time.Time, t2 time.Time) []CostItem {

	var result []CostItem = make([]CostItem, 0)
	for _, con := range c {
		t := con.Time
		fmt.Println("-----")
		fmt.Printf("%v\n", t.String())
		fmt.Printf("%v\n", t1.String())
		fmt.Printf("%v\n", t2.String())
		if (t.After(t1) && t.Before(t2)) || t.Equal(t1) || t.Equal(t2) {
			result = append(result, con)
		}
	}
	return result
}

//CostBetween returns sums all costItems between t1 and t2
func CostBetween(c []CostItem, t1 time.Time, t2 time.Time) float64 {

	current := CostItemBetween(c, t1, t2)
	var cost float64
	for _, con := range current {
		cost = cost + con.Cost
	}
	return cost

}

//CostSortByDate is a list for sorting by date
type CostSortByDate []CostItem

func (a CostSortByDate) Len() int           { return len(a) }
func (a CostSortByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CostSortByDate) Less(i, j int) bool { return a[i].Time.Before(a[j].Time) }
