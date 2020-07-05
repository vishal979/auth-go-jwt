## Authentication in go using jwt access token, jwt refresh token persisting token metadata in redis, and users info in Users table in testdb(MySQL).

Project Structure
```
├── README.md
├── SQLQueries.md
├── cmd
│   └── filehandler
├── docker-compose.yml
├── dockerfile
├── go.mod
├── go.sum
├── logs
├── main.go
├── makefile
├── models
│   └── User.go
├── server
│   ├── auth
│   ├── controller
│   ├── middlewares
│   ├── responses
│   └── server.go
└── vendor
```

The application generates real time logs, in ./logs/server.logs file, which can be put to watch via 
```
tail -f logs/server.logs
```

### to run the project, use makefile
The default goal for make is help, and it looks like this
```
Usage:
  make [target...]

Useful commands:
  build                          to build the project again after making changes
  compose                        to run the containers
  down                           docker-compose down
  pruneVolume                    remove all dangling volumes
  runLocal                       to run the app locally
```

### To run the containers
```
make compose
```
this will remove all the dangling volumes first, then build the project using flag no-cache, and then run the containers.

### expected result
```
docker container ls
```
```
CONTAINER ID        IMAGE                     COMMAND                  CREATED             STATUS                    PORTS                               NAMES
fb91b303db93        docker-contribution_app   "./main"                 40 seconds ago      Up 38 seconds             0.0.0.0:8080->8080/tcp              goapp
19a5b083084d        redis:latest              "docker-entrypoint.s…"   40 seconds ago      Up 39 seconds             0.0.0.0:6378->6379/tcp              goapp_redis
3da52c78e384        mysql/mysql-server:5.7    "/entrypoint.sh mysq…"   40 seconds ago      Up 39 seconds (healthy)   33060/tcp, 0.0.0.0:3307->3306/tcp   goapp_mysql
```




### Data
Redis is used to store metadata for the jwt-token, while mysql is used to store users information. The Users table has 3 entries username, password and id(auto increment), indexing is done on username for fast searches. You can use ```id``` as foreign key in some other table as per requirement.

### Routes 
The application has 5 routes

```
── /login        Post
── /signup       Post
── /logout       Post
── /refreshtoken Post
── validatetoken Post
```

### login
Accepts and returns JSON
```
Input JSON
{
    "username":"<username>",
    "password":"<password>"
}
```
```
Output JSON
{
    "access_token": "eyJhbGciOiJIUzabddIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjRiY2M5M2QzLTU0MmYtNDQyNS05NGUwLWY4MTk0MGE5NjBlZiIsImF1dGhvcml6ZWQiOnRydWUsImV4cCI6MTU5Mzk2Mjc1NywidXNlcl9pZCI6Mn0.ZkQsdfYj0tp4_tALSiIrbGbswjEVfoSYPvKsKveBY",
    "refresh_token": "eyJhbGciOiJIUzIasdsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQ1NjM5NTcsInJlZnJlc2hfdXVpZCI6IjhkNzQ0ZDgxLTdlMzEtNDgyOS1hZmM2LTU5ZWViOGE0OTBiZCIsInVzZXJfaWQiOjJ9.bSorl7aWpU33nvEpnR6POmdsffVfXnBf_mD9Lp6Y"
}
```

### signup
Accepts JSON and returns successful or error
```
{
    "username":"<username>",
    "password":"<password>"
}
```

### logout
Accepts the bearertoken in ```Authorization Header```.

### refreshtoken
Accepts refreshtoken as JSON and returns accesstoken and refresh token as json.
```
Input JSON
{
    "refresh_token": "eyJhbGciOiJIUzIasdsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQ1NjM5NTcsInJlZnJlc2hfdXVpZCI6IjhkNzQ0ZDgxLTdlMzEtNDgyOS1hZmM2LTU5ZWViOGE0OTBiZCIsInVzZXJfaWQiOjJ9.bSorl7aWpU33nvEpnR6POmdsffVfXnBf_mD9Lp6Y"
}
```

```
Output JSON
{
    "access_token": "eyJhbGciOiJIUzabddIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjRiY2M5M2QzLTU0MmYtNDQyNS05NGUwLWY4MTk0MGE5NjBlZiIsImF1dGhvcml6ZWQiOnRydWUsImV4cCI6MTU5Mzk2Mjc1NywidXNlcl9pZCI6Mn0.ZkQsdfYj0tp4_tALSiIrbGbswjEVfoSYPvKsKveBY",
    "refresh_token": "eyJhbGciOiJIUzIasdsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQ1NjM5NTcsInJlZnJlc2hfdXVpZCI6IjhkNzQ0ZDgxLTdlMzEtNDgyOS1hZmM2LTU5ZWViOGE0OTBiZCIsInVzZXJfaWQiOjJ9.bSorl7aWpU33nvEpnR6POmdsffVfXnBf_mD9Lp6Y"
}
```

### validatetoken
This acts as middleware, accepts accesstoken in ```Authorization Header``` and returns success or error if the token could not be verified.
