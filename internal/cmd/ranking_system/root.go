/*
Copyright © 2025 Vladimir Selifanov vladimir.v.selifanov@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package ranking_system

import (
	"fmt"
	"log"
	"os"

	"github.com/SelferRy/ranking_system/internal/config"
	intLogger "github.com/SelferRy/ranking_system/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg     *config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ranking_system",
	Short: "The ranking system chooses which banner to display.",
	Long: `The ranking system implements banner rotation microservice.
This service is designed to select the most effective (clickable) banners in 
conditions of changing user preferences and a set of banners.
It consists of an API and a database that stores information about banners. 
The service provides a gRPC API.

ranking_system has several slots and banners.
A slot is a specific API that a user can interact with.
Each slot can have any number of banners.
Each banner can be in different slots.
Customers are divided into socio-demographic groups. Banners are displayed according to their preferences.

The microservice sends click and impression events to a queue (Kafka) for further processing in analytics systems.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		//_ = cfg  // TODO: make body that will use cfg
		fmt.Println(cfg)
		fmt.Println(cfg.Logger.Level, cfg.Logger.OutputPaths, cfg.Logger.ErrorOutputPaths)
		logger, err := intLogger.New(&cfg.Logger)
		if err != nil {
			log.Fatal("logger didn't initialize.\n", err)
		}
		logger.Info("some")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"./configs/config.yaml",
		"A path to the service configuration file",
	)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().StringP(
	//	"config",
	//	"c",
	//	"./config.yaml",
	//	"A path to the service configuration file.",
	//)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("yaml")
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home + "/configs/")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config.yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Using config file:", viper.ConfigFileUsed())

	cfg = &config.Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
	fmt.Printf("%T\n", cfg)
	fmt.Println(viper.Get("test"))
}
