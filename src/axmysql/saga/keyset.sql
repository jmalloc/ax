--
-- ax_saga_keyset contains the "key sets" that are associated with saga instances
-- that use keyset.Mapper.
--
CREATE TABLE IF NOT EXISTS ax_saga_keyset (
    persistence_key VARBINARY(255) NOT NULL,
    mapping_key     VARBINARY(255) NOT NULL,
    instance_id     VARBINARY(255) NOT NULL,

    PRIMARY KEY (persistence_key, mapping_key),
    INDEX (instance_id)
) ROW_FORMAT=COMPRESSED;
