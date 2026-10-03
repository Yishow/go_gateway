import { AddressParser, getDefaultPlannerStartAddress } from '../../../../utils/addressParser';
import type { ProtocolType } from '../../../../types/datalink';
import type { ProtocolId } from './types';

export interface ProtocolSupport {
  available: boolean;
  /** i18n key (workbench-v2 namespace) explaining why the protocol cannot be selected. */
  reasonKey?: string;
}

const PARSER = new AddressParser();

/**
 * A protocol may be offered only when the whole V2 path works for it. The
 * address parser is the part the frontend can check directly, so availability
 * follows it: an option can never outrun what Step 2 can actually parse.
 * (The collector and setup paths are verified server-side and reject the same
 * protocols.)
 */
export function getProtocolSupport(protocol: ProtocolId): ProtocolSupport {
  try {
    PARSER.parse(getDefaultPlannerStartAddress(protocol as ProtocolType), protocol as ProtocolType);
    return { available: true };
  } catch {
    return { available: false, reasonKey: 'step1.protocols.unavailable_no_v2_path' };
  }
}
