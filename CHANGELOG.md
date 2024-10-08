# Change Log

## [v0.0.6] - 08.10.2024
### Changed
* Added support of lib-errors for config prepare flow:
  * base config
  * JSON config
* Changed MIT License to MIT NON-AI License
* Added linters config and fixed all linter issues

## [v0.0.4, v0.0.5] - 16.04.2024
### Changed
* Changed ldflag manager service
  * Removed version property
  * Added data validation in manager instance creation
* Bump golang version 1.19 -> 1.22

## [v0.0.3] - 09.02.2024 18:29 MSK
### Added
* Added init-flow in envconfig prepare struct flow
* Changed go-namespace
* Content of license file changed-back to MIT

## [v0.0.2] - 26.03.2023 22:33 MSK
### Added
* Added support for JSON-based config files
* Added clear-env flow after successfully config preparation
* Added LD flag manager
* Added support of dependent service-components in envconfig variable-pool mechanism
### Fixed
* Bug with missing secret service-component

## [v0.0.1] - 18.02.2023 18:49 MSK
### Changed
* Lib-config moved to another repository - https://github.com/crypto-bundle/bc-wallet-common-lib-config