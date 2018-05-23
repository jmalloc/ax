--
-- This file contains SQL schema as used by the MessageStore and MessageStream
-- implementations.
--

--
-- messagestore_global stores the next global offset for the entire message store.
--
CREATE TABLE IF NOT EXISTS messagestore_offset (
    _    ENUM('') NOT NULL PRIMARY KEY DEFAULT '', -- ensure there can be only one row
    next BIGINT UNSIGNED NOT NULL DEFAULT 0
);

--
-- messagestore_stream identifies the streams that exist within the message store.
--
CREATE TABLE IF NOT EXISTS messagestore_stream (
    stream_id BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name      VARBINARY(255) NOT NULL UNIQUE,
    next      BIGINT UNSIGNED NOT NULL
) ROW_FORMAT=COMPRESSED;

--
-- messagestore_message contains the messages on each stream.
--
CREATE TABLE IF NOT EXISTS messagestore_message (
    global_offset  BIGINT UNSIGNED NOT NULL,
    insert_time    TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),

    stream_id 	   BIGINT UNSIGNED NOT NULL,
    stream_offset  BIGINT UNSIGNED NOT NULL,
    description    VARBINARY(255) NOT NULL,
    message_id     VARBINARY(255) NOT NULL,
    causation_id   VARBINARY(255) NOT NULL,
    correlation_id VARBINARY(255) NOT NULL,
    time           VARBINARY(255) NOT NULL,
    content_type   VARBINARY(255) NOT NULL,
    data           BLOB NOT NULL,

    PRIMARY KEY (global_offset, insert_time),
    INDEX (stream_id, stream_offset)
) ROW_FORMAT=COMPRESSED;
