package fal

import (
	"api/igm"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"iclient"
)

var gmAddr = ":12587"

func buildGameCreateCmd(parent *cobra.Command) {
	var addr string
	var gtype int64
	var gid string
	var port int64
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a game",
		Run: func(cmd *cobra.Command, args []string) {
			req := &igm.GameCreateRequest{
				GameType: gtype,
				Gid:      gid,
				Port:     port,
			}
			defaultGMRequest(addr, func(c *iclient.GMClient) (interface{}, error) {
				return c.GameCreate(context.Background(), req)
			})
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&gid, "gid", "", "", "")
	flags.Int64VarP(&gtype, "type", "", 0, "")
	flags.Int64VarP(&port, "port", "", 0, "")
	parent.AddCommand(cmd)
}

func buildGameDeleteCmd(parent *cobra.Command) {
	var addr string
	var gid string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a game",
		Run: func(cmd *cobra.Command, args []string) {
			req := &igm.GameDeleteRequest{
				Gid: gid,
			}
			defaultGMRequest(addr, func(c *iclient.GMClient) (interface{}, error) {
				return c.GameDelete(context.Background(), req)
			})
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&gid, "gid", "", "", "")
	parent.AddCommand(cmd)
}

func buildGameListCmd(parent *cobra.Command) {
	var addr string
	var gid string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list game",
		Run: func(cmd *cobra.Command, args []string) {
			req := &igm.GameListRequest{
				Gid: gid,
			}
			defaultGMRequest(addr, func(c *iclient.GMClient) (interface{}, error) {
				return c.GameList(context.Background(), req)
			})
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&gid, "gid", "", "", "")
	parent.AddCommand(cmd)
}

func buildGameAddPlayerCmd(parent *cobra.Command) {
	var addr string
	var gid string
	var pid string
	var pAddr string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "game add player",
		Run: func(cmd *cobra.Command, args []string) {
			req := &igm.AddPlayerRequest{
				Gid: gid,
			}
			defaultGMRequest(addr, func(c *iclient.GMClient) (interface{}, error) {
				return c.GameAddPlayer(context.Background(), req)
			})
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&gid, "gid", "", "", "")
	flags.StringVarP(&pid, "pid", "", "", "")
	flags.StringVarP(&pAddr, "paddr", "", "", "")
	parent.AddCommand(cmd)
}

func buildGameStartCmd(parent *cobra.Command) {
	var addr string
	var gid string
	var pid string
	cmd := &cobra.Command{
		Use:   "start",
		Short: "game start",
		Run: func(cmd *cobra.Command, args []string) {
			req := &igm.GameStartRequest{
				Gid: gid,
				Pid: pid,
			}
			defaultGMRequest(addr, func(c *iclient.GMClient) (interface{}, error) {
				return c.GameStart(context.Background(), req)
			})
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&gid, "gid", "", "", "")
	flags.StringVarP(&pid, "pid", "", "", "")
	parent.AddCommand(cmd)
}

func buildGameStopCmd(parent *cobra.Command) {
	var addr string
	var gid string
	var pid string
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "game stop",
		Run: func(cmd *cobra.Command, args []string) {
			req := &igm.GameStopRequest{
				Gid: gid,
				Pid: pid,
			}
			defaultGMRequest(addr, func(c *iclient.GMClient) (interface{}, error) {
				return c.GameStop(context.Background(), req)
			})
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&gid, "gid", "", "", "")
	flags.StringVarP(&pid, "pid", "", "", "")
	parent.AddCommand(cmd)
}

func buildGameExitCmd(parent *cobra.Command) {
	var addr string
	var gid string
	var pid string
	cmd := &cobra.Command{
		Use:   "exit",
		Short: "game exit",
		Run: func(cmd *cobra.Command, args []string) {
			req := &igm.GameExitRequest{
				Gid:      gid,
				PlayerId: pid,
			}
			defaultGMRequest(addr, func(c *iclient.GMClient) (interface{}, error) {
				return c.GameExit(context.Background(), req)
			})
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&gid, "gid", "", "", "")
	flags.StringVarP(&pid, "pid", "", "", "")
	parent.AddCommand(cmd)
}
