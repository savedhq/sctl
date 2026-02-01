package billing

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/savedhq/sctl/internal/render"
	saved "github.com/savedhq/sdk-go"
	"github.com/spf13/cobra"
)

// BillingInfo is a wrapper for the saved.GetBillingInfo200Response type for rendering.
type BillingInfo saved.GetBillingInfo200Response

// String implements the Stringer interface for BillingInfo.
func (b BillingInfo) String() string {
	data, _ := json.MarshalIndent(b, "", "  ")
	return string(data)
}

// UsageHistory is a wrapper for the saved.GetUsageHistory200Response type for rendering.
type UsageHistory saved.GetUsageHistory200Response

// String implements the Stringer interface for UsageHistory.
func (u UsageHistory) String() string {
	var b strings.Builder
	for _, metric := range u.Metrics {
		b.WriteString(fmt.Sprintf("%s %s\n  Type: %s\n  Value: %d %s\n  Cost: $%.2f\n  Created: %s\n\n",
			color.CyanString("Activity:"),
			metric.ActivityName,
			metric.MetricType,
			metric.Value,
			metric.Unit,
			metric.Cost,
			metric.CreatedAt.Format(time.RFC3339),
		))
	}
	return strings.TrimSuffix(b.String(), "\n\n")
}

// Invoices is a wrapper for the saved.ListInvoices200Response type for rendering.
type Invoices saved.ListInvoices200Response

// String implements the Stringer interface for Invoices.
func (i Invoices) String() string {
	var b strings.Builder
	for _, inv := range i.Invoices {
		b.WriteString(fmt.Sprintf("%s %s\n  Number: %s\n  Amount Due: %d %s\n  Status: %s\n  Total: %d\n\n",
			color.CyanString("ID:"),
			inv.Id,
			inv.Number,
			inv.AmountDue,
			inv.Currency,
			inv.Status,
			inv.Total,
		))
	}
	return strings.TrimSuffix(b.String(), "\n\n")
}

// CreditBalance is a wrapper for the saved.ConfirmCreditPurchase200Response type for rendering.
type CreditBalance saved.ConfirmCreditPurchase200Response

// String implements the Stringer interface for CreditBalance.
func (c CreditBalance) String() string {
	return fmt.Sprintf("Balance: %d", c.Balance)
}

// CreditTransactions is a wrapper for the saved.ListCreditTransactions200Response type for rendering.
type CreditTransactions saved.ListCreditTransactions200Response

// String implements the Stringer interface for CreditTransactions.
func (c CreditTransactions) String() string {
	var b strings.Builder
	for _, tx := range c.Transactions {
		b.WriteString(fmt.Sprintf("%s %s\n  Amount: %d\n  Type: %s\n  Created: %s\n\n",
			color.CyanString("ID:"),
			tx.Id,
			tx.Amount,
			tx.Type,
			tx.CreatedAt.Format(time.RFC3339),
		))
	}
	return strings.TrimSuffix(b.String(), "\n\n")
}

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
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Get billing information",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			if err != nil {
				return err
			}

			resp, r, err := cliCtx.Client.BillingAPI.GetBillingInfo(cliCtx.APICtx, workspaceID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			render.Object(BillingInfo(*resp))
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func newBillingUsageCmd() *cobra.Command {
	var workspaceID string
	cmd := &cobra.Command{
		Use:   "usage",
		Short: "Get usage history",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			if err != nil {
				return err
			}

			resp, r, err := cliCtx.Client.BillingAPI.GetUsageHistory(cliCtx.APICtx, workspaceID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			metrics := resp.GetMetrics()
			if len(metrics) == 0 {
				render.Message(color.YellowString("⚠ No usage data found"))
				return nil
			}

			render.Object(UsageHistory(*resp))
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func newBillingInvoicesCmd() *cobra.Command {
	var workspaceID string
	cmd := &cobra.Command{
		Use:   "invoices",
		Short: "List invoices",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			if err != nil {
				return err
			}

			resp, r, err := cliCtx.Client.BillingAPI.ListInvoices(cliCtx.APICtx, workspaceID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			invoices := resp.GetInvoices()
			if len(invoices) == 0 {
				render.Message(color.YellowString("⚠ No invoices found"))
				return nil
			}

			render.Object(Invoices(*resp))
			return nil
		},
	}
	cmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")
	return cmd
}

func newBillingCreditsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credits",
		Short: "Manage credits",
	}

	var workspaceID string

	balanceCmd := &cobra.Command{
		Use:   "balance",
		Short: "Get credit balance",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			if err != nil {
				return err
			}

			resp, r, err := cliCtx.Client.BillingAPI.GetCreditBalance(cliCtx.APICtx, workspaceID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			render.Object(CreditBalance(*resp))
			return nil
		},
	}
	balanceCmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")

	transactionsCmd := &cobra.Command{
		Use:   "transactions",
		Short: "List credit transactions",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := internal.GetCLIContext(cmd.Context())
			var err error
			workspaceID, err = cliCtx.ResolveWorkspaceID(workspaceID)
			if err != nil {
				return err
			}

			resp, r, err := cliCtx.Client.BillingAPI.ListCreditTransactions(cliCtx.APICtx, workspaceID).Execute()
			if err != nil {
				return internal.PrintAPIError(err)
			}
			defer r.Body.Close()

			transactions := resp.GetTransactions()
			if len(transactions) == 0 {
				render.Message(color.YellowString("⚠ No transactions found"))
				return nil
			}

			render.Object(CreditTransactions(*resp))
			return nil
		},
	}
	transactionsCmd.Flags().StringVarP(&workspaceID, "workspace", "w", "", "Workspace ID")

	cmd.AddCommand(balanceCmd)
	cmd.AddCommand(transactionsCmd)

	return cmd
}
