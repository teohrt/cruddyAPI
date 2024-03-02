# cruddyAPI

[![CircleCI](https://circleci.com/gh/teohrt/cruddyAPI/tree/master.svg?style=svg)](https://circleci.com/gh/teohrt/cruddyAPI)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

Quick note:

This project is an old example created [for someone I used to mentor](https://mentorcruise.com/mentor/TraceOhrt/). The primary objective was to showcase an optimized infrastructure design, specifically addressing concerns about the escalating costs associated with their EC2 instances. The mentee faced challenges justifying the expense due to limited monthly traffic.

By transitioning the server hosting to Lambda, a substantial cost reduction was achieved since the new platform allows them to only pay for the compute that they use. This decision was driven by the need to align infrastructure expenses with actual usage. Notably, the codebase is crafted in a way that Lambda is not a mandatory dependency, offering flexibility. If the need arises or preferences change, reverting to EC2 hosting is straightforward, ensuring adaptability for future requirements.

It's worth noting that as the traffic to the project increases, reevaluating the hosting strategy may be beneficial. At higher scales, API Gateway costs can become a consideration, and the potential latency introduced by Lambda's cold starts might impact customer experience. In such scenarios, switching back to EC2 could be a strategic choice for optimizing both costs and performance.

---

![cruddyAPI](docs/cruddyAPI.jpeg)

This RESTful CRUD API will manage your profile data on AWS at any scale. Just hit deploy and your new microservice will be up and running!

[API Gateway](https://aws.amazon.com/api-gateway/) routes requests, [Lambda](https://aws.amazon.com/lambda/) handles the corresponding logic and your profile data gets stored in [DynamoDB](https://aws.amazon.com/dynamodb/). Should you like to take a deeper look, logs are stored in [CloudWatch](https://aws.amazon.com/cloudwatch/) and [X-Ray](https://aws.amazon.com/xray/) handles the distributed tracing, both of which conveniently make this microservice's internals available to you in dashboards.

## Routes

- [Create profile](docs/createProfileContract.md)
  - `POST api/v1/profiles`
- [Read / retrieve profile](docs/getProfileContract.md)
  - `GET api/v1/profiles/{id}`
- [Update profile](docs/updateProfileContract.md)
  - `PUT api/v1/profiles/{id}`
- [Delete profile](docs/deleteProfileContract.md)
  - `DELETE api/v1/profiles/{id}`
- [Health](docs/healthContract.md)
  - `GET api/v1/health`

## Requirements

### Deployment Requirements

- [Golang](https://golang.org/dl/) >= 1.11
- [Terraform](https://www.terraform.io/downloads.html) >= 0.12.8
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/install-macos.html#awscli-install-osx-pip) >= 1.16.230
  - Infrastructure costs fit comfortably within AWS' Free-Tier with typical development usage

### Local Development Requirements

- All of the above
- [Docker](https://docs.docker.com/v17.12/install/)
- [DynamoDB-local Docker Image](https://hub.docker.com/r/amazon/dynamodb-local/)

## Usage

### Deployment

```bash
git clone https://github.com/teohrt/cruddyAPI.git
cd cruddyAPI
go mod download
make deploy
```

### Local Development

This requires two terminals. One to run the local db, and one to run the API server.
To Set up the local db, run the command

```bash
make db-start
```

Now, initialize the db table in another terminal and start the API server with these commands:

- This app listens to API Gateway requests by default. For local development you will need to make it accessible to normal http requests. Do this by commenting line 37 and uncomment lines 17 and 35 in file "app/app.go".

```bash
make db-table-init
make run-locally
```
### Delete Infrastructure
```bash
make destroy
```
## Licence

See LICENSE.
