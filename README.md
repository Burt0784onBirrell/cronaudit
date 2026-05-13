# cronaudit

> Parses and validates crontab entries across multiple hosts, flagging conflicts and deprecated syntax.

---

## Installation

```bash
go install github.com/youruser/cronaudit@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/cronaudit.git && cd cronaudit && go build ./...
```

---

## Usage

Run `cronaudit` against one or more crontab files:

```bash
cronaudit --hosts hosts.txt --crontabs /etc/cron.d/
```

Example output:

```
[WARN]  host-01: deprecated @reboot syntax in job "backup-task"
[ERROR] host-02: schedule conflict detected between "sync-job" and "report-job" (overlap at 02:00)
[OK]    host-03: all crontab entries valid
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--hosts` | File containing list of target hosts | `hosts.txt` |
| `--crontabs` | Path to crontab files or directory | `/etc/cron.d/` |
| `--strict` | Treat warnings as errors | `false` |
| `--output` | Output format: `text`, `json` | `text` |

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-change`)
3. Commit your changes and open a PR

---

## License

This project is licensed under the [MIT License](LICENSE).