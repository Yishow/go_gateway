/**
 * 將 hex 字串轉成 byte 陣列
 * @param hexData hex 字串（可包含空白分隔）
 */
export const parseHexBytes = (hexData: string): number[] => {
  const tokens = hexData.split(/\s+/).filter(token => token.length > 0)
  const bytes: number[] = []

  for (const token of tokens) {
    const value = Number.parseInt(token, 16)
    if (Number.isNaN(value)) {
      return []
    }
    bytes.push(value)
  }

  return bytes
}

/**
 * 正規化封包的 raw_data，確保可用 byte 陣列輸出
 * @param rawData 後端回傳的 raw_data（可能是陣列、類陣列、或 base64 字串）
 * @param hexData hex 字串（用於回退解析）
 */
export const normalizeRawDataFromPacket = (rawData: unknown, hexData: string): number[] => {
  if (Array.isArray(rawData)) {
    return rawData.map(value => Number(value))
  }

  if (rawData && typeof rawData === 'object' && 'length' in rawData) {
    return Array.from(rawData as ArrayLike<unknown>, value => Number(value))
  }

  if (typeof rawData === 'string' && rawData.length > 0) {
    try {
      const decoded: string = typeof atob === 'function'
        ? atob(rawData)
        : ((globalThis as any).Buffer?.from(rawData, 'base64')?.toString('binary') ?? '')
      const bytes = Array.from(decoded, (char: string) => char.charCodeAt(0))
      if (bytes.length > 0) {
        return bytes
      }
    } catch {
      // 無法解析 base64 時使用 hex_data
    }
  }

  return parseHexBytes(hexData)
}
