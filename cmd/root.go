/*
Copyright © 2020 Manuel Plonski

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
	"os"
	"strings"

	"github.com/dta4/ot3c/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{

	Use:   "ot3c",
	Short: "Open Telekom Cloud Cost Control",
	Long: `OT3C is a tool for automatic control of cloud resources based on cost and budgeting.

	Using the OTC Enterprise Dashboard and the OTC APIs, OT3C analyzes cloud resources based on cost to determine if a specific budget is going to be exceeded. If the predicted cost are exceeding the budget, the tool is able to:

	- create a priority list of resources that are more important than others
	- alert the owners of resources regarding lower priorities and exceeding budgets (WIP)
	- terminate resources with lower priorities to stay under the budget
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	//Load config

	initConfig()
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ot3c.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ot3c" (without extension).
		viper.AddConfigPath(home)
		// Get current working directory
		pwd, _ := os.Getwd()
		viper.AddConfigPath(pwd)
		viper.SetConfigName("ot3c")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		err = parseConfig()

		if err != nil {
			logrus.WithError(err).Fatalln("Error while parsing config. Take a close look to the formating.")
		}
	} else {
		errmsg := err.Error()
		logrus.WithError(err).Warning("Unable to read config.")
		if strings.Contains(errmsg, "Config File \"ot3c\" Not Found") {
			bytes, _ := yaml.Marshal(config.SampleConfig)

			ioutil.WriteFile("ot3c-sample.yaml", bytes, os.FileMode(066))
			logrus.Errorf("No config file has been found. A sample config file named \"ot3c-sample.yaml\" has been generated. Rename the file to \"ot3c.yaml\" and place it ether next to the executable or in your home dir.")
			os.Exit(1)
		} else {
			logrus.WithError(err).Fatalln("Unable to read config.")
			os.Exit(1)
		}

	}
}

func parseConfig() error {
	currentConfig := &config.File{}
	err := viper.UnmarshalExact(currentConfig)

	if err != nil {
		bytes, err := ioutil.ReadFile(viper.ConfigFileUsed())
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(bytes, currentConfig)
		if err != nil {
			return err
		}

	}

	return currentConfig.ApplyConfig()
}
