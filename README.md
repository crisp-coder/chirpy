# Chirpy

###

A simple http server from a course project on boot.dev<br>
The api supports user registrations, posting and getting chirps with user authentication and authorization.<br>

## ENV

Currently the project loads a json config file ".chirpyconfig.json" from users home directory.<br>

You could optionally load the keys from a .env file by modifying main.go<br>

The json file contains:<br>

```
{
    "DB_URL": "URL_TO_POSTGRESQL_SERVER_DATABASE"
    "JWT_SECRET": JWT_SECRET_HERE
    "POLKA_KEY": API_KEY_HERE
}
```

## AUTHORIZATION and AUTHENTICATION

JSON Web Tokens - 1hour<br>
Refresh Tokens - 60 days<br>
Polka API Key - for user upgrade webhook<br>

A request is authenticated by looking up the user associated with the JWT in the database.<br>
A request is authorized if the user has permission to access that resource. This is performed via database lookup.<br>
Users can request a new JWT by logging in again or by requesting the refresh endpoint while the Refresh token is not expired or revoked.<br>

## POSTGRESQL

This api works with a postgresql database.<br>
psql - driver to communicate with the database.<br>
goose - to handle database migrations.<br>
sqlc - to generate go database code from sql query files.<br>

## API ENDPOINTS

GET /admin/metrics - returns number of api accesses<br>
POST /admin/reset - deletes all data in database<br>
GET /api/healthz - returns "OK" if api is running<br>
POST /api/users - creates a new user<br>
Headers: {"Authorization": "Bearer {JWT_TOKEN}"}<br>
Body: {"email": EMAIL, "password": PWD}<br>
PUT /api/users - updates a users email and password<br>
Headers: {"Authorization": "Bearer {JWT_TOKEN}"}<br>
Body: {"email": EMAIL, "password": PWD}<br>
POST /api/login - logs in user and creates refresh token<br>
Headers: {"Authorization": "Bearer {JWT_TOKEN}"}<br>
Body: {"email": EMAIL, "password": PWD}<br>
POST /api/refresh - gets a new jwt if refresh token has not expired<br>
POST /api/revoke - revokes a user's refresh token<br>
POST /api/chirps - creates a new chirp<br>
GET /api/chirps - can provide optional author_id and sort=asc params<br>
GET /api/chirps/{chirpID} - returns a single chirp with matching chirp id<br>
DELETE /api/chirps/{chirpID} - deletes a single chirp with matching chirp id<br>
POST /app/polka/webhooks - handles single "user.upgrade" event from polka payment processor<br>
