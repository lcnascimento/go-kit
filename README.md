[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/)
[![Release Workflow](https://github.com/lcnascimento/go-kit/actions/workflows/release.yaml/badge.svg?branch=main)](https://github.com/lcnascimento/go-kit/actions/workflows/release.yaml)

# Go Kit

Go Kit is a repository of utilitary packages written in Go, designed to improve development experience.


## Internal structure

This repository is designed as a Monorepo. Each folder represents an individual and independent package.

Every package is automatically versioned using [Semantic Release](https://github.com/semantic-release/semantic-release) default rules. Every semantic commit pushed to `main` branch, that updates a package codebase, will generate a new release version for its package. If a commit does not touches the code base of a package, its version must not be increased. A single commit can increase version of multiple packages when it touches them all.

## How to use a package

`go get github.com/lcnascimento/go-kit/<package-name>`

## How to predict next package versions locally

Just run `npm install` and `lerna exec --concurrency 1 -- semantic-release --tag-format='${LERNA_PACKAGE_NAME}/v\${version}'`

## Package Best Practices

- High test coverage;
- Include a well designed `README.md` file;
- Include an example in `examples` folder, that illustrates how to use your package;

## Frequently Asked Questions

### What is the format of a semantic commit?

`<type>(<scope>): <short summary>`

Type must be one of the following:

- build: Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)
- ci: Changes to our CI configuration files and scripts (examples: CircleCi, SauceLabs)
- docs: Documentation only changes
- feat: A new feature
- fix: A bug fix
- perf: A code change that improves performance
- refactor: A code change that neither fixes a bug nor adds a feature
- test: Adding missing tests or correcting existing tests

Commits with a footer message containing the phrase `BREAKING CHANGE: ` indicates a Major release.

You can read more about this [here](https://github.com/angular/angular/blob/master/CONTRIBUTING.md#-commit-message-format).

Notice that, if your changes increases the Major version of your package, you MUST update package's go.mod properly. For instance, if there is a BREAKING CHANGE to package `foo`, that increases its version to `2.0.0`, you MUST update `foo/go.mod` with:

`module github.com/lcnascimento/go-kit/foo` -> `module github.com/lcnascimento/go-kit/foo/v2`

This need is required by Go design.

### What happens if my Pull Request does not have any semantic commit?

No version will be created.

### How can I create release candidates?

Just merge your changes into `beta` branch.
