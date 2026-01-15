# ============================================================================
# Fatek PLC 進階工具集
# 包含：批量操作、循環監控、資料記錄、狀態管理
# ============================================================================

import time
import threading
from typing import List, Dict, Callable, Any, Optional
from collections import deque
from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum
import json

from fatek_plc_library import (
    FatekPLC, ComponentType, FatekException, SerialConfig
)


# ============================================================================
# 資料紀錄類別
# ============================================================================

@dataclass
class DataPoint:
    """資料點"""
    timestamp: datetime
    component: str
    address: int
    value: Any
    
    def to_dict(self) -> Dict:
        """轉換為字典"""
        return {
            "timestamp": self.timestamp.isoformat(),
            "component": self.component,
            "address": self.address,
            "value": self.value
        }


class DataLogger:
    """資料紀錄器"""
    
    def __init__(self, max_records: int = 10000):
        self.max_records = max_records
        self.records: deque = deque(maxlen=max_records)
        self.lock = threading.Lock()
    
    def add_record(
        self,
        component: str,
        address: int,
        value: Any
    ) -> None:
        """添加記錄"""
        with self.lock:
            point = DataPoint(
                timestamp=datetime.now(),
                component=component,
                address=address,
                value=value
            )
            self.records.append(point)
    
    def get_records(self, limit: int = 100) -> List[Dict]:
        """取得最近的記錄"""
        with self.lock:
            records = list(self.records)[-limit:]
            return [r.to_dict() for r in records]
    
    def export_csv(self, filename: str) -> None:
        """匯出為 CSV"""
        import csv
        
        with self.lock:
            with open(filename, 'w', newline='', encoding='utf-8') as f:
                writer = csv.DictWriter(
                    f,
                    fieldnames=['timestamp', 'component', 'address', 'value']
                )
                writer.writeheader()
                for record in self.records:
                    writer.writerow(record.to_dict())
    
    def export_json(self, filename: str) -> None:
        """匯出為 JSON"""
        with self.lock:
            data = [r.to_dict() for r in self.records]
            with open(filename, 'w', encoding='utf-8') as f:
                json.dump(data, f, ensure_ascii=False, indent=2)
    
    def clear(self) -> None:
        """清空記錄"""
        with self.lock:
            self.records.clear()


# ============================================================================
# 狀態監控類別
# ============================================================================

@dataclass
class MonitoringTask:
    """監控任務"""
    component_type: ComponentType
    start_address: int
    quantity: int
    interval: float  # 秒
    callback: Optional[Callable] = None
    enabled: bool = True
    last_values: List[int] = field(default_factory=list)
    read_count: int = 0
    error_count: int = 0


class PLCMonitor:
    """PLC 狀態監控器"""
    
    def __init__(self, plc: FatekPLC, logger=None):
        self.plc = plc
        self.logger = logger
        self.tasks: Dict[str, MonitoringTask] = {}
        self.data_logger = DataLogger()
        self.running = False
        self.monitor_thread: Optional[threading.Thread] = None
        self.lock = threading.Lock()
    
    def add_task(
        self,
        task_id: str,
        component_type: ComponentType,
        start_address: int,
        quantity: int,
        interval: float = 1.0,
        callback: Optional[Callable] = None
    ) -> None:
        """添加監控任務"""
        with self.lock:
            self.tasks[task_id] = MonitoringTask(
                component_type=component_type,
                start_address=start_address,
                quantity=quantity,
                interval=interval,
                callback=callback
            )
            if self.logger:
                self.logger.info(
                    f"添加監控任務: {task_id} "
                    f"({component_type.name}{start_address}~"
                    f"{start_address + quantity - 1})"
                )
    
    def remove_task(self, task_id: str) -> None:
        """移除監控任務"""
        with self.lock:
            if task_id in self.tasks:
                del self.tasks[task_id]
                if self.logger:
                    self.logger.info(f"移除監控任務: {task_id}")
    
    def start(self) -> None:
        """啟動監控"""
        if self.running:
            return
        
        self.running = True
        self.monitor_thread = threading.Thread(target=self._monitor_loop)
        self.monitor_thread.daemon = True
        self.monitor_thread.start()
        
        if self.logger:
            self.logger.info("監控已啟動")
    
    def stop(self) -> None:
        """停止監控"""
        self.running = False
        if self.monitor_thread:
            self.monitor_thread.join(timeout=5)
        
        if self.logger:
            self.logger.info("監控已停止")
    
    def _monitor_loop(self) -> None:
        """監控迴圈"""
        while self.running:
            with self.lock:
                tasks = list(self.tasks.values())
            
            for task in tasks:
                if not task.enabled:
                    continue
                
                try:
                    # 讀取資料
                    if task.component_type in [
                        ComponentType.R, ComponentType.D, ComponentType.W
                    ]:
                        values = self.plc.read_register(
                            task.component_type,
                            task.start_address,
                            task.quantity
                        )
                    else:
                        values = self.plc.read_discrete(
                            task.component_type,
                            task.start_address,
                            task.quantity
                        )
                    
                    # 記錄資料
                    for i, value in enumerate(values):
                        address = task.start_address + i
                        self.data_logger.add_record(
                            task.component_type.name,
                            address,
                            value
                        )
                    
                    # 檢查是否有變化
                    if task.last_values != values:
                        task.last_values = values
                        
                        # 調用回調
                        if task.callback:
                            task.callback(
                                task.component_type,
                                task.start_address,
                                values
                            )
                    
                    task.read_count += 1
                
                except FatekException as e:
                    task.error_count += 1
                    if self.logger:
                        self.logger.error(f"監控錯誤: {e}")
                
                time.sleep(task.interval)
    
    def get_task_status(self, task_id: str) -> Dict[str, Any]:
        """取得任務狀態"""
        with self.lock:
            task = self.tasks.get(task_id)
            if not task:
                return {}
            
            return {
                "task_id": task_id,
                "component": task.component_type.name,
                "start_address": task.start_address,
                "quantity": task.quantity,
                "interval": task.interval,
                "enabled": task.enabled,
                "read_count": task.read_count,
                "error_count": task.error_count,
                "last_values": task.last_values,
                "success_rate": (
                    100 * task.read_count / (task.read_count + task.error_count)
                    if (task.read_count + task.error_count) > 0 else 0
                )
            }
    
    def get_all_task_status(self) -> List[Dict[str, Any]]:
        """取得所有任務狀態"""
        with self.lock:
            task_ids = list(self.tasks.keys())
        
        return [self.get_task_status(task_id) for task_id in task_ids]


# ============================================================================
# 批量操作管理器
# ============================================================================

@dataclass
class BatchOperation:
    """批量操作"""
    operation_id: str
    operation_type: str  # 'read' 或 'write'
    component_type: ComponentType
    start_address: int
    quantity: int
    values: Optional[List[int]] = None
    status: str = 'pending'  # pending, running, completed, failed
    error_message: str = ''
    result: List[int] = field(default_factory=list)
    timestamp: datetime = field(default_factory=datetime.now)


class BatchOperationManager:
    """批量操作管理器"""
    
    def __init__(self, plc: FatekPLC, logger=None):
        self.plc = plc
        self.logger = logger
        self.operations: Dict[str, BatchOperation] = {}
        self.queue: deque = deque()
        self.lock = threading.Lock()
        self.executor_thread: Optional[threading.Thread] = None
        self.running = False
    
    def add_read_operation(
        self,
        operation_id: str,
        component_type: ComponentType,
        start_address: int,
        quantity: int
    ) -> None:
        """添加讀取操作"""
        with self.lock:
            op = BatchOperation(
                operation_id=operation_id,
                operation_type='read',
                component_type=component_type,
                start_address=start_address,
                quantity=quantity
            )
            self.operations[operation_id] = op
            self.queue.append(op)
    
    def add_write_operation(
        self,
        operation_id: str,
        component_type: ComponentType,
        start_address: int,
        values: List[int]
    ) -> None:
        """添加寫入操作"""
        with self.lock:
            op = BatchOperation(
                operation_id=operation_id,
                operation_type='write',
                component_type=component_type,
                start_address=start_address,
                quantity=len(values),
                values=values
            )
            self.operations[operation_id] = op
            self.queue.append(op)
    
    def start_executor(self) -> None:
        """啟動執行器"""
        if self.running:
            return
        
        self.running = True
        self.executor_thread = threading.Thread(target=self._execute_loop)
        self.executor_thread.daemon = True
        self.executor_thread.start()
    
    def stop_executor(self) -> None:
        """停止執行器"""
        self.running = False
        if self.executor_thread:
            self.executor_thread.join(timeout=5)
    
    def _execute_loop(self) -> None:
        """執行迴圈"""
        while self.running:
            with self.lock:
                if not self.queue:
                    op = None
                else:
                    op = self.queue.popleft()
            
            if not op:
                time.sleep(0.1)
                continue
            
            try:
                op.status = 'running'
                
                if op.operation_type == 'read':
                    if op.component_type in [
                        ComponentType.R, ComponentType.D, ComponentType.W
                    ]:
                        op.result = self.plc.read_register(
                            op.component_type,
                            op.start_address,
                            op.quantity
                        )
                    else:
                        op.result = self.plc.read_discrete(
                            op.component_type,
                            op.start_address,
                            op.quantity
                        )
                
                elif op.operation_type == 'write':
                    if op.component_type in [
                        ComponentType.R, ComponentType.D, ComponentType.W
                    ]:
                        self.plc.write_register(
                            op.component_type,
                            op.start_address,
                            op.values
                        )
                    else:
                        self.plc.write_discrete(
                            op.component_type,
                            op.start_address,
                            op.values
                        )
                
                op.status = 'completed'
            
            except FatekException as e:
                op.status = 'failed'
                op.error_message = str(e)
                if self.logger:
                    self.logger.error(f"操作 {op.operation_id} 失敗: {e}")
    
    def get_operation_status(self, operation_id: str) -> Dict[str, Any]:
        """取得操作狀態"""
        with self.lock:
            op = self.operations.get(operation_id)
            if not op:
                return {}
            
            return {
                "operation_id": op.operation_id,
                "operation_type": op.operation_type,
                "component": op.component_type.name,
                "start_address": op.start_address,
                "quantity": op.quantity,
                "status": op.status,
                "error_message": op.error_message,
                "result": op.result,
                "timestamp": op.timestamp.isoformat()
            }
    
    def wait_for_operation(
        self,
        operation_id: str,
        timeout: float = 10.0
    ) -> Dict[str, Any]:
        """等待操作完成"""
        start_time = time.time()
        
        while (time.time() - start_time) < timeout:
            status = self.get_operation_status(operation_id)
            
            if status.get('status') in ['completed', 'failed']:
                return status
            
            time.sleep(0.1)
        
        raise TimeoutError(f"操作 {operation_id} 逾時")


# ============================================================================
# 連線管理器
# ============================================================================

class ConnectionManager:
    """連線管理器 - 自動重連機制"""
    
    def __init__(
        self,
        config: SerialConfig,
        station_id: int = 1,
        max_retries: int = 3,
        retry_delay: float = 1.0,
        logger=None
    ):
        self.config = config
        self.station_id = station_id
        self.max_retries = max_retries
        self.retry_delay = retry_delay
        self.logger = logger
        self.plc: Optional[FatekPLC] = None
        self.connected = False
        self.connection_attempts = 0
        self.last_connection_time: Optional[datetime] = None
    
    def connect(self) -> bool:
        """連線（帶重試機制）"""
        for attempt in range(self.max_retries):
            try:
                self.plc = FatekPLC(
                    config=self.config,
                    station_id=self.station_id,
                    logger=self.logger
                )
                self.plc.connect()
                self.connected = True
                self.last_connection_time = datetime.now()
                self.connection_attempts = 0
                
                if self.logger:
                    self.logger.info("成功連線到 PLC")
                
                return True
            
            except Exception as e:
                self.connection_attempts += 1
                if self.logger:
                    self.logger.warning(
                        f"連線失敗 (嘗試 {attempt + 1}/{self.max_retries}): {e}"
                    )
                
                if attempt < self.max_retries - 1:
                    time.sleep(self.retry_delay)
        
        self.connected = False
        if self.logger:
            self.logger.error("無法連線到 PLC")
        
        return False
    
    def disconnect(self) -> bool:
        """斷線"""
        if self.plc and self.connected:
            self.plc.disconnect()
            self.connected = False
            
            if self.logger:
                self.logger.info("已斷線")
        
        return True
    
    def is_connected(self) -> bool:
        """檢查連線狀態"""
        return self.connected and self.plc and self.plc.is_connected()
    
    def ensure_connected(self) -> bool:
        """確保連線（如果斷線則重連）"""
        if not self.is_connected():
            return self.connect()
        return True


# ============================================================================
# 批量讀寫最佳化類別
# ============================================================================

class OptimizedBatchReader:
    """最佳化批量讀取器 - 合併相鄰地址"""
    
    def __init__(self, plc: FatekPLC, logger=None):
        self.plc = plc
        self.logger = logger
    
    def read_scattered(
        self,
        requests: List[tuple]  # [(component_type, start, quantity), ...]
    ) -> Dict[str, List[int]]:
        """
        讀取分散地址，自動合併相鄰地址
        
        args:
            requests: [(ComponentType.R, 0, 10), (ComponentType.R, 10, 5), ...]
        
        returns:
            {key: [values]}  例: {'R0': [1,2,3...], 'R10': [5,6...]}
        """
        results = {}
        
        # 按元件類型分組
        grouped = {}
        for component_type, start, quantity in requests:
            key = component_type.name
            if key not in grouped:
                grouped[key] = []
            grouped[key].append((start, quantity))
        
        # 對每個元件類型進行合併讀取
        for component_name, ranges in grouped.items():
            component_type = ComponentType[component_name]
            
            # 按開始位址排序
            ranges.sort(key=lambda x: x[0])
            
            # 合併相鄰範圍
            merged = self._merge_ranges(ranges)
            
            # 執行讀取
            for start, end in merged:
                quantity = end - start
                try:
                    values = self.plc.read_register(
                        component_type,
                        start,
                        quantity
                    ) if component_name in ['R', 'D', 'W'] else \
                        self.plc.read_discrete(
                            component_type,
                            start,
                            quantity
                        )
                    
                    # 儲存結果
                    for i, value in enumerate(values):
                        key = f"{component_name}{start + i}"
                        results[key] = value
                
                except Exception as e:
                    if self.logger:
                        self.logger.error(f"讀取失敗: {e}")
        
        return results
    
    def _merge_ranges(self, ranges: List[tuple]) -> List[tuple]:
        """合併相鄰範圍"""
        if not ranges:
            return []
        
        merged = [(ranges[0][0], ranges[0][0] + ranges[0][1])]
        
        for start, quantity in ranges[1:]:
            end = start + quantity
            last_start, last_end = merged[-1]
            
            # 如果相鄰或重疊，則合併
            if start <= last_end:
                merged[-1] = (last_start, max(last_end, end))
            else:
                merged.append((start, end))
        
        return merged


# ============================================================================
# 使用範例
# ============================================================================

def example_monitoring():
    """監控範例"""
    config = SerialConfig(port="COM1", baudrate=9600)
    plc = FatekPLC(config=config, station_id=1)
    
    plc.connect()
    
    # 建立監控器
    monitor = PLCMonitor(plc)
    
    # 添加監控任務
    monitor.add_task(
        'task_y',
        ComponentType.Y,
        0,
        16,
        interval=0.5
    )
    
    monitor.add_task(
        'task_r',
        ComponentType.R,
        0,
        10,
        interval=1.0
    )
    
    # 啟動監控
    monitor.start()
    
    # 執行 10 秒
    try:
        for _ in range(10):
            time.sleep(1)
            status = monitor.get_all_task_status()
            print(f"狀態: {status}")
    finally:
        monitor.stop()
        plc.disconnect()


def example_batch_operations():
    """批量操作範例"""
    config = SerialConfig(port="COM1", baudrate=9600)
    plc = FatekPLC(config=config, station_id=1)
    
    plc.connect()
    
    # 建立批量操作管理器
    batch_manager = BatchOperationManager(plc)
    batch_manager.start_executor()
    
    # 添加操作
    batch_manager.add_read_operation('read1', ComponentType.R, 0, 10)
    batch_manager.add_write_operation('write1', ComponentType.R, 10, [100, 200, 300])
    batch_manager.add_read_operation('read2', ComponentType.Y, 0, 16)
    
    # 等待完成
    try:
        result1 = batch_manager.wait_for_operation('read1', timeout=5)
        print(f"讀取結果: {result1}")
        
        result2 = batch_manager.wait_for_operation('write1', timeout=5)
        print(f"寫入結果: {result2}")
    finally:
        batch_manager.stop_executor()
        plc.disconnect()


if __name__ == "__main__":
    print("Fatek PLC 進階工具集")
    print("支援：監控、批量操作、資料記錄、連線管理")
