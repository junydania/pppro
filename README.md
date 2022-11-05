# PPRO Project Instruction


### Pre-requirements
In your machine install [AWS-CLI](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html) and [CONFIGURE](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html) to save configuration of the your database
Terraform
AWS Account

### App structure Broken Down by Component
- Infrastructure code in Terraform
- Hello World application written in golang and works with DynamoDB
- CI/CD pipeline using GitHub actions


### Setup

- Deploying this application requires first setting up the infrastructure starting with ECR. To deploy ECR, run the command
```bash
    cd environment/prod
    terraform plan -target=module.ecr -out=ecr.out
    terraform apply ecr.out
```

With ECR created, lets create first version of the application and have it pushed to ECR
To execute the commands below, ensure you have make, aws-cli all installed and configured.
Ensure the user account has sufficient permissions to push a docker image to ECR
```bash
    cd hello-app
    make build-container
    make push-container
```

With the image pushed, copy the image URL from ECR and paste in the task-definition section of the terraform code
for deploying the application. With that updated. Run terraform and apply the changes

```bash
    cd environment/prod
    terraform plan -out=hello.out
    terraform apply hello.out
```

With ECS cluster and service created, update the necesary parameters in pipeline with the required values
Navigate to `.github/workflows/ppro-app.yml` file and update the following parameters
```
    AWS_REGION: MY_AWS_REGION                   # set this to your preferred AWS region, e.g. eu-central-1
    ECR_REPOSITORY: MY_ECR_REPOSITORY           # set this to your Amazon ECR repository name
    ECS_SERVICE: MY_ECS_SERVICE                 # set this to your Amazon ECS service name
    ECS_CLUSTER: MY_ECS_CLUSTER                 # set this to your Amazon ECS cluster name
    ECS_TASK_DEFINITION: MY_ECS_TASK_DEFINITION # set this to the path to your Amazon ECS task definition                                  # file, e.g. .aws/task-definition.json
    CONTAINER_NAME: MY_CONTAINER_NAME           # set this to the container name
```

Commit and push your changes and the pipeline should start managing the build and deployment pipeline.

