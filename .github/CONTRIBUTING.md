# Contributing

**Ax** is open source software; contributions from the community are
encouraged and appreciated. Please take a moment to read these guidelines
before submitting changes.

> Please be aware that this project is maintained by a [single individual](https://github.com/jmalloc),
likely for use in commercial products. As such, decisions about the design and
functionality of the software may be wholly or in part governed by the
requirements of those products.

## Requirements

- [Go 1.10](https://golang.org/)
- [GNU make](https://www.gnu.org/software/make/) (or equivalent)

## Running the tests

Then run:

    make

The default target of the make file installs all necessary dependencies and runs
the tests.

Code coverage reports can be built with:

    make coverage

To rebuild coverage reports and open them in a browser, use:

    make coverage-open

## Submitting changes

Change requests are reviewed and accepted via pull-requests on GitHub. If you're
unfamiliar with this process, please read the relevant GitHub documentation
regarding [forking a repository](https://help.github.com/articles/fork-a-repo)
and [using pull-requests](https://help.github.com/articles/using-pull-requests).

**Before submitting a pull-request that includes new features or API changes
please open an issue to discuss the changes.** Include the rationale for the
changes and information about the potential backwards compatibility concerns.

Before submitting any pull-request please run:

    make prepare

This will apply any automated code-style updates, run linting checks, run the
tests and build coverage reports. Please ensure that your changes are tested and
that a high level of code coverage is maintained.

## Branching and versioning

This project uses [semantic versioning](https://semver.org). Releases are made
by creating an annotated tags of commits on the `master` branch. For this
reason, `master` should always be a functional version of the software.
