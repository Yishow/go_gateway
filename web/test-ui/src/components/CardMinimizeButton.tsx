import { useCardMinimize, type CardType, type MinimizedCardInfo } from '../hooks/useCardMinimize'

interface CardMinimizeButtonProps {
  /** Card 類型 */
  cardType: CardType
  /** Card 標題 */
  title: string
  /** Card 的重要資訊摘要 */
  summary: string
  /** Card 的狀態指示 */
  status?: 'connected' | 'disconnected' | 'monitoring' | 'idle'
}

/**
 * Card 最小化按鈕組件
 * 在 card 頂部中間顯示一個漂亮的 3 點按鈕，點擊後最小化 card
 */
export default function CardMinimizeButton({
  cardType,
  title,
  summary,
  status,
}: CardMinimizeButtonProps) {
  const { minimizeCard, isMinimized } = useCardMinimize()
  const minimized = isMinimized(cardType)

  // 如果已經被最小化，不顯示按鈕
  if (minimized) {
    return null
  }

  /**
   * 處理最小化
   */
  const handleMinimize = () => {
    const cardInfo: MinimizedCardInfo = {
      type: cardType,
      title,
      summary,
      status,
    }
    minimizeCard(cardInfo)
  }

  return (
    <button
      onClick={handleMinimize}
      className="
        absolute top-2 right-2 z-50
        p-1.5
        flex items-center justify-center
        transition-all duration-200
        hover:bg-gray-100/80 rounded
        active:scale-95
        group
        pointer-events-auto
      "
      title="最小化此卡片"
    >
      {/* 3 個漂亮的點 */}
      <div className="flex items-center gap-1">
        <span className="w-1.5 h-1.5 rounded-full bg-gray-400 group-hover:bg-gray-600 transition-colors"></span>
        <span className="w-1.5 h-1.5 rounded-full bg-gray-400 group-hover:bg-gray-600 transition-colors"></span>
        <span className="w-1.5 h-1.5 rounded-full bg-gray-400 group-hover:bg-gray-600 transition-colors"></span>
      </div>
    </button>
  )
}
