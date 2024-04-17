
# Stylist Authentication Service

This application provides a secure authentication service for the Freya's stylist addon to be used Stylist Interactive Demo. This emulates what a typical Freya customer would implement in their backend to distribute secure access tokens to their users for the Freya Stylist addon.

## Features

- RSA signing of messages for authentication.
- External service communication for token retrieval.
- Cross-Origin Resource Sharing (CORS) support.

## Requirements

- Go 1.22.2

## Installation

1. Ensure Go 1.22.2 is installed on your system.

## Running the Application

Start the application with Uvicorn:

```bash
cd src/freya_customer_backend_go_demo
go run main.go
```

This will start the application on `localhost:8000`.

## Docker Support

Docker build commands are provided for Linux, Windows, and macOS platforms:

- For Linux and Windows:

```bash
docker build --no-cache -t interactive-demo-backend-api .
```

- For macOS:

```bash
docker buildx build --platform linux/amd64 -t interactive-demo-backend-api .
```


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