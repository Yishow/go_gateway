# V2 Step 3 Live Preview And Target Type Shortcuts Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `/studio/v2` Step 3 use real backend preview results and add fast `target_type` actions for the current row and all rows.

**Architecture:** Keep Step 3 reducer and workspace autosave unchanged. Add a small Step 3-local preview hook to own debounce, loading/error state, and stale response guard; keep payload preview local-only; add `target_type` quick actions in the right preview card and route them back through the existing reducer actions.

**Tech Stack:** React 19, TypeScript 5, Vitest, Testing Library, existing `mappingAPI.preview(...)`

---

## File Map

- `frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx`
  - Existing Step 3 container. Pass selected-row quick-action callbacks into the preview area.
- `frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx`
  - Existing preview shell. Upgrade it from pure display to preview-state container.
- `frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx`
  - Existing operator-visible pipeline renderer. Change it to render backend preview results instead of purely local `runScale/castValue/formatFinal`.
- `frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx`
  - Existing display-only payload JSON renderer. Keep it local and request-free.
- `frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts`
  - New hook. Own debounce, request lifecycle, and stale response guard for preview only.
- `frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx`
  - New component. Render `bool/int16/float64/string` chips plus `套用到全部列`.
- `frontend/src/i18n/locales/zh-TW/workbench-v2.json`
  - Add Step 3 preview loading/error/quick-action strings.
- `frontend/src/i18n/locales/en/workbench-v2.json`
  - Same keys in English.
- `frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx`
  - New focused async preview test file. Keep stale guard and request behavior out of the already-large `step3-mapping.test.tsx`.
- `frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx`
  - New focused quick-action test file.
- `openspec/changes/wire-studio-v2-step3-live-preview-and-target-type-shortcuts/tasks.md`
  - Mark task completion after code/test verification.

## Implementation Notes

- Do **not** add a new V2 preview API. Reuse `mappingAPI.preview(...)` from `frontend/src/services/datalink.ts`.
- Do **not** move preview state into `useWorkbenchV2State` or `mappingReducer`.
- Do **not** make `PayloadPreview` issue `fetch`, `axios`, `useQuery`, or `useMutation` calls.
- Keep the existing footer bulk-apply in `MappingTable.tsx` unchanged. The new preview-card apply-all is `target_type`-only.
- `frontend/tests/unit/workbench-v2/step3-mapping.test.tsx` is already over the soft line limit. Put new async and quick-action coverage in new test files instead of bloating that file.

### Task 1: Add a focused Step 3 live preview hook

**Files:**
- Create: `frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts`
- Create: `frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx`
- Modify: `frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx`

- [ ] **Step 1: Write the failing live preview tests**

Create `frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx` with fake timers and a mocked `mappingAPI.preview`.

```tsx
vi.mock('../../../src/services/datalink', () => ({
  mappingAPI: {
    preview: vi.fn(),
  },
}));

it('sends preview request for the selected row after debounce', async () => {
  vi.useFakeTimers();
  vi.mocked(mappingAPI.preview).mockResolvedValue({
    raw_value: 243,
    final_value: 24.3,
    step_results: [
      { step_index: 0, step_type: 'decode', input_value: 243, output_value: 243, error: '' },
      { step_index: 1, step_type: 'scale', input_value: 243, output_value: 24.3, error: '' },
      { step_index: 2, step_type: 'cast', input_value: 24.3, output_value: 24.3, error: '' },
    ],
  });

  render(<TransformPreview point={point} mapping={mapping} rawSeed={243} />);
  expect(mappingAPI.preview).not.toHaveBeenCalled();

  await vi.advanceTimersByTimeAsync(250);
  expect(mappingAPI.preview).toHaveBeenCalledTimes(1);
});
```

Also add failing expectations for:

- loading state while request is in flight
- backend error state when preview rejects or returns `error`
- no request when `point/mapping/rawSeed` are null

- [ ] **Step 2: Run the new preview test and verify it fails**

Run:

```bash
cd /Users/yishow/prj/go_gateway/frontend
npm run test -- --run tests/unit/workbench-v2/step3-live-preview.test.tsx
```

Expected: FAIL because `TransformPreview` is still a pure display component and there is no live preview hook yet.

- [ ] **Step 3: Write the minimal live preview hook and wire it into TransformPreview**

Create `useStep3LivePreview.ts` as a Step 3-local hook. Keep it tiny and focused.

```ts
export function useStep3LivePreview(
  point: Point | null,
  mapping: Mapping | null,
  rawSeed: number | null
) {
  const [state, setState] = useState<{ status: 'idle' | 'loading' | 'success' | 'error'; data?: MappingPreviewResponse; error?: string }>({ status: 'idle' });
  const requestIdRef = useRef(0);

  useEffect(() => {
    if (!point || !mapping || rawSeed === null) {
      setState({ status: 'idle' });
      return;
    }

    const currentRequestId = ++requestIdRef.current;
    setState({ status: 'loading' });
    const timer = window.setTimeout(async () => {
      try {
        const data = await mappingAPI.preview({
          raw_value: rawSeed,
          transform_pipeline: buildPayload(mapping, point).transform_pipeline,
        });
        if (requestIdRef.current === currentRequestId) {
          setState(data.error ? { status: 'error', error: data.error, data } : { status: 'success', data });
        }
      } catch (error) {
        if (requestIdRef.current === currentRequestId) {
          setState({ status: 'error', error: error instanceof Error ? error.message : 'preview failed' });
        }
      }
    }, 250);

    return () => window.clearTimeout(timer);
  }, [point, mapping, rawSeed]);

  return state;
}
```

Then update `TransformPreview.tsx` to:

- call the hook
- render explicit loading and error blocks
- pass preview data to `PipelineSteps` only on success

- [ ] **Step 4: Run the preview test and verify it passes**

Run:

```bash
cd /Users/yishow/prj/go_gateway/frontend
npm run test -- --run tests/unit/workbench-v2/step3-live-preview.test.tsx
```

Expected: PASS with one request after debounce, explicit loading, explicit error, and no request for the empty state.

- [ ] **Step 5: Commit the preview-hook foundation**

```bash
cd /Users/yishow/prj/go_gateway
git add \
  frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts \
  frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx \
  frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
git commit -m "補齊 Step 3 live preview hook"
```

### Task 2: Render backend preview results and guard against stale responses

**Files:**
- Modify: `frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx`
- Modify: `frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx`
- Modify: `frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx`

- [ ] **Step 1: Extend the preview test with stale-response coverage**

Add a delayed-response test to `step3-live-preview.test.tsx`.

```tsx
it('keeps the latest preview result when an older request resolves later', async () => {
  vi.useFakeTimers();
  let resolveFirst!: (value: MappingPreviewResponse) => void;
  let resolveSecond!: (value: MappingPreviewResponse) => void;

  vi.mocked(mappingAPI.preview)
    .mockImplementationOnce(() => new Promise((resolve) => { resolveFirst = resolve; }))
    .mockImplementationOnce(() => new Promise((resolve) => { resolveSecond = resolve; }));

  const { rerender } = render(<TransformPreview point={point} mapping={mappingA} rawSeed={243} />);
  await vi.advanceTimersByTimeAsync(250);

  rerender(<TransformPreview point={point} mapping={mappingB} rawSeed={243} />);
  await vi.advanceTimersByTimeAsync(250);

  resolveSecond(successB);
  await screen.findByText('99.00');

  resolveFirst(successA);
  expect(screen.getByText('99.00')).toBeInTheDocument();
});
```

- [ ] **Step 2: Run the preview test and verify stale-response coverage fails**

Run:

```bash
cd /Users/yishow/prj/go_gateway/frontend
npm run test -- --run tests/unit/workbench-v2/step3-live-preview.test.tsx
```

Expected: FAIL because the current hook either renders local pipeline math or allows an older response to overwrite the newer one.

- [ ] **Step 3: Change PipelineSteps to render backend preview output instead of local-only math**

Refactor `PipelineSteps.tsx` so it accepts preview data derived from backend results rather than calculating operator-visible values entirely from `runScale/castValue/formatFinal`.

Representative target shape:

```tsx
export interface PipelineStepsProps {
  point: Point;
  mapping: Mapping;
  preview: MappingPreviewResponse;
}
```

Render strategy:

- decode step uses `preview.raw_value`
- scale/cast/final use `preview.step_results` / `preview.final_value`
- remove direct operator-visible dependence on `runScale/castValue/formatFinal`

Keep the order `decode -> scale -> cast -> final`.

- [ ] **Step 4: Run the preview test again and verify it passes**

Run:

```bash
cd /Users/yishow/prj/go_gateway/frontend
npm run test -- --run tests/unit/workbench-v2/step3-live-preview.test.tsx
```

Expected: PASS with stale-response protection and backend-driven output rendering.

- [ ] **Step 5: Commit the backend-result rendering pass**

```bash
cd /Users/yishow/prj/go_gateway
git add \
  frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx \
  frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx \
  frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
git commit -m "改接 Step 3 backend preview 顯示"
```

### Task 3: Add target type quick actions for the selected row and all rows

**Files:**
- Create: `frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx`
- Create: `frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx`
- Modify: `frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx`
- Modify: `frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx`
- Modify: `frontend/src/i18n/locales/zh-TW/workbench-v2.json`
- Modify: `frontend/src/i18n/locales/en/workbench-v2.json`

- [ ] **Step 1: Write the failing quick-action tests**

Create `frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx`.

```tsx
it('updates the selected row target_type when a quick-action chip is clicked', () => {
  const dispatch = vi.fn();
  render(
    <Step3Mapping state={state} dispatch={dispatch} onContinue={vi.fn()} onBack={vi.fn()} />
  );

  fireEvent.click(screen.getByTestId('target-type-chip-int16'));
  expect(dispatch).toHaveBeenCalledWith({
    type: 'updateMapping',
    pointId: 'p-01',
    patch: { target_type: 'int16' },
  });
});

it('applies the selected row target_type to all rows', () => {
  const dispatch = vi.fn();
  render(
    <Step3Mapping state={state} dispatch={dispatch} onContinue={vi.fn()} onBack={vi.fn()} />
  );

  fireEvent.click(screen.getByTestId('btn-apply-target-type-all'));
  expect(dispatch).toHaveBeenCalledWith({
    type: 'bulkApplyTransform',
    fromPointId: 'p-01',
    fields: ['target_type'],
  });
});
```

Also add a no-selected-row expectation so the quick-action controls are not interactive when the preview is empty.

- [ ] **Step 2: Run the quick-action test and verify it fails**

Run:

```bash
cd /Users/yishow/prj/go_gateway/frontend
npm run test -- --run tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
```

Expected: FAIL because the preview card has no quick-action UI and `Step3Mapping` does not pass those callbacks yet.

- [ ] **Step 3: Add the minimal quick-action component and wire callbacks through Step3Mapping**

Create `TargetTypeQuickActions.tsx` and keep it display-focused.

```tsx
const QUICK_TARGET_TYPES = ['bool', 'int16', 'float64', 'string'] as const;

export const TargetTypeQuickActions = ({
  currentType,
  disabled,
  onSelect,
  onApplyAll,
}: {
  currentType: string;
  disabled: boolean;
  onSelect: (value: string) => void;
  onApplyAll: () => void;
}) => (
  <div data-testid="target-type-quick-actions">
    {QUICK_TARGET_TYPES.map((value) => (
      <button
        key={value}
        type="button"
        data-testid={`target-type-chip-${value}`}
        disabled={disabled}
        onClick={() => onSelect(value)}
      >
        {value}
      </button>
    ))}
    <button
      type="button"
      data-testid="btn-apply-target-type-all"
      disabled={disabled}
      onClick={onApplyAll}
    >
      套用到全部列
    </button>
  </div>
);
```

Wire it like this:

- `Step3Mapping.tsx`
  - derive selected point ID
  - pass `onTargetTypeSelect` and `onApplyTargetTypeAll`
- `TransformPreview.tsx`
  - render `TargetTypeQuickActions` only when a selected mapping exists

Use existing reducer actions only:

- current row: `updateMapping`
- all rows: `bulkApplyTransform` with `fields: ['target_type']`

- [ ] **Step 4: Run the quick-action test and verify it passes**

Run:

```bash
cd /Users/yishow/prj/go_gateway/frontend
npm run test -- --run tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
```

Expected: PASS with current-row updates, target-type-only apply-all, and inactive controls when no row is selected.

- [ ] **Step 5: Commit the quick-action pass**

```bash
cd /Users/yishow/prj/go_gateway
git add \
  frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx \
  frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx \
  frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx \
  frontend/src/i18n/locales/zh-TW/workbench-v2.json \
  frontend/src/i18n/locales/en/workbench-v2.json \
  frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
git commit -m "補齊 Step 3 目標型態快捷設定"
```

### Task 4: Final verification and OpenSpec task closure

**Files:**
- Modify: `openspec/changes/wire-studio-v2-step3-live-preview-and-target-type-shortcuts/tasks.md`
- Verify only: `frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx`

- [ ] **Step 1: Run the focused Step 3 test files together**

Run:

```bash
cd /Users/yishow/prj/go_gateway/frontend
npm run test -- --run \
  tests/unit/workbench-v2/step3-live-preview.test.tsx \
  tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx \
  tests/unit/workbench-v2/step3-mapping.test.tsx
```

Expected: PASS. Existing `step3-mapping.test.tsx` should still confirm payload preview rendering and core Step 3 integration behavior.

- [ ] **Step 2: Run the frontend build**

Run:

```bash
cd /Users/yishow/prj/go_gateway/frontend
npm run build
```

Expected: PASS with no TypeScript or Vite build errors.

- [ ] **Step 3: Run diff hygiene**

Run:

```bash
cd /Users/yishow/prj/go_gateway
git diff --check
```

Expected: no output.

- [ ] **Step 4: Mark the Spectra task file complete**

Update `openspec/changes/wire-studio-v2-step3-live-preview-and-target-type-shortcuts/tasks.md` and check off all completed items.

- [ ] **Step 5: Commit the verification and task closure**

```bash
cd /Users/yishow/prj/go_gateway
git add \
  openspec/changes/wire-studio-v2-step3-live-preview-and-target-type-shortcuts/tasks.md
git commit -m "完成 Step 3 live preview 與 target type 快捷設定驗證"
```
