-- Copyright 2020 The Go Authors. All rights reserved.
-- Use of this source code is governed by a BSD-style
-- license that can be found in the LICENSE file.

BEGIN;

DROP INDEX idx_search_documents_version_updated_at;

ALTER TABLE search_documents
	ALTER COLUMN created_at type TIMESTAMP WITHOUT TIME ZONE,
	ALTER COLUMN updated_at type TIMESTAMP WITHOUT TIME ZONE,
	ALTER COLUMN version_updated_at type TIMESTAMP WITHOUT TIME ZONE,
	ALTER COLUMN imported_by_count_updated_at type TIMESTAMP WITHOUT TIME ZONE;

CREATE INDEX idx_search_documents_version_updated_at ON search_documents USING btree (version_updated_at);

COMMENT ON INDEX idx_search_documents_version_updated_at IS
'INDEX idx_search_documents_version_updated_at is used for incremental update of imported_by counts, in order to determine when the latest version of a package was last updated.';


END;