# go-app
Repository for a Go-app 

 
< Better describe the Architecture, UML and stuff >



# Useful links/tools
## Structure 
[Project Structure using Repository Pattern](https://dakaii.medium.com/repository-pattern-in-golang-d22d3fa76d91)

< Better describe the Structure >


## Documentation
For the documentation of the application, [Swag](https://github.com/swaggo/swag#the-swag-formatter) was used.
For a tutorial, see -> [Tutorial](https://martinheinz.dev/blog/9)

## Testing
[Framework](https://apitest.dev/)

Command to test   
Godotenv -> Initializes with .env   
./... -> Tests all directories   
DATABASE_HOST -> Overwrites host since tests are not running in docker
```
DATABASE_HOST=localhost godotenv -f .env go test ./... -v
```

For learning the testing framework, I developed Functional tests that simulate user behaviour on the catalog to trigger success and failure scenarios. 

< Better describe the Tests implemented >

< Implement and describe the Unit tests on the utils >



## Sorting
[Sorting in Golang](https://yourbasic.org/golang/how-to-sort-in-go/)

## Validation
[govalidator](https://github.com/asaskevich/govalidator)

## Possible CLI 
[Cobra](https://github.com/spf13/cobra)

## JSON 
[Decoding JSON Body](https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body)

