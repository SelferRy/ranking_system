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
package cli

import (
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ranking_system",
	Short: "The ranking system (banner rotation) microservices.",
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
	rootCmd.AddCommand(NewServeCommand())
}
