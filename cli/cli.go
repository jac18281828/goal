package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// Options holds the configuration values for the sidecar utility.
type Options struct {
	Chain              string
	EthereumRPCBaseURL string
	EthereumWSURL      string
	SqliteInMemory     bool
	SqliteDBFilePath   string
	RPCGRPCPort        int
	RPCHTTPPort        int
	EtherscanAPIKeys   string
}

// ParseArgs parses command-line arguments and environment variables.
func ParseArgs(args []string, envs map[string]string) (*Options, error) {
	opts := &Options{}

	var rootCmd = &cobra.Command{
		Use:   "sidecar",
		Short: "Sidecar utility",
	}

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the sidecar utility",
		RunE: func(cmd *cobra.Command, cmdArgs []string) error {
			// Validate required flags
			if opts.Chain == "" {
				return fmt.Errorf("--chain is required")
			}
			if opts.EthereumRPCBaseURL == "" {
				return fmt.Errorf("--ethereum.rpc-base-url is required")
			}
			if opts.EthereumWSURL == "" {
				return fmt.Errorf("--ethereum.ws-url is required")
			}
			if opts.EtherscanAPIKeys == "" {
				return fmt.Errorf("--etherscan-api-keys is required")
			}
			return nil
		},
	}

	// Set flags with environment variable support
	runCmd.Flags().StringVar(&opts.Chain, "chain", getEnv(envs, "CHAIN", ""), "<mainnet | holesky | preprod> (required)")
	runCmd.Flags().StringVar(&opts.EthereumRPCBaseURL, "ethereum.rpc-base-url", getEnv(envs, "ETHEREUM_RPC_BASE_URL", ""), "Ethereum RPC base URL (required)")
	runCmd.Flags().StringVar(&opts.EthereumWSURL, "ethereum.ws-url", getEnv(envs, "ETHEREUM_WS_URL", ""), "Ethereum WebSocket URL (required)")
	runCmd.Flags().BoolVar(&opts.SqliteInMemory, "sqlite.in-memory", getEnvBool(envs, "SQLITE_IN_MEMORY", false), "Use SQLite in-memory database")
	runCmd.Flags().StringVar(&opts.SqliteDBFilePath, "sqlite.db-file-path", getEnv(envs, "SQLITE_DB_FILE_PATH", "./sqlite.db"), "SQLite database file path")
	runCmd.Flags().IntVar(&opts.RPCGRPCPort, "rpc.grpc-port", getEnvInt(envs, "RPC_GRPC_PORT", 7100), "gRPC port")
	runCmd.Flags().IntVar(&opts.RPCHTTPPort, "rpc.http-port", getEnvInt(envs, "RPC_HTTP_PORT", 7101), "HTTP port")
	runCmd.Flags().StringVar(&opts.EtherscanAPIKeys, "etherscan-api-keys", getEnv(envs, "ETHERSCAN_API_KEYS", ""), "Etherscan API keys (required)")

	rootCmd.AddCommand(runCmd)

	// Set the command-line arguments
	rootCmd.SetArgs(args)

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		return nil, err
	}

	return opts, nil
}

// Helper functions to get environment variables with default values
func getEnv(envs map[string]string, key string, defaultValue string) string {
	if value, exists := envs[key]; exists {
		return value
	}
	return defaultValue
}

func getEnvInt(envs map[string]string, key string, defaultValue int) int {
	if valueStr, exists := envs[key]; exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvBool(envs map[string]string, key string, defaultValue bool) bool {
	if valueStr, exists := envs[key]; exists {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
