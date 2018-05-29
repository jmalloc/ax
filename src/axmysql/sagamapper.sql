--
-- This file defines the SQL schema used by SagaMapper.
--

--
-- saga_map contains the "mapping key sets" used to map incoming messages to
-- saga instances. It is used for both CRUD and eventsourced saga instances.
--
CREATE TABLE IF NOT EXISTS saga_map (
    saga          VARBINARY(255) NOT NULL,
    mapping_key   VARBINARY(255) NOT NULL,
    instance_id   VARBINARY(255) NOT NULL,

    PRIMARY KEY (saga, mapping_key),
    INDEX (instance_id)
) ROW_FORMAT=COMPRESSED;
