/*
	Copyright © 2022 Funtoo Macaroni OS Linux
	See AUTHORS and LICENSE for the license details and contributors.
*/
package cmd

import (
	. "github.com/geaaru/entities/pkg/entities"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete an entity",
	Args:  cobra.MinimumNArgs(1),
	Long:  `Deletes a entity to your system from a yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := &Parser{}

		entity, err := p.ReadEntity(args[0])
		if err != nil {
			return err
		}

		return entity.Delete(entityFile)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
