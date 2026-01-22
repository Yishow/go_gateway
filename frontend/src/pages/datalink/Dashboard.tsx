import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { dashboardAPI, type DashboardStats, type DeviceStatus } from '../../services/datalink';

// 簡單 SVG 圖標組件
const RefreshIcon = ({ className }: { className?: string }) => (
  <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
  </svg>
);

const AlertCircleIcon = ({ className }: { className?: string }) => (
  <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
  </svg>
);

const CheckCircleIcon = ({ className }: { className?: string }) => (
  <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
  </svg>
);

const XCircleIcon = ({ className }: { className?: string }) => (
  <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
  </svg>
);

const ActivityIcon = ({ className }: { className?: string }) => (
  <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
  </svg>
);

/**
 * 儀表板組件
 * 
 * 顯示系統統計數據、設備狀態和錯誤資訊
 */
export default function Dashboard() {
  const { t } = useTranslation();
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [deviceStatuses, setDeviceStatuses] = useState<DeviceStatus[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lastUpdate, setLastUpdate] = useState<Date | null>(null);

  /**
   * 載入儀表板數據
   */
  const loadDashboardData = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const [statsData, deviceStatusData] = await Promise.all([
        dashboardAPI.getStats(),
        dashboardAPI.getDeviceStatuses(),
      ]);
      
      setStats(statsData);
      setDeviceStatuses(deviceStatusData);
      setLastUpdate(new Date());
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '載入儀表板數據失敗';
      setError(errorMessage);
      console.error('載入儀表板數據失敗:', err);
    } finally {
      setLoading(false);
    }
  };

  // 初始載入和定期刷新
  useEffect(() => {
    loadDashboardData();
    
    // 每30秒自動刷新一次
    const interval = setInterval(loadDashboardData, 30000);
    
    return () => clearInterval(interval);
  }, []);

  /**
   * 格式化數字顯示
   */
  const formatNumber = (num: number): string => {
    if (num >= 1000000) {
      return (num / 1000000).toFixed(1) + 'M';
    }
    if (num >= 1000) {
      return (num / 1000).toFixed(1) + 'K';
    }
    return num.toString();
  };

  /**
   * 解析時間字符串為 Date 對象
   */
  const parseTime = (timeStr: string | null | undefined): Date | null => {
    if (!timeStr) return null;
    
    // 嘗試解析 ISO 8601 格式
    const date = new Date(timeStr);
    
    // 檢查是否為有效日期
    if (isNaN(date.getTime())) {
      console.warn('無法解析時間字符串:', timeStr);
      return null;
    }
    
    return date;
  };

  /**
   * 格式化時間顯示（相對時間）
   */
  const formatTime = (date: Date | string | null | undefined): string => {
    const dateObj = typeof date === 'string' ? parseTime(date) : date;
    
    if (!dateObj) return '從未測試';
    
    const now = new Date();
    const diff = now.getTime() - dateObj.getTime();
    
    // 如果時間差為負數（未來時間），返回格式化日期
    if (diff < 0) {
      return dateObj.toLocaleString('zh-TW');
    }
    
    const seconds = Math.floor(diff / 1000);
    
    if (seconds < 60) {
      return `${seconds} 秒前`;
    }
    
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) {
      return `${minutes} 分鐘前`;
    }
    
    const hours = Math.floor(minutes / 60);
    if (hours < 24) {
      return `${hours} 小時前`;
    }
    
    const days = Math.floor(hours / 24);
    if (days < 30) {
      return `${days} 天前`;
    }
    
    const months = Math.floor(days / 30);
    if (months < 12) {
      return `${months} 個月前`;
    }
    
    // 超過一年，顯示完整日期
    return dateObj.toLocaleDateString('zh-TW', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  if (loading && !stats) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <RefreshIcon className="w-8 h-8 animate-spin text-slate-400 mx-auto mb-4" />
          <p className="text-slate-400">{t('common.loading')}</p>
        </div>
      </div>
    );
  }

  if (error && !stats) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <AlertCircleIcon className="w-8 h-8 text-red-400 mx-auto mb-4" />
          <p className="text-red-400 mb-4">{error}</p>
          <button
            onClick={loadDashboardData}
            className="px-4 py-2 bg-slate-700 hover:bg-slate-600 rounded-lg text-slate-100 transition-colors"
          >
            重試
          </button>
        </div>
      </div>
    );
  }

  const statsCards = stats ? [
    {
      label: t('dashboard.totalDevices'),
      value: stats.total_devices.toString(),
      unit: `${stats.active_devices} ${t('dashboard.active')}`,
      color: 'blue',
      icon: ActivityIcon,
    },
    {
      label: t('dashboard.dataPoints'),
      value: formatNumber(stats.total_points),
      unit: `${stats.enabled_points} ${t('dashboard.monitored')}`,
      color: 'emerald',
      icon: CheckCircleIcon,
    },
    {
      label: t('dashboard.throughput'),
      value: formatNumber(Math.round(stats.estimated_throughput)),
      unit: t('dashboard.opsPerSec'),
      color: 'purple',
      icon: ActivityIcon,
    },
    {
      label: t('dashboard.errors'),
      value: stats.errors_last_24h.toString(),
      unit: t('dashboard.last24h'),
      color: stats.errors_last_24h > 0 ? 'red' : 'green',
      icon: stats.errors_last_24h > 0 ? AlertCircleIcon : CheckCircleIcon,
    },
  ] : [];

  // 按狀態分組設備
  const devicesByStatus = deviceStatuses.reduce((acc, device) => {
    const status = device.status;
    if (!acc[status]) {
      acc[status] = [];
    }
    acc[status].push(device);
    return acc;
  }, {} as Record<string, DeviceStatus[]>);

  const getStatusColor = (status: string): string => {
    switch (status) {
      case 'active':
        return 'text-green-400 bg-green-500/10 border-green-500/20';
      case 'disabled':
        return 'text-slate-400 bg-slate-500/10 border-slate-500/20';
      case 'draft':
        return 'text-yellow-400 bg-yellow-500/10 border-yellow-500/20';
      case 'error':
        return 'text-red-400 bg-red-500/10 border-red-500/20';
      default:
        return 'text-slate-400 bg-slate-500/10 border-slate-500/20';
    }
  };

  /**
   * 翻譯設備狀態
   */
  const translateStatus = (status: string): string => {
    return t(`device.status.${status}`, status);
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'active':
        return <CheckCircleIcon className="w-4 h-4 text-green-400" />;
      case 'disabled':
        return <XCircleIcon className="w-4 h-4 text-slate-400" />;
      case 'draft':
        return <ActivityIcon className="w-4 h-4 text-yellow-400" />;
      case 'error':
        return <AlertCircleIcon className="w-4 h-4 text-red-400" />;
      default:
        return <ActivityIcon className="w-4 h-4 text-slate-400" />;
    }
  };

  return (
    <div className="space-y-8">
      {/* 標題和刷新按鈕 */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-slate-100">{t('dashboard.title')}</h1>
        <div className="flex items-center gap-4">
          {lastUpdate && (
            <span className="text-sm text-slate-400">
              最後更新: {formatTime(lastUpdate)}
            </span>
          )}
          <button
            onClick={loadDashboardData}
            disabled={loading}
            className="flex items-center gap-2 px-4 py-2 bg-slate-700 hover:bg-slate-600 rounded-lg text-slate-100 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <RefreshIcon className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
            刷新
          </button>
        </div>
      </div>

      {/* 統計卡片網格 */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {statsCards.map((stat) => {
          const Icon = stat.icon;
          return (
            <div
              key={stat.label}
              className="bg-slate-800/50 backdrop-blur-sm p-6 rounded-2xl border border-slate-700/50 hover:border-slate-600 transition-colors"
            >
              <div className="flex items-center justify-between mb-2">
                <p className="text-sm text-slate-400 font-medium">{stat.label}</p>
                <Icon className={`w-5 h-5 text-${stat.color}-400`} />
              </div>
              <div className="mt-2 flex items-baseline space-x-2">
                <span className="text-3xl font-bold text-slate-100">{stat.value}</span>
                <span
                  className={`text-xs font-medium px-2 py-0.5 rounded-full bg-${stat.color}-500/10 text-${stat.color}-400 border border-${stat.color}-500/20`}
                >
                  {stat.unit}
                </span>
              </div>
            </div>
          );
        })}
      </div>

      {/* 詳細統計 */}
      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {/* 設備統計 */}
          <div className="bg-slate-800/50 backdrop-blur-sm p-6 rounded-2xl border border-slate-700/50">
            <h2 className="text-lg font-semibold text-slate-100 mb-4">設備統計</h2>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-slate-400">總數</span>
                <span className="text-slate-100 font-semibold">{stats.total_devices}</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-slate-400">啟用</span>
                <span className="text-green-400 font-semibold">{stats.active_devices}</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-slate-400">停用</span>
                <span className="text-slate-400 font-semibold">{stats.disabled_devices}</span>
              </div>
              {stats.error_devices > 0 && (
                <div className="flex justify-between items-center">
                  <span className="text-slate-400">錯誤</span>
                  <span className="text-red-400 font-semibold">{stats.error_devices}</span>
                </div>
              )}
            </div>
          </div>

          {/* 點位統計 */}
          <div className="bg-slate-800/50 backdrop-blur-sm p-6 rounded-2xl border border-slate-700/50">
            <h2 className="text-lg font-semibold text-slate-100 mb-4">點位統計</h2>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-slate-400">總數</span>
                <span className="text-slate-100 font-semibold">{stats.total_points}</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-slate-400">啟用</span>
                <span className="text-green-400 font-semibold">{stats.enabled_points}</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-slate-400">停用</span>
                <span className="text-slate-400 font-semibold">
                  {stats.total_points - stats.enabled_points}
                </span>
              </div>
            </div>
          </div>

          {/* 標籤統計 */}
          <div className="bg-slate-800/50 backdrop-blur-sm p-6 rounded-2xl border border-slate-700/50">
            <h2 className="text-lg font-semibold text-slate-100 mb-4">標籤統計</h2>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-slate-400">總數</span>
                <span className="text-slate-100 font-semibold">{stats.total_tags}</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-slate-400">啟用</span>
                <span className="text-green-400 font-semibold">{stats.active_tags}</span>
              </div>
              {stats.retired_tags > 0 && (
                <div className="flex justify-between items-center">
                  <span className="text-slate-400">退役</span>
                  <span className="text-slate-400 font-semibold">{stats.retired_tags}</span>
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* 設備狀態列表 */}
      {deviceStatuses.length > 0 && (
        <div className="bg-slate-800/50 backdrop-blur-sm p-6 rounded-2xl border border-slate-700/50">
          <h2 className="text-lg font-semibold text-slate-100 mb-4">設備狀態</h2>
          <div className="space-y-3">
            {deviceStatuses.map((device) => (
              <div
                key={device.id}
                className="flex items-center justify-between p-4 bg-slate-700/30 rounded-lg border border-slate-600/30 hover:border-slate-500/50 transition-colors"
              >
                <div className="flex items-center gap-3 flex-1">
                  {getStatusIcon(device.status)}
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <span className="text-slate-100 font-medium">{device.name}</span>
                      <span className={`text-xs px-2 py-0.5 rounded-full ${getStatusColor(device.status)}`}>
                        {translateStatus(device.status)}
                      </span>
                    </div>
                    <div className="flex items-center gap-4 mt-1 text-xs text-slate-400">
                      <span>{device.protocol}</span>
                      <span>最後測試: {formatTime(device.last_test_at)}</span>
                    </div>
                  </div>
                </div>
                {device.last_error && (
                  <div className="text-xs text-red-400 max-w-xs truncate" title={device.last_error}>
                    {device.last_error}
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* 錯誤提示 */}
      {error && stats && (
        <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-4 flex items-center gap-3">
          <AlertCircleIcon className="w-5 h-5 text-red-400 flex-shrink-0" />
          <div className="flex-1">
            <p className="text-red-400 font-medium">載入部分數據時發生錯誤</p>
            <p className="text-red-300/80 text-sm mt-1">{error}</p>
          </div>
        </div>
      )}
    </div>
  );
}
