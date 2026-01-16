import { useState, useEffect, useCallback } from 'react';
import type { Profile } from '../types/profile';
import { PROFILE_STORAGE_KEY, DEFAULT_PROFILE_NAME } from '../types/profile';

/**
 * Profile 管理 Hook
 * 提供 Profile 的 CRUD 操作，並使用 localStorage 進行持久化
 */
export function useProfiles() {
  const [profiles, setProfiles] = useState<Profile[]>([]);
  const [currentProfileId, setCurrentProfileId] = useState<string | null>(null);

  /**
   * 從 localStorage 載入 Profiles
   */
  const loadProfiles = useCallback(() => {
    try {
      const stored = localStorage.getItem(PROFILE_STORAGE_KEY);
      if (stored) {
        const parsed = JSON.parse(stored) as Profile[];
        setProfiles(parsed);
        
        // 如果有當前選中的 profile，確保它存在
        const currentId = localStorage.getItem(`${PROFILE_STORAGE_KEY}_current`);
        if (currentId && parsed.find(p => p.id === currentId)) {
          setCurrentProfileId(currentId);
        } else if (parsed.length > 0) {
          // 如果當前 profile 不存在，選擇第一個
          setCurrentProfileId(parsed[0].id);
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
  }, []);

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
      config: {},
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
   */
  const deleteProfile = useCallback((id: string): void => {
    if (profiles.length <= 1) {
      alert('至少需要保留一個 Profile');
      return;
    }
    
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
   */
  const updateCurrentProfileConfig = useCallback((
    protocol: string,
    connectionMode: string,
    config: Record<string, any>
  ): void => {
    if (!currentProfileId) return;
    
    updateProfile(currentProfileId, {
      protocol,
      connectionMode,
      config,
    });
  }, [currentProfileId, updateProfile]);

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
