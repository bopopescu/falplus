package main

import (
	"github.com/spf13/cobra"
	game_service "gm/game/service"
	gm_service "gm/service"
	player_service "pm/player/service"
	pm_service "pm/service"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "fal",
		Short: "fight against landlord",
	}
	buildGMCmd(rootCmd)
	buildPMCmd(rootCmd)
	buildGameCmd(rootCmd)
	buildPlayerCmd(rootCmd)
	rootCmd.Execute()
}

func buildGMCmd(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "gm",
		Short: "gm command",
	}
	buildGMStartCmd(cmd)
	parent.AddCommand(cmd)
}

func buildGMStartCmd(parent *cobra.Command) {
	var addr string
	var proto string
	var name string
	var conf string
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start gm daemon",
		Run: func(cmd *cobra.Command, args []string) {
			gmServer := gm_service.NewGMServer(conf, name, proto, addr)
			gmServer.Start()
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "a", ":12587", "")
	flags.StringVarP(&proto, "proto", "p", "tcp", "")
	flags.StringVarP(&name, "name", "n", "", "")
	flags.StringVarP(&conf, "conf", "c", "", "")
	parent.AddCommand(cmd)
}

func buildGameCmd(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "game",
		Short: "game command",
	}
	buildGameStartCmd(cmd)
	parent.AddCommand(cmd)
}

func buildGameStartCmd(parent *cobra.Command) {
	var addr string
	var proto string
	var name string
	var conf string
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start game daemon",
		Run: func(cmd *cobra.Command, args []string) {
			gameServer := game_service.NewGameServer(conf, name, proto, addr)
			gameServer.Run()
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "a", "", "")
	flags.StringVarP(&proto, "proto", "p", "tcp", "")
	flags.StringVarP(&name, "name", "n", "", "")
	flags.StringVarP(&conf, "conf", "c", "", "")
	parent.AddCommand(cmd)
}

func buildPMCmd(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "pm",
		Short: "pm command",
	}
	buildPMStartCmd(cmd)
	parent.AddCommand(cmd)
}

func buildPMStartCmd(parent *cobra.Command) {
	var addr string
	var proto string
	var name string
	var conf string
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start pm daemon",
		Run: func(cmd *cobra.Command, args []string) {
			pmServer := pm_service.NewPMServer(conf, name, proto, addr)
			pmServer.Start()
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "a", ":12588", "")
	flags.StringVarP(&proto, "proto", "p", "tcp", "")
	flags.StringVarP(&name, "name", "n", "", "")
	flags.StringVarP(&conf, "conf", "c", "", "")
	parent.AddCommand(cmd)
}

func buildPlayerCmd(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "player",
		Short: "player command",
	}
	buildPlayerStartCmd(cmd)
	parent.AddCommand(cmd)
}

func buildPlayerStartCmd(parent *cobra.Command) {
	var addr string
	var proto string
	var name string
	var conf string
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start player daemon",
		Run: func(cmd *cobra.Command, args []string) {
			playerServer := player_service.NewPlayerServer(conf, name, proto, addr)
			playerServer.Start()
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "a", "", "")
	flags.StringVarP(&proto, "proto", "p", "tcp", "")
	flags.StringVarP(&name, "name", "n", "", "")
	flags.StringVarP(&conf, "conf", "c", "", "")
	parent.AddCommand(cmd)
}
