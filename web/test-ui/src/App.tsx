import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
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

const queryClient = new QueryClient()

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <ToastProvider>
          <CardMinimizeProvider>
            <BrowserRouter>
              <Routes>
                {/* Datalink Routes - New Main UI */}
                <Route path="/datalink" element={<DatalinkLayout />}>
                  <Route index element={<Dashboard />} />
                  <Route path="devices" element={<DevicesPage />} />
                  <Route path="points" element={<div className="text-white">Points Page (Coming Soon)</div>} />
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
                        <Route path="/" element={<TestPage />} />
                        <Route path="/templates" element={<TemplatesPage />} />
                        <Route path="/history" element={<HistoryPage />} />
                        <Route path="/compare" element={<ComparePage />} />
                        <Route path="/analyzer" element={<AnalyzerPage />} />
                      </Routes>
                    </Layout>
                  }
                />
              </Routes>
            </BrowserRouter>
          </CardMinimizeProvider>
        </ToastProvider>
      </ThemeProvider>
    </QueryClientProvider>
  )
}

export default App
