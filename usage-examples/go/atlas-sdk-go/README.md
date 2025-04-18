[//]: # (.. Don't Copy to Target Repo )
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

## Set up 

### Set environment variables and config file

1. Create a `.env` file in the root directory with the following environment variables:
    ```shell
    MONGODB_ATLAS_SERVICE_ACCOUNT_ID=your-service-account-id
    MONGODB_ATLAS_SERVICE_ACCOUNT_SECRET=your-service-account-secret
    ```
2. Update the placeholders in the `configs/config.json` file with your Atlas information:
    ```json
   {
    
    "ATLAS_BASE_URL": "https://cloud.mongodb.com", 
    "ATLAS_ORG_ID": "<your-organization-id>",
    "ATLAS_PROJECT_ID": "<your-project-id>",
    "ATLAS_CLUSTER_NAME": "Cluster0",
    "ATLAS_HOST_NAME": "cluster0-shard-00-00.ab1cd.mongodb.net",
    "ATLAS_PORT": "27017",
    "ATLAS_PROCESS_ID": "cluster0-shard-00-00.ab1cd.mongodb.net:27017"
    
   }
    ```
    > **NOTE: Group ID == Project ID** Groups and projects are synonymous terms. Groups and projects are synonymous terms. Your group id is the same as your project id. 
3. 


---
scratchpad:
- manually test against real infra
- pull the real data

PR Push > run on captured data
& scheduled job to validate? (or bump our sdk version along with the v release &
verify; fix any failing test)
- run the gh action locally
---
1. go sdk testing
2. atlas testing

we'd need to be notified when either change
we can test sdk automatically with the captured data

periodic manual validation step to test the infra side
- release cadence for sdk > gh
- release cadence for atlas/infra > server version? api update?
---
arch center docs versioning considerations
- how to keep arch center version in sync with code example version updates?
---
snippets in generated-examples > push to arch center repo to use in docs?
---
narrating out loud the next step?
