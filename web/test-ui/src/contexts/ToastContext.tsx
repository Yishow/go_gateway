import { createContext, useContext, useState, useCallback } from 'react'
import type { ReactNode } from 'react'
import { ToastContainer } from '../components/Toast'
import type { Toast, ToastType } from '../components/Toast'

/**
 * Toast Context 類型
 */
interface ToastContextType {
  showToast: (message: string, type?: ToastType, duration?: number) => void
  showSuccess: (message: string, duration?: number) => void
  showError: (message: string, duration?: number) => void
  showWarning: (message: string, duration?: number) => void
  showInfo: (message: string, duration?: number) => void
}

const ToastContext = createContext<ToastContextType | undefined>(undefined)

/**
 * Toast Provider 組件
 */
export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([])

  /**
   * 顯示 Toast 通知
   */
  const showToast = useCallback((message: string, type: ToastType = 'info', duration: number = 3000) => {
    const id = `toast-${Date.now()}-${Math.random()}`
    const newToast: Toast = {
      id,
      message,
      type,
      duration,
    }
    setToasts((prev) => [...prev, newToast])
  }, [])

  /**
   * 顯示成功通知
   */
  const showSuccess = useCallback((message: string, duration?: number) => {
    showToast(message, 'success', duration)
  }, [showToast])

  /**
   * 顯示錯誤通知
   */
  const showError = useCallback((message: string, duration?: number) => {
    showToast(message, 'error', duration || 5000) // 錯誤通知預設顯示 5 秒
  }, [showToast])

  /**
   * 顯示警告通知
   */
  const showWarning = useCallback((message: string, duration?: number) => {
    showToast(message, 'warning', duration)
  }, [showToast])

  /**
   * 顯示信息通知
   */
  const showInfo = useCallback((message: string, duration?: number) => {
    showToast(message, 'info', duration)
  }, [showToast])

  /**
   * 關閉 Toast 通知
   */
  const handleClose = useCallback((id: string) => {
    setToasts((prev) => prev.filter((toast) => toast.id !== id))
  }, [])

  return (
    <ToastContext.Provider
      value={{
        showToast,
        showSuccess,
        showError,
        showWarning,
        showInfo,
      }}
    >
      {children}
      <ToastContainer toasts={toasts} onClose={handleClose} />
    </ToastContext.Provider>
  )
}

/**
 * 使用 Toast Hook
 */
export function useToast() {
  const context = useContext(ToastContext)
  if (context === undefined) {
    throw new Error('useToast must be used within a ToastProvider')
  }
  return context
}
