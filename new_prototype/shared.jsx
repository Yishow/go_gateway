// 共用元件 / 圖示 / 工具
// 全部掛到 window 讓其他 babel script 拿得到

const { useState, useEffect, useMemo, useRef, useCallback } = React;

// ──────────────────────────────────────────────────────────────────
// 工具
// ──────────────────────────────────────────────────────────────────
const clsx = (...args) => args.filter(Boolean).join(' ');

const formatAddr = (n) => String(n).padStart(5, '0');

// Modbus 位址 → 函數類型
const fnFromAddr = (addr) => {
  const n = parseInt(String(addr).replace(/[^0-9]/g, ''), 10);
  if (!n) return 'holding_register';
  const first = String(n)[0];
  if (first === '0') return 'coil';
  if (first === '1') return 'discrete_input';
  if (first === '3') return 'input_register';
  return 'holding_register';
};

const dataTypeWidth = (t) => {
  switch (t) {
    case 'bool': return 1;
    case 'int16':
    case 'uint16': return 1;
    case 'int32':
    case 'uint32':
    case 'float32': return 2;
    case 'int64':
    case 'uint64':
    case 'float64': return 4;
    default: return 1;
  }
};

// ──────────────────────────────────────────────────────────────────
// 圖示 (inline SVG)
// ──────────────────────────────────────────────────────────────────
const Icon = ({ name, className = 'w-4 h-4' }) => {
  const paths = {
    device: <><rect x="4" y="3" width="16" height="18" rx="2" /><path d="M9 7h6M9 11h6M9 15h3" /><circle cx="17" cy="15" r="1" fill="currentColor" /></>,
    rule: <><path d="M4 6h12M4 12h16M4 18h8" /><circle cx="20" cy="6" r="2" fill="currentColor" /><circle cx="14" cy="18" r="2" fill="currentColor" /></>,
    map: <><path d="M3 6l6-2 6 2 6-2v14l-6 2-6-2-6 2V6z" /><path d="M9 4v16M15 6v16" /></>,
    db: <><ellipse cx="12" cy="5" rx="8" ry="3" /><path d="M4 5v6c0 1.7 3.6 3 8 3s8-1.3 8-3V5" /><path d="M4 11v6c0 1.7 3.6 3 8 3s8-1.3 8-3v-6" /></>,
    check: <path d="M4 12l5 5L20 6" />,
    chevron: <path d="M9 6l6 6-6 6" />,
    plus: <><path d="M12 5v14M5 12h14" /></>,
    play: <path d="M6 4l14 8-14 8V4z" fill="currentColor" />,
    refresh: <path d="M4 4v6h6M20 20v-6h-6M5 13a8 8 0 0014.5 4.5M19 11a8 8 0 00-14.5-4.5" />,
    cable: <><path d="M5 8a3 3 0 016 0v8a3 3 0 006 0V8" /><path d="M5 4v8M19 12v8" /></>,
    sliders: <><path d="M4 6h10M18 6h2M4 12h2M10 12h10M4 18h14M18 18h2" /><circle cx="16" cy="6" r="2" fill="currentColor" /><circle cx="8" cy="12" r="2" fill="currentColor" /><circle cx="16" cy="18" r="2" fill="currentColor" /></>,
    tag: <><path d="M3 12V5a2 2 0 012-2h7l9 9-9 9-9-9z" /><circle cx="8" cy="8" r="1.5" fill="currentColor" /></>,
    table: <><rect x="3" y="4" width="18" height="16" rx="1.5" /><path d="M3 10h18M3 16h18M10 4v16M16 4v16" /></>,
    alert: <><path d="M12 9v4M12 17h.01" /><path d="M10.3 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.7 3.86a2 2 0 00-3.4 0z" /></>,
    info: <><circle cx="12" cy="12" r="9" /><path d="M12 8h.01M11 12h1v5h1" /></>,
    spark: <><path d="M5 12l4-4 3 6 4-9 3 7" /></>,
    arrow: <path d="M5 12h14M13 5l7 7-7 7" />,
    close: <path d="M6 6l12 12M18 6L6 18" />,
    save: <><path d="M19 21H5a2 2 0 01-2-2V5a2 2 0 012-2h11l5 5v11a2 2 0 01-2 2z" /><path d="M17 21v-8H7v8M7 3v5h8" /></>,
    eye: <><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8S1 12 1 12z" /><circle cx="12" cy="12" r="3" /></>,
    bolt: <path d="M13 2L3 14h7l-1 8 10-12h-7l1-8z" fill="currentColor" />,
    flow: <><path d="M3 6h6l3 6 3-6h6" /><circle cx="3" cy="6" r="1.5" fill="currentColor" /><circle cx="21" cy="6" r="1.5" fill="currentColor" /></>,
    grid: <><rect x="3" y="3" width="7" height="7" rx="1" /><rect x="14" y="3" width="7" height="7" rx="1" /><rect x="3" y="14" width="7" height="7" rx="1" /><rect x="14" y="14" width="7" height="7" rx="1" /></>,
  };
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.75"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      aria-hidden="true"
    >
      {paths[name] ?? null}
    </svg>
  );
};

// ──────────────────────────────────────────────────────────────────
// Status chip
// ──────────────────────────────────────────────────────────────────
const StatusChip = ({ tone = 'neutral', children, dot = true, className = '' }) => {
  const tones = {
    neutral: 'border-slate-700/70 bg-slate-800/60 text-slate-300',
    info: 'border-blue-500/30 bg-blue-500/10 text-blue-200',
    success: 'border-emerald-500/30 bg-emerald-500/10 text-emerald-200',
    warning: 'border-amber-500/30 bg-amber-500/10 text-amber-200',
    error: 'border-red-500/30 bg-red-500/10 text-red-200',
    draft: 'border-slate-500/30 bg-slate-700/40 text-slate-300',
  };
  const dots = {
    neutral: 'bg-slate-400',
    info: 'bg-blue-400',
    success: 'bg-emerald-400',
    warning: 'bg-amber-400',
    error: 'bg-red-400',
    draft: 'bg-slate-400',
  };
  return (
    <span className={clsx('chip', tones[tone], className)}>
      {dot && <span className={clsx('w-1.5 h-1.5 rounded-full', dots[tone], tone === 'info' || tone === 'success' ? 'pulse-dot' : '')} />}
      {children}
    </span>
  );
};

// ──────────────────────────────────────────────────────────────────
// Section card
// ──────────────────────────────────────────────────────────────────
const SectionCard = ({ title, subtitle, icon, aside, children, className = '', contentClassName = '' }) => (
  <section className={clsx('rounded-2xl border border-slate-700/60 bg-slate-900/60 backdrop-blur-sm shadow-xl shadow-black/30', className)}>
    {(title || aside) && (
      <header className="flex items-start justify-between gap-4 border-b border-slate-700/60 px-5 py-4">
        <div className="flex items-start gap-3">
          {icon && (
            <div className="grid place-items-center w-9 h-9 rounded-lg bg-blue-500/10 border border-blue-500/20 text-blue-300">
              {icon}
            </div>
          )}
          <div>
            <h3 className="text-sm font-semibold text-slate-100 leading-tight">{title}</h3>
            {subtitle && <p className="text-xs text-slate-400 mt-0.5">{subtitle}</p>}
          </div>
        </div>
        {aside && <div className="flex items-center gap-2">{aside}</div>}
      </header>
    )}
    <div className={clsx('p-5', contentClassName)}>{children}</div>
  </section>
);

// ──────────────────────────────────────────────────────────────────
// Field 包裝
// ──────────────────────────────────────────────────────────────────
const Field = ({ label, hint, error, required, children, className = '' }) => (
  <label className={clsx('block', className)}>
    <span className={clsx('label', required && 'after:content-["*"] after:text-red-400 after:ml-0.5')}>{label}</span>
    {children}
    {error ? (
      <span className="mt-1 flex items-center gap-1 text-xs text-red-300">
        <Icon name="alert" className="w-3 h-3" />
        {error}
      </span>
    ) : hint ? (
      <span className="mt-1 block text-[11px] text-slate-500">{hint}</span>
    ) : null}
  </label>
);

const Input = React.forwardRef(({ className = '', ...props }, ref) => (
  <input
    ref={ref}
    {...props}
    className={clsx(
      'w-full rounded-lg border border-slate-700/60 bg-slate-950/40 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 transition-colors focus:border-blue-500/70 focus:outline-none focus:ring-2 focus:ring-blue-500/20',
      className,
    )}
  />
));

const Select = ({ className = '', children, ...props }) => (
  <div className="relative">
    <select
      {...props}
      className={clsx(
        'w-full appearance-none rounded-lg border border-slate-700/60 bg-slate-950/40 px-3 py-2 pr-9 text-sm text-slate-100 transition-colors focus:border-blue-500/70 focus:outline-none focus:ring-2 focus:ring-blue-500/20',
        className,
      )}
    >
      {children}
    </select>
    <svg viewBox="0 0 24 24" className="pointer-events-none absolute right-2.5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500">
      <path d="M6 9l6 6 6-6" stroke="currentColor" strokeWidth="1.75" fill="none" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  </div>
);

const Textarea = ({ className = '', ...props }) => (
  <textarea
    {...props}
    className={clsx(
      'w-full rounded-lg border border-slate-700/60 bg-slate-950/40 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 transition-colors focus:border-blue-500/70 focus:outline-none focus:ring-2 focus:ring-blue-500/20',
      className,
    )}
  />
);

const Button = ({ variant = 'primary', size = 'md', icon, children, className = '', ...props }) => {
  const base = 'inline-flex items-center justify-center gap-2 font-medium rounded-lg transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/50 disabled:opacity-50 disabled:cursor-not-allowed';
  const sizes = {
    sm: 'px-2.5 py-1.5 text-xs',
    md: 'px-4 py-2 text-sm',
    lg: 'px-5 py-2.5 text-sm',
  };
  const variants = {
    primary: 'bg-blue-600 text-white hover:bg-blue-500 shadow-lg shadow-blue-500/20',
    secondary: 'bg-slate-800 text-slate-200 hover:bg-slate-700 border border-slate-700/60',
    ghost: 'text-slate-300 hover:bg-slate-800/70',
    danger: 'bg-red-600 text-white hover:bg-red-500',
    success: 'bg-emerald-600 text-white hover:bg-emerald-500 shadow-lg shadow-emerald-500/20',
  };
  return (
    <button {...props} className={clsx(base, sizes[size], variants[variant], className)}>
      {icon}
      {children}
    </button>
  );
};

// ──────────────────────────────────────────────────────────────────
// Toggle - 用精確像素確保 ON 狀態下 dot 不會視覺溢出
// ──────────────────────────────────────────────────────────────────
const Toggle = ({ checked, onChange, label, size = 'md' }) => {
  const dim = size === 'sm'
    ? { w: 32, h: 16, dot: 12, pad: 2 }
    : { w: 40, h: 20, dot: 16, pad: 2 };
  const tx = checked ? dim.w - dim.pad * 2 - dim.dot : 0;
  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      onClick={() => onChange(!checked)}
      className="inline-flex items-center gap-2 group"
    >
      <span
        className={clsx('relative inline-block rounded-full transition-colors flex-shrink-0', checked ? 'bg-blue-600' : 'bg-slate-700')}
        style={{ width: dim.w, height: dim.h }}
      >
        <span
          className="absolute rounded-full bg-white shadow-sm transition-transform will-change-transform"
          style={{
            width: dim.dot,
            height: dim.dot,
            top: dim.pad,
            left: dim.pad,
            transform: `translateX(${tx}px)`,
          }}
        />
      </span>
      {label && <span className="text-xs text-slate-300">{label}</span>}
    </button>
  );
};

// ──────────────────────────────────────────────────────────────────
// 共享預設資料
// ──────────────────────────────────────────────────────────────────
const PROTOCOLS = [
  { id: 'modbus_tcp', name: 'Modbus TCP', desc: '以太網 Modbus，最常見' },
  { id: 'modbus_rtu', name: 'Modbus RTU', desc: '序列埠 Modbus' },
  { id: 'modbus_udp', name: 'Modbus UDP', desc: 'UDP 傳輸' },
  { id: 'fatek_fbs', name: 'Fatek FBS', desc: '永宏 PLC' },
  { id: 'mc_3e', name: 'MC 3E', desc: '三菱 MELSEC' },
  { id: 'mqtt', name: 'MQTT', desc: '訊息佇列' },
];

const DATA_TYPES = ['bool', 'int16', 'int32', 'int64', 'uint16', 'uint32', 'uint64', 'float32', 'float64', 'string'];

const DB_KINDS = [
  { id: 'sqlite', name: 'SQLite', icon: '🗄️' },
  { id: 'postgres', name: 'PostgreSQL', icon: '🐘' },
  { id: 'mysql', name: 'MySQL', icon: '🐬' },
  { id: 'sqlserver', name: 'SQL Server', icon: '🪟' },
];

// 範例已存在的資料庫表欄位 (給 Step 4 自動補齊)
const SAMPLE_DB_TABLES = {
  postgres: [
    {
      schema: 'public',
      name: 'sensor_readings',
      columns: [
        { name: 'ts', data_type: 'timestamptz', nullable: false, primary_key: true },
        { name: 'temp_in_c', data_type: 'double precision', nullable: true, primary_key: false },
        { name: 'temp_out_c', data_type: 'double precision', nullable: true, primary_key: false },
        { name: 'pressure_main_kpa', data_type: 'double precision', nullable: true, primary_key: false },
        { name: 'pressure_sub_kpa', data_type: 'double precision', nullable: true, primary_key: false },
        { name: 'flow_lpm', data_type: 'double precision', nullable: true, primary_key: false },
        { name: 'humidity_pct', data_type: 'double precision', nullable: true, primary_key: false },
        { name: 'vibration_mms', data_type: 'double precision', nullable: true, primary_key: false },
        { name: 'motor_rpm', data_type: 'integer', nullable: true, primary_key: false },
      ],
    },
  ],
};

// 規則 ID 生成器
let __ruleSeq = 0;
const nextRuleId = () => `rule-${Date.now().toString(36)}-${(__ruleSeq++).toString(36)}`;

// 裝置 ID 生成器
let __devSeq = 0;
const nextDeviceId = () => `dev-${Date.now().toString(36)}-${(__devSeq++).toString(36)}`;

// 預設裝置 (template)
const makeDefaultDevice = (overrides = {}) => ({
  id: nextDeviceId(),
  name: '新設備',
  description: '',
  protocol: 'modbus_tcp',
  config: { host: '192.168.1.100', port: 502, slave_id: 1, timeout: 5 },
  status: 'draft',
  test: null,
  ...overrides,
});

// 預設規則 (template)
const makeDefaultRule = (overrides = {}) => ({
  id: nextRuleId(),
  device_id: null,
  name: '規則 1',
  start_address: '40001',
  count: 8,
  data_type: 'int16',
  naming_prefix: 'SENSOR_',
  enabled: true,
  scale_multiplier: 0.1,
  scale_offset: 0,
  data_format: '',
  skipped_addresses: [],
  // Local Modbus Share 對外發布設定 (per-rule)
  share_enabled: true,
  share_start_register: null,
  share_stride: null,
  ...overrides,
});

// 計算每條規則在 Modbus Share 中佔用的 register 起始
// 自動模式：依規則順序連續配置；手動模式：用 rule.share_start_register
const computeShareLayout = (rules, baseRegister = 40001) => {
  const layout = {}; // rule_id => { start, stride, end }
  let cursor = baseRegister;
  (rules || []).forEach((r) => {
    if (!r.share_enabled) {
      layout[r.id] = null;
      return;
    }
    const stride = r.share_stride ?? dataTypeWidth(r.data_type);
    const enabledCount = r.count - (r.skipped_addresses?.length || 0);
    const start = (r.share_start_register != null && r.share_start_register !== '') ? +r.share_start_register : cursor;
    layout[r.id] = { start, stride, end: start + enabledCount * stride, auto: r.share_start_register == null || r.share_start_register === '' };
    cursor = Math.max(cursor, start + enabledCount * stride);
  });
  return layout;
};

// 預設範例值
const DEFAULTS = {
  device: {
    name: 'PLC-生產線-01',
    description: 'Modbus TCP PLC (Line A 主控)',
    protocol: 'modbus_tcp',
    config: { host: '192.168.1.100', port: 502, slave_id: 1, timeout: 5 },
  },
  // 保留 rule (向後相容)，新增 rules 陣列為主資料
  rule: {
    start_address: '40001',
    count: 8,
    data_type: 'int16',
    naming_prefix: 'SENSOR_',
    enabled: true,
    scale_multiplier: 0.1,
    scale_offset: 0,
  },
  pollingGroup: { name: '快速輪詢 (1s)', interval_ms: 1000 },
};

// 給每個點位一個語意 label - 純展示用
const POINT_SEMANTIC = [
  { suffix: 'TEMP_IN', display: '入口溫度', unit: '°C', tag_key: 'line01.temp.inlet' },
  { suffix: 'TEMP_OUT', display: '出口溫度', unit: '°C', tag_key: 'line01.temp.outlet' },
  { suffix: 'PRES_01', display: '主管路壓力', unit: 'kPa', tag_key: 'line01.pressure.main' },
  { suffix: 'PRES_02', display: '次管路壓力', unit: 'kPa', tag_key: 'line01.pressure.sub' },
  { suffix: 'FLOW_01', display: '流量計', unit: 'L/min', tag_key: 'line01.flow.q1' },
  { suffix: 'HUM_01', display: '濕度', unit: '%', tag_key: 'line01.humidity.amb' },
  { suffix: 'VIB_01', display: '振動', unit: 'mm/s', tag_key: 'line01.vibration.motor' },
  { suffix: 'RPM_01', display: '馬達轉速', unit: 'rpm', tag_key: 'line01.motor.rpm' },
];

// ──────────────────────────────────────────────────────────────────
// 衍生：根據 rule 產生 points
// ──────────────────────────────────────────────────────────────────
const derivePoints = (rule, deviceId, skippedSet = new Set()) => {
  const startN = parseInt(String(rule.start_address).replace(/[^0-9]/g, ''), 10) || 40001;
  const width = dataTypeWidth(rule.data_type);
  const points = [];
  for (let i = 0; i < rule.count; i++) {
    const addr = startN + i * width;
    const addrStr = formatAddr(addr);
    const skipped = skippedSet.has(addrStr);
    const semantic = POINT_SEMANTIC[i % POINT_SEMANTIC.length];
    points.push({
      id: `pt-${deviceId}-${rule.id || 'r'}-${addrStr}`,
      device_id: deviceId,
      rule_id: rule.id,
      rule_name: rule.name,
      name: `${rule.naming_prefix}${String(i + 1).padStart(2, '0')}`,
      address: addrStr,
      data_type: rule.data_type,
      function: fnFromAddr(addr),
      width,
      enabled: !skipped,
      skipped,
      display: semantic.display,
      unit: semantic.unit,
      tag_key_suggest: semantic.tag_key,
      // 規則層級預設轉換
      _rule_scale: rule.scale_multiplier ?? 1,
      _rule_offset: rule.scale_offset ?? 0,
    });
  }
  return points;
};

// 多規則：把每條規則展開後串接 (每條規則使用自身的 device_id)
const deriveAllPoints = (rules, fallbackDeviceId) => {
  const all = [];
  (rules || []).forEach((r) => {
    if (r.enabled === false) return;
    const skippedSet = new Set(r.skipped_addresses || []);
    all.push(...derivePoints(r, r.device_id || fallbackDeviceId, skippedSet));
  });
  return all;
};

// ──────────────────────────────────────────────────────────────────
// Export to window
// ──────────────────────────────────────────────────────────────────
Object.assign(window, {
  clsx, formatAddr, fnFromAddr, dataTypeWidth,
  Icon, StatusChip, SectionCard, Field, Input, Select, Textarea, Button, Toggle,
  PROTOCOLS, DATA_TYPES, DB_KINDS, SAMPLE_DB_TABLES, DEFAULTS, POINT_SEMANTIC,
  derivePoints, deriveAllPoints, makeDefaultRule, nextRuleId, computeShareLayout,
  makeDefaultDevice, nextDeviceId,
});
