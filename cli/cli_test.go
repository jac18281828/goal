package cli

import (
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		envs      map[string]string
		want      *Options
		expectErr bool
		errorMsg  string
	}{
		{
			name: "All required flags provided via command line",
			args: []string{
				"run",
				"--chain", "mainnet",
				"--ethereum.rpc-base-url", "http://localhost:8545",
				"--ethereum.ws-url", "ws://localhost:8546",
				"--etherscan-api-keys", "your-api-key",
			},
			envs: map[string]string{},
			want: &Options{
				Chain:              "mainnet",
				EthereumRPCBaseURL: "http://localhost:8545",
				EthereumWSURL:      "ws://localhost:8546",
				SqliteInMemory:     false,
				SqliteDBFilePath:   "./sqlite.db",
				RPCGRPCPort:        7100,
				RPCHTTPPort:        7101,
				EtherscanAPIKeys:   "your-api-key",
			},
			expectErr: false,
		},
		{
			name: "Missing required flag --chain",
			args: []string{
				"run",
				"--ethereum.rpc-base-url", "http://localhost:8545",
				"--ethereum.ws-url", "ws://localhost:8546",
				"--etherscan-api-keys", "your-api-key",
			},
			envs:      map[string]string{},
			expectErr: true,
			errorMsg:  "--chain is required",
		},
		{
			name: "Default values are applied",
			args: []string{
				"run",
				"--chain", "preprod",
				"--ethereum.rpc-base-url", "http://localhost:8545",
				"--ethereum.ws-url", "ws://localhost:8546",
				"--etherscan-api-keys", "your-api-key",
			},
			envs: map[string]string{},
			want: &Options{
				Chain:              "preprod",
				EthereumRPCBaseURL: "http://localhost:8545",
				EthereumWSURL:      "ws://localhost:8546",
				SqliteInMemory:     false,
				SqliteDBFilePath:   "./sqlite.db",
				RPCGRPCPort:        7100,
				RPCHTTPPort:        7101,
				EtherscanAPIKeys:   "your-api-key",
			},
			expectErr: false,
		},
		{
			name: "Flags overridden by environment variables",
			args: []string{
				"run",
				"--chain", "holesky",
				"--etherscan-api-keys", "your-api-key",
			},
			envs: map[string]string{
				"ETHEREUM_RPC_BASE_URL": "http://env-rpc:8545",
				"ETHEREUM_WS_URL":       "ws://env-ws:8546",
				"SQLITE_IN_MEMORY":      "true",
				"RPC_GRPC_PORT":         "8100",
				"RPC_HTTP_PORT":         "8101",
			},
			want: &Options{
				Chain:              "holesky",
				EthereumRPCBaseURL: "http://env-rpc:8545",
				EthereumWSURL:      "ws://env-ws:8546",
				SqliteInMemory:     true,
				SqliteDBFilePath:   "./sqlite.db",
				RPCGRPCPort:        8100,
				RPCHTTPPort:        8101,
				EtherscanAPIKeys:   "your-api-key",
			},
			expectErr: false,
		},
		{
			name: "Invalid integer in environment variable",
			args: []string{
				"run",
				"--chain", "mainnet",
				"--ethereum.rpc-base-url", "http://localhost:8545",
				"--ethereum.ws-url", "ws://localhost:8546",
				"--etherscan-api-keys", "your-api-key",
			},
			envs: map[string]string{
				"RPC_GRPC_PORT": "not-an-int",
			},
			want: &Options{
				Chain:              "mainnet",
				EthereumRPCBaseURL: "http://localhost:8545",
				EthereumWSURL:      "ws://localhost:8546",
				SqliteInMemory:     false,
				SqliteDBFilePath:   "./sqlite.db",
				RPCGRPCPort:        7100, // Default value because of parsing error
				RPCHTTPPort:        7101,
				EtherscanAPIKeys:   "your-api-key",
			},
			expectErr: false,
		},
		{
			name: "Invalid boolean in environment variable",
			args: []string{
				"run",
				"--chain", "mainnet",
				"--ethereum.rpc-base-url", "http://localhost:8545",
				"--ethereum.ws-url", "ws://localhost:8546",
				"--etherscan-api-keys", "your-api-key",
			},
			envs: map[string]string{
				"SQLITE_IN_MEMORY": "not-a-bool",
			},
			want: &Options{
				Chain:              "mainnet",
				EthereumRPCBaseURL: "http://localhost:8545",
				EthereumWSURL:      "ws://localhost:8546",
				SqliteInMemory:     false, // Default value because of parsing error
				SqliteDBFilePath:   "./sqlite.db",
				RPCGRPCPort:        7100,
				RPCHTTPPort:        7101,
				EtherscanAPIKeys:   "your-api-key",
			},
			expectErr: false,
		},
		{
			name: "No arguments or environment variables",
			args: []string{"run"},
			envs: map[string]string{},
			want: &Options{
				SqliteInMemory:   false,
				SqliteDBFilePath: "./sqlite.db",
				RPCGRPCPort:      7100,
				RPCHTTPPort:      7101,
			},
			expectErr: true,
			errorMsg:  "--chain is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := ParseArgs(tt.args, tt.envs)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Fatalf("Expected error message '%s', but got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if !reflect.DeepEqual(opts, tt.want) {
					t.Errorf("Options mismatch.\nGot:  %+v\nWant: %+v", opts, tt.want)
				}
			}
		})
	}
}
