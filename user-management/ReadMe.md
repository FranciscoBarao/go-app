# Marketplace
Readme for the marketplace service


**Objective -** 

**Future -** 


## Entity Relationship

## User API
User JSON
```
{
    "Name": "name",
	"Username": "user",
	"Email": "email@email.com",
	"Password":"123"
}
```


Create
```
curl -X POST localhost:8082/api/register -H 'Content-Type: application/json' -d '{"name": "name", "username": "user", "email": "email@email.com", "password":"123"}'
```




## Links

[Auth in Go](https://codewithmukesh.com/blog/jwt-authentication-in-golang/)
[go chi oauth](https://github.com/go-chi/oauth)