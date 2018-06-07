--
-- ax_saga_snapshot stores snapshots saga.Data instances for eventsourced sagas.
--
CREATE TABLE IF NOT EXISTS ax_saga_snapshot (
    instance_id  VARBINARY(255) NOT NULL,
    revision     BIGINT UNSIGNED NOT NULL,
    description  VARBINARY(255) NOT NULL,
    content_type VARBINARY(255) NOT NULL,
    time         VARBINARY(255) NOT NULL,
    data         BLOB NOT NULL,

    PRIMARY KEY (instance_id, revision)
) ROW_FORMAT=COMPRESSED;
