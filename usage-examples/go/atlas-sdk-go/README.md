# Atlas SDK for Go

## Project Structure
```text
atlas-sdk-go/
│── cmd/                        # Self-contained, runnable scripts
│   ├── get_logs/  
│       ├── main.go             # Get Host logs
│   ├── get_metrics/            
│       ├── main.go             # Get Process and Disk metrics
│── config/                    
│   ├── config.json             # Atlas configuration settings
│── internal/                   
│   ├── config_loader.go        # Loads JSON configs
│   ├── secrets_loader.go       # Loads .env securely
│── .env                        # Secrets file with API keys, database credentials
│── go.mod                     
│── go.sum                           
│── README.md                       
```

## Runnable Scripts
Run a specific script using `run_cmd.sh`. For example, to run `get_logs/main.go`:
```shell
./run_cmd.sh get_logs
```
