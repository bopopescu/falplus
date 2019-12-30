package main

import (
	"api/igm"
	"api/ipm"
	"fmt"
	"golang.org/x/net/context"
	"iclient"
	"net"
)

func create() {
	gmc, err := iclient.NewGMClient(":12587")
	if err != nil {
		panic(err)
	}
	defer gmc.Close()

	pmc, err := iclient.NewPMClient(":12588")
	if err != nil {
		panic(err)
	}
	defer pmc.Close()

	ctx := context.Background()
	gameCreateReq := &igm.GameCreateRequest{GameType: 1}
	gameCreateResp, err := gmc.GameCreate(ctx, gameCreateReq)
	if err != nil || gameCreateResp.Status.Code != 0 {
		panic("create game error")
	}

	playerCreateReq := &ipm.PlayerCreateRequest{Name: "zhang", Password: "123"}
	playerCreateResp, err := pmc.PlayerCreate(ctx, playerCreateReq)
	if err != nil || playerCreateResp.Status.Code != 0 {
		panic("create player error")
	}
	playerSignInReq := &ipm.PlayerSignInRequest{Name: "zhang", Password: "123", Pid: playerCreateResp.Player.Id}
	playerSignInResp, err := pmc.PlayerSignIn(ctx, playerSignInReq)
	if err != nil || playerSignInResp.Status.Code != 0 {
		panic("sign in player error")
	}

	playerCreateReq.Name = "jia"
	playerCreateResp, err = pmc.PlayerCreate(ctx, playerCreateReq)
	if err != nil || playerCreateResp.Status.Code != 0 {
		panic("create player error")
	}
	playerSignInReq.Name = "jia"
	playerSignInReq.Pid = playerCreateResp.Player.Id
	playerSignInResp, err = pmc.PlayerSignIn(ctx, playerSignInReq)
	if err != nil || playerSignInResp.Status.Code != 0 {
		panic("sign in player error")
	}

	playerCreateReq.Name = "hua"
	playerCreateResp, err = pmc.PlayerCreate(ctx, playerCreateReq)
	if err != nil || playerCreateResp.Status.Code != 0 {
		panic("create player error")
	}
	playerSignInReq.Name = "hua"
	playerSignInReq.Pid = playerCreateResp.Player.Id
	playerSignInResp, err = pmc.PlayerSignIn(ctx, playerSignInReq)
	if err != nil || playerSignInResp.Status.Code != 0 {
		panic("sign in player error")
	}

}

func conn() {
	gmc, err := iclient.NewGMClient(":12587")
	if err != nil {
		panic(err)
	}
	defer gmc.Close()

	pmc, err := iclient.NewPMClient(":12588")
	if err != nil {
		panic(err)
	}
	defer pmc.Close()

	ctx := context.Background()
	gameList, err := gmc.GameList(ctx, &igm.GameListRequest{})
	if err != nil || gameList.Status.Code != 0 {
		panic("sign in player error")
	}

	playerList, err := pmc.PlayerList(ctx, &ipm.PlayerListRequest{})
	if err != nil || playerList.Status.Code != 0 {
		panic("sign in player error")
	}

	for _, p := range playerList.Players {
		pc, err := iclient.NewPlayerClient(net.JoinHostPort("", fmt.Sprint(p.Port)))
		if err != nil {
			panic(err)
		}

		addResp, err := gmc.GameAddPlayer(ctx, &igm.AddPlayerRequest{Gid: gameList.Games[0].Gid, Pid: p.Id})
		if err != nil || addResp.Status.Code != 0 {
			panic("GameAddPlayer error")
		}

		attachResp, err := pc.Attach(ctx, &ipm.AttachRequest{Etag: p.Etag, Pid: p.Id, GamePort: addResp.GameAddr})
		if err != nil || attachResp.Status.Code != 0 {
			panic("Attach error")
		}

		pc.Close()
	}

}

func main() {
	create()
	conn()
}
