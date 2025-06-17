# Atlas Go SDK Examples for the MongoDB Atlas Architecture Center

This repository contains runnable examples for the
[Atlas Go SDK](https://www.mongodb.com/docs/atlas/sdk/)
that align with best practices from the MongoDB
[Atlas Architecture Center](https://www.mongodb.com/docs/atlas/architecture/current/).

Use these examples as starting points for your own Atlas integration.

## Features

Currently, the repository includes examples that demonstrate the following:

- Authenticate with service accounts
- Return cluster and database metrics
- Download logs for a specific host
- Return all linked organizations from a specific billing organization
- Programmatically manage Atlas resources

As the Architecture Center documentation evolves, this repository will be updated with new examples 
and improvements to existing code. 

## Project Structure

```text
.
├── cmd                  # Runnable examples by category
│   ├── get_linked_orgs/main.go
│   ├── get_logs/main.go
│   ├── get_metrics_disk/main.go
│   └── get_metrics_process/main.go
├── configs              # Atlas configuration template
│   └── config.json
├── internal             # Shared utilities and helpers
│   ├── auth/
│   ├── config/
│   ├── logs/
│   ├── metrics/
│   └── utils.go
├── go.mod
├── go.sum
├── CHANGELOG.md         # List of major changes to the project 
├── .gitignore           # Ignores .env file and log output
└── .env.example         # Example environment variables
```

## Prerequisites

- Go 1.16 or later
- A MongoDB Atlas project and cluster
- Service account credentials with appropriate permissions. See
  [Service Account Overview](https://www.mongodb.com/docs/atlas/api/service-accounts-overview/).

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/mongodb/atlas-architecture-go-sdk.git
   cd atlas-architecture-go-sdk
   ```

2. Create a `.env` file in the root directory with your MongoDB Atlas service account credentials:
   ```dotenv
   MONGODB_ATLAS_SERVICE_ACCOUNT_ID=your_service_account_id
   MONGODB_ATLAS_SERVICE_ACCOUNT_SECRET=your_service_account_secret
   ```
   > **NOTE:** For production, use a secrets manager (e.g. HashiCorp Vault, AWS Secrets Manager) 
   > instead of environment variables. 
   > See [Secrets management](https://www.mongodb.com/docs/atlas/architecture/current/auth/#secrets-management).

3. Configure Atlas details in `configs/config.json`:
   ```json
   {
     "MONGODB_ATLAS_BASE_URL": "<optional-base-url>",
     "ATLAS_ORG_ID": "<your-organization-id>",
     "ATLAS_PROJECT_ID": "<your-project-id>",
     "ATLAS_CLUSTER_NAME": "<your-cluster-name>",
     "ATLAS_PROCESS_ID": "<cluster-name-shard-00-00.hostsuffix.mongodb.net:port>"
   }
   ```
   > **NOTE:** The base URL defaults to `https://cloud.mongodb.com` if not specified.

4. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running Examples

Examples in this project are intended to be run as individual scripts. 
You can also adjust them to suit your needs:

- Modify time ranges
- Add filtering parameters
- Change output formats

### Billing
#### Get All Linked Organizations 
```bash
go run cmd/get_linked_orgs/main.go
```

### Logs
Logs output to `./logs` as `.gz` and `.txt`.

#### Fetch All Host Logs
```bash
go run cmd/get_logs/main.go
```

### Metrics
Metrics print to the console.

#### Get Disk Measurements
```bash
go run cmd/get_metrics_disk/main.go
```

#### Get Cluster Metrics
```bash
go run cmd/get_metrics_process/main.go
```

## Changelog

For list of major changes to this project, see [CHANGELOG](CHANGELOG.md).

## Reporting Issues

Use the "Rate this page" widget on the
[Atlas Architecture Center](https://www.mongodb.com/docs/atlas/architecture/current/)
docs to leave feedback or file issues.

## License

This project is licensed under Apache 2.0. See [LICENSE](LICENSE.md).
