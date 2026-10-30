
import type { ProtocolType } from '../types/datalink';

const MC3E_MAX_OFFSET = 0xffffff;

export interface ParsedAddress {
  protocol: ProtocolType;
  area: string;
  startNumber: number;
  raw: string;
  type: 'bit' | 'word';
}

export interface ValidationResult {
  valid: boolean;
  error?: string;
}

export function getDefaultPlannerStartAddress(
  protocol: ProtocolType = 'modbus_tcp',
): string {
  switch (protocol) {
    case 'fatek_fbs':
    case 'mc_3e':
      return 'D0';
    case 'modbus_tcp':
    case 'modbus_rtu':
    case 'modbus_udp':
    default:
      return '40001';
  }
}

export class AddressParser {
  /**
   * Parse an address string based on protocol
   */
  parse(address: string, protocol: ProtocolType = 'modbus_tcp'): ParsedAddress {
    const cleanAddress = address.trim();

    if (protocol.startsWith('modbus')) {
      return this.parseModbus(cleanAddress.toUpperCase(), protocol);
    } else if (protocol === 'fatek_fbs') {
      return this.parseFatek(cleanAddress.toUpperCase(), protocol);
    } else if (protocol === 'mc_3e') {
      return this.parseMC3E(cleanAddress.toUpperCase(), protocol);
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
    
    const isBit = prefix === '0' || prefix === '1';
    return {
      protocol,
      area,
      startNumber,
      raw: address,
      type: isBit ? 'bit' : 'word',
    };
  }

  private parseFatek(address: string, protocol: ProtocolType): ParsedAddress {
    // Fatek: X, Y, M, S, T, C, D, R
    const match = address.match(/^([XYMSTCDR])(\d+)$/);
    if (!match) {
      throw new Error(`Invalid Fatek address: ${address}`);
    }
    
    const area = match[1];
    const isBit = ['X', 'Y', 'M', 'S', 'T', 'C'].includes(area);
    return {
      protocol,
      area,
      startNumber: parseInt(match[2], 10),
      raw: address,
      type: isBit ? 'bit' : 'word',
    };
  }

  private parseMC3E(address: string, protocol: ProtocolType): ParsedAddress {
    // MC X/Y contacts use strict octal suffixes.
    const octalMatch = address.match(/^([XY])([0-7]+)$/);
    if (octalMatch) {
      return {
        protocol,
        area: octalMatch[1],
        startNumber: this.parseBoundedMC3EOffset(octalMatch[2], 8, address),
        raw: address,
        type: 'bit',
      };
    }

    // MC W/B areas use hexadecimal suffixes.
    const hexMatch = address.match(/^([WB])([0-9A-F]+)$/);
    if (hexMatch) {
      return {
        protocol,
        area: hexMatch[1],
        startNumber: this.parseBoundedMC3EOffset(hexMatch[2], 16, address),
        raw: address,
        type: hexMatch[1] === 'W' ? 'word' : 'bit',
      };
    }

    // MC D/M areas use decimal suffixes.
    const decMatch = address.match(/^([DM])(\d+)$/);
    if (decMatch) {
      const area = decMatch[1];
      return {
        protocol,
        area,
        startNumber: this.parseBoundedMC3EOffset(decMatch[2], 10, address),
        raw: address,
        type: area === 'M' ? 'bit' : 'word',
      };
    }

    throw new Error(`Invalid MC3E address: ${address}`);
  }

  private parseBoundedMC3EOffset(raw: string, radix: number, address: string): number {
    const value = parseInt(raw, radix);
    if (!Number.isSafeInteger(value) || value < 0 || value > MC3E_MAX_OFFSET) {
      throw new Error(`Invalid MC3E address: ${address} (offset must be 0-${MC3E_MAX_OFFSET})`);
    }
    return value;
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
        if (
          protocol === 'mc_3e' &&
          (!Number.isSafeInteger(currentNumber) || currentNumber > MC3E_MAX_OFFSET)
        ) {
          return [];
        }
        result.push(this.format(parsed.area, currentNumber, protocol, parsed.raw));
      }
      return result;
    } catch {
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
      if (protocol === 'mc_3e') {
        const radix = ['X', 'Y'].includes(area) ? 8 : ['W', 'B'].includes(area) ? 16 : 10;
        return `${area}${number.toString(radix).toUpperCase()}`;
      }
      return `${area}${number}`;
    }
    return `${area}${number}`;
  }

  /**
   * 將位址依協議偏移 delta 格，回傳新位址字串。
   * Modbus 位址最小維持 1；FATEK / MC3E 則允許從 0 起算。
   */
  offset(address: string, delta: number, protocol: ProtocolType = 'modbus_tcp'): string {
    const parsed = this.parse(address, protocol);
    const minNumber =
      protocol === 'fatek_fbs' || protocol === 'mc_3e'
        ? 0
        : 1;
    const nextNumber = parsed.startNumber + delta;
    if (
      protocol === 'mc_3e' &&
      (!Number.isSafeInteger(delta) ||
        !Number.isSafeInteger(nextNumber) ||
        nextNumber > MC3E_MAX_OFFSET)
    ) {
      throw new Error(`MC3E address offset out of range: ${address} + ${delta}`);
    }
    const newNumber = Math.max(minNumber, nextNumber);
    return this.format(parsed.area, newNumber, protocol, parsed.raw);
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

  /**
   * 取得位址區域與暫存器屬性（是否為接點/Bit、型態標籤等）
   */
  getAreaInfo(address: string, protocol: ProtocolType = 'modbus_tcp'): {
    area: string;
    isBit: boolean;
    typeLabel: string;
  } {
    const parsed = this.parse(address, protocol);
    const isBit = parsed.type === 'bit';
    if (protocol === 'fatek_fbs' || protocol === 'mc_3e') {
      return {
        area: parsed.area,
        isBit,
        typeLabel: `${parsed.area} (${isBit ? 'Bit' : 'Word'})`,
      };
    }

    // Modbus 系
    const area = parsed.area;
    if (area === 'CS') return { area, isBit: true, typeLabel: 'coil' };
    if (area === 'IS') return { area, isBit: true, typeLabel: 'discrete_input' };
    if (area === 'IR') return { area, isBit: false, typeLabel: 'input_register' };
    return { area: 'HR', isBit: false, typeLabel: 'holding_register' };
  }
}

export const addressParser = new AddressParser();
