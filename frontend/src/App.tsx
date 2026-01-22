import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
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
import Dashboard from './pages/datalink/Dashboard'
import DevicesPage from './pages/datalink/DevicesPage'
import TagsPage from './pages/datalink/TagsPage'
import MappingsPage from './pages/datalink/MappingsPage'
import MappingWizardPage from './pages/datalink/MappingWizardPage'
import SettingsPage from './pages/datalink/SettingsPage'

/**
 * Points 頁面佔位組件（使用翻譯）
 */
function PointsPage() {
  const { t } = useTranslation();
  return <div className="text-white">{t('points.comingSoon')}</div>;
}

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
        <Route index element={<Dashboard />} />
        <Route path="devices" element={<DevicesPage />} />
        <Route path="points" element={<PointsPage />} />
        <Route path="tags" element={<TagsPage />} />
        <Route path="mappings" element={<MappingsPage />} />
        <Route path="wizard" element={<MappingWizardPage />} />
        <Route path="settings" element={<SettingsPage />} />
      </Route>

      {/* Legacy Test UI Routes */}
      <Route
        path="*"
        element={
          <Layout>
            <Routes>
              <Route path="/test" element={<TestPage />} />
              <Route path="/templates" element={<TemplatesPage />} />
              <Route path="/history" element={<HistoryPage />} />
              <Route path="/compare" element={<ComparePage />} />
              <Route path="/analyzer" element={<AnalyzerPage />} />
            </Routes>
          </Layout>
        }
      />
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
