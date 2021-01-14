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
	"io/ioutil"
	"os"

	"github.com/dta4/ot3c/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// generateConfigCmd represents the generateConfig command
var generateConfigCmd = &cobra.Command{
	Use:   "generateConfig",
	Short: "Generates a sample config to fill out missing values",

	Run: func(cmd *cobra.Command, args []string) {
		bytes, _ := yaml.Marshal(config.SampleConfig)

		ioutil.WriteFile("ot3c-sample.yaml", bytes, os.FileMode(066))
	},
}

func init() {
	rootCmd.AddCommand(generateConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
