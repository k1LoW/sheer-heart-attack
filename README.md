# Sheer Heart Attack [![GitHub release](https://img.shields.io/github/release/k1LoW/sheer-heart-attack.svg)](https://github.com/k1LoW/sheer-heart-attack/releases)

`sheer-heart-attack` is a debugging tool that can execute any command on process/host metrics trigger.

![screencast](screencast.svg)

## Features

- Easy to use (just execute `sheer-heart-attack launch`).
- Track process and/or host metrics. and execute specified command when threshold is exceeded.
- Record the STDOUT and STDERR of the executed command in the structured log.
- Slack notification.

## Quick Start

> This is the recommended usage.

``` console
root@kilr_q:~# source <(curl -sL https://git.io/sheer-heart-attack)
You can use `sheer-heart-attack` command in this session.
root@kilr_q:~# sheer-heart-attack launch
```

**In the case of fish :fish:**

``` console
root@kilr_q:~# curl -sL https://git.io/sheer-heart-attack-fish | source
You can use `sheer-heart-attack` command in this session.
root@kilr_q:~# sheer-heart-attack launch
```

## Install

**manually:**

Download binany from [releases page](https://github.com/k1LoW/sheer-heart-attack/releases)

**go get:**

``` console
go get github.com/k1LoW/sheer-heart-attack
```

## Usage

Just execute `sheer-heart-attack launch`.

**Launch Options (Flags):**

| Option | Default | Purpose |
| --- | --- | --- |
| `pid` | | PID of the process. |
| `threshold` | `cpu > 5 \|\| mem > 10` | Threshold conditions. |
| `interval` | `5` | Interval of checking if the threshold exceeded (seconds). |
| `attempts` | `1` | Maximum number of attempts continuously exceeding the threshold. |
| `command` | | Command to execute when the maximum number of attempts is exceeded. |
| `times` | `1` | Maximum number of command executions. If times < 1, track and execute until timeout. |
| `timeout` | `86400` | Timeout of tracking (seconds). |
| `slack-channel` | | Slack channel to notify. |
| `slack-mention` | | Slack mention. (`@here` or user_id `@UXXXXXXXXX`) |

### Set `threshold` using operators

The following operators can be used to set the threshold:

`+`, `-`, `*`, `/`, `==`, `!=`, `<`, `>`, `<=`, `>=`, `not`, `and`, `or`, `!`, `&&`, `||`

For example, you can set the threshold as follows

- `cpu > 10 and mem > 20`
- `(user + system) > 50 || iowait > 50`
- `load1 > 5 or load15 > 2`

### Slack Notification

`sheer-heart-attack` find and use [Slack Incomming Webhook](https://api.slack.com/incoming-webhooks) URL via envirionment variables ( `SLACK_INCOMMING_WEBHOOK_URL`, `SLACK_WEBHOOK_URL`, `SLACK_URL` )

![slack](slack.png)

## Support Metrics

| Metric | |
| --- | --- |
| `proc_cpu` | Percentage of the CPU time the process uses (percent). |
| `proc_mem` | Percentage of the total RAM the process uses (percent). |
| `proc_rss` | Non-swapped physical memory the process uses (bytes). |
| `proc_vms` | Amount of virtual memory the process uses (bytes). |
| `proc_swap` | Amount of memory that has been swapped out to disk the process uses (bytes). |
| `proc_open_files` | Amount of files and file discripters opend by the process. **linux only** |
| `cpu` | Percentage of cpu used. |
| `mem` | Percentage of RAM used. |
| `swap` | Amount of memory that has been swapped out to disk (bytes). |
| `user` | Percentage of CPU utilization that occurred while executing at the user level. |
| `system` | Percentage of CPU utilization that occurred while executing at the system level. |
| `idle` | Percentage of time that CPUs were idle and the system did not have an outstanding disk I/O request. |
| `nice` | Percentage of CPU utilization that occurred while executing at the user level with nice priority. |
| `iowait` | Percentage of time that CPUs were idle during which the system had an outstanding disk I/O request. |
| `irq` | Percentage of time spent by CPUs to service hardware interrupts. |
| `softirq` | Percentage of time spent by CPUs to service software interrupts. |
| `steal` | Percentage of time spent in involuntary wait by the virtual CPUs while the hypervisor was servicing another virtual processor. |
| `guest` | Percentage of time spent by CPUs to run a virtual processor. |
| `guest_nice` | Percentage of time spent by CPUs to run a virtual processor with nice priority. |
| `load1` | Load avarage for 1 minute. |
| `load5` | Load avarage for 5 minutes. |
| `load15` | Load avarage for 15 minutes. |

## Alternatives

- [ProcDump for Linux](https://github.com/Microsoft/ProcDump-for-Linux) - A Linux version of the ProcDump Sysinternals tool
