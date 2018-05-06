docker rm --force ax_mariadb

docker network create ax_mariadb

docker run \
    --rm \
    --detach \
    --name ax_mariadb \
    --network ax_mariadb \
    --publish 127.0.0.1:3306:3306 \
    --env MYSQL_ROOT_PASSWORD=root \
    --env MYSQL_DATABASE=ax \
    --env MYSQL_USER=ax \
    --env MYSQL_PASSWORD=ax \
    mariadb:10.3
