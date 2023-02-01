# errors <br> [![go.mod version][go-img]][go-url] [![CI][ci-img]][ci-url] [![Codecov][codecov-img]][codecov-url] [![Maintainability][codeclimate-img]][codeclimate-url] [![Go Report Card][goreportcard-img]][goreportcard-url] [![License][license-img]][license-url] [![Go Reference][godoc-img]][godoc-url]

Go library to construct errors with fields for structured logging.

Features:
- Wrap an error with string prefix
- Add custom fields to an error
- Extract all fields from chain of wrapped errors
- Combine several errors into one error (build errors tree)
- Extract paths to each leaf from the errors tree
- Logger agnostic

## Motivation

When structured logger is used, it's better to have constant error messages. For example, message should not contain ID of your entity. Instead, such additional data should be logged in a separate fields. That makes it easier to search, group and analyse logs.

**Bad:**

```json
{"level": "error", "message": "can't find order a881ff5c-ef23-4e6c-a505-9b66ee42b779"}
```

**Good:**

```json
{"level": "error", "message": "can't find order", "order_id": "a881ff5c-ef23-4e6c-a505-9b66ee42b779"}
```

## Installation

```shell
go get github.com/maratori/errors
```

## Usage

TBD

## License

[MIT License][license-url]


[go-img]: https://img.shields.io/github/go-mod/go-version/maratori/errors
[go-url]: /go.mod
[ci-img]: https://github.com/maratori/errors/actions/workflows/ci.yml/badge.svg
[ci-url]: https://github.com/maratori/errors/actions/workflows/ci.yml
[codecov-img]: https://codecov.io/gh/maratori/errors/branch/main/graph/badge.svg?token=LPthtc4wLI
[codecov-url]: https://codecov.io/gh/maratori/errors
[codeclimate-img]: https://api.codeclimate.com/v1/badges/d81e78169c09c31ebe42/maintainability
[codeclimate-url]: https://codeclimate.com/github/maratori/errors/maintainability
[goreportcard-img]: https://goreportcard.com/badge/github.com/maratori/errors
[goreportcard-url]: https://goreportcard.com/report/github.com/maratori/errors
[license-img]: https://img.shields.io/github/license/maratori/errors.svg
[license-url]: /LICENSE
[godoc-img]: https://pkg.go.dev/badge/github.com/maratori/errors.svg
[godoc-url]: https://pkg.go.dev/github.com/maratori/errors
