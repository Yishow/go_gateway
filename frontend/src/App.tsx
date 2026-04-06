import { BrowserRouter, Routes, Route, Navigate, useLocation, useSearchParams } from 'react-router-dom'
import { CardMinimizeProvider } from './components/CardMinimizeProvider'
import { ThemeProvider } from './contexts/ThemeContext'
import { ToastProvider } from './contexts/ToastContext'
import { ErrorBoundary } from './components/ErrorBoundary'
import TestPage from './pages/TestPage'
import TestPageShell from './pages/TestPageShell'
import {
  buildWorkbenchRedirect,
  buildDashboardModalRedirect,
  buildLegacyMigrationRedirect,
  buildLocalModbusCompatRedirect,
} from './features/datalink/legacyRoutes'
import DatalinkWorkbenchPage from './pages/datalink/workbench/DatalinkWorkbenchPage'
import { GatewayCreateEntryRedirect, GatewayEntryRoute, GatewayExpertWorkbenchRoute, GatewayQuickSetupRoute } from './router/gateway'

function LocalModbusCompatRoute() {
  const [searchParams] = useSearchParams()
  return <Navigate to={buildLocalModbusCompatRedirect(searchParams.get('section'))} replace />
}

function LegacyStudioRedirect() {
  const location = useLocation()
  const destination = `${buildWorkbenchRedirect()}${location.search}${location.hash}`
  return <Navigate to={destination} replace />
}

function LegacyTestToolRedirect() {
  return <Navigate to="/test" replace />
}

/**
 * 主應用程式組件
 * 
 * 路由配置：
 * - / -> 重定向到 /studio
 * - /studio -> Datalink 主產品入口
 * - /datalink/* -> legacy 相容路由，統一收斂到 /studio
 * - /test -> 測試工具單頁入口
 */
function AppRoutes() {
  return (
    <Routes>
      {/* 首頁重定向到 studio */}
      <Route path="/" element={<Navigate to="/studio" replace />} />

      {/* Studio Main Route */}
      <Route path="/studio" element={<DatalinkWorkbenchPage />} />

      {/* Datalink legacy routes */}
      <Route path="/datalink" element={<LegacyStudioRedirect />} />
      <Route path="/datalink/workbench" element={<LegacyStudioRedirect />} />
      <Route path="/datalink/workbench/*" element={<LegacyStudioRedirect />} />
      <Route path="/datalink/dashboard-legacy" element={<LegacyStudioRedirect />} />
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
      <Route path="/datalink/local-modbus/legacy" element={<LocalModbusCompatRoute />} />

      {/* Gateway Dual Entry (W1 skeleton) */}
      <Route path="/gateway/entry" element={<GatewayEntryRoute />} />
      <Route path="/gateway/quick-setup" element={<GatewayQuickSetupRoute />} />
      <Route path="/gateway/expert-workbench" element={<GatewayExpertWorkbenchRoute />} />

      {/* Test UI */}
      <Route path="/test" element={<TestPageShell><TestPage /></TestPageShell>} />
      <Route path="/templates" element={<LegacyTestToolRedirect />} />
      <Route path="/history" element={<LegacyTestToolRedirect />} />
      <Route path="/compare" element={<LegacyTestToolRedirect />} />
      <Route path="/analyzer" element={<LegacyTestToolRedirect />} />

      {/* Unknown routes fallback */}
      <Route path="*" element={<Navigate to="/studio" replace />} />
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
    <ErrorBoundary>
      <ThemeProvider>
        <ToastProvider>
          <CardMinimizeProvider>
            <BrowserRouter>
              <AppRoutes />
            </BrowserRouter>
          </CardMinimizeProvider>
        </ToastProvider>
      </ThemeProvider>
    </ErrorBoundary>
  )
}

export default App
