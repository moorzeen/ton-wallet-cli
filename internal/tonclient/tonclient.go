package tonclient

import (
	"context"
	"fmt"

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
)

var (
	status = "Connecting to liteserver... "
	url    = "https://ton.org/global.config.json"
)

func New(testnet bool) (context.Context, ton.APIClientWrapped, error) {
	if testnet {
		status = "Connecting to TESTNET liteserver... "
		url = "https://ton.org/testnet-global.config.json"
	}

	fmt.Print(status)
	defer fmt.Print("Ok\n")

	ctx := context.Background()
	client := liteclient.NewConnectionPool()

	cfg, err := liteclient.GetConfigFromUrl(ctx, url)
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
