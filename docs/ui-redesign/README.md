# UI 重構提案（兩入口）

決策版本：`1A / 2B / 3C`
- 1A：快速入口為預設
- 2B：Feature Flag 灰度上線
- 3C：快速入口與進階入口同時推進

## 文件（完整）
- 規格書（SPEC）：`two-entry-ui-spec.md`
- 實作計畫（PLAN）：`two-entry-ui-implementation-plan.md`
- 實作前準備：`pre-implementation-readiness.md`
- 設計與規劃補充：`two-entry-ui-redesign-plan.md`

## Gemini 生圖資產
- 流程圖：`preview_flow.svg`
- 最終線框預覽：`final_preview_wireframe.svg`

## 生圖標準流程
1) Gemini 產 SVG（指定「純 SVG XML」+ 尺寸 + 畫面結構）
2) 用 agent-browser 開啟 SVG 截圖為 PNG
3) 傳送 PNG 供使用者決策，再回填到 spec/plan
