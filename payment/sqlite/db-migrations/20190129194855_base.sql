-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE payments (
  id TEXT CONSTRAINT ct__payments_id__uuid CHECK (length(id) == 36),
  version INTEGER DEFAULT 0
    CONSTRAINT ct__payments_version__not_null NOT NULL
    CONSTRAINT ct__payments_version__gte_zero CHECK (version >= 0),
  organisation_id TEXT
    CONSTRAINT ct__payments_organisation_id__not_null NOT NULL
    CONSTRAINT ct__payments_organisation_id__uuid CHECK (length(organisation_id) == 36),
  data TEXT
    CONSTRAINT ct__payments_data__not_null NOT NULL
    CONSTRAINT ct__payments_data__json_valid CHECK (length(json(data)) > 1),
  CONSTRAINT uq__payments_id UNIQUE (id, version)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

-- Rollback migrations are not used for 3 main reasons:
--
-- * In a continuous integration and deployment environment, it simplifies the
--   pipeline, just always executing the same procedure in new changes, hence
--   when a migration must be rolled back, the procedure is just the same than
--   when committing one, which is create a new migration for doing the changes
--   which in this case those changes will be the ones required to rollback the
--   ones introduced by a previous migration.
-- * Some times the changes in the schema are aligned with changes in sources,
--   when that happens, the rollback of the migration implies to revert changes
--   in the sources, so it's simpler to revert the changes in the sources and
--   adding a new migration which rollback the schema changes and be applied
--   though the continuous integration as any other change.
-- * The schema is straightforward to be rolled back, but if the changes in the
--   schema has been used in production then it's difficult to predict if what
--   is desired is to just do the opposite changes in the schema without
--   considering data which has been inserted or updated on those new parts of
--   of the schema.
