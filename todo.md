# TODO

## Security
 * Secure the size of [][]byte images
 * Encrypt the environment variables 
 * Better SQL-injection protection (? instances)
 * Docker / k8
 * Allow only masters emails
 * HTTPS w/ Let's Encrypt OR Cloudflare
 * Scale backend sessions (PG? Redis?)
 * Implement refresh tokens
 * Have the backend handle routing
 * Should getUserInfo also check the token DB?

## Features
 * DB & API methods to get all posts from a person, a posts to a person
 * Traefik integration
 * DB & API methods to set bio, profile picture, senior will.
 * Search accounts
 * Make profile picture a link to the google API link

## Bugs
 * Better error checking in Google API call function from google.
 * Add "registered" tag and create un-established users on post create with registered = false. Set registered = true when registering happens.
 * Make an AddUserIfNotExists databsae method
