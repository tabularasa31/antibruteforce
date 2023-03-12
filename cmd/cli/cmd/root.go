package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	proto "github.com/tabularasa31/antibruteforce/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var rootCmd = &cobra.Command{
	Use:   "abf <command> <params>",
	Short: "CLI interface for check bruteforce",
	Long: `CLI interface for check bruteforce by checking login, password and ip. 
White and blacklists included`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Use antibruteforce [command]\nRun 'abf --help' for usage.\n")
	},
}

var (
	client proto.AntiBruteforceClient
	ctx    context.Context
)

func Execute() {
	grpcHost := os.Getenv("GRPC_HOST")
	if grpcHost == "" {
		grpcHost = "localhost:50051"
	}
	conn, err := grpc.Dial(grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Error while starting dial connection: %v", err)
	}

	client = proto.NewAntiBruteforceClient(conn)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)

	err = rootCmd.Execute()
	if err != nil {
		cancel()
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Request(ctx context.Context, a proto.AntiBruteforceClient, login, pass, ip string) (ok, mess string, err error) {
	req := proto.Request{
		Login: login,
		Pass:  pass,
		Ip:    ip,
	}
	resp, err := a.AllowRequest(ctx, &req)
	if err != nil {
		log.Printf("cli - Request - err: %v", err)
		return "", "", err
	}

	return resp.GetOk().String(), resp.GetMessage(), nil
}

func Clear(ctx context.Context, a proto.AntiBruteforceClient, login, ip string) (string, error) {
	req := proto.Request{
		Login: login,
		Ip:    ip,
	}
	resp, err := a.ClearBucket(ctx, &req)
	if err != nil {
		log.Printf("cli - Clear - err: %v", err)
		return "", err
	}

	return resp.GetOk().String(), nil
}

func AddBlacklist(ctx context.Context, a proto.AntiBruteforceClient, subnet string) (ok, mess string, err error) {
	req := proto.Subnet{
		Subnet: subnet,
	}
	resp, err := a.AddToBlackList(ctx, &req)
	if err != nil {
		log.Printf("cli - AddBlacklist - err: %v", err)
		return "", "", err
	}

	return resp.GetOk().String(), resp.GetMessage(), nil
}

func AddWhitelist(ctx context.Context, a proto.AntiBruteforceClient, subnet string) (ok, mess string, err error) {
	req := proto.Subnet{
		Subnet: subnet,
	}
	resp, err := a.AddToWhiteList(ctx, &req)
	if err != nil {
		log.Printf("cli - AddWhitelist - err: %v", err)
		return "", "", err
	}

	return resp.GetOk().String(), resp.GetMessage(), nil
}

func DelBlacklist(ctx context.Context, a proto.AntiBruteforceClient, subnet string) (ok, mess string, err error) {
	req := proto.Subnet{
		Subnet: subnet,
	}
	resp, err := a.RemoveFromBlackList(ctx, &req)
	if err != nil {
		log.Printf("cli - DelBlacklist - err: %v", err)
		return "", "", err
	}

	return resp.GetOk().String(), resp.GetMessage(), nil
}

func DelWhitelist(ctx context.Context, a proto.AntiBruteforceClient, subnet string) (ok, mess string, err error) {
	req := proto.Subnet{
		Subnet: subnet,
	}
	resp, err := a.RemoveFromWhiteList(ctx, &req)
	if err != nil {
		log.Printf("cli - DelWhitelist - err: %v", err)
		return "", "", err
	}

	return resp.GetOk().String(), resp.GetMessage(), nil
}
