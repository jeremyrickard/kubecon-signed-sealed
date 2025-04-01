package cmd

import (
	"github.com/jeremyrickard/kubecon-2025-ssd/cmd/retag"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	rootName             = "azcu"
	rootShortDescription = "A collection of CLI tools for the Azure Container Upstream team"
	rootLongDescription  = "A collection of CLI tools for the Azure Container Upstream team"
)

var (
	debug bool
	trace bool
)
var ConstraintsFile string

// NewRootCmd returns the root command for azcu-cli.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   rootName,
		Short: rootShortDescription,
		Long:  rootLongDescription,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				log.SetLevel(log.DebugLevel)
			}
			if trace {
				log.SetLevel(log.TraceLevel)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	p := rootCmd.PersistentFlags()
	p.BoolVar(&debug, "debug", false, "enable debug level logging")
	p.BoolVar(&trace, "trace", false, "enable trace level logging")

	rootCmd.AddCommand(retag.NewRetagCmd())

	return rootCmd
}
