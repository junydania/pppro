# Stack Management API

API for storing workspace details and provider details in DynamoDB

## Getting Started
### Pre-requirements
In your machine install [AWS-CLI](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html) and [CONFIGURE](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) to save configuration of the your database

### Install
If you already have golang installed you can install by running the command:
```sh
go get -u ./...
```

### Init server
To start the project run the command:
```sh
go run cmd/app/main.go
```
You can see in your terminal this log:
`service running on port  :8082`

## Usage

### Hello

### Get_All
This route return all hello message written into the database
```text
    GET - http://localhost:8082
```

#### Post
This route creates an item in the database
```text
    POST - http://localhost:8082
    {
        "content": "Hello World"
    }
```
