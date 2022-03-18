/*
	Copyright © 2022 Funtoo Macaroni OS Linux
	See AUTHORS and LICENSE for the license details and contributors.
*/
package cmd

import (
	. "github.com/geaaru/entities/pkg/entities"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "applies an entity",
	Args:  cobra.MinimumNArgs(1),
	Long:  `Applies a entity yaml file to your system`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := &Parser{}

		safe, _ := cmd.Flags().GetBool("safe")

		entity, err := p.ReadEntity(args[0])
		if err != nil {
			return err
		}

		return entity.Apply(entityFile, safe)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	var flags = applyCmd.Flags()
	flags.Bool("safe", false,
		"Avoid to override existing entity if it has difference or if the id is used in a different way.")
}
