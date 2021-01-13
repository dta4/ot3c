package otc

import (
	"github.com/avast/retry-go"
	"github.com/dta4/ot3c/data"
	"github.com/dta4/ot3c/enterprise"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack"
	"github.com/huaweicloud/golangsdk/openstack/blockstorage/v3/volumes"
	hVolumes "github.com/huaweicloud/golangsdk/openstack/evs/v3/volumes"
	"github.com/huaweicloud/golangsdk/pagination"
	"github.com/sirupsen/logrus"
)

var evsLog *logrus.Entry

func loadEvsResources(client *golangsdk.ProviderClient) error {
	evsLog = dataLog.WithField("stage", "load_evs")
	evsLog.Info("Loading EVS resources")

	bs, err := openstack.NewBlockStorageV3(client, golangsdk.EndpointOpts{
		Region: "eu-de",
	})

	if err != nil {
		evsLog.WithError(err).Error("Error loading evs resources")
		return err
	}
	var page pagination.Page
	err = retry.Do(
		func() error {
			page, err = volumes.List(bs, volumes.ListOpts{
				Limit: 100,
			}).AllPages()
			return err
		},
	)

	if err != nil {
		evsLog.WithError(err).Error("Error loading evs resources")
		return err
	}

	volumeList, err := volumes.ExtractVolumes(page)

	for _, volume := range volumeList {
		retry.Do(func() error {
			return loadSingleEvsVolume(bs, &volume)
		})

	}

	return nil

}

func loadSingleEvsVolume(evsSC *golangsdk.ServiceClient, volume *volumes.Volume) error {
	//Reload with huawei api
	r := hVolumes.Get(evsSC, volume.ID)

	otcVolume, err := r.Extract()
	if err != nil {
		evsLog.WithError(err).Error("Error loading evs tags")
		return err
	}

	//Get Cost
	cost, err := enterprise.GetAllConsuptionByID(otcVolume.ID)

	if err != nil {
		evsLog.WithError(err).Warn("Error loading evs cost")
		err = nil
	}

	//create serverResource
	evsResource := data.EVSResource{
		ID:     "EVS_" + otcVolume.ID,
		OTCEVS: *otcVolume,
	}
	evsResource.BuildPastCostItems(cost)
	data.AddVirtualResource(&evsResource)

	return nil
}

//FindEVSVolumeWithID returns an EVSResource based on id without prefix
func FindEVSVolumeWithID(id string) (*data.EVSResource, bool) {
	for _, vr := range data.Resources {
		evs, ok := vr.(*data.EVSResource)
		if ok {
			//EVS found
			if evs.GetID() == ("EVS_" + id) {
				//Is evs searched for
				return evs, true
			}
		}
	}
	return &data.EVSResource{}, false
}
