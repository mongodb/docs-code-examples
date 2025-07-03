# C# Driver Code Example Test PoC

This is a PoC to explore testing C# Driver code examples for MongoDB documentation.

The structure of this C# project is as follows:

- `driver.sln`: This solution contains the following projects:
  - `Examples`: This project contains example code and output to validate.
  - `Tests`: This project contains the test infrastructure to actually run
    the tests by invoking the example code.

## To run the tests locally

### Create a MongoDB Docker Deployment

To run these tests locally, you need a MongoDB Docker deployment. Make sure Docker is
running, and then:

1. Pull the MongoDB image from Docker Hub:

   ```shell
   docker pull mongo
   ```
2. Run the container:

   ```shell
   docker run --name mongodb-test -d -p 27017:27017 mongo
   ```

3. Verify the container is running:

   ```shell
   docker ps  
   ```

   The output resembles:

   ```text
   CONTAINER ID   IMAGE   COMMAND                  CREATED        STATUS                  PORTS                        NAMES
   ef70cce38f26   mongo   "/usr/local/bin/runn…"   29 hours ago   Up 29 hours (healthy)   127.0.0.1:63201->27017/tcp   mongodb-test
   ```

   You may note the actual port is different than `27017`, if something else is already running on
   `27017` on your machine. Note the port next to the IP address for running the tests. Alternately, you can get just
   the port info for your container using the following command:

   ```shell
   docker port mongodb-test
   ```
   
   The output resembles:

   ```text
   27017/tcp -> 0.0.0.0:27017
   27017/tcp -> [::]:27017
   ```

### Create a .env file

Create a file named `.env` at the root of the `/driver` directory.
Add the following values to your .env file, similar to the following example:

```
CONNECTION_STRING="mongodb://localhost:27017"
SOLUTION_ROOT="/Users/dachary.carey/workspace/docs-code-examples/usage-examples/csharp/driver/"
```

- `CONNECTION_STRING`: replace the port with the port where your local deployment is running.
- `SOLUTION_ROOT`: insert the path to the `driver` directory on your filesystem.

### Install the dependencies

This test suite requires you to have `.NET` v9.0 installed. If you
do not yet have .NET installed, refer to
[the .NET installation page](https://learn.microsoft.com/en-us/dotnet/core/install/macos)
for details.

From the root of each project directory, run the following command to install
dependencies:

```
dotnet restore
```

### Run the tests

You can run tests from the command line or through your IDE. 

#### Run All Tests from the command line

From the `/drivers` directory, run:

```
dotnet test
```

#### Run Individual Tests from the command line

You can run a single test within a given test suite (file).

From the `/drivers` directory, run:

```
dotnet test --filter "FullyQualifiedName=YourNamespace.YourTestClass.YourTestMethod"  
```
