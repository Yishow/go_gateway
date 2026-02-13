import React, { useState, useEffect } from 'react';
import { tagAPI } from '../../../../services/datalink';
import type { Tag } from '../../../../types/datalink';
import { logger } from '../../../../utils/logger';
import './Steps.css';

/**
 * TagStep 組件屬性
 */
interface TagStepProps {
  /** 已選擇的標籤 ID */
  selectedTagId: string;
  /** 選擇回呼 */
  onSelect: (tagId: string) => void;
  /** 錯誤訊息 */
  error?: string;
}

/**
 * 標籤選擇步驟
 */
export const TagStep: React.FC<TagStepProps> = ({
  selectedTagId,
  onSelect,
  error,
}) => {
  const [tags, setTags] = useState<Tag[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    loadTags();
  }, []);

  const loadTags = async () => {
    try {
      setLoading(true);
      const data = await tagAPI.list();
      setTags(data);
    } catch (err) {
      logger.error('Failed to load tags:', err);
    } finally {
      setLoading(false);
    }
  };

  const filteredTags = tags.filter((tag) =>
    tag.key.toLowerCase().includes(searchQuery.toLowerCase()) ||
    tag.display_name?.toLowerCase().includes(searchQuery.toLowerCase())
  );

  if (loading) {
    return <div className="step-loading">載入中...</div>;
  }

  return (
    <div className="step-tag">
      <div className="step-search">
        <input
          type="text"
          placeholder="搜尋標籤..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="search-input"
        />
      </div>

      {error && <div className="step-error">{error}</div>}

      <div className="tag-list">
        {filteredTags.length === 0 ? (
          <div className="no-data">
            {searchQuery ? '沒有符合的標籤' : '沒有標籤，請先新增標籤'}
          </div>
        ) : (
          filteredTags.map((tag) => (
            <button
              type="button"
              key={tag.id}
              className={`tag-card ${selectedTagId === tag.id ? 'selected' : ''}`}
              onClick={() => onSelect(tag.id)}
              aria-pressed={selectedTagId === tag.id}
              aria-label={`選擇標籤 ${tag.display_name || tag.key}`}
            >
              <div className="tag-icon">🏷️</div>
              <div className="tag-info">
                <h4>{tag.display_name || tag.key}</h4>
                <p>鍵: {tag.key} | 類型: {tag.data_type}</p>
              </div>
              <div className={`tag-status status-${tag.status}`}>
                {tag.status}
              </div>
              {selectedTagId === tag.id && (
                <div className="selected-check">✓</div>
              )}
            </button>
          ))
        )}
      </div>
    </div>
  );
};

export default TagStep;
