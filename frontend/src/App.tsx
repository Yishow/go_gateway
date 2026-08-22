import { lazy, Suspense } from 'react'
import { BrowserRouter, Routes, Route, Navigate, useLocation, useNavigate, useSearchParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { CardMinimizeProvider } from './components/CardMinimizeProvider'
import { ThemeProvider } from './contexts/ThemeContext'
import { ToastProvider } from './contexts/ToastContext'
import { ErrorBoundary } from './components/ErrorBoundary'
import {
  buildWorkbenchRedirect,
  buildWorkbenchV2EntryRedirect,
  buildDashboardModalRedirect,
  buildLegacyMigrationRedirect,
  buildLocalModbusCompatRedirect,
} from './features/datalink/legacyRoutes'

// 路由級代碼分割：各重型頁面改為 lazy chunk，縮小主 bundle 體積
const DatalinkWorkbenchPage = lazy(
  () => import('./pages/datalink/workbench/DatalinkWorkbenchPage'),
)
const TestPage = lazy(() => import('./pages/TestPage'))
const TestPageShell = lazy(() => import('./pages/TestPageShell'))
const DatalinkWorkbenchV2Page = lazy(
  () => import('./pages/datalink/workbench-v2/DatalinkWorkbenchV2Page'),
)
const RuntimeDashboardRoute = lazy(() =>
  import('./features/datalink/runtime-dashboard/RuntimeDashboardRoute').then((m) => ({
    default: m.RuntimeDashboardRoute,
  })),
)
const GatewayCreateEntryRedirect = lazy(() =>
  import('./router/gateway').then((m) => ({ default: m.GatewayCreateEntryRedirect })),
)
const GatewayEntryRoute = lazy(() =>
  import('./router/gateway').then((m) => ({ default: m.GatewayEntryRoute })),
)
const GatewayQuickSetupRoute = lazy(() =>
  import('./router/gateway').then((m) => ({ default: m.GatewayQuickSetupRoute })),
)
const GatewayExpertWorkbenchRoute = lazy(() =>
  import('./router/gateway').then((m) => ({ default: m.GatewayExpertWorkbenchRoute })),
)

function RouteFallback() {
  const { t } = useTranslation()
  return <div className="flex h-full min-h-screen items-center justify-center text-sm text-slate-400">{t('common.loading')}</div>
}

function LocalModbusCompatRoute() {
  const [searchParams] = useSearchParams()
  return <Navigate to={buildLocalModbusCompatRedirect(searchParams.get('section'))} replace />
}

function LegacyStudioRedirect() {
  const location = useLocation()
  const destination = `${buildWorkbenchRedirect()}${location.search}${location.hash}`
  return <Navigate to={destination} replace />
}

function GuidedWorkbenchEntryRedirect() {
  const location = useLocation()
  const destination = `${buildWorkbenchV2EntryRedirect()}${location.search}${location.hash}`
  return <Navigate to={destination} replace />
}

function LegacyTestToolRedirect() {
  return <Navigate to="/test" replace />
}

type RuntimeReturnFocus = {
  step: 1 | 2 | 3 | 4
  issueCode: string | null
}

function parseRuntimeReturnFocus(searchParams: URLSearchParams): RuntimeReturnFocus | null {
  if (searchParams.get('focus') !== 'readiness') {
    return null
  }

  const step = Number(searchParams.get('step'))
  if (step !== 1 && step !== 2 && step !== 3 && step !== 4) {
    return null
  }

  return {
    step,
    issueCode: searchParams.get('issue'),
  }
}

function DatalinkWorkbenchV2Route() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const runtimeReturnFocus = parseRuntimeReturnFocus(searchParams)

  return (
    <DatalinkWorkbenchV2Page
      navigateTo={navigate}
      runtimeReturnFocus={runtimeReturnFocus}
    />
  )
}

/**
 * 主應用程式組件
 * 
 * 路由配置：
 * - / -> 重定向到 /studio/v2
 * - /studio -> Datalink 主產品入口
 * - /datalink/* -> legacy 相容路由，generic landing 收斂到 /studio/v2
 * - /test -> 測試工具單頁入口
 */
export function AppRoutes() {
  return (
    <Suspense fallback={<RouteFallback />}>
      <Routes>
      {/* 首頁重定向到 studio/v2 */}
      <Route path="/" element={<GuidedWorkbenchEntryRedirect />} />

      {/* Studio Main Route */}
      <Route path="/studio/v2" element={<DatalinkWorkbenchV2Route />} />
      <Route path="/studio/runtime" element={<RuntimeDashboardRoute />} />
      <Route path="/studio" element={<DatalinkWorkbenchPage />} />

      {/* Datalink legacy routes */}
      <Route path="/datalink" element={<GuidedWorkbenchEntryRedirect />} />
      <Route path="/datalink/workbench" element={<GuidedWorkbenchEntryRedirect />} />
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
      <Route path="*" element={<GuidedWorkbenchEntryRedirect />} />
      </Routes>
    </Suspense>
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
