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
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
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
					path, err := searchFileAndCompile(f)
					if err != nil {
						panic(err)
					}

					if path != "" {
						fmt.Printf("File %s compiled. Output located in %s\n", args[0], path)
					} else {
						fmt.Printf("Couldn't find source file %s. Please try again.\n", args[0])
					}
				} else {
					fmt.Printf("File %s doesn't have an extension. Example of use: lazy compile myproject.c", f)
					return
				}
			}
		} else {
			// make sure the file has an extension
			if hasExtension(args[0]) {
				path, err := searchFileAndCompile(args[0])
				if err != nil {
					panic(err)
				}

				if path != "" {
					fmt.Printf("File %s compiled. Output located in %s\n", args[0], path)
				} else {
					fmt.Printf("Couldn't find source file %s. Please try again.\n", args[0])
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

/* Searchs for a file and compiles it if exists */
func searchFileAndCompile(file string) (string, error) {
	dir := getDir(file)

	var outputPath string
	// go for every file and subdirectory of the root file dir
	e := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode().IsRegular() && info.Name() == file {
			// file exists
			outputPath = filepath.Dir(path)
			outName := getOutputName(file)
			os.Chdir(outputPath) // change dir to compile in the project dir

			if err := compile(path, outName); err != nil {
				fmt.Printf("Could not compile file %s\n", path)
				return err
			}

			return nil
		}

		return nil
	})

	if e != nil {
		panic(e)
	}

	_, err := filepath.Abs(file)
	return outputPath, err
}

/* Compiles the file given */
func compile(file string, out string) error {
	lang := detectLanguage(file)

	// runs the compiler on terminal
	var cmd *exec.Cmd
	switch lang {
	case "C++":
		cmd = exec.Command("g++", file, "-o", out+".o")
	case "Java":
		cmd = exec.Command("javac", file)
	case "C":
		cmd = exec.Command("gcc", file, "-o", out+".o")
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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

	// get extension of the file
	dot := getExtensionIndex(file)
	ext := file[dot:]

	// check if the extension is on the map and return the language
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
		dot := getExtensionIndex(file) // get until dot

		// ejercicio1 --> ejer1
		return file[:4] + string(file[9:dot])
	} else if strings.Contains(file, "_") {
		return name(file, "_")
	}

	out := name(file, ".")

	return out
}
