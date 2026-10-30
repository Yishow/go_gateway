-- Nullable expand only: never infer per-member intent from legacy reductions.
ALTER TABLE write_groups ADD COLUMN result_shape TEXT
    CHECK (result_shape IS NULL OR result_shape = 'member-calculation-v1');
ALTER TABLE write_group_members ADD COLUMN calculation TEXT
    CHECK (calculation IS NULL OR calculation IN ('latest', 'mean', 'delta'));
