import { z } from 'zod';

/**
 * 環境變數 schema（Vite 僅會注入 VITE_ 前綴的變數，建置後未設定者為 undefined）
 * 在傳入 Zod 前先補預設值，確保絕不傳入 undefined，避免 ZodError
 */
const envSchema = z.object({
  VITE_API_BASE_URL: z.string(),
});

const raw = {
  VITE_API_BASE_URL: import.meta.env.VITE_API_BASE_URL ?? '/api/v1',
};

export const env = envSchema.parse(raw);

export const VITE_API_BASE_URL = env.VITE_API_BASE_URL;
