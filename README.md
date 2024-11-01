# Docs Code Examples

This repository represents the eventual home for all code examples across
[MongoDB documentation](https://www.mongodb.com/docs/) properties.

Details TBD as infrastructure is added and projects related to standardization
and coverage are completed.

## Structure

Each programming language has its own directory, and each directory is split
into examples and tests. Examples live in stand-alone files in the
`language/examples` directory, grouped by topic. Tests live in consolidated
files in the `language/tests` directory, with each test executing example code
to verify that:

- Examples compile and run
- Examples do what we say they do

We use tooling to extract the tested code examples to the `generated-examples`
directory. The content of the `generated-examples` directory is tested,
excerpted code that is ready to use in the documentation.

## Example extraction

Because these examples are designed to be runnable and testable, they contain
additional code related to our infrastructure that is not relevant to the
developer consuming our documentation. For example, the tests use
language-idiomatic methods for reading environment values to get the appropriate
Atlas connection string, or have method signatures or function calls with
return values for testing purposes that are not needed in the docs example.

Docs uses a markup tool called [Bluehawk](https://github.com/mongodb-university/Bluehawk)
to add markup to extract only the relevant parts of the code examples. Each
language directory includes a script that uses this markup tool, `bluehawk.sh`,
to extract relevant code to the `generated-examples` directory. For more
details, refer to the language-specific README.md on generating code examples.

For docs purposes, writers should _only_ use code examples from the
`generated-examples` directory. Using an example directly from the
`language/examples` directory may contain code that is not desired in the docs.

## Automated testing

When any file in any of the `language/examples` directories are updated in a PR,
the GitHub workflows contained in the `.github/workflows` directory automatically
run the test suite. These workflows report the results of the test suite runs
as passes or failures on the PR.

If changing an example causes its test to fail, this should be considered
blocking for merging the example. The PR should not be merged until the example
and/or test is fixed and the test passes.

If changing an example causes an _unrelated_ test to fail, create a Jira ticket
to fix the unrelated test, but this should not block merging an example update.

## Local testing

Each language directory contains a README.md with instructions about how to
run the tests locally. Contributors can and should run tests for any examples
they're changing locally during the update process. Reviewers don't _need_ to
run the tests directly, since the tests run automatically in CI, but if they
_want_ to run the tests, they can check out the PR and run the tests locally.

## Contributing

As we codify workflows and process around this repository, a future update to
this README and the language-specific READMEs will cover contributing. For now,
this repository should be considered a proof-of-concept for discussion and
project planning.
