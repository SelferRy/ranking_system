/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"context"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/app"
	"github.com/SelferRy/ranking_system/internal/config"
	"github.com/SelferRy/ranking_system/internal/infra/logger"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

type ServeOptions struct {
	ConfigFile string
}

// NewServeCommand constructor for runServe cmd
func NewServeCommand() *cobra.Command {
	var opts ServeOptions

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run gRPC server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServe(opts)
		},
	}

	cmd.Flags().StringVar(&opts.ConfigFile, "config", "./configs/config.yaml", "path to config")
	return cmd
}

func runServe(opts ServeOptions) error {
	// Load configuration
	cfg, err := loadConfig(opts.ConfigFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Initialize logger
	logg, err := logger.New(cfg.Logger)
	if err != nil {
		return fmt.Errorf("create logger: %w", err)
	}
	defer func() {
		_ = logg.Sync()
	}()

	logg.Info("Starting ranking system",
		logger.String("version", "1.0.0"),
		logger.Int("pid", os.Getpid()),
	)

	// Create application context with graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// Initialize and run application
	application, err := app.New(ctx, cfg, logg)
	if err != nil {
		return fmt.Errorf("create application: %w", err)
	}

	return application.Run(ctx)
}

func loadConfig(configFile string) (config.Config, error) {
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("logger.level", "DEBUG")
	viper.SetDefault("database.max_connections", 10)

	// Read config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config.Config{}, fmt.Errorf("read config =: %w", err)
		}
		// Config file is optional if using env vars/flags
	}

	viper.AutomaticEnv()

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return config.Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}
