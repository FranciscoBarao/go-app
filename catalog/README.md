# Catalog
Readme for the catalog service


The following service contains the different boardgame related products that can be browsed. These products are the base for the marketplace.
At the moment, we have:
- Boardgames
- Expansions

Each product can be classified with:
- Tags
- Mechanisms
- Categories


## Boardgame API
Boardgame JSON
```
{
    "Name": "name",
	"Publisher": "publisher",
	"Price": 10.0,
	"PlayerNumber": 1,
	"Tags": [
		{ "Name": "A" }
	],
    "Categories": [
		{ "Name": "A" }
	],
    "Mechanisms": [
		{ "Name": "A" }
	]
}
```


Create
```
curl -X POST localhost:8080/api/boardgame -H 'Content-Type: application/json' -d '{ "Name": "DS", "Publisher": "pub", "Price": 10.0, "PlayerNumber": 1, "Tags": [], "Categories: [], "Mechanisms: []}'
```

ReadAll
```
curl -X GET localhost:8080/api/boardgame
```
ReadAll can be filtered and sorted. The filters can have 2 formats, depending on what is being evaluated.
```
filterBy -> Field.Value 
filterBy -> Field.Operator.Value 	

sortBy -> Field.Order
```

Examples of filters that work:
```
	name.a 		   --->   name LIKE ?    %a%
	price.le.10    --->   price <= ?     10
```

Filters will require an update sometime in the future because it doesnt allow floats cause we can't do ```price.lt.10,4```. 

Examples of sorts that work:
```
	name.asc 	  --->    ordered by name in alphabetical ascending order
	price.desc    --->    ordered by price in numerical descending order
```


Read
```
curl -X GET localhost:8080/api/boardgame/<id>
```

Update
```
curl -X PATCH localhost:8080/api/boardgame/<id> -H 'Content-Type: application/json' -d '{ "Name": "O", "Publisher": "pub", "Price": 10.0, "PlayerNumber": 1, "Tags": [], "Categories: [], "Mechanisms: []}'
```

Delete
```
curl -X DELETE localhost:8080/api/boardgame/<id>
```



## Tag/Mechanism/Catagory API

The following three many2many relations all consist of a unique string. These fields are not created in Upscale, which means that when a boardgame is being created, if these fields are added, they must previously exist or the BG creation will fail. The following endpoint description is similar to all three and just vary on the url endpoint possibly being:
```
/tag/
/category/
/mechanism/
```


JSON
```
{
    "Name": "name"
}
```

Create
```
curl -X POST localhost:8080/api/tag -H 'Content-Type: application/json' -d '{ "Name": "name" }'
```

ReadAll
```
curl -X GET localhost:8080/api/tag
```

ReadAll can be sorted. 
```
sortBy -> Field.Order
```

Examples of sorts that work:
```
	name.asc 	  --->    ordered by name in alphabetical ascending order
	price.desc    --->    ordered by price in numerical descending order
```


Read
```
curl -X GET localhost:8080/api/tag/<name>
```

Delete
```
curl -X DELETE localhost:8080/api/tag/<name>
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



## Difference between many2many on 1 table or on both
**Description:** Had Tags[] in boardgame like ```Tags []Tag `gorm:"many2many:boardgame_tags;"` ``` and did not have a list of BGs on the Tags. (only on 1 Table).  
Having this one-sided allowed me to add/delete associations of Tags (with or without Upserting).  
**Issue:** Attempting to delete a Tag that had an association to a BG would fail due to FK constraint.

**Solutions:**
- A) Get all BGs that have that tag and 1 by 1 delete the association
- B) Add BGs list to Tags like ``` Boardgames []Boardgame `gorm:"many2many:boardgame_tags;" json:"-"` ```

**Choice: B** Previously, handling associations was not bidirectional, which means that I was able to handle Tags via BGs but not the other way around.  
**Improvements:** Ability to delete Tags that are already associated. Improved way of returning all BGs with a specific Tag. 