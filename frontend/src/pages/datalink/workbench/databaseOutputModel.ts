import type { DatabaseWriteMode } from '../../../types/datalink';

export type DatabaseOutputScope = {
  connectorId: string;
  tableKey: string;
  columnName: string;
  writeMode: DatabaseWriteMode;
  timestampColumn: string;
};

export const DATABASE_OUTPUT_SCOPE_INITIAL: DatabaseOutputScope = {
  connectorId: '',
  tableKey: '',
  columnName: '',
  writeMode: 'insert',
  timestampColumn: '',
};
