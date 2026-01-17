import { useEffect, useState } from 'react'

/**
 * Toast 通知類型
 */
export type ToastType = 'success' | 'error' | 'warning' | 'info'

/**
 * Toast 通知介面
 */
export interface Toast {
  id: string
  message: string
  type: ToastType
  duration?: number
}

/**
 * Toast 通知組件屬性
 */
interface ToastItemProps {
  toast: Toast
  onClose: (id: string) => void
}

/**
 * 單個 Toast 通知組件
 */
function ToastItem({ toast, onClose }: ToastItemProps) {
  const [progress, setProgress] = useState(100)
  const [isPaused, setIsPaused] = useState(false)
  const duration = toast.duration || 3000

  useEffect(() => {
    if (toast.duration === 0 || isPaused) return

    const startTime = Date.now()
    const interval = setInterval(() => {
      const elapsed = Date.now() - startTime
      const remaining = Math.max(0, duration - elapsed)
      const progressPercent = (remaining / duration) * 100
      setProgress(progressPercent)

      if (remaining <= 0) {
        clearInterval(interval)
        onClose(toast.id)
      }
    }, 50) // 每 50ms 更新一次進度

    return () => clearInterval(interval)
  }, [toast.id, toast.duration, onClose, duration, isPaused])

  const getIcon = () => {
    const iconClass = "w-5 h-5 flex-shrink-0"
    switch (toast.type) {
      case 'success':
        return (
          <div className="flex items-center justify-center w-8 h-8 rounded-full bg-green-100 dark:bg-green-900/30">
            <svg className={iconClass} fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
            </svg>
          </div>
        )
      case 'error':
        return (
          <div className="flex items-center justify-center w-8 h-8 rounded-full bg-red-100 dark:bg-red-900/30">
            <svg className={iconClass} fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </div>
        )
      case 'warning':
        return (
          <div className="flex items-center justify-center w-8 h-8 rounded-full bg-amber-100 dark:bg-amber-900/30">
            <svg className={iconClass} fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
        )
      case 'info':
        return (
          <div className="flex items-center justify-center w-8 h-8 rounded-full bg-blue-100 dark:bg-blue-900/30">
            <svg className={iconClass} fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        )
    }
  }

  const getStyles = () => {
    switch (toast.type) {
      case 'success':
        return {
          container: 'bg-white dark:bg-gray-800 border-green-200 dark:border-green-700/50 text-gray-900 dark:text-gray-100 shadow-xl shadow-green-500/10',
          text: 'text-gray-900 dark:text-gray-100',
          progress: 'bg-green-500'
        }
      case 'error':
        return {
          container: 'bg-white dark:bg-gray-800 border-red-200 dark:border-red-700/50 text-gray-900 dark:text-gray-100 shadow-xl shadow-red-500/10',
          text: 'text-gray-900 dark:text-gray-100',
          progress: 'bg-red-500'
        }
      case 'warning':
        return {
          container: 'bg-white dark:bg-gray-800 border-amber-200 dark:border-amber-700/50 text-gray-900 dark:text-gray-100 shadow-xl shadow-amber-500/10',
          text: 'text-gray-900 dark:text-gray-100',
          progress: 'bg-amber-500'
        }
      case 'info':
        return {
          container: 'bg-white dark:bg-gray-800 border-blue-200 dark:border-blue-700/50 text-gray-900 dark:text-gray-100 shadow-xl shadow-blue-500/10',
          text: 'text-gray-900 dark:text-gray-100',
          progress: 'bg-blue-500'
        }
    }
  }

  const styles = getStyles()

  return (
    <div
      className={`relative flex items-start gap-3 p-4 rounded-xl border-2 ${styles.container} backdrop-blur-sm`}
      role="alert"
      onMouseEnter={() => setIsPaused(true)}
      onMouseLeave={() => setIsPaused(false)}
    >
      {/* 進度條 */}
      {toast.duration !== 0 && (
        <div className="absolute top-0 left-0 right-0 h-1 bg-gray-200 dark:bg-gray-700 rounded-t-xl overflow-hidden">
          <div
            className={`h-full ${styles.progress} transition-all duration-50 ease-linear`}
            style={{ width: `${progress}%` }}
          />
        </div>
      )}
      
      {/* 圖標 */}
      <div className="flex-shrink-0 mt-0.5">
        {getIcon()}
      </div>
      
      {/* 內容 */}
      <div className="flex-1 min-w-0 pt-0.5">
        <p className={`text-sm font-medium leading-relaxed ${styles.text}`}>
          {toast.message}
        </p>
      </div>
      
      {/* 關閉按鈕 */}
      <button
        onClick={() => onClose(toast.id)}
        className="flex-shrink-0 text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300 transition-colors p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700"
        aria-label="關閉通知"
      >
        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  )
}

/**
 * Toast 容器組件屬性
 */
interface ToastContainerProps {
  toasts: Toast[]
  onClose: (id: string) => void
}

/**
 * Toast 容器組件
 */
export function ToastContainer({ toasts, onClose }: ToastContainerProps) {
  if (toasts.length === 0) return null

  return (
    <div className="fixed top-24 left-1/2 -translate-x-1/2 z-[9999] flex flex-col items-center gap-3 w-full max-w-md px-4 pointer-events-none">
      {toasts.map((toast, index) => (
        <div
          key={toast.id}
          className="pointer-events-auto w-full"
          style={{
            animation: `slide-in-down 0.3s cubic-bezier(0.16, 1, 0.3, 1) ${index * 50}ms forwards`,
          }}
        >
          <ToastItem toast={toast} onClose={onClose} />
        </div>
      ))}
    </div>
  )
}
