# Closes running services
# docker-compose down
docker-compose rm -s -v catalog-db 
docker-compose rm -s -v catalog 


# Generates Swag files
export PATH=$(go env GOPATH)/bin:$PATH # GO Path is on the PATH env -> Required for swag init
( cd catalog ; swag init --parseDependency --parseInternal )

# Runs everything
docker-compose up -d --build catalog catalog-db
#docker-compose up -d --build marketplace marketplace-db


# Test Everything
#sleep 2  # Wait for it to boot before testing
# ( cd catalog ; DATABASE_HOST=localhost godotenv -f environment/dev/.env go test ./... )


# Pretty Console
#sleep 4 # Sleep for checking tests/errors
#printf "\033c"


#Docker logging
#docker logs --follow marketplace