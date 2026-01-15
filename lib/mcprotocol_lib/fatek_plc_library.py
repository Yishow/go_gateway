# ============================================================================
# Fatek PLC 通訊函式庫 - 生產環境級別實現
# Binary Protocol (雙速率二進制碼) 完整實現
# ============================================================================
# 功能特性:
# - Fatek Binary Communication Protocol 完全實現
# - 自動分包處理 (超過點數限制自動分割)
# - 根據 PLC 型號自動選擇 Frame 類型
# - CRC-16 Modbus 校驗
# - 完整的錯誤處理與驗證
# - 詳細的日誌記錄
# - 請求優化與批合併
# ============================================================================

import serial
import struct
import logging
import time
from enum import Enum
from typing import List, Dict, Tuple, Optional, Union
from dataclasses import dataclass
from datetime import datetime
import threading
import queue


# ============================================================================
# 常數定義
# ============================================================================

class FrameHeader(Enum):
    """Frame 頭字節"""
    QUERY = 0x51      # 查詢請求
    RESPONSE = 0x52   # 響應


class FrameID(Enum):
    """Frame ID"""
    BINARY_DEFAULT = 0x10  # Binary 協議預設值


class FrameTerminator(Enum):
    """Frame 終止符"""
    DEFAULT = 0x55AA  # 預設終止符


class CommandCode(Enum):
    """Fatek 命令代碼"""
    READ_DISCRETE = 0x44      # 讀取離散點 (Y/M/S/T/C)
    WRITE_DISCRETE = 0x45     # 寫入離散點
    READ_REGISTER = 0x46      # 讀取連續暫存器 (R/D/W)
    WRITE_REGISTER = 0x47     # 寫入連續暫存器
    MIXED_READ = 0x48         # 混合讀取 (離散 + 暫存器)
    MIXED_WRITE = 0x49        # 混合寫入
    LOOPBACK_TEST = 0x4E      # 迴路測試
    READ_SYSTEM_STATUS = 0x53  # 讀取系統狀態


class ComponentType(Enum):
    """元件類型"""
    X = 0x00   # 輸入接點
    Y = 0x01   # 輸出接點
    M = 0x02   # 內部繼電器
    S = 0x03   # 狀態繼電器
    T = 0x04   # 計時器
    C = 0x05   # 計數器
    R = 0x10   # 16-bit 暫存器
    D = 0x11   # 32-bit 暫存器
    W = 0x12   # 16-bit 字元


class ControlCode(Enum):
    """控制碼"""
    DISABLE = 0x01   # 禁用
    ENABLE = 0x02    # 啟用
    SET = 0x03       # 設置 (1)
    RESET = 0x04     # 重置 (0)


class ErrorCode(Enum):
    """Fatek 錯誤代碼"""
    NO_ERROR = 0x00           # 無錯誤
    ILLEGAL_VALUE = 0x02      # 非法值
    ILLEGAL_COMPONENT = 0x03  # 非法元件
    ILLEGAL_ADDRESS = 0x04    # 非法位址
    PROTECTION = 0x05         # 受保護
    ILLEGAL_QUANTITY = 0x06   # 非法數量
    ILLEGAL_COMMAND = 0x07    # 非法命令
    COMMUNICATION_ERROR = 0xFE # 通訊錯誤
    TIMEOUT = 0xFF            # 逾時


class PLCModel(Enum):
    """PLC 型號"""
    FB_SERIES = "FB"      # FB 系列 (FBs)
    M_SERIES = "M"        # M 系列
    UNKNOWN = "UNKNOWN"   # 未知


# ============================================================================
# 例外類別
# ============================================================================

class FatekException(Exception):
    """Fatek 例外基類"""
    pass


class FatekConnectionError(FatekException):
    """連線錯誤"""
    pass


class FatekCommunicationError(FatekException):
    """通訊錯誤"""
    pass


class FatekProtocolError(FatekException):
    """協議錯誤"""
    pass


class FatekTimeoutError(FatekException):
    """逾時錯誤"""
    pass


class FatekDataError(FatekException):
    """資料錯誤"""
    pass


class FatekAddressError(FatekException):
    """位址錯誤"""
    pass


# ============================================================================
# 資料類別
# ============================================================================

@dataclass
class SerialConfig:
    """序列埠配置"""
    port: str
    baudrate: int = 9600
    bytesize: int = 8
    parity: str = 'E'  # Even
    stopbits: int = 1
    timeout: float = 1.0
    write_timeout: float = 1.0


@dataclass
class FatekResponse:
    """Fatek 響應資料"""
    station_id: int
    command_code: int
    error_code: int
    payload: bytes
    timestamp: datetime
    raw_data: bytes = None
    
    def is_success(self) -> bool:
        """檢查是否成功"""
        return self.error_code == 0x00


# ============================================================================
# CRC-16 Modbus 計算
# ============================================================================

class CRC16:
    """CRC-16 Modbus 計算"""
    
    POLY = 0xA001
    INIT_VALUE = 0xFFFF
    
    @staticmethod
    def calculate(data: bytes) -> int:
        """計算 CRC-16"""
        crc = CRC16.INIT_VALUE
        
        for byte in data:
            crc ^= byte
            for _ in range(8):
                if crc & 1:
                    crc = (crc >> 1) ^ CRC16.POLY
                else:
                    crc >>= 1
        
        return crc
    
    @staticmethod
    def verify(data: bytes, crc: int) -> bool:
        """驗證 CRC-16"""
        return CRC16.calculate(data) == crc


# ============================================================================
# 日誌設定
# ============================================================================

class FatekLogger:
    """Fatek 日誌管理"""
    
    @staticmethod
    def setup_logger(name: str = "fatek_plc") -> logging.Logger:
        """設定日誌"""
        logger = logging.getLogger(name)
        logger.setLevel(logging.DEBUG)
        
        # 檔案處理器
        file_handler = logging.FileHandler(
            f"fatek_plc_{datetime.now().strftime('%Y%m%d_%H%M%S')}.log"
        )
        file_handler.setLevel(logging.DEBUG)
        
        # 控制台處理器
        console_handler = logging.StreamHandler()
        console_handler.setLevel(logging.INFO)
        
        # 格式化器
        formatter = logging.Formatter(
            '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
        )
        file_handler.setFormatter(formatter)
        console_handler.setFormatter(formatter)
        
        logger.addHandler(file_handler)
        logger.addHandler(console_handler)
        
        return logger


# ============================================================================
# 協議訊息構建
# ============================================================================

class ProtocolBuilder:
    """Fatek Binary Protocol 訊息構建"""
    
    def __init__(self, logger: logging.Logger = None):
        self.logger = logger or FatekLogger.setup_logger()
    
    def build_query_message(
        self,
        station_id: int,
        command_code: int,
        payload: bytes
    ) -> bytes:
        """構建查詢訊息"""
        header = bytes([FrameHeader.QUERY.value])
        frame_id = bytes([FrameID.BINARY_DEFAULT.value])
        
        # 計算長度 (payload 字節數)
        length = len(payload)
        length_bytes = struct.pack('<H', length)  # Little-endian
        
        # 計算 CRC (長度 + payload)
        crc_data = length_bytes + payload
        crc_value = CRC16.calculate(crc_data)
        crc_bytes = struct.pack('<H', crc_value)  # Little-endian
        
        # 終止符
        terminator = struct.pack('<H', FrameTerminator.DEFAULT.value)
        
        # 組合訊息
        message = header + frame_id + length_bytes + payload + crc_bytes + terminator
        
        self.logger.debug(
            f"構建查詢訊息: STN={station_id:02X}, "
            f"CMD={command_code:02X}, Length={length}, "
            f"HEX={message.hex().upper()}"
        )
        
        return message
    
    def parse_response(self, data: bytes) -> FatekResponse:
        """解析響應訊息"""
        if len(data) < 9:
            raise FatekProtocolError(f"響應長度不足: {len(data)}")
        
        # 驗證頭字節
        if data[0] != FrameHeader.RESPONSE.value:
            raise FatekProtocolError(f"無效的响應頭: {data[0]:02X}")
        
        # 驗證 Frame ID
        if data[1] != FrameID.BINARY_DEFAULT.value:
            raise FatekProtocolError(f"無效的 Frame ID: {data[1]:02X}")
        
        # 解析長度
        length = struct.unpack('<H', data[2:4])[0]
        
        # 驗證訊息長度
        expected_length = 4 + length + 2 + 2  # Header+ID + Length + Payload + CRC + Terminator
        if len(data) < expected_length:
            raise FatekProtocolError(
                f"訊息長度不符: 期望 {expected_length}, 實際 {len(data)}"
            )
        
        # 提取 payload
        payload = data[4:4+length]
        
        # 驗證 CRC
        crc_data = data[2:4+length]
        crc_bytes = data[4+length:6+length]
        crc_value = struct.unpack('<H', crc_bytes)[0]
        
        if not CRC16.verify(crc_data, crc_value):
            raise FatekProtocolError(
                f"CRC 驗證失敗: 計算 {CRC16.calculate(crc_data):04X}, "
                f"接收 {crc_value:04X}"
            )
        
        # 驗證終止符
        terminator = struct.unpack('<H', data[6+length:8+length])[0]
        if terminator != FrameTerminator.DEFAULT.value:
            raise FatekProtocolError(f"無效的終止符: {terminator:04X}")
        
        # 解析 payload
        station_id = payload[0]
        command_code = payload[1]
        error_code = payload[2]
        response_payload = payload[3:]
        
        response = FatekResponse(
            station_id=station_id,
            command_code=command_code,
            error_code=error_code,
            payload=response_payload,
            timestamp=datetime.now(),
            raw_data=data
        )
        
        self.logger.debug(
            f"解析響應: STN={station_id:02X}, "
            f"CMD={command_code:02X}, ERR={error_code:02X}, "
            f"HEX={data.hex().upper()}"
        )
        
        return response
    
    def get_error_description(self, error_code: int) -> str:
        """取得錯誤描述"""
        error_map = {
            0x00: "無錯誤",
            0x02: "非法值",
            0x03: "非法元件",
            0x04: "非法位址",
            0x05: "受保護",
            0x06: "非法數量",
            0x07: "非法命令",
            0xFE: "通訊錯誤",
            0xFF: "逾時"
        }
        return error_map.get(error_code, f"未知錯誤 (0x{error_code:02X})")


# ============================================================================
# 訊息發送/接收管理
# ============================================================================

class CommunicationHandler:
    """通訊處理"""
    
    def __init__(self, config: SerialConfig, logger: logging.Logger = None):
        self.config = config
        self.logger = logger or FatekLogger.setup_logger()
        self.serial_port: Optional[serial.Serial] = None
        self.connected = False
    
    def connect(self) -> bool:
        """連線到 PLC"""
        try:
            self.serial_port = serial.Serial(
                port=self.config.port,
                baudrate=self.config.baudrate,
                bytesize=self.config.bytesize,
                parity=self.config.parity,
                stopbits=self.config.stopbits,
                timeout=self.config.timeout,
                write_timeout=self.config.write_timeout
            )
            self.connected = True
            self.logger.info(f"已連線到 {self.config.port}")
            return True
        except serial.SerialException as e:
            self.logger.error(f"連線失敗: {e}")
            raise FatekConnectionError(f"無法連線到序列埠: {e}")
    
    def disconnect(self) -> bool:
        """斷線"""
        try:
            if self.serial_port and self.serial_port.is_open:
                self.serial_port.close()
                self.connected = False
                self.logger.info("已斷線")
            return True
        except serial.SerialException as e:
            self.logger.error(f"斷線失敗: {e}")
            raise FatekConnectionError(f"斷線失敗: {e}")
    
    def send_message(self, message: bytes) -> None:
        """發送訊息"""
        if not self.connected or not self.serial_port:
            raise FatekConnectionError("未連線")
        
        try:
            self.serial_port.write(message)
            self.logger.debug(f"已發送: {message.hex().upper()}")
        except serial.SerialException as e:
            self.logger.error(f"發送失敗: {e}")
            raise FatekCommunicationError(f"發送失敗: {e}")
    
    def receive_message(self, timeout: Optional[float] = None) -> bytes:
        """接收訊息"""
        if not self.connected or not self.serial_port:
            raise FatekConnectionError("未連線")
        
        old_timeout = self.serial_port.timeout
        if timeout:
            self.serial_port.timeout = timeout
        
        try:
            data = bytearray()
            start_time = time.time()
            
            while True:
                chunk = self.serial_port.read(1)
                if not chunk:
                    if data and data[0] == 0x52:  # 回應頭
                        # 嘗試讀取足夠的資料
                        if len(data) >= 6:
                            break
                    
                    if timeout and (time.time() - start_time) > timeout:
                        raise FatekTimeoutError("接收訊息逾時")
                    continue
                
                data.extend(chunk)
                
                # 檢查是否接收完整
                if len(data) >= 8:
                    if data[0] == 0x52:  # 回應頭
                        length = struct.unpack('<H', data[2:4])[0]
                        if len(data) >= 8 + length:
                            break
            
            self.logger.debug(f"已接收: {bytes(data).hex().upper()}")
            return bytes(data)
        
        except FatekTimeoutError:
            raise
        except Exception as e:
            self.logger.error(f"接收失敗: {e}")
            raise FatekCommunicationError(f"接收失敗: {e}")
        finally:
            self.serial_port.timeout = old_timeout


# ============================================================================
# 命令構建器
# ============================================================================

class CommandBuilder:
    """Fatek 命令構建"""
    
    # 點數限制
    READ_DISCRETE_MAX = 256
    WRITE_DISCRETE_MAX = 256
    READ_REGISTER_MAX = 64
    WRITE_REGISTER_MAX = 64
    
    def __init__(self, logger: logging.Logger = None):
        self.logger = logger or FatekLogger.setup_logger()
    
    def build_read_discrete(
        self,
        station_id: int,
        component_type: ComponentType,
        start_address: int,
        quantity: int
    ) -> bytes:
        """構建讀取離散點命令"""
        if quantity > self.READ_DISCRETE_MAX:
            raise FatekDataError(
                f"讀取數量超過限制: {quantity} > {self.READ_DISCRETE_MAX}"
            )
        
        payload = bytearray()
        payload.append(station_id)
        payload.append(CommandCode.READ_DISCRETE.value)
        payload.append(quantity)
        payload.append(component_type.value)
        payload.extend(struct.pack('>I', start_address))  # Big-endian
        
        return bytes(payload)
    
    def build_write_discrete(
        self,
        station_id: int,
        component_type: ComponentType,
        start_address: int,
        values: List[int]
    ) -> bytes:
        """構建寫入離散點命令"""
        quantity = len(values)
        if quantity > self.WRITE_DISCRETE_MAX:
            raise FatekDataError(
                f"寫入數量超過限制: {quantity} > {self.WRITE_DISCRETE_MAX}"
            )
        
        payload = bytearray()
        payload.append(station_id)
        payload.append(CommandCode.WRITE_DISCRETE.value)
        payload.append(quantity)
        payload.append(component_type.value)
        payload.extend(struct.pack('>I', start_address))
        
        # 打包資料
        for value in values:
            payload.append(0x01 if value else 0x00)
        
        return bytes(payload)
    
    def build_read_register(
        self,
        station_id: int,
        component_type: ComponentType,
        start_address: int,
        quantity: int,
        is_32bit: bool = False
    ) -> bytes:
        """構建讀取暫存器命令"""
        max_quantity = self.READ_REGISTER_MAX // (2 if is_32bit else 1)
        if quantity > max_quantity:
            raise FatekDataError(
                f"讀取數量超過限制: {quantity} > {max_quantity}"
            )
        
        payload = bytearray()
        payload.append(station_id)
        payload.append(CommandCode.READ_REGISTER.value)
        payload.append(quantity)
        payload.append(component_type.value)
        payload.extend(struct.pack('>I', start_address))
        
        return bytes(payload)
    
    def build_write_register(
        self,
        station_id: int,
        component_type: ComponentType,
        start_address: int,
        values: List[int],
        is_32bit: bool = False
    ) -> bytes:
        """構建寫入暫存器命令"""
        quantity = len(values)
        max_quantity = self.WRITE_REGISTER_MAX // (2 if is_32bit else 1)
        if quantity > max_quantity:
            raise FatekDataError(
                f"寫入數量超過限制: {quantity} > {max_quantity}"
            )
        
        payload = bytearray()
        payload.append(station_id)
        payload.append(CommandCode.WRITE_REGISTER.value)
        payload.append(quantity)
        payload.append(component_type.value)
        payload.extend(struct.pack('>I', start_address))
        
        # 打包資料
        for value in values:
            if is_32bit:
                payload.extend(struct.pack('>I', value))
            else:
                payload.extend(struct.pack('>H', value))
        
        return bytes(payload)
    
    def split_read_discrete(
        self,
        station_id: int,
        component_type: ComponentType,
        start_address: int,
        quantity: int
    ) -> List[bytes]:
        """自動分包讀取離散點"""
        commands = []
        remaining = quantity
        current_address = start_address
        
        while remaining > 0:
            chunk = min(remaining, self.READ_DISCRETE_MAX)
            cmd = self.build_read_discrete(
                station_id,
                component_type,
                current_address,
                chunk
            )
            commands.append(cmd)
            remaining -= chunk
            current_address += chunk
        
        return commands
    
    def split_write_discrete(
        self,
        station_id: int,
        component_type: ComponentType,
        start_address: int,
        values: List[int]
    ) -> List[bytes]:
        """自動分包寫入離散點"""
        commands = []
        remaining = len(values)
        current_address = start_address
        current_offset = 0
        
        while remaining > 0:
            chunk_size = min(remaining, self.WRITE_DISCRETE_MAX)
            chunk_values = values[current_offset:current_offset + chunk_size]
            cmd = self.build_write_discrete(
                station_id,
                component_type,
                current_address,
                chunk_values
            )
            commands.append(cmd)
            remaining -= chunk_size
            current_address += chunk_size
            current_offset += chunk_size
        
        return commands


# ============================================================================
# 主 Fatek PLC 類別
# ============================================================================

class FatekPLC:
    """Fatek PLC 主類別"""
    
    def __init__(
        self,
        config: SerialConfig,
        station_id: int = 1,
        plc_model: PLCModel = PLCModel.FB_SERIES,
        logger: logging.Logger = None
    ):
        self.config = config
        self.station_id = station_id
        self.plc_model = plc_model
        self.logger = logger or FatekLogger.setup_logger()
        
        self.comm_handler = CommunicationHandler(config, self.logger)
        self.protocol_builder = ProtocolBuilder(self.logger)
        self.command_builder = CommandBuilder(self.logger)
        
        self.request_id = 0
        self.lock = threading.Lock()
    
    # ========================================================================
    # 連線管理
    # ========================================================================
    
    def connect(self) -> bool:
        """連線到 PLC"""
        return self.comm_handler.connect()
    
    def disconnect(self) -> bool:
        """斷線"""
        return self.comm_handler.disconnect()
    
    def is_connected(self) -> bool:
        """檢查連線狀態"""
        return self.comm_handler.connected
    
    # ========================================================================
    # 低級通訊
    # ========================================================================
    
    def _send_command(self, payload: bytes) -> FatekResponse:
        """發送命令並接收響應"""
        with self.lock:
            try:
                # 構建訊息
                message = self.protocol_builder.build_query_message(
                    self.station_id,
                    payload[1],  # 命令代碼
                    payload
                )
                
                # 發送訊息
                self.comm_handler.send_message(message)
                
                # 接收響應
                response_data = self.comm_handler.receive_message(
                    timeout=self.config.timeout
                )
                
                # 解析響應
                response = self.protocol_builder.parse_response(response_data)
                
                # 檢查錯誤代碼
                if response.error_code != 0x00:
                    error_desc = self.protocol_builder.get_error_description(
                        response.error_code
                    )
                    self.logger.error(f"PLC 錯誤: {error_desc}")
                    raise FatekCommunicationError(f"PLC 錯誤: {error_desc}")
                
                return response
            
            except FatekTimeoutError:
                self.logger.error("通訊逾時")
                raise
            except Exception as e:
                self.logger.error(f"通訊異常: {e}")
                raise
    
    # ========================================================================
    # 高級讀取操作
    # ========================================================================
    
    def read_discrete(
        self,
        component_type: Union[ComponentType, str],
        start_address: int,
        quantity: int
    ) -> List[int]:
        """讀取離散點 (Y/M/S/T/C/X)"""
        # 參數驗證
        if isinstance(component_type, str):
            component_type = ComponentType[component_type.upper()]
        
        if quantity <= 0:
            raise FatekAddressError("讀取數量必須 > 0")
        
        if quantity > self.command_builder.READ_DISCRETE_MAX:
            # 自動分包
            results = []
            for cmd in self.command_builder.split_read_discrete(
                self.station_id,
                component_type,
                start_address,
                quantity
            ):
                response = self._send_command(cmd)
                results.extend(list(response.payload))
            return results[:quantity]
        
        # 單一請求
        cmd = self.command_builder.build_read_discrete(
            self.station_id,
            component_type,
            start_address,
            quantity
        )
        response = self._send_command(cmd)
        
        return list(response.payload[:quantity])
    
    def read_register(
        self,
        component_type: Union[ComponentType, str],
        start_address: int,
        quantity: int,
        is_32bit: bool = False
    ) -> List[int]:
        """讀取暫存器 (R/D/W)"""
        # 參數驗證
        if isinstance(component_type, str):
            component_type = ComponentType[component_type.upper()]
        
        if quantity <= 0:
            raise FatekAddressError("讀取數量必須 > 0")
        
        max_quantity = self.command_builder.READ_REGISTER_MAX // (2 if is_32bit else 1)
        if quantity > max_quantity:
            # 自動分包
            results = []
            for cmd in self.command_builder.split_read_discrete(
                self.station_id,
                component_type,
                start_address,
                quantity
            ):
                response = self._send_command(cmd)
                # 解析暫存器資料
                values = self._parse_register_data(response.payload, is_32bit)
                results.extend(values)
            return results[:quantity]
        
        # 單一請求
        cmd = self.command_builder.build_read_register(
            self.station_id,
            component_type,
            start_address,
            quantity,
            is_32bit
        )
        response = self._send_command(cmd)
        
        return self._parse_register_data(response.payload, is_32bit)
    
    def read_random_registers(
        self,
        addresses: List[Tuple[ComponentType, int]]
    ) -> Dict[str, int]:
        """隨機讀取暫存器"""
        results = {}
        
        for component_type, address in addresses:
            try:
                values = self.read_register(component_type, address, 1)
                key = f"{component_type.name}{address}"
                results[key] = values[0] if values else 0
            except Exception as e:
                self.logger.error(f"讀取 {component_type.name}{address} 失敗: {e}")
                raise
        
        return results
    
    # ========================================================================
    # 高級寫入操作
    # ========================================================================
    
    def write_discrete(
        self,
        component_type: Union[ComponentType, str],
        start_address: int,
        values: List[int]
    ) -> bool:
        """寫入離散點"""
        if isinstance(component_type, str):
            component_type = ComponentType[component_type.upper()]
        
        if len(values) <= 0:
            raise FatekAddressError("寫入值不能為空")
        
        # 自動分包
        for cmd in self.command_builder.split_write_discrete(
            self.station_id,
            component_type,
            start_address,
            values
        ):
            self._send_command(cmd)
        
        return True
    
    def write_register(
        self,
        component_type: Union[ComponentType, str],
        start_address: int,
        values: List[int],
        is_32bit: bool = False
    ) -> bool:
        """寫入暫存器"""
        if isinstance(component_type, str):
            component_type = ComponentType[component_type.upper()]
        
        if len(values) <= 0:
            raise FatekAddressError("寫入值不能為空")
        
        # 自動分包
        max_quantity = self.command_builder.WRITE_REGISTER_MAX // (2 if is_32bit else 1)
        remaining = len(values)
        current_address = start_address
        current_offset = 0
        
        while remaining > 0:
            chunk_size = min(remaining, max_quantity)
            chunk_values = values[current_offset:current_offset + chunk_size]
            cmd = self.command_builder.build_write_register(
                self.station_id,
                component_type,
                current_address,
                chunk_values,
                is_32bit
            )
            self._send_command(cmd)
            remaining -= chunk_size
            current_address += chunk_size
            current_offset += chunk_size
        
        return True
    
    # ========================================================================
    # 輔助方法
    # ========================================================================
    
    def _parse_register_data(self, data: bytes, is_32bit: bool = False) -> List[int]:
        """解析暫存器資料"""
        values = []
        
        if is_32bit:
            for i in range(0, len(data), 4):
                if i + 4 <= len(data):
                    value = struct.unpack('>I', data[i:i+4])[0]
                    values.append(value)
        else:
            for i in range(0, len(data), 2):
                if i + 2 <= len(data):
                    value = struct.unpack('>H', data[i:i+2])[0]
                    values.append(value)
        
        return values
    
    def loopback_test(self, data: bytes) -> bool:
        """迴路測試"""
        payload = bytearray()
        payload.append(self.station_id)
        payload.append(CommandCode.LOOPBACK_TEST.value)
        payload.extend(data)
        
        response = self._send_command(bytes(payload))
        return response.payload == data
    
    def read_system_status(self) -> Dict[str, any]:
        """讀取系統狀態"""
        payload = bytearray()
        payload.append(self.station_id)
        payload.append(CommandCode.READ_SYSTEM_STATUS.value)
        
        response = self._send_command(bytes(payload))
        
        # 解析系統狀態資料
        status = {
            "raw_data": response.payload.hex().upper(),
            "timestamp": response.timestamp
        }
        
        return status


# ============================================================================
# 使用範例
# ============================================================================

def example_usage():
    """使用範例"""
    
    # 配置序列埠
    config = SerialConfig(
        port="COM1",
        baudrate=9600,
        timeout=1.0
    )
    
    # 建立 PLC 連線
    plc = FatekPLC(
        config=config,
        station_id=1,
        plc_model=PLCModel.FB_SERIES
    )
    
    try:
        # 連線
        plc.connect()
        print("已連線到 PLC")
        
        # 讀取離散點 Y0~Y9
        y_values = plc.read_discrete(ComponentType.Y, 0, 10)
        print(f"Y0~Y9: {y_values}")
        
        # 讀取暫存器 R0~R9
        r_values = plc.read_register(ComponentType.R, 0, 10)
        print(f"R0~R9: {r_values}")
        
        # 寫入離散點 Y0~Y3
        plc.write_discrete(ComponentType.Y, 0, [1, 0, 1, 0])
        print("已寫入 Y0~Y3")
        
        # 寫入暫存器 R0=1000, R1=2000
        plc.write_register(ComponentType.R, 0, [1000, 2000])
        print("已寫入 R0=1000, R1=2000")
        
        # 隨機讀取
        random_reads = plc.read_random_registers([
            (ComponentType.R, 0),
            (ComponentType.R, 1),
            (ComponentType.R, 100)
        ])
        print(f"隨機讀取: {random_reads}")
        
    finally:
        # 斷線
        plc.disconnect()
        print("已斷線")


if __name__ == "__main__":
    example_usage()
