import socket
import struct
import time
from typing import List, Union, Tuple, Dict, Optional

class MCProtocolError(Exception):
    """MC Protocol 通訊異常基類"""
    pass

class MCConnectionError(MCProtocolError):
    """連線相關異常"""
    pass

class MCResponseError(MCProtocolError):
    """PLC 回應錯誤異常 (包含錯誤碼)"""
    def __init__(self, error_code: int, message: str):
        self.error_code = error_code
        super().__init__(f"{message} (Error Code: 0x{error_code:04X})")

class MCProtocol:
    """
    Mitsubishi MELSEC MC Protocol (3E Binary Frame) 用戶端實作。
    
    支援功能：
    - 批量讀取位元 (Batch Read Bit)
    - 批量寫入位元 (Batch Write Bit)
    - 批量讀取字組 (Batch Read Word)
    - 批量寫入字組 (Batch Write Word)
    - 隨機讀取字組 (Random Read Word)
    - 連線狀態檢查 (Check Connection)
    - 自動重連機制
    
    支援的軟元件 (Device) 代碼：
    - 'D': 資料暫存器 (Data Register)
    - 'W': 連結暫存器 (Link Register)
    - 'R': 檔案暫存器 (File Register)
    - 'M': 內部繼電器 (Internal Relay)
    - 'X': 輸入 (Input)
    - 'Y': 輸出 (Output)
    - 'L': 鎖存繼電器 (Latch Relay)
    - 'B': 連結繼電器 (Link Relay)
    - 'F': 報警器 (Annunciator)
    - 'TN': 計時器 (目前值)
    - 'CN': 計數器 (目前值)
    - 'SB': 特殊連結繼電器
    - 'SW': 特殊連結暫存器
    - 'S': 步進繼電器
    - 'TS': 計時器 (觸點)
    - 'TC': 計時器 (線圈)
    - 'CS': 計數器 (觸點)
    - 'CC': 計數器 (線圈)
    """

    # 軟元件類型映射 (Device Name -> (Device Code, Is Bit Device?))
    # 參考: MELSEC Communication Protocol Reference Manual
    DEVICE_MAP: Dict[str, Tuple[int, bool]] = {
        'D':  (0xA8, False), # Word
        'W':  (0xB4, False), # Word
        'R':  (0xAF, False), # Word
        'M':  (0x90, True),  # Bit
        'X':  (0x9C, True),  # Bit
        'Y':  (0x9D, True),  # Bit
        'L':  (0x92, True),  # Bit
        'B':  (0xA0, True),  # Bit
        'F':  (0x93, True),  # Bit
        'TN': (0xC2, False), # Word (Timer Current Value)
        'CN': (0xC5, False), # Word (Counter Current Value)
        'SB': (0xA1, True),  # Bit
        'SW': (0xB5, False), # Word
        'S':  (0x98, True),  # Bit
        'TS': (0xC1, True),  # Bit (Timer Contact)
        'TC': (0xC0, True),  # Bit (Timer Coil)
        'CS': (0xC4, True),  # Bit (Counter Contact)
        'CC': (0xC3, True),  # Bit (Counter Coil)
        'Z':  (0xCC, False), # Word (Index Register)
    }

    def __init__(self, host: str, port: int, network_no: int = 0, pc_no: int = 255, timeout: float = 2.0):
        """
        初始化 MC Protocol 用戶端。

        Args:
            host (str): PLC IP 位址。
            port (int): PLC MC Protocol 埠號 (通常為 5000 或 5002)。
            network_no (int): 網路編號 (預設 0)。
            pc_no (int): PC 編號 (預設 0xFF/255)。
            timeout (float): 通訊超時時間 (秒)。
        """
        self.host = host
        self.port = port
        self.network_no = network_no
        self.pc_no = pc_no
        self.timeout = timeout
        self.sock: Optional[socket.socket] = None
        self.is_connected = False

    def connect(self):
        """建立 TCP 連線。"""
        try:
            if self.sock:
                self.sock.close()
            self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.sock.settimeout(self.timeout)
            self.sock.connect((self.host, self.port))
            self.is_connected = True
        except socket.error as e:
            self.is_connected = False
            raise MCConnectionError(f"無法連線至 PLC ({self.host}:{self.port}): {e}")

    def close(self):
        """關閉連線。"""
        if self.sock:
            try:
                self.sock.close()
            except socket.error:
                pass
        self.sock = None
        self.is_connected = False
    
    def check_connection(self) -> bool:
        """
        檢查目前連線狀態。
        會嘗試發送一個輕量級的讀取請求 (讀取 D0) 來驗證連線是否有效。
        
        Returns:
            bool: 若連線有效則回傳 True，否則 False。
        """
        if not self.is_connected or self.sock is None:
            return False
            
        try:
            # 嘗試讀取 D0 一個 Word 來測試連線
            # 注意：這裡不調用 batch_read_word 是為了避免遞迴調用 _ensure_connection
            # 但為了簡單與一致性，我們直接使用 read_d 並捕捉異常
            self.read_d(0, 1)
            return True
        except Exception:
            self.is_connected = False
            return False

    def _ensure_connection(self):
        """確保連線存在，若斷線則嘗試重連。"""
        if not self.is_connected or self.sock is None:
            self.connect()

    def _build_header(self, data_len: int) -> bytes:
        """
        建立 3E Frame 請求標頭 (Request Header)。
        Structure:
        - Subheader (2 bytes): 0x5000 (Request)
        - Network No (1 byte)
        - PC No (1 byte)
        - Request Dest Module IO No (2 bytes): 0x03FF (Own station)
        - Request Dest Module Station No (1 byte): 0x00
        - Data Length (2 bytes): Length of following data
        - CPU Monitoring Timer (2 bytes): 0x0010 (Example)
        """
        header = bytearray()
        header += b'\x50\x00'                  # Subheader (Request)
        header += struct.pack('B', self.network_no)
        header += struct.pack('B', self.pc_no)
        header += b'\xFF\x03'                  # IO No (Self)
        header += b'\x00'                      # Station No (Self)
        header += struct.pack('<H', data_len + 2) # Data Length (Timer + Command + Subcommand + Data)
        header += b'\x10\x00'                  # CPU Timer (Low byte first)
        return header

    def _parse_device_code(self, device: str) -> Tuple[int, bool]:
        """解析軟元件代碼。"""
        device = device.upper()
        if device not in self.DEVICE_MAP:
            raise ValueError(f"不支援的軟元件類型: {device}")
        return self.DEVICE_MAP[device]

    def _send_command(self, command: bytes, subcommand: bytes, data: bytes) -> bytes:
        """
        發送指令並接收回應。
        
        Args:
            command: 指令碼 (2 bytes)
            subcommand: 子指令碼 (2 bytes)
            data: 請求數據內容
            
        Returns:
            bytes: 回應數據內容 (不含標頭與結束碼)
        """
        self._ensure_connection()

        # 構建完整封包
        # 數據長度 = Timer(2) + Command(2) + Subcommand(2) + Data(N)
        payload_len = 2 + 2 + 2 + len(data)
        header = b'\x50\x00'  # Subheader (Request)
        header += struct.pack('B', self.network_no)
        header += struct.pack('B', self.pc_no)
        header += b'\xFF\x03' # IO No
        header += b'\x00'     # Station No
        header += struct.pack('<H', payload_len) # Data Length
        header += b'\x10\x00' # CPU Monitoring Timer (0x0010 = 4s approximately, unit is 250ms)
        
        request_packet = header + command + subcommand + data

        try:
            self.sock.sendall(request_packet)
            
            # 接收回應標頭 (11 bytes)
            # Subheader(2) + Net(1) + PC(1) + IO(2) + Station(1) + Len(2) + EndCode(2)
            # 注意: 成功回應標頭為 0xD000 (Subheader) + ... + Len(2) + EndCode(2) + Data
            # 失敗回應標頭略有不同，但 EndCode 位於固定位置
            
            # 讀取前 9 個位元組以獲取數據長度
            # Response: 0xD000(2) + Net(1) + PC(1) + IO(2) + Station(1) + Len(2)
            header_res = self._recv_exact(9)
            
            # 檢查 Subheader
            if header_res[:2] != b'\xD0\x00':
                raise MCProtocolError(f"無效的回應標頭: {header_res[:2].hex()}")
            
            data_len = struct.unpack('<H', header_res[7:9])[0]
            
            # 讀取剩餘數據 (EndCode + Data)
            body_res = self._recv_exact(data_len)
            
            end_code = struct.unpack('<H', body_res[:2])[0]
            
            if end_code != 0:
                raise MCResponseError(end_code, "PLC 回傳錯誤代碼")
            
            return body_res[2:] # 返回實際數據部分

        except (socket.timeout, socket.error) as e:
            self.is_connected = False
            raise MCConnectionError(f"通訊錯誤: {e}")

    def _recv_exact(self, length: int) -> bytes:
        """從 socket 接收指定長度的數據。"""
        data = b''
        while len(data) < length:
            chunk = self.sock.recv(length - len(data))
            if not chunk:
                raise MCConnectionError("連線被遠端關閉")
            data += chunk
        return data

    def batch_read_word(self, device_type: str, head_device: int, count: int) -> List[int]:
        """
        批量讀取字組數據 (Word)。
        
        Args:
            device_type (str): 軟元件類型 (如 'D', 'W')。
            head_device (int): 起始位址。
            count (int): 讀取數量 (Word 數)。
            
        Returns:
            List[int]: 讀取到的數值列表 (16-bit unsigned integer)。
        """
        dev_code, is_bit = self._parse_device_code(device_type)
        
        # Command: 0401 (Batch Read)
        # Subcommand: 0000 (Word Access)
        command = b'\x01\x04'
        subcommand = b'\x00\x00'
        
        # Request Data:
        # Head Device No (3 bytes)
        # Device Code (1 byte)
        # Number of devices (2 bytes)
        
        # 3 bytes Head Device (Little Endian)
        head_dev_bytes = struct.pack('<I', head_device)[:3]
        data = head_dev_bytes + struct.pack('B', dev_code) + struct.pack('<H', count)
        
        response = self._send_command(command, subcommand, data)
        
        # Parse Response (Each word is 2 bytes)
        # data length should be count * 2
        if len(response) != count * 2:
            raise MCProtocolError(f"回應長度不符: 預期 {count*2}, 實際 {len(response)}")
            
        # Unpack as unsigned short (H)
        format_str = f'<{count}H'
        values = struct.unpack(format_str, response)
        return list(values)

    def batch_write_word(self, device_type: str, head_device: int, values: List[int]):
        """
        批量寫入字組數據 (Word)。
        
        Args:
            device_type (str): 軟元件類型 (如 'D', 'W')。
            head_device (int): 起始位址。
            values (List[int]): 要寫入的數值列表 (16-bit int)。
        """
        dev_code, is_bit = self._parse_device_code(device_type)
        count = len(values)
        
        # Command: 1401 (Batch Write)
        # Subcommand: 0000 (Word Access)
        command = b'\x01\x14'
        subcommand = b'\x00\x00'
        
        head_dev_bytes = struct.pack('<I', head_device)[:3]
        data_header = head_dev_bytes + struct.pack('B', dev_code) + struct.pack('<H', count)
        
        # Pack values
        data_body = struct.pack(f'<{count}H', *values)
        
        self._send_command(command, subcommand, data_header + data_body)

    def batch_read_bit(self, device_type: str, head_device: int, count: int) -> List[int]:
        """
        批量讀取位元數據 (Bit)。
        
        Args:
            device_type (str): 軟元件類型 (如 'M', 'X', 'Y')。
            head_device (int): 起始位址。
            count (int): 讀取數量 (Bit 數)。
            
        Returns:
            List[int]: 讀取到的狀態列表 (0 或 1)。
        """
        dev_code, is_bit = self._parse_device_code(device_type)
        if not is_bit:
            raise ValueError(f"設備類型 {device_type} 不是位元設備 (Bit Device)") 
        
        # Command: 0401 (Batch Read)
        # Subcommand: 0001 (Bit Access)
        command = b'\x01\x04'
        subcommand = b'\x01\x00' # 0001 Little Endian
        
        head_dev_bytes = struct.pack('<I', head_device)[:3]
        data = head_dev_bytes + struct.pack('B', dev_code) + struct.pack('<H', count)
        
        response = self._send_command(command, subcommand, data)
        
        expected_bytes = (count + 1) // 2
        if len(response) != expected_bytes:
            raise MCProtocolError(f"回應長度不符: 預期 {expected_bytes}, 實際 {len(response)}")

        result = []
        for i in range(count):
            byte_idx = i // 2
            # High nibble is first device (i=0), Low nibble is second device (i=1)
            # 與寫入邏輯一致：val1 << 4 | val2，其中 val1 是 i，val2 是 i+1
            is_high_nibble = (i % 2) == 0
            byte_val = response[byte_idx]
            
            if is_high_nibble:
                val = (byte_val >> 4) & 0x0F
            else:
                val = byte_val & 0x0F
            
            result.append(1 if val == 1 else 0)
            
        return result

    def batch_write_bit(self, device_type: str, head_device: int, values: List[int]):
        """
        批量寫入位元數據 (Bit)。
        
        Args:
            device_type (str): 軟元件類型 (如 'M', 'X', 'Y')。
            head_device (int): 起始位址。
            values (List[int]): 要寫入的狀態列表 (0 或 1)。
        """
        dev_code, is_bit = self._parse_device_code(device_type)
        count = len(values)
        
        # Command: 1401 (Batch Write)
        # Subcommand: 0001 (Bit Access)
        command = b'\x01\x14'
        subcommand = b'\x01\x00'
        
        head_dev_bytes = struct.pack('<I', head_device)[:3]
        data_header = head_dev_bytes + struct.pack('B', dev_code) + struct.pack('<H', count)
        
        data_bytes = bytearray()
        for i in range(0, count, 2):
            val1 = 1 if values[i] else 0
            val2 = 0
            if i + 1 < count:
                val2 = 1 if values[i+1] else 0
            
            # High nibble is first device (i), Low nibble is second device (i+1)
            byte_val = (val1 << 4) | val2
            data_bytes.append(byte_val)
            
        self._send_command(command, subcommand, data_header + data_bytes)

    def random_read_word(self, word_devices: List[Tuple[str, int]]) -> List[int]:
        """
        隨機讀取字組 (Random Read Word)。
        一次請求中讀取多個不連續的字組軟元件。

        Args:
            word_devices (List[Tuple[str, int]]): 設備列表，例如 [('D', 100), ('D', 105), ('W', 200)]。
                                                 注意：目前實作限制每次請求最多 192 個 Word。

        Returns:
            List[int]: 對應設備的數值列表。
        """
        count = len(word_devices)
        if count == 0:
            return []
        
        # Command: 0403 (Random Read)
        # Subcommand: 0000 (Word Access)
        command = b'\x03\x04'
        subcommand = b'\x00\x00'

        # Request Data:
        # Word Access Points (1 byte)
        # Double Word Access Points (1 byte) (目前設為 0)
        # Word Access Data [DevCode(1) + HeadNo(3)] * N
        
        data = bytearray()
        data += struct.pack('B', count) # Word Access Points
        data += b'\x00' # Double Word Access Points
        
        for dev_type, dev_addr in word_devices:
            dev_code, _ = self._parse_device_code(dev_type)
            # 3 bytes Head Device (Little Endian)
            addr_bytes = struct.pack('<I', dev_addr)[:3]
            data += struct.pack('B', dev_code) + addr_bytes
            
        response = self._send_command(command, subcommand, data)
        
        # Response Data: Value (2 bytes) * Count
        if len(response) != count * 2:
            raise MCProtocolError(f"回應長度不符: 預期 {count*2}, 實際 {len(response)}")
            
        values = struct.unpack(f'<{count}H', response)
        return list(values)

    # 便捷方法 (Alias)
    def read_d(self, address: int, count: int = 1) -> List[int]:
        """讀取 D 暫存器 (Word)"""
        return self.batch_read_word('D', address, count)

    def write_d(self, address: int, values: Union[int, List[int]]):
        """寫入 D 暫存器 (Word)。values 可以是單個整數或整數列表。"""
        if isinstance(values, int):
            values = [values]
        self.batch_write_word('D', address, values)

    def read_m(self, address: int, count: int = 1) -> List[int]:
        """讀取 M 繼電器 (Bit)"""
        return self.batch_read_bit('M', address, count)

    def write_m(self, address: int, values: Union[int, List[int]]):
        """寫入 M 繼電器 (Bit)。values 可以是單個整數(0/1)或列表。"""
        if isinstance(values, int):
            values = [values]
        self.batch_write_bit('M', address, values)

if __name__ == "__main__":
    # 使用範例
    print("MC Protocol Library Loaded.")
    print("Example usage:")
    print("  plc = MCProtocol('192.168.1.10', 5000)")
    print("  if plc.check_connection():")
    print("      data = plc.read_d(100, 10)")
    print("      random_data = plc.random_read_word([('D', 100), ('D', 200)])")
    print("      print(random_data)")