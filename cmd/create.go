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
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	op "github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create FILE1 [FILE2...]",
	Short: "Create a file",
	Long: `Create a file in a new directory with the name of the extension, or adding to it if it was already created. 
If both flags -o and -t are given, the operating system will open the file with the OS preferred application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			for _, f := range args {
				// make sure the file has an extension
				if hasExtension(f) {
					err := createFile(f, false, false)
					if err != nil {
						panic(err)
					}
				} else {
					fmt.Printf("File %s doesn't have an extension. Example of use: lazy create -o myproject.go", f)
					return
				}
			}
		} else {
			/* returns the value of the given flag */
			getFlag := func(f string) bool {
				flag, err := cmd.Flags().GetBool(f)
				if err != nil {
					panic(err)
				}

				return flag
			}

			// check for flags -o and -t
			open := getFlag("open")
			terminal := getFlag("open-in-terminal")

			if hasExtension(args[0]) {
				err := createFile(args[0], open, terminal)
				if err != nil {
					panic(err)
				}
			} else {
				fmt.Println("The file must have an extension. Example: lazy create -o myproject.go")
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().BoolP("open", "o", false, "open the file after creating it, with the OS preferred application")
	createCmd.Flags().BoolP("open-in-terminal", "t", false, "open the file after creating it, on the current terminal")
}

/* Creates a file and opens it if the flag is true */
func createFile(name string, open bool, withTerminal bool) error {
	dir := createDir(name)
	file := dir + name // append the name of the file to the directory

	err := os.WriteFile(file, nil, 0644)
	if err != nil {
		return err
	}

	// if flag -o or -t, open the file
	if open {
		// run with the os preferred app
		op.Run(file)
	} else if withTerminal || (withTerminal && runtime.GOOS == "windows") {
		// set variables for bash script
		os.Setenv("PROYECT_PATH", dir)
		os.Setenv("PROYECT", name)

		// get root directory of the project
		basepath := getBasePath()
		// execute script
		cmd := exec.Command("/bin/bash", basepath+"/scripts/open_dir.sh")
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

	usr, _ := user.Current()
	dir := usr.HomeDir

	var path string
	if runtime.GOOS == "windows" {
		path = "H:\\code\\" + name[dot+1:] + "_projects\\"
	} else {
		// get the documents path and append the new name to it
		path = filepath.Join(dir, "Documents", name[dot+1:]+"_projects")
	}

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
	index := dot

	hasMoreThanOneDot := strings.Count(filepath, string(filepath[dot])) > 1
	if hasMoreThanOneDot {
		// get the index of the other dot and add to the total index count
		dot = getExtensionIndex(filepath[dot+1:])
		index += dot
	}

	return index
}

/* Get root directory of the file */
func getBasePath() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Join(filepath.Dir(b), "..")

	return basepath
}

/* returns if the file has an extension or not */
func hasExtension(f string) bool {
	if err := getExtensionIndex(f); err != -1 {
		return true
	}

	return false
}
