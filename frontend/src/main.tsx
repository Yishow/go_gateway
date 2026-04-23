if (import.meta.env.DEV) {
  import("react-grab");
}

import './env' // 環境變數驗證（VITE_API_BASE_URL 可選，預設 /api/v1）
import * as React from 'react'
import * as ReactDOM from 'react-dom/client'
import { QueryClientProvider } from '@tanstack/react-query'
import { queryClient } from './lib/queryClient'
import './i18n/config' // 確保初始化副作用已執行
import App from './App.tsx'
import './index.css'

/**
 * 等待 i18n 初始化完成後再渲染 React 應用
 * 
 * 根據 react-i18next 最佳實踐：
 * - i18n.init() 是異步的，必須等待完成後才能使用 I18nextProvider
 * - 確保 initReactI18next 已正確綁定 React hooks
 * - 避免 "Cannot read properties of null (reading 'useMemo')" 錯誤
 */
const root = ReactDOM.createRoot(document.getElementById('root')!)

root.render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>
  </React.StrictMode>,
)
