[//]: # (.. Don't Copy to Target Repo )
# Atlas SDK for Go

This project demonstrates how to script specific functionality using the Atlas SDK for Go. Code examples are included in the Atlas Architecture Center docs and available in a user-facing version of the project. 

- Generated usage examples, which are included directly in the docs
- Within a 

## Project Structure
```text
atlas-sdk-go/
│── bluehawk/             # Bluehawk scripts to snip and copy code examples
│   ├── copy.sh  
│   ├── snip.sh  
│── cmd/                  # Self-contained, runnable scripts
│   ├── get_logs/  
│       ├── main.go             
│   ├── get_metrics/ 
│       ├── dev/ 
│           ├── main.go            
│       ├── prod/    
│           ├── main.go         
│── config/                # Atlas configuration settings
│   ├── config.json             
│── internal/              # Shared internal logic
    │── auth/              
        ├── auth.go                   
│   ├── config_loader.go        
│   ├── secrets_loader.go       
│── .env                   # Secrets file (excluded from Git)
│── go.mod                     
│── go.sum                           
│── README.md              # Internal-only README (do not copy)         
```

## Runnable Scripts
You can run individual scripts using `run_cmd.sh` and specifying the script's action (i.e. the parent directory for the `main.go` you want to run). 

For example, to run `get_logs/main.go`:
```shell
./run_cmd.sh get_logs
```

## Set up 

### Prerequisites

- A [service account](https://www.mongodb.com/docs/atlas/api/service-accounts-overview/#std-label-service-accounts-overview) with access to your Atlas project

### Set environment variables and config file

1. Set the following variable values, either as a `.env` file in the root directory or through your IDE
    ```shell
    MONGODB_ATLAS_SERVICE_ACCOUNT_ID=your-service-account-id
    MONGODB_ATLAS_SERVICE_ACCOUNT_SECRET=your-service-account-secret
    ```
2. Update the placeholders in the `configs/config.json` file with your Atlas cluster information:
    ```json
   {
    
    "ATLAS_BASE_URL": "https://cloud.mongodb.com", 
    "ATLAS_ORG_ID": "<your-organization-id>",
    "ATLAS_PROJECT_ID": "<your-project-id>",
    "ATLAS_CLUSTER_NAME": "Cluster0",
    "ATLAS_PROCESS_ID": "cluster0-shard-00-00.ab1cd.mongodb.net:27017"
    
   }
    ```
    > **NOTE: Group ID == Project ID** Groups and projects are synonymous terms. Groups and projects are synonymous terms. Your group id is the same as your project id. 
3. 

## Generate Examples

### Generate code usage examples for docs
This project uses the following Bluehawk commands to generate the code examples:

- Usage examples for the docs. These are generated using the `bluehawk snip` command based on the `snippet` markup in the code file.
  ```shell
   ./bluehawk/snip.sh
   ```
  
### Copy project files for user-facing project repo

To copy the full project files for the user-facing artifact repo. These are generated using the `bluehawk copy` command, and any specified files are ignored.
  ```shell
   ./bluehawk/copy.sh
   ```

> **NOTE: "Copy" State** This project uses a state named "copy" specifically for any manipulations needed for code copied to the artifact repo. 

## Copy Generated Examples to Other Repos



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
