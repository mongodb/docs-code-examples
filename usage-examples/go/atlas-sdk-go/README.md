# Atlas SDK for Go

## Project Structure
```text
atlas-sdk-go/
│── cmd/                  # Self-contained, runnable scripts
│   ├── get_logs/  
│       ├── main.go             
│   ├── get_metrics/            
│       ├── main.go             
│── config/                # Atlas configuration settings
│   ├── config.json             
│── internal/              # Shared internal logic
    │── auth/              
        ├── auth.go                   
│   ├── api_client.go  
│   ├── config_loader.go        
│   ├── secrets_loader.go       
│── .env                   # Secrets file (excluded from Git)
│── go.mod                     
│── go.sum                           
│── README.md                       
```

## Runnable Scripts
You can run individual scripts using `run_cmd.sh` and specifying the script's action (i.e. the parent directory for the `main.go` you want to run). 

For example, to run `get_logs/main.go`:
```shell
./run_cmd.sh get_logs
```
