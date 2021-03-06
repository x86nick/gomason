// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/nikogura/gomason/pkg/gomason"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build your code in a clean environment.",
	Long: `
Build your code in a clean environment.

Includes 'test'.  It aint gonna build if the tests don't pass.

You could run 'test' separately, but 'build' is nice enough to do it for you.

Binaries are dropped into the current working directory.
`,
	Run: func(cmd *cobra.Command, args []string) {
		workDir, err := ioutil.TempDir("", "gomason")
		if err != nil {
			log.Fatalf("Failed to create temp dir: %s", err)
		}

		if verbose {
			log.Printf("Created temp dir %s", workDir)
		}

		defer os.RemoveAll(workDir)

		gopath, err := gomason.CreateGoPath(workDir)
		if err != nil {
			log.Fatalf("Failed to create ephemeral GOPATH: %s", err)
		}

		meta, err := gomason.ReadMetadata("metadata.json")

		err = gomason.GovendorInstall(gopath, verbose)
		if err != nil {
			log.Fatalf("Failed to install Govendor: %s", err)
		}

		if err != nil {
			log.Fatalf("couldn't read package information from metadata.json: %s", err)

		}

		err = gomason.Checkout(gopath, meta, branch, verbose)
		if err != nil {
			log.Fatalf("failed to checkout package %s at branch %s: %s", meta.Package, branch, err)
		}

		err = gomason.GovendorSync(gopath, meta, verbose)
		if err != nil {
			log.Fatalf("error running govendor sync: %s", err)
		}

		err = gomason.GoTest(gopath, meta.Package, verbose)
		if err != nil {
			log.Fatalf("error running go test: %s", err)
		}

		log.Printf("Tests Succeeded!\n\n")

		err = gomason.Build(gopath, meta, branch, verbose)
		if err != nil {
			log.Fatalf("build failed: %s", err)
		}

		log.Printf("Build Succeeded!\n\n")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
