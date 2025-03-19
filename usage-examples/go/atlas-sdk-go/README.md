```text
atlas-sdk-go/
│── cmd/                        # ✅ Self-contained, runnable scripts
│   ├── get_logs/              # Script to fetch invoices
│       ├── main.go                # Script to retrieve logs
│   ├── get_metrics/              # Script to fetch invoices
│       ├── main.go            # Script to delete data
│── config/                          # ✅ Configuration directory
│   ├── config.json                  # General app settings (non-sensitive)
│── internal/                        # ✅ Shared internal logic
│   ├── config_loader.go              # Loads JSON configs
│   ├── secrets_loader.go             # Loads secrets securely
│── .gitignore                        # Exclude secrets.json
│── go.mod                            # Go module file
│── go.sum                            # Dependency checksums
│── README.md                         # Documentation


```shell
./bluehawk.sh
```

```shell
./run_cmd.sh get_logs
```
