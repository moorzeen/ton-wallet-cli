package account

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/xssnick/tonutils-go/ton/wallet"
)

func NewVanity(suffix string, withSeed bool) (*Account, error) {
	ctx, cancel := context.WithCancel(context.Background())

	ws := make(chan Account, 1)
	threads := runtime.NumCPU()
	var counter uint64

	if !withSeed {
		for x := 0; x < threads; x++ {
			go func(ctx context.Context, cancel context.CancelFunc, ws chan<- Account) {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						atomic.AddUint64(&counter, 1)

						_, pk, _ := ed25519.GenerateKey(nil)
						w, err := wallet.FromPrivateKey(nil, pk, wallet.V4R2)
						if err != nil {
							continue
						}

						if strings.HasSuffix(w.WalletAddress().String(), suffix) {
							ws <- Account{
								Wallet: w,
								Seed:   "was not used",
								Key:    hex.EncodeToString(w.PrivateKey().Seed()),
								State:  nil,
							}

							cancel()
							break
						}
					}
				}
			}(ctx, cancel, ws)
		}
	} else {
		for x := 0; x < threads; x++ {
			go func(ctx context.Context, cancel context.CancelFunc, ws chan<- Account) {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						atomic.AddUint64(&counter, 1)

						seed := wallet.NewSeed()
						w, _ := wallet.FromSeed(nil, seed, wallet.V4R2)

						if strings.HasSuffix(w.WalletAddress().String(), suffix) {
							ws <- Account{
								Wallet: w,
								Seed:   strings.Join(seed, " "),
								Key:    hex.EncodeToString(w.PrivateKey().Seed()),
								State:  nil,
							}

							cancel()
							return
						}
					}
				}
			}(ctx, cancel, ws)
		}
	}

	fmt.Printf("searching vanity address on %d CPUs...\n", threads)
	for {
		time.Sleep(1 * time.Second)
		fmt.Printf("checked %d per second\n", atomic.LoadUint64(&counter))
		atomic.StoreUint64(&counter, 0)
		select {
		case <-ctx.Done():
			a := <-ws
			return &a, nil
		default:

		}
	}
}
