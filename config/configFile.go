package config

import (
	"errors"
	"time"

	"github.com/google/martian/log"
	huaweiotc "github.com/huaweicloud/golangsdk"
)

//File is a struct that represents the config file for OT3C. When loaded it can be transphered into the different static Vars in this package to be accessed by other packages.
type File struct {
	EnterpriseDashboardAPISecret string `yaml:"dashboard-secret"`
	OTCCreds                     huaweiotc.AKSKAuthOptions
	BillingBegin                 string        `yaml:"billing-begin"`
	TargetBudget                 float64       `yaml:"target-budget"`
	TerminationGrace             time.Duration `yaml:"termination-grace"`
}

//SampleConfig is just a sample config
var SampleConfig File = File{
	EnterpriseDashboardAPISecret: "<very-long-api-secret>",
	OTCCreds: huaweiotc.AKSKAuthOptions{
		IdentityEndpoint: "https://iam.eu-de.otc.t-systems.com/v3",
		AccessKey:        "<very-pretty-access-key>",
		SecretKey:        "<very-pretty-secret-key>",
		ProjectId:        "<very-long-project-id>",
	},
	BillingBegin:     "2016-01-05",
	TargetBudget:     100.0,
	TerminationGrace: time.Hour * 24,
}

//ApplyConfig applies the config from the file to OT3C
func (file *File) ApplyConfig() error {
	APIKey = file.EnterpriseDashboardAPISecret
	OTCCreds = file.OTCCreds
	OTCUsingOTP = false
	begin, err := time.Parse("2006-01-02", file.BillingBegin)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}
	BillingBegin = begin
	TargetBudget = file.TargetBudget
	if TargetBudget <= 0 {
		err = errors.New("Target Budget is not valid (budget<0)")
		log.Errorf(err.Error())
		return err
	}
	TerminateASAPGracePeriod = file.TerminationGrace
	if TerminateASAPGracePeriod == 0 {
		TerminateASAPGracePeriod = time.Hour * 24
	}
	return nil
}
