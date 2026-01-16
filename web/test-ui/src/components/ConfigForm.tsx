import { useState } from "react";
import { useTestAPI } from "../services/api";

interface ConfigFormProps {
  protocol: string;
  mode: string;
  config: Record<string, any>;
  onConfigChange: (config: Record<string, any>) => void;
  connectionId: string | null;
  onConnectionChange: (id: string | null) => void;
  onMinimize?: () => void;
}

/**
 * 輸入群組組件（移到組件外部以避免重新創建）
 */
const InputGroup = ({ label, children }: { label: string; children: React.ReactNode }) => (
  <div className="flex flex-col space-y-1.5">
    <label className="text-sm font-medium text-gray-700">{label}</label>
    {children}
  </div>
);

/**
 * 樣式化輸入框組件（移到組件外部以避免重新創建）
 */
const StyledInput = (props: React.InputHTMLAttributes<HTMLInputElement>) => (
  <input
    {...props}
    className="w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all duration-200"
  />
);

/**
 * 樣式化選擇框組件（移到組件外部以避免重新創建）
 */
const StyledSelect = (props: React.SelectHTMLAttributes<HTMLSelectElement>) => (
  <div className="relative">
    <select
      {...props}
      className="w-full px-4 py-2 appearance-none border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white transition-all duration-200"
    />
    <div className="absolute inset-y-0 right-0 flex items-center px-2 pointer-events-none text-gray-500">
      <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
      </svg>
    </div>
  </div>
);

export default function ConfigForm({
  protocol,
  mode,
  config,
  onConfigChange,
  connectionId,
  onConnectionChange,
  onMinimize,
}: ConfigFormProps) {
  const [loading, setLoading] = useState(false);
  const { connect, disconnect } = useTestAPI();

  const isTCP = mode === "tcp" || mode === "udp";
  const isSerial = mode === "serial";

  /**
   * 處理配置變更
   */
  const handleConfigChange = (key: string, value: any) => {
    onConfigChange({ ...config, [key]: value });
  };

  /**
   * 處理數字輸入框變更（保持字符串格式，避免失去焦點）
   * 使用受控組件模式，確保輸入框不會失去焦點
   */
  const handleNumberChange = (key: string, value: string) => {
    // 允許空值或純數字字符串
    if (value === "" || /^\d+$/.test(value)) {
      // 空值時設置為 undefined，數字時轉換為整數
      handleConfigChange(key, value === "" ? undefined : parseInt(value, 10));
    }
  };

  /**
   * 處理連線
   */
  const handleConnect = async () => {
    // 驗證 TCP/UDP 配置
    if (isTCP) {
      if (!config.host || config.host.trim() === "") {
        alert("請輸入主機位址");
        return;
      }
      if (!config.port || config.port <= 0) {
        alert("請輸入有效的埠號");
        return;
      }
    }

    // 驗證 Serial 配置
    if (isSerial) {
      if (!config.port || config.port.trim() === "") {
        alert("請輸入串列埠名稱");
        return;
      }
    }

    setLoading(true);
    try {
      const result = await connect(protocol, config);
      if (result && result.connection_id) {
        onConnectionChange(result.connection_id);
      } else {
        alert("連線成功，但未收到連線 ID");
      }
    } catch (error: any) {
      // 提取錯誤訊息
      let errorMessage = "連線失敗";
      if (error.response && error.response.data && error.response.data.error) {
        errorMessage = error.response.data.error;
      } else if (error.message) {
        errorMessage = error.message;
      } else if (typeof error === "string") {
        errorMessage = error;
      }
      alert(`連線失敗: ${errorMessage}`);
      console.error("連線錯誤詳情:", error);
    } finally {
      setLoading(false);
    }
  };

  /**
   * 處理斷線
   */
  const handleDisconnect = async () => {
    if (!connectionId) return;

    setLoading(true);
    try {
      await disconnect(connectionId);
      onConnectionChange(null);
    } catch (error: any) {
      alert(`斷線失敗: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 gap-4 animate-fade-in">
        {/* TCP/UDP 配置 */}
        {isTCP && (
          <div className="space-y-4">
            <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider border-b pb-2">網路設定</h3>
            <div className="grid grid-cols-2 gap-4">
              <div className="col-span-2 sm:col-span-1">
                <InputGroup label="主機位址 (IP)">
                  <StyledInput
                    type="text"
                    value={config.host || ""}
                    onChange={(e) => handleConfigChange("host", e.target.value)}
                    placeholder="192.168.1.100"
                  />
                </InputGroup>
              </div>
              <div className="col-span-2 sm:col-span-1">
                <InputGroup label="連接埠 (Port)">
                  <StyledInput
                    type="text"
                    inputMode="numeric"
                    value={config.port !== undefined ? String(config.port) : ""}
                    onChange={(e) => handleNumberChange("port", e.target.value)}
                    placeholder={protocol.includes("modbus") ? "502" : protocol.includes("fatek") ? "500" : "5000"}
                  />
                </InputGroup>
              </div>
            </div>

            {(protocol.includes("modbus") || protocol.includes("fatek")) && (
              <div className="col-span-2">
                <InputGroup label={protocol.includes("modbus") ? "單元 ID (Unit ID)" : "站號 (Station)"}>
                  <StyledInput
                    type="text"
                    inputMode="numeric"
                    value={
                      config.unitID !== undefined
                        ? String(config.unitID)
                        : config.station !== undefined
                        ? String(config.station)
                        : ""
                    }
                    onChange={(e) => {
                      const key = protocol.includes("modbus") ? "unitID" : "station";
                      handleNumberChange(key, e.target.value);
                    }}
                    placeholder="1"
                  />
                </InputGroup>
              </div>
            )}
          </div>
        )}

        {/* Serial 配置 */}
        {isSerial && (
          <div className="space-y-4">
            <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider border-b pb-2">串列埠設定</h3>
            <div className="col-span-2">
              <InputGroup label="串列埠名稱 (COM/TTY)">
                <StyledInput
                  type="text"
                  value={config.port || ""}
                  onChange={(e) => handleConfigChange("port", e.target.value)}
                  placeholder="COM3 或 /dev/ttyUSB0"
                />
              </InputGroup>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <InputGroup label="波特率 (Baud)">
                <StyledSelect
                  value={config.baudRate || 9600}
                  onChange={(e) => handleConfigChange("baudRate", parseInt(e.target.value))}
                >
                  <option value={9600}>9600</option>
                  <option value={19200}>19200</option>
                  <option value={38400}>38400</option>
                  <option value={57600}>57600</option>
                  <option value={115200}>115200</option>
                </StyledSelect>
              </InputGroup>

              <InputGroup label="資料位元 (Data Bits)">
                <StyledSelect
                  value={config.dataBits || (protocol.includes("fatek") || protocol.includes("mcprotocol") ? 7 : 8)}
                  onChange={(e) => handleConfigChange("dataBits", parseInt(e.target.value))}
                >
                  <option value={7}>7 Bits</option>
                  <option value={8}>8 Bits</option>
                </StyledSelect>
              </InputGroup>

              <InputGroup label="停止位元 (Stop Bits)">
                <StyledSelect
                  value={config.stopBits || (protocol.includes("mcprotocol") ? 2 : 1)}
                  onChange={(e) => handleConfigChange("stopBits", parseInt(e.target.value))}
                >
                  <option value={1}>1 Bit</option>
                  <option value={2}>2 Bits</option>
                </StyledSelect>
              </InputGroup>

              <InputGroup label="同位檢查 (Parity)">
                <StyledSelect
                  value={config.parity || (protocol.includes("fatek") || protocol.includes("mcprotocol") ? "E" : "N")}
                  onChange={(e) => handleConfigChange("parity", e.target.value)}
                >
                  <option value="N">None</option>
                  <option value="E">Even</option>
                  <option value="O">Odd</option>
                </StyledSelect>
              </InputGroup>
            </div>
            {(protocol.includes("modbus") || protocol.includes("fatek")) && (
              <div className="col-span-2 pt-2">
                <InputGroup label={protocol.includes("modbus") ? "單元 ID (Unit ID)" : "站號 (Station)"}>
                  <StyledInput
                    type="text"
                    inputMode="numeric"
                    value={
                      config.unitID !== undefined
                        ? String(config.unitID)
                        : config.station !== undefined
                        ? String(config.station)
                        : ""
                    }
                    onChange={(e) => {
                      const key = protocol.includes("modbus") ? "unitID" : "station";
                      handleNumberChange(key, e.target.value);
                    }}
                    placeholder="1"
                  />
                </InputGroup>
              </div>
            )}
          </div>
        )}
      </div>

      {/* 連線按鈕 */}
      <div className="pt-2">
        {!connectionId ? (
          <button
            onClick={handleConnect}
            disabled={loading}
            className="w-full bg-blue-600 text-white px-4 py-3 rounded-xl font-semibold shadow-lg shadow-blue-500/30 hover:bg-blue-700 hover:shadow-blue-500/40 active:scale-[0.98] transition-all disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? (
              <span className="flex items-center justify-center gap-2">
                <svg className="animate-spin h-5 w-5 text-white" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  ></path>
                </svg>
                連線中...
              </span>
            ) : (
              "建立連線"
            )}
          </button>
        ) : (
          <div className="space-y-3 animate-fade-in">
            <div className="flex items-center gap-2 p-3 bg-green-50 text-green-700 rounded-lg border border-green-200">
              <div className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></div>
              <span className="text-sm font-medium">已連線 ID: {connectionId.slice(0, 8)}...</span>
            </div>
            <div className="grid grid-cols-2 gap-3">
              <button
                onClick={handleDisconnect}
                disabled={loading}
                className="bg-white text-red-600 border border-red-200 px-4 py-3 rounded-xl font-semibold hover:bg-red-50 hover:border-red-300 active:scale-[0.98] transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {loading ? "斷線中..." : "斷開連線"}
              </button>
              {onMinimize && (
                <button
                  onClick={onMinimize}
                  disabled={loading}
                  className="bg-white text-gray-600 border border-gray-200 px-4 py-3 rounded-xl font-semibold hover:bg-gray-50 hover:border-gray-300 active:scale-[0.98] transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                >
                  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
                  </svg>
                  最小化
                </button>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
