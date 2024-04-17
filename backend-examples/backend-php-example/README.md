
# Stylist Authentication Service

This application provides a secure authentication service for the Freya's stylist addon to be used Stylist Interactive Demo. This emulates what a typical Freya customer would implement in their backend to distribute secure access tokens to their users for the Freya Stylist addon.

## Features

- RSA signing of messages for authentication.
- External service communication for token retrieval.
- Cross-Origin Resource Sharing (CORS) support.

## Installation

1. Ensure Docker is installed on your system

```bash
docker build --no-cache -t interactive-demo-backend-api .
```

## Running the Application

Start the application with Docker:

```bash
docker run -p 80:80 --name interactive-demo-backend-container interactive-demo-backend-api
```

This will start the application on `localhost:80`.

## API Usage

The service exposes one endpoint `/demo/v1/authenticate` which accepts POST requests containing JSON payloads with `user_id` and `company_code`.

Example request:

```json
{
  "user_id": "your_user_id",
  "company_code": "your_company_code"
}
```

On success, the service responds with a JSON containing the authentication token in this format:

```json
{
    "token": "eyJhbGci...du7QUyqxvyBII"
}
```