# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Added proxy api for delegates

## [0.5.4] - 2024-06-28

### Changed
- Updated readme.md

## [0.5.3] - 2024-06-25

### Added
- Added space update

## [0.5.2] - 2024-06-18

### Added
- Added cli command building in the docker container

### Fixed
- Updated go version up to 1.22

## [0.5.1] - 2024-06-13

### Added
- Integration with Ipfs Fetcher
- Consuming delete proposal events based on snapshot messages

### Fixed
- GetProposalIDsForUpdate with list of spaces

## [0.5.0-draft] - 2024-03-12

### Added
- Added worker for fetching metadata messages from spanshot

### Changed
- Update votes for spaces with real updates only

## [0.4.4] - 2024-03-12

### Added
- Get vote by id endpoint

## [0.4.3] - 2024-03-06

### Fixed
- Fixed Dockerfile

## [0.4.2] - 2024-03-06

### Added
- Collecting snapshot sdk key remaining and limit values

## [0.4.1] - 2024-03-01

### Changed
- Moved gRPC protocol to the separated go module

## [0.4.0] - 2024-03-01

### Changed
- Changed the path name of the go module
- Updated dependencies for self developed libraries
- Added badges with link to the license and passed workflows

### Added
- Added LICENSE information
- Added info for contributing
- Added github issues templates
- Added linter and unit-tests workflows for github actions

### Fixed
- Fixed linter warnings

## [0.3.8] - 2024-02-28

### Added
- Metrics

## [0.3.7] - 2024-02-13

### Added
- Added different api key for voting requests

## [0.3.6] - 2024-02-07

### Added
- Verified field mapping

## [0.3.5] - 2024-02-06

### Fixed
- Fixed proposal creation event

## [0.3.4] - 2023-12-29

### Fixed
- Fixed shutter voting for quadratic and ranked voting

## [0.3.3] - 2023-12-20

### Changed
- Changed app name from snapshot to goverland

## [0.3.2] - 2023-12-20

### Added
- Flagged field mapping

### Added
- Publish deleted proposals

## [0.3.1] - 2023-12-04

### Fixed
- Fixed docker build

## [0.3.0] - 2023-12-04

### Added
- Added voting implementation

## [0.2.10] - 2023-11-16

### Fixed
- Fixed selecting active proposals for voting. It was stuck based on updated_at field

## [0.2.9] - 2023-11-15

### Changed
- Decrease vote gap from 5 to 1 minutes
- Actualize DB schema

## [0.2.8] - 2023-10-02

### Changed
- SpeedUp votes processing

## [0.2.7] - 2023-09-29

### Added
- CLI commands for importing spaces, proposals and votes from CSV dumps

## [0.2.6] - 2023-09-12

### Added
- Collecting proposal votes

## [0.2.5] - 2023-08-25

### Changed
- Updated SDK version

### Added
- Added API key for snapshot.org

## [0.2.4] - 2023-07-26

### Fixed
- Fixed active proposals worker - delete not fetched proposals

## [0.2.3] - 2023-07-26

### Fixed
- Fixed illegal chars in the snapshot as string format

## [0.2.2] - 2023-07-26

### Fixed
- Fixed illegal chars in the snapshot

## [0.2.1] - 2023-07-16

### Fixed
- Fixed active proposals update worker config
- Fixed saving/updating proposal in the repository

## [0.2.0] - 2023-07-16

### Added
- Added worker for updating active proposals

## [0.1.4] - 2023-07-15

### Fixed
- Updated platform-events dependency to v0.0.17

## [0.1.3] - 2023-07-12

### Fixed
- Updated platform-events dependency to v0.0.12

## [0.1.2] - 2023-07-11

### Fixed
- Updated platform-events dependency to v0.0.11

## [0.1.1] - 2023-07-06

### Fixed
- Fixed dockerfile

### Added
- Added auto migrations

## [0.1.0] - 2023-07-06

### Added
- Added worker for fetching new proposals from snapshot.org
- Added worker for fetching unknown spaces from fetched proposals
- Added init command for fetching all ranked space identifiers
- Added init command for filter unique space identifiers
- Added init command for fetching spaces by identifiers
