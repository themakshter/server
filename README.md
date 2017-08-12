# Impactasaurus Server

Impactasaurus is changing the way charities measure and report on social impact. We are building a free, open source, easy to use, configurable impact measure tool, which is compatible with any CRM. Read more about Impactasaurus at https://impactasaurus.org.

This project is the backend for Impactasaurus. It is composed of a single golang application which uses a mongo database. A graphql API is exposed for the web app to consume.

## Getting Started

Currently impactasaurus is an invite only application. To get access please email admin@impactasaurus.org or visit our [gitter chat room](https://gitter.im/impactasaurus) and ask for an invite.

A docker compose file is available which will start the server locally with a linked mongodb. Ensure you have docker and docker compose installed, then run the following command:
```
docker-compose build && docker-compose up
```
This will start the golang application, a mongodb database and an [in browser IDE](https://github.com/graphql/graphiql) for interacting with the graphql API.

The following URLs are of interest:

 - http://localhost:8082 : The graphql IDE
 - http://localhost:8081/v1/graphql : The graphql API
 - mongodb://localhost:27017 : The mongodb database

To use the graphql IDE, you must first obtain a JWT. This can be achieved by logging into the web app and running the following javascript in the developer console:
```
localStorage.getItem("token")
```
An Authorization header should be added to requests to the graphql API, this should be of the form:
```
Authorization: Bearer {jwt}
```

You can configure the web app to communicate to your locally hosted server instance. This is detailed more in the [app project's readme](https://github.com/impactasaurus/app).

## API Documentation

GraphQL APIs include documentation, to view this, please navigate to the graphql IDE listed above. The API documentation will be visible on the right hand side of the web site.

## Configuration

The golang application is configured using environmental variables. The details of the available env vars can be found at `cmd/config.go`. Environmental variables can be added or adjusted, when using docker-compose, by editing `server.environment` within the `docker-compose.yml` file.

## Contributing

Please read the [contribution guidelines](https://github.com/impactasaurus/server/blob/master/CONTRIBUTING.md) to find out how to contribute.