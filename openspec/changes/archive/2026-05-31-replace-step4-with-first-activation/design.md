## Context

在前面的 autosave changes 完成後，Step 4 不再是「把資料真正寫進後端」的地方。使用者最新決策是：Step 4 只負責第一次啟動所有合法、可用、尚未啟動的設備；已經在跑的設備之後改設定直接靠 autosave 套用。

## Goals / Non-Goals

**Goals:**

- 讓 Step 4 只做 first activation
- 啟動目前所有合法、可用、尚未啟動的設備
- 逐台回傳啟動成功/失敗
- 即使部分失敗，仍能前往 runtime

**Non-Goals:**

- 不再把 Step 4 當成資料 commit
- 不重做 runtime 工作區頁面本身
- 不在這個 change 內處理 invalid edit stop

## Decisions

### First activation targets all eligible devices

Step 4 會一次對 singleton workspace 內所有合法、可用、尚未啟動的設備做 activation。這直接符合使用者最後決定。

### Activation returns per-device results

逐台成功/失敗必須可觀察，不能只有整批成功或整批失敗。這對應使用者指定的 partial success。

### Navigation remains available after partial failure

只要至少有成功設備，使用者就能前往 runtime 觀察成功的那些設備；失敗結果則留在 Step 4 顯示。

## Implementation Contract

- Behavior:
  - Step 4 最後按鈕只會啟動尚未啟動且目前合法/可用的設備
  - 每台設備各自回傳 `success` 或 `failed`
  - 部分失敗不阻止使用者前往 runtime
- Interface / data shape:
  - `POST /api/v1/datalink/studio-v2/workspace/activate`
  - response 至少包含：
    - `workspace_id`
    - `results[]`
    - each result: `device_id`, `status`, `message`
- Failure modes:
  - 若沒有任何符合條件的設備，回應需明確告知「無可啟動設備」
  - 單台 activation 失敗不得回滾已成功啟動的其它設備
- Acceptance criteria:
  - backend activation tests 覆蓋多台 eligible device、部分失敗、無可啟動設備
  - frontend Step 4 tests 覆蓋結果列表與 partial success 導頁行為
- Scope boundaries:
  - In scope: first activation API、Step 4 UX 語意切換、partial success 結果
  - Out of scope: runtime workspace device switching、live update apply
- Not complete if:
  - Step 4 還在做整批資料 commit
  - partial failure 仍會卡死不能進 runtime
  - 已啟動設備再次被 Step 4 重覆啟動

## Risks / Trade-offs

- [Risk] 舊 commit UI 文案和測試會和新 activation 語意衝突 → Mitigation: 在這個 change 內一併改成功卡/按鈕語意與相關測試
- [Risk] eligible device 判斷若與 availability contract 不一致會誤啟動 → Mitigation: 直接依賴 workspace availability state，不再各自判斷
