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
	"strings"

	"github.com/dta4/ot3c/data"
	"github.com/dta4/ot3c/otc"
	"github.com/dta4/ot3c/planning"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all resources controlled by OT3C",
	Run: func(cmd *cobra.Command, args []string) {
		listCommand()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//List prints the current resources on the terminal
func List() {

	fmt.Println()
	tree := treeprint.New()
	for _, resource := range data.Resources {
		CreateTree(tree, resource)
	}
	fmt.Println(tree.String())
	fmt.Println("-----------------------------------------------")
	current := planning.TotalCurrentCost()
	rest := planning.TotalRestCost()
	fmt.Printf("Current: %v €\n", current)
	fmt.Printf("Predicted: %v €\n", rest)
	fmt.Println("-----------------------------------------------")
	fmt.Printf("Predicted total: %v € (%v € left in budget)\n", current+rest, planning.BudgetAtEndOfMonth())

}

func listCommand() {
	if Preflight() {
		otc.RunDefaultDataChain()
		List()
	}
}

func printVR(vr data.VirtualResource) string {
	var svr strings.Builder
	svr.WriteString(fmt.Sprintf("%v\n", vr.GetID()))
	svr.WriteString(fmt.Sprintf("%v\n", vr.GetTags()))
	svr.WriteString(fmt.Sprintf("Current Cost: %v €\n", data.GetCurrentCost(vr)))
	svr.WriteString(fmt.Sprintf("Future Cost: %v €\n", data.GetCostTillBillingEnd(vr)))
	return svr.String()
}

//CreateTree decends recusivly and adds nodes to the treeprint
func CreateTree(tree treeprint.Tree, vr data.VirtualResource) {

	local := tree.AddBranch(fmt.Sprintf("%v\n", vr.GetID())).AddNode(fmt.Sprintf("%v\n", vr.GetTags())).AddNode(fmt.Sprintf("Current Cost: %v €\n", data.GetCurrentCost(vr))).AddNode(fmt.Sprintf("Future Cost: %v €\n", data.GetCostTillBillingEnd(vr)))
	gr, ok := vr.(data.GroupResource)
	if ok {
		tree = local.AddBranch("list")
		for _, vrg := range gr.List {
			CreateTree(tree, vrg)
		}
	}

}
