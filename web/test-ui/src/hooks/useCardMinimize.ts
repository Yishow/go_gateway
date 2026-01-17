import { createContext, useContext } from 'react'

/**
 * Card 類型定義
 */
export type CardType = 'config' | 'operations' | 'monitor' | 'debug' | 'scanner'

/**
 * 最小化 Card 的資訊介面
 */
export interface MinimizedCardInfo {
  /** Card 類型 */
  type: CardType
  /** Card 標題 */
  title: string
  /** Card 的重要資訊摘要 */
  summary: string
  /** Card 的狀態指示（如連線狀態） */
  status?: 'connected' | 'disconnected' | 'monitoring' | 'idle'
}

/**
 * Card 最小化 Context 的介面
 */
export interface CardMinimizeContextType {
  minimizedCards: Map<CardType, MinimizedCardInfo>
  minimizeCard: (cardInfo: MinimizedCardInfo) => void
  minimizeCardsBatch: (cards: MinimizedCardInfo[]) => void
  restoreCard: (cardType: CardType) => void
  isMinimized: (cardType: CardType) => boolean
  getMinimizedCardInfo: (cardType: CardType) => MinimizedCardInfo | undefined
  getAllMinimizedCards: () => MinimizedCardInfo[]
}

/**
 * Card 最小化 Context
 */
export const CardMinimizeContext = createContext<CardMinimizeContextType | undefined>(undefined)

/**
 * Card 最小化管理 Hook
 * 提供全局狀態管理來追蹤哪些 card 被最小化了
 */
export function useCardMinimize() {
  const context = useContext(CardMinimizeContext)
  if (context === undefined) {
    throw new Error('useCardMinimize must be used within a CardMinimizeProvider')
  }
  return context
}
