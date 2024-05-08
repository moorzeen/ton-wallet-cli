package account

import (
	"bufio"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/moorzeen/ton-wallet-cli/internal/tonclient"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

func walletFromSecret(client ton.APIClientWrapped, secret string) (*wallet.Wallet, error) {
	parts := strings.Split(secret, " ")

	if len(parts) == 24 {
		seed := strings.Split(secret, " ")

		w, err := wallet.FromSeed(client, seed, wallet.V4R2)
		if err != nil {
			return nil, err
		}

		return w, nil
	}

	key, err := hex.DecodeString(secret)
	if err != nil || len(key) != ed25519.SeedSize {
		return nil, errors.New("invalid private key")
	}

	w, err := wallet.FromPrivateKey(client, ed25519.NewKeyFromSeed(key), wallet.V4R2)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func buildMessage(to, comment, amount string) (*wallet.Message, error) {
	receiver, err := address.ParseAddr(to)
	if err != nil {
		return nil, fmt.Errorf("parse wallet address error: %w", err)
	}

	var body *cell.Cell
	if comment != "" {
		body, err = wallet.CreateCommentCell(comment)
		if err != nil {
			return nil, fmt.Errorf("create comment error: %w", err)
		}
	}

	var (
		mode uint8
		coin tlb.Coins
	)
	if amount == "ALL" {
		mode = 128
	} else {
		mode = 1 + 2
		amount = strings.ReplaceAll(amount, " ", "")
		amount = strings.ReplaceAll(amount, ",", ".")
		coin, err = tlb.FromTON(amount)
		if err != nil {
			return nil, fmt.Errorf("convert string to coins error: %w", err)
		}
	}

	return &wallet.Message{
		Mode: mode,
		InternalMessage: &tlb.InternalMessage{
			IHRDisabled: true,
			Bounce:      false,
			DstAddr:     receiver,
			Amount:      coin,
			Body:        body,
		},
	}, nil
}

func Send(to, msg, amount, key string) error {
	ctx, client, err := tonclient.New()
	if err != nil {
		return fmt.Errorf("new client error: %w", err)
	}

	wlt, err := walletFromSecret(client, key)
	if err != nil {
		return fmt.Errorf("wallet from secret error: %w", err)
	}

	tx, err := buildMessage(to, msg, amount)
	if err != nil {
		return fmt.Errorf("build transaction error: %w", err)
	}

	var a string
	if tx.Mode == 128 {
		a = "all balance"
	} else {
		a = tx.InternalMessage.Amount.String() + " TON"
	}
	text := fmt.Sprintf("You want to send %s to %s with comment \"%s\". Are you sure? (\"yes\" or \"no\")", a, to, msg)

	for {
		fmt.Printf("%s\n%s", text, "> ")

		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		err = input.Err()
		if err != nil {
			return fmt.Errorf("scan input error: %w", err)
		}

		switch input.Text() {
		case "yes":
			break
		case "no":
			return nil
		default:
			fmt.Println("Please, type \"yes\" or \"no\"")
			continue
		}

		break
	}

	fmt.Print("Sending transaction and waiting for confirmation... ")

	transaction, _, err := wlt.SendWaitTransaction(ctx, tx)
	if err != nil {
		fmt.Print("Warning\n")
		return fmt.Errorf("sending transaction error: %w", err)
	}

	block, err := client.CurrentMasterchainInfo(ctx)
	if err != nil {
		fmt.Print("Warning\n")
		return fmt.Errorf("get masterchain info err: %w", err)
	}

	balance, err := wlt.GetBalance(ctx, block)
	if err != nil {
		fmt.Print("Warning\n")
		return fmt.Errorf("GetBalance err: %w", err)
	}

	fmt.Printf("Ok\nTransaction confirmed at block %d, hash: %s, balance left: %s\n", block.SeqNo,
		base64.StdEncoding.EncodeToString(transaction.Hash), balance.String())

	return nil
}
