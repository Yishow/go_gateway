import { useEffect, useState, useRef, useCallback } from "react";
import { useTranslation } from "react-i18next";
import {
  dashboardAPI,
  type DashboardStats,
  type DeviceStatus,
} from "../../services/datalink";
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

// --- Icons (Simple SVG Components) ---

const LoadingIcon = ({ className }: { className?: string }) => (
  <svg
    className={className}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
    />
  </svg>
);

const RefreshIcon = ({
  className,
  isLoading,
}: {
  className?: string;
  isLoading?: boolean;
}) => (
  <svg
    className={`${className} ${isLoading ? "animate-spin" : ""}`}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
    />
  </svg>
);

const AlertCircleIcon = ({ className }: { className?: string }) => (
  <svg
    className={className}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
    />
  </svg>
);

const ActivityIcon = ({ className }: { className?: string }) => (
  <svg
    className={className}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M13 10V3L4 14h7v7l9-11h-7z"
    />
  </svg>
);

const ServerIcon = ({ className }: { className?: string }) => (
  <svg
    className={className}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"
    />
  </svg>
);

const DatabaseIcon = ({ className }: { className?: string }) => (
  <svg
    className={className}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4"
    />
  </svg>
);

const TagIcon = ({ className }: { className?: string }) => (
  <svg
    className={className}
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"
    />
  </svg>
);

// --- Helper Functions ---

const formatNumber = (num: number): string => {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + "M";
  if (num >= 1000) return (num / 1000).toFixed(1) + "K";
  return num.toString();
};

const formatTimeAgo = (dateStr: string | undefined): string => {
  if (!dateStr) return "從未測試";
  const date = new Date(dateStr);
  if (isNaN(date.getTime())) return "無效日期";

  const now = new Date();
  const diff = now.getTime() - date.getTime();
  if (diff < 0) return date.toLocaleString("zh-TW"); // Future time

  const seconds = Math.floor(diff / 1000);
  if (seconds < 60) return `${seconds} 秒前`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes} 分鐘前`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours} 小時前`;
  return date.toLocaleDateString("zh-TW");
};

// --- Components ---

import type { TFunction } from 'i18next';

interface StatsCardProps {
    title: string;
    value: string | number;
    unit?: string;
    icon: React.ComponentType<{ className?: string }>;
    color: string;
    trend?: {
        value: string;
        isPositive: boolean;
        label: string;
    };
}

const StatsCard = ({ title, value, unit, icon: Icon, color, trend }: StatsCardProps) => (
  <div
    className={`
        relative overflow-hidden
        bg-slate-800/50 backdrop-blur-sm p-6 rounded-2xl 
        border border-slate-700/50 hover:border-${color}-500/30 
        transition-all duration-300 group
    `}
  >
    {/* Background Glow */}
    <div
      className={`absolute -right-4 -top-4 w-24 h-24 bg-${color}-500/10 rounded-full blur-2xl group-hover:bg-${color}-500/20 transition-all`}
    ></div>

    <div className="flex justify-between items-start mb-4">
      <div>
        <p className="text-sm text-slate-400 font-medium mb-1">{title}</p>
        <div className="flex items-baseline gap-2">
          <h3 className="text-3xl font-bold text-slate-100 tracking-tight">
            {value}
          </h3>
          {unit && (
            <span className="text-sm text-slate-500 font-medium">{unit}</span>
          )}
        </div>
      </div>
      <div
        className={`p-3 rounded-xl bg-${color}-500/10 text-${color}-400 ring-1 ring-${color}-500/20`}
      >
        <Icon className="w-6 h-6" />
      </div>
    </div>

    {trend && (
      <div className="flex items-center gap-2 text-xs">
        <span
          className={`font-medium ${trend.isPositive ? "text-green-400" : "text-slate-400"}`}
        >
          {trend.value}
        </span>
        <span className="text-slate-500">{trend.label}</span>
      </div>
    )}
  </div>
);

const DeviceCard = ({ device, t }: { device: DeviceStatus, t: TFunction }) => {
  const isError = device.status === "error";
  const isActive = device.status === "active";

  let statusColor = "slate";
  if (isActive) statusColor = "green";
  else if (isError) statusColor = "red";
  else if (device.status === "draft") statusColor = "yellow";

  return (
    <div className="bg-slate-800/40 border border-slate-700/50 rounded-xl p-4 hover:bg-slate-800/60 transition-colors group">
      <div className="flex justify-between items-start mb-3">
        <div className="flex items-center gap-3">
          <div
            className={`
                        w-10 h-10 rounded-lg flex items-center justify-center
                        bg-${statusColor}-500/10 text-${statusColor}-400
                        ring-1 ring-${statusColor}-500/20
                     `}
          >
            <ServerIcon className="w-5 h-5" />
          </div>
          <div>
            <h4 className="font-semibold text-slate-200 group-hover:text-white transition-colors">
              {device.name}
            </h4>
            <span className="text-xs text-slate-500 font-mono">
              {device.protocol}
            </span>
          </div>
        </div>
        <span
          className={`
                    text-[10px] uppercase font-bold px-2 py-0.5 rounded-full
                    bg-${statusColor}-500/10 text-${statusColor}-400 border border-${statusColor}-500/20
                `}
        >
          {t(`device.status.${device.status}`, device.status)}
        </span>
      </div>

      <div className="space-y-2">
        <div className="flex justify-between text-xs">
          <span className="text-slate-500">最後回應</span>
          <span className="text-slate-300">
            {formatTimeAgo(device.last_test_at)}
          </span>
        </div>
        {/* Mini Status Bar */}
        <div className="w-full h-1.5 bg-slate-700/50 rounded-full overflow-hidden flex">
          {/* Simulated "health" bar for visual effect */}
          <div
            className={`h-full bg-${statusColor}-500/60 w-full`}
            style={{ width: isActive ? "100%" : "10%" }}
          ></div>
        </div>
      </div>

      {device.last_error && (
        <div
          className="mt-3 pt-3 border-t border-slate-700/50 text-xs text-red-400 truncate"
          title={device.last_error}
        >
          ⚠️ {device.last_error}
        </div>
      )}
    </div>
  );
};

interface HistoryPoint {
  time: string;
  value: number;
}

// --- Main Dashboard Component ---

export default function Dashboard() {
  const { t } = useTranslation();
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [deviceStatuses, setDeviceStatuses] = useState<DeviceStatus[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lastUpdate, setLastUpdate] = useState<Date | null>(null);
  
  // Track if we have loaded data at least once to avoid flickering loading state
  const hasDataRef = useRef(false);
  
  // Local history for the chart (Client-side buffer)
  const [throughputHistory, setThroughputHistory] = useState<HistoryPoint[]>([]);

  const loadDashboardData = useCallback(async () => {
    try {
      // Don't set global loading true on refresh to avoid flickering entire UI
      // Use ref to check if we already have data
      if (!hasDataRef.current) setLoading(true);
      setError(null);
      
      const [statsData, deviceStatusData] = await Promise.all([
        dashboardAPI.getStats(),
        dashboardAPI.getDeviceStatuses(),
      ]);
      
      setStats(statsData);
      setDeviceStatuses(deviceStatusData);
      setLastUpdate(new Date());
      hasDataRef.current = true;

      // Update Chart History
      const now = new Date();
      const timeStr = now.toLocaleTimeString('zh-TW', { hour12: false, hour: '2-digit', minute: '2-digit', second: '2-digit' });
      
      setThroughputHistory(prev => {
          const newHistory = [...prev, { time: timeStr, value: Math.round(statsData.estimated_throughput) }];
          // Keep last 20 points (approx 100 seconds at 5s interval)
          if (newHistory.length > 20) return newHistory.slice(newHistory.length - 20);
          return newHistory;
      });

    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '載入儀表板數據失敗';
      // Only show full error if we have no data at all
      if (!hasDataRef.current) setError(errorMessage);
      // eslint-disable-next-line no-console
      console.error('載入儀表板數據失敗:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadDashboardData();
    // Faster poll for "live" feel
    const interval = setInterval(loadDashboardData, 5000);
    return () => clearInterval(interval);
  }, [loadDashboardData]);

  if (loading && !stats) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh]">
        <LoadingIcon className="w-10 h-10 animate-spin text-blue-500 mb-4" />
        <p className="text-slate-400 font-medium animate-pulse">
          {t("common.loading")}
        </p>
      </div>
    );
  }

  if (error && !stats) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh]">
        <div className="bg-red-500/10 p-4 rounded-full mb-4 ring-1 ring-red-500/20">
          <AlertCircleIcon className="w-8 h-8 text-red-400" />
        </div>
        <h3 className="text-xl font-bold text-slate-200 mb-2">無法載入數據</h3>
        <p className="text-slate-400 mb-6 max-w-md text-center">{error}</p>
        <button
          onClick={loadDashboardData}
          className="px-6 py-2.5 bg-blue-600 hover:bg-blue-500 text-white rounded-lg font-medium transition-all shadow-lg shadow-blue-500/20"
        >
          重新嘗試
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      {/* Header Section */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-100 mb-1">
            {t("dashboard.title")}
          </h1>
          <p className="text-slate-400 text-sm">系統運作概況與即時監控</p>
        </div>
        <div className="flex items-center gap-3 bg-slate-800/50 p-1.5 rounded-lg border border-slate-700/50">
          <span className="text-xs text-slate-500 px-2">
            更新於 {lastUpdate?.toLocaleTimeString()}
          </span>
          <button
            onClick={loadDashboardData}
            disabled={loading}
            className="p-1.5 hover:bg-slate-700 rounded-md text-slate-400 hover:text-white transition-colors"
            title="刷新數據"
          >
            <RefreshIcon className="w-4 h-4" isLoading={loading} />
          </button>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {stats && (
          <>
            <StatsCard
              title={t("dashboard.throughput")}
              value={formatNumber(Math.round(stats.estimated_throughput))}
              unit="ops/sec"
              icon={ActivityIcon}
              color="blue"
              trend={{ value: "系統運作中", isPositive: true, label: "即時" }}
            />
            <StatsCard
              title={t("dashboard.totalDevices")}
              value={stats.total_devices}
              unit={`/ ${stats.active_devices} ${t("dashboard.active")}`}
              icon={ServerIcon}
              color="indigo"
              trend={{
                value: `${((stats.active_devices / (stats.total_devices || 1)) * 100).toFixed(0)}%`,
                isPositive: true,
                label: "運行率",
              }}
            />
            <StatsCard
              title={t("dashboard.dataPoints")}
              value={formatNumber(stats.total_points)}
              icon={DatabaseIcon}
              color="emerald"
              trend={{
                value: `${stats.enabled_points}`,
                isPositive: true,
                label: "監控中",
              }}
            />
            <StatsCard
              title={t("dashboard.errors")}
              value={stats.errors_last_24h}
              icon={AlertCircleIcon}
              color={stats.errors_last_24h > 0 ? "red" : "green"}
              trend={{
                value: "24H",
                isPositive: stats.errors_last_24h === 0,
                label: "錯誤統計",
              }}
            />
          </>
        )}
      </div>

      {/* Main Content Grid (Chart + Devices) */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Chart Section (Span 2) */}
        <div className="lg:col-span-2 bg-slate-800/50 backdrop-blur-sm border border-slate-700/50 rounded-2xl p-6">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-lg font-semibold text-slate-100 flex items-center gap-2">
              <ActivityIcon className="w-5 h-5 text-blue-400" />
              流量趨勢
            </h2>
            <div className="text-xs font-mono text-slate-500 bg-slate-800 px-2 py-1 rounded">
              Live / 5s
            </div>
          </div>

          <div className="h-[300px] w-full">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={throughputHistory}>
                <defs>
                  <linearGradient id="colorValue" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.3} />
                    <stop offset="95%" stopColor="#3b82f6" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <XAxis
                  dataKey="time"
                  stroke="#64748b"
                  fontSize={12}
                  tickLine={false}
                  axisLine={false}
                />
                <YAxis
                  stroke="#64748b"
                  fontSize={12}
                  tickLine={false}
                  axisLine={false}
                  tickFormatter={(value) => `${value}`}
                />
                <Tooltip
                  contentStyle={{
                    backgroundColor: "#1e293b",
                    borderColor: "#334155",
                    color: "#f8fafc",
                  }}
                  itemStyle={{ color: "#93c5fd" }}
                />
                <Area
                  type="monotone"
                  dataKey="value"
                  stroke="#3b82f6"
                  strokeWidth={2}
                  fillOpacity={1}
                  fill="url(#colorValue)"
                  isAnimationActive={false} // Disable animation for smoother live updates
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Quick Stats / System Health (Span 1) */}
        <div className="bg-slate-800/50 backdrop-blur-sm border border-slate-700/50 rounded-2xl p-6">
          <h2 className="text-lg font-semibold text-slate-100 mb-4 flex items-center gap-2">
            <TagIcon className="w-5 h-5 text-emerald-400" />
            系統概況
          </h2>

          <div className="space-y-4">
            <div className="flex justify-between items-center p-3 bg-slate-700/30 rounded-lg">
              <span className="text-slate-400 text-sm">標籤總數</span>
              <span className="text-slate-200 font-mono font-bold">
                {stats?.total_tags || 0}
              </span>
            </div>
            <div className="flex justify-between items-center p-3 bg-slate-700/30 rounded-lg">
              <span className="text-slate-400 text-sm">活躍標籤</span>
              <span className="text-emerald-400 font-mono font-bold">
                {stats?.active_tags || 0}
              </span>
            </div>
            <div className="flex justify-between items-center p-3 bg-slate-700/30 rounded-lg">
              <span className="text-slate-400 text-sm">退役標籤</span>
              <span className="text-slate-500 font-mono font-bold">
                {stats?.retired_tags || 0}
              </span>
            </div>

            <div className="my-4 border-t border-slate-700/50"></div>

            <h3 className="text-sm font-semibold text-slate-300 mb-2">
              狀態分佈
            </h3>
            {/* Fake Distribution Bar */}
            <div className="w-full h-2 bg-slate-700 rounded-full overflow-hidden flex">
              <div
                style={{
                  width: `${stats ? (stats.active_devices / (stats.total_devices || 1)) * 100 : 0}%`,
                }}
                className="bg-green-500 h-full"
              ></div>
              <div
                style={{
                  width: `${stats ? (stats.error_devices / (stats.total_devices || 1)) * 100 : 0}%`,
                }}
                className="bg-red-500 h-full"
              ></div>
              <div className="flex-1 bg-slate-600 h-full"></div>
            </div>
            <div className="flex justify-between text-xs text-slate-500 mt-1">
              <span className="flex items-center gap-1">
                <div className="w-2 h-2 rounded-full bg-green-500"></div> 正常
              </span>
              <span className="flex items-center gap-1">
                <div className="w-2 h-2 rounded-full bg-red-500"></div> 異常
              </span>
              <span className="flex items-center gap-1">
                <div className="w-2 h-2 rounded-full bg-slate-600"></div> 停用
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Device Status Grid */}
      <div>
        <h2 className="text-xl font-bold text-slate-100 mb-4 flex items-center gap-2">
          <ServerIcon className="w-6 h-6 text-indigo-400" />
          設備監控
        </h2>

        {deviceStatuses.length > 0 ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {deviceStatuses.map((device) => (
              <DeviceCard key={device.id} device={device} t={t} />
            ))}
          </div>
        ) : (
          <div className="text-center py-12 bg-slate-800/30 rounded-2xl border border-slate-700/30 border-dashed">
            <p className="text-slate-500">目前沒有已連接的設備</p>
          </div>
        )}
      </div>
    </div>
  );
}
