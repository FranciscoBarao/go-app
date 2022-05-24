# go-app
Repository for a Go-app 


[Project Structure using Repository Pattern](https://dakaii.medium.com/repository-pattern-in-golang-d22d3fa76d91)

[Decoding JSON Body](https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body)



## Boardgame API

Boardgame JSON
```
{
    "Name": "Test",
	"Dealer": "Dealer",
	"Price": 10.0,
	"PlayerNumber": 1,
	"Tags": [
		{ "Name": "A" },
		{ "Name": "B" }
	]
}
```


Create
```
curl -X POST localhost:8080/api/boardgame -H 'Content-Type: application/json' -d '{ "Name": "DS", "Publisher": "pub", "Price": 10.0, "PlayerNumber": 1, "Tags": [ { "Name": "A" }, { "Name": "B" }, { "Name": "C" }]}'
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
curl -X PATCH localhost:8080/api/boardgame/2 -H 'Content-Type: application/json' -d '{ "Name": "O", "Publisher": "pub", "Price": 10.0, "PlayerNumber": 1, "Tags": [ { "Name": "A" }]}'
```

Delete
```
curl -X DELETE localhost:8080/api/boardgame/2
```



## Tag API



Tag JSON
```
{
    "Name": "Test"
}
```



Create
```
curl -X POST localhost:8080/api/tag -H 'Content-Type: application/json' -d '{ "Name": "B" }'
```

ReadAll
```
curl -X GET localhost:8080/api/tag
```


Read
```
curl -X GET localhost:8080/api/tag/Test
```

Delete
```
curl -X DELETE localhost:8080/api/tag/Test
```



# GORM Learning Examples

Finds tags associated with a certain model
```
instance.db.Model(test).Association("Tags").Find(association)
```

Finds Boardgame with everything using Eager Loading
```
instance.db.Preload(clause.Associations).First(&bg, "name = ?", "bg")
```

Create while skipping all associations (Just creates boardgame and not tags/relations) --> Works but have to specify which Omits
```
err := instance.db.Omit("Tags.*").Create(&temp) 
```

value and not &value ->  You want to pass a pointer of the struct not the interface  
Omit() 				 -> skip the upserting of associations  
omits... 			 -> Pass each omit value as a separate argument  
Goal of omits 		 -> To receive many2many relations like 'tags.*' or 'expansions.*' and it should not create them but just add them to relational table  


# Documentation - Swagger - Swagon
For the documentation of the application, [Swag](https://github.com/swaggo/swag#the-swag-formatter) was used.
For a tutorial, see -> [Tutorial](https://martinheinz.dev/blog/9)