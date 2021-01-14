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
	"fmt"
	"io/ioutil"
	"math"

	"github.com/dta4/ot3c/data"
	"github.com/dta4/ot3c/otc"
	"github.com/dta4/ot3c/planning"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log *logrus.Entry = logrus.WithField("command", "plan")
var outfile *string

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Creates a Termination Plan",
	Long:  `Creates a Termination Plan and prints it on the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		plan()
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
	outfile = planCmd.Flags().StringP("outfile", "o", "", "Creates a Termination Plan file which can be used to start an OT3C T-Plan executor")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// planCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// planCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func plan() {
	if Preflight() {
		err := otc.RunDefaultDataChain()
		if err != nil {
			log.WithError(err).Error("Error on loading data")
			return
		}
		List()
		fmt.Print("\n\n\n\n")
		budget := planning.BudgetAtEndOfMonth()
		if planning.BudgetAtEndOfMonth() <= 0 {
			toSave := math.Abs(budget)
			plan, err := planning.CreateTerminationPlan(toSave)
			if err != nil {
				log.WithError(err).Error("Error on Termination Plan creation")
				return
			}
			data.SetTerminationItems(plan)
			data.UpdateTerminationPlan()
			fmt.Println(data.CurrentTerminationPlanToString())
			if *outfile != "" {
				str, err := data.RenderTerminationPlanToYAMLFile()
				if err != nil {
					log.WithError(err).Error("Error during TPlan file render")
				}
				err = ioutil.WriteFile(*outfile, []byte(str), 0644)
				if err != nil {
					log.WithError(err).Error("Error during TPlan file write")
				}
			}
		} else {
			log.Info("No Termination Plan needet")
		}
	}
}
