# Online Cinema Service
The project is an online cinema service that allows 
users to view the schedule, purchase tickets, view the list of movies they have watched, and generate PDF documents for purchased tickets. 

## Description
#### 1 - Authentication

The application has a secure
authentication and authorization system to protect the endpoints, with JWT tokens
used to authenticate and authorize user requests.

#### 2 - API
The service has a RESTful API that supports
CRUD operations for admins to add, edit, and delete halls, movies, and cinema
sessions.

## Usage
**1.** [Install golang](https://go.dev/doc/install)  
**2.** Download repository from Bitbucket
```shell
go get "bitbucket.org/Ernst_Dzeravianka/cinemago-app"
```
**2.** Set up environment variables:
- Create .env file
- Set `PORT` variable to port on which the server will listen for incoming connections.

**3.** Run web service using Makefile:
```shell
make run
```
## Testing
Test the service by using Makefile:
```shell
make test
```

