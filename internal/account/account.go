package account

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/mdp/qrterminal/v3"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

type Account struct {
	Wallet *wallet.Wallet
	Seed   string
	Key    string
	State  *State
}

func New() (*Account, error) {
	seed := wallet.NewSeed()

	w, err := wallet.FromSeed(nil, seed, wallet.V4R2)
	if err != nil {
		return nil, fmt.Errorf("wallet from seed error: %w", err)
	}

	return &Account{
		Wallet: w,
		Seed:   strings.Join(seed, " "),
		Key:    hex.EncodeToString(w.PrivateKey().Seed()),
		State:  nil,
	}, nil
}

func (a *Account) Print() {
	addr := a.Wallet.WalletAddress().String()
	seed := a.Seed
	key := a.Key

	fmt.Println()
	qrterminal.Generate(addr, qrterminal.L, os.Stdout)
	fmt.Printf("Account address: %s\n"+
		"Seed phrase: %s\n"+
		"Private key: %s\n\n"+
		"Please, keep seed phrase and private key safe. Check wallet address when scan QR code.\n", addr, seed, key)
}
