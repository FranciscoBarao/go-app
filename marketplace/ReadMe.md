# Marketplace
Readme for the marketplace service


**Objective -** 

**Future -** 


## Entity Relationship

## Offer API
Offer JSON
```
{
    "Name": "name",
	"Price": 10.0
}
```


Create
```
curl -X POST localhost:8081/api/offer -H 'Content-Type: application/json' -d '{ "type": "Boardgame","name": "name", "price": 10.0}'
```

GetAll
```
curl -X GET localhost:8081/api/offer
```

Get
```
curl -X GET localhost:8081/api/offer/{id}
```

Update
```
curl -X PATCH localhost:8081/api/offer/{id} -H 'Content-Type: application/json' -d '{ "type": "Boardgame","name": "name", "price": 20.0}'
```

Delete
```
curl -X DELETE localhost:8081/api/offer/{id}
```

