# Author Management Service

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/1f7bd19322da40fbb15afd12f154ce14)](https://www.codacy.com/gh/wcodesoft/author-management-service/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=wcodesoft/author-management-service&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/1f7bd19322da40fbb15afd12f154ce14)](https://www.codacy.com/gh/wcodesoft/author-management-service/dashboard?utm_source=github.com&utm_medium=referral&utm_content=wcodesoft/author-management-service&utm_campaign=Badge_Coverage)

Author management reusable microservice used on all Mosha projects. This component
can be split in two modules:

* **proto**: module that will hold all proto definitions.
* **service**: actual implementation of the service.

## Proto

Proto was generated using the package: https://github.com/wcodesoft/proto-builder

To add the `Go` library to your service run the command:

```bash
go get -u github.com/wcodesoft/author-management-service/grpc/go/author-management.proto
```

## Run Service

On the `service` folder execute the following command to run the service:

```bash
go run service
```

## Test service

Tests are already implemented for the service, to run the tests and get the coverage report run the following command
on the `service` folder:

```bash
go test ./... -v coverageprofile="coverage.out" 
```
