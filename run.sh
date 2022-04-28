docker-compose down
docker-compose up -d --build

printf "\033c"
echo "running on 8080"
docker logs --follow boardgame-app