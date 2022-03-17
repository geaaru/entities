// Copyright © 2020 Ettore Di Giacinto <mudler@gentoo.org>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var entityFile string

const (
	ENTITIES_VERSION = `0.6.6`
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
