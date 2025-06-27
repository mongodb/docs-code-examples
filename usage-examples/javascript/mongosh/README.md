# mongosh Code Example Test PoC

This is a PoC to explore testing mongosh code examples for MongoDB documentation.

The structure of this JavaScript project is as follows:

- `/examples`: This directory contains example code and output to validate.
- `/tests`: This directory contains the test infrastructure to actually run
  the tests by invoking the example code.

While running the tests, this test suite creates files in a `/temp` directory. We concatenate
the code example with the necessary commands to connect to the deployment and use the correct
database.

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

```shell
CONTAINER ID   IMAGE   COMMAND                  CREATED        STATUS                  PORTS                        NAMES
ef70cce38f26   mongo   "/usr/local/bin/runn…"   29 hours ago   Up 29 hours (healthy)   127.0.0.1:63201->27017/tcp   mongodb-test
```

You may note the actual port is different than `27017`, if something else is already running on
`27017` on your machine. Note the port next to the IP address for running the tests.

### Create a .env file

Create a file named `.env` at the root of the `/mongosh` directory.
Add the following values to your .env file, substituting the port where your local deployment
is running:

```
CONNECTION_STRING="mongodb://localhost:63201"
CONNECTION_PORT="63201"
```

### Install the dependencies

This test suite requires you to have `Node.js` v20 or newer installed. If you
do not yet have Node installed, refer to
[the Node.js installation page](https://nodejs.org/en/download/package-manager)
for details.

From the root of the `/mongosh` directory, run the following command to install
dependencies:

```
npm install
```

### Run the tests

You can run tests from the command line or through your IDE. 

#### Run All Tests from the command line

From the `/mongosh` directory, run:

```
npm test
```

This invokes the following command from the `package.json` `test` key:

```
jest --runInBand --detectOpenHandles
```

In the above command:

- `jest` is the command to run the test suite
- `--runInBand` is a flag that specifies only running one test at a time
  to avoid collisions when creating/editing/dropping indexes. Otherwise, Jest
  defaults to running tests in parallel.
- `--detectOpenHandles` is a flag that tells Jest to track resource handles or async
  operations that remain open after the tests are complete. These can cause the test suite
  to hang, and this flag tells Jest to report info about these instances.

#### Run Test Suites from the command line

You can run all the tests in a given test suite (file).

From the `/tests` directory, run:

```
export $(xargs < ../.env) && jest -t '<text string from the 'describe' block you want to run>' --runInBand
```

In the above command:

- `export $(xargs < ../.env)` is Linux flag to make the contents of the `.env`
  file available to the test suite
- `jest` is the command to run the test suite
- `-t '<text string from the 'describe' block you want to run>'` is the way to
  specify to run all tests in test suite, which in this test, is a single file
- `--runInBand` is a flag that specifies only running one test at a time
  to avoid collisions when creating/editing/dropping indexes. Otherwise, Jest
  defaults to running tests in parallel.

#### Run Individual Tests from the command line

You can run a single test within a given test suite (file).

From the `/tests` directory, run:

```
export $(xargs < ../.env) && jest -t '<text string from the 'it' block of the test you want to run>'
```

In the above command:

- `export $(xargs < ../.env)` is Linux flag to make the contents of the `.env`
  file available to the test suite
- `jest` is the command to run the test suite
- `-t '<text string from the 'it' block of the test you want to run>'` is the
  way to specify to run a single test matching your text

Since you are only running a single test, there is no chance of colliding
with the other tests, so the `--runInBand` flag isn't needed.
