import { useState } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { ThemeToggle } from './ThemeToggle'

// Simple Icons
const Icons = {
  Home: () => <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" /></svg>,
  Test: () => <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19.428 15.428a2 2 0 00-1.022-.547l-2.384-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z" /></svg>,
  Template: () => <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7v8a2 2 0 002 2h6M8 7V5a2 2 0 012-2h4.586a1 1 0 01.707.293l4.414 4.414a1 1 0 01.293.707V15a2 2 0 01-2 2h-2M8 7H6a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2v-2" /></svg>,
  History: () => <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>,
  Compare: () => <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" /></svg>,
  Analyzer: () => <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>,
  Menu: () => <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" /></svg>,
  Close: () => <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>,
}

interface LayoutProps {
  children: React.ReactNode
}

export default function Layout({ children }: LayoutProps) {
  const location = useLocation()
  const [isSidebarOpen, setIsSidebarOpen] = useState(false)

  const navItems = [
    { path: '/datalink', label: '首頁', icon: Icons.Home },
    { path: '/test', label: '測試工具', icon: Icons.Test },
    { path: '/templates', label: '模板管理', icon: Icons.Template },
    { path: '/history', label: '測試歷史', icon: Icons.History },
    { path: '/compare', label: '對比模式', icon: Icons.Compare },
    { path: '/analyzer', label: '協議分析器', icon: Icons.Analyzer },
  ]

  return (
    <div className="min-h-screen flex overflow-hidden bg-gray-50 dark:bg-[#0f172a]">
      {/* Mobile Sidebar Overlay */}
      {isSidebarOpen && (
        <div 
          className="fixed inset-0 bg-gray-800 bg-opacity-50 z-40 lg:hidden"
          onClick={() => setIsSidebarOpen(false)}
        />
      )}

      {/* Sidebar - Consistent with DatalinkLayout */}
      <aside 
        className={`
          fixed top-0 left-0 z-50 w-64 border-r border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900/95 backdrop-blur-md relative overflow-hidden transform transition-transform duration-200 ease-in-out h-screen
          ${isSidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}
        `}
      >
        {/* Sidebar Gradient Glow */}
        <div className="absolute top-0 left-0 w-full h-96 bg-blue-900/20 blur-3xl pointer-events-none -translate-y-1/2 hidden dark:block"></div>
        
        <div className="h-full flex flex-col overflow-hidden relative z-10">
          {/* Logo */}
          <div className="h-16 flex-shrink-0 flex items-center justify-between px-6 border-b border-gray-100 dark:border-gray-800">
            <h1 className="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-indigo-600 dark:from-blue-400 dark:to-indigo-400 flex items-center gap-2">
              <span className="p-1 bg-blue-100 dark:bg-blue-900/30 rounded-lg text-blue-600 dark:text-blue-400">🚀</span>
              協議測試
            </h1>
            <ThemeToggle />
          </div>

          {/* Navigation */}
          <nav className="flex-1 min-h-0 px-4 py-6 space-y-1 overflow-y-auto">
            {navItems.map((item) => {
              // 特殊處理首頁：當路徑是 /datalink 或其子路由時，都視為 active
              const isActive = item.path === '/datalink' 
                ? location.pathname === '/datalink' || location.pathname.startsWith('/datalink/')
                : location.pathname === item.path || location.pathname.startsWith(`${item.path}/`)
              return (
                <Link
                  key={item.path}
                  to={item.path}
                  onClick={() => setIsSidebarOpen(false)}
                  className={`
                    flex items-center gap-3 px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 group relative
                    ${isActive 
                      ? 'bg-blue-50 text-blue-700 shadow-sm border border-blue-100 dark:bg-blue-500/10 dark:text-blue-400 dark:border-blue-500/20' 
                      : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-800/50 dark:hover:text-gray-200'
                    }
                  `}
                >
                  <item.icon />
                  {item.label}
                  {/* Active Indicator Dot (Optional, like DatalinkLayout) */}
                  {isActive && (
                    <span className="absolute right-3 w-1.5 h-1.5 rounded-full bg-blue-500 dark:bg-blue-400 shadow-[0_0_8px_rgba(59,130,246,0.6)]"></span>
                  )}
                </Link>
              )
            })}
          </nav>

          {/* Footer */}
          <div className="flex-shrink-0 p-4 border-t border-gray-100 dark:border-gray-800">
            <div className="flex items-center gap-3 px-4 py-2 text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-800/50 rounded-xl border border-gray-100 dark:border-gray-700/50">
              <div className="w-2 h-2 rounded-full bg-green-500 shadow-[0_0_4px_rgba(34,197,94,0.6)]"></div>
              <span>系統正常運行</span>
            </div>
          </div>
        </div>
      </aside>

      {/* Main Content */}
      <div className="flex-1 flex flex-col min-w-0 overflow-hidden lg:ml-64 bg-gray-50 dark:bg-[#0f172a] min-h-screen">
        {/* Mobile Header */}
        <header className="lg:hidden bg-white dark:bg-gray-900/80 backdrop-blur-lg shadow-sm border-b border-gray-200 dark:border-gray-800 h-16 flex items-center justify-between px-4 flex-shrink-0 relative z-30">
          <button
            onClick={() => setIsSidebarOpen(true)}
            className="p-2 rounded-md text-gray-500 hover:text-gray-600 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-300 focus:outline-none"
          >
            <Icons.Menu />
          </button>
          <span className="font-semibold text-gray-900 dark:text-white">協議測試工具</span>
          <ThemeToggle />
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-y-auto p-4 sm:p-6 lg:p-8 pb-20 lg:pb-32 min-h-full">
          <div className="max-w-7xl mx-auto">
            {children}
          </div>
        </main>
      </div>
    </div>
  )
}
