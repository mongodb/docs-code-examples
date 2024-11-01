# Java Vector Search Code Example Test PoC

This is a PoC to explore testing Java Vector Search code examples for
MongoDB documentation.

The structure of this Java project is as follows:

- `/src/main/java`: This directory contains example code, marked up with
  Bluehawk, that will be outputted to the `/generated` directory when we run
  the Bluehawk script.
- `/src/test/java`: This directory contains the test infrastructure to actually
  run the tests by invoking the example code.

## To run the tests locally

### Create an Atlas cluster

To run these tests locally, you need an Atlas cluster with the sample data set
loaded and no search indexes. For best results, create a fresh M0 cluster, add
sample data, and save the connection string for use in the next step.

### Create environment variables

Create environment variables to use when running the tests. Depending on
how you run the tests, you may do this in a few ways:

- Set them in your system through the terminal
- Add them in the IDE

The test suite requires two environment variables:

- `ATLAS_CONNECTION_STRING` - your Atlas or local connection string
- `ENV` - whether you're running the tests against Atlas or a local MDB in
  Docker. Some of the tests require us to distinguish between environents
  because, for example, Vector Search produces different ANN scores in different
  environments.
  - Valid values are:
    - `"Atlas"`: Use when running the test suite against an Atlas deployment
    - `"local"`: Use when running the test suite against a local MDB in Docker

#### Add them to your process through the terminal

You can manually set these environment values through the terminal using the
`export` command:

```console
export ATLAS_CONNECTION_STRING="<your-connection-string>"
export ENV="Atlas"
```

This is how we set the environment variables in CI. This sets the environment
variables system-wide.

#### Add them in the IDE

You can set the environment variables in a run configuration in your IDE.
In JetBrains IDEA, for example, you can do this through:

- Right click on the `Play` button to the left of a test class to run all of
  the tests in the class, or a single test.
- Select `Modify run configuration`
- Add them in the `Environment variables` line
  - You can supply them inline in a semicolon-separated-format:
    ```
    ATLAS_CONNECTION_STRING=<your-connection-string>;ENV=Atlas
    ```
  - You can press what looks like a document icon to the right side of the field,
    which brings up a modal where you can use plus and minus buttons to add,
    edit, and remove environment variables

NOTE: If you add environment variables through the JetBrains IDE:
- Omit any `"` around strings
- You must add environment variables separately for every run configuration/play
  button. i.e. if you add environment variables for the test class, but then
  later only want to run a single test, you must add them again for the single
  test.

### Install the dependencies

This test suite was written with Java JDK Zulu v17 installed. If you
do not yet have Java installed, install it. There are many ways to install Java,
which is out of the scope of this README. The simplest method to manage Java
installs is using the IntelliJ IDEA IDE: 
[JetBrains: Download a JDK](https://www.jetbrains.com/guide/java/tips/download-jdk/)

Separately, you can install `Maven` if you want to run the test suite from the
command line. For details about installing Maven, refer to
[Installing Apache Maven](https://maven.apache.org/install.html). You don't
need to install Maven if you plan to run the tests from within the IDE.

#### Install Project Dependencies in the IDE

Open the `pom.xml` file in the root of the `/java` directory in your IDE. The
first time you load it, or when you change any dependencies, a `Sync` button
appears. Press this button to download dependencies.

#### Install Project Dependencies from the Command Line

If you install `Maven`, you can install project dependencies from the command
line. From the root of the `/java` directory, run the following command to
install dependencies:

```
mvn install
```

### Run the tests

#### Run Tests from the IDE

You can press the play button next to a TestClass to run all of the tests in
the class, or next to a test name to run a single test in the IDE. Because this
test suite relies on an Atlas connection string and environment value passed
in from the environment, running tests in the IDE will fail unless you
configure the IDE with the appropriate environment variables.

See the instructions above to add the environment variables in the IDE.

#### Run All Tests from the command line

If you have `Maven` installed, you can run tests from the command line. From
the `/java` directory, run:

```
mvn test
```

#### Run a Test Class from the command line

You can run all the tests in a given test class (file).

From the `/java` directory, run:

```
mvn -Dtest=YourTestClassName test
```

For example:

```
mvn -Dtest=QueryTests test
```

#### Run Individual Tests from the command line

You can run a single test within a given test suite (file).

From the `/java` directory, run:

```
mvn -Dtest=YourTestClassName#YourTestMethodName test
```

For example:

```
mvn -Dtest=QueryTests#TestAnnQueryBasic test
```

## To run the tests in CI

Two GitHub workflows run these tests in CI automatically when you change any
example:

- `.github/workflows/test-java-examples-in-atlas.yml`
- `.github/workflows/test-java-examples-in-docker.yml`

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

Install Bluehawk, and then run the Bluehawk script at the root of the `/java`
directory to generate updated example files:

```
./bluehawk.sh
```
