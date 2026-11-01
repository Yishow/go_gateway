# 記錄編輯器 1.3：來源範圍保存與返回交接

## 完成邊界

2026-10-31：完成 `streamline-recording-task-editor` 的 1.3，不代表整個 change、七案、正式預設入口或現場驗收完成。沿用原 rule/mapping 保存 queue、validator 與 canonical identity 讀回，不新增 decoder、持久化草稿或採集引擎。

- 同 workspace、device、request、modal session 與 buffer generation 的保存結果，必須有正式 rule、mapping、point、Tag identity 才能返回；單一 POST 成功不是完成。
- 未知結果保留輸入並查原請求，不自動重送；晚回覆不關新表單、不抹除較新輸入。已存設備／範圍不因取消 picker 或記錄草稿而刪除。
- 新點返回同設備的有界選點列表，明示清除遮擋搜尋／類型條件的原因，不自動勾選或加入記錄；原成員與尚未保存意圖保留。
- 範圍格式走既有 Rule.data_format；單位走原 mapping queue 與 Tag 讀回，絕不塞入 Rule API。此項不含 1.6 的逐點 inherit／override／reset controls。
- 採集排程只顯示精確 point/group 對應的讀回配置，不代表正在採集；無資料誠實顯示未取得，不猜預設週期。唯讀摘要不阻塞保存、不授予開始權限。
- 全域排程入口沿既有進階工作區，在新分頁依「設定 → 排程器」到達；不是本範圍群組週期的直接編輯器。原分頁未保存內容保留，沒有重新掛載舊 PollingGroup 管理頁。

## 實際版本與歸屬

本批去重共 76 個產品／測試檔；完整 paths、逐檔 hash 及 provenance 保存在本次驗證的 `editor13-final-scope-audit/`。不含 OpenSpec artifacts 或 collector WIP。

| 來源 | 已發布／驗證基準 |
| --- | --- |
| 原 range22 | `60b346fe` |
| readback8 | `e8936de9` |
| offline14 | `4de7c221` |
| validation7 | `e639a467` |
| save-readback24 | `09369187` |
| 本次隔離 QA | `f6916e9c` 加 exact52：inactive7、locator3、autosave10、roundtrip11、ownsave7、parity14 |

52 檔符合最終 run5 union hash，其餘 24 檔符合既有 checkpoint。舊版本測試只代表其當時版本，不能加總或代替後續修改的驗證。跨元件修正仍有原 owner 的獨立 review，不以此文件取得其他工作者的檔案所有權。

## 測試、建置及獨立覆核

- 原具名 range-owner、submission、readback、queue-observation、device-read、source-flow、selection-stack、picker-state 測試涵蓋 scope/session/generation、late／unknown、防重送、取消、同名／同位址 identity、MC compatibility 及 offline preparation。
- 正式 HTTP＋SQLite readback 契約位於 `recording_editor_range_readback_contract_test.go`；proposal 的空 identity 不能冒充已存 mapping，Tag 空單位可合法省略。
- 參數補齊先取得真正行為失敗：兩檔 11 例，10 FAIL／1 PASS；非缺 module 或語法失敗。
- 參數相關最終選測：21 files／205 tests PASS；精準 ESLint、diff-check、行數檢查 PASS。與歷史 300／195 例等結果重疊，不相加。
- exact52 首次完整 strict TypeScript 失敗 TS18048；保留失敗紀錄。僅增加 observed 顯式 guard，未使用 assertion／any；之後 schedule 4 tests、該元件 ESLint 及獨立 diff review PASS。
- 2026-10-31 10:15:47 UTC：最終 exact52 原 strict tsc、Vite、嵌入靜態資源同步、Go production build PASS；建置前後與 Browser 後 hash 全相同。
- 最終 parity patch SHA256：`ac92ad028a486b4c9a27ce406d998fbde1c6bd112f0ceac8ae0b5f14cf655fa0`；14 檔 hash manifest SHA256：`47bd6f5b42bb4556fa9d595c611b21ffc0f6796e12bfcce3538441ea0a7fedcd`。
- 最後逐條需求覆核未發現阻止 1.3 結束的具體缺口；不是其他 tasks 的完成證據。

## 真實 Browser 結果

Linux 隔離 Chromium，產品操作使用官方 chrome-devtools-mcp 1.12.0；agent-browser 不可用時依指定順序回退。原生桌面只啟動隔離控制程序，沒有用 DOM evaluate 偽造操作。所有 fixture 變更均經 UI；SQL 僅唯讀驗證與一致備份。HSL 是本地試用模擬器，不是實體 PLC。

- run4：全部來源停用仍可由首頁編輯原群組；保存名稱後不需重新載入即可選點／新增來源。非空 picker 的 pending、query、device filter、焦點在來源返回後保留；Cancel 不改群組。
- run4：已存零點設備新增範圍，真 source POST201＋兩筆 mapping POST201，返回兩個未選取新點；取消仍保留來源，不改原成員。
- run4：connector 名稱／長路徑與真晚 ACK 保留較新輸入；重新載入及 SQLite 一致。真 stale revision409 顯示安全錯誤並保留草稿與原記錄內容；名稱更新不移動群組目的地。
- run5：float32 10 項顯示 20 個位址（41501–41520）；單次 source POST201 帶 CDAB、倍率 0.1、偏移 -2，不帶 unit；10 筆 mapping POST201 與 canonical Tag GET 均為 kWh，SQLite 的規則、點位與單位一致。
- run5：保存後直接返回 picker、0 selected／Add0 disabled；Cancel 後原群組與完整成員列不變。全域設定新分頁可到 System Settings／Scheduler Core，本頁格式、單位與其他草稿保留；Keep editing／Discard 均依原契約運作。
- run5：所有自有服務正常停止、兩個隔離 DB integrity OK 並保留一致備份，沒有 Apply、start 或 target test-write。

完整命令、原始 UI／HTTP／SQL 證據與 screenshot 索引：本次驗證目錄 `linux-qa-tools/evidence/gateway-native-4/RUN4-QA.md` 與 `gateway-native-5/RUN5-QA.txt`。此文件不把本機 evidence 路徑當成其他機器已存在的檔案。

## 明確保留的限制與後續

- run5 未額外執行 browser unknown／accepted-write ACK-loss，也未單獨操作 saved-inactive 的已填值排程摘要；對應既有具名測試仍是這些邊界的證據，不宣稱 Browser 全矩陣通過。
- 後續 run7 真 fresh Step3 重現舊名稱前綴：第一筆 POST 已帶前綴，17 筆 request/ACK 名稱一致；先聚焦再填、Ctrl+A 鍵盤替換及重載／SQL 對照正確。本觀察處置為 Chrome DevTools MCP 1.12.0 填表順序相容性限制，不修改產品；工具清空後才 focus 與受控欄位重繪的互動是最有支持的解釋，但非逐事件因果證明，也不代表所有 fresh 手動情境已驗。後續 QA 依 [測試指引](../../testing.md) 逐欄鍵盤替換與核對 payload/ACK；正常鍵盤亦失敗時重開產品調查。
- 雙主題／全部 viewport／鍵盤視覺矩陣、參數 help 與密度屬後續 editor 3.3 等任務；目前白色 modal 與深色 shell 的不一致未被隱藏。
- 1.4 協議／planner preview、1.6 逐點格式、2.1 逐成員計算 UI、lifecycle、adoption 及七案完成後全套測試均未由本項解除。未宣稱 Windows、實體 PLC、正式 DB、部署或預設切換驗收。
