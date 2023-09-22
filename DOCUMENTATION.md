# Team Demerzel Events API Documentation

## Table of Contents

* [Introduction](#introduction)
   
* [API Usage and Features](#api-usage-and-features)
   * [How to Call the API](#how-to-call-the-api)
   * [Authenticating to the API](#authenticating-to-the-api)
   * [Security Definitions](#security-definitions)

* [API Endpoints](#api-endpoints)
   * [API Health](#api-health)
   * [Authentication](#authentication)
   * [Groups](#groups)
   * [Users]()
   * [Events]()
   * [Comments]()
   * [Images](#images)
  
* [Request and Response Formats]()

## Introduction

...

## API Usage and Features

### How to Call the API

The API can be accessed via HTTP requests. It exposes endpoints for various CRUD
operations on events, users and groups. An Authentication Bearer Token needs to 
be set in the request header, to get the JWT token, read how to Authenticate in the next section.  
* **local Host**: `localHost:5005/`  
* **Current (Active) Host**: `https://api-s65g.onrender.com/`  
* **API Base Path**: `/api`
```
Authorization : Bearer <replace_with_jwt_token>
```

### Authenticating to the API

The API supports the Google OAuth Authentication Workflow:  
To initialize the OAuthentication workflow, the endpoint: `GET /oauth/initialize`,
should be accessed, upon success a redirection_url is sent along with the response. 
If successful, the JWT Token to be used in subsequent request
header should be gotten from the response body.

* **Path**: `{host}/oauth/initialize`
* **Sample Request**
   ```
   Request URL: {host}/oauth/initialize
   Request Method: GET
   ```
* **Sample Response**
   * Status Code: 200
   * Body: 
      ```JSON
      {
         "status": "success"
         "data": {
            "redirectUrl": "https://accounts.google.com/o/oauth2/auth?client_id=1024268478019-v7ebrdubfni6qsvtpiicv9eb2le6m4al.apps.googleusercontent.com&redirect_uri=https%3A%2F%2Fapi-s65g.onrender.com%2Foauth%2Fcallback&response_type=code&scope=openid+profile+email&state=somerandomoauthstri"
         },
         "message": "Authentication initialized",
      }
      ```
* **After Oauth Redirection**
   * Response Body:
      ```JSON
      {
         "status":"success",
         "data":{"token": "jwt_token"},
         "message":"Authentication successful"
      }
### Security Definitions
   ```Json
   {
      "securityDefinitions": {
         "Bearer": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header",
            "description": "Enter your Bearer token in the format: Bearer <token>"
         }
      }
   }
   ```


## API Endpoints
### API Health
* **GET /health**
   * **Sample Request URL**: `{host}/health `
   * **Response**:  
   Status Code: 200  
   Body:
      ```Json
      {
         "data": null,
         "message": "Team Demerzel Events API",
         "status": "success"
      }
### Authentication
* **GET /oauth/initialize**  
   * **Summary**: Initialize the Oauth workflow  
   * **Despcription**: Authenticate a user using Google's Oauth and get a JWT token for 
   subsequent requests.  
   **NB**: Read the [section above](#authenticating-to-the-api) to get more information
### Groups
* **POST /api/groups**
   * **Summary**: Create a group
   * **Description**: Creates a new group and returns it, supply the name of the group 
   in the request body.
   * **Sample Request URL**: `{host}/api/groups`
   * **Parameters**:   
      Body:
      ```Json
      {
      "name": "group_name"
      }
      ```
   * **Responses**:  
      Status Code: 201  
      Body:
      ```Json
      {
         "status": "success",
         "data": {
            "id": "a2218d8f-4cdb-4114-a847-4cf8fcbd2e54",
            "name": "group_name",
            "created_at": "2023-09-20T18:28:42.523+01:00",
            "updated_at": "2023-09-20T18:28:42.523+01:00",
            "members": null
         },
         "message": "Group created successfully"
      }
      ```
      Status Code: 400  
      Body:
      ```Json
      {
         "status": "error",
         "message": "Invalid request body format",
         "data": null
      }
      ```
* **PUT /api/groups/{id}**
   * **Summary**: Update a group
   * **Description**: Update an existing group by ID, supply the new name of the
   group in the request body.
   * **Sample Request URL**: `{host}/api/groups/9d8be7c8-d2f7-4e4a-8897-9a08e5a18340`
   * **Parameters**:  
      Body:  
      ```Json
      {
      "name": "new_group_name"
      }
      ```
   * **Response**:  
   Status Code: 200  
   Body:
      ```Json
      {
         "data": {
            "id": "9d8be7c8-d2f7-4e4a-8897-9a08e5a18340",
            "name": "new_group_name",
            "created_at": "2023-09-22T05:09:44.999Z",
            "updated_at": "2023-09-22T05:15:07.97Z",
            "members": null
         },
         "message": "Group updated successfully",
         "status": "success"
      }
      ```
* **POST /api/groups/{id}/subscribe**  
   * **Summary**: Subscribe to a group
   * **Description**: Subscribe a user to a group by ID, it is required that the
   user bearer token is added to the request header.
   * **Sample Request URL**: `{host}/api/groups/9d8be7c8-d2f7-4e4a-8897-9a08e5a18340/subscribe`
   * **Security**: 
      ```Json
      "security": [
            { "Bearer": [] }
            ]
      ```
   * **Responses**:  
      Status Code: 200  
      Body:  
      ```Json
      {
         "data": {
            "id": "7078c94c-0138-4488-bc74-e3be3f53271e",
            "user_id": "864f1310-b0a0-45uc-9407-90ea7e1871zf",
            "group_id": "9d8be7c8-d2f7-4e4a-8897-9a08e5a18340",
            "created_at": "2023-09-22T05:50:13.262Z",
            "updated_at": "2023-09-22T05:50:13.262Z"
         },
         "message": "User successfully subscribed to group",
         "status": "success"
      }
      ```   
      Status Code: 404  
      Body:  
      ```Json
      {
         "message": "Invalid user or group ID. Please check the values and try again",
         "status": "error"
      }
      ```
      Status Code: 409  
      Body:
      ```Json
      {
         "message": "User already subscribed to group",
         "status": "error"
      }
      ```

* **POST /api/groups/{id}/unsubscribe**
   * **Summary**: Unsubscribe to a group
   * **Description**: Unsubscribe a user from a group by ID, it is required that the
   user bearer token is added to the request header.
   * **Sample Request URL**: `{host}/api/groups/9d8be7c8-d2f7-4e4a-8897-9a08e5a18340/unsubscribe`
   * **Security**: 
      ```Json
      "security": [
            { "Bearer": [] }
            ]
      ```
   * **Responses**:
      Status Code: 200  
      Body: 
      ```Json
      {
         "data": null,
         "message": "User successfully unsubscribed to group",
         "status": "success"
      }
      ```
      Status Code: 404  
      Body:
      ```Json
      {
         "message": "User not subscribed to this group",
         "status": "error"
      }
      ```
      Status Code: 409  
      Body:
      ```Json
       {
         "message": "Failed to unsubscribe user from group",
         "status": "error"
      }
      ```
* **GET /api/groups/user**
   * **Summary**: Get groups a user is subscribed to.
   * **Description**: Get a list of groups a user is subscribed to, it is required that the user bearer token is added to the request header.
   * **Sample Request URL**: `{host}/api/groups/user`
   * **Security**: 
      ```Json
      "security": [
            { "Bearer": [] }
            ]
      ```
   * **Response**:  
      Status Code: 200  
      Body: 
      ```Json
      {
         "data": [
            {
                  "id": "9d8be7c8-d2f7-4e4a-8897-9a08e5a18340",
                  "name": "new api group",
                  "created_at": "2023-09-22T05:09:44.999Z",
                  "updated_at": "2023-09-22T05:15:07.97Z",
                  "members": null
            }
         ],
         "message": "Fetched all user groups",
         "status": "success"
      }
      ```
* **GET /api/groups?name={name}**
   * **Summary**: Fetch a group by name
   * **Description**: Get the details of a group by it's name
   * **Sample Request URL**: `{host}/api/groups?name=Demerzel`
   * **Parameters**:  
      In Query:  
      ```Json
      {
      "name": "Demerzel"
      }
      ```
   * **Responses**:  
      Status Code: 200  
      Body:  
      ```Json
      {
         "data": [
            {
                  "id": "dc12ee4e-cc8c-4231-a0a7-7ae3148c6392",
                  "name": "Demerzel",
                  "created_at": "2023-09-22T06:36:49.415Z",
                  "updated_at": "2023-09-22T06:36:49.415Z",
                  "members": null
            }
         ],
         "message": "Groups retrieved successfully",
         "status": "success"
      }
      ```   
      Status Code: 200  
      Body:
      ```Json
      {
         "data": [],
         "message": "No groups",
         "status": "success"
      }
      ```
### Images
* **POST api/images/upload**
   * **Summary**: Upload an Image
   * **Description**: Upload an image and get a Cloudinary url of the hosted image. Set the Content-Type header to `multipart/form-data`. Name the form field containing the image file `file`.
   * **Sample Request URL**: `{host}/api/images/upload`
   * **Parameters**:  
   Body
      ```Json
      "file": "url_to_file/image.jpeg"
      ```
   * **Responses**:  
      Status Code: 200  
      Body:
      ```Json
      {
         "status": "success",
         "url": "cloudinary_url",
         "message": "File uploaded"
      }
      ```
      Status Code: 400  
      Body:
      ```Json
      {
         "status": "error",
         "message": "Unable to upload file: <error_message>" 
      }
      ```
