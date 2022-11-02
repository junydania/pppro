## Deploying the application and Infrastructure

This guide provides instruction to deploy the infrastructure and application to AWS


### Instruction

A bunch of modules were created and used in getting the environment ready.
Infrastructure code is deployed from environments/prod folder

## What does the code deploy ?

    - VPC
    - Subnets
    - ECR
    - ECS Cluster
    - ECS Task definition
    - ECS Service
    - EFS

## Application Deployment

The bitcoin application is deployed as a daemon within the ECS cluster. Built using GitHub Actions and deployed to ECR. 
A Makefile for building and deploying the docker image was also added to the bitcoin-app repository