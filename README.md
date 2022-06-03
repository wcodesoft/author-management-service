# Author Management Service

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