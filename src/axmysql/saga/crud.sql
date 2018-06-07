--
-- ax_saga_instance stores saga.Data instances for each instance of a CRUD saga.
--
CREATE TABLE IF NOT EXISTS ax_saga_instance (
    instance_id     VARBINARY(255) NOT NULL,
    revision        BIGINT UNSIGNED NOT NULL,
    persistence_key VARBINARY(255) NOT NULL,
    description     VARBINARY(255) NOT NULL,
    content_type    VARBINARY(255) NOT NULL,
    data            BLOB NOT NULL,
    insert_time     TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6),
    update_time     TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),

    PRIMARY KEY (instance_id)
) ROW_FORMAT=COMPRESSED;
