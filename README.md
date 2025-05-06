# Payments Api

Comprehensive API offering REST services developed in GO for the administration and control of payment operations. This module allows efficient connection with various payment gateways through DEUNA, thus facilitating the processing of payments related to ClubHub clients' debts and reservations.

## Getting started

```
cd existing_repo
git remote add origin https://gitlab.com/clubhub.ai1/organization/backend/payments-api.git
git branch -M main
git push -uf origin main
```

## Running Base Dependencies

To run Temporal in Docker, follow these steps:
Navigate to the `docker/` directory and execute the following command to base dependencies:

```sh
docker-compose up
```

## Running Temporal in Docker

To run Temporal in Docker, follow these steps:
Navigate to the `docker/temporal` directory and execute the following command to start Temporal and its dependencies:

```sh
docker-compose up
```

## Create payments namespace

```sh
docker exec temporal-admin-tools tctl --namespace payments namespace register
```