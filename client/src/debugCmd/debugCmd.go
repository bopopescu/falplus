package debugCmd

import (
	"api/idebug"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"util"
)

func BuildDebugCmd(rootCmd *cobra.Command, addr string) {
	var debugCmd = &cobra.Command{
		Use:   "debug",
		Short: "Manage Debug",
	}

	//添加到上级命令
	rootCmd.AddCommand(debugCmd)

	//添加子命令

	buildDebugProfCmd(debugCmd, addr)
	buildDebugStatsCmd(debugCmd, addr)
	buildDebugLogCmd(debugCmd, addr)
	buildDebugStartCmd(debugCmd, addr)
	buildDebugStopCmd(debugCmd, addr)
}

func buildDebugStartCmd(parentCmd *cobra.Command, addr string) {
	var address string

	var ProfCmd = &cobra.Command{
		Use:   "start",
		Short: "start prof service",
		Example: "use Example:\n" +
			"go tool pprof http://localhost:10000/debug/pprof/profile\n" +
			"go tool pprof http://localhost:10000/debug/pprof/heap\n" +
			"go tool pprof http://localhost:10000/debug/pprof/goroutine",
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()

			c := idebug.NewDebugClient(conn)
			resp, err := c.StartPprof(context.Background(), &idebug.StartPprofRequest{Addr: address})
			if err != nil {
				fmt.Println(err)
				return
			}
			util.PrintStructObject(resp)
		},
	}

	parentCmd.AddCommand(ProfCmd)

	ProfCmd.Flags().StringVarP(&address, "listen", "l", ":10000", "the listen address")
	ProfCmd.Flags().StringVarP(&addr, "addr", "a", addr, "the service address")
}

func buildDebugStopCmd(parentCmd *cobra.Command, addr string) {
	var ProfCmd = &cobra.Command{
		Use:   "stop",
		Short: "stop prof service",
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()

			c := idebug.NewDebugClient(conn)
			resp, err := c.StopPprof(context.Background(), &empty.Empty{})
			if err != nil {
				fmt.Println(err)
				return
			}
			util.PrintStructObject(resp)
		},
	}

	parentCmd.AddCommand(ProfCmd)
	ProfCmd.Flags().StringVarP(&addr, "addr", "a", addr, "the service address")
}

func buildDebugProfCmd(parentCmd *cobra.Command, addr string) {
	var address string
	var name string
	var path string
	var time int64

	var ProfCmd = &cobra.Command{
		Use:   "prof",
		Short: "prof command ",
		Long: `	exampel :
        access storage debug prof --path=/root/default.prof
        go tool pprof -top /usr/local/bin/storage  /root/default.prof
        go tool pprof -svg /usr/local/bin/storage  /root/default.prof  > 123.svg
        /var/lib/thci/tool/bin/pprof -top /usr/local/bin/storage  /default.prof
        /var/lib/thci/tool/bin/pprof -svg /usr/local/bin/storage  /default.prof  > 123.svg
		`,
		Run: func(cmd *cobra.Command, args []string) {
			opts := []grpc.DialOption{grpc.WithInsecure()}
			conn, err := grpc.Dial(address, opts...)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()

			c := idebug.NewDebugClient(conn)
			if name != "cpu" && name != "heap" && name != "goruntine" && name != "threadcreate" && name != "block" {
				fmt.Printf("name error")
				return
			}
			resp, err := c.Prof(context.Background(), &idebug.ProfRequest{Name: name, Path: path, Time: time})
			if err != nil {
				fmt.Println(err)
				return
			}
			util.PrintStructObject(resp)
		},
	}

	parentCmd.AddCommand(ProfCmd)

	ProfCmd.Flags().StringVarP(&address, "addr", "a", addr, "")
	ProfCmd.Flags().StringVarP(&name, "name", "n", "cpu", "cpu,heap, goroutine,threadcreate,block")
	ProfCmd.Flags().StringVarP(&path, "path", "p", "default.prof", "prof file")
	ProfCmd.Flags().Int64VarP(&time, "time", "t", 20, "time")
}

func buildDebugStatsCmd(parentCmd *cobra.Command, addr string) {
	var StatsCmd = &cobra.Command{
		Use:   "stats",
		Short: "stats command",
	}

	parentCmd.AddCommand(StatsCmd)

	buildStatsMemstatsCmd(StatsCmd, addr)
	buildStatsRoutineCmd(StatsCmd, addr)
	buildStatsStackCmd(StatsCmd, addr)
}

func buildDebugLogCmd(parentCmd *cobra.Command, addr string) {
	var LogCmd = &cobra.Command{
		Use:   "log",
		Short: "set and get log level",
	}

	parentCmd.AddCommand(LogCmd)

	buildDebugGetLogLevelCmd(LogCmd, addr)
	buildDebugSetLogLevelCmd(LogCmd, addr)
}

func buildDebugGetLogLevelCmd(parentCmd *cobra.Command, addr string) {
	var address string

	var GetLogLevelCmd = &cobra.Command{
		Use:   "get",
		Short: "get log level",
		Run: func(cmd *cobra.Command, args []string) {
			opts := []grpc.DialOption{grpc.WithInsecure()}
			conn, err := grpc.Dial(address, opts...)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()

			c := idebug.NewDebugClient(conn)
			resp, err := c.GetLogLevel(context.Background(), &empty.Empty{})
			if err != nil {
				fmt.Println(err)
				return
			}
			util.PrintStructObject(resp)
		},
	}

	parentCmd.AddCommand(GetLogLevelCmd)

	GetLogLevelCmd.Flags().StringVarP(&address, "addr", "a", addr, "")
}

func buildDebugSetLogLevelCmd(parentCmd *cobra.Command, addr string) {
	var address string
	var level string

	var SetLogLevelCmd = &cobra.Command{
		Use:   "set",
		Short: "set log level",
		Run: func(cmd *cobra.Command, args []string) {
			opts := []grpc.DialOption{grpc.WithInsecure()}
			conn, err := grpc.Dial(address, opts...)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()

			c := idebug.NewDebugClient(conn)
			resp, err := c.SetLogLevel(context.Background(), &idebug.LogRequest{Level: level})
			if err != nil {
				fmt.Println(err)
				return
			}
			util.PrintStructObject(resp)
		},
	}

	parentCmd.AddCommand(SetLogLevelCmd)

	SetLogLevelCmd.Flags().StringVarP(&level, "level", "l", "", "value:[debug,info,warn,error,fatal,panic]")
	SetLogLevelCmd.Flags().StringVarP(&address, "addr", "a", addr, "")
}

func buildStatsStackCmd(parentCmd *cobra.Command, addr string) {
	var address string

	var StackCmd = &cobra.Command{
		Use:   "stack",
		Short: "stack controller stack",
		Run: func(cmd *cobra.Command, args []string) {
			opts := []grpc.DialOption{grpc.WithInsecure()}
			conn, err := grpc.Dial(address, opts...)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()

			c := idebug.NewDebugClient(conn)
			resp, err := c.Stats(context.Background(), &idebug.StatsRequest{Name: "stack"}, grpc.MaxCallRecvMsgSize(1024*1024*1024))
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(resp.Data))
		},
	}

	parentCmd.AddCommand(StackCmd)

	StackCmd.Flags().StringVarP(&address, "addr", "a", addr, "")
}

func buildStatsMemstatsCmd(parentCmd *cobra.Command, addr string) {
	var address string

	var MemstatsCmd = &cobra.Command{
		Use:   "memstats",
		Short: "memstats controller stack",
		Run: func(cmd *cobra.Command, args []string) {
			opts := []grpc.DialOption{grpc.WithInsecure()}
			conn, err := grpc.Dial(address, opts...)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()

			c := idebug.NewDebugClient(conn)
			resp, err := c.Stats(context.Background(), &idebug.StatsRequest{Name: "memstats"}, grpc.MaxCallRecvMsgSize(1024*1024*1024))
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(resp.Data))
		},
	}

	parentCmd.AddCommand(MemstatsCmd)

	MemstatsCmd.Flags().StringVarP(&address, "addr", "a", addr, "")
}

func buildStatsRoutineCmd(parentCmd *cobra.Command, addr string) {
	var address string

	var RoutineCmd = &cobra.Command{
		Use:   "routine",
		Short: "routine controller stack",
		Run: func(cmd *cobra.Command, args []string) {
			opts := []grpc.DialOption{grpc.WithInsecure()}
			conn, err := grpc.Dial(address, opts...)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()

			c := idebug.NewDebugClient(conn)
			resp, err := c.Stats(context.Background(), &idebug.StatsRequest{Name: "routine"})
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(resp.Data))
		},
	}

	parentCmd.AddCommand(RoutineCmd)

	RoutineCmd.Flags().StringVarP(&address, "addr", "a", addr, "")
}
