docker run \
    --rm \
    --detach \
    --name ax_rmq \
    --publish 127.0.0.1:5672:5672 \
    --publish 127.0.0.1:15672:15672 \
    rabbitmq:management
