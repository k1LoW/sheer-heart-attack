// Copyright © 2019 Ken'ichiro Oyama <k1lowxb@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/k1LoW/exec"
	"github.com/k1LoW/sheer-heart-attack/options"

	"github.com/labstack/gommon/color"
	"github.com/spf13/cobra"
)

var (
	nonInteractive bool
	lang           string
)

// launchCmd represents the launch command
var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch 'track' command in background.",
	Long:  `Launch 'track' command in background.`,
	Run: func(cmd *cobra.Command, args []string) {
		exe, err := os.Executable()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		if lang == "" {
			lang = os.Getenv("LANG")
		}
		o, err := options.NewOptions(
			pid,
			name,
			threshold,
			interval,
			attempts,
			command,
			times,
			timeout,
			slackChannel,
			slackMention,
			nonInteractive,
			lang,
		)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		trackCommand := []string{exe, "track"}
		trackCommand = append(trackCommand, o.Get()...)

		envs := os.Environ()
		c := exec.Command(trackCommand[0], trackCommand[1:]...)
		if pid > 0 {
			envs = append(envs, fmt.Sprintf("PID=%d", pid))
		}

		c.Env = envs
		err = c.Start()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s %s\n", color.Magenta("Launched:", color.B), strings.Join(trackCommand, " "))
	},
}

func init() {
	rootCmd.AddCommand(launchCmd)
	launchCmd.Flags().BoolVarP(&nonInteractive, "non-interactive", "", false, "Disables all interactive prompting.")

	launchCmd.Flags().Int32VarP(&pid, "pid", "", 0, "PID of the process")
	launchCmd.Flags().StringVarP(&name, "name", "", "", "name of the process")
	launchCmd.Flags().StringVarP(&threshold, "threshold", "", "cpu > 5 || mem > 10", "Threshold conditions")
	launchCmd.Flags().StringVarP(&interval, "interval", "", "5s", "Interval of checking if the threshold exceeded''")
	launchCmd.Flags().IntVarP(&attempts, "attempts", "", 1, "Maximum number of attempts continuously exceeding the threshold")
	launchCmd.Flags().StringVarP(&command, "command", "", "", "Command to execute when the maximum number of attempts is exceeded")
	launchCmd.Flags().IntVarP(&times, "times", "", 1, "Maximum number of command executions. If times < 1, track and execute until timeout")
	launchCmd.Flags().StringVarP(&timeout, "timeout", "", "1day", "Timeout of tracking''")
	launchCmd.Flags().StringVarP(&slackChannel, "slack-channel", "", "", "Slack channel to notify")
	launchCmd.Flags().StringVarP(&slackMention, "slack-mention", "", "", "Slack mention")
	launchCmd.Flags().StringVarP(&lang, "lang", "", "", "Language")
}
