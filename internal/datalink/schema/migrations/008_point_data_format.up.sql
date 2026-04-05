-- -----------------------------------------------------------------------------
-- 點位表新增字節序（可空；空字串／NULL 表示歷史 Modbus 預設 ABCD）
-- -----------------------------------------------------------------------------

ALTER TABLE points ADD COLUMN data_format TEXT
    CHECK (
        data_format IS NULL
        OR data_format = ''
        OR data_format IN ('ABCD', 'BADC', 'CDAB', 'DCBA')
    );
