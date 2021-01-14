package termination

import (
	"fmt"
	"time"

	"github.com/dta4/ot3c/data"
	"github.com/dta4/ot3c/otc"
	"github.com/hashicorp/go-uuid"
	"github.com/sirupsen/logrus"
)

var termLog *logrus.Entry = logrus.WithField("module", "termination")
var daemonLock chan bool = make(chan bool, 1)

//RunTerminationPlan starts a long running process for executing a Termination Plan. Returns on update
func RunTerminationPlan(update *chan bool, waitDuration time.Duration) {
	//locking daemon
	daemonLock <- true
	id, _ := uuid.GenerateUUID()
	runLog := termLog.WithField("id", id).WithField("component", "daemon")
	runLog.Info("Starting Termination Daemon...")
	if len(data.ActiveTerminationPlan.Items) == 0 {
		runLog.Error("Can not execute Termination Plan. No Items available.")
		return
	}
	done := false

	timer1 := time.NewTimer(time.Millisecond * 1)

	for !done {
		select {
		case <-*update:
			//Termination Plan Updated
			done = true
			timer1.Stop()
		case <-timer1.C:
			//Wait time is up
			ExecTerminationRun()
			timer1.Reset(waitDuration)
			runLog.Infof("Sleeping until %v for next run... (%v)", time.Now().Add(waitDuration).String(), waitDuration.String())
		}
	}
	runLog.Info("Termination Daemon done")
	//Releasing Lock
	<-daemonLock
}

func getValidItems() []data.TerminationItem {
	var validItems []data.TerminationItem = make([]data.TerminationItem, 0)
	for _, item := range data.ActiveTerminationPlan.Items {
		vr := data.FindVRByID(item.TVR.GetID())
		if vr != nil {
			//VR is still existing
			validItems = append(validItems, item)
		}
	}
	return validItems

}

//ExecTerminationRun performs a one time termination run
func ExecTerminationRun() {
	id, _ := uuid.GenerateUUID()
	runLog := termLog.WithField("id", id).WithField("component", "run")
	runLog.WithField("time", time.Now().String()).Info("Termination Run started...")
	fmt.Println()
	err := updateVRs()
	if err != nil {
		runLog.WithError(err).Error("Error while updating VRs")
		return
	}
	fmt.Println()
	runLog.Info("Current Termination Plan...")
	fmt.Println(data.CurrentTerminationPlanToString())
	fmt.Println()
	//Creating OTC client
	client, _ := otc.Login()

	items := getValidItems()
	current := time.Now()
	count := 0
	for _, item := range items {
		//if termination date is due
		if current.After(item.TDate) {
			runLog.Infof("Terminating %v", item.TVR.GetID())
			//execute termination
			err := item.TVR.Terminate(client)
			if err != nil {
				runLog.WithError(err).Error("Error on VR Termination")
			}
			count++
			err = terminationNotify(item.TVR)
			if err != nil {
				runLog.WithError(err).Warn("Error on Notification")
			}

		}
	}
	if count == 0 {
		runLog.Info("No VR has been terminated. (Have not reached termination date)")
	} else {
		runLog.Infof("%v VRs have been terminated.", count)
	}
	runLog.Info("Termination Run done")
	fmt.Println()
}

func terminationNotify(vr data.VirtualResource) error {
	termLog.Warningf("[MOCK] Notifing VR Owner with id %v", vr.GetID())
	return nil
}

func updateVRs() error {
	return otc.RunDefaultDataChain()
}
