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

	"github.com/dta4/ot3c/enterprise"
	"github.com/spf13/cobra"
)

// getIDsCmd represents the getIDs command
var getIDsCmd = &cobra.Command{
	Use:   "dumpData",
	Short: "Dumps all cost data into a .csv file",
	Run: func(cmd *cobra.Command, args []string) {
		printAllIDs()
	},
}

func init() {
	rootCmd.AddCommand(getIDsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getIDsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getIDsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func printAllIDs() {
	var idDesc map[string]string
	idDesc = make(map[string]string)
	cons, _ := enterprise.GetAllConsumptions()
	for _, as := range cons {
		idDesc[as.Product] = as.ProductDescription
	}
	var list string = ""
	for _, asdf := range cons {
		list = list + fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,\n", asdf.Kind, asdf.Contract, asdf.BusinessPartnerID, asdf.ResellerID, asdf.QuantityUnit, asdf.ConsumptionDate, asdf.UnitPrice, asdf.ListpriceAmount, asdf.Amount, asdf.Product, asdf.ProductDescription, asdf.ProjectID, asdf.ResourceID, asdf.BillingQuantity, asdf.ProjectName)
	}
	ioutil.WriteFile("output.csv", []byte(list), 0777)
}
