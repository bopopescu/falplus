package fal

import (
	"api/igm"
	"api/ipm"
	"card"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"iclient"
	"sort"
	"strconv"
	"strings"
)

var pmAddr = ":12588"

func buildPlayerCreateCmd(parent *cobra.Command) {
	var addr string
	var pid string
	var name string
	var pwd string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a player",
		Run: func(cmd *cobra.Command, args []string) {
			req := &ipm.PlayerCreateRequest{
				Name:     name,
				Pid:      pid,
				Password: pwd,
			}
			defaultPMRequest(addr, func(c *iclient.PMClient) (interface{}, error) {
				return c.PlayerCreate(context.Background(), req)
			})
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&pid, "id", "", "", "")
	flags.StringVarP(&name, "name", "n", "", "")
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
			req := &ipm.PlayerDeleteRequest{
				Pid: pid,
			}
			defaultPMRequest(addr, func(c *iclient.PMClient) (interface{}, error) {
				return c.PlayerDelete(context.Background(), req)
			})
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
			req := &ipm.PlayerListRequest{
				Pid: pid,
			}
			defaultPMRequest(addr, func(c *iclient.PMClient) (interface{}, error) {
				return c.PlayerList(context.Background(), req)
			})
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&pid, "pid", "", "", "")
	parent.AddCommand(cmd)
}

func buildPlayerSignIn(parent *cobra.Command) {
	var addr string
	var pid string
	var name string
	var pwd string
	cmd := &cobra.Command{
		Use:   "signin",
		Short: "sign in player",
		Run: func(cmd *cobra.Command, args []string) {
			req := &ipm.PlayerSignInRequest{
				Pid:      pid,
				Name:     name,
				Password: pwd,
			}
			defaultPMRequest(addr, func(c *iclient.PMClient) (interface{}, error) {
				return c.PlayerSignIn(context.Background(), req)
			})
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&pid, "pid", "", "", "")
	flags.StringVarP(&name, "name", "n", "", "")
	flags.StringVarP(&pwd, "password", "", "", "")
	parent.AddCommand(cmd)
}

func buildPlayerSignOut(parent *cobra.Command) {
	var addr string
	var pid string
	var etag string
	cmd := &cobra.Command{
		Use:   "signout",
		Short: "sign out player",
		Run: func(cmd *cobra.Command, args []string) {
			req := &ipm.PlayerSignOutRequest{
				Pid:  pid,
				Etag: etag,
			}
			defaultPMRequest(addr, func(c *iclient.PMClient) (interface{}, error) {
				return c.PlayerSignOut(context.Background(), req)
			})
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&pid, "pid", "", "", "")
	flags.StringVarP(&etag, "etag", "", "", "")
	parent.AddCommand(cmd)
}

func buildPlayerAttach(parent *cobra.Command) {
	var addr string
	var pid string
	var etag string
	var gAddr string
	cmd := &cobra.Command{
		Use:   "attach",
		Short: "player attach game",
		Run: func(cmd *cobra.Command, args []string) {
			req := &ipm.AttachRequest{
				Pid:      pid,
				Etag:     etag,
				GamePort: gmAddr,
			}
			defaultPlayerRequest(addr, func(c *iclient.PlayerClient) (interface{}, error) {
				return c.Attach(context.Background(), req)
			})
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&pid, "pid", "", "", "")
	flags.StringVarP(&etag, "etag", "", "", "")
	flags.StringVarP(&gAddr, "gaddr", "", "", "")
	parent.AddCommand(cmd)
}

func buildPlayerDetach(parent *cobra.Command) {
	var addr string
	var etag string
	cmd := &cobra.Command{
		Use:   "detach",
		Short: "player detach game",
		Run: func(cmd *cobra.Command, args []string) {
			req := &ipm.DetachRequest{
				Etag: etag,
			}
			defaultPlayerRequest(addr, func(c *iclient.PlayerClient) (interface{}, error) {
				return c.Detach(context.Background(), req)
			})
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&etag, "etag", "", "", "")
	parent.AddCommand(cmd)
}

func buildPlayerGetMsg(parent *cobra.Command) {
	var addr string
	var etag string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get message from game",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := iclient.NewPlayerClient(addr)
			if err != nil {
				fmt.Println(err)
				return
			}
			req := &ipm.GetMessageRequest{
				Etag: etag,
			}

			resp, err := c.GetMessage(context.Background(), req)
			if err != nil {
				fmt.Println(err)
				return
			}
			if resp.Status.Code != 0 {
				fmt.Println(resp.Status)
				return
			}

			fmt.Printf("message type:%s\n", igm.MsgType[resp.Gmsg.MsgType])
			var cards []card.Card
			for _, req := range resp.Gmsg.YourCards {
				cards = append(cards, card.Cards[req])
			}
			sort.Slice(cards, func(i, j int) bool {
				return cards[i].GetCardValue() > cards[j].GetCardValue()
			})
			fmt.Println("your cards:")
			for _, c := range cards {
				fmt.Printf("%d-%s-%d\t", c.CardSeq, c.CardType, c.CardNumber)
			}
			fmt.Println()

			var lastCards []card.Card
			for _, req := range resp.Gmsg.LastCards {
				lastCards = append(lastCards, card.Cards[req])
			}
			sort.Slice(lastCards, func(i, j int) bool {
				return lastCards[i].GetCardValue() > lastCards[j].GetCardValue()
			})
			fmt.Printf("Last round owner %s cards:\n", resp.Gmsg.LastId)
			for _, c := range lastCards {
				fmt.Printf("%d-%s-%d\t", c.CardSeq, c.CardType, c.CardNumber)
			}
			fmt.Println()

			fmt.Printf("player:%s is round owner\n", resp.Gmsg.RoundOwner)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&etag, "etag", "", "", "")
	parent.AddCommand(cmd)
}

func buildPlayerPutMsg(parent *cobra.Command) {
	var addr string
	var etag string
	var msgtype int64
	var seqs string
	cmd := &cobra.Command{
		Use:   "put",
		Short: "put message to game",
		Run: func(cmd *cobra.Command, args []string) {
			msg := &igm.PlayerMessage{
				MsgType: msgtype,
			}
			if seqs != "" {
				cards := strings.Split(seqs, ",")
				for _, c := range cards {
					seq, err := strconv.ParseInt(c, 10, 64)
					if err != nil {
						fmt.Println(err)
						return
					}
					msg.PutCards = append(msg.PutCards, seq)
				}
			}
			req := &ipm.PutMessageRequest{
				Etag: etag,
				Pmsg: msg,
			}
			defaultPlayerRequest(addr, func(c *iclient.PlayerClient) (interface{}, error) {
				return c.PutMessage(context.Background(), req)
			})
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "", "", "")
	flags.StringVarP(&etag, "etag", "", "", "")
	flags.StringVarP(&seqs, "seqs", "", "", "")
	flags.Int64VarP(&msgtype, "type", "", 0, "0:msg 1:get 2:put 3:pass")

	parent.AddCommand(cmd)
}
