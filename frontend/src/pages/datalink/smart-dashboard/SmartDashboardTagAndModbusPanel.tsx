/**
 * 組合層：將 TagPanel 與 ModbusPanel 組合為單一面板。
 * SmartDashboardPage 仍透過此元件統一傳遞所有 props，
 * 業務拆分在 SmartDashboardTagPanel / SmartDashboardModbusPanel 各自維護。
 */
import SmartDashboardTagPanel from './SmartDashboardTagPanel';
import SmartDashboardModbusPanel from './SmartDashboardModbusPanel';
import type { SmartDashboardTagPanelProps, PendingTagEdit } from './SmartDashboardTagPanel';
import type { SmartDashboardModbusPanelProps } from './SmartDashboardModbusPanel';

export type { PendingTagEdit };

type SmartDashboardTagAndModbusPanelProps =
  SmartDashboardTagPanelProps & SmartDashboardModbusPanelProps;

export default function SmartDashboardTagAndModbusPanel(props: SmartDashboardTagAndModbusPanelProps) {
  const {
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
  } = props;

  return (
    <>
      <SmartDashboardTagPanel
        selectedSourceAddress={selectedSourceAddress}
        activePointForLink={activePointForLink}
        linkedTag={linkedTag}
        tagLinkMode={tagLinkMode}
        setTagLinkMode={setTagLinkMode}
        selectedTagIdForLink={selectedTagIdForLink}
        setSelectedTagIdForLink={setSelectedTagIdForLink}
        tags={tags}
        handleLinkTagToSelectedAddress={handleLinkTagToSelectedAddress}
        updateMappingPending={updateMappingPending}
        createMappingPending={createMappingPending}
        newTagKey={newTagKey}
        setNewTagKey={setNewTagKey}
        newTagDisplayName={newTagDisplayName}
        setNewTagDisplayName={setNewTagDisplayName}
        handleCreateTagAndLink={handleCreateTagAndLink}
        createTagPending={createTagPending}
        tagLinkActionMessage={tagLinkActionMessage}
        linkedTagAffectedMappingsCount={linkedTagAffectedMappingsCount}
        tagEditDisplayName={tagEditDisplayName}
        setTagEditDisplayName={setTagEditDisplayName}
        tagEditUnit={tagEditUnit}
        setTagEditUnit={setTagEditUnit}
        tagEditDescription={tagEditDescription}
        setTagEditDescription={setTagEditDescription}
        handleSaveLinkedTagEdit={handleSaveLinkedTagEdit}
        updateTagPending={updateTagPending}
        pendingTagEdit={pendingTagEdit}
        handleConfirmTagEdit={handleConfirmTagEdit}
        clearPendingTagEdit={clearPendingTagEdit}
        tagEditMessage={tagEditMessage}
      />
      <SmartDashboardModbusPanel
        goToLocalModbusWorkbench={goToLocalModbusWorkbench}
        loadModbusStatus={loadModbusStatus}
        modbusStatus={modbusStatus}
        modbusRegister={modbusRegister}
        setModbusRegister={setModbusRegister}
        handleBindTagToModbus={handleBindTagToModbus}
        handlePushCurrentValueToModbus={handlePushCurrentValueToModbus}
        handleSyncModbusFromMappings={handleSyncModbusFromMappings}
      />
    </>
  );
}
