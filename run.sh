# Closes running services
docker-compose down


# GO Path is on the PATH environment variable -> Required for swag init
export PATH=$(go env GOPATH)/bin:$PATH

# Generates Swag files
( cd catalog ; swag init --parseDependency --parseInternal )

# Runs everything
docker-compose up -d --build

# Pretty Console
# printf "\033c"
echo "running on 8080"

#Docker logging
docker logs --follow catalog