package fal

import (
	"api/igm"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"util"
)

var gmAddr = ":12587"

func buildGameCreateCmd(parent *cobra.Command){
	var addr string
	var gtype int64
	var gid string
	var port int64
	cmd := &cobra.Command{
		Use:"create",
		Short:"create a game",
		Run: func(cmd *cobra.Command, args []string) {
			if addr == "" {
				addr = gmAddr
			}
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()
			c := igm.NewGMClient(conn)
			req := &igm.GameCreateRequest{
				GameType:gtype,
				Gid:gid,
				Port:port,
			}
			resp, err := c.GameCreate(context.Background(), req)
			if err != nil {
				fmt.Println(err)
				return
			}
			util.PrintStructObject(resp)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&gid, "gid", "", "", "")
	flags.Int64VarP(&gtype, "type", "", 0, "")
	flags.Int64VarP(&port, "port", "", 0, "")
	parent.AddCommand(cmd)
}

func buildGameDeleteCmd(parent *cobra.Command){
	var addr string
	var gid string
	cmd := &cobra.Command{
		Use:"delete",
		Short:"delete a game",
		Run: func(cmd *cobra.Command, args []string) {
			if addr == "" {
				addr = gmAddr
			}
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()
			c := igm.NewGMClient(conn)
			req := &igm.GameDeleteRequest{
				Gid:gid,
			}
			resp, err := c.GameDelete(context.Background(), req)
			if err != nil {
				fmt.Println(err)
				return
			}
			util.PrintStructObject(resp)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&gid, "gid", "", "", "")
	parent.AddCommand(cmd)
}

func buildGameListCmd(parent *cobra.Command){
	var addr string
	var gid string
	cmd := &cobra.Command{
		Use:"list",
		Short:"list game",
		Run: func(cmd *cobra.Command, args []string) {
			if addr == "" {
				addr = gmAddr
			}
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()
			c := igm.NewGMClient(conn)
			req := &igm.GameListRequest{
				Gid:gid,
			}
			resp, err := c.GameList(context.Background(), req)
			if err != nil {
				fmt.Println(err)
				return
			}
			util.PrintStructObject(resp)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&gid, "gid", "", "", "")
	parent.AddCommand(cmd)
}

