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
`service running on port  :8080`

## Usage

### Health
This route return life of the project
```text
    GET - http://localhost:8080/health
```

### Get_All
This route return all workspace data in your database
```text
    GET - http://localhost:8080/workspace
```

#### Get_One
This route return specific data in your database
```text
    GET - http://localhost:8080/workspace/{ID}
```

#### Post
This route creates an item in the database
```text
    POST - http://localhost:8080/workspace
    {
        "name": "zain-workspace705",
        "account_name": "coda-payment",
        "stack_name": "codapay",
        "region": "us-west-1",
        "uptime_hours": 0,
        "workspace_details": {
            "instance_id": "i-223232323888995552",
            "instance_ip": "10.12.33.19",
            "security_groups": ["sg-984884848", "sg-232323232"],
            "subnet_ids": ["subnet-049494", "subnet-9044543"],
            "username": "odania",
            "vpc_id": "vpc-sds782737237"
        }
    }
```

#### Put
This route updates a workspace in the database
```text
    PUT - http://localhost:8080/workspace/{ID}
    
    {
        "name": "workspace1"
    }
```

#### Delete
This route removes a workspace from the database
```text
    DELETE - http://localhost:8080/workspace/{ID}
```