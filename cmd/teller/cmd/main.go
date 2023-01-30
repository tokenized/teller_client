package cmd

import (
	"github.com/spf13/cobra"
)

var scCmd = &cobra.Command{
	Use:   "teller",
	Short: "Teller CLI",
}

func Execute() {
	scCmd.AddCommand(cmdCreateInstrument)
	scCmd.AddCommand(cmdSendTokens)
	scCmd.AddCommand(cmdReceiveTokens)
	scCmd.AddCommand(cmdReclaimBitcoin)
	scCmd.AddCommand(cmdListen)
	scCmd.Execute()
}
