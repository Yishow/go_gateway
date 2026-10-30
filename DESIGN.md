# Gateway 記錄介面設計契約

本文件是新簡易入口、編輯頁及相關置中 modal／對話框的設計 authority；目前是待實作規格，不宣稱產品已對齊。沿用同一 Studio V2、ThemeProvider、服務及設定，不重寫全站。

## 來源與取捨

- 產品基準：6f4a5fc。字型／品牌映射來源：`frontend/src/styles/tokens.ts`、`frontend/src/features/datalink/workbench-v2/tokens.ts`、`frontend/src/index.css`。
- Taste Skill v2.1.5，固定 [commit 12af471](https://github.com/Leonxlnx/taste-skill/tree/12af471f627f866411b4bd785b19646d3672784c)，採 `skills/taste-skill/references/app-screens.md` 的 app 路徑，不採 landing-page hero、照片或入場動畫。
- 新版候選方向是 Quiet mineral：暖中性亮底、非navy的中性深底、青綠單一重點色；使用者尚未驗收v5，不能聲稱已核准。保留字型載入、CJK/system fallback與功能語意，不保留已被否決的blue/slate palette。
- **Supersedes**: v4 待核准 blue/slate 品牌映射。v4被使用者否決，改用下表v5候選tokens；舊版只保留在歷史artifact，不作目前實作色值authority。新介面不用玻璃、裝飾漸層、發光或backdrop blur。
- 舊淺色 demo 僅歷史視覺來源；本契約取代 light-only 配色要求。mock 行為、固定成功或簡化協議不得作正式契約。
- Taste draw/lint/inspect scripts 及三種寬度截圖均 NOT RUN；只採純文字設計規範，不冒稱完整 Taste gate 通過。

## 頁面與引導

首頁回答「哪些資料可靠留下，現在需處理什麼？」；編輯頁回答「從哪裡來、何時記錄、存到哪裡、變更影響什麼？」。保留工作區與物件名稱，單一主要下一步；異常、草稿與正常記錄可同時到達。

來源→記錄方式→目的地逐段展開；進階參數保留於詳細設定或精確進階交接。保存、開始、採集中、SQL receipt、readback、cleanup 不混稱成功。零、缺值、等待、過期與未知分開。

## 先以使用者任務組織點位流程

永遠以使用者友善為出發點：先讓人知道正在選什麼、按下去改什麼、取消保留什麼，再決定版面。依 Taste app-screens 的單一問題／主要動作、同軸對齊、資料與計數同源、平面表格與真實空態原則，不以多加提示、卡片或換色掩蓋流程混亂；本輪不改既有 tokens。

- 主記錄頁回答「這組記錄有哪些點、各自存什麼、多久存一次？」。「記錄點位 · N 點」下只列已加入點及每member method／interval，唯一「選擇點位」入口；移除僅移出草稿，提供復原，不刪共用來源。來源設定是明示影響共用來源的次要動作，不混成記錄算法。
- 選點modal回答「這次要再加入哪些點？」。固定標題後只有同高的一組搜尋／設備／點位類型控制，搜尋最寬，三者共同過濾同一列表與符合數；「點位類型」使用量測值／累積值／狀態值，不混同float32等資料型別。identity欄最寬，安靜row rules，不放逐列設定按鈕或競爭建立toolbar。
- 結果摘要「符合 X 點 · 其中 Y 點已加入」，既有成員標「已加入」且不能在此取消。可捲body下固定footer把「待加入 N 點」、取消、唯一primary「加入 N 個點位」放一起；N僅計本次pending，0禁用；跨條件隱藏的pending另說明數量。計數與同一canonical集合推導，不與下方列表脫勾。
- 「找不到點位？新增來源」為單一次要流程入口；有條件無結果先清條件，單台零點設備顯示在該設備新增點位，完全無設備才引導新增設備。不把載入失敗顯示為空資料。設備filter必含已載入零點設備，切換立即改結果，不放第二個無關target selector。
- setup在同一active modal替換內容，不疊dialog。單設備scope直接標明唯讀目標；全部設備時在setup才明確選既有設備或新增設備，不暗選第一台。設備與範圍分別明確保存，已保存來源不因取消後一步或取消picker而刪除；未決操作查證，不聲稱取消成功。
- 主頁draft、picker pending及source form buffer分開；「加入」僅合併新identity到draft，不保存整份記錄、不啟動；「取消／關閉／Escape」只捨棄本次pending，不跳整份記錄捨棄。來源未送出buffer另經dirty guard；背景點擊不默默關閉。主頁尚未通過驗證的輸入文字、既有method／interval及同ID復原／重加意圖均保留。
- 同scope返回恢復query/filter/page/scroll/focus與pending；保存新範圍成功是明示例外：定位新點並說明清除阻擋條件原因，不自動勾選。移除／revision失效仍保留可辨識原因，重驗後才可加入，不拿同名／同位址替代ID。
- 一般modal尺寸沿75vw × 75dvh及既有小視窗例外。手機搜尋一列、兩filter同列、第一結果列在上方60%內（設計條件，以1440/768/390實際畫面確認），footer可換行但不遮內容；短列表不硬放pagination，大量資料沿25/50有界頁。主頁與picker的主要動作屬不同active context，不用背景primary爭搶注意力。

這是本輪流程與層級修訂，不是owner視覺核准、browser結果或後端完成證據；固定新版demo仍需owner checkpoint。

## 主題權威與相容遷移

- 單一 `ThemeProvider` 保存 preference `system | light | dark`，另推導 resolved `light | dark`；既有 consumer 的 theme 保持 resolved，不增加第二 provider。
- settings `auto` 映射 system。首次 app/backend 設定 hydration，只有原始 SettingItem[] 真正存在且有效的 theme key 才是 server preference；mapper 補出的 dark 不是意圖。
- 初始優先序：已保存有效 global server preference → 新版本機明確 preference → legacy localStorage.theme 的有效 light/dark → system。legacy 值可能只是舊程式記下 OS 結果，不能聲稱可辨識手動意圖；保守保留並可明確選回「跟隨系統」。
- 新本機 preference 使用獨立版本化 key（設計名稱 `gateway.theme.preference.v1`），不得把 resolved theme 寫回 preference；不覆寫舊 key 來製造手動選擇證據。server preference 與 UI 設定下拉使用同一 theme-only adapter；原設定頁明確整頁Save語意不改。OS change只更新resolved，不寫任何settings key。
- 開頁後的使用者明確選擇立即更新 provider 並只呼叫現有 updateKey('theme', value) 或等價 theme-only adapter，不呼叫 persistSettingsSnapshot／整頁 onSave，也不flush其他未保存業務設定；在同一 app/backend context 較舊 hydration 回覆不能覆蓋新選擇。保存 pending/error 明示，不假稱跨 reload 已保存；切 backend/app context 清除舊 request epoch；現有server theme位於global system_settings，不是workspace設定，切workspace本身不重設或另存theme。
- system 模式訂閱 matchMedia change，卸載清理；light/dark override 不受 OS 變更影響，回 system 立即採當下 OS。無 storage/matchMedia 或拒絕存取時不中斷頁面，system 無可讀 OS 時 fallback light，不自動存成 override。
- html.light/html.dark 與 native color-scheme 同步，V2 新 surfaces 不再強制 dark；body portal、dialog/menu/select/focus 使用相同 resolved theme。主題變化不能送設備／DB 命令、重設草稿或重新啟動記錄。
- 本輪遷移新首頁、編輯器及重用的受影響 modal/dialog/error controls；未觸及舊頁不全面換皮，但跨頁 provider 選擇保持一致，不新增相互矛盾的 theme 控制。

## Modal 尺寸與未決選項

使用者已確認：modal 本身 width=目前 browser viewport 的75%、height=75%，置中且內容內捲（75vw × 75dvh）。主頁正常 RWD，不是將頁面縮至75%。不使用側邊滑出設定頁；focus trap、返回焦點、Escape與dirty guard保留。

已選1／2／4＋A／B：modal固定標題／操作區、情境說明、同scope返回保留搜尋／篩選／捲動；手機近滿版，小型確認框較小。未選3 density、5 CmdK、6額外動效、7 maximize；8開發驗收元件預覽頁已選，但不放主要選單、不新增第三operator主頁。

尺寸分類（數值是設計條件，須以本批viewport驗收）：
- 一般設定／需要表單的impact/detail popup：寬度≥640px且高度≥480px時，modal自身75vw × 75dvh、置中；recording editor主畫面仍是第二個正常RWD頁，不搬入modal。
- 手機width<640px，或height<480px造成操作區擠壓的橫向小視窗：近滿版，寬calc(100vw - 16px)、高calc(100dvh - 16px)，扣safe-area且不超出可見viewport；這是已選A的響應式設計邊界。
- 小型確認（停止／移除／捨棄等短內容、無複雜表單）：width=min(420px, calc(100vw - 32px))，height:auto，max-height:calc(100dvh - 32px)，仍置中，長文body內捲。不可把schema impact或多欄設定偽裝小confirm省去資訊。
- modal分header/body/footer三區，僅body捲動，標題與主要操作可見；窄／短視窗header/footer可換行但不能遮住body、focus或錯誤；瀏覽器縮放／螢幕鍵盤時依可見viewport避遮蔽，焦點可捲入body。
- 關閉、Escape、Back和取消同dirty guard；有已送出操作不聲稱取消。焦點trap／恢復到觸發點；觸發物件已消失時回到同scope安全標題。
- context help放欄位附近，可展開補充名詞／範例，但必填、驗證失敗、停用原因、資料表／歷史影響、不可逆或待決風險與consent不能藏在help內，help本身可鍵盤開關。
- 返回保留搜尋／篩選／捲動只在同workspace＋entity/recording identity＋view context；它是短暫UI狀態不是第二資料authority。不跨scope恢復舊selection/receipt/action availability；scope/revision改變先重驗選取，已不存在物件解除可用選取權威，保留可辨識的失效列與原因供處理。回返後先套篩選再恢復scroll與focus，不用舊async回覆重建已失效項目。

這份修訂仍需新固定版本demo確認，不將需求選項同意視為整頁視覺核准。

## Primitive tokens 與語意映射

依 [designtoken.md spec](https://designtoken.md/spec) 組織 primitive → semantic → component，不另造第二個人讀token文件。以下primitive是本次Quiet mineral候選的有名值；新component只用下一節semantic角色，不直接取primitive。50–900各列為設計映射值。

| family | 50 | 100 | 200 | 300 | 400 | 500 | 600 | 700 | 800 | 900 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| mineral | #F7F8F5 | #EFEFEB | #E0E5DF | #C8D1C9 | #B8BEB9 | #929B94 | #586762 | #394A43 | #242624 | #191A19 |
| teal | #F1F8F4 | #E2F1EA | #BDE0D3 | #83C8B7 | #59AA98 | #398F7E | #267F70 | #1D6D63 | #175A52 | #124B44 |
| green | #F0FDF4 | #DCFCE7 | #BBF7D0 | #86EFAC | #4ADE80 | #22C55E | #16A34A | #15803D | #166534 | #14532D |
| yellow | #FEFCE8 | #FEF9C3 | #FEF08A | #FDE047 | #FACC15 | #EAB308 | #CA8A04 | #A16207 | #854D0E | #713F12 |
| red | #FEF2F2 | #FEE2E2 | #FECACA | #FCA5A5 | #F87171 | #EF4444 | #DC2626 | #B91C1C | #991B1B | #7F1D1D |

另有white #FFFFFF、black #000000及表面專用中性色#202220、#203332與selected dark #26473E。語意角色由下表唯一映射；舊blue/slate/cyan不混入新surfaces。
## 語意色彩 tokens

以下是設計值，不是已完成渲染驗收。正式實作在既有 `frontend/src/index.css` 的 light/dark selector 定義 CSS custom properties；`styles/tokens.ts` 與 V2 tokens 僅作同名 `var(--gw-...)` adapter，不重複另一套 raw hex。新 UI 與 portal 消費同一組。Demo 可複製這張表，但必須逐值比對且標明 mock-only。

| CSS 名稱（前綴 --gw-） | light | dark | 用途 |
| --- | --- | --- | --- |
| ground | #F7F8F5 | #191A19 | 畫布 |
| surface | #FFFFFF | #202220 | 靜態面板 |
| raised | #FFFFFF | #242624 | 浮動面板 |
| active | #E2F1EA | #26473E | 選取背景 |
| text | #203332 | #EFEFEB | 主要文字 |
| muted | #586762 | #B8BEB9 | 次要資訊 |
| faint | #586762 | #B8BEB9 | 輔助／placeholder，非不可讀淡字 |
| border | #C8D1C9 | #586762 | 裝飾分隔，不單獨表示可操作邊界 |
| control | #586762 | #929B94 | 必要控制項邊界 |
| focus | #1D6D63 | #83C8B7 | 焦點環 |
| accent | #1D6D63 | #83C8B7 | 連結／選取文字 |
| primary | #1D6D63 | #83C8B7 | 主要按鈕背景 |
| on-primary | #FFFFFF | #191A19 | 主要按鈕文字 |
| primary-hover | #175A52 | #BDE0D3 | 主要按鈕 hover |
| success-fg | #166534 | #BBF7D0 | 成功文字 |
| success-bg | #F0FDF4 | #14532D | 成功背景 |
| success-border | #15803D | #4ADE80 | 成功邊界 |
| warning-fg | #854D0E | #FEF08A | 警告文字 |
| warning-bg | #FEFCE8 | #713F12 | 警告背景 |
| warning-border | #A16207 | #FACC15 | 警告邊界 |
| danger-fg | #991B1B | #FECACA | 錯誤／危險文字 |
| danger-bg | #FEF2F2 | #7F1D1D | 錯誤背景 |
| danger-border | #B91C1C | #F87171 | 錯誤邊界 |
| unknown-fg | #394A43 | #EFEFEB | 未知／等待文字 |
| unknown-bg | #E0E5DF | #394A43 | 未知背景 |
| unknown-border | #586762 | #929B94 | 未知邊界 |
| overlay | rgba(25,26,25,0.45) | rgba(0,0,0,0.65) | modal 遮罩，不作文字底色 |
| shadow | 0 8px 24px rgba(25,26,25,0.16) | 0 8px 24px rgba(0,0,0,0.32) | 僅浮動面板 |

狀態必須伴隨文字／圖示，不只靠色彩。disabled 不降整體 opacity 至不可讀，沿用 muted 與明確理由。文字／placeholder／資訊性弱化字目標至少 4.5:1；控制邊界與焦點至少 3:1（依實際相鄰背景測試）。不得把 faint 放在 active/status 背景而不驗證。表格列不鋪滿彩色。

## 尺寸、字型與互動 tokens

- `--gw-font-sans`: Inter, "Noto Sans TC", "PingFang TC", "Microsoft JhengHei", system-ui, sans-serif；`--gw-font-mono`: "JetBrains Mono", monospace。沿用既有載入與 fallback，不增加 CDN 或字型依賴。等寬只用 ID／code，數值 tabular-nums、數值欄右對齊。
- 字級 small/body/title：12/14/24px，字重400/400/600，行高 1.5/1.5/1.25，letter-spacing均0（CJK不增加字距）；label14px/500/1.5/0，numeric14px/500/1.5/0；手機 label/data 至少14px。這些數字是設計條件，非量測結果。
- spacing scale：4/8/12/16/24/32/48px；半徑 control/panel/dialog：6/8/12px，pill9999px。
- control height desktop40px、compact32px（僅token，不實作未選density switch）；手機可點擊目標至少44px。輸入 label 置頂；placeholder 不代替 label。
- focus ring2px、offset2px；sticky/sidebar/dialog/toast 層級10/20/40/50，overlay 位於 dialog 下，互不遮蔽焦點。
- motion feedback150ms ease（僅基本操作回饋，不代表候選額外動效已批准）；不做入場、shimmer 或無限 pulse。prefers-reduced-motion: reduce 時移除 transform／循環／非必要 transition，狀態不靠動畫才看得到。
- 使用現有 lucide-react 同一 outline icon set，有文字或 accessible name；不使用裝飾 emoji 或每卡彩色 icon tile。
- 桌面雙區可用，手機單區與可返回 modal；表格必要橫捲只在自身容器，整頁不溢出。長名稱可換行、點位可搜尋；不為窄畫面刪參數。

## Typography / Spacing / Radius / Elevation

名稱為固定token identifier，CSS以 `--gw-` 前綴及連字號映射（如 typography.body.size → --gw-typography-body-size）。大型字級只定義詞彙，不授權增加hero或展示大標。

| Name | Size | Weight | Line-height | Letter-spacing | Use |
| --- | --- | --- | --- | --- | --- |
| typography.xs | 12px | 400 | 1.5 | 0 | 補充說明 |
| typography.sm | 14px | 400 | 1.5 | 0 | body，手機label/data下限 |
| typography.base | 16px | 400 | 1.5 | 0 | 長段落 |
| typography.lg | 18px | 500 | 1.4 | 0 | section |
| typography.xl | 20px | 600 | 1.4 | 0 | modal標題 |
| typography.2xl | 24px | 600 | 1.25 | 0 | page標題 |
| typography.3xl | 30px | 600 | 1.25 | 0 | 保留，不用於目前操作畫面 |
| typography.4xl | 36px | 600 | 1.2 | 0 | 保留，不用於目前操作畫面 |
| typography.5xl | 48px | 600 | 1.2 | 0 | 保留，不用於目前操作畫面 |

| Spacing token | Value | Use |
| --- | --- | --- |
| space.2xs | 4px | icon與短label |
| space.xs | 8px | 控制項內距 |
| space.sm | 12px | field群間距 |
| space.md | 16px | row/panel內距 |
| space.lg | 24px | section/card內距 |
| space.xl | 32px | 大section間距 |
| space.2xl | 48px | 頁面區塊 |
| space.3xl | 64px | 保留較大間距，不增加hero |

| Radius token | Value | Use |
| --- | --- | --- |
| radius.sm | 6px | button/input |
| radius.md | 8px | panel |
| radius.lg | 12px | modal |
| radius.xl | 16px | 保留，不預設套所有卡 |
| radius.full | 9999px | pill |

| Elevation token | Light | Dark | Use |
| --- | --- | --- | --- |
| elevation.none | none | none | 靜態面板 |
| elevation.sm | 0 1px 2px rgba(25,26,25,.12) | 0 1px 2px rgba(0,0,0,.24) | 小浮層 |
| elevation.md | 0 4px 12px rgba(25,26,25,.16) | 0 4px 12px rgba(0,0,0,.28) | menu |
| elevation.lg | var(--gw-shadow) | var(--gw-shadow) | modal |
| elevation.xl | 0 16px 48px rgba(25,26,25,.20) | 0 16px 48px rgba(0,0,0,.40) | 保留浮層，不疊靜態卡 |

- button.primary: background=primary; text=on-primary; border=none; radius=radius.sm; padding=8px 16px; hover:background=primary-hover; focus:ring=2px focus; disabled:background=surface; disabled:text=muted。
- card: background=surface; text=text; border=1px border; radius=radius.md; padding=space.md; shadow=elevation.none。
- input: background=surface; text=text; border=1px control; radius=radius.sm; padding=8px 12px; focus:ring=2px focus; placeholder=faint。
- navigation: background=surface; text=muted; selected:background=active; selected:text=accent; focus:ring=2px focus。

簡要視覺參考：`ground [ surface: text / muted | control input | primary CTA ]`；modal是raised+shadow，status只在相應內容用fg/bg/border；light與dark切換同角色不換資訊階層。

## Component states 與視覺參考

| 元件 | default | hover / selected | focus / disabled / error |
| --- | --- | --- | --- |
| primary button | primary + on-primary | primary-hover；pressed不移動版面 | focus環；disabled surface+muted+control、可讀原因 |
| secondary button | surface+text+control | ground；selected active+accent | focus環；disabled不隱藏label |
| input/select | surface+text+control、faint placeholder | control邊界，保留label | focus環；error danger邊界+文字原因+aria-describedby |
| list item | surface+text | ground；selected active+accent+非色彩選取標記 | 可鍵盤到達、明示disabled理由 |
| modal | raised+text+shadow，overlay | 不做裝飾動畫 | focus trap、返回focus、內捲可達錯誤 |
| status | 對應 fg/bg/border+文字 | 無假可點擊hover | unknown不套success；更新不只用動態顏色 |

視覺參考是本次固定版本 mock demo（待修訂／待owner確認），必須有相同 tokens 的 light/dark 真實畫面證據後才補版本與截圖；舊v2快取只供追溯，不當最終visual reference。元件預覽頁（選項8）已選，僅dev/test隱藏入口，重用正式components/tokens，禁止資料mutation、啟用與真credential；production build不bundle此開發入口（使用既有buildmodeguard，不增服務）。demo提供專用非主選單連結／hash，只含mock fixtures。CSS implementation與本文同步，token drift測試逐值比對；不靠口頭宣稱一致。

## 協議保真與示範邊界

Modbus TCP/UDP/RTU、FATEK FBs、MC 3E 是目前 V2 可達的五類；MQTT picker disabled，不因 backend 有 schema 就宣稱可完整記錄。完整欄位／驗證矩陣由 `streamline-recording-task-editor/design.md` 管理。

每個既有欄位只能標示「簡易頁可編輯」、「進階可達且原值保留」或「目前不支援」，不能消失。新 UI 不以顯示欄位重建 object 丟失未呈現設定。讀取項目數不等於 word 數；讀取排程不等於 SQL 寫入時間窗。Demo 每範圍 poll 不是正式 Rule 能力，不可冒稱已實作。

## 驗收與交接

- 元件／provider 測試：server key 有／無、mapper default、legacy light/dark、system 初載、OS live change、explicit override、回 system、storage/matchMedia 不可用、較舊 hydration、backend context切換、workspace切換不重設及保存失敗。
- 每種主題驗證首頁五態／混合、editor、modal、portal/native controls、錯誤/disabled/focus，主題切換保留 dirty form 與 operation；CSS token 與本文逐值一致。
- 靜態計算 token 文字對比可作輔助證據，不取代真實相鄰背景／focus／遮罩的渲染驗收。實際1440/768/390寬度、鍵盤與 reduced-motion 須記錄本批 screenshot/browser evidence；做不到明列 NOT RUN。
- 本輪 owner 必須重新看固定版本 demo 的來源／參數、生命週期、亮暗主題，之前 checkpoint 不視為同意。確認後才可能 apply 第一案；demo同意不代表正式 browser 或設備驗收。
- 每 change 跑相關測試，完整功能最後跑全套；不新增環境、不連真PLC/正式DB、不執行第三方 Taste scripts。沒有必要驗收不得自行切預設。

## 參數說明入口（已選 context help 的具體呈現）

每個設定參數的input/select/checkbox等控制旁提供圓圈ⓘ說明按鈕，具可讀accessible name（如「timeout的說明」），使用button type=button，不能誤submit表單。點擊／觸控／鍵盤Enter或Space開啟，不只靠hover。欄位label、必填、錯誤、停用理由、criticalrisk及consent仍直接可見，不藏ⓘ。

說明依參數需要呈現用途、單位／合法範圍、例子、何時需改；沒有source證據就不杜撰範圍、預設通用性或設備安全保證。沿既有validator、schema與protocol adapter，明示例如IP/serial是backend主機所用，採集週期不等於寫入窗。

ⓘ使用輕量help popover，不是設定表單modal，不套75%尺寸；遵守同tokens、對比、keyboard focus、可見viewport邊界。點外或Escape關閉、鍵盤關閉焦點回觸發ⓘ；點外關閉不搶走新點擊目標焦點。內容可讀、觸控可關，不讓popovers遮住必須處理的錯誤或主要操作。開／關不改欄位、不保存／probe／啟用。

coverage map對每個實際呈現參數列canonical field、ⓘ內容來源、demo呈現位置及formal差異；正式尚未於demo展示的欄位列在協議parity說明／進階交接，不假稱全欄位demo完成。這是已選#2的細化，不新增第三頁或新引擎。

## v4 未驗收與下一版密度修正

v4離線測試通過不代表使用者接受；已知邏輯／樣式缺陷現在修demo，不留正式版。本文件已替換成Quiet mineral候選，仍需新版視覺驗收；owner gate保持停止。

設備常用欄位緊湊對齊且常顯：protocol/mode、endpoint/serial、station/slave、timeout與必要連線參數，不把使用者常改項目藏advanced。桌面可3欄、窄畫面2欄／1欄（設計條件，不固定每裝置相同欄數）；ⓘ及label不犧牲。點位用compact rows並常顯name/address/readonly datatype/per-point wordorder/semantic，清楚呈現跨度、有效字序與inherit/override狀態，不每點大卡；低頻額外欄位才折疊。

首頁以記錄管理列表為主，設備短摘要最多4列（設計條件）＋同頁管理模式；完整列表可搜尋、依設備/狀態/目的地篩選、穩定排序與25/50分頁（設計條件），不無限鋪卡。設備完整管理模式也有自身搜尋/filter/page狀態；100設備/300規則fixtures只是測試設計，不是performance承諾。pagination是已載入資料的client projection，不宣稱已有跨設備後端索引。返回維持同scope的query/filter/page/scroll/focus，revision變化清失效selection並clamp page。

generic dialog不設定尺寸，以明確data-size=form|confirm等語意分類；短確認不可只檢查單一dialog id；「尚未套用的修改」必屬短確認，height:auto/fit-content、grid-template-rows:auto auto auto並解除native dialog inset stretching、bounded max-height、長文僅body內捲，不被regular75%覆寫。實際render驗證不可由CSS字串或jsdom替代。
