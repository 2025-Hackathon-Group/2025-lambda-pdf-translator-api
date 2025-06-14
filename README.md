# 2025-lambda-pdf-translator
Our entry for the 2025 AWS Lambda hackathon.

#### Setup

For development, use the following command to spin up and connect to the api docker.

```
docker compose -f docker-compose-dev.yml up -d && docker compose -f docker-compose-dev.yml exec api bash
```

then once in the container run use the makefile:
 - `make run` starts the server
 - `make seed` seeds the database with the following credentials:


**User**
```
Name:           "Default User"
Email:          "saml@everbit.dev"
Password:       "password"
```
**Organisation**
```
Name:  "Default Organization",
Email: "admin@default.org",
```

To run in prod mode, spin up the api with `docker compose up -d`.