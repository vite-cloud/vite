package proxy

import (
	"github.com/spf13/cobra"
	"github.com/vite-cloud/vite/core/handler/cli/cli"
	"net/http"
)

func runRunCommand(cli *cli.CLI) error {
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello world"))
		}),
	}

	return server.ListenAndServe()
}

func NewRunCommand(cli *cli.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "run",
		Short:  "run the proxy",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRunCommand(cli)
		},
	}

	return cmd
}
