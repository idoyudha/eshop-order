# Order Service
Part of [eshop](https://github.com/idoyudha/eshop) Microservices Architecture.

## Overview
This service handles order management of user and admin. Using postgres as main database with Command and Query Responsibility Segregation (CQRS). A pattern that separates read and update operations for a data store. Write database for transactional (command) and read database mostly for get operation (query). Read database will get updated data from event that generated by command service.

## Architecture
```
eshop-auth
├── .github/
│   └── workflows/          # github workflows to automatically test, build, and push
├── cmd/
│   └── app/                # configuration and log initialization
├── config/                 # configuration
├── internal/   
│   ├── app/                # one run function in the `app.go`
│   ├── constant/           # global constant
│   ├── controller/         # serve handler layer
│   │   ├── http/
│   │   |   └── v1/         # rest http
│   │   ├── kafka           # kafka consumers
│   │   └── redis           # redis subscriber (ex: subscribe to expired key)
│   ├── dto/                # data transfer object global (ex: kafka publisher and consumer)
│   ├── entity/             # entities of business logic (models) can be used in any layer
│   └── usecase/            # business logic
│       ├── commandrepo/    # database that business logic works with, only command database
│       └── queryrepo/      # database that business logic works with, only query database
│   
├── migrations/             # sql migration
└── pkg/
    ├── httpserver/         # http server initialization
    ├── kafka/              # kafka initialization
    ├── logger/             # logger initialization
    └── postgresql/         # postgresql initialization
```

## Tech Stack
- Programming Language: Go
- CI/CD: Github Actions
- Framework: Gin
- Database: PostgreSQL
- Identity and Access Management: AWS Cognito
- Message Broker: Apache Kafka and Redis Pub/Sub
- Container: Docker

## API Documentation
tbd