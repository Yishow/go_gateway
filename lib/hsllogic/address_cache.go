package hsllogic

import (
	"fmt"
	"sync"
)

// =============================================================================
// AddressCache 位址快取
// =============================================================================

/**
 * AddressCache 位址解析快取，避免熱迴圈中重複解析相同的位址字串
 *
 * 使用場景：
 * - 高頻輪詢相同的標籤位址時，可避免重複解析
 * - 建議在 Connector 或 PollingGroup 層級使用
 *
 * 執行緒安全：使用 sync.Map 實作，支援並發讀寫
 */
type AddressCache struct {
	// cache 快取 map，key 為 "protocol:address"，value 為 *ParsedAddress
	cache sync.Map
	// hits 快取命中次數（用於監控）
	hits uint64
	// misses 快取未命中次數（用於監控）
	misses uint64
	// mu 保護統計計數器的互斥鎖
	mu sync.RWMutex
}

// =============================================================================
// 建構函數
// =============================================================================

/**
 * NewAddressCache 建立新的位址快取
 * @returns *AddressCache 位址快取實例
 */
func NewAddressCache() *AddressCache {
	return &AddressCache{}
}

// =============================================================================
// 公開方法
// =============================================================================

/**
 * Get 取得已解析的位址，若快取未命中則解析並快取
 * @param protocol 協議類型
 * @param address 原始位址字串 (如 "D100", "HR1000")
 * @returns *ParsedAddress 解析後的位址結構
 * @returns error 解析錯誤
 */
func (c *AddressCache) Get(protocol ProtocolType, address string) (*ParsedAddress, error) {
	// 建立快取鍵 (protocol:address)
	key := c.makeCacheKey(protocol, address)

	// 嘗試從快取取得
	if cached, ok := c.cache.Load(key); ok {
		c.recordHit()
		parsed, typeOK := cached.(*ParsedAddress)
		if !typeOK {
			return nil, fmt.Errorf("快取資料型別錯誤: %s", key)
		}
		return parsed, nil
	}

	// 快取未命中，進行解析
	c.recordMiss()
	parsed, err := ParseAddress(protocol, address)
	if err != nil {
		return nil, fmt.Errorf("解析位址失敗 [%s:%s]: %w", protocol, address, err)
	}

	// 快取解析結果
	c.cache.Store(key, parsed)
	return parsed, nil
}

/**
 * Clear 清空快取
 */
func (c *AddressCache) Clear() {
	c.cache.Range(func(key, value interface{}) bool {
		c.cache.Delete(key)
		return true
	})

	c.mu.Lock()
	c.hits = 0
	c.misses = 0
	c.mu.Unlock()
}

/**
 * Stats 取得快取統計資訊
 * @returns hits 快取命中次數
 * @returns misses 快取未命中次數
 * @returns hitRate 命中率 (0.0 ~ 1.0)
 */
func (c *AddressCache) Stats() (hits, misses uint64, hitRate float64) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	hits = c.hits
	misses = c.misses
	total := hits + misses
	if total > 0 {
		hitRate = float64(hits) / float64(total)
	}
	return
}

/**
 * Size 取得快取中的項目數量
 * @returns int 快取項目數量
 */
func (c *AddressCache) Size() int {
	count := 0
	c.cache.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

// =============================================================================
// 內部方法
// =============================================================================

/**
 * makeCacheKey 產生快取鍵
 * @param protocol 協議類型
 * @param address 位址字串
 * @returns string 快取鍵 (格式: "protocol:address")
 */
func (c *AddressCache) makeCacheKey(protocol ProtocolType, address string) string {
	return string(protocol) + ":" + address
}

/**
 * recordHit 記錄快取命中
 */
func (c *AddressCache) recordHit() {
	c.mu.Lock()
	c.hits++
	c.mu.Unlock()
}

/**
 * recordMiss 記錄快取未命中
 */
func (c *AddressCache) recordMiss() {
	c.mu.Lock()
	c.misses++
	c.mu.Unlock()
}

// =============================================================================
// 全域快取實例（可選使用）
// =============================================================================

var (
	globalAddressCache     *AddressCache
	globalAddressCacheOnce sync.Once
)

/**
 * GetGlobalAddressCache 取得全域位址快取
 * @returns *AddressCache 全域快取實例
 */
func GetGlobalAddressCache() *AddressCache {
	globalAddressCacheOnce.Do(func() {
		globalAddressCache = NewAddressCache()
	})
	return globalAddressCache
}

/**
 * SetGlobalAddressCache 設定全域位址快取
 * @param cache 快取實例
 */
func SetGlobalAddressCache(cache *AddressCache) {
	globalAddressCache = cache
}
