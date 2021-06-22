# Notes

## Authentication by x509 certificates using Golang

Heimdallr uses public key (or asymmetric) authentication:

- The client uses a private key to sign a **JSON Web Token (JWT)**. The token is passed to the service as a proof of the client identity.
- The service uses the client public key (uploaded before the JWT is sent) to verify the client's identity.

Steps:

- Generate the RSA key
- Generate an RSA key with a self signed x509 certificate
- Generate keys
- Store public key in database and redis

Operations:

- Create client (receives id and generates key pair, returns private key)
- Authenticate (receives jwt and authenticate with public key, returns result)
- Black listing (receives id and returns result)
- Delete client (receives id and returns result)
- Regenerate keys (receive id and regenerate key pair, returns private key)

DB:

- Redis
