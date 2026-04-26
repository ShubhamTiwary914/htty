<div align="center">

<img src="./assets/htty.svg" alt="htty" width="480"/>

<br/><br/>

keyboard-driven HTTP client that lives in your terminal.

<br/>

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)
![License](https://img.shields.io/github/license/ShubhamTiwary914/htty?style=flat-square)
![Release](https://img.shields.io/github/v/release/ShubhamTiwary914/htty?style=flat-square)

</div>

---

#### Preview:
<img width="1880" height="1089" alt="image" src="https://github.com/user-attachments/assets/9676f4ef-0aa6-462f-9089-0fd7d35a5497" />


## Install

```bash
curl -fsSL https://raw.githubusercontent.com/ShubhamTiwary914/htty/master/setup/setup.sh | bash
```

Requires: `curl`, `unzip`, `go`, `make`

---

## Usage

```
htty
```

Navigate panels with `Tab`. Set your method, URL, headers, and body — hit `Enter` to send. Response renders inline.

| Key | Action |
|---|---|
| `Tab` / `Shift+Tab` | Move between panels |
| `Enter` | Confirm / send request |
| `Ctrl+C` | Quit |

---

## Configuration

The install script sets up config at `~/.config/htty/config.json`. Override paths with environment variables:

| Variable | Default | Description |
|---|---|---|
| `CONFIG_FILE` | `~/.config/htty/config.json` | Config file path |
| `CACHE_PREFIX` | `~/.cache/htty` | Completions cache directory |
| `LOGLEVEL` | `info` | Log level (`debug`, `info`, `error`) |
| `LOGFILE` | `/var/log/htty/htty.log` | File for logging I/O |

---

## Build from source

```bash
git clone https://github.com/ShubhamTiwary914/htty
cd htty
make build
```

---

## Contributing

Bug reports and PRs welcome. Open an [issue](https://github.com/ShubhamTiwary914/htty/issues) first for anything beyond small fixes.

---
