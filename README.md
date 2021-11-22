# image-converter
[![Lint](https://github.com/Konstantsiy/image-converter/actions/workflows/lint.yml/badge.svg)](https://github.com/Konstantsiy/image-converter/actions/workflows/lint.yml)
[![CircleCI](https://circleci.com/gh/circleci/circleci-docs.svg?style=shield)](https://circleci.com/gh/circleci/circleci-docs)

Service that expose a RESTful API to convert JPEG to PNG and vice versa and compress the image 
with the compression ratio specified by the user. The user has the ability to view
the history and status of their requests (for example, queued, processed, completed) and upload 
the original image and the processed one.
# Endpoints
- /user/login - user authorization [POST]
- /user/signup - user registration [POST]
- /conversion - convert needed image [POST]
- /images/{id} - get needed image [GET]
- /requests - get the user's requests history [GET]
# Architecture Diagram
![alt text](./docs/architecture-diagram-2.jpg)
# Database Scheme
![alt text](./docs/db.png)

# Setting up the environment

Before launching the application, you need to configure the project environment on CircleCI.

Configuring the database (AWS RDS - Postgres):
```text
DB_USER
DB_PASSWORD
DB_NAME
DB_HOST
DB_PORT
DB_SSL_MODE (optional)
```
Configuring JWT:
```text
JWT_SIGNING_KEY
```
Configuring a Message Broker (Amazon MQ - RabbitMQ):
```text
RABBITMQ_QUEUE_NAME
RABBITMQ_AMQP_CONNECTION_URL
```
Configuring AWS account and AWS S3:
```text
AWS_REGION
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
AWS_BUCKET_NAME
AWS_CLUSTER_NAME
AWS_PRIVATE_ECR_ACCOUNT_URL_API
AWS_PRIVATE_ECR_ACCOUNT_URL_WORKER
AWS_TASK_DEFINITION_FAMILY
```

