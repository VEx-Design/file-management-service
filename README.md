# File Management Service

## Introduction
File Management Service is a system designed to efficiently handle file storage, retrieval, and management. This project provides a structured approach to organizing files, making it easier to store, access, and manage different types of files securely.

## Features
- Upload and download files easily
- Organized file storage system
- Secure access and retrieval
- Open-source and customizable

## Installation Dependencies
Ensure you have Go installed. After cloning the project, install the necessary packages and libraries with the following command:
```bash
go mod tidy
```

## Environment Variables
Create a `.env` file in the root directory and configure it based on `.env.example`.

## Running the Project
To run the project, use the following command:

```bash
make restart
```

## Docker
### Building the Docker Image and Running the Container
```bash
make up
```


