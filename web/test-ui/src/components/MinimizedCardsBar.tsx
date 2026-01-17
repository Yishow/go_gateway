import { useState, useEffect } from 'react'
import { useCardMinimize, type CardType } from '../hooks/useCardMinimize'

/**
 * Card 類型對應的顏色配置
 */
const cardColors: Record<CardType, { bg: string; border: string; text: string; dot: string }> = {
  config: { 
    bg: 'bg-blue-50 dark:bg-blue-900/20', 
    border: 'border-blue-200 dark:border-blue-800', 
    text: 'text-blue-700 dark:text-blue-400', 
    dot: 'bg-blue-500' 
  },
  operations: { 
    bg: 'bg-amber-50 dark:bg-amber-900/20', 
    border: 'border-amber-200 dark:border-amber-800', 
    text: 'text-amber-700 dark:text-amber-400', 
    dot: 'bg-amber-500' 
  },
  monitor: { 
    bg: 'bg-purple-50 dark:bg-purple-900/20', 
    border: 'border-purple-200 dark:border-purple-800', 
    text: 'text-purple-700 dark:text-purple-400', 
    dot: 'bg-purple-500' 
  },
  debug: { 
    bg: 'bg-green-50 dark:bg-green-900/20', 
    border: 'border-green-200 dark:border-green-800', 
    text: 'text-green-700 dark:text-green-400', 
    dot: 'bg-green-500' 
  },
  scanner: { 
    bg: 'bg-indigo-50 dark:bg-indigo-900/20', 
    border: 'border-indigo-200 dark:border-indigo-800', 
    text: 'text-indigo-700 dark:text-indigo-400', 
    dot: 'bg-indigo-500' 
  },
}

/**
 * 狀態指示器的顏色
 */
const statusColors: Record<string, string> = {
  connected: 'bg-green-500',
  disconnected: 'bg-gray-400 dark:bg-gray-500',
  monitoring: 'bg-purple-500',
  idle: 'bg-gray-300 dark:bg-gray-600',
}

/**
 * 集中管理的 Card Bar 組件
 * 顯示所有被最小化的 card，並提供恢復功能
 */
export default function MinimizedCardsBar() {
  const { getAllMinimizedCards, restoreCard } = useCardMinimize()
  const minimizedCards = getAllMinimizedCards()
  const [isVisible, setIsVisible] = useState(false)

  // 在組件掛載後立即觸發顯示，確保定位正確
  useEffect(() => {
    if (minimizedCards.length > 0) {
      // 使用微任務確保在當前渲染完成後再顯示
      Promise.resolve().then(() => {
        setIsVisible(true)
      })
    }
  }, [minimizedCards.length])

  // 如果沒有最小化的 card，不顯示 bar
  if (minimizedCards.length === 0) {
    return null
  }

  /**
   * 處理 card 恢復
   */
  const handleRestore = (cardType: CardType) => {
    restoreCard(cardType)
  }

  return (
    <div className="fixed bottom-4 left-1/2 -translate-x-1/2 z-50">
      <div className={`bg-white dark:bg-gray-800 shadow-lg border border-gray-200 dark:border-gray-700 rounded-xl p-3 flex items-center gap-2 max-w-4xl overflow-x-auto ${isVisible ? 'animate-slide-up' : 'opacity-0'}`}>
        {minimizedCards.map((cardInfo) => {
          const colors = cardColors[cardInfo.type]
          const statusColor = cardInfo.status ? statusColors[cardInfo.status] : 'bg-gray-300 dark:bg-gray-600'

          return (
            <button
              key={cardInfo.type}
              onClick={() => handleRestore(cardInfo.type)}
              className={`
                flex items-center gap-2 px-4 py-2 rounded-lg border transition-all
                ${colors.bg} ${colors.border} ${colors.text}
                hover:shadow-md hover:scale-105 active:scale-95
                flex-shrink-0
              `}
              title={`點擊恢復 ${cardInfo.title}`}
            >
              {/* 狀態指示點 */}
              <div className={`w-2 h-2 rounded-full ${statusColor} ${cardInfo.status === 'monitoring' || cardInfo.status === 'connected' ? 'animate-pulse' : ''}`}></div>
              
              {/* Card 標題 */}
              <span className="font-semibold text-sm whitespace-nowrap">{cardInfo.title}</span>
              
              {/* 重要資訊摘要 */}
              {cardInfo.summary && (
                <span className="text-xs opacity-75 whitespace-nowrap truncate max-w-[200px]">
                  {cardInfo.summary}
                </span>
              )}
              
              {/* 恢復圖標 */}
              <svg
                className="w-4 h-4 opacity-60 hover:opacity-100 transition-opacity"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4"
                />
              </svg>
            </button>
          )
        })}
      </div>
    </div>
  )
}