import { useState, useCallback, type ReactNode } from 'react'
import { CardMinimizeContext, type CardMinimizeContextType, type CardType, type MinimizedCardInfo } from '../hooks/useCardMinimize'

/**
 * Card 最小化 Provider 組件
 * 提供全局狀態管理來追蹤哪些 card 被最小化了
 */
export function CardMinimizeProvider({ children }: { children: ReactNode }) {
  const [minimizedCards, setMinimizedCards] = useState<Map<CardType, MinimizedCardInfo>>(new Map())

  /**
   * 最小化一個 card
   */
  const minimizeCard = useCallback((cardInfo: MinimizedCardInfo) => {
    setMinimizedCards(prev => {
      const newMap = new Map(prev)
      newMap.set(cardInfo.type, cardInfo)
      return newMap
    })
  }, [])

  /**
   * 恢復一個 card
   */
  const restoreCard = useCallback((cardType: CardType) => {
    setMinimizedCards(prev => {
      const newMap = new Map(prev)
      newMap.delete(cardType)
      return newMap
    })
  }, [])

  /**
   * 檢查 card 是否被最小化
   */
  const isMinimized = useCallback((cardType: CardType) => {
    return minimizedCards.has(cardType)
  }, [minimizedCards])

  /**
   * 獲取最小化的 card 資訊
   */
  const getMinimizedCardInfo = useCallback((cardType: CardType): MinimizedCardInfo | undefined => {
    return minimizedCards.get(cardType)
  }, [minimizedCards])

  /**
   * 獲取所有最小化的 card
   */
  const getAllMinimizedCards = useCallback((): MinimizedCardInfo[] => {
    return Array.from(minimizedCards.values())
  }, [minimizedCards])

  const value: CardMinimizeContextType = {
    minimizedCards,
    minimizeCard,
    restoreCard,
    isMinimized,
    getMinimizedCardInfo,
    getAllMinimizedCards,
  }

  return (
    <CardMinimizeContext.Provider value={value}>
      {children}
    </CardMinimizeContext.Provider>
  )
}
