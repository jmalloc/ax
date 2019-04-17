--
-- ax_delayed_message stores the messages that are not yet ready to be sent.
--
CREATE TABLE IF NOT EXISTS ax_delayed_message (
    message_id     VARBINARY(255) NOT NULL,
    causation_id   VARBINARY(255) NOT NULL,
    correlation_id VARBINARY(255) NOT NULL,
    created_at     VARBINARY(255) NOT NULL,
    send_at        VARBINARY(255) NOT NULL,
    content_type   VARBINARY(255) NOT NULL,
    data           LONGBLOB NOT NULL,
    operation      INTEGER NOT NULL,
    destination    VARBINARY(255) NOT NULL,

    PRIMARY KEY (message_id),
    INDEX (send_at),
    INDEX (causation_id),
    INDEX (correlation_id)
) ROW_FORMAT=COMPRESSED;
