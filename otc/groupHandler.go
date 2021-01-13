package otc

import (
	"errors"

	"github.com/dta4/ot3c/data"
)

//GroupEcsEvsResource creates groupResource from ecs and evs pairs
func GroupEcsEvsResource(ecs *data.ECSResource) (data.GroupResource, error) {
	var VRList []data.VirtualResource = make([]data.VirtualResource, 0)
	VRList = append(VRList, ecs)
	//ECS found
	for _, val := range ecs.OTCServer.VolumesAttached {
		//Find EVS with id
		evs, ok := FindEVSVolumeWithID(val["id"])
		if !ok {
			ecsLog.Errorf("Cant find EVS with ID %v", val["id"])
			return data.GroupResource{}, errors.New("Cant find EVS VR")
		}
		if _, ok := evs.GetTags()["ot3c_group"]; !ok {
			VRList = append(VRList, evs)
		} else {
			ecsLog.Debug("Ignoring EVS due to grouping tag")
		}

	}
	tags := make(map[string]string)
	if group, ok := ecs.GetTags()["ot3c_group"]; ok {
		tags["ot3c_group"] = group
	}
	result := data.GroupResource{
		List: VRList,
		Tags: tags,
	}
	result.PatchGroup()
	return result, nil
}

//GroupVRByTags creates GR based on ot3c_group tags.
func GroupVRByTags() {
	var groupMap map[string][]data.VirtualResource = make(map[string][]data.VirtualResource)

	for _, vr := range data.Resources {
		groupString, ok := vr.GetTags()["ot3c_group"]
		if !ok {
			continue
		}
		_, ok = groupMap[groupString]
		if !ok {
			groupMap[groupString] = make([]data.VirtualResource, 0)
		}
		groupMap[groupString] = append(groupMap[groupString], vr)

	}

	for groupString, vrlist := range groupMap {
		tags := make(map[string]string)
		tags["ot3c_group"] = groupString
		gr := data.GroupResource{
			List: vrlist,
			Tags: tags,
		}
		gr.PatchGroup()
		for _, vr := range vrlist {
			data.RemoveResourceWithID(vr.GetID())
		}
		data.AddVirtualResource(gr)
	}
}
