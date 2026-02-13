import { useState } from 'react';
import type { ReactNode } from 'react';
import { Outlet, NavLink, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useTheme } from '../contexts/ThemeContext';
import LanguageSwitcher from '../components/LanguageSwitcher';

type NavItem = {
  path: string;
  label: string;
  icon: ReactNode;
};

const SidebarToggleIcon = ({ open }: { open: boolean }) => (
  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
    {open ? (
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
    ) : (
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
    )}
  </svg>
);

export default function DatalinkLayout() {
  const { theme } = useTheme();
  const { t } = useTranslation();
  const location = useLocation();
  const [mobileSidebarOpen, setMobileSidebarOpen] = useState(false);

  const isDarkMode = theme === 'dark';

  const navItems: NavItem[] = [
    {
      path: '/datalink',
      label: t('nav.dashboard'),
      icon: (
        <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
        </svg>
      ),
    },
    {
      path: '/datalink/devices',
      label: t('nav.devices'),
      icon: (
        <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
        </svg>
      ),
    },
    {
      path: '/datalink/points',
      label: t('nav.points'),
      icon: (
        <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
      ),
    },
    {
      path: '/datalink/polling-groups',
      label: t('nav.pollingGroups'),
      icon: (
        <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      ),
    },
    {
      path: '/datalink/tags',
      label: t('nav.tags'),
      icon: (
        <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
        </svg>
      ),
    },
    {
      path: '/datalink/mappings',
      label: t('nav.mappings'),
      icon: (
        <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
        </svg>
      ),
    },
    {
      path: '/datalink/wizard',
      label: t('nav.mappingWizard'),
      icon: (
        <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
        </svg>
      ),
    },
    {
      path: '/datalink/settings',
      label: t('nav.settings'),
      icon: (
        <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
      ),
    },
    {
      path: '/test',
      label: t('nav.protocolTestTool'),
      icon: (
        <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19.428 15.428a2 2 0 00-1.022-.547l-2.384-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z" />
        </svg>
      ),
    },
  ];

  const matchedNavItem = navItems
    .filter((item) => location.pathname === item.path || location.pathname.startsWith(`${item.path}/`))
    .sort((a, b) => b.path.length - a.path.length)[0];
  const currentLabel = matchedNavItem?.label || t('nav.dashboard');

  const sidebarBaseClass = isDarkMode
    ? 'border-slate-800 bg-slate-900/95 text-slate-200'
    : 'border-slate-200 bg-white/95 text-slate-800';

  return (
    <div className={`min-h-screen ${isDarkMode ? 'bg-[#0f172a]' : 'bg-slate-50'} lg:flex`}>
      {mobileSidebarOpen && (
        <button
          type="button"
          className="fixed inset-0 z-30 bg-slate-900/50 lg:hidden"
          aria-label={t('layout.closeSidebar')}
          onClick={() => setMobileSidebarOpen(false)}
        />
      )}

      <aside
        className={`
          fixed inset-y-0 left-0 z-40 w-64 border-r backdrop-blur-md overflow-hidden transition-transform duration-200
          ${sidebarBaseClass}
          ${mobileSidebarOpen ? 'translate-x-0' : '-translate-x-full'}
          lg:static lg:translate-x-0 lg:flex-shrink-0
        `}
      >
        <div className="absolute top-0 left-0 w-full h-96 bg-blue-900/20 blur-3xl pointer-events-none -translate-y-1/2" />
        <div className="p-6 relative z-10 flex h-full flex-col">
          <div className="flex items-center space-x-3 mb-8">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center shadow-lg shadow-blue-500/20">
              <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <h1 className={`text-xl font-bold tracking-tight ${isDarkMode ? 'text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-indigo-400' : 'text-slate-900'}`}>
              Datalink
            </h1>
          </div>

          <nav className="space-y-1 flex-1">
            {navItems.map((item) => (
              <NavLink
                key={item.path}
                to={item.path}
                end={item.path === '/datalink'}
                onClick={() => setMobileSidebarOpen(false)}
                className={({ isActive }) => `
                  flex items-center space-x-3 px-4 py-3 rounded-xl transition-all duration-200 group
                  focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2
                  ${isDarkMode ? 'focus-visible:ring-offset-slate-900' : 'focus-visible:ring-offset-white'}
                  ${isActive
                    ? (isDarkMode
                      ? 'bg-blue-500/10 text-blue-400 shadow-sm border border-blue-500/20'
                      : 'bg-blue-50 text-blue-700 shadow-sm border border-blue-200')
                    : (isDarkMode
                      ? 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50'
                      : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100')}
                `}
              >
                <span>{item.icon}</span>
                <span className="font-medium">{item.label}</span>
              </NavLink>
            ))}
          </nav>

          <div className={`rounded-xl p-4 border ${isDarkMode ? 'bg-slate-800/50 border-slate-700/50' : 'bg-slate-100 border-slate-200'}`}>
            <div className="flex items-center space-x-3">
              <div className={`w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold ${isDarkMode ? 'bg-slate-700 text-slate-300' : 'bg-slate-200 text-slate-700'}`}>
                {t('system.admin')}
              </div>
              <div>
                <div className={`text-xs font-medium ${isDarkMode ? 'text-slate-300' : 'text-slate-700'}`}>{t('system.systemAdmin')}</div>
                <div className={`text-[10px] ${isDarkMode ? 'text-slate-500' : 'text-slate-500'}`}>{t('system.version')}</div>
              </div>
            </div>
          </div>
        </div>
      </aside>

      <main className="min-w-0 flex-1 lg:ml-0">
        <header
          className={`
            sticky top-0 z-20 border-b px-4 py-3 sm:px-6 sm:py-4 flex items-center justify-between gap-3
            ${isDarkMode ? 'bg-slate-900/80 border-slate-800' : 'bg-white/90 border-slate-200'}
            backdrop-blur-lg
          `}
        >
          <div className="flex items-center gap-3 min-w-0">
            <button
              type="button"
              className={`lg:hidden p-2 rounded-lg transition-colors ${isDarkMode ? 'text-slate-300 hover:bg-slate-800' : 'text-slate-700 hover:bg-slate-100'}`}
              aria-label={mobileSidebarOpen ? t('layout.closeSidebar') : t('layout.openSidebar')}
              onClick={() => setMobileSidebarOpen((prev) => !prev)}
            >
              <SidebarToggleIcon open={mobileSidebarOpen} />
            </button>
            <div className={`text-sm flex items-center gap-2 ${isDarkMode ? 'text-slate-400' : 'text-slate-600'} min-w-0`}>
              <span>Datalink</span>
              <span>/</span>
              <span className={`truncate font-medium ${isDarkMode ? 'text-slate-100' : 'text-slate-900'}`}>{currentLabel}</span>
            </div>
          </div>

          <div className="flex items-center gap-2 sm:gap-4">
            <LanguageSwitcher />
            <button
              type="button"
              aria-label={t('layout.notifications')}
              className={`p-2 rounded-lg transition-colors relative ${isDarkMode ? 'text-slate-400 hover:text-white hover:bg-slate-800/60' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'}`}
            >
              <span className={`absolute top-1.5 right-1.5 w-2 h-2 bg-red-500 rounded-full border-2 ${isDarkMode ? 'border-slate-900' : 'border-white'}`} />
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
              </svg>
            </button>
          </div>
        </header>

        <div className="p-4 sm:p-6 lg:p-8 max-w-7xl mx-auto">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
