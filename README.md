# Chirpy

###

A simple http server from a course project on boot.dev
The api supports user registrations, posting and getting chirps with user authentication and authorization.

## ENV

Currently the project loads a json config file ".chirpyconfig.json" from users home directory.

You could optionally load the keys from a .env file by modifying main.go

The json file contains:

```
{
    "DB_URL": "URL_TO_POSTGRESQL_SERVER_DATABASE"
    "JWT_SECRET": JWT_SECRET_HERE
    "POLKA_KEY": API_KEY_HERE
}
```

## AUTHORIZATION and AUTHENTICATION

JSON Web Tokens - 1hour
Refresh Tokens - 60 days
Polka API Key - for user upgrade webhook

A request is authenticated by looking up the user associated with the JWT in the database.
A request is authorized if the user has permission to access that resource. This is performed via database lookup.
Users can request a new JWT by logging in again or by requesting the refresh endpoint while the Refresh token is not expired or revoked.

## POSTGRESQL

This api works with a postgresql database.
psql - driver to communicate with the database.
goose - to handle database migrations.
sqlc - to generate go database code from sql query files.

## API ENDPOINTS

GET /admin/metrics - returns number of api accesses
POST /admin/reset - deletes all data in database
GET /api/healthz - returns "OK" if api is running
POST /api/users - creates a new user
PUT /api/users - updates a users email and password
POST /api/login - logs in user and creates refresh token
POST /api/refresh - gets a new jwt if refresh token has not expired
POST /api/revoke - revokes a user's refresh token
POST /api/chirps - creates a new chirp
GET /api/chirps - can provide optional author_id and sort=asc params
GET /api/chirps/{chirpID} - returns a single chirp with matching chirp id
DELETE /api/chirps/{chirpID} - deletes a single chirp with matching chirp id
POST /app/polka/webhooks - handles single "user.upgrade" event from polka payment processor
