# GBA (Go base API)

The Go Base API is a secure, robust foundation for building mobile application servers and other Go-based projects. It features comprehensive user management and authentication systems, ensuring secure access and data integrity. Initially developed for the Penpals app, this API serves as a versatile starting point for any Go developer looking to accelerate their project development with pre-built, reliable backend functionality.
## Author

- [@enzo-gbd](https://github.com/enzo-gbd)


## Used By

This project is used by the following projects:

- `Penpals`: An AI-based app designed to make learning new languages easy and intuitive through conversations with artificial intelligence pen pals


## Environment Variables ($ROOT/local.env)

To configure your application correctly, you will need to set the following environment variables. These variables allow the application to connect to external services like PostgreSQL and Redis, and to configure security parameters for token management.

### PostgreSQL Variables

`POSTGRES_HOST`: Hostname of the PostgreSQL server. Default is localhost.

`POSTGRES_PORT`: Port on which the PostgreSQL server is running. Replace port with the actual port number.

`POSTGRES_USER`: Username for accessing the PostgreSQL database.

`POSTGRES_PASSWORD`: Password for the PostgreSQL user.

`POSTGRES_DATABASE`: Name of the database to connect to.

`POSTGRES_SSLMODE`: SSL mode for the PostgreSQL connection. Typical values are disable, require, and verify-full.

### Redis Variables

`REDIS_HOST`: Hostname of the Redis server. Default is localhost.

`REDIS_PORT`: Port on which the Redis server is running. Replace
port with the actual port number.

`REDIS_PASSWORD`: Password to access the Redis server, if applicable.

### Client Origin

`CLIENT_ORIGIN`: Specifies the origin of the clients that are allowed to access the server, typically used for CORS settings. Default is localhost.

### Access Token Variables

`ACCESS_TOKEN_PRIVATE_KEY`: Private key for signing JWT access tokens.

`ACCESS_TOKEN_PUBLIC_KEY`: Public key for verifying JWT access tokens. Should match the private key.

`ACCESS_TOKEN_EXPIRED_IN`: Lifespan of the access token. Default is 15m (15 minutes).

`ACCESS_TOKEN_MAXAGE`: Maximum age (in seconds) that the access token is considered valid. Default is 15.

### Refresh Token Variables

`REFRESH_TOKEN_PRIVATE_KEY`: Private key for signing JWT refresh tokens.

`REFRESH_TOKEN_PUBLIC_KEY`: Public key for verifying JWT refresh tokens. Should match the private key.

`REFRESH_TOKEN_EXPIRED_IN`: Lifespan of the refresh token. Default is 60m (60 minutes).

`REFRESH_TOKEN_MAXAGE`: Maximum age (in seconds) that the refresh token is considered valid. Default is 60.

## Environment Variables ($ROOT/docker/.env)

### PostgreSQL Variables

`POSTGRES_USER`: Username for accessing the PostgreSQL database.

`POSTGRES_PASSWORD`: Password for the PostgreSQL user.

`POSTGRES_DATABASE`: Name of the database to connect to.

### Redis Variables

`REDIS_PASSWORD`: Password to access the Redis server, if applicable.

### Environment target

`TARGET`: The environment targeted (development, production, ...)

## TLS Configuration

Create the TLS directory

```bash
  mkdir -p tls
```

Add your SSL/TLS certificate and private key

```bash
  cp path/to/your/cert.pem tls/cert.pem
  cp path/to/your/key.pem tls/key.pem
```

Verify the files are correctly placed

```bash
  ls tls
```
## Run Locally

Clone the project

```bash
  git clone https://github.com/enzo-gbd/GBA
```

Go to the project directory

```bash
  cd GBA
```

Build the docker

```bash
  cd docker
  docker-compose up -d
```


