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
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/antonmedv/expr"
	"github.com/k1LoW/duration"
	"github.com/k1LoW/exec"
	"github.com/k1LoW/metr/metrics"
	"github.com/k1LoW/sheer-heart-attack/logger"
	"github.com/k1LoW/sheer-heart-attack/options"
	"github.com/mattn/go-isatty"
	"github.com/shirou/gopsutil/process"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	startMessage         = "Tracking start"
	timeoutMessage       = "Tracking timeout"
	executeMessage       = "Threshold exceeded and execute command"
	noExecuteMessage     = "Threshold exceeded ( no command )"
	endMessage           = "Tracking ended"
	errorMessage         = "Error"
	executeFailedMessage = "Failed to execute command"
)

var force bool

type trackField struct {
	key   string
	value interface{}
}

// trackCmd represents the track command
var trackCmd = &cobra.Command{
	Use:   "track",
	Short: "Track the process metrics and execute command when the threshold is exceeded.",
	Long:  `Track the process metrics and execute command when the threshold is exceeded.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return errors.New("track require no args")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if isatty.IsTerminal(os.Stdout.Fd()) && !force {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "can not execute `track` directly. execute via `launch`, or use `--force` option")
			os.Exit(1)
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		numRe := regexp.MustCompile(`^\d+$`)
		if numRe.MatchString(timeout) {
			timeout = timeout + "s"
		}
		timeoutDuration, err := duration.Parse(timeout)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		timer := time.NewTimer(timeoutDuration)
		if numRe.MatchString(interval) {
			interval = interval + "s"
		}
		intervalDuration, err := duration.Parse(interval)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		ticker := time.NewTicker(intervalDuration)
		envs := os.Environ()
		executionTimeout := intervalDuration * 3
		logPath, err := filepath.Abs(fmt.Sprintf("sheer-heart-attack-%s.log", time.Now().Format(time.RFC3339)))
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		l := logger.NewLogger(logPath)
		if pid > 0 {
			p, err := process.NewProcess(pid)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			name, err = p.Name()
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		}
		hostname, err := os.Hostname()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		trackFields := []trackField{
			{"pid", fmt.Sprintf("%d", pid)},
			{"name", name},
			{"threshold", threshold},
			{"interval", interval},
			{"attempts", attempts},
			{"times", times},
			{"timeout", timeout},
			{"slack-channel", slackChannel},
			{"command", commands},
			{"hostname", hostname},
			{"log-path", logPath},
		}

		lNoHook := l
		if slackChannel != "" {
			webhookURL, err := options.GetEnvSlackIncommingWebhook()
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			l = l.WithOptions(zap.Hooks(notifySlack(webhookURL, slackChannel, slackMention, trackFields)))
		}
		fields := []zap.Field{}
		for _, o := range trackFields {
			fields = append(fields, zap.Any(o.key, o.value))
		}
		exceeded := 0
		execution := 0
		l.Info(startMessage, fields...)
		sg := sync.WaitGroup{}

	L:
		for {
			select {
			case <-timer.C:
				l.Info(timeoutMessage)
				break L
			case <-ticker.C:
				var (
					m   *metrics.Metrics
					err error
				)
				switch {
				case pid > 0:
					m, err = metrics.GetMetrics(options.CollectInterval, pid)
					if err != nil {
						l.Error(errorMessage, zap.Error(err))
						break L
					}
				case name != "":
					m, err = metrics.GetMetricsByName(options.CollectInterval, name)
					if err != nil {
						l.Error(errorMessage, zap.Error(err))
						break L
					}
				default:
					m, err = metrics.GetMetrics(options.CollectInterval)
					if err != nil {
						l.Error(errorMessage, zap.Error(err))
						break L
					}
				}
				got, err := expr.Eval(fmt.Sprintf("(%s) == true", threshold), m.Raw())
				if err != nil {
					l.Error(errorMessage, zap.Error(err))
					break L
				}
				if got.(bool) {
					exceeded++
				} else {
					exceeded = 0
				}
				if times > 0 && execution >= times {
					break L
				}
				if exceeded >= attempts {
					sg.Add(1)
					go func(ctx context.Context) {
						execution++
						exceeded = 0
						fields := []zap.Field{}
						m.Each(func(metric metrics.Metric, value interface{}) {
							fields = append(fields, zap.Any(metric.Name, value))
						})
						if len(commands) > 0 {
							cmdSg := sync.WaitGroup{}
							hooked := false
							for _, c := range commands {
								cmdSg.Add(1)
								go func(ctx context.Context, c string) {
									stdout, stderr, err := execute(ctx, c, envs, executionTimeout)
									fields = []zap.Field{
										zap.String("command", c),
										zap.ByteString("stdout", stdout),
										zap.ByteString("stderr", stderr),
									}
									if hooked {
										lNoHook.Info(executeMessage, fields...)
									} else {
										hooked = true
										l.Info(executeMessage, fields...)
									}
									if err != nil {
										l.Error(executeFailedMessage, zap.Error(err))
										// do not break
									}
									cmdSg.Done()
								}(ctx, c)
							}
							cmdSg.Wait()
						} else {
							l.Info(noExecuteMessage, fields...)
						}
						sg.Done()
					}(ctx)
				}
			case <-ctx.Done():
				break L
			}
		}

		sg.Wait()
		l.Info(endMessage)
	},
}

func init() {
	rootCmd.AddCommand(trackCmd)
	trackCmd.Flags().Int32VarP(&pid, "pid", "", 0, "PID of the process")
	trackCmd.Flags().StringVarP(&name, "name", "", "", "name of the process")
	trackCmd.Flags().StringVarP(&threshold, "threshold", "", "cpu > 5 || mem > 10", "Threshold conditions")
	trackCmd.Flags().StringVarP(&interval, "interval", "", "5s", "Interval of checking if the threshold exceeded")
	trackCmd.Flags().IntVarP(&attempts, "attempts", "", 1, "Maximum number of attempts continuously exceeding the threshold")
	trackCmd.Flags().StringArrayVarP(&commands, "command", "", []string{}, "Command to execute when the maximum number of attempts is exceeded")
	trackCmd.Flags().IntVarP(&times, "times", "", 1, "Maximum number of command executions. If times < 1, track and execute until timeout")
	trackCmd.Flags().StringVarP(&timeout, "timeout", "", "1day", "Timeout of tracking")
	trackCmd.Flags().StringVarP(&slackChannel, "slack-channel", "", "", "Slack channel to notify")
	trackCmd.Flags().StringVarP(&slackMention, "slack-mention", "", "", "Slack mention")
	trackCmd.Flags().BoolVarP(&force, "force", "", false, "Force execute 'track' command on tty")
}

func execute(ctx context.Context, command string, envs []string, timeout time.Duration) ([]byte, []byte, error) {
	innerCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	c := exec.CommandContext(innerCtx, "sh", "-c", command)
	c.Env = envs
	c.Stdout = &stdout
	c.Stderr = &stderr
	if err := c.Start(); err != nil {
		return stdout.Bytes(), stderr.Bytes(), err
	}
	if err := c.Wait(); err != nil {
		return stdout.Bytes(), stderr.Bytes(), err
	}
	return stdout.Bytes(), stderr.Bytes(), nil
}
