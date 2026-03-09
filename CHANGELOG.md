# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Common Changelog](https://common-changelog.org/),
and adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 0.2.0 - 2026-03-09

### Changed

- handlers: ensure length and valid_for meta parameters are positive
- handlers: refactor cache to use generics

### Added

- ssh|tls: add support for alternate crypto implementations
- handlers: add support for JWT-based secrets
- handler: add support for API key secrets
- handlers: add support for OTP-based secrets
- jwt: include jti claim in generated tokens
- test: implement handler unit tests

### Fixed

- jwt: fix duplicate use of subject meta param

## 0.1.0 - 2026-02-16

### Added

- Initial release
