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
	"path/filepath"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run FILE1.x",
	Short: "Runs a file",
	Long:  `Runs a file given. If it's not compiled, compiles it and then runs it.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkAndRun(args[0])
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

// checkAndRuns checks for errors and runs the file if none.
func checkAndRun(file string) {
	// make sure the file has an extension
	if hasExtension(file) {
		filePath, rootPath, err := searchFile(file)
		if err != nil {
			panic(err)
		}

		if filePath == "" {
			fmt.Printf("Couldn't find source file %s. Please try again.\n", file)
			return
		}

		os.Chdir(rootPath) // run on the root dir of the file

		toRun := filepath.Base(filePath) // get only the name of the file
		if err := run(toRun); err != nil {
			fmt.Printf("Could not run file %s\n", filePath)
			panic(err)
		}
	} else {
		fmt.Println("The file must have an extension. Example: lazy run myproject.c")
		return
	}
}

// run runs the file. If the file isn't compiled, compiles it and then runs it.
func run(file string) error {
	outName := getOutputName(file)

	compiled, outPath := isAlreadyCompiled(file)
	if !compiled {
		compile(file, outName)
		outRoot := filepath.Dir(outPath)
		// wait until it finishes compiling
		for {
			if _, err := os.Stat(outPath); !os.IsNotExist(err) {
				// file created
				break
			}
		}

		fmt.Printf("File %s compiled. Output located in %s\n", file, outRoot)
	}

	// run the file
	cmd := exec.Command("./" + outName + ".o")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Run exec failed: %s\n", err)
		return err
	}

	return nil
}

// isAlreadyCompiled checks if the given file is already compiled by checking if exists .o files in the current directory.
//
// Additionally returns the output path (the path where it will be located if not compiled, or the actual path of the output).
func isAlreadyCompiled(file string) (bool, string) {
	out := getOutputName(file) + ".o"
	cwd, err := os.Getwd() // check in the current directory of calling
	if err != nil {
		panic(err)
	}

	outPath := filepath.Join(cwd, out)

	if _, err := os.Stat(outPath); !os.IsNotExist(err) {
		// out exists
		return true, outPath
	}

	return false, outPath
}
