package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	teller "github.com/tokenized/teller_client"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var cmdSendTokens = &cobra.Command{
	Use:   "send <currency code> <paymail> <amount>",
	Short: "Send tokens to an address.",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) != 3 {
			return errors.New("Incorrect parameter count \"send <instrument id> <receiver address> <amount>>\"")
		}

		code := args[0]
		paymail := args[1]
		amountText := args[2]

		request := &teller.SendTokens{
			ID: uuid.New(),
		}

		if len(code) != 3 || strings.ToUpper(code) != code {
			fmt.Printf("Invalid currency code : %s\n", code)
			return nil
		}
		request.CurrencyCode = code

		if strings.Contains(paymail, "@") {
			request.Paymail = paymail
		} else {
			fmt.Printf("Invalid paymail : %s\n", paymail)
			return nil
		}

		amount, err := strconv.Atoi(amountText)
		if err != nil {
			fmt.Printf("Invalid amount : %s\n", err)
			return nil
		}
		request.Quantity = uint64(amount)

		ctx := context.Background()
		if err := teller.ProcessRequest(ctx, request.ID, request); err != nil {
			fmt.Printf("Failed to send tokens : %s\n", err)
			return nil
		}

		return nil
	},
}
