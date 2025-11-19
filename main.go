package main

import (
	"loan_system/cmd/loan"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "Amartha Loan Service",
		Short: "Amartha Loan Service",
		Long:  `Backend Service for Amartha Loan Service Project`,
		Run: func(cmd *cobra.Command, args []string) {
			loan.Execute()
		},
	}

	_ = rootCmd.Execute()
}
