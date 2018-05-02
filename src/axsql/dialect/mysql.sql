CREATE TABLE outbox (
    `message_id`     VARBINARY(255) PRIMARY KEY,
);

CREATE TABLE outbox_envelope (
    `message_id`     VARBINARY(255) NOT NULL PRIMARY KEY,
    `causation_id`   VARBINARY(255) NOT NULL, -- outbox.message_id
    `correlation_id` VARBINARY(255) NOT NULL
    `time`           TIMESTAMP(6) NOT NULL,
    `content_type`   VARBINARY(255) NOT NULL,
    `body`           BLOB NOT NULL,
    `operation`      INTEGER NOT NULL,
    `destination`    VARBINARY(255) NOT NULL,

    INDEX(`causation_id`)
);
