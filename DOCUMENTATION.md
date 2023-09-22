# REST API Documentation

## Table of Contents

1. Introduction
   
2. API Usage
   - 2.1 How to Call the API
   - 2.2 Authenticating to the API
   - 2.3 Supported Endpoints
  
3. Request and Response Formats
   - 3.1 Endpoint group 1
  
4. Sample API Calls
   - 4.1 Endpoint group 1
   
5. Setting Up and Running the API in a Container
   - 5.1 Environment Variables
   - 5.2 Docker Container Setup
  
6. Additional Info
---



## 1. Introduction

API Intro


## 2. API Usage

### 2.1 How to Call the API

The API can be accessed via HTTP requests. It exposes endpoints for various CRUD operations.

Sample API base URL: `http://example.com/api`

Current [Active] base URL: `url`


### 2.2 Authenticating to the API

The API supports the following Authentication Workflow(s):

- **Google OAuth**: `Endpoint URL`
  - **Sample Request**
  - **Sample Response**


### 2.3 Supported Endpoints

The API supports the following Endpoints:

- **Endpoint Group 1**: `!`


## 3. Request and Response Formats

### 3.1 Endpoint Group 1



## 4. Sample API Calls

Here are some sample API calls:



## 5. Setting Up and Running the API in a Container

### 5.1 Environment Variables

To run the API in a container, you'll need to set the following environment variables:

- `APP_ENV`: The current Application environment ('local' or 'prod').
- `PORT`: Access port for the API App.
- `MYSQL_HOST`: Hostname of the MySQL database server.
- `MYSQL_PORT`: Port of the MySQL database server.
- `MYSQL_USERNAME`: MySQL database username.
- `MYSQL_PASSWORD`: MySQL database password.
- `MYSQL_DBNAME`: Name of the MySQL database.
- `GOOGLE_CLIENT_ID`: Google OAuth Client_ID
- `GOOGLE_CLIENT_SECRET`: Google OAuth Client Secret
- `GOOGLE_CALLBACK_URL`: Google OAuth Callback URL
- `JWT_SECRET`: JWT Secret key.

### 5.2 Docker Container Setup

1. Build the Docker image: Replace {repo_access_key} with the approriate Repository access key.
   
   ```sh
   docker build -t your-api-image --build-arg="ACCESS_KEY={repo_access_key}" .

   docker run -d -p 8080:8080 \
   -e APP_ENV=app-env \
   -e PORT=app-port \
   -e MYSQL_HOST=your-db-hostname \
   -e MYSQL_PORT=your-db-port \
   -e MYSQL_USERNAME=your-db-username \
   -e MYSQL_PASSWORD=your-db-password \
   -e MYSQL_DBNAME=your-db-name \
   -e GOOGLE_CLIENT_ID=google-client-id \
   -e GOOGLE_CLIENT_SECRET=google-client-secret \ 
   -e GOOGLE_CALLBACK_URL=google-callback-url \ 
   -e JWT_SECRET=jwt-secret \
   --name your-container-name \
   your-api-image

The API should now be accessible at http://localhost:8080


## 6. Additional Info

- This repository contains Github Actions workflow that builds the App into a Docker image and uploads it to a private Azure Container Registry instance. 