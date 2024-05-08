package account

import (
	"fmt"
	"sort"

	"github.com/moorzeen/ton-wallet-cli/internal/tonclient"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
)

var (
	statePrint         = "\nIs active: %v, status: %s\nBalance: %s TON\n%s\n"
	stateInactivePrint = "\nAccount is not active, balance: 0 TON\n"
)

type State struct {
	IsActive     bool
	Status       string
	Balance      tlb.Coins
	Transactions []*tlb.Transaction
}

func GetState(addr string, testnet bool) (*State, error) {
	a, err := address.ParseAddr(addr)
	if err != nil {
		return nil, fmt.Errorf("parse address error: %w", err)
	}

	ctx, client, err := tonclient.New(testnet)
	if err != nil {
		return nil, fmt.Errorf("new client error: %w", err)
	}

	fmt.Print("Fetching and verifying TON blockchain data... ")

	block, err := client.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("get masterchain info error: %w", err)
	}

	res, err := client.WaitForBlock(block.SeqNo).GetAccount(ctx, block, a)
	if err != nil {
		return nil, fmt.Errorf("get account error: %w", err)
	}

	fmt.Print("Ok\n")

	if !res.IsActive {
		return &State{IsActive: res.IsActive}, nil
	}

	lastHash := res.LastTxHash
	lastLt := res.LastTxLT

	fmt.Print("Getting transaction history... ")

	list, err := client.ListTransactions(ctx, a, 5, lastLt, lastHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of transactions: %w", err)
	}

	fmt.Print("Ok\n")

	lastHash = list[0].PrevTxHash
	lastLt = list[0].PrevTxLT

	sort.Slice(list, func(i, j int) bool {
		return list[i].LT > list[j].LT
	})

	txs := make([]*tlb.Transaction, 0)
	for _, t := range list {
		txs = append(txs, t)
	}

	return &State{
		IsActive:     res.IsActive,
		Status:       string(res.State.Status),
		Balance:      res.State.Balance,
		Transactions: txs,
	}, nil
}

func (s *State) Print() {
	if !s.IsActive {
		fmt.Print(stateInactivePrint)
		return
	}

	txs := "no transactions"
	if len(s.Transactions) > 0 {
		txs = "Last 5 transactions:"

		for _, tx := range s.Transactions {
			txs += "\n" + tx.String()
		}
	}

	fmt.Printf(statePrint, s.IsActive, s.Status, s.Balance.String(), txs)
}
