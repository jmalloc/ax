package axmysql

// SagaSchema is a collection of DDL queries that create the schema
// used by OutboxRepository.
var SagaSchema = []string{
	`CREATE TABLE IF NOT EXISTS saga_instance (
		instance_id           VARBINARY(255) NOT NULL PRIMARY KEY,
		current_revision      BIGINT UNSIGNED NOT NULL,

		snapshot_revision     BIGINT UNSIGNED,
		snapshot_description  VARBINARY(255),
		snapshot_content_type VARBINARY(255),
		snapshot_data         BLOB,

		INDEX (saga)
	)`,
	`CREATE TABLE IF NOT EXISTS saga_keyset (
		saga          VARBINARY(255) NOT NULL,
		mapping_key   VARBINARY(255) NOT NULL,
		instance_id   VARBINARY(255) NOT NULL,

		PRIMARY KEY (saga, mapping_key),
		INDEX (instance_id)
	)`,
}
