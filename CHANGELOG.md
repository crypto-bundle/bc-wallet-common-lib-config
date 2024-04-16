# Change Log

## [v0.0.1] - 18.02.2023 18:49 MSK
### Changed
* Lib-config moved to another repository - https://github.com/crypto-bundle/bc-wallet-common-lib-config

## [v0.0.2] - 26.03.2023 22:33 MSK
### Added
* Added support for json-based config files
* Added clear-env flow after successfully config preparation
* Added LD flag manager
* Added support of dependent service-components in envconfig variable-pool mechanism
### Fixed
* Bug with missing secret service-component

## [v0.0.3] - 09.02.2024 18:29 MSK
### Added
* Added init-flow in envconfig prepare struct flow
* Changed go-namespace
* Content of license filed changed-back to MIT

## [v0.0.4] - 16.04.2024
### Changed
* Changed ldflag manager service
  * Removed version property
  * Added data validation in manager instance creation