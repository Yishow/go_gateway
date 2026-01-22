
import { ProtocolType } from '../types/datalink';

export interface ParsedAddress {
  protocol: ProtocolType;
  area: string;
  startNumber: number;
  raw: string;
}

export interface ValidationResult {
  valid: boolean;
  error?: string;
}

export class AddressParser {
  /**
   * Parse an address string based on protocol
   */
  parse(address: string, protocol: ProtocolType = 'modbus_tcp'): ParsedAddress {
    const cleanAddress = address.trim().toUpperCase();
    
    if (protocol.startsWith('modbus')) {
      return this.parseModbus(cleanAddress, protocol);
    } else if (protocol === 'fatek_fbs') {
      return this.parseFatek(cleanAddress, protocol);
    } else if (protocol === 'mc_3e') {
      return this.parseMC3E(cleanAddress, protocol);
    }
    
    throw new Error(`Unsupported protocol: ${protocol}`);
  }

  private parseModbus(address: string, protocol: ProtocolType): ParsedAddress {
    // Modbus: 0xxxx, 1xxxx, 3xxxx, 4xxxx
    // 5-6 digits. Leading digit indicates area.
    const match = address.match(/^([0134])(\d{4,5})$/);
    if (!match) {
      throw new Error(`Invalid Modbus address: ${address}`);
    }
    
    const prefix = match[1];
    const offsetStr = match[2];
    const startNumber = parseInt(offsetStr, 10);
    
    let area = 'UNKNOWN';
    switch (prefix) {
      case '0': area = 'CS'; break; // Coil Status
      case '1': area = 'IS'; break; // Input Status
      case '3': area = 'IR'; break; // Input Register
      case '4': area = 'HR'; break; // Holding Register
    }
    
    return {
      protocol,
      area,
      startNumber, // Keep as is (e.g. 40001 -> 1) logic handled by driver usually, but for UI we parse what user sees
      raw: address
    };
  }

  private parseFatek(address: string, protocol: ProtocolType): ParsedAddress {
    // Fatek: X, Y, M, S, T, C, D, R
    const match = address.match(/^([XYMSTCDR])(\d+)$/);
    if (!match) {
      throw new Error(`Invalid Fatek address: ${address}`);
    }
    
    return {
      protocol,
      area: match[1],
      startNumber: parseInt(match[2], 10),
      raw: address
    };
  }

  private parseMC3E(address: string, protocol: ProtocolType): ParsedAddress {
    // MC3E: D, W, M, X, Y, B, etc.
    const match = address.match(/^([DWMXYB])(\d+)$/); // Simplified regex
    if (!match) {
        // Try hex for X/Y/B if needed, but for MVP sticking to digit
        // Mitsubishi X/Y are usually hex. 
        // Let's support Hex for X/Y/B if it matches
        const hexMatch = address.match(/^([XYB])([0-9A-F]+)$/);
        if (hexMatch) {
             return {
                protocol,
                area: hexMatch[1],
                startNumber: parseInt(hexMatch[2], 16), // Treat as hex
                raw: address
            };
        }
        throw new Error(`Invalid MC3E address: ${address}`);
    }
    
    return {
      protocol,
      area: match[1],
      startNumber: parseInt(match[2], 10),
      raw: address
    };
  }

  /**
   * Expand a start address to a list of addresses
   */
  expand(startAddress: string, count: number, protocol: ProtocolType = 'modbus_tcp'): string[] {
    try {
      const parsed = this.parse(startAddress, protocol);
      const result: string[] = [];
      
      for (let i = 0; i < count; i++) {
        const currentNumber = parsed.startNumber + i;
        result.push(this.format(parsed.area, currentNumber, protocol, parsed.raw));
      }
      return result;
    } catch (e) {
      return [];
    }
  }

  private format(area: string, number: number, protocol: ProtocolType, _originalRaw: string): string {
    if (protocol.startsWith('modbus')) {
      // Reconstruct Modbus address
      // Need to know which prefix to use based on Area
      let prefix = '4';
      if (area === 'CS') prefix = '0';
      if (area === 'IS') prefix = '1';
      if (area === 'IR') prefix = '3';
      if (area === 'HR') prefix = '4';
      
      // Pad to match original length if possible, or standard 5 digits
      // originalRaw: 40001 (5 digits). 4 + 0001. 
      // If we have 1, we need 0001.
      // Simple padding:
      return `${prefix}${number.toString().padStart(4, '0')}`;
    } else if (protocol === 'fatek_fbs' || protocol === 'mc_3e') {
        // Just Area + Number
        // For MC3E X/Y/B, format as Hex? 
        // If the parsed raw input was Hex, we should output Hex.
        // We can check if area is X/Y/B
        if (protocol === 'mc_3e' && ['X', 'Y', 'B'].includes(area)) {
             return `${area}${number.toString(16).toUpperCase()}`;
        }
        return `${area}${number}`;
    }
    return `${area}${number}`;
  }

  validate(address: string, protocol: ProtocolType = 'modbus_tcp'): ValidationResult {
    try {
      this.parse(address, protocol);
      return { valid: true };
    } catch (e) {
        if (e instanceof Error) {
            return { valid: false, error: e.message };
        }
        return { valid: false, error: 'Unknown error' };
    }
  }
}

export const addressParser = new AddressParser();
