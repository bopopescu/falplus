package fal

import (
	"fdb"
	"fmt"
	"github.com/spf13/cobra"
)

func buildBucketCmd(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:"bucket",
		Short:"bucket operate",
	}
	buildBucketListCmd(cmd)
	parent.AddCommand(cmd)
}

func buildBucketListCmd(parent *cobra.Command) {
	var file string
	cmd := &cobra.Command{
		Use:"list",
		Short:"list buckets",
		Run: func(cmd *cobra.Command, args []string) {
			db := fdb.NewDB(file)
			buckets, err := db.GetAllBucket()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(buckets)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&file, "file", "f", "/var/lib/fal/gm.db", "")
	parent.AddCommand(cmd)
}

func buildKeyCmd(parent *cobra.Command) {
	cmd := &cobra.Command{
		Use:"key",
		Short:"key operate",
	}
	buildKVList(cmd)
	parent.AddCommand(cmd)
}

func buildKVList(parent *cobra.Command) {
	var file string
	var bucket string
	cmd := &cobra.Command{
		Use:"list",
		Short:"list kv",
		Run: func(cmd *cobra.Command, args []string) {
			db := fdb.NewDB(file)
			buckets, err := db.GetAllKV(bucket)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(buckets)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&file, "file", "f", "/var/lib/fal/gm.db", "")
	flags.StringVarP(&bucket, "bucket", "b", "", "")
	parent.AddCommand(cmd)
}