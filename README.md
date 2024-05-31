# Go-based JWT Auth Server

This project provides a simple and efficient skeleton for a JWT (JSON Web Token) authenticate & authorization server implemented in Go. It is designed to be lightweight, making it an excellent tool for learning how authentication & authorization servers work.

## Features

- Zero Dependencies: No external libraries are required.
- RS256 Algorithm: Utilizes the RS256 (RSA Signature with SHA-256) algorithm for signing JWTs.

## Prerequisites

- Go (version 1.18 or higher)
- RSA Private Key for signing (have a sample `private-key.pem`)

## Installation & Run

```sh
# Clone the Repository:
git clone https://github.com/yourusername/jwt-encryption-server.git

cd jwt-encryption-server

# Generate private-key.pem
openssl genrsa -f4 -out private-key.pem 2048

# Run the Server:
go run cmd/serverd/main.go
```

## APIs

### To get JWT given claims

```bash
curl --location 'http://localhost:8080/token' \
--header 'content-type: application/json' \
--data '{
    "client_id": "sample-client-id",
    "client_secret": "sample-client-secret"
}'
```

Response:

```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJodHRwOi8vbG9jYWxob3N0Ojk5OTkvIiwiYXpwIjoic2FtcGxlX3N1YmplY3QiLCJleHAiOjE3MTcyNjMyNTYsImd0eSI6ImNsaWVudC1jcmVkZW50aWFscyIsImlhdCI6MTcxNzE3Njg1NiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwLyIsInN1YiI6InNhbXBsZV9zdWJqZWN0QGNsaWVudHMifQ.l59hG6zj9nwRPQrwN6WFQL7QyZH-MykRIXMPdVCyONcuVOyOEgNBhRmq49-c-uMmoiBPMYP97UE043VjIEAC5-nzOpODf9tl-jVlcSk2mVpyaIFRffUFsMDXS-qdOVmx1OVm7Cgux3IJX8gW8NblrF2I3NSW5vlEV4eI9wMMOnaKsOH-exAqKwMCZEBx86VQ2etb86u5x4FZiV86jGnQqB1iMdSqZ2DfAXWYg-ETJb-S5mXyltZBrZDAVn1KN85Oiqm7egSYUpljwA0ydXRg47IxLGlN1ZEgf9pZBaT8dU9yc1I6qoeAc6q0XAfVA9sq8Or_eT7mHI73v6fY5Lu77A"
}
```

### To get JWKS

```bash
curl --location 'http://localhost:8080/.well-known/jwks.json'
```

Response:
```json
{
  "keys": [
    {
      "kty": "RSA",
      "x5c": [
        "MIIDBjCCAe6gAwIBAgIBATANBgkqhkiG9w0BAQsFADArMRMwEQYDVQQKEwpFeGFtcGxlIENvMRQwEgYDVQQDEwtleGFtcGxlLmNvbTAeFw0yNDA1MzExODQ5NTdaFw0yNTA1MzExODQ5NTdaMCsxEzARBgNVBAoTCkV4YW1wbGUgQ28xFDASBgNVBAMTC2V4YW1wbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAslodEe4cM/0q2PzCKTkM8pGhKivPfYr6Gx5kHrj0ATL3zyczZ4ZMjycE5gBZ57CJ0t2sszr8UVzhL2geHhTWzFfKV4kHi+56HSTrfQ8Q25UJ4s8XbFNHtwwSipBpXBT5Ft5O1bEiij78plxt5eu1H0boO7ViRZqOQs+fyzPJEAUXkd+BwtQ+lEEAJZv9qiK5hWGIqLvPawq1ukKCBCBBokyG3w5JY3ZehFYcT/2ZGu3FHd0tYw0oM4IDQ7k5sv7AQeSvrvmfxHkgNul/X3OaePDBjGCZRKJ6NUJ/UqXaq1qgMmBcSniP0rja2tJwfi06tgLCRKO28XV4XFjfzmgj/QIDAQABozUwMzAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDAYDVR0TAQH/BAIwADANBgkqhkiG9w0BAQsFAAOCAQEAo5i7ryYV3+IMERJWF9fEqhPGEq1yylG2t1T6oNHTuNvLTSedtB6FIPNlMA3a3hzMDaMn/cu7I5jNOzS1WgSPgBZhg+yhHUu9hEl662rpZ3rSmHnYHA5X30k6jQTdwu5Oq8U2FMT9HOH3X9lNyTi6lIR8tYLOUiAqlSit1t8wMO381eWHUazv9wBCmEtwq0wA64uAbRsVY2nlDJuGCfCTBEUMlR6gQgZoVeNjOSU05ZoZDwvLCxsPIyg8N3KXksPqHfGSfbeqwytOUvsodGWk12CsES1O6aA3QW5eqIQc9WDf0kj4UAblSYPwGT0Uc4gfg6PiWpa/8+zGsBKsKT2zIg=="
      ],
      "n": "slodEe4cM_0q2PzCKTkM8pGhKivPfYr6Gx5kHrj0ATL3zyczZ4ZMjycE5gBZ57CJ0t2sszr8UVzhL2geHhTWzFfKV4kHi-56HSTrfQ8Q25UJ4s8XbFNHtwwSipBpXBT5Ft5O1bEiij78plxt5eu1H0boO7ViRZqOQs-fyzPJEAUXkd-BwtQ-lEEAJZv9qiK5hWGIqLvPawq1ukKCBCBBokyG3w5JY3ZehFYcT_2ZGu3FHd0tYw0oM4IDQ7k5sv7AQeSvrvmfxHkgNul_X3OaePDBjGCZRKJ6NUJ_UqXaq1qgMmBcSniP0rja2tJwfi06tgLCRKO28XV4XFjfzmgj_Q",
      "e": "AQAB",
      "use": "sig",
      "alg": "RS256"
    }
  ]
}
```

## How to verify the access_token

1. Open [https://jwt.io](https://jwt.io)
2. Pass the `access_token` to left panel
3. Copy only one JWK and pass to the `public key panel`

For example:

![img.png](asset/images/img.png)