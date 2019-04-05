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
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

// trackCmd represents the track command
var trackCmd = &cobra.Command{
	Use:   "track",
	Short: "track",
	Long:  `track.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return errors.New("track require no args")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		timer := time.NewTimer(time.Duration(timeout) * time.Second)
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		envs := os.Environ()

	L:
		for {
			select {
			case <-timer.C:
				fmt.Printf("%s\n", "timeout")
				break L
			case <-ticker.C:
				c := exec.CommandContext(ctx, "sh", "-c", command)
				c.Env = envs
				out, err := c.Output()
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
					break L
				}
				fmt.Printf("%s", out)
			case <-ctx.Done():
				break L
			}
		}
		fmt.Printf("%s\n", "end")
	},
}

func init() {
	rootCmd.AddCommand(trackCmd)
	trackCmd.Flags().StringVarP(&command, "command", "c", "", "command tor track.")
	trackCmd.Flags().IntVarP(&interval, "interval", "n", 5, "execution interval of tracking command. (seconds)")
	trackCmd.Flags().IntVarP(&timeout, "timeout", "d", 60*60*24, "execution duration of tracking command. (seconds)")
}
