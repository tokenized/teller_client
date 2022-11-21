package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/tokenized/specification/dist/golang/protocol"
	teller "github.com/tokenized/teller_client"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var cmdReceiveTokens = &cobra.Command{
	Use:   "receive <instrument id> <amount>",
	Short: "Tell teller to receive tokens.",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Incorrect parameter count \"receive <instrument id> <amount>>\"")
		}

		instrumentID := args[0]
		if _, _, err := protocol.DecodeInstrumentID(instrumentID); err != nil {
			fmt.Printf("Invalid instrument ID : %s : %s\n", instrumentID, err)
			return nil
		}

		amount, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Invalid amount : %s\n", err)
			return nil
		}

		request := &teller.ReceiveTokens{
			ID:           uuid.New(),
			InstrumentID: instrumentID,
			Quantity:     uint64(amount),
		}

		ctx := context.Background()
		if err := teller.ProcessRequest(ctx, request.ID, request); err != nil {
			fmt.Printf("Failed to receive tokens : %s\n", err)
			return nil
		}

		return nil
	},
}
