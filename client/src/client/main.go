package main

import (
	"fal"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// export GOPATH="/mnt/d/workspace/common1:/mnt/d/workspace/dependency1:/mnt/d/workspace/client"

func main() {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "client is fal manager",
	}
	fal.BuildFALCmd(cmd)
	BuildCompletionCmd(cmd)
	cmd.Execute()
}

func BuildCompletionCmd(root *cobra.Command) {

	var completionCmd = &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long: `To load completion run

. <(bitbucket completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(bitbucket completion)
`,
		Run: func(cmd *cobra.Command, args []string) {
			f, err := os.Create("/etc/bash_completion.d/client.sh")
			if err != nil {
				fmt.Println(err)
				return
			}
			root.GenBashCompletion(f)
			fmt.Println("Please run 'source /etc/bash_completion.d/client.sh&&source /etc/bash_completion' or reopen shell")
		},
	}
	root.AddCommand(completionCmd)
}
