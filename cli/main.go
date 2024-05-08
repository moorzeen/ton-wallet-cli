package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/moorzeen/ton-wallet-cli/internal/account"
)

var (
	usage = "Usage: ./twc <command> [arguments]\n\n" +
		"Commands:\n" +
		"\tnew\tcreate regular wallet\n" +
		"\tvanity\tcreate wallet with address suffix\n" +
		"\tbalance\tget wallet balance and more\n" +
		"\tsend\tsend TON to another account\n" +
		"\tversion\tshow build version and exit\n" +
		"\thelp\tshow this text\n\n" +
		"Use \"./twc <command> -h\" for more information about a command.\n\n" +
		"To run in testnet add \"--testnet\" as the last flag in the command line."

	gitCommit string
)

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	isTestnet := fs.Bool("testnet", false, "testnet mode")

	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "new":
		acc, err := account.New()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		acc.Print()

	case "vanity":
		suffix := fs.String("suffix", "", "desired address suffix")
		withSeed := fs.Bool("seed", false, "use seed phrase (slow)")

		err := fs.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if *suffix == "" {
			fmt.Println(errors.New("suffix is not set, usage: ./twc vanity --suffix=[your suffix] --seed(bool)"))
			os.Exit(1)
		}

		acc, err := account.NewVanity(*suffix, *withSeed)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		acc.Print()

	case "balance":
		address := fs.String("address", "", "wallet address to check balance")

		err := fs.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if *address == "" {
			fmt.Println(errors.New("address is not set, usage: ./twc balance --address=[wallet address]"))
			os.Exit(1)
		}

		st, err := account.GetState(*address, *isTestnet)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		st.Print()

	case "send":
		to := fs.String("to", "", "destination address")
		msg := fs.String("msg", "", "message")
		amount := fs.String("amount", "", "amount of TON to send, example: \"2,33\" or \"100\" or \"44.4\" or \"ALL\"")
		key := fs.String("key", "", "seed phrase or private key")

		err := fs.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if *to == "" || *amount == "" || *key == "" {
			fmt.Println("flags --to, --amount, --key are required\nusage: ./twc send --to=[destination address] --msg=\"[comment (optional)]\" --amount=[amount] --key=\"[seed phrase or private key]\"")
			os.Exit(1)
		}

		err = account.Send(*to, *msg, *amount, *key, *isTestnet)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	case "version":
		fmt.Println("Build version:", gitCommit)
		os.Exit(0)

	case "help":
		fmt.Println(usage)
		os.Exit(0)

	default:
		fmt.Println("Unknown command:", os.Args[1], "\n", usage)
		os.Exit(1)
	}

}
