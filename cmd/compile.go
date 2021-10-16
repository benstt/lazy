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

	"github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			for _, f := range args {
				// make sure the file has an extension
				if hasExtension(f) {
					err := compile(f)
					if err != nil {
						panic(err)
					}
				} else {
					fmt.Printf("File %s doesn't have an extension. Example of use: lazy compile myproject.c", f)
					return
				}
			}
		} else {
			// make sure the file has an extension
			if hasExtension(args[0]) {
				err := compile(args[0])
				if err != nil {
					panic(err)
				}
			} else {
				fmt.Println("The file must have an extension. Example: lazy compile myproject.c")
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}

/* Compiles the file given */
func compile(file string) error {
	lang := detectLanguage(file)

	// runs the compiler on terminal
	var cmd *exec.Cmd
	switch lang {
	case "C++":
		cmd = exec.Command("g++", file, "-o", getOutputName(file)+".o")
	case "Java":
		cmd = exec.Command("javac", file)
	case "C":
		cmd = exec.Command("gcc", file, "-o", getOutputName(file)+".o")
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Start failed: %s\n", err)
		return err
	}

	return nil
}

/* Detects the language of given file
For now this only applies for CPP, Java and C files
*/
func detectLanguage(file string) string {
	languages := map[string]string{
		".cpp":  "C++",
		".java": "Java",
		".c":    "C",
	}

	dot := getExtensionIndex(file)
	ext := file[dot:]

	if val, ok := languages[ext]; ok {
		return val
	}

	return ""
}

/* Gets the .out name of the file as an abbreviation of the original */
func getOutputName(file string) string {
	/* Gets the substr until the given char */
	name := func(f string, char string) string {
		i := strings.Index(f, char)
		return f[:i]
	}

	if strings.HasPrefix(file, "ejercicio") {
		// ejercicio1 --> ejer1
		return file[:4] + string(file[9:])
	} else if strings.Contains(file, "_") {
		return name(file, "_")
	}

	return file
}
