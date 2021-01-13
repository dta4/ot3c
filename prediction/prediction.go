package prediction

import (
	"github.com/dta4/ot3c/data"
	"github.com/sirupsen/logrus"
)

var predLog *logrus.Entry = logrus.WithField("module", "prediction")

//RunCostPrediction selects and then runs a Cost Prediction model that fits to the current data coverage
func RunCostPrediction(data []data.CostItem) []data.CostItem {
	return runSteadyState(data)
}

//RunCostPredictionOnVR runs a cost Prediction on a VR. If VR is grouped then it will decend recursively until the leaves are reached.
func RunCostPredictionOnVR(vr data.VirtualResource) error {
	didDoPred := false
	predLog.Debugf("Running Cost prediction on %v", vr.GetID())
	gr, gok := vr.(data.GroupResource)
	ecs, ecok := vr.(*data.ECSResource)
	evs, evok := vr.(*data.EVSResource)
	var err error = nil
	if gok {
		for _, gvr := range gr.List {
			err = RunCostPredictionOnVR(gvr)
		}
		didDoPred = true
	} else if ecok {
		pred := RunCostPrediction(ecs.GetCostItems())
		ecs.SetPredCostItems(pred)
		didDoPred = true
	} else if evok {
		pred := RunCostPrediction(evs.GetCostItems())
		evs.SetPredCostItems(pred)
		didDoPred = true
	}

	if didDoPred {
		predLog.Debugf("Cost prediction on %v done", vr.GetID())
	}
	return err
}

//RunCostPredictionOnAll Iterates over all ressources and runs a cost prediction
func RunCostPredictionOnAll() error {
	compleat := len(data.Resources)
	index := 0
	predLog.Infof("Running cost prediction on %v resources", compleat)
	for _, vrs := range data.Resources {
		err := RunCostPredictionOnVR(vrs)
		if err != nil {
			return err
		}
		index++
		a := float32(index) / float32(compleat)
		predLog.Infof("%v%% predicted (%v/%v)", a*100, index, compleat)
	}
	return nil
}
