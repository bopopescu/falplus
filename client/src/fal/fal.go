package fal

import (
	"debugCmd"
	"github.com/spf13/cobra"
)

func BuildFALCmd(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "fal",
		Short: "manage fal system",
	}
	buildGMCmd(cmd)
	buildPMCmd(cmd)
	buildDBCmd(cmd)
	debugCmd.BuildDebugCmd(cmd, ":12587")
	parent.AddCommand(cmd)
}

func buildDBCmd(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "manage db",
	}
	buildBucketCmd(cmd)
	buildKeyCmd(cmd)
	parent.AddCommand(cmd)
}

func buildGMCmd(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "game",
		Short: "manage game",
	}
	buildGameCreateCmd(cmd)
	buildGameDeleteCmd(cmd)
	buildGameListCmd(cmd)

	parent.AddCommand(cmd)
}

func buildPMCmd(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "player",
		Short: "manage player",
	}

	buildPlayerCreateCmd(cmd)
	buildPlayerDeleteCmd(cmd)
	buildPlayerListCmd(cmd)

	parent.AddCommand(cmd)
}
