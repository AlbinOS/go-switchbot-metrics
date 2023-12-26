/*
Copyright © 2023 Albin Gilles <gilles.albin@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"github.com/AlbinOS/go-switchbot-metrics/serve"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		serve.Init()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	serveCmd.PersistentFlags().String("bind_ip", "127.0.0.1", "App ip to listen to")
	viper.BindPFlag("bind_ip", serveCmd.PersistentFlags().Lookup("bind_ip"))

	serveCmd.PersistentFlags().String("bind_port", "3000", "App port to listen to")
	viper.BindPFlag("bind_port", serveCmd.PersistentFlags().Lookup("bind_port"))

	serveCmd.PersistentFlags().String("switchbot_openapi_token", "", "SwitchBot Open API token")
	viper.BindPFlag("switchbot_openapi_token", serveCmd.PersistentFlags().Lookup("switchbot_openapi_token"))

	serveCmd.PersistentFlags().String("switchbot_secret_key", "", "SwitchBot API secery key")
	viper.BindPFlag("switchbot_secret_key", serveCmd.PersistentFlags().Lookup("switchbot_secret_key"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
