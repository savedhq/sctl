package billing

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/spf13/cobra"
)

func NewBillingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "billing",
		Aliases: []string{"bill"},
		Short:   "Manage billing and credits",
	}

	cmd.AddCommand(newBillingInfoCmd())
	cmd.AddCommand(newBillingCreditsCmd())
	cmd.AddCommand(newBillingUsageCmd())
	cmd.AddCommand(newBillingInvoicesCmd())

	return cmd
}

func newBillingInfoCmd() *cobra.Command {
	var workspaceID string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Get billing information",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			if workspaceID == "" {
				workspaceID = cliCtx.GetWorkspaceID()
			}
			if workspaceID == "" {
				return fmt.Errorf("workspace_id required (use --workspace or set in config)")
			}

			resp, r, err := cliCtx.Client.BillingAPI.GetBillingInfo(cliCtx.APICtx, workspaceID).Execute()
			if err != nil {
				return fmt.Errorf("API error: %v", err)
			}
			defer r.Body.Close()

			if jsonOutput {
				data, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Println(string(data))
				return nil
			}

			// Human-readable output
			data, err := json.Marshal(resp)
			if err != nil {
				return fmt.Errorf("failed to marshal json: %w", err)
			}

			var infoMap map[string]interface{}
			if err := json.Unmarshal(data, &infoMap); err != nil {
				return fmt.Errorf("failed to unmarshal json into map: %w", err)
			}

			color.Cyan("Billing Information:")
			for key, value := range infoMap {
				// Simple pretty printing for nested objects
				if v, ok := value.(map[string]interface{}); ok {
					fmt.Printf("  %s:\n", key)
					for k, val := range v {
						fmt.Printf("    %s: %v\n", k, val)
					}
				} else {
					fmt.Printf("  %s: %v\n", key, value)
				}
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	return cmd
}

func newBillingUsageCmd() *cobra.Command {
	var workspaceID string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "usage",
		Short: "Get usage history",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			if workspaceID == "" {
				workspaceID = cliCtx.GetWorkspaceID()
			}
			if workspaceID == "" {
				return fmt.Errorf("workspace_id required (use --workspace or set in config)")
			}

			resp, r, err := cliCtx.Client.BillingAPI.GetUsageHistory(cliCtx.APICtx, workspaceID).Execute()
			if err != nil {
				return fmt.Errorf("API error: %v", err)
			}
			defer r.Body.Close()

			if jsonOutput {
				data, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Println(string(data))
				return nil
			}

			metrics := resp.GetMetrics()
			if len(metrics) == 0 {
				color.Yellow("⚠ No usage data found")
				return nil
			}

			for _, metric := range metrics {
				color.Cyan("Activity: %s", metric.GetActivityName())
				fmt.Printf("  Type: %s\n", metric.GetMetricType())
				fmt.Printf("  Value: %d %s\n", metric.GetValue(), metric.GetUnit())
				fmt.Printf("  Cost: $%.2f\n", metric.GetCost())
				fmt.Printf("  Created: %s\n\n", metric.GetCreatedAt().Format("2006-01-02 15:04:05"))
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	return cmd
}

func newBillingInvoicesCmd() *cobra.Command {
	var workspaceID string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "invoices",
		Short: "List invoices",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			if workspaceID == "" {
				workspaceID = cliCtx.GetWorkspaceID()
			}
			if workspaceID == "" {
				return fmt.Errorf("workspace_id required (use --workspace or set in config)")
			}

			resp, r, err := cliCtx.Client.BillingAPI.ListInvoices(cliCtx.APICtx, workspaceID).Execute()
			if err != nil {
				return fmt.Errorf("API error: %v", err)
			}
			defer r.Body.Close()

			if jsonOutput {
				data, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Println(string(data))
				return nil
			}

			invoices := resp.GetInvoices()
			if len(invoices) == 0 {
				color.Yellow("⚠ No invoices found")
				return nil
			}

			for _, inv := range invoices {
				color.Cyan("ID: %s", inv.GetId())
				fmt.Printf("  Number: %s\n", inv.GetNumber())
				fmt.Printf("  Amount Due: %d %s\n", inv.GetAmountDue(), inv.GetCurrency())
				fmt.Printf("  Status: %s\n", inv.GetStatus())
				fmt.Printf("  Total: %d\n\n", inv.GetTotal())
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	return cmd
}

func newBillingCreditsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credits",
		Short: "Manage credits",
	}

	var workspaceID string

	var jsonOutput bool
	balanceCmd := &cobra.Command{
		Use:   "balance",
		Short: "Get credit balance",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			if workspaceID == "" {
				workspaceID = cliCtx.GetWorkspaceID()
			}
			if workspaceID == "" {
				return fmt.Errorf("workspace_id required (use --workspace or set in config)")
			}

			resp, r, err := cliCtx.Client.BillingAPI.GetCreditBalance(cliCtx.APICtx, workspaceID).Execute()
			if err != nil {
				return fmt.Errorf("API error: %v", err)
			}
			defer r.Body.Close()

			if jsonOutput {
				data, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Println(string(data))
				return nil
			}

			// Human-readable output
			color.Cyan("Credit Balance:")
			fmt.Printf("  Available: %d %s\n", resp.GetBalance(), resp.GetCurrency())

			return nil
		},
	}
	balanceCmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	balanceCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	transactionsCmd := &cobra.Command{
		Use:   "transactions",
		Short: "List credit transactions",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			if cliCtx == nil {
				return fmt.Errorf("CLI context not initialized")
			}

			if workspaceID == "" {
				workspaceID = cliCtx.GetWorkspaceID()
			}
			if workspaceID == "" {
				return fmt.Errorf("workspace_id required (use --workspace or set in config)")
			}

			resp, r, err := cliCtx.Client.BillingAPI.ListCreditTransactions(cliCtx.APICtx, workspaceID).Execute()
			if err != nil {
				return fmt.Errorf("API error: %v", err)
			}
			defer r.Body.Close()

			if jsonOutput {
				data, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal json: %w", err)
				}
				fmt.Println(string(data))
				return nil
			}

			transactions := resp.GetTransactions()
			if len(transactions) == 0 {
				color.Yellow("⚠ No transactions found")
				return nil
			}

			for _, tx := range transactions {
				color.Cyan("ID: %s", tx.GetId())
				fmt.Printf("  Amount: %d\n", tx.GetAmount())
				fmt.Printf("  Type: %s\n", tx.GetType())
				fmt.Printf("  Created: %s\n\n", tx.GetCreatedAt().Format("2006-01-02 15:04:05"))
			}
			return nil
		},
	}
	transactionsCmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	transactionsCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	cmd.AddCommand(balanceCmd)
	cmd.AddCommand(transactionsCmd)

	return cmd
}
