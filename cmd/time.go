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

	"github.com/dta4/ot3c/config"
	"github.com/spf13/cobra"
)

// timeCmd represents the time command
var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "Displays current time and date configuration",
	Long:  `Displays current time, date set for billing, billing period start, billing period end and time left`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Configured Billing Time: %v\nCurrent period start: %v\nCurrent period end: %v\nTime left: %v\n", config.BillingBegin.String(), config.CurrentBillingStart().String(), config.CurrentBillingEnd().String(), config.CalcTimeLeftInBillPeriod().String())
	},
}

func init() {
	rootCmd.AddCommand(timeCmd)

}
