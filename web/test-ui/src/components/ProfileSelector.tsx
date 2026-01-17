import { useState } from 'react';
import { useProfiles } from '../hooks/useProfiles';
import type { Profile } from '../types/profile';
import { useToast } from '../contexts/ToastContext';

interface ProfileSelectorProps {
  /** 當前選中的通訊協定 */
  currentProtocol: string;
  /** 當前選中的連線模式 */
  currentConnectionMode: string;
  /** 當前配置 */
  currentConfig: Record<string, any>;
  /** 當 Profile 切換時的回調 */
  onProfileChange: (profile: Profile) => void;
}

/**
 * Profile 選擇器組件
 * 提供 Profile 的選擇、新增、刪除、重命名功能
 */
export default function ProfileSelector({
  currentProtocol,
  currentConnectionMode,
  currentConfig,
  onProfileChange,
}: ProfileSelectorProps) {
  const {
    profiles,
    currentProfileId,
    createProfile,
    updateProfile,
    deleteProfile,
    switchProfile,
    updateCurrentProfileConfig,
  } = useProfiles();
  const { showError } = useToast();

  const [isAdding, setIsAdding] = useState(false);
  const [newProfileName, setNewProfileName] = useState('');
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editingName, setEditingName] = useState('');
  const [deleteConfirmId, setDeleteConfirmId] = useState<string | null>(null);

  /**
   * 處理 Profile 切換
   */
  const handleSwitchProfile = (profileId: string) => {
    // 先保存當前配置到當前 profile
    if (currentProfileId) {
      updateCurrentProfileConfig(currentProtocol, currentConnectionMode, currentConfig);
    }
    
    // 切換到新 profile
    switchProfile(profileId);
    const profile = profiles.find(p => p.id === profileId);
    if (profile) {
      onProfileChange(profile);
    }
  };

  /**
   * 處理新增 Profile
   */
  const handleAddProfile = () => {
    if (!newProfileName.trim()) {
      showError('請輸入 Profile 名稱');
      return;
    }

    // 檢查名稱是否重複
    if (profiles.some(p => p.name === newProfileName.trim())) {
      showError('Profile 名稱已存在');
      return;
    }

    // 先保存當前配置
    if (currentProfileId) {
      updateCurrentProfileConfig(currentProtocol, currentConnectionMode, currentConfig);
    }

    // 創建新 profile
    const newProfile = createProfile(newProfileName.trim());
    onProfileChange(newProfile);
    setIsAdding(false);
    setNewProfileName('');
  };

  /**
   * 處理刪除 Profile
   */
  const handleDeleteProfile = (profileId: string) => {
    // 檢查是否至少保留一個 Profile
    if (profiles.length <= 1) {
      showError('至少需要保留一個 Profile');
      return;
    }
    
    setDeleteConfirmId(profileId);
  };

  /**
   * 確認刪除 Profile
   */
  const confirmDeleteProfile = (profileId: string) => {
    // 先保存當前配置
    if (currentProfileId) {
      updateCurrentProfileConfig(currentProtocol, currentConnectionMode, currentConfig);
    }

    deleteProfile(profileId);
    setDeleteConfirmId(null);
    
    // 如果刪除的是當前 profile，切換到新的當前 profile
    if (profileId === currentProfileId) {
      const remaining = profiles.filter(p => p.id !== profileId);
      if (remaining.length > 0) {
        handleSwitchProfile(remaining[0].id);
      }
    }
  };

  /**
   * 開始編輯 Profile 名稱
   */
  const handleStartEdit = (profile: Profile) => {
    setEditingId(profile.id);
    setEditingName(profile.name);
  };

  /**
   * 保存編輯的 Profile 名稱
   */
  const handleSaveEdit = (profileId: string) => {
    if (!editingName.trim()) {
      showError('Profile 名稱不能為空');
      return;
    }

    // 檢查名稱是否重複（排除自己）
    if (profiles.some(p => p.id !== profileId && p.name === editingName.trim())) {
      showError('Profile 名稱已存在');
      return;
    }

    updateProfile(profileId, { name: editingName.trim() });
    setEditingId(null);
    setEditingName('');
  };

  /**
   * 取消編輯
   */
  const handleCancelEdit = () => {
    setEditingId(null);
    setEditingName('');
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <label className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
          配置 Profile
        </label>
        {!isAdding && (
          <button
            onClick={() => setIsAdding(true)}
            className="px-3 py-1.5 text-xs font-medium text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/20 rounded-lg hover:bg-blue-100 dark:hover:bg-blue-900/30 transition-colors flex items-center gap-1.5"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            新增
          </button>
        )}
      </div>

      {/* 新增 Profile 輸入框 */}
      {isAdding && (
        <div className="flex items-center gap-2 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800 animate-fade-in">
          <input
            type="text"
            value={newProfileName}
            onChange={(e) => setNewProfileName(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleAddProfile();
              } else if (e.key === 'Escape') {
                setIsAdding(false);
                setNewProfileName('');
              }
            }}
            placeholder="輸入 Profile 名稱..."
            className="flex-1 px-3 py-1.5 text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-white border border-blue-300 dark:border-blue-700 rounded-md hover:bg-blue-50 dark:hover:bg-blue-900/10 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:bg-white dark:focus:bg-gray-800 transition-all"
            autoFocus
          />
          <button
            onClick={handleAddProfile}
            className="px-3 py-1.5 text-xs font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 transition-colors"
          >
            確定
          </button>
          <button
            onClick={() => {
              setIsAdding(false);
              setNewProfileName('');
            }}
            className="px-3 py-1.5 text-xs font-medium text-gray-600 dark:text-gray-300 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
          >
            取消
          </button>
        </div>
      )}

      {/* Profile 列表 */}
      <div className="space-y-2">
        {profiles.length === 0 ? (
          <div className="p-3 text-sm text-gray-500 dark:text-gray-400 text-center">
            載入中...
          </div>
        ) : (
          profiles.map((profile) => {
          const isCurrent = profile.id === currentProfileId;
          const isEditing = editingId === profile.id;

          return (
            <div
              key={profile.id}
              className={`
                flex items-center gap-2 p-3 rounded-lg border transition-all
                ${isCurrent
                  ? 'bg-blue-50 dark:bg-blue-900/20 border-blue-300 dark:border-blue-700 shadow-sm'
                  : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'
                }
              `}
            >
              {/* Profile 選擇按鈕 */}
              <button
                onClick={() => handleSwitchProfile(profile.id)}
                className={`
                  flex-1 text-left flex items-center gap-2
                  ${isCurrent ? 'text-blue-700 dark:text-blue-400 font-medium' : 'text-gray-700 dark:text-gray-300'}
                `}
              >
                <div className={`
                  w-2 h-2 rounded-full
                  ${isCurrent ? 'bg-blue-600' : 'bg-gray-300 dark:bg-gray-600'}
                `} />
                {isEditing ? (
                  <input
                    type="text"
                    value={editingName}
                    onChange={(e) => setEditingName(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') {
                        handleSaveEdit(profile.id);
                      } else if (e.key === 'Escape') {
                        handleCancelEdit();
                      }
                    }}
                    onClick={(e) => e.stopPropagation()}
                    className="flex-1 px-2 py-1 text-sm bg-white dark:bg-gray-800 text-gray-900 dark:text-white border border-blue-300 dark:border-blue-700 rounded hover:bg-blue-50 dark:hover:bg-blue-900/10 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:bg-white dark:focus:bg-gray-800 transition-all"
                    autoFocus
                  />
                ) : (
                  <span className="flex-1">{profile.name}</span>
                )}
              </button>

              {/* 操作按鈕 */}
              {!isEditing && (
                <div className="flex items-center gap-1">
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      handleStartEdit(profile);
                    }}
                    className="p-1.5 text-gray-400 dark:text-gray-500 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
                    title="重新命名"
                  >
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  {profiles.length > 1 && (
                    deleteConfirmId === profile.id ? (
                      <div className="flex items-center gap-1">
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            confirmDeleteProfile(profile.id);
                          }}
                          className="p-1.5 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors"
                          title="確認刪除"
                        >
                          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                          </svg>
                        </button>
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            setDeleteConfirmId(null);
                          }}
                          className="p-1.5 text-gray-400 dark:text-gray-500 hover:bg-gray-50 dark:hover:bg-gray-800 rounded transition-colors"
                          title="取消"
                        >
                          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                          </svg>
                        </button>
                      </div>
                    ) : (
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleDeleteProfile(profile.id);
                        }}
                        className="p-1.5 text-gray-400 dark:text-gray-500 hover:text-red-600 dark:hover:text-red-400 transition-colors"
                        title="刪除"
                      >
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                      </button>
                    )
                  )}
                </div>
              )}

              {/* 編輯模式的操作按鈕 */}
              {isEditing && (
                <div className="flex items-center gap-1">
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      handleSaveEdit(profile.id);
                    }}
                    className="p-1.5 text-green-600 dark:text-green-400 hover:text-green-700 dark:hover:text-green-300 transition-colors"
                    title="儲存"
                  >
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                  </button>
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      handleCancelEdit();
                    }}
                    className="p-1.5 text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
                    title="取消"
                  >
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
              )}
            </div>
          );
        })
        )}
      </div>
    </div>
  );
}