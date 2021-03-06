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
	"github.com/sebnyberg/gchat/pkg/client"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "client commands",
	Long: `gchat client
	
Connect to the server with: 
	> gchat client connect`,
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connects to the chat server",
	Long:  `Connects to the chat server`,
	Run: func(cmd *cobra.Command, args []string) {
		client.StartClient(username)
	},
}

var username string

func init() {
	rootCmd.AddCommand(clientCmd)

	connectCmd.Flags().StringVarP(&username, "username", "u", "", "Username (required)")
	connectCmd.MarkFlagRequired("username")

	clientCmd.AddCommand(connectCmd)
}
