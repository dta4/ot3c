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
	"time"

	"github.com/dta4/ot3c/data"
	"github.com/dta4/ot3c/otc"
	"github.com/dta4/ot3c/termination"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var tfile *string
var execLog *logrus.Entry = logrus.WithField("command", "exec")

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute Termination Plan from file",
	Long:  "Starts a long running process that will execute a file-based Termination Plan. When the designated file changes, the process restarts and uses the new file.",
	Run: func(cmd *cobra.Command, args []string) {
		runExec()
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.
	tfile = execCmd.Flags().StringP("file", "f", "", "Termination Plan file to use")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runExec() {
	if !Preflight() {
		execLog.Error("Preflight check failed")
	}
	if *tfile == "" {
		execLog.Error("No Termination Plan File specified")
		return
	}

	err := otc.RunDefaultDataChain()
	if err != nil {
		execLog.WithError(err).Error("Error on VR loading")
	}
	stop := make(chan bool)

	update := runFileWatcher(*tfile, &stop)
	updateNext := make(chan bool, 1)
	//Start termination routine

	//force first update
	firstTime := true
	*update <- true
	for {
		select {
		case <-*update:
			fmt.Println()
			execLog.WithField("time", time.Now().String()).Info("Termination Plan file change detected")
			//file update detected

			blob, err := ioutil.ReadFile(*tfile)
			if err != nil {
				execLog.WithError(err).Error("Error on reading file")
				return
			}
			str := string(blob)
			err = data.ParseTerminationPlanFromString(str)
			if err != nil {
				execLog.WithError(err).Error("Error on reading file")
				break
			} else {
				//stop old daemon
				if !firstTime {
					//only notify if not first run
					updateNext <- true
				}
				firstTime = false
				go termination.RunTerminationPlan(&updateNext, time.Hour)
			}

		}
	}

}

//runFileWatcher starts go routing to watch for changes
func runFileWatcher(file string, stop *chan bool) *chan bool {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	update := make(chan bool, 1)

	go func() {
		stopBool := false
		for !stopBool {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					//notify
					update <- true
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			case <-*stop:

				stopBool = true
				watcher.Close()
			}
		}
	}()

	err = watcher.Add(file)
	if err != nil {
		log.Fatal(err)
	}
	return &update
}
