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
	cancel context.CancelFunc
)

func Execute() {
	if err := execNoExit(); err != nil {
		log.Fatal(err)
	}

	if err := rootCmd.Execute(); err != nil {
		cancel()
		log.Fatal(err)
	}
}

func execNoExit() error {
	grpcHost := os.Getenv("GRPC_HOST")
	if grpcHost == "" {
		grpcHost = "localhost:50051"
	}
	conn, err := grpc.Dial(grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error while starting dial connection: %w", err)
	}

	client = proto.NewAntiBruteforceClient(conn)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return nil
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
