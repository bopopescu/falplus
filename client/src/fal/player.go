package fal

import (
	"api/ipm"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"util"
)

var pmAddr = ":12588"

func buildPlayerCreateCmd(parent *cobra.Command) {
	var addr string
	var pid string
	var port int64
	var name string
	var pwd string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a player",
		Run: func(cmd *cobra.Command, args []string) {
			if addr == "" {
				addr = pmAddr
			}
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()
			c := ipm.NewPMClient(conn)
			req := &ipm.PlayerCreateRequest{
				Name:     name,
				Pid:      pid,
				Port:     port,
				Password: pwd,
			}
			resp, err := c.PlayerCreate(context.Background(), req)
			if err != nil {
				fmt.Println(err)
				return
			}
			util.PrintStructObject(resp)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&pid, "pid", "", "", "")
	flags.StringVarP(&name, "name", "n", "", "")
	flags.Int64VarP(&port, "port", "p", 0, "")
	flags.StringVarP(&pwd, "password", "", "", "")
	parent.AddCommand(cmd)
}

func buildPlayerDeleteCmd(parent *cobra.Command) {
	var addr string
	var pid string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a game",
		Run: func(cmd *cobra.Command, args []string) {
			if addr == "" {
				addr = pmAddr
			}
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()
			c := ipm.NewPMClient(conn)
			req := &ipm.PlayerDeleteRequest{
				Pid: pid,
			}
			resp, err := c.PlayerDelete(context.Background(), req)
			if err != nil {
				fmt.Println(err)
				return
			}
			util.PrintStructObject(resp)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&pid, "pid", "", "", "")
	parent.AddCommand(cmd)
}

func buildPlayerListCmd(parent *cobra.Command) {
	var addr string
	var pid string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list game",
		Run: func(cmd *cobra.Command, args []string) {
			if addr == "" {
				addr = pmAddr
			}
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()
			c := ipm.NewPMClient(conn)
			req := &ipm.PlayerListRequest{
				Pid: pid,
			}
			resp, err := c.PlayerList(context.Background(), req)
			if err != nil {
				fmt.Println(err)
				return
			}
			util.PrintStructObject(resp)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&pid, "pid", "", "", "")
	parent.AddCommand(cmd)
}
