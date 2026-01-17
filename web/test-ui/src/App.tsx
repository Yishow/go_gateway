import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { CardMinimizeProvider } from './components/CardMinimizeProvider'
import { ThemeProvider } from './contexts/ThemeContext'
import TestPage from './pages/TestPage'
import TemplatesPage from './pages/TemplatesPage'
import HistoryPage from './pages/HistoryPage'
import ComparePage from './pages/ComparePage'
import AnalyzerPage from './pages/AnalyzerPage'
import Layout from './components/Layout'

const queryClient = new QueryClient()

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <CardMinimizeProvider>
          <BrowserRouter>
            <Layout>
              <Routes>
                <Route path="/" element={<TestPage />} />
                <Route path="/templates" element={<TemplatesPage />} />
                <Route path="/history" element={<HistoryPage />} />
                <Route path="/compare" element={<ComparePage />} />
                <Route path="/analyzer" element={<AnalyzerPage />} />
              </Routes>
            </Layout>
          </BrowserRouter>
        </CardMinimizeProvider>
      </ThemeProvider>
    </QueryClientProvider>
  )
}

export default App
