package fal

import (
	"debugCmd"
	"fmt"
	"github.com/spf13/cobra"
	"iclient"
	"util"
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
	buildGameAddPlayerCmd(cmd)
	buildGameStartCmd(cmd)
	buildGameStopCmd(cmd)
	buildGameExitCmd(cmd)

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
	buildPlayerSignIn(cmd)
	buildPlayerSignOut(cmd)
	buildPlayerAttach(cmd)
	buildPlayerDetach(cmd)
	buildPlayerGetMsg(cmd)
	buildPlayerPutMsg(cmd)

	parent.AddCommand(cmd)
}

func defaultGMRequest(addr string, f func(c *iclient.GMClient) (interface{}, error)) {
	if addr == "" {
		addr = gmAddr
	}
	c, err := iclient.NewGMClient(addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := f(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	util.PrintStructObject(resp)
}

func defaultPMRequest(addr string, f func(c *iclient.PMClient) (interface{}, error)) {
	if addr == "" {
		addr = pmAddr
	}
	c, err := iclient.NewPMClient(addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := f(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	util.PrintStructObject(resp)
}

func defaultPlayerRequest(addr string, f func(c *iclient.PlayerClient) (interface{}, error)) {
	c, err := iclient.NewPlayerClient(addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := f(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	util.PrintStructObject(resp)
}
