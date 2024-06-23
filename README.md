# Goverland Datasource Snapshot

<a href="https://github.com/goverland-labs/goverland-datasource-snapshot?tab=License-1-ov-file" rel="nofollow"><img src="https://img.shields.io/github/license/goverland-labs/goverland-datasource-snapshot" alt="GPL 3.0" style="max-width:100%;"></a>
![unit-tests](https://github.com/goverland-labs/goverland-datasource-snapshot/workflows/unit-tests/badge.svg)
![golangci-lint](https://github.com/goverland-labs/goverland-datasource-snapshot/workflows/golangci-lint/badge.svg)

Goverland Datasource Snapshot is a micro service designed to manage and synchronize various data sources for Goverland ecosystem.
At the moment the focus is on Snapshot API.
The microservice also leverages technologies such as gRPC and NATS.
PostgreSQL is used as database.
The application exposes Prometheus metrics.

## Usage

- The application initializes various services and workers to manage proposals, spaces, votes, and messages.
- It sets up gRPC servers for remote procedure calls and integrates with NATS for event publishing.

## Features

- Snapshot SDK Integration: facilitates interaction with Snapshot services for proposals, spaces, votes, and messages.
- Event Publishing: utilizes NATS for event-driven architecture.
- gRPC Services: provides gRPC-based APIs for interacting with the Snapshot from internal services.
- Database: Integrates with PostgreSQL for storing and managing data.
- Metrics and Monitoring: exposes Prometheus metrics for monitoring.
- Health Checks: Implements health check endpoints to monitor the application's status.

## For Developers

### Pre-requisites

- GoLang 1.21
- Docker

### Run Locally

1. Clone the repository:
   ```bash
   git clone https://github.com/goverland-labs/goverland-datasource-snapshot.git
   cd goverland-datasource-snapshot
   ```
2. Install the required Go dependencies:
   ```bash
    go mod tidy
   ```
3. Spin up the necessary containers:
   ```bash
   docker-compose up -d
   ```
4. Run the application. Don't forget to provide your own Snapshot API key:
   ```bash
    POSTGRES_DSN="host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable" \
    SNAPSHOT_API_KEY=YOUR_KEY \
    go run main.go
   ```

#### To tear down docker compose environment

```bash
docker-compose down
```

### Prometheus

Prometheus metrics are exposed at [http://localhost:2112/metrics](http://localhost:2112/metrics).

## Docker Image

1. Build the Docker Image:

   - Navigate to the root directory containing a `Dockerfile`.
   - Run the following command to build the Docker image:
     ```bash
     docker build -t goverland-datasource-snapshot .
     ```

2. Run the Docker Container:

   - After the image is built, you can run it as a container. Use the following command to start the container:
     ```bash
     docker run -d --name goverland-datasource-snapshot -p 8080:8080 goverland-datasource-snapshot
     ```

## Contribution Rules

To request or propose new features or bug fixes, [create a new issue](https://github.com/goverland-labs/goverland-datasource-snapshot/issues/new/choose) using specific template.
If you wish to contribute to the project, please read [CONTRIBUTING](CONTRIBUTING.md) guide carefully
or contact us via [Discord](https://discord.gg/uerWdwtGkQ).

## Changelog

Check all updates in our [CHANGELOG](CHANGELOG.md).

## License

Goverland Datasource Snapshot is [GPL-3.0 licensed](./LICENSE).
