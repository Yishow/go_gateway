import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { CardMinimizeProvider } from './components/CardMinimizeProvider'
import { ThemeProvider } from './contexts/ThemeContext'
import { ToastProvider } from './contexts/ToastContext'
import TestPage from './pages/TestPage'
import TemplatesPage from './pages/TemplatesPage'
import HistoryPage from './pages/HistoryPage'
import ComparePage from './pages/ComparePage'
import AnalyzerPage from './pages/AnalyzerPage'
import Layout from './components/Layout'
import DatalinkLayout from './layouts/DatalinkLayout'
import SmartDashboard from './pages/datalink/SmartDashboard'
import Dashboard from './pages/datalink/Dashboard'
import DevicesPage from './pages/datalink/DevicesPage'
import DeviceOnboardingPage from './pages/datalink/DeviceOnboardingPage'
import TagsPage from './pages/datalink/TagsPage'
import SettingsPage from './pages/datalink/SettingsPage'
import PollingGroupsPage from './pages/datalink/PollingGroupsPage'

/**
 * 主應用程式組件
 * 
 * 路由配置：
 * - / -> 重定向到 /datalink
 * - /datalink -> Datalink 主界面（Dashboard）
 * - /datalink/* -> Datalink 子路由
 * - /test, /templates, /history, /compare, /analyzer -> 舊版測試工具路由
 */
function AppRoutes() {
  return (
    <Routes>
      {/* 首頁重定向到 datalink */}
      <Route path="/" element={<Navigate to="/datalink" replace />} />

      {/* Datalink Routes - New Main UI */}
      <Route path="/datalink" element={<DatalinkLayout />}>
        <Route index element={<SmartDashboard />} />
        <Route path="dashboard-legacy" element={<Dashboard />} />
        <Route path="devices" element={<DevicesPage />} />
        <Route path="devices/new" element={<DeviceOnboardingPage />} />
        <Route path="points" element={<Navigate to="/datalink?legacy=points" replace />} />
        <Route path="polling-groups" element={<PollingGroupsPage />} />
        <Route path="tags" element={<TagsPage />} />
        <Route path="mappings" element={<Navigate to="/datalink?legacy=mappings" replace />} />
        <Route path="wizard" element={<Navigate to="/datalink?legacy=wizard" replace />} />
        <Route path="settings" element={<SettingsPage />} />
      </Route>

      {/* Legacy Test UI Routes */}
      <Route path="/test" element={<Layout><TestPage /></Layout>} />
      <Route path="/templates" element={<Layout><TemplatesPage /></Layout>} />
      <Route path="/history" element={<Layout><HistoryPage /></Layout>} />
      <Route path="/compare" element={<Layout><ComparePage /></Layout>} />
      <Route path="/analyzer" element={<Layout><AnalyzerPage /></Layout>} />

      {/* Unknown routes fallback */}
      <Route path="*" element={<Navigate to="/datalink" replace />} />
    </Routes>
  );
}

/**
 * 主應用程式組件
 * 
 * 根據 react-i18next 最佳實作：
 * 1. 使用 initReactI18next 初始化插件（在 i18n/config.ts 中）
 * 2. initReactI18next 會自動綁定 i18n 實例，不需要 I18nextProvider
 * 3. 設置 react.useSuspense: false 以兼容 React 19
 * 4. 確保 i18n 初始化完成後才渲染應用（在 main.tsx 中處理）
 */
function App() {
  return (
    <ThemeProvider>
      <ToastProvider>
        <CardMinimizeProvider>
          <BrowserRouter>
            <AppRoutes />
          </BrowserRouter>
        </CardMinimizeProvider>
      </ToastProvider>
    </ThemeProvider>
  )
}

export default App
