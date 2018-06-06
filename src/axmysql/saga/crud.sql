--
-- ax_saga_instance stores saga.Data instances for each instance of a CRUD saga.
--
CREATE TABLE IF NOT EXISTS ax_saga_instance (
    instance_id  VARBINARY(255) NOT NULL,
    revision     BIGINT UNSIGNED NOT NULL,
    description  VARBINARY(255) NOT NULL,
    content_type VARBINARY(255) NOT NULL,
    data         BLOB NOT NULL,

    PRIMARY KEY (instance_id)
) ROW_FORMAT=COMPRESSED;
