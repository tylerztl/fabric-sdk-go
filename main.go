package main

import (
	"fabric-sdk-go/server"
	_ "fabric-sdk-go/server/helpers"
	"fmt"
	"log"
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
				log.Printf("Recover error : %v", err)
			}
		}()

		err := server.Run()
		log.Printf("server run error : %v", err)
	},
}

func init() {
	serverCmd.Flags().StringVarP(&server.ServerPort, "port", "p", "50052", "server port")
	serverCmd.Flags().StringVarP(&server.SwaggerDir, "swagger-dir", "", "swagger", "path to the directory which contains swagger definitions")

	rootCmd.AddCommand(serverCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func main() {
	Execute()
}
