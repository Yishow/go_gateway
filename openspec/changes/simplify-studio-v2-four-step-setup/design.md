## Context

前置A/B/C/D全部驗證，以穩定API再改UI。已讀proposal。Step1 allTested與inferHydratedProgress不同、PROTOCOLS含MQTT但AddressParser不支援、detectAddressConflicts忽略device、Step4 selectedPlan與enabledTargetCount跨authority。前三步既有內容可保留，不需要新V3。

詳細來源見 [evidence](../../../docs/plans/studio-v2-write-groups/evidence.md)，共用欄位與政策見 [跨案契約](../../../docs/plans/studio-v2-write-groups/contracts.md)。這是擬議設計，不是完成報告。

## Goals / Non-Goals

**Goals:** 四步以相同persisted truth導航；操作簡單但保留安全gate；輸出只編輯WriteGroup。

**Non-Goals:** 新framework／視覺套件、MQTT協議實作、full runtime dashboard改版、報表/aggregation配置、恢復legacy studio。

## Decisions

1. **Step1連線。** form僅顯示選定協議必要欄位；保存/connection test/probe各有狀態。backend提供supported setup/parser/collector capabilities；MQTT不具完整V2路徑就disabled並附原因。initial/reload讀相同device validity及revision；改host/protocol後舊probe失效，取消/返回不偽造成功。未連線／probe未解決仍可保存draft並前後導航設定，視覺狀態標pending/unverified，僅相關activation/read操作受阻；不能把allTested沿用成所有Next/Save硬門檻。
2. **Step2採集點。** 顯示來源device、protocol、normalized area/address/type、polling period及少量真值preview，進階批量與geometry收折。conflict key含device與normalized address space，另外保留data width overlap檢查；不同device同40001合法，同device真重疊有定位到rule/point的錯誤。preview不是persisted完成；save/cancel/reload不丟未送草稿。
3. **Step3 Point-to-Tag。** 預設聚焦tag name/type/scale與sample preview，明示原值→轉換後值→quality/time。一次保存真正point/tag/mapping identity，基本group可直接使用；measurement physical semantics僅於選擇進階功能時確認。逐列error、IME、Enter/blur防重送及dirty draft保護沿用。
4. **Step4群組。** 只顯示group list/editor：哪些Tags一起寫、到哪裡、怎麼成一row、現在可否開始。managed/custom是同group switch且保留可相容bindings；legacy plan全不匹配時可create new或明確review migration，不能只剩disabled select。使用real metadata與reviewable suggestions，沒有欄位或permission failure不能填sample。共享column要同row不衝突且identity有證據；unsafe跨group/upsert阻擋。
5. **readiness/提交。** server依persisted group/connector/schema/member/applied revisions產生issue set，保留workspace/settings/readiness token。managed不另要求db.targets enabledCount；DB未選的Share-only不要求DB。改設定使相關preview/probe失效，Save是draft，Apply後才影響new intake；不能用全devices running或動畫結束宣布DB成功。
6. **可用性。** 既有元件/色彩延用；搜尋/問題篩選只影響顯示，bulk action說明作用scope/count。390/768/1440 CSS px是預定驗收viewport，不是已量測結果。鍵盤可到所有修復action、focus visible、欄位error連結label、table自有橫向scroll，en/zh-TW所有新文案入locale。runtime僅補所需group delivery資訊，維持原focused monitor路由。

### Migration Plan

先readiness/capability contract與focused tests，再逐步切前三步及Step4；舊editor只讀相容projection直到切換完成，不把兩個editor同時顯示成可修改。保留legacy compatibility API及draft；UI回退也須遵守server canonical authority，不能以舊客戶端繞過CAS。

### Open Questions

無外部視覺稿，延用既有元件。不虛構像素設計或性能成效；viewport/keyboard和錯誤路徑由F收集實際證據。

## Risks / Trade-offs

[更多安全檢查讓表面步驟變多] → 四步不增，進階設定收折、blocking issue有直接修復位置。

[frontend/backend不同版本] → server fail-closed、capability/version check及legacy adapter，不放寬readiness。

[空清單/讀取錯誤混淆] → 明確pending/error/empty/mismatch，均有下一步。
