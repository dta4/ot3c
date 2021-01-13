package config

import (
	"time"

	"github.com/rickb777/date/period"
)

//BillingBegin is the Time/Date the billing cicle beginns
var BillingBegin time.Time

//TerminateASAPGracePeriod is a grace period that a VR gets if it needs to be terminated ASAP.
var TerminateASAPGracePeriod time.Duration

//CalcTimeLeftInBillPeriod calculates the time left in the current Billing period
func CalcTimeLeftInBillPeriod() time.Duration {
	now := time.Now()

	timeLeft, _ := period.Between(now, CurrentBillingEnd()).Duration()
	return timeLeft
}

//CurrentBillingEnd returns the End Date of the current Billing Mounth
func CurrentBillingEnd() time.Time {
	now := time.Now()
	billingEnd := BillingBegin

	for !billingEnd.After(now) {
		billingEnd = billingEnd.AddDate(0, 1, 0)

	}
	return billingEnd
}

//CurrentBillingStart returns the start Date of current Billing Mounth
func CurrentBillingStart() time.Time {
	return CurrentBillingEnd().AddDate(0, -1, 0)
}

//GetTerminateASAPDate calculates the time for the Termination Date
func GetTerminateASAPDate() time.Time {
	return time.Now().Add(TerminateASAPGracePeriod)
}
