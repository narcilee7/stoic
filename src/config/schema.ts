import { z } from 'zod';

// --- Agent Config Schema ---
export const agentConfigSchema = z.object({
  enabled: z.boolean().default(true),
  eventBufferSize: z.number().int().positive().default(1000),
  processInterval: z.number().int().positive().default(5000), // ms
  maxEventsPerBatch: z.number().int().positive().default(50),
  cpuThreshold: z.number().min(0).max(100).default(80), // %
  cpuWarningThreshold: z.number().min(0).max(100).default(70), // %
  cpuCriticalThreshold: z.number().min(0).max(100).default(90), // %
  cooldownPeriod: z.number().int().positive().default(30000), // ms
  privacyLevel: z.enum(['standard', 'strict', 'minimal']).default('standard'),
  notificationsEnabled: z.boolean().default(true), // 新增
});
export type AgentConfig = z.infer<typeof agentConfigSchema>;


// --- Database Config Schema ---
export const databaseConfigSchema = z.object({
  provider: z.enum(['sqlite']).default('sqlite'),
  url: z.string().default('database/stoic.db'),
  maxRetries: z.number().int().min(0).default(3),
  retryDelay: z.number().int().positive().default(1000), // ms
});
export type DatabaseConfig = z.infer<typeof databaseConfigSchema>;


// --- Widget Config Schema ---
export const widgetConfigSchema = z.object({
  theme: z.enum(['light', 'dark', 'system']).default('system'),
  refreshInterval: z.number().int().positive().default(60000), // ms
  showBreathing: z.boolean().default(true),
  showWordCloud: z.boolean().default(true),
});
export type WidgetConfig = z.infer<typeof widgetConfigSchema>;


// --- Main App Config Schema ---
export const appConfigSchema = z.object({
  logLevel: z.enum(['debug', 'info', 'warn', 'error']).default('info'),
  enableAnalytics: z.boolean().default(false),
  defaultUsername: z.string().default('default_user'),
});
export type AppConfig = z.infer<typeof appConfigSchema>;


// --- 合并所有配置 ---
export const fullConfigSchema = z.object({
  app: appConfigSchema,
  agent: agentConfigSchema,
  database: databaseConfigSchema,
  widget: widgetConfigSchema,
});
export type FullConfig = z.infer<typeof fullConfigSchema>;