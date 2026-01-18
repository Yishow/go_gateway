import { useNavigate } from 'react-router-dom';

export default function Dashboard() {
  const navigate = useNavigate();

  const cards = [
    {
      title: 'Devices',
      description: 'Manage industrial devices and PLCs',
      path: '/datalink/devices',
      color: 'bg-blue-600',
      icon: (
        <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 01-2 2v4a2 2 0 012 2h14a2 2 0 012-2v-4a2 2 0 01-2-2m-2-4h.01M17 16h.01" />
        </svg>
      ),
    },
    {
      title: 'Points',
      description: 'Configure data points and polling',
      path: '/datalink/points',
      color: 'bg-indigo-600',
      icon: (
        <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
      ),
    },
    {
      title: 'Tags',
      description: 'Define global tag dictionary',
      path: '/datalink/tags',
      color: 'bg-emerald-600',
      icon: (
        <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
        </svg>
      ),
    },
    {
      title: 'Mappings',
      description: 'Design data transformation pipelines',
      path: '/datalink/mappings',
      color: 'bg-purple-600',
      icon: (
        <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
        </svg>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-slate-100">Overview</h2>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {cards.map((card) => (
          <div
            key={card.path}
            onClick={() => navigate(card.path)}
            className="group relative bg-slate-800 rounded-xl p-6 cursor-pointer border border-slate-700 hover:border-blue-500/50 transition-all duration-200 hover:shadow-lg hover:shadow-blue-500/10"
          >
            <div className={`inline-flex items-center justify-center p-3 rounded-lg ${card.color} text-white mb-4`}>
              {card.icon}
            </div>
            <h3 className="text-lg font-semibold text-slate-100 mb-2 group-hover:text-blue-400 transition-colors">
              {card.title}
            </h3>
            <p className="text-sm text-slate-400">
              {card.description}
            </p>
          </div>
        ))}
      </div>

      <div className="mt-8 bg-slate-800 rounded-xl border border-slate-700 p-6">
        <h3 className="text-lg font-semibold text-slate-100 mb-4">System Status</h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="p-4 bg-slate-900/50 rounded-lg border border-slate-800">
            <div className="text-sm text-slate-400 mb-1">Active Connections</div>
            <div className="text-2xl font-bold text-emerald-400">0</div>
          </div>
          <div className="p-4 bg-slate-900/50 rounded-lg border border-slate-800">
            <div className="text-sm text-slate-400 mb-1">Points Polled / Sec</div>
            <div className="text-2xl font-bold text-blue-400">0</div>
          </div>
          <div className="p-4 bg-slate-900/50 rounded-lg border border-slate-800">
            <div className="text-sm text-slate-400 mb-1">Pipeline Errors</div>
            <div className="text-2xl font-bold text-slate-400">0</div>
          </div>
        </div>
      </div>
    </div>
  );
}
