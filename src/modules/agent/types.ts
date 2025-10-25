import { z } from 'zod';

// -- Event --
export const eventTypeSchema = z.enum([
  // 系统事件
  'cpu_usage_high',
  'memory_usage_high',
  // 用户行为事件
  'keyboard_burst',
  'mouse_rapid',
  'idle_detected',
  // 开发行为事件
  'git_reset_frequent',
  'build_failed',
]);

export type EventType = z.infer<typeof eventTypeSchema>;

// 事件严重等级
export const eventSeveritySchema = z.enum(['low', 'medium', 'high', 'critical']);
export type EventSeverity = z.infer<typeof eventSeveritySchema>;

export const agentEventSchema = z.object({
  id: z.string().uuid(),
  type: eventTypeSchema,
  source: z.string(), // e.g., 'cpu-listener'
  severity: eventSeveritySchema,
  timestamp: z.date(),
  value: z.number().optional(), // 可选的标准化数值 (0-1)
  metadata: z.record(z.any(), z.string()).optional(),
});

export type AgentEvent = z.infer<typeof agentEventSchema>;

// 干预
export const interventionTypeSchema = z.enum([
  'suggest_breathing_exercise',
  'suggest_scream_session',
  'ask_cognitive_question',
  'show_motivational_quote',
]);

export type InterventionType = z.infer<typeof interventionTypeSchema>;

export const interventionSchema = z.object({
  id: z.string().uuid(),
  type: interventionTypeSchema,
  source: z.string(), // e.g., 'stress-planner'
  reason: z.string(), // 为什么触发这个干预
  timestamp: z.date(),
  urgency: z.number().min(0).max(1), // 紧急程度 (0-1)
  parameters: z.record(z.any(), z.string()).optional(), // 执行干预所需的参数
});
export type Intervention = z.infer<typeof interventionSchema>;