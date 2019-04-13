# Sheer Heart Attack [![GitHub release](https://img.shields.io/github/release/k1LoW/sheer-heart-attack.svg)](https://github.com/k1LoW/sheer-heart-attack/releases)

`sheer-heart-attack` is a debugging tool that can execute any command on the process/host metrics trigger.

![screencast](screencast.svg)

## Features

- Easy to use.
- Track the process or host metrics. and execute specified command when threshold is exceeded.
- Record the STDOUT and STDERR of the executed command in the structured log.
- Slack notification.

## Quick Start

> This is the recommended usage.

``` console
root@kilr_q:~# source <(curl -sL https://git.io/sheer-heart-attack)
You can use `sheer-heart-attack` command in this session.
root@kilr_q:~# sheer-heart-attack launch
```

### In the case of fish :fish:

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

## Alternatives

- [ProcDump for Linux](https://github.com/Microsoft/ProcDump-for-Linux) - A Linux version of the ProcDump Sysinternals tool
