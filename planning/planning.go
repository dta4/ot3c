package planning

import (
	"errors"
	"sort"
	"time"

	"github.com/dta4/ot3c/config"
	"github.com/dta4/ot3c/data"
	"github.com/sirupsen/logrus"
)

var planLog *logrus.Entry = logrus.WithField("module", "planning")

//BudgetAtEndOfMonth calculates rest budget at end of the cost month. (Positive = under the budget | Negative = over the budget)
//Requires cost prediction to be done before calling.
func BudgetAtEndOfMonth() float64 {
	pred := TotalRestCost()
	curr := TotalCurrentCost()
	budget := config.TargetBudget
	budget = budget - (pred + curr)
	return budget
}

//TotalRestCost calculates rest cost of the month
func TotalRestCost() float64 {
	var total float64
	for _, vrs := range data.Resources {
		total = total + data.GetCostTillBillingEnd(vrs)
	}
	return total
}

//TotalRestCost calculates rest cost of the month
func TotalCurrentCost() float64 {
	var total float64
	for _, vrs := range data.Resources {
		total = total + data.GetCurrentCost(vrs)
	}
	return total
}

//CreateTerminationPlan creates a Termination Plan given a specific ammount of money to save
func CreateTerminationPlan(budgetToReduce float64) ([]data.TerminationItem, error) {
	planLog.Info("Creating Termination Plan")
	data.SortResourcesByPriority()

	var terminationPlan []data.TerminationItem = make([]data.TerminationItem, 0)
	var toSave float64 = budgetToReduce
	for _, vr := range data.Resources {
		prio := data.GetPrio(vr)
		if prio == 0 {
			//if vr list is now non touchable
			return terminationPlan, errors.New("End of manageable Resource list")
		}
		ti := createTerminationItem(vr, toSave)
		toSave = toSave - ti.TSave
		terminationPlan = append(terminationPlan, ti)
		if toSave <= 0 {
			break
		}
	}
	planLog.Info("Termination Plan created")
	return terminationPlan, nil
}

//createTerminationItem creates a termination item.
func createTerminationItem(vr data.VirtualResource, toSave float64) data.TerminationItem {
	tillEnd := data.GetCostTillBillingEnd(vr)
	if tillEnd > toSave {
		//can terminate later
		costItemsCurrent := data.CostItemBetween(vr.GetPredCostItems(), time.Now(), config.CurrentBillingEnd())
		costItemSort := data.CostSortByDate(costItemsCurrent)
		sort.Sort(sort.Reverse(costItemSort))

		var terminationDate time.Time = config.GetTerminateASAPDate()
		for _, ci := range costItemSort {
			toSave = toSave - ci.Cost
			if toSave <= 0 {
				//Termination date found
				terminationDate = ci.Time
				break
			}
		}
		costReduction := data.CostBetween(vr.GetPredCostItems(), terminationDate, config.CurrentBillingEnd())
		item := data.TerminationItem{
			TVR:   vr,
			TDate: terminationDate,
			TSave: costReduction,
			TASAP: false,
		}
		return item

	} else {
		costReduction := data.CostBetween(vr.GetPredCostItems(), config.GetTerminateASAPDate(), config.CurrentBillingEnd())
		//must terminate ASAP
		item := data.TerminationItem{
			TVR:   vr,
			TDate: config.GetTerminateASAPDate(),
			TSave: costReduction,
			TASAP: true,
		}
		return item
	}

}
