/*
Copyright © 2021 Benjamín García Roqués <benjamingarciaroques@gmail.com>

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
	"os"
	// "os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a file",
	Long:  `Create a file in a new directory with the name of the extension, or adding to it if it was already created.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, f := range args {
			open, err := cmd.Flags().GetBool("open")
			if err != nil {
				panic(err)
			}
			createFile(f, open)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().BoolP("open", "o", false, "open the file after creating it")
}

/* Creates a file and opens it if the flag is true */
func createFile(name string, open bool) error {
	dir := createDir(name)
	file := dir + name

	err := os.WriteFile(file, nil, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("File created at %s\n", file)

	// TODO: handle open flag
	// if open {
	// 	cmd := exec.Command("cmd", file)
	// 	err = cmd.Start()
	// 	if err != nil {
	// 		fmt.Printf("Start failed: %s", err)
	// 	}
	// 	fmt.Printf("Waiting for command to finish.\n")
	// 	err = cmd.Wait()
	// 	fmt.Printf("Command finished with error: %v\n", err)
	// }

	return nil
}

/* Create a directory with the name of the extension followed by a '_projects' */
func createDir(name string) string {
	dot := getExtensionIndex(name)
	path := "H:/code/" + name[dot+1:] + "_projects/"

	// check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			panic(err)
		}
	}

	return path
}

/* Get position of the dot in the extension */
func getExtensionIndex(filepath string) int {
	dot := strings.Index(filepath, ".")

	return dot
}
