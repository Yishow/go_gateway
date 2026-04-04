-- -----------------------------------------------------------------------------
-- 來源規則表新增目標資料型態與縮放欄位
-- -----------------------------------------------------------------------------

-- 新增 target_data_type 欄位（可空，預設 NULL 表示與 Point 讀取型別相同）
ALTER TABLE source_rules ADD COLUMN target_data_type TEXT
    CHECK (target_data_type IS NULL OR target_data_type IN ('bool', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'string'));

-- 新增 scale_multiplier 欄位（縮放倍率，可空）
ALTER TABLE source_rules ADD COLUMN scale_multiplier REAL;

-- 新增 scale_offset 欄位（偏移量，可空）
ALTER TABLE source_rules ADD COLUMN scale_offset REAL;

-- 註解：target_data_type 為 NULL 表示與 data_type 相同
-- scale_multiplier/scale_offset 為 NULL 表示無縮放