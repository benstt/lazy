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
	"os/exec"
	"strings"
    "path/filepath"
    "runtime"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
    Use:   "create FILE1 [FILE2...]",
	Short: "Create a file",
	Long:  `Create a file in a new directory with the name of the extension, or adding to it if it was already created.`,
    Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, f := range args {
			open, err := cmd.Flags().GetBool("open")
			if err != nil {
				panic(err)
			}

            // make sure the file has an extension
            if err := getExtensionIndex(f); err != -1 {
                createFile(f, open)
            } else {
                fmt.Println("The file must have an extension. Example: lazy create -o myproject.go")
                return
            }
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
	file := dir + name // append the name of the file to the directory

	err := os.WriteFile(file, nil, 0644)
	if err != nil {
		return err
	}

    // if flag -o, open the file
	if open {
        // set variables for bash script
        os.Setenv("PROYECT_PATH", dir)
        os.Setenv("PROYECT", name)

        // get root directory of the project
        basepath := getBasePath()

        // execute bash script
        cmd := exec.Command("/bin/bash", basepath + "/scripts/open_dir.sh")
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        err := cmd.Run()
        if err != nil {
            fmt.Printf("Start failed: %s\n", err)
        }
	}

	fmt.Printf("File created at %s\n", file)

	return nil
}

/* Create a directory with the name of the extension followed by '_projects' */
func createDir(name string) string {
	dot := getExtensionIndex(name)
    basepath := getBasePath() // get root path of this file (cmd)

    // get the root path of the lazy directory and append the new name to it
	path := basepath + "/../" + name[dot+1:] + "_projects/"

	// check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			panic(err)
		}
	}
    // clean up path for simplicity
    path = filepath.Dir(path) + "/"

	return path
}

/* Get position of the dot in the extension */
func getExtensionIndex(filepath string) int {
	dot := strings.Index(filepath, ".")
	return dot
}

/* Get root directory of the file */
func getBasePath() string {
    _, b, _, _ := runtime.Caller(0)
    basepath   := filepath.Join(filepath.Dir(b), "..")

    return basepath
}
