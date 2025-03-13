package cli

import (
	"go-clean-api/pkg/infrastructure/chi_router"
	"go-clean-api/pkg/infrastructure/logger"
	"log"
	"runtime"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "run",
	Short: "Start server",
	Long:  `Start server`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func startServer() {
	config, err := initConfig()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := initDatabase(config)
	if err != nil {
		log.Fatalln(err)
	}

	l, err := logger.NewZapLogger(*config)
	if err != nil {
		log.Fatalln(err)
	}

	// Set the max number of CPUs
	runtime.GOMAXPROCS(config.Server.MaxCPU)

	server := chi_router.NewChiServer(*config, db, l)
	if err = server.Start(); err != nil {
		log.Fatalln(err)
	}
}
