--
-- saga_instance stores saga.Data instances for each instance of a CRUD saga.
--
CREATE TABLE IF NOT EXISTS projection_offset (
    projection  VARBINARY(255) NOT NULL,
    next_offset BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (projection)
) ROW_FORMAT=COMPRESSED;
