interface SmartDashboardTagAndModbusPanelProps {
  selectedSourceAddress: string;
  activePointForLink: any;
  linkedTag: any;
  tagLinkMode: 'existing' | 'create';
  setTagLinkMode: (mode: 'existing' | 'create') => void;
  selectedTagIdForLink: string;
  setSelectedTagIdForLink: (tagId: string) => void;
  tags: any[];
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
  pendingTagEdit: any;
  handleConfirmTagEdit: () => void;
  clearPendingTagEdit: () => void;
  tagEditMessage: string;
  goToLocalModbusWorkbench: () => void;
  loadModbusStatus: () => void;
  modbusStatus: any;
  modbusRegister: string;
  setModbusRegister: (value: string) => void;
  handleBindTagToModbus: () => void;
  handlePushCurrentValueToModbus: () => void;
  handleSyncModbusFromMappings: () => void;
}

export default function SmartDashboardTagAndModbusPanel({
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
  goToLocalModbusWorkbench,
  loadModbusStatus,
  modbusStatus,
  modbusRegister,
  setModbusRegister,
  handleBindTagToModbus,
  handlePushCurrentValueToModbus,
  handleSyncModbusFromMappings,
}: SmartDashboardTagAndModbusPanelProps) {
  return (
    <>
      <div className="p-4 border-t border-white/5 space-y-3">
        <div className="flex items-center justify-between">
          <p className="text-xs font-semibold tracking-wide text-slate-200">Tag Linkage</p>
          <span className="text-[10px] text-slate-400 font-mono">
            {selectedSourceAddress || '-'}
          </span>
        </div>
        <div className="rounded-lg border border-white/10 bg-slate-900/60 p-3 text-[11px] text-slate-300 space-y-2">
          <p>選取格位: {selectedSourceAddress || '尚未選取'}</p>
          <p>點位: {activePointForLink?.name || '尚未建立點位'}</p>
          <p>目前 Tag: {linkedTag?.key || '未連結'}</p>
        </div>
        <div className="grid grid-cols-2 gap-2">
          <button
            type="button"
            onClick={() => setTagLinkMode('existing')}
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
            className={`rounded-lg border px-3 py-2 text-xs font-medium focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 ${
              tagLinkMode === 'create'
                ? 'border-indigo-500/40 bg-indigo-500/20 text-indigo-100'
                : 'border-slate-700 bg-slate-800/70 text-slate-300'
            }`}
          >
            新建 Tag
          </button>
        </div>
        {tagLinkMode === 'existing' ? (
          <div className="space-y-2">
            <label className="block text-[11px] text-slate-300">
              已有 Tag
              <select
                value={selectedTagIdForLink}
                onChange={(e) => setSelectedTagIdForLink(e.target.value)}
                className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-blue-500"
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
              className="w-full rounded-lg border border-blue-500/40 bg-blue-500/20 px-3 py-2 text-xs font-medium text-blue-100 hover:bg-blue-500/30 disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
            >
              套用 Tag 連結
            </button>
          </div>
        ) : (
          <div className="space-y-2">
            <label className="block text-[11px] text-slate-300">
              Tag Key
              <input
                value={newTagKey}
                onChange={(e) => setNewTagKey(e.target.value)}
                placeholder="例如: line_a_temp"
                className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs font-mono focus:outline-none focus:ring-2 focus:ring-indigo-500"
              />
            </label>
            <label className="block text-[11px] text-slate-300">
              Display Name
              <input
                value={newTagDisplayName}
                onChange={(e) => setNewTagDisplayName(e.target.value)}
                placeholder="例如: Line A Temperature"
                className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-indigo-500"
              />
            </label>
            <button
              type="button"
              onClick={handleCreateTagAndLink}
              disabled={!activePointForLink || !newTagKey.trim() || createTagPending || createMappingPending || updateMappingPending}
              className="w-full rounded-lg border border-indigo-500/40 bg-indigo-500/20 px-3 py-2 text-xs font-medium text-indigo-100 hover:bg-indigo-500/30 disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500"
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
        <div className="rounded-lg border border-white/10 bg-slate-900/60 p-3 space-y-2">
          <div className="flex items-center justify-between">
            <p className="text-[11px] font-semibold tracking-wide text-slate-200">全域 Tag 內嵌編輯</p>
            <span className="text-[10px] text-slate-400">
              影響映射: {linkedTagAffectedMappingsCount}
            </span>
          </div>
          <label className="block text-[11px] text-slate-300">
            Display Name
            <input
              value={tagEditDisplayName}
              onChange={(e) => setTagEditDisplayName(e.target.value)}
              disabled={!linkedTag}
              className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-amber-500 disabled:opacity-50"
            />
          </label>
          <label className="block text-[11px] text-slate-300">
            Unit
            <input
              value={tagEditUnit}
              onChange={(e) => setTagEditUnit(e.target.value)}
              disabled={!linkedTag}
              className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-amber-500 disabled:opacity-50"
            />
          </label>
          <label className="block text-[11px] text-slate-300">
            Description
            <input
              value={tagEditDescription}
              onChange={(e) => setTagEditDescription(e.target.value)}
              disabled={!linkedTag}
              className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-amber-500 disabled:opacity-50"
            />
          </label>
          <button
            type="button"
            onClick={handleSaveLinkedTagEdit}
            disabled={!linkedTag || updateTagPending}
            className="w-full rounded-lg border border-amber-500/40 bg-amber-500/20 px-3 py-2 text-xs font-medium text-amber-100 hover:bg-amber-500/30 disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500"
          >
            儲存全域 Tag 變更
          </button>
          {pendingTagEdit && linkedTag && (
            <div className="rounded-lg border border-amber-500/30 bg-amber-500/10 p-2.5 space-y-2">
              <p className="text-[11px] text-amber-100">
                變更預覽（第二次確認）: 將影響 {linkedTagAffectedMappingsCount} 個映射
              </p>
              <div className="space-y-1 text-[11px] text-slate-200">
                <p>
                  display_name: <span className="text-slate-400">{linkedTag.display_name || '-'}</span> →{' '}
                  <span className="text-amber-100">{pendingTagEdit.display_name || '-'}</span>
                </p>
                <p>
                  unit: <span className="text-slate-400">{linkedTag.unit || '-'}</span> →{' '}
                  <span className="text-amber-100">{pendingTagEdit.unit || '-'}</span>
                </p>
                <p>
                  description: <span className="text-slate-400">{linkedTag.description || '-'}</span> →{' '}
                  <span className="text-amber-100">{pendingTagEdit.description || '-'}</span>
                </p>
              </div>
              <div className="grid grid-cols-2 gap-2">
                <button
                  type="button"
                  onClick={handleConfirmTagEdit}
                  disabled={updateTagPending}
                  className="rounded-md border border-amber-400/40 bg-amber-500/20 px-2 py-1.5 text-[11px] font-medium text-amber-100 hover:bg-amber-500/30 disabled:opacity-50"
                >
                  確認寫入
                </button>
                <button
                  type="button"
                  onClick={clearPendingTagEdit}
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
      <div className="p-4 border-t border-white/5 space-y-3">
        <div className="flex items-center justify-between">
          <p className="text-xs font-semibold tracking-wide text-slate-200">Local Modbus Share</p>
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={goToLocalModbusWorkbench}
              className="rounded border border-blue-500/30 bg-blue-500/10 px-2 py-1 text-[10px] text-blue-200 hover:bg-blue-500/20 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
            >
              完整工作台
            </button>
            <button
              type="button"
              onClick={loadModbusStatus}
              className="rounded border border-slate-700 px-2 py-1 text-[10px] text-slate-300 hover:bg-slate-800/70 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
            >
              Refresh
            </button>
          </div>
        </div>
        <div className="rounded-lg border border-white/10 bg-slate-900/60 p-3 text-[11px] text-slate-300">
          <p>狀態: {modbusStatus?.enabled ? 'Running' : 'Stopped'}</p>
          <p>Address: {modbusStatus?.address || '-'} (Port 5020)</p>
          <p>Mappings: {modbusStatus?.mapping_count ?? 0}</p>
        </div>
        <label className="block text-[11px] text-slate-300">
          Register (Holding)
          <input
            value={modbusRegister}
            onChange={(e) => setModbusRegister(e.target.value)}
            className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs font-mono focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </label>
        <div className="grid grid-cols-1 gap-2">
          <button
            type="button"
            onClick={handleBindTagToModbus}
            className="w-full rounded-lg border border-indigo-500/40 bg-indigo-500/20 px-3 py-2 text-xs font-medium text-indigo-100 hover:bg-indigo-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500"
          >
            綁定目前 Tag 到 Register
          </button>
          <button
            type="button"
            onClick={handlePushCurrentValueToModbus}
            className="w-full rounded-lg border border-emerald-500/40 bg-emerald-500/20 px-3 py-2 text-xs font-medium text-emerald-100 hover:bg-emerald-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
          >
            推送目前值到 Modbus
          </button>
          <button
            type="button"
            onClick={handleSyncModbusFromMappings}
            className="w-full rounded-lg border border-blue-500/40 bg-blue-500/20 px-3 py-2 text-xs font-medium text-blue-100 hover:bg-blue-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
          >
            同步全部啟用映射
          </button>
        </div>
      </div>
    </>
  );
}
