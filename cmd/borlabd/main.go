// Command borlabd is the bor-lab simulation node daemon. It exposes a small
// cobra command tree: `version`, `start`, and `tx send`.
package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"

	"github.com/devcon-dutralab/bor-lab/app"
	"github.com/devcon-dutralab/bor-lab/x/transfer/types"
)

const version = "0.1.0"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "borlabd:", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "borlabd",
		Short:         "bor-lab simulation node",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(newVersionCmd(), newStartCmd(), newTxCmd())
	return root
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the node version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Println(version)
			return nil
		},
	}
}

func newStartCmd() *cobra.Command {
	var blocks int
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Run the node against a fresh in-memory app",
		RunE: func(cmd *cobra.Command, _ []string) error {
			application := app.NewApp()
			application.Transfer.SetBalance("alice", 1000)
			application.Transfer.SetBalance("bob", 1000)

			// The node fans block production and the mempool drain out
			// across two workers.
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				for h := int64(1); h <= int64(blocks); h++ {
					_ = application.BeginBlock(h)
					_ = application.EndBlock(h)
				}
			}()
			go func() {
				defer wg.Done()
				for h := 1; h <= blocks; h++ {
					_ = application.Transfer.Send(types.MsgSend{
						FromAddress: "alice",
						ToAddress:   "bob",
						Amount:      types.Coin{Denom: types.DefaultDenom, Amount: 1},
					})
				}
			}()
			wg.Wait()

			cmd.Printf("produced %d blocks; alice=%d bob=%d\n",
				blocks, application.Transfer.Balance("alice"), application.Transfer.Balance("bob"))
			return nil
		},
	}
	cmd.Flags().IntVar(&blocks, "blocks", 10, "number of blocks to produce")
	return cmd
}

func newTxCmd() *cobra.Command {
	tx := &cobra.Command{Use: "tx", Short: "Transaction subcommands"}

	var from, to, amount string
	send := &cobra.Command{
		Use:   "send",
		Short: "Send tokens between accounts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := types.ParseCoin(amount)
			if err != nil {
				return err
			}
			application := app.NewApp()
			application.Transfer.SetBalance(from, 1000)
			if err := application.Transfer.Send(types.MsgSend{
				FromAddress: from,
				ToAddress:   to,
				Amount:      c,
			}); err != nil {
				return err
			}
			cmd.Printf("sent %s from %s to %s\n", c, from, to)
			return nil
		},
	}
	send.Flags().StringVar(&from, "from", "", "sender address")
	send.Flags().StringVar(&to, "to", "", "recipient address")
	send.Flags().StringVar(&amount, "amount", "", "amount, e.g. 100ubor")
	tx.AddCommand(send)
	return tx
}
