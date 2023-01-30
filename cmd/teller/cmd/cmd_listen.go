package cmd

import (
	"context"
	"fmt"

	teller "github.com/tokenized/teller_client"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var cmdListen = &cobra.Command{
	Use:   "listen request_uuid",
	Short: "Listen for responses. Most likely after a timeout",
	Args:  cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		ctx := context.Background()

		requestID, err := uuid.Parse(args[0])
		if err != nil {
			return errors.Wrap(err, "request_uuid")
		}

		if err := teller.Listen(ctx, requestID); err != nil {
			return fmt.Errorf("Failed to listen : %s\n", err)
		}

		return nil
	},
}
