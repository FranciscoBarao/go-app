# go-app
Repository for a Go-app 


[Project Structure using Repository Pattern](https://dakaii.medium.com/repository-pattern-in-golang-d22d3fa76d91)

[Decoding JSON Body](https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body)




Boardgame JSON
```
{
    "Name": "Test",
	"Dealer": "Dealer",
	"Price": 10.0,
	"PlayerNumber": 1
}
```



Create
```
curl -X POST localhost:8080/api/boardgame -H 'Content-Type: application/json' -d '{ "Name": "Test", "Dealer": "Dealer", "Price": 10.0, "PlayerNumber": 1 }'
```

ReadAll
```
curl -X GET localhost:8080/api/boardgame
```


Read
```
curl -X GET localhost:8080/api/boardgame/Test
```

Update
```
curl -X PATCH localhost:8080/api/boardgame/2 -H 'Content-Type: application/json' -d '{ "Name": "Test", "Dealer": "Dealer", "Price": 10.0, "PlayerNumber": 1 }'
```

Delete
```
curl -X DELETE localhost:8080/api/boardgame/2
```