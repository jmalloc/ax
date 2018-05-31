--
-- This file defines the SQL schema used by SagaKeySetRepository.
--

--
-- saga_keyset contains the "key sets" that are associated with saga instances
-- that use key-based mapping.
--
CREATE TABLE IF NOT EXISTS saga_keyset (
    saga        VARBINARY(255) NOT NULL,
    mapping_key VARBINARY(255) NOT NULL,
    instance_id VARBINARY(255) NOT NULL,

    PRIMARY KEY (saga, mapping_key),
    INDEX (instance_id)
) ROW_FORMAT=COMPRESSED;
