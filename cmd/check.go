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
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"github.com/slashfast/tondns/internal/check"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var checkCmd = &cobra.Command{
	Use:   "check [domains...]",
	Short: "Check TON DNS domains expiration status",
	Long: `The check command allows you to verify the expiration status of one or multiple TON DNS domains.

When run without specifying domains and lite mode disabled - checks all domains in the wallet.


Examples:
tondns check durov.ton
tondns check example.ton durov.ton
tondns check --pretty
tondns check -o domains.json example.ton durov.ton`,
	Run: func(cmd *cobra.Command, domains []string) {
		isCheckAllMode := len(domains) == 0

		if err := validateInput(domains, isCheckAllMode); err != nil {
			cmd.PrintErrln(err)
			return
		}

		checker, err := check.NewChecker(cfg)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		var allResults []check.Result
		if isCheckAllMode {
			allResults = checker.SmartCheckAll()
		} else {
			allResults = checkDomains(checker, domains)
		}

		if err := processResults(cmd, allResults); err != nil {
			cmd.PrintErrln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().String("delay", "100ms", "delay between API requests (e.g., 100ms, 1s)")
	checkCmd.Flags().BoolP("renew", "r", false, "automatically renew domains when they are about to expire")
	checkCmd.Flags().String("threshold", "720h", "threshold for determining renewal necessity (e.g., 720h = 30 days)")
	checkCmd.Flags().Bool("force", false, "force domain renewal regardless of expiration status")
	checkCmd.Flags().BoolP("pretty", "p", false, "format JSON output for readability")
	checkCmd.Flags().StringP("output", "o", "", "file path to save the results")

	viper.BindPFlags(checkCmd.Flags())
}

func validateInput(domains []string, isCheckAllMode bool) error {
	if cfg.Lite && isCheckAllMode {
		return fmt.Errorf("lite mode doesn't support check all mode")
	}

	for _, arg := range domains {
		if !strings.HasSuffix(arg, ".ton") {
			return fmt.Errorf("invalid domain: %s", arg)
		}
	}
	return nil
}

func processResults(cmd *cobra.Command, results []check.Result) error {
	var jsonResults []byte
	var err error

	if cfg.Pretty {
		jsonResults, err = json.MarshalIndent(results, "", "  ")
	} else {
		jsonResults, err = json.Marshal(results)
	}

	if err != nil {
		return err
	}

	if cfg.Output == "" {
		cmd.OutOrStdout().Write(jsonResults)
		return nil
	}

	return os.WriteFile(cfg.Output, jsonResults, 0644)
}

func checkDomains(checker *check.Checker, domains []string) []check.Result {
	var results []check.Result
	var mutex sync.Mutex
	wg := sync.WaitGroup{}

	checkMethod := checker.Check
	if !cfg.Lite {
		checkMethod = checker.SmartCheck
	}

	for _, domain := range domains {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			result := checkMethod(d)

			mutex.Lock()
			results = append(results, result)
			mutex.Unlock()
		}(domain)
	}
	wg.Wait()

	return results
}
