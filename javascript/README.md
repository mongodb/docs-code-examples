# JavaScript Vector Search Code Example Test PoC

This is a PoC to explore testing JavaScript Vector Search code examples for
MongoDB documentation.

The structure of this JavaScript project is as follows:

- `/example`: This directory contains example code, marked up with Bluehawk,
  that will be outputted to the `/generated` directory when we run the Bluehawk
  script.
- `/tests`: This directory contains the test infrastructure to actually run
  the tests by invoking the example code.

## To run the tests locally

### Create an Atlas cluster

To run these tests locally, you need an Atlas cluster with the sample data set
loaded and no search indexes. For best results, create a fresh M0 cluster, add
sample data, and save the connection string for use in the next step.

### Create a .env file

Create a file named '.env' at the root of the '/javascript' directory within
this project. Add your Atlas connection string as an environment value named
`ATLAS_CONNECTION_STRING`:

```
ATLAS_CONNECTION_STRING="<your-connection-string>"
```

Replace the `<your-connection-string>` placeholder with the connection
string from the Atlas cluster you created in the prior step.

Add an `ENV` environment value whose value is `"Atlas"`. This denotes that
you are running tests against Atlas instead of a local instance. Some functionality
is not supported in local deployment, and some query results vary between the
two environments, so specifying the environment gives the test suite info about
which tests to run and which outputs to expect.

```
ENV="Atlas"
```

### Install the dependencies

This test suite requires you to have `Node.js` v20 or newer installed. If you
do not yet have Node installed, refer to
[the Node.js installation page](https://nodejs.org/en/download/package-manager)
for details.

From the root of the `/javascript` directory, run the following command to install
dependencies:

```
npm install
```

### Run the tests

#### Run Tests from the IDE

Normally, you could press the play button next to a test name to run a test
in the IDE. Because this test suite relies on an Atlas connection string and
environment value passed in from the environment, running tests in the IDE
will fail unless you configure the IDE with the appropriate environment
variables.

In JetBrains IDEs, you can do the following:

- Click the play button next to the test suite name
- Select the `Modify Run Configuration` option
- In the `Environment Variables` field, supply the appropriate environment variables
  - Note: you do not need to use quotes around the connection string in this field
    i.e. it should resemble:
    ATLAS_CONNECTION_STRING=mongodb+srv://your-connection-string

#### Run All Tests from the command line

From the `/javascript` directory, run:

```
npm test
```

This invokes the following command from the `package.json` `test` key:

```
export $(xargs < .env) && jest  --run-in-band
```

In the above command:

- `export $(xargs < .env)` is Linux flag to make the contents of the `.env`
  file available to the test suite
- `jest` is the command to run the test suite
- `--runInBand` is a flag that specifies only running one test at a time
  to avoid collisions when creating/editing/dropping indexes. Otherwise, Jest
  defaults to running tests in parallel.

#### Run Test Suites from the command line

You can run all the tests in a given test suite (file).

From the `/tests` directory, run:

```
export $(xargs < .env) && jest -t '<text string from the 'describe' block you want to run>' --runInBand
```

In the above command:

- `export $(xargs < .env)` is Linux flag to make the contents of the `.env`
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
export $(xargs < .env) && jest -t '<text string from the 'it' block of the test you want to run>'
```

In the above command:

- `export $(xargs < .env)` is Linux flag to make the contents of the `.env`
  file available to the test suite
- `jest` is the command to run the test suite
- `-t '<text string from the 'it' block of the test you want to run>'` is the
  way to specify to run a single test matching your text

Since you are only running a single test, there is no chance of colliding
with the other tests, so the `--runInBand` flag isn't needed.

## To run the tests in CI

Two GitHub workflows run these tests in CI automatically when you change any
files in the `examples` directory:

- `.github/workflows/test-javascript-examples-in-atlas.yml`
- `.github/workflows/test-javascript-examples-in-docker.yml`

GitHub reports the results as passing or failing checks on any PR that changes
an example.

If changing an example causes its test to fail, this should be considered
blocking to merge the example.

If changing an example causes an _unrelated_ test to fail, create a Jira ticket
to fix the unrelated test, but this should not block merging an example update.

## To generate tested code examples

This test suite uses [Bluehawk](https://github.com/mongodb-university/Bluehawk)
to generate code examples from the test files. Bluehawk contains functionality
to replace content that we do not want to show verbatim to users, remove test
functionality from the outputted code examples, etc.

Install Bluehawk, and then run the Bluehawk script at the root of the `/javascript`
directory to generate updated example files:

```
./bluehawk.sh
```
