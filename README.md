# Yearbook

A full-featured social media platform for congratulating the class of 2020. This repository provides the backend code, written in Go. The frontend code is found
[here](https://github.com/mattnappo/react-yearbook).

# Architecture

The core database is a PostgreSQL instance. To interface with the database, an ORM layer is added for safer transactions. The database stores all posts, user
metadata, and current OAuth tokens. All authentication/authorization is handled through Google OAuth for better security. A bearer token is sent and validated
with every request from the frontend. In front of the ORM lies a REST API running on a web server which abstracts all database CRUD operations and authentication/authorization methods. The frontend can then easily make calls to this API to perform any necessary operation.

# Usage

I deployed this during the spring of my junior year of high school. It was a big success and had over a hundred daily active users!
