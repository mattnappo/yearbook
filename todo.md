# TODO

## Security
 * Secure the size of images (set a cap. compression?)
 * Better SQL-injection protection
 * Have the backend handle routing?
 * What's going on with sessions?

## Scalability
 * Scale the backend sessions (PG? Redis?)
 * Optimize the UpdateUser and AddToAndFrom database methods (one DB call)

## Features
 * Profile pic should be either base64 OR a Google link
 * Add all seniors to the DB

## Network
 * Docker / k8
 * HTTPS w/ Let's Encrypt OR Cloudflare
 * Traefik integration
 * Caching layer

## Bugs
 * Better error handling (stop ignoring some errors)
 * Write more unit tests
