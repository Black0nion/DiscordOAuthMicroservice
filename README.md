## DiscordOAuthMicroservice
A lightweight microservice for authenticating users with Discord OAuth2.
Authenticates users with Discord OAuth2 and stores their tokens in a database.
Users then authenticate with a session id that is linked with those tokens.

### Configuration
To configure the project create a `.env` file in the root directory of the project with the following variables:
```dotenv
# Discord application client id
CLIENT_ID=
# Discord application client secret
CLIENT_SECRET=
# Database IP - can contain port e.g. localhost:5432
DB_IP=
# Database name
DB_NAME=
# Datbase Username
DB_USERNAME=
# Database Password
DB_PASSWORD=
# Server port
PORT=
# Redirect URL for Discord OAuth2 - even if you only use the exchange endpoint this still has to be set!
REDIRECT_URL=
# OPTIONAL: CORS origin - can be a list of origins separated by a comma - if not set, all origins are allowed (*)
CORS_ORIGIN=
```

### Running
To run the project run `go run ./src` in the root directory of the project.

### Docker
There's a docker image available on [GHCR](https://github.com/Black0nion/DiscordOAuthMicroservice/pkgs/container/discordoauthmicroservice). 