# Go Vector Search Code Example Test PoC

This is a PoC to explore testing Python Vector Search code examples for MongoDB
documentation.

The structure of this Python project is as follows:

- `/examples`: This directory contains example code, marked up with Bluehawk,
  that will be outputted to the `/generated` directory when we run the Bluehawk
  script.
- `/tests_package`: This directory contains the test infrastructure to actually
  run the tests by invoking the example code. (This directory can't be named
  simply `tests` as this is a protected namespace in Python.)

## To run the tests locally

### Create an Atlas cluster

To run these tests locally, you need an Atlas cluster with the sample data set
loaded and no search indexes. For best results, create a fresh M0 cluster, add
sample data, and save the connection string for use in the next step.

### Create a .env file

Create a file named '.env' at the root of the '/python' directory within this
project. Add your Atlas connection string as an environment value named
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

### Create and/or activate a Python Virtual Environment

This test suite requires you to have `Python` installed.

We strongly recommend you use `venv` to manage Python dependencies specific to
this project.

#### Create a virtual environment

In the root of the `/python` directory, if you have Python 3.3 or later
installed, you can create a virtual environment with the following command:

```
python3 -m venv ./venv
```

#### At the start of each session

When you want to work with Python examples in this project, run the
following command to activate the virtual environment:

```
/venv/bin/activate
```

Among other things, this creates a shell script called `deactivate` that you
can run when you're ready to exit the virtual environment.

#### When you're done working with the Python code

When you want to exit the virtual environment, in the same terminal where you
activated the virtual environment, run the following command:

```
deactivate
```

If you have other terminal sessions already open when you activate the virtual
environment, these other sessions may not have access to the `deactivate`
script.

You must repeat this process any time you want to work with Python examples
in this project.

### Install the dependencies

Run the following command in your virtual environment to install the required
dependencies:

```
pip install pymongo python-dotenv
```

### Run the tests

In your IDE, navigate to the `/tests_package` directory, and find the test file you want
to run. For example, `/tests_package/test_manage_indexes.py`.

Use your IDE to run all tests in the file, or run a specific test.

- File tests: run all tests in the given file
- Individual tests: press the Run button next to the name of the test you want to run.

Alternately, if you prefer to run tests from the command line:

#### Run All Tests from the command line

From the root of the `/python` directory, run:

```
python -m unittest discover tests_package
```

In this command, `tests_package` is the name of the directory that contains the tests.

#### Run Individual Tests from the command line

```
python3 -m unittest tests_package/FILENAME -k TEST_METHOD_NAME
```

For example:

```
python3 -m unittest tests_package/test_manage_indexes.py -k test_create_basic
```

## To run the tests in CI

Two GitHub workflows run these tests in CI automatically when you change any
files in the `examples` directory:

- `.github/workflows/test-python-examples-in-atlas.yml`
- `.github/workflows/test-python-examples-in-docker.yml`

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

Install Bluehawk, and then run the Bluehawk script at the root of the `/python`
directory to generate updated example files:

```
./bluehawk.sh
```
