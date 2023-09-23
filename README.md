# Team Demerzel - Backend

## 1. About Project
Team Demerzel is a HNG X Backend team of Go Developers.
This project is an API for an 'Events' Mobile Application built by Team Demerzel - Mobile.

## 2. Documentation
>The API from this project is fully Documented [HERE](./DOCUMENTATION.md)

## 3. Stacks and Technologies used
- [Golang](https://go.dev/)
- [MySQL](https://www.mysql.com/)
  
## 4. Prerequisites 
 - Golang v >= 1.19
 - Mysql v >= 8.0

## 5. Getting Access
### 5.1 As a Contributor

#### Fork repo to personal github account
Your Repo's found at https://github.com/hngx-org/Demerzel-events-backend
So to work with Forks you basically:
1. Fork your Team Repo to your personal Github account
2. Pull the code back to your local Machine
3. Checkout to your assignment branch
4. Do your thing
5. Push back to your Personal Github Repo. That'll be your 'ORIGIN' remote (not the 'UPSTREAM' remote)
6. You head over to Github and Create a Pull request to the Main Repository's branch
* Remember a Pull Request can contain multiple commits.
That's it. Here's a [video](https://youtu.be/nT8KGYVurIU) to help further.


> Additional help
- To clone repo from your account
```bash
git clone https://github.com/{your username}/Demerzel-events-backend
```

## 6. Run locally 
### To Run API Server Locally
 - Create a .env file
```bash
cp .env.example .env
```
 - Edit File with the following approriate values
```text
APP_ENV=app-env
PORT=app-port
MYSQL_HOST=your-db-hostname
MYSQL_PORT=your-db-port
MYSQL_USERNAME=your-db-username
MYSQL_PASSWORD=your-db-password
MYSQL_DBNAME=your-db-name
GOOGLE_CLIENT_ID=google-client-id
GOOGLE_CLIENT_SECRET=google-client-secret
GOOGLE_CALLBACK_URL=google-callback-url
JWT_SECRET=jwt-secret
```
 - Run server
```bash
go run main.go
```

## 7. Setting Up and Running the API in a Container

### 7.1 Environment Variables

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

### 7.2 Docker Container Setup

1. Build the Docker image: Replace {repo_access_key} with the approriate Repository access key.
   
   ```bash
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


## 8. Additional Info(s)

- This repository contains Github Actions workflow that builds the App into a Docker image and uploads it to a private Azure Container Registry instance.