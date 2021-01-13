package otc

import (
	"github.com/dta4/ot3c/data"
	"github.com/dta4/ot3c/prediction"
	"github.com/sirupsen/logrus"
)

var dataLog logrus.Entry

//LoadVirtualResources loads all resources from the OTC and links them with cost. After that it will look for other costs that are not represented by VRs and add them as Ghosts.
func LoadVirtualResources() error {
	dataLog = *logrus.WithField("module", "dataHandler")
	client, err := Login()
	if err != nil {
		dataLog.WithError(err).Error("OTC login error")
		return err
	}
	if len(data.Resources) > 0 {
		dataLog.Debug("Clearing old resources")
		data.Resources = make([]data.VirtualResource, 0)
	}
	//Loading OTC and Cost data
	dataLog.Info("Loading Resources...")
	loadEcsResources(client)
	loadEvsResources(client)

	//Load Ghosts
	dataLog.Info("Loading Ghost Resources")
	LoadGhostResources()

	//Post done
	dataLog.Info("Loading completed")
	return nil
}

//RunPostProcessing performs postprocessing work after all resources are loaded
func RunPostProcessing() {
	//Linking and Grouping of data
	postLog := dataLog.WithField("stage", "postprocessing")
	postLog.Info("Starting Post Processing")
	postLog.Info("Grouping ECS-EVS Resources")
	GroupEcsEvsResources()
	postLog.Info("Grouping Resources by Tags")
	GroupVRByTags()
	postLog.Info("Post Processing done")
}

//RunDefaultDataChain combines VR func that are regularly used together to load and prepare data.
func RunDefaultDataChain() error {
	err := LoadVirtualResources()
	if err != nil {
		return err
	}
	RunPostProcessing()
	data.FilterNonOT3C()
	err = prediction.RunCostPredictionOnAll()
	if err != nil {
		return err
	}
	data.SortResourcesByPriority()
	return nil
}
