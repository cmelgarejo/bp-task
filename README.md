# bp-task - IPFS Token Metadata API

This is a **Golang-based RESTful API** that serves as a middleware service for fetching and storing IPFS (InterPlanetary File System) token metadata.

## Core Functionality

1. **IPFS Data Processing**: 
   - Accepts a CSV file containing IPFS Content Identifiers (CIDs) via a POST endpoint
   - Fetches JSON metadata from IPFS using the public gateway (`https://ipfs.io/ipfs/{cid}`)
   - Stores this metadata in a PostgreSQL database

2. **RESTful API Endpoints**:
   - `POST /ipfs` - Upload a CSV file with IPFS CIDs to process and store metadata
   - `GET /tokens` - Retrieve all stored token metadata from the database
   - `GET /tokens/{cid}` - Retrieve specific token metadata by IPFS CID

## Key Features

- **Asynchronous Processing**: Uses goroutines with throttling (max 5 concurrent) to process IPFS CIDs
- **Database Storage**: Uses PostgreSQL with pgx driver for efficient data storage
- **Security**: Basic Auth middleware for API protection
- **Logging**: Structured logging with Go's `slog` package
- **Graceful Shutdown**: Proper server shutdown handling
- **Timeouts**: HTTP server with read/write timeouts to prevent DoS attacks

## Technical Stack

- **Language**: Go 1.21.6
- **Database**: PostgreSQL (via pgx/v5)
- **Web Framework**: Native Go HTTP server (no external routing framework)
- **Data Format**: JSON for API responses, CSV for input
- **Authentication**: Basic Auth

## Use Case
This API serves as a caching layer for blockchain/NFT-related applications where IPFS is used to store token metadata. It fetches metadata from IPFS and stores it locally for faster access, rather than hitting IPFS every time metadata is needed.

- GET /tokens: This endpoint fetches all data stored in the database and returns it in JSON format.
- GET /tokens/{cid}: This endpoint fetches only the one record for that individual IPFS cid

## Running the API

Example values are in the .env.example file, so please copy it over to `.env` and fill in the values needed.

Run `go mod tidy` to get all the dependencies

Then run `make` it will default to `make run` which will run the API on port 8080, you can change the port by setting
the `API_PORT` environment variable, or changing that info in the .env file.

`DATABASE_URL` is also is need to be set in the env or `.env` file.

## Running the tests

using REST Client extension for VSCode, you can run the tests in the `e2e-tests.rest` file.
<https://marketplace.visualstudio.com/items?itemName=humao.rest-client>

## Notes

- Used pgx for db access, it's a fast and simple library for postgres.
- Used go's `slog` for logging, it's a simple logging library now part of the stdlib.
- Used a `makefile` for running the api, it's simple and easy to use.
- Added Basic Auth middleware so its a little bit more "secure"
- Didn't used anything else for routing could've gone with gorilla/mux or chi, stdlib but its enough for this task.
- Didn't used any ORM, I like to keep it simple and use simple queries.
- [TODO] Add more tests, actual runnable `go test` tests.
- [TODO] Containerize the api, and use docker-compose to run the tests and all dependencies
