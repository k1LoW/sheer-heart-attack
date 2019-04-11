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
	"os/exec"
	"strconv"
	"strings"

	"github.com/labstack/gommon/color"
	"github.com/spf13/cobra"
)

var nonInteractive bool

// launchCmd represents the launch command
var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch 'track' command in background.",
	Long:  `Launch 'track' command in background.`,
	Run: func(cmd *cobra.Command, args []string) {
		envs := os.Environ()
		exe, err := os.Executable()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		trackCommand := []string{exe, "track"}

		// pid
		optPID, err := optionPID(pid, nonInteractive)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		pidStr := optPID[1]
		trackCommand = append(trackCommand, optPID...)
		// threshold
		pidInt32, err := strconv.ParseInt(pidStr, 10, 32)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		optThreshold, err := optionThreshold(threshold, int32(pidInt32), nonInteractive)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		trackCommand = append(trackCommand, optThreshold...)
		// interval
		optInterval, err := optionInterval(interval, nonInteractive)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		trackCommand = append(trackCommand, optInterval...)
		// attempts
		optAttempts, err := optionAttempts(attempts, nonInteractive)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		trackCommand = append(trackCommand, optAttempts...)
		// command
		optCommand, err := optionCommand(command, nonInteractive)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		trackCommand = append(trackCommand, optCommand...)
		// times
		optTimes, err := optionTimes(times, nonInteractive)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		trackCommand = append(trackCommand, optTimes...)
		// timeout
		optTimeout, err := optionTimeout(timeout, nonInteractive)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		trackCommand = append(trackCommand, optTimeout...)

		c := exec.Command(trackCommand[0], trackCommand[1:]...)
		envs = append(envs, "PID=%s", pidStr)
		c.Env = envs
		c.Start()

		fmt.Printf("%s %s\n", color.Magenta("Launched:", color.B), strings.Join(trackCommand, " "))
	},
}

func init() {
	rootCmd.AddCommand(launchCmd)
	launchCmd.Flags().BoolVarP(&nonInteractive, "non-interactive", "", false, "Disables all interactive prompting.")

	launchCmd.Flags().Int32VarP(&pid, "pid", "", 0, "PID of the process")
	launchCmd.Flags().StringVarP(&threshold, "threshold", "", "cpu > 5 || mem > 10", "Threshold conditions")
	launchCmd.Flags().IntVarP(&interval, "interval", "", 5, "Interval of checking if the threshold exceeded (seconds)")
	launchCmd.Flags().IntVarP(&attempts, "attempts", "", 1, "Maximum number of attempts continuously exceeding the threshold")
	launchCmd.Flags().StringVarP(&command, "command", "", "", "Command to execute when the maximum number of attempts is exceeded")
	launchCmd.Flags().IntVarP(&times, "times", "", 1, "Maximum number of command executions. If times < 1, track and execute until timeout")
	launchCmd.Flags().IntVarP(&timeout, "timeout", "", 60*60*24, "Timeout of tracking (seconds)")
}
