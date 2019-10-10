--
-- ax_projection_offset stores the next offset to be read by a projection consumer.
--
CREATE TABLE IF NOT EXISTS ax_projection_offset (
    persistence_key VARBINARY(255) NOT NULL,
    next_offset     BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (persistence_key)
) ROW_FORMAT=COMPRESSED;
