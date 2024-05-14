# Marketplace
Readme for the marketplace service


**Objective -** 

**Future -** 


## Entity Relationship

## User API
User JSON
```
{
	"Username": "user",
	"Email": "email@email.com",
	"Password":"123"
}
```


Register a user 
```
curl -X POST localhost:8082/api/register -H 'Content-Type: application/json' -d '{"username": "user", "email": "email@email.com", "password":"123"}'
```

Login 
```
curl -X POST localhost:8082/api/login -H 'Content-Type: application/json' -d '{"username": "user", "email": "email@email.com", "password":"123"}'
```



## Links

[Auth in Go](https://codewithmukesh.com/blog/jwt-authentication-in-golang/)
[go chi oauth](https://github.com/go-chi/oauth)