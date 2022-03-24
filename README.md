# fukaeri

Open-source 4chan archiver written in Go using MongoDB as data storage

- respects 4chan's API usage policies
- mimics real-life website usage patterns

## Usage
### Getting started
Recommended method of running *fukaeri* is to use a Docker container. Image is available on [Docker Hub](https://hub.docker.com/r/k0mmsussert0d/fukaeri).

Fukaeri requires MongoDB deployment to archive data to. `docker-compose.yml` provides an example of running fukaeri and MongoDB within separate containers. Default configuration file `conf.yml` can be used in this setup, making first run as simple as running `docker-compose up`.

Furthermore, [mongo-express](https://github.com/mongo-express/mongo-express) instance is configured and can be accessed at `http://localhost:8081`.

### Configuration
Configuration entries are stored in `conf.yml` file, which should always be stored alongside fukaeri executable. The purpose of most options is self-explanatory, while others are explained with comments.

## Development
This project is compatible with [Visual Studio Code Dev Containers](https://code.visualstudio.com/docs/remote/containers). Opening it inside a Dev Container brings up full development environment, including:
- Go runtime container with numerous tools provided by [microsoft/vscode-dev-containers](https://github.com/Microsoft/vscode-dev-containers)
- MongoDB container with **no data volumes**
- mongo-express instance connected to MongoDB container

When running in a Dev Container, fukaeri uses `conf.yml` file as a configuration source as well. It's suitable for running fukaeri both in dev and production-ready container, as described in [Usage](#usage).
