/*
	Copyright © 2022 Funtoo Macaroni OS Linux
	See AUTHORS and LICENSE for the license details and contributors.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var entityFile string

const (
	ENTITIES_VERSION = `0.9.1`
)

var (
	BuildTime      string
	BuildCommit    string
	BuildGoVersion string
)

func version() string {
	ans := fmt.Sprintf("%s-g%s %s", ENTITIES_VERSION, BuildCommit, BuildTime)
	if BuildGoVersion != "" {
		ans += " " + BuildGoVersion
	}
	return ans
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "entities",
	Version: version(),
	Short:   "Modern go identity manager for UNIX systems",
	Long: `Entities is a modern groups and user manager for Unix system. It allows to create/delete user and groups 
in a system given policies following the entities yaml format.

For example:

	$> entities apply <entity.yaml>
	$> entities delete <entity.yaml>
	$> entities create <entity.yaml>
`,
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
	rootCmd.PersistentFlags().StringVarP(&entityFile, "file", "f", "", "File to manipulate ( e.g. /etc/passwd ) ")
}
