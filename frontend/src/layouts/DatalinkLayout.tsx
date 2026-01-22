import { Outlet, NavLink, useLocation } from 'react-router-dom';
import { useTheme } from '../contexts/ThemeContext';

export default function DatalinkLayout() {
  const { theme } = useTheme();
  const location = useLocation();

  const isDarkMode = theme === 'dark';

  const navItems = [
    { path: '/datalink', label: 'Dashboard', icon: (
      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
      </svg>
    )},
    { path: '/datalink/devices', label: 'Devices', icon: (
      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
      </svg>
    )},
    { path: '/datalink/tags', label: 'Tags', icon: (
      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
      </svg>
    )},
    { path: '/datalink/mappings', label: 'Mappings', icon: (
      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
         <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
      </svg>
    )},
    { path: '/datalink/wizard', label: 'Mapping Wizard', icon: (
      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
      </svg>
    )},
    { path: '/datalink/settings', label: 'Settings', icon: (
      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
      </svg>
    )},
    { path: '/test', label: 'Protocol Test Tool', icon: (
      <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19.428 15.428a2 2 0 00-1.022-.547l-2.384-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z" />
      </svg>
    )},
  ];

  return (
    <div className={`min-h-screen flex ${isDarkMode ? 'bg-[#0f172a]' : 'bg-slate-50'}`}>
      {/* Sidebar - Glassmorphism & Gradient */}
      <aside className="w-64 flex-shrink-0 border-r border-slate-800 bg-slate-900/95 backdrop-blur-md relative overflow-hidden z-20">
         {/* Sidebar Gradient Glow */}
         <div className="absolute top-0 left-0 w-full h-96 bg-blue-900/20 blur-3xl pointer-events-none -translate-y-1/2"></div>
         
         <div className="p-6 relative z-10">
            <div className="flex items-center space-x-3 mb-8">
              <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center shadow-lg shadow-blue-500/20">
                  <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                  </svg>
              </div>
              <h1 className="text-xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-indigo-400 tracking-tight">
                Datalink
              </h1>
            </div>

            <nav className="space-y-1">
              {navItems.map(item => (
                <NavLink
                  key={item.path}
                  to={item.path}
                  end={item.path === '/datalink'} // Only match exact for Dashboard, prefix for others
                  className={({ isActive }) => `
                    flex items-center space-x-3 px-4 py-3 rounded-xl transition-all duration-200 group
                    ${isActive 
                      ? 'bg-blue-500/10 text-blue-400 shadow-sm border border-blue-500/20' 
                      : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50'
                    }
                  `}
                >
                  <span className={`transition-colors duration-200`}>
                    {item.icon}
                  </span>
                  <span className="font-medium">{item.label}</span>
                  
                  {/* Active Indicator Dot */}
                  <span className={`ml-auto w-1.5 h-1.5 rounded-full bg-blue-400 shadow-[0_0_8px_rgba(96,165,250,0.6)] opacity-0 transition-opacity duration-200 ${
                       location.pathname === item.path || (item.path !== '/datalink' && location.pathname.startsWith(item.path)) 
                       ? 'opacity-100' : ''
                  }`}></span>
                </NavLink>
              ))}
            </nav>
         </div>
         
         <div className="absolute bottom-6 left-6 right-6">
            <div className="bg-slate-800/50 backdrop-blur-sm rounded-xl p-4 border border-slate-700/50">
               <div className="flex items-center space-x-3">
                   <div className="w-8 h-8 rounded-full bg-slate-700 flex items-center justify-center text-xs font-bold text-slate-300">
                       Admin
                   </div>
                   <div>
                       <div className="text-xs font-medium text-slate-300">System Admin</div>
                       <div className="text-[10px] text-slate-500">v1.2.0-beta</div>
                   </div>
               </div>
            </div>
         </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 overflow-auto relative">
         {/* Top Header */}
         <header className="sticky top-0 z-10 bg-slate-900/80 backdrop-blur-lg border-b border-slate-800 px-8 py-4 flex justify-between items-center">
             <div>
                 {/* Breadcrumb-ish */}
                 <div className="text-sm text-slate-500 flex items-center gap-2">
                     <span>Datalink</span>
                     <span className="text-slate-700">/</span>
                     <span className="text-slate-200 capitalize">
                        {location.pathname === '/datalink' ? 'Dashboard' : location.pathname.split('/').pop()}
                     </span>
                 </div>
             </div>
             <div className="flex items-center space-x-4">
                 <button className="p-2 text-slate-400 hover:text-white transition-colors relative">
                    <span className="absolute top-2 right-2 w-2 h-2 bg-red-500 rounded-full border-2 border-slate-900"></span>
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                       <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                    </svg>
                 </button>
             </div>
         </header>

         <div className="p-8 max-w-7xl mx-auto">
           <Outlet />
         </div>
      </main>
    </div>
  );
}
