package main

import (
	"fabric-sdk-go/server"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run the gRPC fabric-sdk server",
}

var serverCmd = &cobra.Command{
	Use:   "start",
	Short: "Run the gRPC fabric-sdk server",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Recover error : %v", err)
			}
		}()

		err := server.Run()
		fmt.Printf("server run error : %v", err)
	},
}

func init() {
	serverCmd.Flags().StringVarP(&server.ServerPort, "port", "p", "8080", "server port")

	rootCmd.AddCommand(serverCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
