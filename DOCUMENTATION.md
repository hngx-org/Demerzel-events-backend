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
   * [Users](#users)
   * [Events](#events)
   * [Comments]()
   * [Images](#images)
  
* [Request and Response Formats](#request-and-response-format)

## Introduction

Welcome to team Demerzel Event API documentation. Our API is designed to empower your mobile app with the ability to seamlessly manage events, groups, and user interactions. With this API, you can effortlessly create, read, update, and delete events, allowing your users to stay engaged and informed. Additionally, users can create groups, subscribe to their favorite communities, leave comments on events, and express their interest in attending. Dive into our comprehensive documentation to unlock the full potential of enhancing user experiences with event management and community engagement in your app.

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

* **GET /api/groups/{id}**
   * **Summary**: Get a group by Id
   * **Description**: Fetch a group details with Id
   * **Sample Request Url**: `{host}/api/groups/bcbb8231-e97d-461e-aa95-e2c426b68b28`
   * **Response**:  
      Status Code: 200
      Body: 
      ```Json
      {
         "data": {
            "id": "bcbb8231-e97d-461e-aa95-e2c426b68b28",
            "name": "My Group",
            "created_at": "2023-09-22T09:34:44.376Z",
            "updated_at": "2023-09-22T09:34:44.376Z",
            "members": [
                  {
                     "id": "a750e6c1-50e5-4ba3-84fc-73dbbb684acc",
                     "user_id": "6be00466-af9f-444c-81bc-969d575342be",
                     "group_id": "bcbb8231-e97d-461e-aa95-e2c426b68b28",
                     "created_at": "2023-09-22T09:55:59.673Z",
                     "updated_at": "2023-09-22T09:55:59.673Z",
                     "user": {
                        "id": "6be00466-af9f-444c-81bc-969d575342be",
                        "name": "user_name",
                        "email": "example@gmail.com",
                        "avatar": "avarter_url"
                     }
                  }
               ]
         },
         "message": "Group retrieved successfully",
         "status": "success"
      }

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
### Users
* **GET /api/users/current**
   * **Summary**: Retrieve details of the current User
   * **Description**: Get the detials of the current user whose bearer token is in the Authorization key in the header.
   * **Sample Request URL**: `{host}/api/users/current`
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
         "data": {
            "id": "664f1610-b0a0-45bc-9807-96ea7d1872af",
            "name": "user_name",
            "email": "example@gmail.com",
            "avatar": "avater_url"
         },
         "message": "User retrieved successfully",
         "status": "success"
      }
   
* **GET /api/users/{id}**
   * **Summary**: Retrieve a User by id
   * **Description**: Get the details of a User with User Id (uuid)
   * **Sample Request URL**: `{host}/api/users/664f1610-b0a0-45bc-2807-96ea7d1872ai`
   * **Response**:   
   Status Code: 200  
   Body:
      ```Json
      {
         "data": {
            "user": {
                  "id": "664f1610-b0a0-45bc-2807-96ea7d1872ai",
                  "name": "user_name",
                  "email": "example@gmail.com",
                  "avatar": "https://avater_url"
            }
         },
         "message": "User retrieved successfully",
         "status": "success"
      }
      ```
      Status Code: 404  
      Body:
      ```Json
      {
         "message": "User with the ID (664f1610-b0a0-45bc-9807-96ea7d1872a) does not exist",
         "status": "error"
      }
      ```
* **PUT /api/users/{id}**
   * **Summary**: Update details of a user
   * **Description**: Update the Name and Avater of a user, remember to set the bearer token in the request header.
   * **Sample Request URL**: `{host}/api/users/664f1610-b0a0-45bc-2807-96ea7d1872ai`
   * **Parameters**:  
   Body:  
      ```Json
      {
         "name": "new_user_name",
         "avater": "avater_url",
      }
      ```
      > **NB:** Both parameters (name and avater) are both optional, but either must be supplied.
   * **Response**:  
   Status Code: 200  
   Body:  
      ```Json
      {
         "data": null,
         "message": "User updated successfully",
         "status": "success"
      }
      ```
      Status Code: 400  
      Body:
      ```Json
      {
         "message": "Invalid request body format",
         "status": "error"
      }
      ```
* **GET /api/users**
   * **Summary**: Retrive all users
   * **Description**: Get a list of all users details: name, email and avatar
   * **Sample Request URL**: `{host}/api/image/upload`
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
         "data": {
            "user": [
                  {
                     "id": "664f1610-b0a0-45bc-9807-96ea7d4562af",
                     "name": "Username_1",
                     "email": "example1@gmail.com",
                     "avatar": "https://avatar.url"
                  },
                  {
                     "id": "6be00466-af9f-444c-81bc-969d523442be",
                     "name": "Username_2",
                     "email": "example2@gmail.com",
                     "avatar": "https://avatar.url"
                  }
            ]
         },
         "message": "Users Retrieved Successfully",
         "status": "success"
      }

* **POST api/users/logout**
   * **Summary**: Logout the current user
   * **Sample Request URL**: `{host}/api/users/logout`
   * **Response**:   
   Status Code: 200  
   Body:
      ```Json
      {
         "data": null,
         "message": "logged out successfully",
         "status": "success"
      }
      ```

### Events
* **POST /api/event/**
   * **Summary**: Create an Event
   * **Description**: Creates a new Event and returns it, supply the title of the event, event location, event description, Start Date, End Date, Start Time and End Time.
   * **Sample Request URL**: `{host}/event`
   * **Parameters**:  
   Body:  
      ```Json
      {
         "CreatorId": "user_uuid",
         "Thumbnail": "thumbnail",
         "Location": "location",
         "Title": "title",
         "Description": "event_description",
         "Start_Date": "2023-10-01",
         "End_Date": "2023-10-01",
         "Start_Time": "16:00",
         "End_Time": "18:00"
      }
      ```
   * **Response**:  
   Status Code: 200    
   Body:
   ```Json
   {
      "data": {
         "event": {
               "id": "d7515060-9276-40f0-a294-6f64c6cead54",
               "creator_id": "664f1610-b0a0-45bc-9807-96ea7d1872af",
               "thumbnail": "",
               "location": "location",
               "title": "event",
               "description": "event_description",
               "start_date": "2023-10-01",
               "end_date": "2023-10-01",
               "start_time": "16:00",
               "end_time": "18:00",
               "created_at": "2023-09-22T14:52:54.785Z",
               "updated_at": "2023-09-22T14:52:54.785Z",
               "creator": null
         }
      },
      "message": "Event Created",
      "status": "success"
   }
   ```
* **GET api/events**
   * **Summary**: Get a list of all Events
   * **Sample Request URL**: `{host}/api/events`
   * **Response**:  
   Status Code: 200  
   Body:
   ```Json
   {
      "data": {
         "events": [
               {
                  "id": "3d5c9a49-a918-4421-9ce0-e4e15abd3389",
                  "creator_id": "a6d0f5bf-b923-47a5-8406-746551dca395",
                  "thumbnail": "",
                  "location": "Uyo, Nigeria",
                  "title": "Event Two",
                  "description": "event_description",
                  "start_date": "2023-09-21",
                  "end_date": "2023-09-25",
                  "start_time": "15:1",
                  "end_time": "17:1",
                  "created_at": "2023-09-22T14:01:22.326Z",
                  "updated_at": "2023-09-22T14:01:22.326Z",
                  "creator": {
                     "id": "a6d0f5bf-b923-47a5-8406-746551dca395",
                     "name": "user_name",
                     "email": "example@gmail.com",
                     "avatar": "https://lh3.googleusercontent.com/a/ACg8ocJ9wGidl0HuJH6tw3cg42Zvd_OoNUbfSaiPg3VDED1aGpI=s96-c"
                  }
               }
         ]
      },
      "message": "All Events",
      "status": "success"
   }
* **GET api/events/{id}**
   * **Summary**: Fetch an Event by Id
   * **Sample Request URL**: `{host}/api/events/d7515060-9276-40f0-a294-6f64c6cead54
   * **Response**:   
   Status Code: 200  
   Body: 
      ```Json
      {
         "data": {
            "events": [
                  {
                     "id": "3d5c9a49-a918-4421-9ce0-e4e15abd3389",
                     "creator_id": "a6d0f5bf-b923-47a5-8406-746551dca395",
                     "thumbnail": "",
                     "location": "Uyo, Nigeria",
                     "title": "Event Two",
                     "description": "event_description",
                     "start_date": "2023-09-21",
                     "end_date": "2023-09-25",
                     "start_time": "15:1",
                     "end_time": "17:1",
                     "created_at": "2023-09-22T14:01:22.326Z",
                     "updated_at": "2023-09-22T14:01:22.326Z",
                     "creator": {
                        "id": "a6d0f5bf-b923-47a5-8406-746551dca395",
                        "name": "user_name",
                        "email": "example@gmail.com",
                        "avatar": "https://lh3.googleusercontent.com/a/ACg8ocJ9wGidl0HuJH6tw3cg42Zvd_OoNUbfSaiPg3VDED1aGpI=s96-c"
                     }
                  }
            ]
         },
         "message": "All Events",
         "status": "success"
      }
      ```

* **GET api/events/comments/{id}**
   * **Summary**: Fetch all comments tied to an event.
   * **Sample Request URL**: `{host}/api/events/comments/d7515060-9276-40f0-a294-6f64c6cead54`
   * **Response**:  
   Status Code: 200  
   Body:
      ```Json
      {
         "data": {
            "comments": []
         },
         "message": "Comments Successfully retrieved",
         "status": "success"
      }
      ```



### Images
* **POST api/images/upload**
   * **Summary**: Upload an Image
   * **Description**: Upload an image and get a Cloudinary url of the hosted image. Set the Content-Type header to `multipart/form-data`. Name the form field containing the image file `file`.
   * **Sample Request URL**: `{host}/api/images/upload`
   * **Parameters**:  
   Body:  
      ```Json
      "file": "url_to_file/image.jpeg"
      ```
   * **Responses**:  
      Status Code: 200  
      Body:
      ```Json
      {
         "status": "success"
         "data": {
            "url": "https://res.cloudinary.com/jamesjohn/image/upload/v1695389554/1695389551-1660372349874.jpeg.jpg"
         },
         "message": "File uploaded",
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

## Request and Response Format
### Response:
The Api Response follow the JSend format, whcih have a `status`, `data` and `message` key, the status falls under either `success` or `error` respectfully.  
Read more at:  [The JSend Specification](https://github.com/omniti-labs/jsend)
