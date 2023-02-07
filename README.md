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
# Database IP - can contain port
DB_IP=
# Datbase Username
DB_USERNAME=
# Database Password
DB_PASSWORD=
# Server port
PORT=
# Public URL of the server - will get suffixed with `/auth/callback` - this is where the user will be redirected after authenticating
PUBLIC_URL=
```

### Running
To run the project run `go run ./src` in the root directory of the project.
