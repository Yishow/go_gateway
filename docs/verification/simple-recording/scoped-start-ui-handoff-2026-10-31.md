# 獨立記錄 Scoped Start：前端與真窗口交接

歸屬 streamline-recording-task-editor 3.2。承接 preparation-flow 與 scoped-start-backend 交接；沒有另一個啟動 engine、ledger、store 或隱藏 Basic 記錄。

## Scope 與恢復契約

- 同原 useCircuitProfileStart／session store／completion owner，以明確 saved_record 分支驗一筆已保存 group、exact 唯一 device 集合；舊 Basic 沒有 discriminator 時仍單設備／最多兩個 profile group。DB setup head 與 Share runtime head／settings revision／readiness token 分域驗證。
- 新 host 先讀 canonical saved group、readiness 與 live source／connector，不把 draft 當已保存 scope。原 request 保存及 roundtrip 可證後才 POST，double click／late owner 不新增命令。
- Lost ACK 無 operation ID 時，明確原 request 唯讀 lookup；找到 ID 後接原 GET／completion。unknown／partial／retry_same_request 不換意圖；malformed／unreadable／competing retained owner 拒覆寫。Start 保留原 recovery fail-closed，並非要求先執行 test-write；3.1 schema 的 optional-store 例外不外推。
- 已證歷史 terminal 若與新 saved revision／members 不符，必須重新 review 才能開始新意圖；不 rebase 舊 request、不 auto-stop 他筆。歷史成功不能觸發 completion reload 撤掉新 review lease。
- First SQL 使用同 applied revision／target 的 scoped receipt 與 precise timestamp，不借同毫秒舊 receipt、全域寫入計數或另一個窗口。Window/outbox row.record_id 不等 configuration record_id，不能猜相等。

## Normal／locator 入口

兩入口共用 SavedRecordingStartReview。Locator-only 只接受 WriteGroupSection 當前已驗 editorEntry／hydration、live route lease；自己 Save 後用更新 entry，不借外層初始 URL snapshot 或完整 workspace state。foreign／stale／suspended／A→B→A／late read 均拒絕。

明確開 locator review 才讀原 Modbus Share status／CheckHydration，核同 workspace 與原 runtime barrier；沒有 bootstrap、global autosave 或 DB-head 代替 Share-head。成功後沿原 scoped reload；舊 review readonly 不取得新寫入權。

QA31 發現舊 CTA 文案誤稱 Basic，以及 Start 成功刷新後 Close 列表又顯 outer Draft。最窄修為通用「開始前檢查」，原路由不變；同 route 最後已驗 receipt 只作列表顯示，route／loader／suspension 變即失效。原 fresh Open／CAS 保留，不由展示 receipt 授權新操作。

## 精確測試與獨立審查

兩設備 owner roundtrip／completion、legacy／forged 拒絕、真 exact POST、double click／late、原 request lookup、terminal 新意圖／unknown、receipt 與不可讀 owner 均有精確 RED／GREEN。初期五檔20案例通過；後續原 caller 111／113，兩個舊導航案例在單 hook baseline 也同點失敗，僅補既有2.3 pending／Discard 操作，保原 late／foreign／no-resubmit 斷言後通過。

Locator／normal／原 caller 五檔91案例通過；歷史 completion 缺陷真2RED→18選測通過。最後兩個 caller 缺陷真6RED→65相關案例，另3個 label／原路由 no-mutation 案例通過。精確 lint、typecheck、差異檢查通過；各階段不同命令／hash 證據保留，沒有整包重跑或宣稱一次全綠。

獨立 v4 24-path review 無剩餘 findings；backend 七檔另有獨立 review／race／vet／lint，不能將前端 review 當後端或全案 approval。

## QA31 真兩設備完整窗口

Safe 未開始設定建立全新 fixture31／target，正式 UI 建一筆 managed-wide mixed：A.VA mean15、B.gauge mean30、A.Energy delta60，兩個 configured logical devices 使用本機 HSL simulator，不稱兩台實體 PLC。原 counter 確認沒有建立 Basic group。

Normal review 精確兩設備／三 members 後 Cancel，零 Start；再由原 Home locator 進同 group，唯一 POST345 Start200。保留原 Share proof，ACK 先等待，後取得 matching committed receipt，不把 ACK 當 SQL。

僅對自有 loopback simulator 白名單詞對以標準 FC16／FC03 seed/readback：gauge12.5、25，counter100。完整 UTC 窗口22:31:00–22:32:00，counter 於 :15.034→107、:40.008→119；真 SQL 同 row delta19、mean15=12.5、mean30=25。目標／outbox 的 effect、digest、window record ID 一致；first-reading required 失敗整 row skip、not_due NULL 及 first-partial 原始證據均保留。

Native scope 為 A兩個來源／4 words、B一個來源／2 words，非全部配置九點。API 原生 fixture另證第三未選設備保持 draft。本輪 Browser 不為重複證明多造第三設備。

## QA32 Fresh 窄 caller 回歸

沒有複製 active31 config 指向舊 target；另用安全未開始備份建立 fixture32／新 target，一筆 VA latest5s。Normal review Cancel→locator 唯一新 Start→同 session Close 列表 Saved 而非 Draft；reopen 先 canonical list 再 fresh editor，saved／applied revision 相符。

13筆真SQL12.5／13張receipt及 GET218 effect／digest／record ID join 通過。這輪只重驗新 CTA／成功返回，不重做 QA31 兩設備／60秒計算；兩層證據分開。

兩輪最後導航／console 無 React #185 或 JS exception；setup 中間非法地址與舊 runtime head 曾回400／409，最後合法值／新 head200恢復，沒有冒稱刻意負向。外部 font 拒連仍明列。自有服務已全停，active configs／targets 成對一致備份；QA32其他32份舊DB（含active31 pair）未變。

## 證據與後續

- Frontend v4 manifest：62352fc7ffadaa22d3a59f9b565bb8682a9d347cd067a6e89b76357432a6a209。
- QA31 report：aecd7bafe36513ce67c6d061da475e56f1f04606b2242d9d4cb5060fbeb715e2。
- QA32 report：a234a8366a43fd4d9e6c2f831f12e9b3f816da882b745060e1797db531620935；exact391 union：aeb433f1d327d86d33d0521ddc8b061e42b7a6ddf77d6d087ed7a4263c79b6a0。
- editor32-scoped-start-evidence、mixed-start-evidence、gateway-native-31／32 保存完整 baseline／patch／命令／network／SQL／截圖。QA32部分前置 schema response body 已 expired，只以當時 request/status／UI／SQL證據陳述；Start／fresh GET／delivery 全 body 保留。
- Submitted lost ACK、corrupt store 等是單元／原生 API 證據，不製造 Browser 丟包。尚未跑全套；原3.3主題／sticky／參數help、4.1整體驗收、lifecycle／adoption仍未完成，未同步 main。
