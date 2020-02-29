# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
- Remove dependency on github.com/pkg/errors
- Update to Go 1.14

## [v0.5.0] - 2019-10-08
- Add an (optional) trusted header for inferring client IP address behind reverse proxies, load balancers, etc.

## [v0.4.1] - 2019-04-19
- Fix wrong default logger name

## [v0.4.0] - 2019-03-24
- Introduce options to control TLS and timeout configuration
- Write unit tests

## [v0.3.0] - 2019-03-18
- Update to the changed Module interface of joe v0.4.0

## [v0.2.0] - 2019-03-10
- Add HTTP headers to emitted `RequestEvent`

## [v0.1.0] - 2019-03-03

Initial alpha release

[Unreleased]: https://github.com/go-joe/http-server/compare/v0.5.0...HEAD
[v0.5.0]: https://github.com/go-joe/http-server/compare/v0.4.1...v0.5.0
[v0.4.1]: https://github.com/go-joe/http-server/compare/v0.4.0...v0.4.1
[v0.4.0]: https://github.com/go-joe/http-server/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/go-joe/http-server/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/go-joe/http-server/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/go-joe/http-server/releases/tag/v0.1.0
