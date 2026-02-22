import type { Tag, Point } from '../../../types/datalink';

/** Tag 編輯暫存結構（僅可編輯欄位） */
export interface PendingTagEdit {
  display_name: string;
  unit: string;
  description: string;
}

export interface SmartDashboardTagPanelProps {
  selectedSourceAddress: string;
  activePointForLink: Point | null;
  linkedTag: Tag | null;
  tagLinkMode: 'existing' | 'create';
  setTagLinkMode: (mode: 'existing' | 'create') => void;
  selectedTagIdForLink: string;
  setSelectedTagIdForLink: (tagId: string) => void;
  tags: Tag[];
  handleLinkTagToSelectedAddress: () => void;
  updateMappingPending: boolean;
  createMappingPending: boolean;
  newTagKey: string;
  setNewTagKey: (value: string) => void;
  newTagDisplayName: string;
  setNewTagDisplayName: (value: string) => void;
  handleCreateTagAndLink: () => void;
  createTagPending: boolean;
  tagLinkActionMessage: string;
  linkedTagAffectedMappingsCount: number;
  tagEditDisplayName: string;
  setTagEditDisplayName: (value: string) => void;
  tagEditUnit: string;
  setTagEditUnit: (value: string) => void;
  tagEditDescription: string;
  setTagEditDescription: (value: string) => void;
  handleSaveLinkedTagEdit: () => void;
  updateTagPending: boolean;
  pendingTagEdit: PendingTagEdit | null;
  handleConfirmTagEdit: () => void;
  clearPendingTagEdit: () => void;
  tagEditMessage: string;
}

/** Tag 連結與全域 Tag 編輯面板 */
export default function SmartDashboardTagPanel({
  selectedSourceAddress,
  activePointForLink,
  linkedTag,
  tagLinkMode,
  setTagLinkMode,
  selectedTagIdForLink,
  setSelectedTagIdForLink,
  tags,
  handleLinkTagToSelectedAddress,
  updateMappingPending,
  createMappingPending,
  newTagKey,
  setNewTagKey,
  newTagDisplayName,
  setNewTagDisplayName,
  handleCreateTagAndLink,
  createTagPending,
  tagLinkActionMessage,
  linkedTagAffectedMappingsCount,
  tagEditDisplayName,
  setTagEditDisplayName,
  tagEditUnit,
  setTagEditUnit,
  tagEditDescription,
  setTagEditDescription,
  handleSaveLinkedTagEdit,
  updateTagPending,
  pendingTagEdit,
  handleConfirmTagEdit,
  clearPendingTagEdit,
  tagEditMessage,
}: SmartDashboardTagPanelProps) {
  return (
    <div className="space-y-3 border-t border-white/5 p-4">
      {/* 區塊說明與標題列 */}
      <p className="text-[11px] leading-relaxed text-slate-400">
        將記憶體格位上的點位與 Tag 建立映射，採集值會寫入 Tag 供下游使用。請先選取網格上「已建立點位」的格位，再選擇或新建 Tag 並套用連結。
      </p>
      <div className="flex items-center justify-between" title="Tag Linkage（標籤連結）：此區塊用於建立或更新「點位 → Tag」的映射關係。左側為目前選取的記憶體格位位址，選定點位並選擇 Tag 後，點位採集到的資料會經由映射寫入該 Tag。">
        <p className="text-xs font-semibold tracking-wide text-slate-200">Tag Linkage</p>
        <span className="font-mono text-[10px] text-slate-400" title="目前選取的記憶體格位位址（例如 Modbus 40001、FATEK D0）。點擊網格或側欄點位會更新此值。">{selectedSourceAddress || '—'}</span>
      </div>

      {/* 目前選取狀態 */}
      <div
        className="grid grid-cols-3 gap-1.5 rounded-lg border border-white/10 bg-slate-900/60 p-2.5 text-[11px]"
        title="目前操作上下文：格位＝網格上選到的位址；點位＝該格位已建立的點位名稱（若無則顯示 —）；Tag＝該點位目前連結的 Tag key（若無則顯示 —）。"
      >
        <div title="記憶體格位的位址（如 40001、D100）。同一設備內唯一，用於識別採集來源。">
          <p className="text-slate-500">格位</p>
          <p className="mt-0.5 font-mono text-slate-200">{selectedSourceAddress || '—'}</p>
        </div>
        <div title="點位（Point）：代表該格位已建立的採集點，含名稱、資料型別、輪詢群組等。未建立點位時無法進行 Tag 連結。">
          <p className="text-slate-500">點位</p>
          <p className="mt-0.5 truncate text-slate-200">{activePointForLink?.name || '—'}</p>
        </div>
        <div title="Tag：目前與該點位連結的標籤 key。採集值會經由映射寫入此 Tag；若顯示 — 表示尚未建立映射。">
          <p className="text-slate-500">Tag</p>
          <p className="mt-0.5 truncate font-mono text-slate-200">{linkedTag?.key || '—'}</p>
        </div>
      </div>

      {/* 模式切換 */}
      <div className="grid grid-cols-2 gap-2">
        <button
          type="button"
          onClick={() => setTagLinkMode('existing')}
          title="從系統中已存在的 Tag 清單選擇一個 Tag，再按「套用 Tag 連結」建立或更新「目前點位 → 該 Tag」的映射。適用於多個點位共用同一個 Tag 或沿用既有 Tag 的情境。"
          className={`rounded-lg border px-3 py-2 text-xs font-medium focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 ${
            tagLinkMode === 'existing'
              ? 'border-blue-500/40 bg-blue-500/20 text-blue-100'
              : 'border-slate-700 bg-slate-800/70 text-slate-300'
          }`}
        >
          選擇既有 Tag
        </button>
        <button
          type="button"
          onClick={() => setTagLinkMode('create')}
          title="建立一個全新的 Tag（需輸入 Tag Key 與選填 Display Name），並同時建立「目前點位 → 新 Tag」的映射。新 Tag 的 data_type 會與目前點位一致。若 Key 已存在請改用「選擇既有 Tag」。"
          className={`rounded-lg border px-3 py-2 text-xs font-medium focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 ${
            tagLinkMode === 'create'
              ? 'border-indigo-500/40 bg-indigo-500/20 text-indigo-100'
              : 'border-slate-700 bg-slate-800/70 text-slate-300'
          }`}
        >
          新建 Tag
        </button>
      </div>

      {/* 既有 Tag 選擇 / 新建 Tag */}
      {tagLinkMode === 'existing' ? (
        <div className="space-y-2">
          <label className="block text-[11px] text-slate-300" title="列出系統中所有既有 Tag；選單顯示格式為 key (data_type)。選擇後需按「套用 Tag 連結」才會真正建立或更新映射。">
            已有 Tag
            <select
              value={selectedTagIdForLink}
              onChange={(e) => setSelectedTagIdForLink(e.target.value)}
              className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-blue-500"
              title="選擇要與目前點位建立映射的 Tag。若該點位已有映射，套用後會更新為此 Tag；若無則會建立新映射。"
              aria-label="選擇既有 Tag"
            >
              <option value="">請選擇 Tag</option>
              {tags.map((tag) => (
                <option key={tag.id} value={tag.id}>
                  {tag.key} ({tag.data_type})
                </option>
              ))}
            </select>
          </label>
          <button
            type="button"
            onClick={handleLinkTagToSelectedAddress}
            disabled={!activePointForLink || !selectedTagIdForLink || updateMappingPending || createMappingPending}
            title="將「目前選取的點位」與「上方選定的 Tag」建立或更新映射。若該點位尚無映射則建立一筆新映射；若已有映射則更新其 tag_id 為所選 Tag，採集值之後會寫入此 Tag。需先選取有點位的格位並選擇一個 Tag 才能按。"
            className="w-full rounded-lg border border-blue-500/40 bg-blue-500/20 px-3 py-2 text-xs font-medium text-blue-100 hover:bg-blue-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 disabled:cursor-not-allowed disabled:opacity-50"
          >
            套用 Tag 連結
          </button>
        </div>
      ) : (
        <div className="space-y-2">
          <label className="block text-[11px] text-slate-300" title="新 Tag 的唯一識別 key，通常使用英文、數字、底線（如 line_a_temp）。建立後無法在此處改名，請謹慎填寫。若 key 已存在會提示改用既有 Tag 模式。">
            Tag Key
            <input
              value={newTagKey}
              onChange={(e) => setNewTagKey(e.target.value)}
              placeholder="例如: line_a_temp"
              className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs font-mono focus:outline-none focus:ring-2 focus:ring-indigo-500"
              title="必填。新 Tag 的 key，系統內不可重複。"
              aria-label="新 Tag 的 Key"
            />
          </label>
          <label className="block text-[11px] text-slate-300" title="Tag 的顯示名稱，供報表或 UI 顯示用。可留白，留白時會以 Tag Key 作為顯示名稱。">
            Display Name
            <input
              value={newTagDisplayName}
              onChange={(e) => setNewTagDisplayName(e.target.value)}
              placeholder="例如: Line A Temperature"
              className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-indigo-500"
              title="選填。建立後仍可在「全域 Tag 編輯」中修改。"
              aria-label="新 Tag 的 Display Name"
            />
          </label>
          <button
            type="button"
            onClick={handleCreateTagAndLink}
            disabled={!activePointForLink || !newTagKey.trim() || createTagPending || createMappingPending || updateMappingPending}
            title="依上方輸入建立一個新 Tag（data_type 與目前點位相同），並同時建立「目前點位 → 新 Tag」的映射。完成後會自動切回「選擇既有 Tag」且選中剛建立的 Tag。"
            className="w-full rounded-lg border border-indigo-500/40 bg-indigo-500/20 px-3 py-2 text-xs font-medium text-indigo-100 hover:bg-indigo-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500 disabled:cursor-not-allowed disabled:opacity-50"
          >
            建立並連結 Tag
          </button>
        </div>
      )}

      {tagLinkActionMessage && (
        <p className="rounded border border-slate-700 bg-slate-900/70 px-2 py-1.5 text-[11px] text-slate-300">
          {tagLinkActionMessage}
        </p>
      )}

      {/* 全域 Tag 內嵌編輯 */}
      <div
        className="space-y-2 rounded-lg border border-white/10 bg-slate-900/60 p-3"
        title="當目前選取的格位已連結某個 Tag 時，可在此直接編輯該 Tag 的顯示名稱、單位、描述。變更會套用到「整個 Tag」，因此所有使用此 Tag 的映射都會一併受影響；儲存前會顯示預覽與影響的映射數量。"
      >
        <div className="flex items-center justify-between">
          <p
            className="text-[11px] font-semibold tracking-wide text-slate-200"
            title="編輯目前連結的 Tag 的全域屬性（Display Name、Unit、Description）。此為「Tag 本身」的編輯，不是只改單一映射；因此所有指向此 Tag 的點位都會看到更新後的顯示名稱、單位與描述。僅在已連結 Tag 時欄位才可編輯。"
          >
            全域 Tag 編輯
          </p>
          <span
            className="text-[10px] text-slate-400"
            title="目前連結的這個 Tag 被多少筆「點位→Tag」映射使用。若為 2 表示有 2 個點位連結到此 Tag；儲存此處的變更時，這 2 筆映射所顯示的 Tag 資訊都會一併更新。"
          >
            影響映射: {linkedTagAffectedMappingsCount}
          </span>
        </div>
        <label className="block text-[11px] text-slate-300" title="Tag 的顯示名稱，用於報表或介面顯示。修改後所有使用此 Tag 的映射都會顯示新名稱。">
          Display Name
          <input
            value={tagEditDisplayName}
            onChange={(e) => setTagEditDisplayName(e.target.value)}
            disabled={!linkedTag}
            className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-amber-500 disabled:opacity-50"
            title="目前連結 Tag 的顯示名稱。無連結 Tag 時無法編輯。"
            aria-label="Tag Display Name"
          />
        </label>
        <div className="grid grid-cols-2 gap-2">
          <label className="block text-[11px] text-slate-300" title="Tag 的單位（如 °C、kPa），供數值顯示或報表使用。">
            Unit
            <input
              value={tagEditUnit}
              onChange={(e) => setTagEditUnit(e.target.value)}
              disabled={!linkedTag}
              className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-amber-500 disabled:opacity-50"
              title="目前連結 Tag 的單位。"
              aria-label="Tag Unit"
            />
          </label>
          <label className="block text-[11px] text-slate-300" title="Tag 的說明文字，可記錄用途或備註。">
            Description
            <input
              value={tagEditDescription}
              onChange={(e) => setTagEditDescription(e.target.value)}
              disabled={!linkedTag}
              className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-amber-500 disabled:opacity-50"
              title="目前連結 Tag 的描述。"
              aria-label="Tag Description"
            />
          </label>
        </div>
        <button
          type="button"
          onClick={handleSaveLinkedTagEdit}
          disabled={!linkedTag || updateTagPending}
          title="將上方 Display Name / Unit / Description 的變更寫入目前連結的 Tag。按下後會先出現「變更預覽」與影響的映射數量，需再按「確認寫入」才會真正呼叫 API 更新；按「取消」則捨棄此次編輯。因會改動全域 Tag，所有使用此 Tag 的映射都會受影響。"
          className="w-full rounded-lg border border-amber-500/40 bg-amber-500/20 px-3 py-2 text-xs font-medium text-amber-100 hover:bg-amber-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500 disabled:cursor-not-allowed disabled:opacity-50"
        >
          儲存全域 Tag 變更
        </button>

        {/* 二次確認預覽 */}
        {pendingTagEdit && linkedTag && (
          <div
            className="space-y-2 rounded-lg border border-amber-500/30 bg-amber-500/10 p-2.5"
            title="按下「儲存全域 Tag 變更」後會顯示此預覽。請確認欄位前後值與影響的映射數量，再按「確認寫入」真正更新 Tag，或「取消」捨棄變更。"
          >
            <p className="text-[11px] text-amber-100">
              變更預覽：將影響 {linkedTagAffectedMappingsCount} 個映射
            </p>
            <div className="space-y-1 text-[11px] text-slate-200">
              {(
                [
                  ['display_name', linkedTag.display_name, pendingTagEdit.display_name],
                  ['unit',         linkedTag.unit,         pendingTagEdit.unit],
                  ['description',  linkedTag.description,  pendingTagEdit.description],
                ] as const
              ).map(([field, before, after]) => (
                <p key={field}>
                  {field}:{' '}
                  <span className="text-slate-400">{before || '—'}</span>
                  {' → '}
                  <span className="text-amber-100">{after || '—'}</span>
                </p>
              ))}
            </div>
            <div className="grid grid-cols-2 gap-2">
              <button
                type="button"
                onClick={handleConfirmTagEdit}
                disabled={updateTagPending}
                title="將預覽中的變更正式寫入 Tag；會呼叫 API 更新 Tag 的 display_name、unit、description，所有使用此 Tag 的映射都會看到新值。"
                className="rounded-md border border-amber-400/40 bg-amber-500/20 px-2 py-1.5 text-[11px] font-medium text-amber-100 hover:bg-amber-500/30 disabled:opacity-50"
              >
                確認寫入
              </button>
              <button
                type="button"
                onClick={clearPendingTagEdit}
                title="捨棄此次編輯，不寫入 Tag；預覽區會關閉，可再修改欄位後重新按「儲存全域 Tag 變更」。"
                className="rounded-md border border-slate-700 bg-slate-800/70 px-2 py-1.5 text-[11px] text-slate-200 hover:bg-slate-700"
              >
                取消
              </button>
            </div>
          </div>
        )}

        {tagEditMessage && (
          <p className="rounded border border-slate-700 bg-slate-900/70 px-2 py-1.5 text-[11px] text-slate-300">
            {tagEditMessage}
          </p>
        )}
      </div>
    </div>
  );
}
