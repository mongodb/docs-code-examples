# MongoDB Atlas Architecture Center Go SDK Code Examples

This repository contains [Atlas Go SDK](https://www.mongodb.com/docs/atlas/sdk/)
code examples that follow recommendations in MongoDB's official
[Atlas Architecture Center documentation](https://www.mongodb.com/docs/atlas/architecture/current/).
You can run, download, and modify these code examples as starting points for
configuring your MongoDB Atlas architecture for your use case.

## Overview

### Project Structure

```text
Project Root  
├── cmd  
│   ├── get_logs/main.go  
│   ├── get_metrics_disk/main.go  
│   ├── get_metrics_process/main.go  
├── internal  
│   ├── auth  
│   │   ├── auth.go  
│   ├── logs  
│   │   ├── downloader.go  
│   ├── metrics  
│   │   ├── metrics.go  
├── go.mod  
├── go.sum  
├── configs  
│   ├── config.json 
├── .env        #             
```

> NOTE: In a production environment, you are likely to use 

note in the README that you'll most likely be using a secrets manager in prod
## License

This project is licensed under the [Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0).

## Issues

To report an issue with any of these code examples, please leave feedback
through the corresponding documentation page in the
[MongoDB Atlas Architecture Center](https://www.mongodb.com/docs/atlas/architecture/current/).
Using the `Rate This Page` button, you can add a comment about the issue after
leaving a star rating.

## Contributing

We are not currently accepting public contributions to this repository at this
time.****
