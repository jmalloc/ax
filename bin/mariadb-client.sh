docker run \
    --rm \
    --interactive \
    --tty \
    --network ax_mariadb \
    mariadb:10.3 \
    mysql -h ax_mariadb -u root --password=root --database=ax "$@"
