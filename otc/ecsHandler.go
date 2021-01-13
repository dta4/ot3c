package otc

import (
	"github.com/dta4/ot3c/data"
	"github.com/dta4/ot3c/enterprise"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack"
	"github.com/huaweicloud/golangsdk/pagination"
	"github.com/sirupsen/logrus"

	"github.com/avast/retry-go"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/servers"
	"github.com/huaweicloud/golangsdk/openstack/ecs/v1/tags"
)

var ecsLog *logrus.Entry

const otcRegion string = "eu-de"

func loadEcsResources(client *golangsdk.ProviderClient) error {
	ecsLog = dataLog.WithField("stage", "load_ecs")
	ecsLog.Info("Loading ECS resources")

	compute, err := openstack.NewComputeV2(client, golangsdk.EndpointOpts{
		Region: otcRegion,
	})
	ecs, err := openstack.NewEcsV1(client, golangsdk.EndpointOpts{
		Region: otcRegion,
	})

	if err != nil {
		ecsLog.WithError(err).Error("Error loading ecs resources")
		return err
	}
	var page pagination.Page
	err = retry.Do(func() error {
		page, err = servers.List(compute, servers.ListOpts{
			Limit: 100,
		}).AllPages()
		return err
	})

	if err != nil {
		ecsLog.WithError(err).Error("Error loading ecs resources")
		return err
	}
	serversList, err := servers.ExtractServers(page)

	for _, server := range serversList {

		err := retry.Do(
			func() error {
				return loadSingleEcsServer(ecs, &server)
			},
		)
		if err != nil {
			ecsLog.WithError(err).Error("Error loading ecs resources")
			return err
		}

	}

	return nil

}

func loadSingleEcsServer(ecsSC *golangsdk.ServiceClient, server *servers.Server) error {
	//Get Cost
	cost, err := enterprise.GetAllConsuptionByID(server.ID)

	if err != nil {
		ecsLog.WithError(err).Warn("Error loading ecs cost")
		err = nil
	}

	//Get Tags

	t := tags.Get(ecsSC, server.ID)
	tas, err := t.Extract()
	if err != nil {
		ecsLog.WithError(err).Error("Error loading ecs tags")
		return err
	}
	var tagsMap map[string]string = make(map[string]string)
	for _, tagOTC := range tas.Tags {
		tagsMap[tagOTC.Key] = tagOTC.Value
	}
	//create serverResource
	serverResource := data.ECSResource{
		ID:        "ECS_" + server.ID,
		OTCServer: *server,
		Tags:      tagsMap,
	}
	serverResource.SetPastCostItems(cost)
	data.AddVirtualResource(&serverResource)
	return nil
}

//GroupEcsEvsResources creates groupResources from ecs and evs pairs
func GroupEcsEvsResources() error {
	ecsLog = dataLog.WithField("stage", "process_ecs-evs")

	for _, vr := range data.Resources {
		//Search for ecs resource
		ecs, ok := vr.(*data.ECSResource)
		if ok {
			//ECS found
			gr, err := GroupEcsEvsResource(ecs)
			if err != nil {
				ecsLog.WithError(err).Error("Error on ECS-EVS grouping")
				return err
			}
			data.AddVirtualResource(gr)

		}
	}

	return nil
}
