# Task

Using the data from your database, set up a Golang-based RESTful API with the following endpoints:

- GET /tokens: This endpoint should fetch all data stored in the database and return it in JSON format.
- GET /tokens/{cid}: This endpoint should fetch only the one record for that individual IPFS cid

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
