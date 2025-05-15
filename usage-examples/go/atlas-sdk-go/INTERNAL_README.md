# Atlas SDK for Go

This project demonstrates how to script specific functionality using the Atlas
SDK for Go. Code examples are used in the Atlas Architecture Center docs, and
the project is made available in a user-facing repo.

## Project Structure
```text
atlas-sdk-go/
│── cmd/                  # Self-contained, runnable scripts
│   ├── get_logs/  
│       ├── main.go             
│   ├── get_metrics_disk/
│       ├── main.go
│   ├── get_metrics_process/
│       ├── main.go
│── config/                # Atlas configuration settings
│   ├── config.json             
│── internal/              # Shared internal logic
│   ├── auth/
|       ├── client.go
│   ├── config/
|       ├── json.go
|       ├── secrets.go
|       ├── loader.go
│── .env                   # Secrets file (excluded from Git)
│── go.mod                     
│── go.sum                           
│── README.md              # Internal-only README (do not copy with Copier Tool)
│── scripts/               # Internal-only Bluehawk scripts to snip and copy code examples
│   ├── bluehawk.sh
```

## Runnable Scripts
You can run individual scripts from the terminal. For example, to run `get_logs/main.go`:
```shell
go run cmd/get_logs/main.go
```

## Set up 

### Prerequisites

- A [service account](https://www.mongodb.com/docs/atlas/api/service-accounts-overview/#std-label-service-accounts-overview) with access to your Atlas project

> **NOTE:** Some scripts require an M10+ cluster

### Set environment variables and config file

1. Set the following variable values, either as a `.env` file in the root directory or through your IDE:
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

## Write Tests

# TODO

## Generate Examples

This project uses Bluehawk to generate code examples from the source code.

- Usage examples for the docs. These are generated using the `bluehawk snip`
  command based on the `snippet` markup in the code file.
- Full project files for the user-facing project repo. These are generated using
  the `bluehawk copy` command.
  
Run the bluehawk script and enter either `snip` or `copy`. The selected command
runs with the defined defaults.

  ```shell
   ./scripts/bluehawk.sh
   ```

> **NOTE: "Copy" State** This project uses a state named "copy" specifically for any manipulations needed for code copied to the artifact repo. 

