--
-- This file defines the SQL schema used by SnapshotRepository.
--

--
-- saga_snapshot stores snapshots saga.Data instances for eventsourced sagas.
--
CREATE TABLE IF NOT EXISTS saga_snapshot (
    instance_id  VARBINARY(255) NOT NULL,
    revision     BIGINT UNSIGNED NOT NULL,
    description  VARBINARY(255) NOT NULL,
    content_type VARBINARY(255) NOT NULL,
    data         BLOB NOT NULL,

    PRIMARY KEY (instance_id, revision)
) ROW_FORMAT=COMPRESSED;
