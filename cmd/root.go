/*
 * Copyright (c) 2025 Vladislav Trofimenko <com@slashfast.dev>
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of this software and associated documentation files (the
 * "Software"), to deal in the Software without restriction, including
 * without limitation the rights to use, copy, modify, merge, publish,
 * distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to
 * the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
 * LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
 * OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
 * WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package cmd

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strings"
	"tondns/internal/config"
	"tondns/pkg/ton"

	"github.com/mitchellh/mapstructure"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg config.Config

var rootCmd = &cobra.Command{
	Use:   "tondns",
	Short: "TON DNS domain management tool",
	Long: `A command-line tool for managing TON DNS domains.
	
Allows checking domain expiration dates, managing domain renewals,
and monitoring domain status. Supports both interactive and automated usage.

For detailed documentation and examples, visit:
https://github.com/slashfast/tonft`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(unmarshalArgs)

	cfg.WalletVersion = ton.V5R1

	rootCmd.PersistentFlags().StringArrayP("seed", "s", []string{}, "wallet seed phrase for authentication")
	rootCmd.PersistentFlags().Var(&cfg.WalletVersion, "version", "wallet version (e.g., v3r2, v4r2)")
	rootCmd.PersistentFlags().Var(&cfg.Proxy, "proxy", "HTTP/SOCKS5 proxy URL (e.g., socks5://127.0.0.1:9050)")
	rootCmd.PersistentFlags().BoolP("lite", "l", false, "enable lite mode for read-only operations (no wallet required)")
	rootCmd.PersistentFlags().Bool("wait-tx", false, "wait for transaction confirmations (slower but more reliable)")

	viper.BindPFlags(rootCmd.PersistentFlags())
}

func unmarshalArgs() {
	viper.AutomaticEnv()
	viper.ReadInConfig()
	viper.SetTypeByDefaultValue(true)

	if err := viper.Unmarshal(
		&cfg,
		viper.DecodeHook(
			mapstructure.ComposeDecodeHookFunc(
				func(f, t reflect.Type, data any) (any, error) {
					if f.Kind() != reflect.String || t != reflect.TypeOf(config.ProxyURL{}) {
						return data, nil
					}

					return url.Parse(data.(string))
				},
				mapstructure.StringToTimeDurationHookFunc(),
			),
		),
	); err != nil {
		cobra.CheckErr(err)
	}

	if !cfg.Lite {
		if len(cfg.Seed) == 0 {
			cobra.CheckErr(fmt.Errorf("seed is empty"))
		} else {
			cfg.Seed = strings.Fields(cfg.Seed[0])
		}
	}
}
