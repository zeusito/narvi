# Narvi Backend API

This project contains the backend API for the Narvi platform.

## Features

- Chi Router for HTTP-based endpoints
- Zerolog for logging capabilities
- Koanf for configuration, supports files and env vars
- PGX and Bun for PostgreSQL database access
- DBMate for database migrations
- Session management (using Opaque tokens) and HTTP filter to protect endpoints
- Hashing algorithms, including argon2id
- Makefile with the most common tasks
- Multi-stage Dockerfile for building and running the application
- A basic authentication module

## Getting Started

- Check out the `Makefile` for more information about available commands
- Check out the AGENTS.md file for more information about coding standards

## Folder Structure

- cmd/main.go - main entry point
- internal - application-specific business logic
- pkg – shared packages that might be used in multiple modules (follows the Unix philosophy of simple tools that make one thing)
- resources - application-specific resources, such as config files, databases, etc.
