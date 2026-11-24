# Go-app
Repository for a Go-app 

The following project intends to be a learning example of a marketplace for boardgame related products. 

 
# Architecture 

<!-- TODO -->

Microservices:
- Catalog - Restricted service that contains a pre-defined set of products (Boardgames only for now). Main objective is to impose other components a more restrictive choice of products to use. Restricts the objects in the architecture.

<!-- TODO -->
- Marketplace - Basically a OLX for Boardgames where users can create offers. This is where users can re-sell products to each other.
- User Management - Authentication and Authorization service that uses Oauth2. All other services are locked behind the tokens generated in this service. This service also has the user information to be used in other services (E.g: Ratings have a reference to a username which is present due to the authz token)
- Comments - Service that contains comments left by users. These comments can reference several types (Products, Lists, Comments, etc) and therefore, aims to abstract and decouple the connections to other objects. 
- Ratings - Service that contains Ratings left by users. These ratings can reference several types (Products, Lists, Comments, etc) and therefore, aims to abstract and decouple the connections to other objects. 
- Lists - Service that aims to hold all product related lists of the user. Given that the Catalog is abstracted to the notion of a user, this svc aims to encapsulate all that logic. This allows the user to still create wishlist, ranking lists, etc..
- Notifications - Service that aims to notify the user of something (Marketplace offer for example)

Note: 
Comments and Ratings require two fields namely reference_id and reference_namespace. These fields basically refer to what object we are refering too and what Unique ID it has on the architecture.

## Microservices communication
The communications between the FE and the several services can be seen as domain-related.
If we require to perform boardgame related actions, all requests should pass through the Catalog since we require business logic related to it (E.g: BG exists so we can rate it?)


## Documentation
For the documentation of the application, [Swag](https://github.com/swaggo/swag#the-swag-formatter) was used.

## Testing
* [Integration tests framework](https://apitest.dev/)
* [mock-gen](https://github.com/uber-go/mock)
* [embedded-postgres](https://github.com/fergusstrange/embedded-postgres)