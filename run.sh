# Closes running services
docker-compose down


# GO Path is on the PATH environment variable -> Required for swag init
export PATH=$(go env GOPATH)/bin:$PATH

# Generates Swag files
#( cd catalog ; swag init --parseDependency --parseInternal )
( cd marketplace ; swag init --parseDependency --parseInternal )

# Runs everything
docker-compose up -d --build marketplace marketplace-db

# Wait for it to boot before testing
sleep 2 

# Test Everything
# ( cd catalog ; DATABASE_HOST=localhost godotenv -f .env go test ./... )

# Sleep for checking tests/errors
#sleep 4

# Pretty Console
printf "\033c"
echo "running on 8080"

#Docker logging
docker logs --follow marketplace