package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/tokenized/specification/dist/golang/instruments"
	teller "github.com/tokenized/teller_client"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var cmdCreateInstrument = &cobra.Command{
	Use:   "create instrument_type currency_code precision use_identity_oracle enforcement_orders_permitted",
	Short: "Send a request to create a new instrument.",
	Args:  cobra.MinimumNArgs(3),
	RunE: func(c *cobra.Command, args []string) error {

		instrumentType := args[0]

		currencyCode := args[1]
		currencyData := instruments.CurrenciesData(currencyCode)
		if currencyData == nil {
			return fmt.Errorf("Unsupported currency code : %s", currencyCode)
		}

		precision, err := strconv.ParseUint(args[2], 10, 64)
		if err != nil {
			return errors.Wrap(err, "precision")
		}
		if precision < 2 || precision > 6 {
			return fmt.Errorf("Precision out of range (2 is cents, 6 is 1/10000 of cents): %d", precision)
		}

		useIdentityOracle := false
		if len(args) > 3 {
			switch strings.ToLower(args[3]) {
			case "true":
				useIdentityOracle = true
			case "false":
			default:
				return fmt.Errorf("use_identity_oracle must be \"true\" or \"false\" : %s\n",
					args[3])
			}
		}

		allowEnforcementOrders := false
		if len(args) > 4 {
			switch strings.ToLower(args[4]) {
			case "true":
				allowEnforcementOrders = true
			case "false":
			default:
				return fmt.Errorf("enforcement_orders_permitted must be \"true\" or \"false\" : %s\n",
					args[4])
			}
		}

		request := &teller.CreateInstrument{
			ID:                         uuid.New(),
			CurrencyCode:               currencyCode,
			UseIdentityOracle:          useIdentityOracle,
			EnforcementOrdersPermitted: allowEnforcementOrders,
			InstrumentType:             instrumentType,
			Precision:                  uint(precision),
			// EntityContract *bitcoin.RawAddress `bsor:"7" json:"entity_contract"`
		}

		ctx := context.Background()
		if err := teller.ProcessRequest(ctx, request.ID, request); err != nil {
			return fmt.Errorf("Failed to create instrument : %s\n", err)
		}

		return nil
	},
}
