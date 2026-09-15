// Package datalink 提供設備資料收集管線的核心功能。
//
// 本套件實作以下子模組：
//   - device: 設備註冊與管理
//   - point: 點位目錄與配置
//   - tag: 全域標籤字典
//   - mapping: 映射管線與轉換
//   - collector: 資料收集排程器
//   - storage: 時序資料儲存
//
// # 標籤命名規則
//
// 標籤鍵 (Tag Key) 遵循以下規則：
//   - 允許的字元：ASCII 字母 (a-z, A-Z)、數字 (0-9)、底線 (_)、連字號 (-)、點號 (.)、斜線 (/)
//   - 長度限制：1-128 個字元
//   - 唯一性：大小寫不敏感 (case-insensitive)
//   - 範例有效鍵：site1/area_a/temp-sensor, PLC.D0100, motor_speed_rpm
//
// # 時間精度 (Timestamp Precision)
//
// 系統支援以下時間精度設定：
//   - TimePrecisionSecond: 秒級精度，時戳截斷至秒
//   - TimePrecisionMillisecond: 毫秒級精度，時戳保留至毫秒
//
// 精度設定為全域配置，透過 API 或 UI 進行調整。
//
// # 轉換表達式語法 (Transform Expression Syntax)
//
// 轉換管線支援以下表達式語法：
//
// ## 運算子 (Operators)
//   - 算術運算：+, -, *, /, % (取模)
//   - 比較運算：>, <, >=, <=, ==, !=
//   - 布林運算：and, or, not (使用英文關鍵字)
//
// ## 函數集 (Functions)
//
// 條件函數：
//   - if(condition, true_value, false_value)
//
// 數學函數：
//   - abs(x): 絕對值
//   - min(a, b, ...): 最小值
//   - max(a, b, ...): 最大值
//   - clamp(value, min, max): 限制值在範圍內
//   - round(x): 四捨五入
//   - floor(x): 向下取整
//   - ceil(x): 向上取整
//   - sqrt(x): 平方根
//   - pow(base, exp): 次方
//
// 位元運算函數：
//   - bitand(a, b): 位元 AND
//   - bitor(a, b): 位元 OR
//   - bitxor(a, b): 位元 XOR
//   - shiftl(value, bits): 左移
//   - shiftr(value, bits): 右移
//
// 型別轉換函數：
//   - int(x): 轉換為整數
//   - float(x): 轉換為浮點數
//   - bool(x): 轉換為布林值
//   - string(x): 轉換為字串
//   - coalesce(a, b, ...): 返回第一個非空值
//
// ## 語法規則
//   - 函數名稱使用 snake_case 格式
//   - 運算子使用中綴表示法 (infix notation)
//   - 支援括號改變優先序
//   - 變數參考使用 $raw (原始值) 或 $prev (上一步結果)
//
// ## 範例表達式
//   - 簡單縮放： $raw * 0.1 + 10
//   - 條件判斷： if($raw > 100, 100, $raw)
//   - 位元擷取： bitand(shiftr($raw, 8), 0xFF)
//   - 型別轉換： round(float($raw) / 10.0)
package datalink
