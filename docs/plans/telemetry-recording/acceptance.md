# 全組驗收矩陣

以下是待實作驗收，非已通過產品測試。每個案例必須有來源、觀測時間、期望結果與實際證據。
自動測試優先使用固定時鐘；現場測試不得寫入未授權 PLC 或生產資料表。

## 數值與語意

| ID | 案例 | 期望結果 | 主責 |
| --- | --- | --- | --- |
| M01 | Modbus：多台，各3A+3V+1kW+1kWh | 七個量測曲線角色加一個counter，各台身份獨立 | 1、5 |
| M02 | MC：D0、D1、D2、D3-D4、D5、D6、D7 | 八word七主要項目，不重複占用D4 | 1 |
| M03 | raw -123 ×0.1 | -12.3°C，只換算一次 | 1 |
| M04 | uint64 9007199254740993 | API、DB、CSV每位數相同 | 1、2、6 |
| C01 | 12000.0→12000.4→12000.7kWh | 初筆只建基準，0.4+0.3=0.7 | 3 |
| C02 | 123→123 | 有效零用量，不是missing | 3 |
| C03 | counter 999→4，無reset證據 | uncertain，不用abs，不自動把4算用量 | 3 |
| C04 | 已確認新epoch 4→6 | 新段2，舊新中間仍不完整 | 3 |
| C05 | 65530→4，M=65536，可證單次回捲 | 10；不能排除多次則unknown | 3 |
| C06 | signed net100→98 | -2，不當reset | 3 |
| C07 | delta1重送同identity | 只算1 | 3、4 |
| C08 | 不同非重疊interval，各delta1 | 合計2，不按相同值去重 | 3 |
| C09 | 輪詢held delta寄存器，無sequence | 阻擋用量累加，仍可留原值 | 1、3 |
| W01 | 10持續10秒、20持續50秒 | mean1100/60，min10 max20 coverage60秒 | 3 |
| W02 | 60秒全缺資料 | count0，mean/min/max null，非0值 | 3 |
| W03 | 10:00讀1000，10:10讀1002且連續性可信 | 整段2，十分鐘各段未知，不全部記在10:10 | 3、6 |
| W04 | 午夜缺可信讀值 | 日量未知，顯示已知小計；不以當日首末筆假裝完整 | 3、6 |
| W05 | source clock倒退或同時刻衝突值 | uncertain／需修正，不以到達順序當真實先後 | 3 |
| R01 | 2kW有效持續1800秒 | 積分估算1kWh；不另加到權威counter | 3 |
| R02 | 120L/min有效持續60秒 | 積分估算120L=0.12m³ | 3 |
| S01 | 初筆run，30秒後stop，中間有效 | initial snapshot及轉態，觀察run30秒；無虚構先前事件 | 3 |
| S02 | run後長斷線 | 超過max_hold的時間unknown，不能全部算run | 3 |
| S03 | alarm bit0出現再解除 | 兩個具名轉態，原bitmask可查 | 3、6 |
| S04 | batch trigger持續high且壓力缺值 | 一筆partial batch，無舊壓力補值 | 2 |
| S05 | 文本批號或enum未知值 | 不做平均；原值與unknown標記可查 | 1、3 |
| D01 | 衍生公式除零/單位不合/循環 | 阻擋或invalid，不吐假數字 | 3 |
| D02 | 總表+其分表欲自動加總 | 警告／阻擋重複計量，需明確表達式 | 3、6 |

## 正確送達與恢復

| ID | 案例 | 期望結果 | 主責 |
| --- | --- | --- | --- |
| P01 | 12次採集後才flush | every_sample仍12筆，摘要用12筆，不只最後一筆 | 2、4 |
| P02 | 外部DB斷線，採集仍好 | 本機durable增加、delivered不假綠；恢復後原時間補送 | 4 |
| P03 | 外部commit後丟失response | 收到相同record id，最後只有一份正式紀錄 | 4 |
| P04 | checkpoint各交易邊界crash | 不漏已接受樣本、不重複計量、resume可重入 | 3、4 |
| P05 | 一慢一快兩目的地 | 快者照送，慢者獨立backlog，無雙倍報表 | 4、6 |
| P06 | pending期間改DB endpoint或disable | 舊record保持原destination；disable只停止新intake | 4 |
| P07 | 舊摘要revision晚送 | 不蓋掉新版修正摘要 | 3、4 |
| P08 | 證據已過期要求重算 | 拒絕並說明限制，不編造舊raw | 3 |
| P09 | 磁碟警戒／滿額 | 提前告警，受影響記錄暫停；不靜默刪未送資料 | 4 |
| P10 | retention時有落後consumer | 保留必要資料或告警，不破壞checkpoint/receipt | 4 |
| P11 | schema變更／權限失敗／不支援kind | typed blocked，可修復；不以sample columns fallback | 2、5 |
| P12 | gateway正常關機／強制終止／雙worker | 有界drain、恢復與唯一消費所有權 | 4 |

## 使用者與安全驗收

| ID | 案例 | 期望結果 | 主責 |
| --- | --- | --- | --- |
| U01 | 新用戶套三相電表範本 | 無SQL、無Row group知识完成真實保存 | 5 |
| U02 | MC混合型別與用途 | 格式占用、原值、換算、用途一眼可核對 | 1、5 |
| U03 | 已存連線密碼遮蔽 | 後端安全沿用；換身份才要求新驗證 | 2、5 |
| U04 | 試寫/建表preview後換目標 | token過期，後端拒絕舊操作 | 2 |
| U05 | test-write重按/逾時 | 查原結果不重複副作用；無read權限不稱verified | 2、5 |
| U06 | autosave fail、重新整理、兩分頁 | 草稿不假稱已存，不覆蓋別版，修復可續作 | 5 |
| U07 | 運行中改用途/倍率/來源 | 顯示影響與epoch，確認後切換，舊資料可追查 | 1、4、5 |
| U08 | 只選Local Modbus | 不需要DB；仍通過既有Share權限/probe/revision gate | 5 |
| U09 | 圖表、表格、CSV相同cutoff | 數值與完整度一致，文字formula injection被防護 | 6 |
| U10 | 沒有SELECT權限／primary offline | 可看送出狀態，不用local pending偽造DB歷史 | 6 |
| U11 | 窄視窗、鍵盤、色覺非依賴 | label/focus/error與主要操作均可用 | 5、6 |
| U12 | 部分設備測試失敗且明確未選 | 已選正常設備可繼續，未選設備不自動啟用 | 5 |
| U13 | 跨workspace存取/惡意欄名/外部路徑 | ownership與identifier/path驗證阻擋；日志無secrets | 1、2、6 |

## 環境與發行閘門

SQLite、PostgreSQL、MySQL 各測：建表／已有表、exact numeric、NULL、unique receipt、
commit未知、晚到更正、權限不足與schema drift。SQL Server不先宣稱新模式可用。
Windows/Linux/ARM 各列執行狀態；真MC PLC與真Modbus電表分開列出型号與韌體證據。
模擬器通過不等於現場通過；寫入生產資料前需部署者確認。
最小負載100設備×8項÷5秒=160樣本/秒；測無界增長、斷線配額與查詢p95，記錄參考硬體。
性能目標不是既有保證，未達標需限制可設定負載並說明，不允許靜默丟樣本。
發行需全部自動案例完成、必要實機項目被具名接受、rollback與backup/restore演練完成。
