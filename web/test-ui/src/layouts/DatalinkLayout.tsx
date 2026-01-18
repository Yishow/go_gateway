import { Link, Outlet, useLocation } from 'react-router-dom';

const navItems = [
  { name: 'Dashboard', path: '/datalink' },
  { name: 'Devices', path: '/datalink/devices' },
  { name: 'Points', path: '/datalink/points' },
  { name: 'Tags', path: '/datalink/tags' },
  { name: 'Mappings', path: '/datalink/mappings' },
  { name: 'Settings', path: '/datalink/settings' },
];

export default function DatalinkLayout() {
  const location = useLocation();

  return (
    <div className="flex h-screen bg-slate-900 text-slate-100 font-sans">
      {/* Sidebar */}
      <aside className="w-64 flex-shrink-0 bg-slate-950 border-r border-slate-800">
        <div className="h-16 flex items-center px-6 border-b border-slate-800">
          <h1 className="text-xl font-bold bg-gradient-to-r from-blue-400 to-blue-600 bg-clip-text text-transparent">
            Datalink
          </h1>
        </div>
        <nav className="p-4 space-y-1">
          {navItems.map((item) => {
            const isActive =
              item.path === '/datalink'
                ? location.pathname === '/datalink'
                : location.pathname.startsWith(item.path);

            return (
              <Link
                key={item.path}
                to={item.path}
                className={`
                  block px-4 py-2 rounded-md text-sm font-medium transition-colors
                  ${
                    isActive
                      ? 'bg-blue-600/10 text-blue-400'
                      : 'text-slate-400 hover:text-slate-200 hover:bg-slate-900'
                  }
                `}
              >
                {item.name}
              </Link>
            );
          })}
        </nav>
        <div className="absolute bottom-0 w-64 p-4 border-t border-slate-800">
          <Link
            to="/"
            className="block text-center text-xs text-slate-500 hover:text-slate-300 transition-colors"
          >
            ← Back to Test UI
          </Link>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 flex flex-col min-w-0 overflow-hidden">
        {/* Top Header */}
        <header className="h-16 flex items-center justify-between px-6 bg-slate-900 border-b border-slate-800">
          <div className="flex items-center text-sm text-slate-400">
            <span className="font-medium text-slate-200">
              {navItems.find((i) =>
                i.path === '/datalink'
                  ? location.pathname === '/datalink'
                  : location.pathname.startsWith(i.path)
              )?.name || 'Datalink'}
            </span>
          </div>
          <div className="flex items-center space-x-4">
            <div className="text-xs text-slate-500">
              Go Gateway v0.1.0-alpha
            </div>
          </div>
        </header>

        {/* Scrollable Content */}
        <div className="flex-1 overflow-auto bg-slate-900 p-6">
          <div className="max-w-7xl mx-auto">
            <Outlet />
          </div>
        </div>
      </main>
    </div>
  );
}
