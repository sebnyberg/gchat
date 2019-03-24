// Copyright © 2019 SEBASTIAN NYBERG <sebastian@sebnyberg.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/sebnyberg/gchat/pkg/server"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server commands",
	Long: `gchat server
	
Start the server with: 
	> gchat server start`,
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the server",
	Long: `Starts the gchat server
	
Connect to the chat with:
	> gchat client connect`,
	Run: func(cmd *cobra.Command, args []string) {
		err := server.StartServer()
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(startCmd)
}
