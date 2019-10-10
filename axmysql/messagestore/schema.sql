--
-- ax_messagestore_offset stores the next global message offset across all streams.
--
-- It uses an ENUM field with a single value as the primary key to
-- ensure there can only ever be a single row.
--
CREATE TABLE IF NOT EXISTS ax_messagestore_offset (
    _    ENUM('') NOT NULL PRIMARY KEY DEFAULT '',
    next BIGINT UNSIGNED NOT NULL DEFAULT 0
);

--
-- ax_messagestore_stream contains the streams that exist within the message store.
--
-- next is the next unused offset on the stream.
--
CREATE TABLE IF NOT EXISTS ax_messagestore_stream (
    stream_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name      VARBINARY(255) NOT NULL UNIQUE,
    next      BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (stream_id)
) ROW_FORMAT=COMPRESSED;

--
-- ax_messagestore_message contains the messages on each stream.
--
-- The primary key includes the insert_time to allow for time-based
-- partitioning.
--
CREATE TABLE IF NOT EXISTS ax_messagestore_message (
    global_offset  BIGINT UNSIGNED NOT NULL,
    insert_time    TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),

    stream_id      BIGINT UNSIGNED NOT NULL,
    stream_offset  BIGINT UNSIGNED NOT NULL,
    description    VARBINARY(255) NOT NULL,
    message_id     VARBINARY(255) NOT NULL,
    causation_id   VARBINARY(255) NOT NULL,
    correlation_id VARBINARY(255) NOT NULL,
    created_at     VARBINARY(255) NOT NULL,
    send_at        VARBINARY(255) NOT NULL,
    content_type   VARBINARY(255) NOT NULL,
    data           LONGBLOB NOT NULL,

    PRIMARY KEY (global_offset, insert_time),
    INDEX (stream_id, stream_offset),
    INDEX (message_id),
    INDEX (causation_id),
    INDEX (correlation_id)
) ROW_FORMAT=COMPRESSED;
