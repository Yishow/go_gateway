"""
Fatek PLC Communication Driver
==============================

This module provides a complete interface for communicating with FATEK FBs series PLCs.
It supports both Serial (RS-232/485) and TCP/IP (Ethernet) transport layers.

The driver implements the FATEK ASCII Protocol, including checksum calculation (LRC),
frame assembly, and error handling.

Dependencies:
    - pyserial
    - socket (standard library)

Usage:
    >>> from fatek_driver import create_tcp_client
    >>> client = create_tcp_client('192.168.1.5')
    >>> client.connect()
    >>> # Read X0-X5
    >>> status = client.read_status('X', 0, 5)
    >>> print(status)
    >>> client.close()
"""

import socket
import time
import serial
import re
from abc import ABC, abstractmethod
from typing import List, Union, Tuple, Optional, Dict

# Constants
STX = b'\x02'
ETX = b'\x03'
DEFAULT_STATION = 1

class FatekException(Exception):
    """Base exception for all Fatek driver errors."""
    pass

class FatekCommunicationError(FatekException):
    """
    Raised when low-level communication fails.

    Examples:
        - Connection timeout.
        - Checksum (LRC) mismatch.
        - Malformed frames (missing STX/ETX).
    """
    pass

class FatekProtocolError(FatekException):
    """
    Raised when the PLC receives the command but returns a logical error code.

    Attributes:
        error_code (str): The error character returned by PLC ('2', '4', 'A', etc.).
        command (str): The command ID that caused the error (e.g., '44').
    """
    def __init__(self, error_code: str, command: str):
        self.error_code = error_code
        self.command = command
        msg = f"PLC Error {error_code} on Command {command}"
        if error_code == '2':
            msg += ": Illegal Value (Value out of range or non-hex)"
        elif error_code == '4':
            msg += ": Illegal Format/Command (Unsupported command or LRC error)"
        elif error_code == 'A':
            msg += ": Illegal Address (Address boundary exceeded)"
        super().__init__(msg)

class FatekTransport(ABC):
    """
    Abstract base class for transport layer implementations.
    Defines the interface for sending and receiving raw bytes.
    """
    
    @abstractmethod
    def connect(self):
        """Establishes the connection to the PLC."""
        pass
        
    @abstractmethod
    def close(self):
        """Closes the connection."""
        pass
        
    @abstractmethod
    def send_receive(self, data: bytes) -> bytes:
        """
        Sends data to the PLC and waits for a response.

        Args:
            data (bytes): The raw request frame including STX, content, LRC, and ETX.

        Returns:
            bytes: The raw response frame.

        Raises:
            FatekCommunicationError: If send/receive fails or times out.
        """
        pass

class FatekSerial(FatekTransport):
    """
    Serial transport implementation using pyserial.
    
    Args:
        port (str): COM port name (e.g., 'COM1', '/dev/ttyUSB0').
        baudrate (int): Baud rate. Defaults to 9600.
        data_bits (int): Number of data bits. Defaults to 7 (Standard Fatek ASCII).
        parity (str): Parity check ('N', 'E', 'O'). Defaults to 'E' (Even).
        stop_bits (int): Number of stop bits. Defaults to 1.
        timeout (float): Read timeout in seconds. Defaults to 1.0.
    """
    def __init__(self, port: str, baudrate=9600, data_bits=7, parity='E', stop_bits=1, timeout=1.0):
        self.port = port
        self.baudrate = baudrate
        self.data_bits = data_bits
        self.parity = parity
        self.stop_bits = stop_bits
        self.timeout = timeout
        self._serial = None

    def connect(self):
        """Opens the serial port."""
        parity_map = {'N': serial.PARITY_NONE, 'E': serial.PARITY_EVEN, 'O': serial.PARITY_ODD}
        self._serial = serial.Serial(
            port=self.port,
            baudrate=self.baudrate,
            bytesize=self.data_bits,
            parity=parity_map.get(self.parity, serial.PARITY_EVEN),
            stopbits=self.stop_bits,
            timeout=self.timeout
        )

    def close(self):
        """Closes the serial port."""
        if self._serial and self._serial.is_open:
            self._serial.close()

    def send_receive(self, data: bytes) -> bytes:
        """Sends data over serial and reads until ETX is found."""
        if not self._serial or not self._serial.is_open:
            raise FatekCommunicationError("Serial port not open")
        
        # Clear buffer to ensure no old data interferes
        self._serial.reset_input_buffer()
        
        self._serial.write(data)
        
        # Read byte by byte until ETX (0x03)
        response = b""
        while True:
            char = self._serial.read(1)
            if not char:
                raise FatekCommunicationError("Timeout waiting for response")
            response += char
            if char == ETX:
                break
        return response

class FatekTCP(FatekTransport):
    """
    TCP/IP transport implementation.
    Used for FATEK FBs PLCs with Ethernet modules or FATEK-Ethernet converters.
    
    Args:
        host (str): PLC IP address.
        port (int): TCP port. Defaults to 500.
        timeout (float): Socket timeout in seconds. Defaults to 2.0.
    """
    def __init__(self, host: str, port=500, timeout=2.0):
        self.host = host
        self.port = port
        self.timeout = timeout
        self._socket = None

    def connect(self):
        """Establishes a TCP socket connection."""
        self._socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._socket.settimeout(self.timeout)
        self._socket.connect((self.host, self.port))

    def close(self):
        """Closes the TCP socket."""
        if self._socket:
            self._socket.close()
            self._socket = None

    def send_receive(self, data: bytes) -> bytes:
        """Sends data over TCP and reads until ETX is found."""
        if not self._socket:
            raise FatekCommunicationError("Socket not connected")
        
        try:
            self._socket.sendall(data)
            
            # Read until ETX
            response = b""
            while True:
                chunk = self._socket.recv(1)
                if not chunk:
                    raise FatekCommunicationError("Connection closed by peer")
                response += chunk
                if chunk == ETX:
                    break
            return response
        except socket.timeout:
            raise FatekCommunicationError("Socket timeout")
        except socket.error as e:
            raise FatekCommunicationError(f"Socket error: {e}")

class FatekClient:
    """
    High-level FATEK PLC Client.
    
    This class handles the logic of FATEK ASCII Protocol commands,
    address formatting, and response parsing.
    
    Args:
        transport (FatekTransport): The transport instance (Serial or TCP).
        station (int): The PLC station number (0-255). Defaults to 1.
    """
    
    def __init__(self, transport: FatekTransport, station=DEFAULT_STATION):
        self.transport = transport
        self.station = station

    def _calculate_lrc(self, data: bytes) -> str:
        """
        Calculates the Longitudinal Redundancy Check (LRC) checksum.
        
        Args:
            data (bytes): The byte sequence to checksum (including STX).
            
        Returns:
            str: Two-character hexadecimal string of the checksum.
        """
        lrc = 0
        for b in data:
            lrc = (lrc + b) & 0xFF
        return f"{lrc:02X}"

    def _build_frame(self, command: str, body: str = "") -> bytes:
        """
        Constructs a complete FATEK ASCII frame.
        
        Structure: STX + Station(2) + Command(2) + Body + LRC(2) + ETX
        
        Args:
            command (str): The 2-character command code (e.g., '44').
            body (str): The data payload for the command.
            
        Returns:
            bytes: The encoded ASCII frame ready for transmission.
        """
        station_str = f"{self.station:02X}"
        content = f"{station_str}{command}{body}"
        
        # Calculate LRC. Note: According to documentation, LRC includes STX.
        stx_val = ord(STX)
        lrc_sum = stx_val
        for char in content:
            lrc_sum += ord(char)
        
        lrc_hex = f"{lrc_sum & 0xFF:02X}"
        
        frame = f"{STX.decode()}{content}{lrc_hex}{ETX.decode()}"
        return frame.encode('ascii')

    def _parse_response(self, response: bytes, expected_cmd: str) -> str:
        """
        Validates and parses the response frame.
        
        Checks STX, ETX, LRC, and the error code returned by the PLC.
        
        Args:
            response (bytes): The raw response frame.
            expected_cmd (str): The command code that was sent.
            
        Returns:
            str: The body of the response (data payload).
            
        Raises:
            FatekCommunicationError: If structure or checksum is invalid.
            FatekProtocolError: If PLC returns an error code (non-zero status).
        """
        if len(response) < 6:
             pass

        if response[0:1] != STX:
            raise FatekCommunicationError("Invalid STX")
        if response[-1:] != ETX:
            raise FatekCommunicationError("Invalid ETX")

        # Validate Checksum
        # Everything before the last 3 bytes (LRC(2) + ETX(1)) is content
        content_with_stx = response[:-3] 
        received_lrc = response[-3:-1].decode()
        
        lrc_sum = 0
        for b in content_with_stx:
            lrc_sum += b
        
        calculated_lrc = f"{lrc_sum & 0xFF:02X}"
        
        if received_lrc != calculated_lrc:
            raise FatekCommunicationError(f"Checksum Mismatch. Recv: {received_lrc}, Calc: {calculated_lrc}")

        # Parse Fields
        # Station(2) -> [1:3]
        # Command(2) -> [3:5]
        cmd = response[3:5].decode()
        
        if cmd != expected_cmd:
             raise FatekCommunicationError(f"Unexpected command in response: {cmd}")
             
        # Status Code / Error Code (Byte after Cmd, index 5)
        # For Read commands (e.g., 44, 46, 48), success is '0' followed by data.
        # For Write commands (e.g., 45, 47), success is just '0' (no body).
        # For Errors, it is an error code '2', '4', 'A' instead of '0'.
        
        status_code = response[5:6].decode()
        
        if status_code != '0':
            raise FatekProtocolError(status_code, cmd)
            
        # Return Body (Skip STX(1)+Station(2)+Cmd(2)+StatusCode(1) = 6 bytes)
        # We start from index 6. Checksum is at -3.
        body = response[6:-3].decode()
        return body

    def _normalize_address(self, symbol_type: str, number: int) -> str:
        """
        Formats the component address to the fixed length required by the protocol.
        
        Args:
            symbol_type (str): The component type (e.g., 'X', 'R', 'D').
            number (int): The address number.
            
        Returns:
            str: Formatted address string (e.g., 'X00010', 'D00100').
            
        Raises:
            ValueError: If the symbol type is unknown.
        """
        symbol = symbol_type.upper()
        
        # 5 chars: X, Y, M, S, T, C
        if symbol in ['X', 'Y', 'M', 'S', 'T', 'C']:
            return f"{symbol}{number:04d}"
            
        # 6 chars: R, D, RT, RC, F
        if symbol in ['R', 'D', 'RT', 'RC', 'F']:
            return f"{symbol}{number:05d}"
            
        # 6 chars: WX, WY, WM, WS, WT, WC (16-bit access to discrete)
        if symbol in ['WX', 'WY', 'WM', 'WS', 'WT', 'WC']:
            return f"{symbol}{number:04d}"

        # 7 chars: 32-bit registers (DWX, DR, DD...)
        if symbol in ['DR', 'DD', 'DF'] or symbol.startswith('DW') or symbol.startswith('D'):
             if len(symbol) >= 2: # DWX, DD, etc.
                 if symbol in ['DWX', 'DWY', 'DWM', 'DWS', 'DWT', 'DWC']:
                     return f"{symbol}{number:04d}"
                 else:
                     return f"{symbol}{number:05d}"

        raise ValueError(f"Unknown component type: {symbol}")

    def connect(self):
        """Connects the underlying transport."""
        self.transport.connect()
        
    def close(self):
        """Closes the underlying transport."""
        self.transport.close()

    # --- Commands ---

    def read_status(self, symbol: str, start_addr: int, count: int) -> List[bool]:
        """
        Reads continuous discrete status (Cmd 44).
        
        Used for X, Y, M, S, T(status), C(status).

        Args:
            symbol (str): Component type (e.g., 'X', 'Y').
            start_addr (int): Starting address.
            count (int): Number of points to read (Max 255).

        Returns:
            List[bool]: List of statuses (True for ON, False for OFF).
        """
        if count > 255:
            raise ValueError("Max count is 255")
            
        addr_str = self._normalize_address(symbol, start_addr)
        count_hex = f"{count:02X}"
        
        body = f"{count_hex}{addr_str}"
        req = self._build_frame("44", body)
        resp = self.transport.send_receive(req)
        data_str = self._parse_response(resp, "44")
        
        # Data is string of '0' and '1'
        return [True if c == '1' else False for c in data_str]

    def write_status(self, symbol: str, start_addr: int, data: List[bool]):
        """
        Writes continuous discrete status (Cmd 45).

        Args:
            symbol (str): Component type.
            start_addr (int): Starting address.
            data (List[bool]): List of boolean values to write.
        """
        count = len(data)
        if count > 255:
            raise ValueError("Max count is 255")
            
        addr_str = self._normalize_address(symbol, start_addr)
        count_hex = f"{count:02X}"
        data_str = "".join(['1' if b else '0' for b in data])
        
        body = f"{count_hex}{addr_str}{data_str}"
        req = self._build_frame("45", body)
        resp = self.transport.send_receive(req)
        self._parse_response(resp, "45")

    def read_registers(self, symbol: str, start_addr: int, count: int) -> List[int]:
        """
        Reads continuous 16-bit registers (Cmd 46).
        
        Used for R, D, RT, RC.

        Args:
            symbol (str): Component type (e.g., 'D', 'R').
            start_addr (int): Starting address.
            count (int): Number of registers to read (Max 64).

        Returns:
            List[int]: List of register values (unsigned 16-bit integers).
        """
        if count > 64:
             raise ValueError("Max count is 64")
        
        addr_str = self._normalize_address(symbol, start_addr)
        count_hex = f"{count:02X}"
        
        body = f"{count_hex}{addr_str}"
        req = self._build_frame("46", body)
        resp = self.transport.send_receive(req)
        data_str = self._parse_response(resp, "46")
        
        # Data is 4 hex chars per word
        results = []
        for i in range(0, len(data_str), 4):
            val_hex = data_str[i:i+4]
            results.append(int(val_hex, 16))
        return results

    def write_registers(self, symbol: str, start_addr: int, data: List[int]):
        """
        Writes continuous 16-bit registers (Cmd 47).

        Args:
            symbol (str): Component type.
            start_addr (int): Starting address.
            data (List[int]): List of integer values to write.
        """
        count = len(data)
        if count > 64:
             raise ValueError("Max count is 64")
             
        addr_str = self._normalize_address(symbol, start_addr)
        count_hex = f"{count:02X}"
        
        data_body = ""
        for val in data:
            # Ensure 16-bit limit
            val = val & 0xFFFF
            data_body += f"{val:04X}"
            
        body = f"{count_hex}{addr_str}{data_body}"
        req = self._build_frame("47", body)
        resp = self.transport.send_receive(req)
        self._parse_response(resp, "47")

    def read_random(self, items: List[Tuple[str, int]]) -> Dict[str, Union[bool, int]]:
        """
        Mixed/Random Read (Cmd 48).
        Allows reading multiple arbitrary addresses (Discrete or Register) in one go.

        Args:
            items (List[Tuple[str, int]]): List of (symbol, address) tuples.
                                           e.g. [('X', 0), ('R', 100), ('M', 20)]

        Returns:
            Dict[str, Union[bool, int]]: Dictionary mapping 'SymbolAddress' to value.
                                         e.g. {'X0': True, 'R100': 1234}
        """
        if len(items) > 64: # Protocol limit approx
             raise ValueError("Too many items")
             
        count_hex = f"{len(items):02X}"
        body = count_hex
        
        for sym, addr in items:
            body += self._normalize_address(sym, addr)
            
        req = self._build_frame("48", body)
        resp = self.transport.send_receive(req)
        raw_data = self._parse_response(resp, "48")
        
        # Parse mixed response
        results = {}
        ptr = 0
        for sym, addr in items:
            key = f"{sym}{addr}"
            # Determine if it's discrete or register based on symbol
            if sym.upper() in ['X', 'Y', 'M', 'S', 'T', 'C']:
                # Discrete: 1 char
                val_str = raw_data[ptr:ptr+1]
                results[key] = (val_str == '1')
                ptr += 1
            else:
                # Register: 4 chars
                val_str = raw_data[ptr:ptr+4]
                results[key] = int(val_str, 16)
                ptr += 4
        
        return results

    def loopback_test(self, data: str = "AABB") -> bool:
        """
        Performs a Loopback Test (Cmd 4E).
        Sends data to the PLC and checks if it echoes back exactly.

        Args:
            data (str): The string data to send (Hex string usually).

        Returns:
            bool: True if echo matches, False otherwise.
        """
        req = self._build_frame("4E", data)
        resp = self.transport.send_receive(req)
        
        try:
            # Cmd 4E response structure is special: it just echoes.
            # No status code '0' field.
            # Response: STX + Station + 4E + Data + LRC + ETX
            
            if len(resp) < 6: return False
            if resp[0:1] != STX or resp[-1:] != ETX: return False
            
            # Extract content
            received_cmd = resp[3:5].decode()
            if received_cmd != "4E": return False
            
            received_data = resp[5:-3].decode()
            return received_data == data
            
        except Exception:
            return False

    def run(self):
        """Starts the PLC (RUN mode) using Cmd 41."""
        self._control_run_stop(run=True)

    def stop(self):
        """Stops the PLC (STOP mode) using Cmd 41."""
        self._control_run_stop(run=False)

    def _control_run_stop(self, run: bool):
        """Helper for RUN/STOP commands."""
        ctrl_code = '1' if run else '0'
        req = self._build_frame("41", ctrl_code)
        resp = self.transport.send_receive(req)
        self._parse_response(resp, "41")

    def single_action(self, symbol: str, addr: int, action: str):
        """
        Performs a single discrete control action (Cmd 42).

        Args:
            symbol (str): Component type (e.g., 'Y', 'M').
            addr (int): Component address.
            action (str): Operation to perform.
                          Options: 'DISABLE', 'ENABLE', 'SET', 'RESET'.

        Raises:
            ValueError: If action is invalid.
        """
        action_map = {
            'DISABLE': '1',
            'ENABLE': '2',
            'SET': '3',
            'RESET': '4'
        }
        
        code = action_map.get(action.upper())
        if not code:
            raise ValueError("Invalid action. Use: DISABLE, ENABLE, SET, RESET")
            
        addr_str = self._normalize_address(symbol, addr)
        body = f"{code}{addr_str}"
        
        req = self._build_frame("42", body)
        resp = self.transport.send_receive(req)
        self._parse_response(resp, "42")


# Factory Functions

def create_serial_client(port: str, station=1) -> FatekClient:
    """
    Creates a FatekClient configured for Serial communication.

    Args:
        port (str): Serial port name.
        station (int): PLC station ID.

    Returns:
        FatekClient: Configured client instance.
    """
    transport = FatekSerial(port)
    return FatekClient(transport, station)

def create_tcp_client(host: str, port=500, station=1) -> FatekClient:
    """
    Creates a FatekClient configured for TCP communication.

    Args:
        host (str): IP address of the PLC.
        port (int): TCP port (usually 500).
        station (int): PLC station ID.

    Returns:
        FatekClient: Configured client instance.
    """
    transport = FatekTCP(host, port)
    return FatekClient(transport, station)

if __name__ == "__main__":
    print("Fatek Library Loaded.")
    print("Use help(fatek_driver) to see documentation.")