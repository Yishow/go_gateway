## Context

`fix-studio-v2-setup-and-db-flow` 的實作已完成但尚未提交（工作區差異約 3,150 行，涵蓋 `internal/datalink/dbtarget`、`internal/datalink/device`、`internal/api` 與 Studio V2 前端）。程式碼審查在該差異上確認 12 項缺陷，橫跨三個子系統：

- **Studio V2 前端設定 / 啟用流程**：設定頁存檔狀態機新增了 ownership（世代）機制與 Modbus Share 存檔狀態欄位，但狀態進入 `saving` 的入口比離開 `saving` 的出口多，且 autosave barrier 由同步讀取改為非同步等待後失去了「一定會有結果」的性質。
- **MySQL 資料庫目標**：本次差異首度讓 MySQL 成為 `GenerateSchema` 的支援型別（新增 `internal/datalink/dbtarget/service_mysql_inspection.go`），使原本 MySQL 走不到的 schema 產生路徑變成可達，暴露出方言語法與預設 schema 名稱兩處不相容；同時 MySQL DSN 組裝加入了 cleartext 認證與明文退回的預設值。
- **Modbus Share 啟用前置檢查**：barrier validator 的檢查順序被調整，Share 停用時由「回報停用錯誤」改為「直接放行」。

限制條件：

- 差異尚未提交，修正應與原差異一併進入主線，不另外拆分為多個提交序列。
- Go 端已有針對本次差異新增的測試（例如 Share 停用仍可啟用工作區、MySQL DSN 預設值），修正必須明確說明哪些既有測試預期需要一併更新，而非讓測試「碰巧」失敗。
- 前端測試使用 Vitest + Testing Library，逾時類行為必須可由測試控制，不得依賴實際等待。

## Goals / Non-Goals

**Goals:**

- 讓 Studio V2 的存檔狀態機在所有路徑上收斂到終態（`saved`、`save-error` 或 `idle`），使 Step 4 啟用永遠會得到成功或可重試的失敗。
- 讓 MySQL 資料庫目標的 schema 產生輸出方言合法、schema 名稱正確的語句，並保證 upsert 映射所需的唯一鍵確實存在。
- 讓 MySQL 連線在無法建立加密通道時失敗關閉，而非以明文送出密碼。
- 讓 Share 停用時的啟用前置檢查範圍成為明確的契約，而非早退所造成的副作用。

**Non-Goals:**

- 不重寫設定頁的 ownership / 佇列機制，也不改變既有 `SettingsOperationOwnership` 的 API 形狀；本次僅修正狀態收斂與失效範圍。
- 不新增前端全域的存檔狀態機抽象層，也不將 autosave barrier 改為 React Query 或其他狀態庫管理。
- 不為 MySQL 以外的資料庫類型調整 schema 產生語法或連線安全預設值。
- 不處理「使用者從未執行 schema 產生就直接寫入 MySQL」情境下的執行期唯一鍵偵測；本次僅保證 schema 產生路徑會建立唯一鍵。
- 不變更 Share 啟用的錯誤碼集合（沿用 `modbus_share_save_incomplete`、`modbus_share_revision_conflict`、`modbus_share_hydration_required`）。
- 不調整前端既有的協議預設值表（`protocols.ts` 的預設值內容維持不變），僅改變這些預設值何時被寫入裝置設定。

## Decisions

### D1 存檔中狀態僅由實際發出的存檔請求設定

`onUpdateModbusShare` 與 `onReset` 目前在每次欄位編輯時就把 `settings.modbus_share.save_state` 設為 `saving`，但這兩個進入點都不會發出任何請求，唯一能清除該狀態的 `persistSettingsSnapshot` 只有在使用者按下儲存後才執行。結果是任何一次 Modbus Share 編輯都會讓狀態永久停在 `saving`。

決策：欄位編輯與重置只更新資料本身與清除既有 `save_error`，不得寫入 `save_state: 'saving'`；`saving` 僅由 `onSave` 在實際排入存檔佇列時設定。編輯後的「有未儲存變更」語意由既有的儲存列（save bar）呈現，不借用存檔狀態欄位。

替代方案：讓編輯時寫入 `dirty` 之類的第三種狀態。否決原因為 `save_state` 同時被 autosave barrier 消費，新增狀態值會擴大 barrier 的判斷面，且 barrier 對未知狀態的行為需另行定義，成本高於收益。

### D2 autosave settlement 等待改為有界

`useStudioV2AutosaveSettlement` 回傳的等待函式在偵測到有存檔進行中時建立一個只能由「下一次無存檔進行中的狀態變更」解決的承諾，沒有逾時、不會被拒絕、也不會在卸載時清空等待者。任何殘留的 `saving` 都會讓 `activateWorkspace` 永久等待，而呼叫端 `useStep4Activation.start` 已先把畫面切到 `activating`。

決策：等待加入固定上限（10 秒，以具名常數表示並可由測試注入覆寫）。逾時後以「當下狀態的 barrier 快照」解決該承諾，而非拒絕；由於此時仍有 `saving`，快照的 `pending_saves` 必然大於 0，呼叫端會走既有的 `modbus_share_save_incomplete` 失敗路徑，畫面顯示可重試的錯誤。等待者在解決或逾時後一律從等待清單移除，元件卸載時清除計時器。

替代方案一：逾時後 reject。否決原因為呼叫端的 catch 會把它歸到 `activation_failed` 泛用碼，喪失「存檔未完成」這個可操作的語意。
替代方案二：恢復先前的同步讀取。否決原因為同步讀取會讓「存檔正在進行中但很快就會完成」的正常情況誤報失敗，這正是本次差異改為等待的動機。

### D3 存檔擁有權失效不得吞掉終態

`invalidateSettingsSave` 會在每次一般 / 時序 / 排程 / Modbus 設定編輯時推進 `settings-save` 世代。若使用者在存檔進行中編輯任一欄位，`persistSettingsSnapshot` 的成功與失敗分支都會因 `ownership.isCurrent` 為偽而全部略過：錯誤不顯示、不提供重試，`save_state` 停在 `saving`。

決策：把「是否為最新請求」與「是否需要收斂終態」分離。ownership 檢查僅用於決定是否覆寫使用者可見的欄位資料（例如 `settings_revision` 回填）與是否設定重試動作；`save_state` 與 `save_error` 的終態寫入不受 ownership 影響，一定執行。失敗時即使世代已推進，也必須設定 `operationError` 與可重試狀態，讓使用者知道剛才那次存檔失敗。

替代方案：改成編輯欄位時不推進世代，只有真正的新存檔請求才推進。否決原因為世代機制同時保護 `settings_revision` 回填不被舊回應覆寫，取消該保護會重新引入舊回應污染新編輯的問題。

### D4 連接器成功清單以生命週期而非目前狀態修剪

`successfulAdds` 這個 ref 只以「目前 state 中存在的連接器 id」修剪。新增 A、刪除 A、再新增 B 時，第三步的既有 id 集合為空，A 通過修剪並被重新寫回列表，畫面出現一個後端已不存在的連接器。

決策：刪除連接器時明確把該 id 從 `successfulAdds` 移除（以刪除事件為準），並在寫回列表時同時排除已知的已刪除 id。修剪條件由「目前不存在於 state」改為「已被寫回列表或已被刪除」。

替代方案：改用一次性的 refetch 取代本地合併。否決原因為新增流程刻意採用樂觀寫入以避免多次新增時的閃爍，改為 refetch 會改變既有互動體感且超出本次修正範圍。

### D5 站號預設值在載入正規化階段寫入設定

`ConnectionConfigForm` 對 TCP 系協議以顯示用的預設站號（`mc_3e` 為 0、其餘為 1）補上空值，但該值從未寫回裝置設定。經由 API 建立、未帶 `station_no` 的 MC-3E 裝置在畫面上顯示已填好的「0」，`isStudioV2DeviceValid` 卻因 `station_no ?? station` 為 undefined 而判定 `draft-invalid`，使用者看不到任何空欄位可修正，啟用被擋住。

決策：站號預設值改在裝置載入正規化階段（`studioV2DeviceAutosave.ts` 既有的協議設定正規化流程）寫入設定；表單不再自行捏造顯示值，欄位值直接反映設定內容，缺值時顯示為空。正規化沿用 `protocols.ts` 既有的協議預設值，不新增第二份預設值來源。

替代方案：放寬 `isStudioV2DeviceValid` 對 `mc_3e` 站號的要求。否決原因為站號是後端 driver 的必要參數，放寬會把錯誤延後到執行期。

### D6 MySQL 唯一索引改用方言合法語法

`buildEnsureUniqueIndexStatement` 一律輸出 `CREATE UNIQUE INDEX IF NOT EXISTS`，這在 MySQL 不是合法語法。本次差異之前 MySQL 走不到此路徑（`inspectTables` 直接回報不支援的資料庫類型），差異讓 MySQL 成為支援型別後此路徑變成可達。連帶影響是 MySQL 的 `ON DUPLICATE KEY UPDATE` 只在時間戳欄位確實具備唯一鍵時才會生效，索引建不起來時 upsert 會無聲退化為每個時間桶重複寫入。

決策：依資料庫類型輸出語法——MySQL 輸出不含 `IF NOT EXISTS` 的 `CREATE UNIQUE INDEX`，SQLite 與 PostgreSQL 維持既有含 `IF NOT EXISTS` 的語法。重複建立的防護沿用既有機制：欄位已是主鍵或唯一鍵時（由資料表檢查結果的 `PrimaryKey` / `Unique` 判定）不產生語句，同一批次內以索引名稱去重。

替代方案：對 MySQL 先查 `information_schema.statistics` 再決定是否建立。否決原因為資料表檢查結果已帶有欄位的主鍵 / 唯一鍵資訊，額外查詢會讓 schema 產生從單純的語句組裝變成需要連線的流程，且 dry-run 模式將無法純離線產生語句。

### D7 schema 產生器改以連接器解析預設 schema

`buildSchemaGenerateStatements` 仍以 `defaultSchemaForKind(kind)` 補齊空白的 `table_schema`（MySQL 得到 `main`），而映射的建立與更新已改用 `defaultSchemaForConnector`（MySQL 取連線設定中的資料庫名稱）。任何 `table_schema` 為空的映射資料列，其鍵值會對不上以真實資料庫名稱建立的資料表清單，並產生指向 `main` 的建表語句，執行時失敗於未知資料庫。

決策：`buildSchemaGenerateStatements` 改為接收連接器本身而非僅接收資料庫類型，空白 schema 一律以 `defaultSchemaForConnector` 解析，與映射寫入端共用同一個預設值來源。

替代方案：在讀取映射時就把空白 schema 正規化補齊。否決原因為那會改動映射投影的語意（讓讀取端產生與資料庫中不同的值），影響面大於單一產生器的參數調整。

### D8 MySQL cleartext 認證改為失敗關閉

MySQL DSN 組裝（正式連線與探測用的管理連線兩處）同時具備三個設定：`allow_cleartext_passwords` 預設為 true、`AllowFallbackToPlaintext` 固定為 true、以及 cleartext 啟用時 TLS 模式回傳 `preferred`。驅動的 `preferred` 本身即隱含允許退回未加密連線，因此伺服器未提供 TLS 時，`mysql_clear_password` 外掛會把目標資料庫密碼以明文送上線路，且不產生任何錯誤或警告。

決策：

- `allow_cleartext_passwords` 預設值改為 false。MySQL 8 的 `caching_sha2_password` 在未加密連線上會由驅動自動走 RSA 公鑰交換，不需要 cleartext 外掛，因此此預設值變更不影響一般帳號登入。
- 連線設定明確提供 `tls` 值時，以該值為準（視為操作者的明示選擇）。
- 未提供 `tls` 但要求 `use_tls` 時，TLS 模式為 `skip-verify`（實際加密、不驗證憑證），不再使用會靜默退回明文的 `preferred`。
- 未提供 `tls` 但明確啟用 cleartext 時，TLS 模式為 `skip-verify` 且不允許退回明文；伺服器未提供 TLS 時連線失敗並回報可辨識的錯誤。
- `AllowFallbackToPlaintext` 不再無條件設為 true。

替代方案一：維持預設值但加上警告日誌。否決原因為密碼外洩是一次性且不可回溯的事件，日誌無法阻止它發生。
替代方案二：cleartext 啟用時要求 `tls=true`（完整憑證驗證）。否決原因為工控現場的 MySQL 多使用自簽憑證，強制完整驗證會讓多數合法部署無法連線，實務上會被操作者以更不安全的方式繞過。

### D9 Share 停用時保留工作區層級前置檢查

barrier validator 目前在 `settings.Enabled` 為偽時直接回傳 nil，略過的不只是 Share 專屬檢查，還包含 readiness token、workspace revision 與 settings revision 三項檢查。只要 Modbus Share 恰好停用，帶著未完成或過期 autosave 的工作區也能通過啟用。另一方面，本次差異新增的測試明確要求「Share 從未完成 hydration 時，停用狀態下工作區仍應可啟用」，這是刻意的設計意圖，不能倒退。

決策：把「Share 專屬檢查」與「工作區層級檢查」分開：

- Share 啟用時：維持現行完整檢查順序（hydration 就緒、readiness token、settings revision、workspace revision）。
- Share 停用且 hydration 已就緒時：仍檢查 readiness token 與 workspace revision，略過 settings revision 這項 Share 設定專屬檢查。
- Share 停用且 hydration 未就緒時：直接放行，保留「Share 未完成 bootstrap 不得綁架工作區啟用」的既有行為。

同時修正錯誤分類：`GetSettings` 回傳錯誤時不再回報為「Share 已停用」，而是回報為前置檢查的內部讀取失敗（可重試），並且不從零值讀取 settings revision。

替代方案：Share 停用時完全不檢查，改由其他機制保障工作區層級 barrier。否決原因為目前正式接線中並未安裝 `WithRevisionValidator`，Share 的 barrier validator 是唯一的伺服器端 revision 再驗證來源，移除後只剩前端自報的 `pending_saves`。

### D10 移除探測路徑的重複自動建庫分支

`probeConnector` 對 MySQL 先呼叫已內建「建庫後重試」的 `openExternalDBManagerWithMySQLDatabaseEnsure`，之後又有一段針對 MySQL 的 else-if 分支再次呼叫建庫函式。資料庫不存在時該分支永遠進不去（前一步已處理），其他錯誤時則以已被轉換過的錯誤再次觸發建庫嘗試。

決策：移除該 MySQL else-if 分支，PostgreSQL 分支維持不變（PostgreSQL 沒有對應的內建重試包裝）。

替代方案：把 PostgreSQL 也改為內建重試包裝以求對稱。否決原因為那是行為重構而非缺陷修正，超出本次範圍。

## Implementation Contract

**行為契約（前端）**

- 編輯 Modbus Share 任一欄位或按下重置後，`settings.modbus_share.save_state` 不會變成 `saving`；Step 4 的啟用按鈕在沒有實際存檔進行中的情況下可立即進入啟用流程。
- 按下儲存後 `save_state` 進入 `saving`，並在請求結束後必定收斂為 `saved` 或 `save-error`，即使期間使用者編輯了其他設定欄位；`save-error` 時使用者可見錯誤訊息且可重試。
- 存在殘留的 `saving` 時，啟用流程最多等待 10 秒後即以 `modbus_share_save_incomplete` 失敗，Step 4 離開 `activating` 並顯示可重試錯誤。
- 新增連接器 A、刪除 A、再新增 B 後，連接器列表只包含 B。
- 未帶站號的 MC-3E 裝置在載入後，其設定中的 `station_no` 已被補為協議預設值，畫面顯示值與 `isStudioV2DeviceValid` 的判定一致（同為有效）。

**行為契約（後端）**

- MySQL 連接器對含 upsert 映射且時間戳欄位非主鍵 / 非唯一鍵的既有資料表執行 schema 產生時，產生的語句可被 MySQL 接受並成功建立唯一索引；dry-run 模式輸出同一組語句而不需連線。
- 映射的 `table_schema` 為空時，MySQL 的 schema 產生語句以連接器連線設定中的資料庫名稱為 schema，不再出現 `main`。
- MySQL 連接器在未提供 `allow_cleartext_passwords` 時，組出的 DSN 不啟用 cleartext 認證；明確啟用 cleartext 時，DSN 的 TLS 模式為加密模式且不允許退回明文。
- Modbus Share 停用且 hydration 已就緒時，帶著錯誤或缺漏 readiness token 的啟用請求被拒絕；hydration 未就緒時啟用請求通過。
- Share 設定讀取失敗時，啟用回應的錯誤碼不是「已停用」。

**介面與資料形狀**

- `buildSchemaGenerateStatements` 的參數由資料庫類型改為連接器（型別 `*schema.DatabaseConnector`），內部改以 `defaultSchemaForConnector` 解析空白 schema；呼叫端同步調整。
- `buildEnsureUniqueIndexStatement` 依資料庫類型輸出語句，回傳值形狀（語句字串與索引鍵）不變。
- `mysqlTLSConfigValue` 的判斷改依 D8 的四條規則；MySQL DSN 組裝的兩處（正式連線與探測管理連線）共用同一組規則，不得各自複製一份。
- `useStudioV2AutosaveSettlement` 的回傳型別不變（仍回傳取得 settlement 的函式）；新增可注入的逾時參數供測試使用，預設值為具名常數。
- barrier validator 的簽章與錯誤型別不變，僅調整檢查順序與條件。

**失敗模式**

- autosave 逾時：以 barrier 快照收斂，呼叫端得到 `modbus_share_save_incomplete`，可重試。
- MySQL 唯一索引建立失敗（例如既有資料已有重複時間戳）：schema 產生整體回報失敗並保留原始資料庫錯誤訊息，屬刻意浮現而非靜默。
- MySQL 啟用 cleartext 但伺服器無 TLS：連線失敗，錯誤沿用既有的連線錯誤分類與訊息路徑，不新增錯誤碼。
- Share 設定讀取失敗：回報可重試的內部前置檢查失敗。

**驗收條件**

- 前端：`npx vitest run` 於 `frontend` 目錄下全部通過，且新增測試覆蓋——編輯 Modbus Share 後 `save_state` 不為 `saving`；存檔進行中編輯欄位後失敗仍收斂為 `save-error` 且顯示錯誤；殘留 `saving` 時啟用在注入的短逾時後得到 `modbus_share_save_incomplete`；新增 / 刪除 / 再新增後列表只含最後一筆；未帶站號的 MC-3E 裝置載入後為有效。
- 後端：`go build ./...`、`go vet ./...` 通過；`go test ./internal/datalink/dbtarget/... ./internal/api/...` 通過，且新增測試覆蓋——MySQL 唯一索引語句不含 `IF NOT EXISTS`；空白 schema 的 MySQL 映射產生的語句使用連線設定的資料庫名稱；MySQL 預設 DSN 不啟用 cleartext；明確啟用 cleartext 的 DSN 為加密且不退回明文；Share 停用但 hydration 就緒時錯誤 readiness token 被拒。
- 既有測試更新：`internal/datalink/dbtarget/service_mysql_test.go` 中斷言預設啟用 cleartext 與 `preferred` 的案例需依 D8 改寫；`internal/api/router_modbus_share_test.go` 中 Share 停用的案例需明確保持 hydration 未就緒，以對應 D9 的分支。
- 文件：依 repo 規則，若本次改動觸及 studio surface inventory 目錄則需同步 changelog；本次預期不觸及該目錄，若實作中發現觸及則必須補上。

**範圍邊界**

- 在範圍內：上述 D1 至 D10 十項決策所對應的程式碼與測試調整。
- 不在範圍內：設定頁存檔佇列 / ownership 機制的重構、其他資料庫類型的語法或安全預設值、執行期（非 schema 產生期）的唯一鍵偵測、Share 錯誤碼集合變更、協議預設值表內容變更、`WithRevisionValidator` 的正式接線。

## Risks / Trade-offs

- **[cleartext 預設值變更為破壞性改動]** → 既有依賴 cleartext 且伺服器無 TLS 的 MySQL 連接器升級後會連線失敗。緩解：連線設定仍可明確設定 `allow_cleartext_passwords` 與 `tls` 兩個鍵回到舊行為；proposal 的 Impact 已標示此為破壞性改動，實作時需在錯誤訊息中保留驅動原始錯誤以便診斷。
- **[逾時值選擇]** → 10 秒對慢速後端可能過短，造成偽陽性的存檔未完成錯誤。緩解：逾時後走的是可重試路徑而非終局失敗，且 D1 / D3 修正後正常流程不應觸發逾時；逾時以具名常數表示，日後調整成本低。
- **[D3 讓非最新請求也能寫入終態]** → 舊請求的失敗可能覆蓋新請求剛設定的 `saving`。緩解：`save_state` 的終態寫入限定為 `saved` / `save-error`，佇列本身已序列化存檔（同一時間僅一個請求在飛行中），因此不會出現舊請求在新請求之後結束的情形；測試需覆蓋連續兩次存檔的狀態序列。
- **[D9 的三分支條件較複雜]** → 分支條件寫錯會讓 Share 停用時整段檢查再次失效，且不易由現有測試察覺。緩解：三個分支各自需有對應測試（停用且就緒時拒絕錯誤 token、停用且未就緒時放行、啟用時維持完整檢查）。
- **[D5 改變裝置載入時的設定內容]** → 正規化寫入 `station_no` 會讓裝置在使用者未編輯的情況下產生一次差異，可能觸發 autosave。緩解：正規化沿用既有的協議設定正規化流程（該流程已會改寫設定），實作時需確認不會因此把未修改的裝置標記為待存檔；測試需覆蓋「載入後不產生非預期的存檔請求」。
