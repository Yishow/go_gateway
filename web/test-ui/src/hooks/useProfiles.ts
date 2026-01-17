import { useState, useEffect, useCallback } from 'react';
import type { Profile, ConnectionModeConfigs } from '../types/profile';
import { PROFILE_STORAGE_KEY, DEFAULT_PROFILE_NAME } from '../types/profile';

/**
 * Profile 管理 Hook
 * 提供 Profile 的 CRUD 操作，並使用 localStorage 進行持久化
 */
export function useProfiles() {
  const [profiles, setProfiles] = useState<Profile[]>([]);
  const [currentProfileId, setCurrentProfileId] = useState<string | null>(null);

  /**
   * 保存 Profiles 到 localStorage
   */
  const saveProfiles = useCallback((profilesToSave: Profile[]) => {
    try {
      localStorage.setItem(PROFILE_STORAGE_KEY, JSON.stringify(profilesToSave));
    } catch (error) {
      console.error('保存 Profiles 失敗:', error);
    }
  }, []);

  /**
   * 遷移舊格式的 Profile 到新格式（向後兼容）
   * 如果 config 是舊格式（Record<string, any>），轉換為新格式（ConnectionModeConfigs）
   */
  const migrateProfileIfNeeded = useCallback((profile: any): Profile => {
    // 檢查是否為舊格式（config 不是 ConnectionModeConfigs 格式）
    if (profile.config && typeof profile.config === 'object') {
      const configKeys = Object.keys(profile.config);
      // 如果 config 有 tcp/udp/serial/rtu 鍵，則已經是 new format
      const isNewFormat = ['tcp', 'udp', 'serial', 'rtu'].some(key => configKeys.includes(key));
      
      if (!isNewFormat && configKeys.length > 0) {
        // 這是舊格式，需要遷移
        const oldConfig = profile.config;
        const connectionMode = profile.connectionMode || 'tcp';
        const newConfig: ConnectionModeConfigs = {
          [connectionMode]: oldConfig,
        };
        return {
          ...profile,
          config: newConfig,
        };
      }
    }
    // 如果 config 為空或已經是 new format，直接返回
    if (!profile.config || Object.keys(profile.config).length === 0) {
      return {
        ...profile,
        config: {},
      };
    }
    return profile as Profile;
  }, []);

  /**
   * 從 localStorage 載入 Profiles
   */
  const loadProfiles = useCallback(() => {
    try {
      const stored = localStorage.getItem(PROFILE_STORAGE_KEY);
      if (stored) {
        const parsed = JSON.parse(stored) as any[];
        // 遷移舊格式的 Profile（向後兼容）
        const migratedProfiles = parsed.map(p => migrateProfileIfNeeded(p));
        setProfiles(migratedProfiles);
        
        // 如果有當前選中的 profile，確保它存在
        const currentId = localStorage.getItem(`${PROFILE_STORAGE_KEY}_current`);
        if (currentId && migratedProfiles.find(p => p.id === currentId)) {
          setCurrentProfileId(currentId);
        } else if (migratedProfiles.length > 0) {
          // 如果當前 profile 不存在，選擇第一個
          setCurrentProfileId(migratedProfiles[0].id);
        }
        
        // 如果有遷移，保存遷移後的資料
        if (JSON.stringify(migratedProfiles) !== JSON.stringify(parsed)) {
          saveProfiles(migratedProfiles);
        }
      } else {
        // 如果沒有存儲的 profiles，創建一個預設的
        const defaultProfile: Profile = {
          id: `profile_${Date.now()}`,
          name: DEFAULT_PROFILE_NAME,
          protocol: 'modbus_tcp',
          connectionMode: 'tcp',
          config: {},
          createdAt: Date.now(),
          updatedAt: Date.now(),
        };
        setProfiles([defaultProfile]);
        setCurrentProfileId(defaultProfile.id);
        saveProfiles([defaultProfile]);
      }
    } catch (error) {
      console.error('載入 Profiles 失敗:', error);
      // 發生錯誤時創建預設 profile
      const defaultProfile: Profile = {
        id: `profile_${Date.now()}`,
        name: DEFAULT_PROFILE_NAME,
        protocol: 'modbus_tcp',
        connectionMode: 'tcp',
        config: {},
        createdAt: Date.now(),
        updatedAt: Date.now(),
      };
      setProfiles([defaultProfile]);
      setCurrentProfileId(defaultProfile.id);
    }
  }, [migrateProfileIfNeeded, saveProfiles]);

  /**
   * 初始化載入
   */
  useEffect(() => {
    loadProfiles();
  }, [loadProfiles]);

  /**
   * 獲取當前選中的 Profile
   */
  const getCurrentProfile = useCallback((): Profile | null => {
    if (!currentProfileId) return null;
    return profiles.find(p => p.id === currentProfileId) || null;
  }, [profiles, currentProfileId]);

  /**
   * 創建新的 Profile
   */
  const createProfile = useCallback((name: string): Profile => {
    const newProfile: Profile = {
      id: `profile_${Date.now()}`,
      name,
      protocol: 'modbus_tcp',
      connectionMode: 'tcp',
      config: {} as ConnectionModeConfigs,
      createdAt: Date.now(),
      updatedAt: Date.now(),
    };
    
    const updated = [...profiles, newProfile];
    setProfiles(updated);
    saveProfiles(updated);
    setCurrentProfileId(newProfile.id);
    localStorage.setItem(`${PROFILE_STORAGE_KEY}_current`, newProfile.id);
    
    return newProfile;
  }, [profiles, saveProfiles]);

  /**
   * 更新 Profile
   */
  const updateProfile = useCallback((
    id: string,
    updates: Partial<Omit<Profile, 'id' | 'createdAt'>>
  ): void => {
    const updated = profiles.map(p => {
      if (p.id === id) {
        return {
          ...p,
          ...updates,
          updatedAt: Date.now(),
        };
      }
      return p;
    });
    
    setProfiles(updated);
    saveProfiles(updated);
  }, [profiles, saveProfiles]);

  /**
   * 刪除 Profile
   * 注意：調用者應該先檢查 profiles.length > 1
   */
  const deleteProfile = useCallback((id: string): void => {
    const updated = profiles.filter(p => p.id !== id);
    setProfiles(updated);
    saveProfiles(updated);
    
    // 如果刪除的是當前 profile，切換到第一個
    if (currentProfileId === id) {
      const newCurrentId = updated.length > 0 ? updated[0].id : null;
      setCurrentProfileId(newCurrentId);
      if (newCurrentId) {
        localStorage.setItem(`${PROFILE_STORAGE_KEY}_current`, newCurrentId);
      } else {
        localStorage.removeItem(`${PROFILE_STORAGE_KEY}_current`);
      }
    }
  }, [profiles, currentProfileId, saveProfiles]);

  /**
   * 切換到指定的 Profile
   */
  const switchProfile = useCallback((id: string): void => {
    if (profiles.find(p => p.id === id)) {
      setCurrentProfileId(id);
      localStorage.setItem(`${PROFILE_STORAGE_KEY}_current`, id);
    }
  }, [profiles]);

  /**
   * 更新當前 Profile 的配置
   * 按連線模式分開儲存配置
   */
  const updateCurrentProfileConfig = useCallback((
    protocol: string,
    connectionMode: string,
    config: Record<string, any>
  ): void => {
    if (!currentProfileId) return;
    
    const currentProfile = profiles.find(p => p.id === currentProfileId);
    if (!currentProfile) return;
    
    // 獲取現有的配置，按模式分開儲存
    const existingConfig: ConnectionModeConfigs = currentProfile.config || {};
    
    // 更新當前模式的配置
    const updatedConfig: ConnectionModeConfigs = {
      ...existingConfig,
      [connectionMode]: config,
    };
    
    updateProfile(currentProfileId, {
      protocol,
      connectionMode,
      config: updatedConfig,
    });
  }, [currentProfileId, profiles, updateProfile]);

  return {
    profiles,
    currentProfileId,
    currentProfile: getCurrentProfile(),
    createProfile,
    updateProfile,
    deleteProfile,
    switchProfile,
    updateCurrentProfileConfig,
    reload: loadProfiles,
  };
}
