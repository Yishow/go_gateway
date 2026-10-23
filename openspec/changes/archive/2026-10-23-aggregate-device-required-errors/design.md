## Context

基準 f21e2a6b；原12個未提交檔已精確複製到既有 review clone，原工作樹暫不修改。此 change 僅負責 device validation 與相關說明，Windows logs 由既有 active change 負責。

## Goals / Non-Goals

**Goals:** 同時呈現 missing／empty required，保留單類錯誤格式及 schema 條件。準確提交既有文件意圖。

**Non-Goals:** 不重寫 schema validator、不改協定、driver、HTTP envelope、optional欄位或型別政策；不把 cross-build 當 Windows 驗收，不執行現場／正式 DB。

## Decisions

### 保留既有驗證順序並合併必填問題

missing 在前、empty 在後，各依 schema 順序；只有其中一類時保留原訊息。兩類同時存在以分隔段落表示。range／enum 錯誤仍在 required 問題修正後回報。替代第一類優先 return 仍需重送，故拒絕。

## Implementation Contract

validateConnectionConfig 的 error 文字不得包含缺欄位的設定值，只列欄名與缺少／空值原因。Modbus {host:"",slave_id:1} 同時顯示 port missing 與 host empty；FATEK {mode:""} 顯示 station_no missing 與 mode empty。保留單類 message、合法設定、nil map、條件 serial/network 及既有 range／enum checks。先新增失敗混合測試，再最小修正；原草稿的 aggregation 保護以 mutation／完整測試核對。文件去除私有現場位址、group/資料表執行識別，保留有來源日期的歷史摘要與本輪未重驗註記。

## Risks / Trade-offs

- 呼叫端若依文字判斷錯誤 → 單類格式不變，HTTP envelope 不變，僅混合案例擴充內容。
- 歷史現場成功不能本輪重驗 → 明示來源與限制；durability 僅描述 accepted journal 與可證明 scope，unknown 保持 blocked。
