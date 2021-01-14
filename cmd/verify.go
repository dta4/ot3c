/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/dta4/ot3c/enterprise"
	"github.com/dta4/ot3c/otc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var verifyLog *logrus.Entry
var preflightResult bool = true
var preflightLog *logrus.Entry

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verifies the current config and performs a preflight check.",
	Long:  `This command checks the provided config regarding credentials/sanity and performs a preflight check.`,
	Run: func(cmd *cobra.Command, args []string) {
		Preflight()
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	verifyLog = logrus.WithField("command", "verify")
}

//Preflight checks if the config contains corrent credentails and pings different apis for verification.
func Preflight() bool {
	preflightLog = logrus.WithField("phase", "preflight")
	preflightLog.Info("Performing preflight checks")
	checkOTC()
	checkED()
	preflightLog.WithField("step", "preflight_check")
	if preflightResult {
		preflightLog.Info("Preflight checks succeeded")
		return true
	} else {
		preflightLog.Warn("Preflight checks failed")
		return false
	}

}

//checkOTC retruns true if authentication is correct and if access is given(LAST PART IS WIP)
func checkOTC() {
	otcLog := preflightLog.WithField("step", "OTC_API_cred_check")
	otcLog.Debugln("Checking OTC creds...")
	_, err := otc.Login()

	if err != nil {
		otcLog.WithError(err).Error("OTC NO-GO")
		preflightResult = false
		return
	}
	otcLog.Info("OTC GO")
	return
}

func checkED() {
	dashLog := preflightLog.WithField("step", "ED_API_cred_check")
	dashLog.Debug("Checking Enterprise Dashboard creds...")
	_, err := enterprise.GetConsumption(1, 0)
	if err != nil {
		dashLog.WithError(err).Error("ED NO-GO")
		preflightResult = false
		return
	}
	dashLog.Info("ED GO")
	return
}
