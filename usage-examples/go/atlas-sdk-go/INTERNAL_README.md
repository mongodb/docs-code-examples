# Atlas SDK for Go 

> NOTE: This is an internal-only file and refers to the internal project details.  
> The external project details are documented in the README.md file.

This project demonstrates how to script specific functionality using the Atlas
SDK for Go. Code examples are used in the Atlas Architecture Center docs, and
a sanitized copy of the project is available in a user-facing repo: 
https://github.com/mongodb/atlas-architecture-go-sdk.

## Project Structure
```text
.
├── cmd/                   # Self-contained, runnable examples by category
├── configs/               # Atlas details 
├── internal               # Shared utilities and helpers (NOTE: ANY TEST FILES ARE INTERNAL ONLY)
├── go.mod
├── CHANGELOG.md           # User-facing list of major project changes
│── README.md              # User-facing README for copied project
│── INTERNAL_README.md     # (NOTE: INTERNAL ONLY - DON'T COPY TO ARTIFACT REPO)
│── scripts/               # (NOTE: INTERNAL ONLY) snip and copy code examples
│   └── bluehawk.sh 
├── .gitignore             
└── .env.example           
```

## Adding Examples
To add examples to the project: 



## Runnable Scripts
You can run individual scripts from the terminal. For example, to run `get_logs/main.go`:
```shell
go run cmd/get_logs/main.go
```

## Set up 

### Prerequisites

Contact the Developer Docs team with any setup questions or issues.

- A [service account](https://www.mongodb.com/docs/atlas/api/service-accounts-overview/#std-label-service-accounts-overview) with access to your Atlas project. 

> **NOTE:** Some scripts require an M10+ cluster

### Set environment variables and config file

1. Set the following variable values, either as a `.env` file in the root directory or through your IDE:
   ```dotenv
   MONGODB_ATLAS_SERVICE_ACCOUNT_ID=your_service_account_id
   MONGODB_ATLAS_SERVICE_ACCOUNT_SECRET=your_service_account_secret
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

... TODO

## Generate Examples

This project uses [Bluehawk](https://mongodb-university.github.io/Bluehawk/) markup to generate code examples from 
the source code.

We generate two types of code examples, intended for different use and destination: code snippets and full copied files

- To generate a new code snippet to use directly in the docs, use Bluehawk's [`snippet` tag](https://mongodb-university.github.io/Bluehawk/reference/tags#snippet) to mark the snippet content, then run the following script:
  ```bash
  ./scripts/bluehawk.sh snip
  ```
  Generated snippets output to `generated-usage-examples/go/atlas-sdk-go`
- To copy an entire file for the user-facing artifact repo, run the following script:
  ```bash
  ./scripts/bluehawk.sh copy
  ```
  Copied files output to `generated-usage-examples/go/atlas-sdk-go/project-copy` in their original directory structure.


  ```shell
   ./scripts/bluehawk.sh
   ```

> **TIP: "Copy" State** This project uses a state named "copy" specifically for any manipulations needed for outputted copied files intended for the artifact repo. See `cmd/get_linked_orgs/main.go` for an example of using the "copy" state to remove lines from the outputted file. Refer to Bluehawk's [`state` tag](https://mongodb-university.github.io/Bluehawk/reference/tags#state) docs for more info.
