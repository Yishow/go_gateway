export interface DatabaseDeliveryOutcomeFields {
  last_schema_ensure_at?: string | null;
  last_schema_ensure_status?: string;
  last_schema_ensure_error?: string;
  last_write_at?: string | null;
  last_write_status?: string;
  last_write_error?: string;
  last_flush_at?: string | null;
  last_flush_status?: string;
  last_flush_error?: string;
}
