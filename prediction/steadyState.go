package prediction

import (
	"time"

	"github.com/dta4/ot3c/data"
)

//runs a steadyState Prediction for 30 days in the future
func runSteadyState(cost []data.CostItem) []data.CostItem {

	daysLeft := 30
	if len(cost) == 0 {
		return make([]data.CostItem, 0)
	}
	lastState := cost[len(cost)-1]
	result := make([]data.CostItem, int(daysLeft))
	currentTime := lastState.Time
	for i := 0; i < int(daysLeft); i++ {

		result[i] = data.CostItem{
			Cost: lastState.Cost,
			Time: currentTime.Add(time.Hour * 24),
		}
		currentTime = currentTime.Add(time.Hour * 24)
	}
	return result
}
