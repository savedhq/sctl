package auth

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/savedhq/sctl/internal"
	"github.com/savedhq/sctl/internal/render"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "auth",
		Short:       "Authenticate with Saved API",
		Annotations: map[string]string{"auth": "none"},
	}

	cmd.AddCommand(newAuthLoginCmd())
	cmd.AddCommand(newAuthLogoutCmd())
	cmd.AddCommand(newAuthStatusCmd())

	return cmd
}

func newAuthLoginCmd() *cobra.Command {
	var issuer, clientID, audience, scope string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Saved",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")

			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}
			viper.AddConfigPath(homeDir + "/.sctl")
			viper.ReadInConfig()

			// 0. Resolve Server URL (default to prod if not set)
			serverURL := viper.GetString("server_url")
			if serverURL == "" {
				serverURL = "https://api.saved.sh"
			}

			// 1. Resolve Issuer
			if issuer == "" {
				issuer = os.Getenv("AUTH_ISSUER")
				if issuer == "" {
					issuer = viper.GetString("auth_issuer")
				}
			}

			// 2. Resolve Client ID
			if clientID == "" {
				clientID = os.Getenv("AUTH_CLIENT_ID")
				if clientID == "" {
					clientID = viper.GetString("auth_client_id")
				}
			}

			// 3. Resolve Audience
			if audience == "" {
				audience = os.Getenv("AUTH_AUDIENCE")
				if audience == "" {
					audience = viper.GetString("auth_audience")
				}
			}

			if scope == "" {
				scope = "openid profile email offline_access"
			}

			// Dynamic Config Fetch
			if issuer == "" || clientID == "" || audience == "" {
				render.Message(color.CyanString("Fetching authentication configuration from %s...", serverURL))
				authConfig, err := internal.FetchAuthConfig(serverURL)
				if err != nil {
					return fmt.Errorf("failed to fetch auth config (and no flags provided): %w", err)
				}
				issuer = authConfig.Issuer
				clientID = authConfig.ClientID
				audience = authConfig.Audience
				render.Message(color.GreenString("✓ Configuration fetched"))
			}

			if issuer == "" || clientID == "" || audience == "" {
				return fmt.Errorf("authentication configuration requirement invalid after fetch")
			}

			token, err := internal.LoginWithDeviceFlow(issuer, clientID, audience, scope)
			if err != nil {
				return err
			}

			if err := internal.InitConfig(); err != nil {
				return fmt.Errorf("failed to initialize config: %w", err)
			}

			// Save Token and Auth Config
			viper.Set("api_key", token.AccessToken)
			viper.Set("refresh_token", token.RefreshToken)
			viper.Set("auth_issuer", issuer)
			viper.Set("auth_client_id", clientID)
			viper.Set("auth_audience", audience)

			home, _ := os.UserHomeDir()
			configPath := home + "/.sctl/config.yaml"

			if err := viper.WriteConfigAs(configPath); err != nil {
				return fmt.Errorf("failed to save token: %w", err)
			}

			render.Message(color.GreenString("✓ Logged in successfully"))
			render.Message(color.CyanString("Token and auth config saved to ~/.sctl/config.yaml"))

			return nil
		},
	}

	cmd.Flags().StringVar(&issuer, "issuer", "", "OIDC Issuer URL")
	cmd.Flags().StringVar(&clientID, "client-id", "", "OIDC Client ID")
	cmd.Flags().StringVar(&audience, "audience", "", "OIDC Audience")
	cmd.Flags().StringVar(&scope, "scope", "", "OIDC Scopes")

	return cmd
}

func newAuthLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout and remove stored credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
			viper.AddConfigPath("$HOME/.sctl")
			viper.ReadInConfig()

			viper.Set("api_key", "")
			viper.Set("refresh_token", "")

			home, _ := os.UserHomeDir()
			configPath := home + "/.sctl/config.yaml"

			if err := viper.WriteConfigAs(configPath); err != nil {
				return fmt.Errorf("failed to clear credentials: %w", err)
			}

			render.Message(color.GreenString("✓ Logged out successfully"))
			return nil
		},
	}
}

// AuthStatus represents the authentication status.
type AuthStatus struct {
	LoggedIn bool   `json:"logged_in"`
	Message  string `json:"message"`
	Token    string `json:"token,omitempty"`
}

// String implements the Stringer interface for AuthStatus.
func (s AuthStatus) String() string {
	if !s.LoggedIn {
		return fmt.Sprintf("%s %s\n%s",
			color.YellowString("⚠"),
			s.Message,
			color.CyanString("Run 'sctl auth login' to authenticate"),
		)
	}
	tokenStr := ""
	if s.Token != "" {
		tokenLen := len(s.Token)
		if tokenLen > 20 {
			tokenStr = fmt.Sprintf("\nToken: %s...%s", s.Token[:8], s.Token[tokenLen-8:])
		}
	}
	return fmt.Sprintf("%s %s%s", color.GreenString("✓"), s.Message, tokenStr)
}

func newAuthStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
			viper.AddConfigPath("$HOME/.sctl")
			viper.ReadInConfig()

			apiKey := viper.GetString("api_key")

			if apiKey == "" {
				render.Object(AuthStatus{
					LoggedIn: false,
					Message:  "Not logged in",
				})
				return nil
			}

			render.Object(AuthStatus{
				LoggedIn: true,
				Message:  "Logged in",
				Token:    apiKey,
			})

			return nil
		},
	}
}
