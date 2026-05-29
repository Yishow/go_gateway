const DOC = {
  decisions: [
    { label: "/studio", tone: "paused", title: "完整版主線", detail: "功能應全面覆蓋需求面，但這一輪先暫停，只保留完整 inventory。" },
    { label: "/studio/v2", tone: "focus", title: "預設入口與重點施作", detail: "面向使用者，簡單、有指引、好觀察，應作為網站服務預設入口。" },
    { label: "/studio/runtime", tone: "info", title: "V2 的 post-setup observer", detail: "與 V2 一起補齊 setup 後觀察與 lifecycle 契約。" },
    { label: "/test + /gateway/*", tone: "experimental", title: "暫停 / 僅記錄", detail: "/test 先暫停，/gateway/* 保留 prototype 記錄，不列為近期主產品主線。" }
  ],
  tabs: [
    { id: "overview", label: "Overview", hint: "總覽、決策、圖表、驗證路徑" },
    { id: "studio", label: "/studio", hint: "完整版主線 inventory，先暫停" },
    { id: "studio-v2", label: "/studio/v2", hint: "重點施作、預設入口" },
    { id: "runtime", label: "/studio/runtime", hint: "post-setup focused observer" },
    { id: "test", label: "/test", hint: "工程工具 inventory，先暫停" },
    { id: "gateway", label: "/gateway/*", hint: "experimental surfaces" },
    { id: "api", label: "API Registry", hint: "頁面與 backend contract 索引" },
    { id: "roadmap", label: "Roadmap", hint: "change priority、缺口與順序" }
  ],
  surfaces: {
    studio: {
      title: "/studio 完整版主線",
      tone: "paused",
      audience: "進階操作者、維運人員、需要完整控制的人",
      strategy: "先暫停，但保留完整盤點；未來要回來做完整版整理。",
      role: "正式 datalink 主產品主線，功能最完整。",
      sourceOfTruth: [
        "WorkbenchProvider cross-step context",
        "datalink persisted CRUD state",
        "runtime snapshot / runtime SSE",
        "source rule candidate snapshots / review decisions"
      ],
      handoffs: [
        "<code>?step=</code> 與 <code>?target=</code> 深連結",
        "<code>selectedDeviceId</code> 跨步驟共享",
        "Step 2/3/4 共用 active rule / review context",
        "Output 與 runtime 觀察已深入主線，不是獨立頁外掛"
      ],
      keyFiles: [
        "frontend/src/pages/datalink/workbench/DatalinkWorkbenchPage.tsx",
        "frontend/src/pages/datalink/workbench/WorkbenchProvider.tsx",
        "frontend/src/pages/datalink/workbench/MuiDeviceStep.tsx",
        "frontend/src/pages/datalink/workbench/SourceCanvasSection.tsx",
        "frontend/src/pages/datalink/workbench/TagBindingStudio.tsx",
        "frontend/src/pages/datalink/workbench/LocalModbusBoard.tsx",
        "frontend/src/pages/datalink/workbench/DatabaseTargetBoard.tsx"
      ],
      actions: [
        { scope: "Shell", action: "step / target deep-link sync", frontend: "DatalinkWorkbenchPage + WorkbenchProvider", state: "WorkbenchProvider context", api: "none", status: "wired", notes: "以 <code>?step=</code>、<code>?target=</code> 維持主線切換與深連結。" },
        { scope: "Device", action: "列設備、建立設備、編輯設備、clone、draft 測試", frontend: "MuiDeviceStep / WorkbenchDeviceStep", state: "TanStack Query + form draft", api: "GET/POST/PUT /datalink/devices, POST /datalink/devices/test-draft", status: "wired", notes: "完整設備 CRUD 與 draft diagnostics 主線。" },
        { scope: "Device", action: "persisted activate / disable", frontend: "WorkbenchSourceRuntimeCluster 等", state: "mutation + optimistic cache", api: "POST /datalink/devices/:id/activate|disable", status: "wired", notes: "已存在，但語言上尚未全面收斂成 readiness-first。" },
        { scope: "Source", action: "source rule create/update/delete/enable/disable", frontend: "SourceCanvasSection", state: "queries + mutations", api: "GET/POST/PUT/DELETE /datalink/source-rules + enable/disable", status: "wired", notes: "Step 2 已經是 persisted source planning，不只是 UI 草稿。" },
        { scope: "Source", action: "create/delete points", frontend: "SourceCanvasSection", state: "query + mutations", api: "GET/POST/DELETE /datalink/points", status: "wired", notes: "Point lifecycle 與 source rule planning 深度綁定。" },
        { scope: "Source", action: "source step runtime live preview", frontend: "SourceCanvasSection + WorkbenchSourceRuntimeCluster", state: "runtime snapshot + SSE", api: "GET /datalink/runtime/status + GET /datalink/runtime/stream", status: "wired", notes: "runtime 已經深入主線，不應被視為額外頁面能力。" },
        { scope: "Tag", action: "review candidates / decisions", frontend: "SourceRuleTagReviewSurface + MuiTagIncidentDesk", state: "candidate snapshot + review decisions", api: "GET /source-rules/:id/candidates, GET/POST /tag-review-decisions", status: "wired", notes: "review 路徑存在，但正式 apply contract 尚未成為主線。" },
        { scope: "Tag", action: "tag/mapping CRUD 與 batch create", frontend: "TagBindingStudio", state: "queries + mutations", api: "GET/POST/DELETE /tags, POST /tags/batch, GET/POST/DELETE /mappings", status: "wired", notes: "目前仍以細粒度 CRUD 組裝主流程。" },
        { scope: "Output", action: "Local Modbus output 規劃、share server control、mapping sync", frontend: "LocalModbusBoard", state: "queries + imperative API calls", api: "modbus-share/*", status: "wired", notes: "Output 工具面成熟，但 review/apply contract 尚未完全收斂。" },
        { scope: "Output", action: "Database connector / table / schema / mapping / dry-run", frontend: "DatabaseTargetBoard", state: "query + imperative API calls", api: "db-targets/connectors/* + db-targets/mappings/*", status: "wired", notes: "是最完整的 operator output surface 之一。" },
        { scope: "Apply Contract", action: "source-rule tags apply / output apply", frontend: "review surfaces only", state: "目前多由 CRUD 組裝", api: "POST /source-rules/:id/tags/apply, /database-outputs/apply, /local-modbus/apply", status: "exists-not-wired", notes: "未來簡化 /studio 的關鍵收斂點。" }
      ],
      gaps: [
        "功能完整但 interaction density 太高，複雜度問題大於 API 缺口。",
        "Tag / Output 雖有 apply endpoints，但主線仍大量靠細粒度 CRUD 組裝。",
        "這一輪先暫停，不建議為配合 V2 先大改 /studio IA。"
      ]
    },
    "studio-v2": {
      title: "/studio/v2 面向使用者主線",
      tone: "focus",
      audience: "一般使用者、需要簡單引導式 setup 的操作者",
      strategy: "這一輪重點施作，並應成為網站服務預設入口。",
      role: "簡化、有指引、好觀察的 setup flow。",
      sourceOfTruth: [
        "useWorkbenchV2State reducer",
        "Step-based local draft state",
        "Step 4 handoff route state"
      ],
      handoffs: [
        "<code>/</code> 應預設導向 <code>/studio/v2</code>",
        "Step 4 commit 成功後 handoff 到 <code>/studio/runtime</code>",
        "若沒有 persisted <code>device_id</code>，目前 fallback 到 <code>/studio/runtime</code> 無 query"
      ],
      keyFiles: [
        "frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx",
        "frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx",
        "frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts",
        "frontend/src/features/datalink/workbench-v2/state/commitLog.ts",
        "frontend/src/features/datalink/workbench-v2/shell/resolveRuntimeDashboardDevice.ts"
      ],
      actions: [
        { scope: "Shell", action: "step rail / summary rail / settings 切換", frontend: "WorkbenchV2Shell", state: "local reducer", api: "none", status: "frontend-local-only", notes: "頁面導覽已成形，但尚未接正式 persisted workflow。" },
        { scope: "Step 1", action: "設備草稿編輯、協定切換、staged test", frontend: "Step1Device", state: "local reducer", api: "none", status: "frontend-local-only", notes: "connect / probe 階段視覺完整，但不是正式 datalink device lifecycle。" },
        { scope: "Step 2", action: "rule 草稿規劃", frontend: "Step2Rule", state: "local reducer", api: "none", status: "frontend-local-only", notes: "仍未落 persisted source rule contract。" },
        { scope: "Step 3", action: "mapping / transform 草稿調整", frontend: "Step3Mapping", state: "local reducer", api: "none", status: "frontend-local-only", notes: "尚未綁正式 tag/mapping APIs。" },
        { scope: "Step 4", action: "DB / output 草稿與 commit progress", frontend: "Step4Database + buildCommitLogSequence", state: "local reducer + simulated commit logs", api: "endpoint labels only", status: "frontend-local-only", notes: "最大的產品缺口；現在只是模擬 orchestration。" },
        { scope: "Settings", action: "connector pool / scheduler / modbus-share 設定草稿", frontend: "SettingsPage", state: "local reducer", api: "none", status: "frontend-local-only", notes: "適合作為未來後端設定 contract 的前端面。" },
        { scope: "Handoff", action: "Step 4 成功後前往 runtime dashboard", frontend: "resolveRuntimeDashboardDevice + navigate", state: "route query", api: "none", status: "semantics-gap", notes: "已修正不再傳草稿 device id，但 commit contract 仍不存在。" },
        { scope: "Default Entry", action: "作為網站服務預設入口", frontend: "App router", state: "routing contract", api: "none", status: "exists-not-wired", notes: "方向已確定，但 router 尚未改成預設落到 V2。" }
      ],
      gaps: [
        "缺正式 commit endpoint / orchestration contract。",
        "缺 persisted device_id / runtime handoff truth。",
        "缺 setup 完成後可直接觀察的最小 persisted summary。",
        "router 尚未改為預設入口。"
      ]
    },
    runtime: {
      title: "/studio/runtime Focused Runtime Dashboard",
      tone: "info",
      audience: "setup 後操作者、需要快速確認這台設備有沒有真的在跑的人",
      strategy: "跟隨 V2 一起施作，作為 post-setup focused observer。",
      role: "snapshot + SSE runtime 觀察面，而不是 fleet dashboard 首頁。",
      sourceOfTruth: [
        "runtime status snapshot",
        "runtime SSE value/status/heartbeat",
        "selected device_id search param",
        "points list for selected device"
      ],
      handoffs: [
        "Step 4 commit 後由 V2 handoff",
        "可直接打 <code>/studio/runtime?device_id=...</code>",
        "缺少 persisted context 時退回 missing-device-context"
      ],
      keyFiles: [
        "frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx",
        "frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts",
        "frontend/src/features/datalink/runtime-dashboard/useRuntimeStatus.ts",
        "frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts",
        "internal/api/handlers/runtime_handler.go",
        "internal/api/handlers/runtime_stream_handler.go"
      ],
      actions: [
        { scope: "Route", action: "解析 device_id 與 focused device", frontend: "useRuntimeDashboardState", state: "search params + devices query", api: "GET /datalink/devices", status: "wired", notes: "device context bar 與 focused device 來自這裡。" },
        { scope: "Snapshot", action: "讀 runtime snapshot", frontend: "useRuntimeStatus", state: "TanStack Query", api: "GET /datalink/runtime/status?device_id=...", status: "wired", notes: "首屏 truth。" },
        { scope: "Points", action: "讀 selected device points", frontend: "usePointsQuery", state: "TanStack Query", api: "GET /datalink/points?device_id=...", status: "wired", notes: "用來決定 stream point_ids 與表格內容。" },
        { scope: "Live Stream", action: "接 SSE value / status / heartbeat", frontend: "useRuntimeDashboardStream", state: "EventSource + buffered state", api: "GET /datalink/runtime/stream?device_id=...&point_ids=...", status: "wired", notes: "已具備 value/status/heartbeat 事件與 log panel。" },
        { scope: "Degraded Mode", action: "stream 掛掉後退回 polling", frontend: "useRuntimeDashboardState", state: "routeState machine", api: "runtime/status polling", status: "wired", notes: "已有 degraded path。" },
        { scope: "Missing Device Context", action: "無 device_id 或無法解析 device 時顯示 choose-device surface", frontend: "useRuntimeDashboardState + RuntimeDashboardPage", state: "routeState", api: "snapshot error parsing", status: "semantics-gap", notes: "目前部分仍依賴 device not found 錯誤字串。" },
        { scope: "Device Switch", action: "從 dashboard 重新切 device", frontend: "FocusedDeviceHeader", state: "search params mutation", api: "GET /datalink/devices + GET /datalink/points + runtime APIs", status: "wired", notes: "仍保持單台 focused context，不轉成 fleet console。" }
      ],
      gaps: [
        "runtime lifecycle 缺 starting / committed but not running 等正式語意。",
        "前端仍有少量字串判斷風險。",
        "與 V2 setup commit 之間還沒有正式 orchestration handoff contract。"
      ]
    },
    test: {
      title: "/test 工程工具",
      tone: "paused",
      audience: "工程測試使用者",
      strategy: "先暫停，但保留完整 API / 行為盤點。",
      role: "協議連線、讀寫、monitor、debug 的工程工具，不是產品主線。",
      sourceOfTruth: [
        "local config/profile state",
        "/api/v1/test/* APIs",
        "/api/v1/debug/* APIs"
      ],
      handoffs: [
        "與 /gateway/* prototype submit 有間接依賴",
        "不應被 /studio/v2 當正式 persisted contract 使用"
      ],
      keyFiles: [
        "frontend/src/pages/TestPage.tsx",
        "frontend/src/services/api.ts",
        "frontend/src/components/ConfigForm.tsx",
        "frontend/src/components/TestOperations.tsx",
        "frontend/src/components/MonitorControl.tsx",
        "frontend/src/components/DebugPanel.tsx"
      ],
      actions: [
        { scope: "Connection", action: "connect / disconnect / status", frontend: "useTestAPI", state: "hook + local page state", api: "POST /test/connect, POST /test/disconnect, GET /test/status", status: "wired", notes: "工程連線主線。" },
        { scope: "Operations", action: "read / write / batch", frontend: "useTestAPI + operations components", state: "local form state", api: "POST /test/read, /test/write, /test/batch", status: "wired", notes: "測協定能力最核心的工具面。" },
        { scope: "Monitor", action: "start / stop monitor", frontend: "MonitorControl", state: "local state", api: "POST /test/monitor/start|stop", status: "wired", notes: "monitor stream API 尚未全面盤進頁面。" },
        { scope: "Scripts", action: "script CRUD", frontend: "未全面接線", state: "n/a", api: "POST /test/script, GET/POST/DELETE /test/scripts/*", status: "exists-not-wired", notes: "API 存在，但目前頁面文件只先記錄。" },
        { scope: "Debug", action: "refresh / clear packets and logs", frontend: "useDebugAPI", state: "hook state", api: "GET /debug/packets, GET /debug/logs, DELETE /debug/clear", status: "wired", notes: "工程調試已可用。" },
        { scope: "Raw Debug", action: "send raw packet", frontend: "未全面接線", state: "n/a", api: "POST /debug/send-raw", status: "exists-not-wired", notes: "保留工程調試擴充點。" }
      ],
      gaps: [
        "script / monitor stream 還沒完整盤成前端工具面。",
        "與 /gateway/* 的依賴需要持續記錄，避免誤判為正式產品契約。"
      ]
    },
    gateway: {
      title: "/gateway/* Experimental Surfaces",
      tone: "experimental",
      audience: "模式實驗與原型驗證",
      strategy: "只保留記錄，不列為近期主產品重點。",
      role: "feature-flag 保護下的 dual entry prototype，不是正式 datalink persisted flow。",
      sourceOfTruth: [
        "gatewayDraftStore",
        "settings feature flag",
        "/test/connect"
      ],
      handoffs: [
        "/gateway/entry -> quick/expert",
        "quick/expert submit -> /test/connect",
        "未與 datalink persisted device/source/runtime contract 對齊"
      ],
      keyFiles: [
        "frontend/src/router/gateway.tsx",
        "frontend/src/pages/gateway/GatewayEntryPage.tsx",
        "frontend/src/pages/gateway/GatewayQuickSetupPage.tsx",
        "frontend/src/pages/gateway/GatewayExpertWorkbenchPage.tsx",
        "frontend/src/features/gateway/useGatewayDualEntryFlag.ts"
      ],
      actions: [
        { scope: "Entry", action: "讀 dual-entry feature flag 並決定是否可進入", frontend: "GatewayDualEntryGuard", state: "React Query flag", api: "GET /datalink/settings", status: "wired", notes: "ENABLE_GATEWAY_DUAL_ENTRY 是入口守門條件。" },
        { scope: "Quick Setup", action: "五步驟本地驗證與 draft 編輯", frontend: "GatewayQuickSetupPage", state: "gatewayDraftStore", api: "none until submit", status: "frontend-local-only", notes: "主要是 prototype UX，不是正式 persisted setup。" },
        { scope: "Quick Submit", action: "提交 connect test", frontend: "useTestAPI().connect()", state: "adapter payload", api: "POST /test/connect", status: "connect-only", notes: "成功不代表 datalink setup 完成。" },
        { scope: "Expert Workbench", action: "編輯 route table / plugin chain / raw manifest", frontend: "GatewayExpertWorkbenchPage", state: "gatewayDraftStore", api: "none until submit", status: "frontend-local-only", notes: "builder + raw manifest 仍是 prototype。" },
        { scope: "Expert Submit", action: "提交 connect test", frontend: "useTestAPI().connect()", state: "adapter payload", api: "POST /test/connect", status: "connect-only", notes: "也是 connect-only，不是 datalink commit。" }
      ],
      gaps: [
        "若未來要轉正，必須先決定它和 /studio/v2 的關係。",
        "目前不能被誤認為正式產品預設入口。"
      ]
    }
  },
  apiGroups: [
    {
      name: "Device / Protocol",
      rows: [
        { api: "GET /api/v1/datalink/devices", owners: "/studio, /studio/runtime", use: "列設備、解析 focused device", status: "wired" },
        { api: "POST /api/v1/datalink/devices", owners: "/studio", use: "建立 persisted device", status: "wired" },
        { api: "PUT /api/v1/datalink/devices/:id", owners: "/studio", use: "更新 device", status: "wired" },
        { api: "POST /api/v1/datalink/devices/test-draft", owners: "/studio", use: "測試 draft connection", status: "wired" },
        { api: "POST /api/v1/datalink/devices/:id/test", owners: "/studio", use: "測試 persisted connection", status: "exists-not-wired" },
        { api: "POST /api/v1/datalink/devices/:id/readiness", owners: "未來 /studio, /studio/v2", use: "readiness-first staged contract", status: "exists-not-wired" },
        { api: "POST /api/v1/datalink/devices/:id/activate|disable", owners: "/studio", use: "collector 啟停", status: "wired" },
        { api: "GET /api/v1/datalink/protocols", owners: "未來 /studio", use: "列 protocol 能力", status: "exists-not-wired" }
      ]
    },
    {
      name: "Source / Point / Review",
      rows: [
        { api: "GET/POST/PUT/DELETE /api/v1/datalink/source-rules", owners: "/studio", use: "persisted source planning", status: "wired" },
        { api: "POST /api/v1/datalink/source-rules/:id/enable|disable", owners: "/studio", use: "rule lifecycle", status: "wired" },
        { api: "GET/POST/DELETE /api/v1/datalink/points", owners: "/studio, /studio/runtime", use: "point list and create/delete", status: "wired" },
        { api: "POST /api/v1/datalink/points/:id/poll, /points/poll", owners: "未來 /studio", use: "即時輪詢", status: "exists-not-wired" },
        { api: "GET /api/v1/datalink/source-rules/:id/candidates", owners: "/studio", use: "review snapshot", status: "wired" },
        { api: "GET/POST /api/v1/datalink/source-rules/:id/tag-review-decisions", owners: "/studio", use: "review decisions", status: "wired" },
        { api: "POST /api/v1/datalink/source-rules/:id/tags/apply", owners: "未來 /studio", use: "正式 apply tags", status: "exists-not-wired" }
      ]
    },
    {
      name: "Tag / Mapping",
      rows: [
        { api: "GET/POST/DELETE /api/v1/datalink/tags", owners: "/studio", use: "tag CRUD", status: "wired" },
        { api: "POST /api/v1/datalink/tags/batch", owners: "/studio", use: "batch create tags", status: "wired" },
        { api: "POST /api/v1/datalink/tags/validate-key", owners: "未來 /studio", use: "tag key validation", status: "exists-not-wired" },
        { api: "GET/POST/PUT/DELETE /api/v1/datalink/mappings", owners: "/studio", use: "mapping CRUD", status: "wired" },
        { api: "POST /api/v1/datalink/mappings/preview", owners: "未來 /studio", use: "pipeline preview", status: "exists-not-wired" },
        { api: "POST /api/v1/datalink/mappings/validate-pipeline", owners: "未來 /studio", use: "pipeline validation", status: "exists-not-wired" }
      ]
    },
    {
      name: "Runtime / Monitoring",
      rows: [
        { api: "GET /api/v1/datalink/runtime/status", owners: "/studio, /studio/runtime", use: "runtime snapshot", status: "wired" },
        { api: "GET /api/v1/datalink/runtime/stream", owners: "/studio, /studio/runtime", use: "value/status/heartbeat SSE", status: "wired" },
        { api: "GET /api/v1/datalink/dashboard/stats", owners: "舊 dashboard contracts", use: "fleet style stats", status: "exists-not-wired" },
        { api: "GET /api/v1/datalink/dashboard/device-statuses", owners: "舊 dashboard contracts", use: "fleet device status", status: "exists-not-wired" }
      ]
    },
    {
      name: "Output",
      rows: [
        { api: "GET/POST/PUT/DELETE /api/v1/datalink/modbus-share/*", owners: "/studio", use: "local modbus share control and mappings", status: "wired" },
        { api: "GET/POST/PUT/DELETE /api/v1/datalink/db-targets/connectors/*", owners: "/studio", use: "connector management", status: "wired" },
        { api: "GET/POST/PUT/DELETE /api/v1/datalink/db-targets/mappings/*", owners: "/studio", use: "database target mappings", status: "wired" },
        { api: "POST /api/v1/datalink/source-rules/:id/database-outputs/apply", owners: "未來 /studio", use: "formal DB output apply", status: "exists-not-wired" },
        { api: "POST /api/v1/datalink/source-rules/:id/local-modbus/apply", owners: "未來 /studio", use: "formal modbus output apply", status: "exists-not-wired" }
      ]
    },
    {
      name: "Test / Debug / Experimental",
      rows: [
        { api: "GET /api/v1/datalink/settings", owners: "/gateway/*, future /studio/v2", use: "feature flags / settings", status: "wired" },
        { api: "POST /api/v1/test/connect", owners: "/test, /gateway/*", use: "connect-only test", status: "connect-only" },
        { api: "POST /api/v1/test/disconnect", owners: "/test", use: "disconnect", status: "wired" },
        { api: "GET /api/v1/test/status", owners: "/test", use: "connection state", status: "wired" },
        { api: "POST /api/v1/test/read|write|batch", owners: "/test", use: "engineering read/write ops", status: "wired" },
        { api: "POST /api/v1/test/monitor/start|stop", owners: "/test", use: "monitor control", status: "wired" },
        { api: "GET /api/v1/test/monitor/stream", owners: "/test", use: "monitor stream", status: "exists-not-wired" },
        { api: "GET /api/v1/debug/packets|logs, DELETE /api/v1/debug/clear", owners: "/test", use: "debug panels", status: "wired" }
      ]
    }
  ],
  roadmap: [
    { priority: "P0", name: "make-studio-v2-default-entry", target: "/ -> /studio/v2", why: "入口要先對齊使用者方向", status: "focus", notes: "調整 router default、fallback、legacy redirects。" },
    { priority: "P0", name: "integrate-studio-v2-commit-with-runtime-lifecycle", target: "/studio/v2 -> /studio/runtime", why: "V2 必須從 local-only setup 變成正式主線", status: "focus", notes: "補 commit endpoint / orchestration contract / persisted handoff context。" },
    { priority: "P0", name: "normalize-runtime-lifecycle-status-contract", target: "/studio/runtime", why: "runtime 不能再靠錯誤字串判斷部分狀態", status: "focus", notes: "補 starting、not_committed、device_not_found 等正式語意。" },
    { priority: "P1", name: "promote-runtime-dashboard-to-post-setup-truth", target: "V2 post-setup observer", why: "讓 runtime 真正成為 V2 setup 後的真實觀察面", status: "info", notes: "補 focused summary 與 setup summary handoff。" },
    { priority: "P1", name: "refactor-studio-mainline-surface-boundaries", target: "/studio", why: "完整版仍需要整理，但這輪先暫停", status: "paused", notes: "之後再處理 step boundary 與 review/apply 主線。" },
    { priority: "P1", name: "promote-source-rule-apply-contracts-in-studio", target: "/studio", why: "減少 CRUD 組裝複雜度", status: "paused", notes: "存在但目前不列為主優先。" },
    { priority: "P2", name: "align-gateway-pages-with-datalink-product-contract", target: "/gateway/*", why: "決定它是 prototype 還是正式產品線", status: "experimental", notes: "目前只保留記錄。" },
    { priority: "P2", name: "promote-readiness-first-device-flow", target: "/studio, /studio/v2", why: "讓 staged readiness 成為正式語言", status: "info", notes: "設備 readiness API 已存在但未全面接線。" }
  ],
  verificationPaths: [
    "<code>/</code> -> <code>/studio/v2</code> -> Step 1~4 -> commit -> <code>/studio/runtime?device_id=...</code> -> snapshot -> SSE connected",
    "<code>/studio</code> -> device CRUD -> source rules -> points -> tag review -> output planning -> runtime preview",
    "<code>/test</code> -> connect -> read/write/batch -> monitor -> debug packets/logs",
    "<code>/gateway/entry</code> -> quick/expert -> <code>/test/connect</code>"
  ],
  mdLinks: [
    { label: "README.md", href: "./README.md" },
    { label: "studio-mainline.md", href: "./studio-mainline.md" },
    { label: "studio-v2-runtime.md", href: "./studio-v2-runtime.md" },
    { label: "test-tooling.md", href: "./test-tooling.md" },
    { label: "gateway-experiments.md", href: "./gateway-experiments.md" },
    { label: "backend-api-registry.md", href: "./backend-api-registry.md" },
    { label: "gap-roadmap.md", href: "./gap-roadmap.md" }
  ]
};

const state = {
  activeTab: "overview",
  search: "",
  status: "all"
};

const tabListEl = document.getElementById("tab-list");
const tabContentEl = document.getElementById("tab-content");
const decisionBoardEl = document.getElementById("decision-board");
const searchInputEl = document.getElementById("search-input");
const statusFilterEl = document.getElementById("status-filter");

function chipClass(status) {
  const map = {
    focus: "focus",
    paused: "paused",
    info: "info",
    gap: "gap",
    experimental: "experimental",
    wired: "wired",
    "exists-not-wired": "exists-not-wired",
    "frontend-local-only": "frontend-local-only",
    "semantics-gap": "semantics-gap",
    "connect-only": "connect-only",
    "experimental-status": "experimental-status"
  };
  return map[status] || "info";
}

function matchesFilter(row) {
  const haystack = JSON.stringify(row).toLowerCase();
  const searchOk = state.search === "" || haystack.includes(state.search);
  const statusOk = state.status === "all" || row.status === state.status;
  return searchOk && statusOk;
}

function renderDecisionBoard() {
  decisionBoardEl.innerHTML = DOC.decisions.map((item) => `
    <article class="decision-card">
      <span class="chip ${chipClass(item.tone)}">${item.label}</span>
      <strong>${item.title}</strong>
      <p>${item.detail}</p>
    </article>
  `).join("");
}

function renderTabs() {
  tabListEl.innerHTML = DOC.tabs.map((tab) => `
    <button class="tab-btn ${tab.id === state.activeTab ? "active" : ""}" data-tab="${tab.id}">
      ${tab.label}
      <small>${tab.hint}</small>
    </button>
  `).join("");

  tabListEl.querySelectorAll(".tab-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      state.activeTab = btn.dataset.tab;
      render();
    });
  });
}

function renderOverview() {
  const surfaceMetrics = [
    { label: "/studio", maturity: 88, priority: "暫停", color: "linear-gradient(90deg, #f0a338, #ffd390)" },
    { label: "/studio/v2", maturity: 44, priority: "重點施作", color: "linear-gradient(90deg, #18b97d, #84f2c6)" },
    { label: "/studio/runtime", maturity: 62, priority: "跟隨 V2", color: "linear-gradient(90deg, #57a8ff, #abd4ff)" },
    { label: "/test", maturity: 68, priority: "暫停", color: "linear-gradient(90deg, #f0a338, #ffd390)" },
    { label: "/gateway/*", maturity: 26, priority: "僅記錄", color: "linear-gradient(90deg, #a986ff, #d2c0ff)" }
  ];

  return `
    <section class="section-card">
      <h3>目前決策摘要</h3>
      <div class="summary-grid">
        <div class="summary-block">
          <strong>產品主線策略</strong>
          <p class="muted"><code>/studio</code> 保留完整版定位但先暫停；<code>/studio/v2</code> 成為近期重點與預設入口；<code>/studio/runtime</code> 作為 V2 的 post-setup observer；<code>/test</code> 暫停；<code>/gateway/*</code> 僅保留實驗記錄。</p>
        </div>
        <div class="summary-block">
          <strong>最重要的技術缺口</strong>
          <p class="muted">V2 還缺正式 commit lifecycle；runtime 還缺完整 lifecycle semantics；router 還沒改成預設導向 V2。</p>
        </div>
      </div>
    </section>

    <section class="section-card">
      <h3>成熟度 / 優先順序圖表</h3>
      <div class="bar-list">
        ${surfaceMetrics.map((item) => `
          <div class="bar-row">
            <div><strong>${item.label}</strong><div class="muted">${item.priority}</div></div>
            <div class="bar-track"><div class="bar-fill" style="width:${item.maturity}%; background:${item.color};"></div></div>
            <div class="muted">${item.maturity}%</div>
          </div>
        `).join("")}
      </div>
    </section>

    <section class="section-card">
      <h3>驗證路徑</h3>
      <ul>${DOC.verificationPaths.map((path) => `<li>${path}</li>`).join("")}</ul>
    </section>

    <section class="section-card">
      <h3>對應 Markdown</h3>
      <div class="link-list">${DOC.mdLinks.map((item) => `<a href="${item.href}">${item.label}</a>`).join("")}</div>
    </section>
  `;
}

function renderSurface(surface) {
  const actions = surface.actions.filter(matchesFilter);
  const bannerTone = surface.tone === "focus" ? "focus" : surface.tone === "paused" ? "paused" : "gap";
  const bannerLabel = surface.tone === "focus" ? "重點施作" : surface.tone === "paused" ? "先暫停" : surface.tone === "experimental" ? "experimental" : "追隨施作";

  return `
    <section class="banner ${bannerTone}">
      <span class="chip ${chipClass(surface.tone)}">${bannerLabel}</span>
      <p style="margin:10px 0 0;">${surface.strategy}</p>
    </section>

    <section class="section-card">
      <h3>Surface 摘要</h3>
      <div class="summary-grid">
        <div class="summary-block"><strong>受眾</strong><p>${surface.audience}</p></div>
        <div class="summary-block"><strong>定位</strong><p>${surface.role}</p></div>
        <div class="summary-block"><strong>Source of Truth</strong><ul>${surface.sourceOfTruth.map((x) => `<li>${x}</li>`).join("")}</ul></div>
        <div class="summary-block"><strong>Route / Handoff</strong><ul>${surface.handoffs.map((x) => `<li>${x}</li>`).join("")}</ul></div>
      </div>
    </section>

    <section class="section-card">
      <h3>關鍵檔案</h3>
      <div class="path-list">${surface.keyFiles.map((file) => `<code>${file}</code>`).join("")}</div>
    </section>

    <section class="section-card">
      <h3>前端動作 / 行為 / API 對照</h3>
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Scope</th>
              <th>使用者行為</th>
              <th>前端實作</th>
              <th>State / Source</th>
              <th>API / Contract</th>
              <th>狀態</th>
              <th>備註</th>
            </tr>
          </thead>
          <tbody>
            ${actions.length > 0 ? actions.map((row) => `
              <tr>
                <td>${row.scope}</td>
                <td>${row.action}</td>
                <td>${row.frontend}</td>
                <td>${row.state}</td>
                <td><code>${row.api}</code></td>
                <td><span class="chip ${chipClass(row.status)}">${row.status}</span></td>
                <td>${row.notes}</td>
              </tr>
            `).join("") : '<tr><td colspan="7" class="muted">目前篩選條件下沒有資料。</td></tr>'}
          </tbody>
        </table>
      </div>
      <div class="table-note">可用上方搜尋與狀態篩選器縮小範圍。</div>
    </section>

    <section class="section-card">
      <h3>主要缺口 / 維護提醒</h3>
      <ul>${surface.gaps.map((gap) => `<li>${gap}</li>`).join("")}</ul>
    </section>
  `;
}

function renderApiRegistry() {
  return DOC.apiGroups.map((group) => {
    const rows = group.rows.filter(matchesFilter);
    return `
      <section class="section-card">
        <h3>${group.name}</h3>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th>API</th>
                <th>主用頁面</th>
                <th>用途</th>
                <th>狀態</th>
              </tr>
            </thead>
            <tbody>
              ${rows.length > 0 ? rows.map((row) => `
                <tr>
                  <td><code>${row.api}</code></td>
                  <td>${row.owners}</td>
                  <td>${row.use}</td>
                  <td><span class="chip ${chipClass(row.status)}">${row.status}</span></td>
                </tr>
              `).join("") : '<tr><td colspan="4" class="muted">目前篩選條件下沒有資料。</td></tr>'}
            </tbody>
          </table>
        </div>
      </section>
    `;
  }).join("");
}

function renderRoadmap() {
  const rows = DOC.roadmap.filter(matchesFilter);
  return `
    <section class="section-card">
      <h3>Change Priority</h3>
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Priority</th>
              <th>Change</th>
              <th>主要目標</th>
              <th>為什麼</th>
              <th>目前狀態</th>
              <th>備註</th>
            </tr>
          </thead>
          <tbody>
            ${rows.length > 0 ? rows.map((row) => `
              <tr>
                <td>${row.priority}</td>
                <td><code>${row.name}</code></td>
                <td>${row.target}</td>
                <td>${row.why}</td>
                <td><span class="chip ${chipClass(row.status)}">${row.status}</span></td>
                <td>${row.notes}</td>
              </tr>
            `).join("") : '<tr><td colspan="6" class="muted">目前篩選條件下沒有資料。</td></tr>'}
          </tbody>
        </table>
      </div>
    </section>

    <section class="section-card">
      <h3>目前順序</h3>
      <ol class="muted">
        <li>先把預設入口切到 <code>/studio/v2</code>。</li>
        <li>再讓 V2 commit 變成正式 persisted lifecycle。</li>
        <li>接著把 runtime lifecycle semantics 補成產品契約。</li>
        <li><code>/studio</code>、<code>/test</code>、<code>/gateway/*</code> 這輪只保留完整記錄，不擴張 scope。</li>
      </ol>
    </section>
  `;
}

function renderContent() {
  if (state.activeTab === "overview") return renderOverview();
  if (state.activeTab === "studio") return renderSurface(DOC.surfaces.studio);
  if (state.activeTab === "studio-v2") return renderSurface(DOC.surfaces["studio-v2"]);
  if (state.activeTab === "runtime") return renderSurface(DOC.surfaces.runtime);
  if (state.activeTab === "test") return renderSurface(DOC.surfaces.test);
  if (state.activeTab === "gateway") return renderSurface(DOC.surfaces.gateway);
  if (state.activeTab === "api") return renderApiRegistry();
  if (state.activeTab === "roadmap") return renderRoadmap();
  return "";
}

function render() {
  renderTabs();
  tabContentEl.innerHTML = renderContent();
}

searchInputEl.addEventListener("input", (event) => {
  state.search = event.target.value.trim().toLowerCase();
  render();
});

statusFilterEl.addEventListener("change", (event) => {
  state.status = event.target.value;
  render();
});

renderDecisionBoard();
render();
