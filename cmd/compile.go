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
	Use:   "compile FILE1.x [FILE2.x [FILE3.x]]...",
	Short: "Compile a file",
	Long: `Compile a given C, C++ or Java file, whether or not it is located on the current directory.
The output will be stored in the file's project folder.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			for _, f := range args {
				checkAndRunCompile(f)
			}
		} else {
			checkAndRunCompile(args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}

// checkAndRunCompile checks for errors and compile if there's none.
func checkAndRunCompile(file string) {
	// the file must have an extension
	if hasExtension(file) {
		filePath, rootPath, err := searchFile(file)
		if err != nil {
			panic(err)
		}

		if filePath == "" {
			fmt.Printf("Couldn't find source file %s. Please try again.\n", file)
			return
		}

		outName := getOutputName(file)
		os.Chdir(rootPath) // change dir to compile in the project dir

		if err := compile(filePath, outName); err != nil {
			fmt.Printf("Could not compile file %s\n", filePath)
			panic(err)
		}

		fmt.Printf("File %s compiled. Output located in %s\n", file, rootPath)
	} else {
		fmt.Println("The file must have an extension. Example: lazy compile myproject.c")
		return
	}
}

// searchFile searchs for a file and returns the path and its root path, along with an error if any.
func searchFile(file string) (string, string, error) {
	dir := getDir(file)

	var filePath string
	var rootPath string
	// go for every file and subdirectory of the root file dir
	e := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode().IsRegular() && info.Name() == file {
			// file exists
			filePath = path
			rootPath = filepath.Dir(path)

			return nil
		}

		return nil
	})

	return filePath, rootPath, e
}

// compile compiles the file given with the output name of the out parameter and returns an error on failure.
//
// compile uses GCC for C compiling, G++ for C++ compiling and Javac for Java compiling.
func compile(path string, out string) error {
	lang := detectLanguage(path)

	// run the compiler on terminal
	var cmd *exec.Cmd
	switch lang {
	case "C++":
		cmd = exec.Command("g++", path, "-o", out+".o")
	case "Java":
		cmd = exec.Command("javac", path)
	case "C":
		cmd = exec.Command("gcc", path, "-o", out+".o")
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Compiling exec failed: %s\n", err)
		return err
	}

	return nil
}

// detectLanguage detects the programming language of given file and returns the name of it.
//
// If the language is not found, it returns "".
//
// * For now this only applies for CPP, Java and C files. *
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

// getOutputName gets the .out name of the file as an abbreviation of the original.
//
// If the file string is in the form of "ejercicioX.y", it will return it in format "ejerX" (without extension).
// If the file contains an underscore, it will return a slice until that underscore.
// Else it will return a slice of the string until a dot (".").
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
