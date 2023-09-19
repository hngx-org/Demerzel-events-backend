# Team Demerzel

## About Project

## Stack and Technologies used
- [Golang](https://go.dev/)
- [MySQL](https://www.mysql.com/)
## Prerequisites 
Golang v >= 1.19
Mysql v >= 8.0
## Installing
    
### As a maintainer

### Fork repo to personal github account
Your Repo's found at https://github.com/hngx-org/Demerzel-events-backend
So to work with Forks you basically:
1. Fork your Team Repo to your personal Github account
2. Pull the code back to your local Machine
3. Checkout to your assignment branch
4. Do your thing
5. Push back to your Personal Github Repo. That'll be your 'ORIGIN' remote (not the 'UPSTREAM' remote)
6. You head over to Github and Create a Pull request to the Main Repository's branch
* Remember that a Pull Request can contain multiple commits.
That's basically it. Here's a [video](https://youtu.be/nT8KGYVurIU) to help further.

>After forking go over the step.
- Clone repo from your account
```bash
git clone https://github.com/{your username}/Demerzel-events-backend
```

## Run locally 
### Steps on how to run server locally
> Make sure you have setup MySQl env variables
```bash
cp .env.example > .env

# Make sure you change the values to your local MySQL values
MYSQL_USERNAME=root
MYSQL_PASSWORD=password
MYSQL_HOST=127.0.0.1
MYSQL_PORT=3306
MYSQL_DBNAME=
```
Run server
```bash
cd {filename}

go get

go run main.go
```