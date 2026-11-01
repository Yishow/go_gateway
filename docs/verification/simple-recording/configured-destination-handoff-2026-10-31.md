# 已設定目的地選擇驗收

歸屬 streamline-recording-task-editor 2.4。沿用原 workspace write-groups 查詢，不是實際資料庫全表清單；不新增 inventory API、schema 掃描、credential 讀取或配置權威。

## 行為與保護

- 候選保留 workspace／connector／database／schema／table 完整識別，顯示原記錄名稱、表及資料庫。相同名稱不同連線／DB 不合併。
- 只按後端既有長表 scope（workspace、connector、database、schema、storage mode、layout）判斷相容，不在前端重算 canonical table hash。別組寬表 owner 拒絕重用，直接顯示原因，提供原記錄的 scoped 編輯入口。
- 選取只填四個 intent 欄位：connector ID／revision、schema、table。沒有複製來源 record/group ID、schema digest、readiness 或 ownership proof。
- 原 Save owner 在 fresh fetchLatest 後核對所選 reference 仍存在且 revision／destination identity 相同，否則拒絕 mutation。原 schema／CAS／readiness／ownership 不變。
- 同頁原記錄入口走既有 open／dirty guard；locator 走原 scoped target／router guard，保留 Keep／明確 Discard。開啟或選取不自動 Save／Apply／Start。

## 精確驗證與修正

- Model 四個真 RED 後，新 model／chooser 11 案例通過；直接相關七檔 79 案例通過，覆蓋 identity、long compatibility、foreign-wide、loading/error/stale、fresh Save reference 拒絕及同 guard 原記錄入口。
- 八個歸屬檔 lint／typecheck、locale JSON／舊 tokens 不變及差異檢查通過；GroupEditor 最終 488 行。獨立 exact-baseline／frozen hash review 通過，未跑全專案測試。
- run25 真 Browser 證 long intent-only 選取、new draft 身份不變、canonical table 唯讀、切 layout 拒選及真 foreign-wide owner／原記錄入口。合法寬表前置僅 UI Save-only，未 Apply／Start，沒有 SQL／API 塞造 identity。
- run25 發現切 long→wide 時總警告誤稱 stale。精確 RED 後改為顯示原判斷的實際 layout／scope／revision 原因，既有 selection／Save gates 未改；受影響 10 案例通過。
- 真截圖發現新 chooser 的 gw-secondary 沒有定義，且技術識別字壓過記錄名稱。改用既有 shared Button／asChild 與已定義 gw tokens；record 名稱為主標，表／DB／schema 明列，UUID／layout／mode／revision 放技術詳細，關鍵拒絕原因不折疊。
- 共享 Button 的 secondary 色彩 token 在此 repo 亦未定義，因此再以新 chooser 局部、既有 light/dark gw tokens 補齊正常／hover／disabled／focus；不修改 global CSS／Button。該 style 精確 RED → 2 PASS，兩檔 lint／typecheck 通過。保留 v2 已建置但未跑 Browser 的紀錄，不冒稱已驗視覺。

## run26 窄 Browser closure

v3 真 Dark／Light 設定下的桌面與 390px 窄版截圖，root 已看實際像素：新 chooser 的人讀名稱、動作邊界、禁選、鍵盤焦點可辨，關鍵原因可見，長識別字沒有橫向溢出。長表選取仍純 intent；切 wide 正確顯 layout 原因；原 wide 入口 Keep 保 draft／focus，Discard 開啟正式原記錄。

只有明確的兩次 theme PUT，沒有額外 group／source mutation、Apply 或 Start；最後導航及 console 無 React #185／uncaught。自有服務 20:36:59 UTC 停止。原 run25 功能及 failed 視覺／原因證據保留，不把修後截圖替換原失敗。

同名不同 DB／connector 及 stale／error 邊界依相關單元／整合測試，不假稱每項都有實際 Browser 故障注入。測試副本沿未啟動的合法設定備份，原資料及目標檔保護由 QA 原始紀錄核對。

## 後續 3.3／4.1 必須閉合的既有頁面問題

本項不是整頁雙主題完成。真 Light 下外殼及舊卡仍深色、部分新元件變淺色，形成混色；須在原 3.3 統一頁面 tokens 並真 Browser 重驗。窄版 schema／main 等人讀短字不應因 break-all 被拆字，亦在同項修正。

run25 原 Advanced 的 Choose per-point results 曾被 sticky header 遮住，兩次工具點擊未切換，PageUp 使它可見後一次成功。保留幾何／截圖及鍵盤可達性追蹤，不把此舊入口問題歸因新 chooser，也不當作已修；原 3.3 的 scroll／focus／RWD 及 4.1 終驗須確認。

## 凍結識別

- v3 八檔 manifest SHA256：92f550fd4db1075bf35656c5b25e5ea43ed4a540b123dae56eca20394c2f2271。
- 證據組：editor24-configured-destinations-evidence、gateway-native-25／26；保存 task baseline、每版差異／hash、精確命令、截圖、network、SQL 及服務清理紀錄。
- 下一項 3.1 接原 save／check／schema impact；3.2 才接明確 scoped source Start。完整七案尚未完成，不同步 main 或宣稱全部通過。
