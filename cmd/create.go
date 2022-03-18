/*
	Copyright © 2022 Funtoo Macaroni OS Linux
	See AUTHORS and LICENSE for the license details and contributors.
*/
package cmd

import (
	. "github.com/geaaru/entities/pkg/entities"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create an entity",
	Args:  cobra.MinimumNArgs(1),
	Long:  `Create a entity to your system from yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := &Parser{}

		entity, err := p.ReadEntity(args[0])
		if err != nil {
			return err
		}

		return entity.Create(entityFile)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
