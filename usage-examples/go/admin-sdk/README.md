```text
my-go-project/
│── scripts/                        # ✅ Self-contained, runnable scripts
│   ├── get_invoices.go              # Script to fetch invoices
│   ├── fetch_logs.go                # Script to retrieve logs
│   ├── delete_old_data.go           # Script to delete data
│── config/                          # ✅ Configuration directory
│   ├── config.json                  # General app settings (non-sensitive)
│   ├── logging.json                  # Logging settings
│   ├── database.json                 # Database settings (no credentials)
│   ├── features.json                 # Feature flags and toggles
│── secrets/                         # ✅ Secure secrets directory (excluded from Git)
│   ├── secrets.json                   # 🔐 API keys, database credentials
│── internal/                        # ✅ Shared internal logic
│   ├── config_loader.go              # Loads JSON configs
│   ├── secrets_loader.go             # Loads secrets securely
│── .gitignore                        # Exclude secrets.json
│── go.mod                            # Go module file
│── go.sum                            # Dependency checksums
│── README.md                         # Documentation

```
