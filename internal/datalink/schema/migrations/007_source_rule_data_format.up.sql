-- -----------------------------------------------------------------------------
-- 來源規則表新增字節序格式（可空；空字串／NULL 表示使用連線預設）
-- -----------------------------------------------------------------------------

ALTER TABLE source_rules ADD COLUMN data_format TEXT
    CHECK (
        data_format IS NULL
        OR data_format = ''
        OR data_format IN ('ABCD', 'BADC', 'CDAB', 'DCBA')
    );
