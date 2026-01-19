import { Link } from 'react-router-dom';

export default function Dashboard() {
  const stats = [
    { label: 'Total Devices', value: '12', unit: 'active', color: 'blue' },
    { label: 'Data Points', value: '1,240', unit: 'monitored', color: 'emerald' },
    { label: 'Throughput', value: '450', unit: 'ops/sec', color: 'purple' },
    { label: 'Errors', value: '0', unit: 'last 24h', color: 'red' },
  ];

  const features = [
    { 
      title: 'Device Manager', 
      desc: 'Connect Modbus, FATEK, and MQTT devices.',
      path: '/datalink/devices',
      icon: (
        <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
        </svg>
      ),
      gradient: 'from-blue-500 to-cyan-500'
    },
    { 
      title: 'Tag Dictionary', 
      desc: 'Manage standard tags and metadata.',
      path: '/datalink/tags',
      icon: (
        <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
        </svg>
      ),
      gradient: 'from-purple-500 to-pink-500'
    },
    { 
      title: 'Data Mapping', 
      desc: 'Transform and route data pipelines.',
      path: '/datalink/mappings',
      icon: (
        <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
           <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
        </svg>
      ),
      gradient: 'from-emerald-500 to-teal-500'
    },
    { 
      title: 'Mapping Wizard', 
      desc: 'Guided step-by-step workflow to create mappings.',
      path: '/datalink/wizard',
      icon: (
        <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
        </svg>
      ),
      gradient: 'from-cyan-500 to-blue-500',
      highlight: true
    },
     { 
      title: 'System Settings', 
      desc: 'Configure storage and performance.',
      path: '/datalink/settings',
      icon: (
        <svg className="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
      ),
      gradient: 'from-orange-500 to-amber-500'
    },
  ];

  return (
    <div className="space-y-8">
      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {stats.map((stat) => (
          <div key={stat.label} className="bg-slate-800/50 backdrop-blur-sm p-6 rounded-2xl border border-slate-700/50 hover:border-slate-600 transition-colors">
            <p className="text-sm text-slate-400 font-medium">{stat.label}</p>
            <div className="mt-2 flex items-baseline space-x-2">
              <span className="text-3xl font-bold text-slate-100">{stat.value}</span>
              <span className={`text-xs font-medium px-2 py-0.5 rounded-full bg-${stat.color}-500/10 text-${stat.color}-400 border border-${stat.color}-500/20`}>
                {stat.unit}
              </span>
            </div>
          </div>
        ))}
      </div>

      {/* Feature Cards */}
      <div>
        <h2 className="text-xl font-bold text-slate-100 mb-6 flex items-center gap-2">
            <span>Quick Access</span>
            <div className="h-px bg-slate-800 flex-1 ml-4"></div>
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-6">
          {features.map((feature) => (
            <Link
              key={feature.title}
              to={feature.path}
              className="group relative bg-slate-800 p-6 rounded-2xl border border-slate-700 hover:border-slate-600 transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:shadow-blue-500/5 overflow-hidden"
            >
              {/* Background Glow */}
              <div className={`absolute top-0 right-0 w-32 h-32 bg-gradient-to-br ${feature.gradient} opacity-0 group-hover:opacity-10 blur-3xl transition-opacity duration-500 -translate-y-1/2 translate-x-1/2`}></div>

              <div className="flex items-start space-x-6 relative z-10">
                <div className={`w-14 h-14 rounded-xl bg-gradient-to-br ${feature.gradient} flex items-center justify-center shadow-lg transform group-hover:scale-110 transition-transform duration-300`}>
                   {feature.icon}
                </div>
                <div className="flex-1">
                  <h3 className="text-lg font-bold text-slate-100 group-hover:text-blue-400 transition-colors">
                    {feature.title}
                  </h3>
                  <p className="mt-2 text-slate-400 leading-relaxed">
                    {feature.desc}
                  </p>
                </div>
                <div className="text-slate-600 group-hover:text-blue-400 transition-colors transform group-hover:translate-x-1 duration-300">
                  <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                  </svg>
                </div>
              </div>
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}
