# open a mysql client for the banking example database
docker run \
    --rm \
    --interactive \
    --tty \
    --network banking_backend \
    mariadb:10.3 \
    mysql -h db -u root --password=root --database=banking "$@"
