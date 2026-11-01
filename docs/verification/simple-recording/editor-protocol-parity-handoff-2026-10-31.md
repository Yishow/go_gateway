# Editor 1.4 協議欄位局部交接

此文件只證明 1.4 的 protocol/parameter 部分，**1.4 不勾完成**；canonical acquisition preview 仍等 batching 4.1。沒有自行計畫合批、補 actual=0 或新增 planner。

## 已固定的局部版本

- 8 檔 frozen archive SHA256：`a874b0ea17afa77189ce40a8bcfa7ad97f3f18318be9b64751e2475cf2824f6a`。
- 完整 path/hash 與相關 logs 由本次 `editor14-protocol-evidence/{manifest.txt,frozen.sha256,HANDOFF.md}` 保存；8 檔是在已驗 run5 基礎精準替換，不用整份 f691 diff 覆蓋前一批。
- ConnectionConfigForm 兩 consumer 共用原 Input/Select/Field 與原 backend adapter，不另造保存 queue。
- RTU 補完整 serial bits/stop/timeout/設備預設格式，baud 依 server enum；非法 raw text 保留並阻止保存，沒有靜默改值。
- FATEK 明示 TCP/serial。明確切到 serial 才補缺值，render 不寫入；原 serial 配置、隱藏 TCP 欄位與未知欄位跨單欄修改保留。
- TCP/UDP/MC 設備預設格式與 MC routing 欄位可達；原 canonical/legacy alias 優先序沿既有 adapter，shared validator 只取原 canonical 值、不放寬數字型別。
- MQTT picker disabled，不新增密碼或宣稱來源可記錄；來源保存、取消、dirty guard、診斷與啟用邊界不變。

## 已執行與未執行

- 真 RED→GREEN；final related 13 files／132 tests PASS，含 shared two-consumer、MC aliases、RTU baud、FATEK mode 與原 device autosave guards。
- 精準 ESLint 0 issue；兩 locale JSON parse/key parity、line、diff、frozen hashes PASS；獨立 review findings 已修復。
- run6 exact58 原 strict tsc/Vite/static sync/Go build PASS；官方 chrome-devtools-mcp 1.12.0 在 Linux 的五協議設定操作、HTTP readback 與 SQLite 核對 PASS。
- FATEK serial station 單欄更改及 TCP/serial 往返、RTU重開後只改port、TCP/UDP/MC改host皆保原參數。Serial 的 last_test_at 為null，沒有硬體 test/read，不宣稱串列通訊成功。
- Legacy module aliases 由 unit/caller tests 承擔，Browser僅canonical配置；完整theme/viewport、preview及全部產品套件未驗。
- Browser報告：`linux-qa-tools/evidence/gateway-native-6/RUN6-QA.txt`，SHA256 `52ce51b44f5d05c3704492ae33bdfd7f9ac22a8d55f9a20d96615b56d02ad0a2`。本機evidence不是其他機器已存在的檔案。

## 對下一項的交接邊界

1.5 recent-connection defaults 只需要此已驗 protocol/parameter 契約及 1.3 matching device save/readback，不依賴未完成 preview。可以依正式 tasks/design 的分期 gate 先做；這不代表 cache 行為已驗，也不解除 1.4 或 editor 終驗。

1.4 共用表單保留既有設備的 inactive FATEK 欄位；1.5 新表單跨 subtype 的 whitelist/不相容欄位清除應放新設備 adapter，不能破壞此既有設備保真。Cache 必須 page-session、非權威、只由 matching confirmed snapshot 更新，不複製 secret/identity/proof，也不新增保存 owner。


## 後續 canonical preview 完成

原協議局部證據保持不變。Run17 已完成正式 scoped Start、真 available API→DOM、Refresh、布局與最後 console／導航驗收；前置五協議欄位與新預覽共同閉合 editor 1.4。完整邊界、修前失敗與修後證據見 `acquisition-preview-handoff-2026-10-31.md`。本機 HSL 模擬器不等於實體串口或 PLC 現場驗收，4.2 效能矩陣與七案整體終驗仍待完成。
