import { BrowserRouter, Routes, Route, Navigate, useSearchParams } from 'react-router-dom'
import { CardMinimizeProvider } from './components/CardMinimizeProvider'
import { ThemeProvider } from './contexts/ThemeContext'
import { ToastProvider } from './contexts/ToastContext'
import TestPage from './pages/TestPage'
import TemplatesPage from './pages/TemplatesPage'
import HistoryPage from './pages/HistoryPage'
import ComparePage from './pages/ComparePage'
import AnalyzerPage from './pages/AnalyzerPage'
import Layout from './components/Layout'
import SmartDashboard from './pages/datalink/SmartDashboard'
import {
  buildDashboardModalRedirect,
  buildLegacyMigrationRedirect,
  buildLocalModbusCompatRedirect,
} from './features/datalink/legacyRoutes'
import LocalModbusWorkbenchPage from './pages/datalink/LocalModbusWorkbenchPage'
import DatalinkWorkbenchPage from './pages/datalink/workbench/DatalinkWorkbenchPage'
import { GatewayCreateEntryRedirect, GatewayEntryRoute, GatewayExpertWorkbenchRoute, GatewayQuickSetupRoute } from './router/gateway'

function LocalModbusCompatRoute() {
  const [searchParams] = useSearchParams()
  return <Navigate to={buildLocalModbusCompatRedirect(searchParams.get('section'))} replace />
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

      {/* Datalink Routes - Dashboard First */}
      <Route path="/datalink" element={<SmartDashboard />} />
      <Route path="/datalink/dashboard-legacy" element={<Navigate to="/datalink" replace />} />
      <Route path="/datalink/devices-legacy" element={<Navigate to={buildDashboardModalRedirect('devices')} replace />} />
      <Route path="/datalink/devices" element={<Navigate to={buildDashboardModalRedirect('devices')} replace />} />
      <Route path="/datalink/devices/new" element={<GatewayCreateEntryRedirect />} />
      <Route path="/datalink/points" element={<Navigate to={buildLegacyMigrationRedirect('points')} replace />} />
      <Route path="/datalink/polling-groups" element={<Navigate to={buildDashboardModalRedirect('polling-groups')} replace />} />
      <Route path="/datalink/tags" element={<Navigate to={buildDashboardModalRedirect('tags')} replace />} />
      <Route path="/datalink/mappings" element={<Navigate to={buildLegacyMigrationRedirect('mappings')} replace />} />
      <Route path="/datalink/wizard" element={<Navigate to={buildLegacyMigrationRedirect('wizard')} replace />} />
      <Route path="/datalink/settings" element={<Navigate to="/datalink?modal=settings&section=settings" replace />} />
      <Route path="/datalink/local-modbus" element={<LocalModbusCompatRoute />} />
      <Route path="/datalink/local-modbus/legacy" element={<LocalModbusWorkbenchPage />} />
      <Route path="/datalink/workbench" element={<DatalinkWorkbenchPage />} />

      {/* Gateway Dual Entry (W1 skeleton) */}
      <Route path="/gateway/entry" element={<GatewayEntryRoute />} />
      <Route path="/gateway/quick-setup" element={<GatewayQuickSetupRoute />} />
      <Route path="/gateway/expert-workbench" element={<GatewayExpertWorkbenchRoute />} />

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
