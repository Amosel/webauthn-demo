# WebAuthn Server in Go

This repository contains a Golang implementation of a WebAuthn server. It provides functionality for user registration and authentication using the WebAuthn protocol.

## Features

- User registration with WebAuthn credentials
- User authentication using WebAuthn
- Session management for WebAuthn operations
- Pluggable storage for user credentials (currently using in-memory storage)

## Prerequisites

- Go 1.22 or higher
- Required dependencies (see `go.mod` file)

## Getting Started

1. Clone the repository:
   `git clone https://github.com/amosel/webauthn-demo.git`
2. Install the dependencies:
   `go mod download`
3. Generate TLS certificates:
   `openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt`
4. Run the server:
   `go run .`

5. Access the server at `https://localhost:443` in your web browser.

## Project Structure

- `main.go`: Entry point of the server application
- `models.go`: Defines the data models used in the application
- `utils.go`: Contains utility functions
- `webauthn_handlers.go`: Implements the HTTP handlers for WebAuthn operations
- `webauthn_store.go`: Provides the storage interface and implementation for WebAuthn sessions
- `webauthn_user.go`: Defines the WebAuthn user interface and related functions
- `static/`: Directory containing static files (e.g., HTML, CSS, JavaScript)

## Implementation Details

### main.go

- Initializes the WebAuthn session store
- Sets up router and handlers for WebAuthn operations
- Starts the HTTPS server

### models.go

- Defines `User`, `WebAuthnConfiguration`, `WebAuthn`, and `SessionRequest` structs

### utils.go

- Implements `isValidEmail` function to validate email addresses

### webauthn_handlers.go

- Defines handlers for beginning and finishing WebAuthn registration and login
- Includes functions to add credentials to users

### webauthn_store.go

- Manages in-progress WebAuthn registrations and logins
- Implements functions to begin and finish WebAuthn registration and login

### webauthn_user.go

- Implements the `WebAuthnUser` struct and required methods for WebAuthn user interface
- Includes functions to load WebAuthn credentials and create `WebAuthnUser` from `User`

### static/index.html

- Provides a simple front-end interface for testing WebAuthn registration and login
- Includes JavaScript functions to handle WebAuthn registration and login processes

## TODO

- Implement persistent storage for user credentials (e.g., database)
- Add support for additional WebAuthn options and configurations
- Enhance error handling and logging
- Implement user management functionality
- Add unit tests and integration tests
- Improve documentation and code comments

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
