--
-- projection_offset stores the next offset to be read by a projection consumer.
--
CREATE TABLE IF NOT EXISTS projection_offset (
    projection  VARBINARY(255) NOT NULL,
    next_offset BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (projection)
) ROW_FORMAT=COMPRESSED;
