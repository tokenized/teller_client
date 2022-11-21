package cmd

import (
	"context"
	"fmt"

	"github.com/tokenized/pkg/bitcoin"
	teller "github.com/tokenized/teller_client"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var cmdReclaimBitcoin = &cobra.Command{
	Use:   "reclaim <address>",
	Short: "Send bitcoin to an address.",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Incorrect parameter count \"reclaim <address>\"")
		}

		address, err := bitcoin.DecodeAddress(args[0])
		if err != nil {
			fmt.Printf("Invalid receive address : %s\n", err)
			return nil
		}

		request := &teller.ReclaimBitcoin{
			ID:      uuid.New(),
			Address: address,
		}

		ctx := context.Background()
		if err := teller.ProcessRequest(ctx, request.ID, request); err != nil {
			fmt.Printf("Failed to reclaim bitcoin : %s\n", err)
			return nil
		}

		return nil
	},
}
