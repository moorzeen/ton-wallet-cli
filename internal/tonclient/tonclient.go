package tonclient

import (
	"context"
	"fmt"

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
)

func New() (context.Context, ton.APIClientWrapped, error) {
	fmt.Print("Connecting to liteserver... ")
	defer fmt.Print("Ok\n")

	ctx := context.Background()
	client := liteclient.NewConnectionPool()

	cfg, err := liteclient.GetConfigFromUrl(ctx, "https://ton.org/global.config.json")
	if err != nil {
		return nil, nil, fmt.Errorf("get config err: %w", err)
	}

	err = client.AddConnectionsFromConfig(ctx, cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("connection err: %w", err)
	}

	apiClient := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	apiClient.SetTrustedBlockFromConfig(cfg)

	ctx = client.StickyContext(ctx)

	return ctx, apiClient, nil
}
