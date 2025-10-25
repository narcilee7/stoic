import { integer, sqliteTable, text, real, primaryKey } from 'drizzle-orm/sqlite-core';
import { relations } from 'drizzle-orm';

// --- 核心模型 ---

export const users = sqliteTable('users', {
  id: integer('id').primaryKey({ autoIncrement: true }),
  username: text('username').notNull().unique(),
  createdAt: integer('created_at', { mode: 'timestamp' }).default(new Date()),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).default(new Date()),
});

export const moodEntries = sqliteTable('mood_entries', {
  id: integer('id').primaryKey({ autoIncrement: true }),
  userId: integer('user_id').notNull().references(() => users.id, { onDelete: 'cascade' }),
  moodScore: real('mood_score').notNull(),
  moodLabel: text('mood_label').notNull(),
  confidence: real('confidence'),
  context: text('context'),
  triggerEvent: text('trigger_event'),
  triggerValue: real('trigger_value'),
  createdAt: integer('created_at', { mode: 'timestamp' }).default(new Date()),
  updatedAt: integer('updated_at', { mode: 'timestamp' }).default(new Date()),
});

export const chatMessages = sqliteTable('chat_messages', {
    id: integer('id').primaryKey({ autoIncrement: true }),
    userId: integer('user_id').notNull().references(() => users.id, { onDelete: 'cascade' }),
    philosopherId: text('philosopher_id').notNull(),
    role: text('role').notNull(),
    content: text('content').notNull(),
    moodAtTime: real('mood_at_time'),
    emotionDetected: text('emotion_detected'),
    tokenCount: integer('token_count'),
    modelUsed: text('model_used'),
    responseTimeMs: integer('response_time_ms'),
    createdAt: integer('created_at', { mode: 'timestamp' }).default(new Date()),
});


// --- 功能模块模型 ---

export const breatheSessions = sqliteTable('breathe_sessions', {
    id: integer('id').primaryKey({ autoIncrement: true }),
    userId: integer('user_id').notNull().references(() => users.id, { onDelete: 'cascade' }),
    durationSeconds: integer('duration_seconds').notNull(),
    patternType: text('pattern_type').notNull(),
    completedCycles: integer('completed_cycles'),
    avgHeartRate: integer('avg_heart_rate'),
    completionRate: real('completion_rate'),
    sessionQuality: text('session_quality'),
    notes: text('notes'),
    createdAt: integer('created_at', { mode: 'timestamp' }).default(new Date()),
});

export const screamSessions = sqliteTable('scream_sessions', {
    id: integer('id').primaryKey({ autoIncrement: true }),
    userId: integer('user_id').notNull().references(() => users.id, { onDelete: 'cascade' }),
    screamText: text('scream_text'),
    screamIntensity: real('scream_intensity'),
    screamMode: text('scream_mode').notNull(),
    durationSeconds: integer('duration_seconds'),
    characterCount: integer('character_count'),
    reliefScore: real('relief_score'),
    createdAt: integer('created_at', { mode: 'timestamp' }).default(new Date()),
});

export const cognitiveSessions = sqliteTable('cognitive_sessions', {
    id: integer('id').primaryKey({ autoIncrement: true }),
    userId: integer('user_id').notNull().references(() => users.id, { onDelete: 'cascade' }),
    sessionType: text('session_type').notNull(),
    problemContext: text('problem_context'),
    worstCase: text('worst_case'),
    bestCase: text('best_case'),
    mostRealistic: text('most_realistic'),
    finalThought: text('final_thought'),
    anxietyBefore: real('anxiety_before'),
    anxietyAfter: real('anxiety_after'),
    helpfulnessScore: real('helpfulness_score'),
    createdAt: integer('created_at', { mode: 'timestamp' }).default(new Date()),
});

export const diaryEntries = sqliteTable('diary_entries', {
    id: integer('id').primaryKey({ autoIncrement: true }),
    userId: integer('user_id').notNull().references(() => users.id, { onDelete: 'cascade' }),
    encryptedContent: text('encrypted_content').notNull(),
    contentHash: text('content_hash').notNull(),
    wordCount: integer('word_count'),
    moodSummary: real('mood_summary'),
    keyTopics: text('key_topics'), // JSON
    isEncrypted: integer('is_encrypted', { mode: 'boolean' }).default(true),
    createdAt: integer('created_at', { mode: 'timestamp' }).default(new Date()),
    updatedAt: integer('updated_at', { mode: 'timestamp' }).default(new Date()),
});

export const gameSessions = sqliteTable('game_sessions', {
    id: integer('id').primaryKey({ autoIncrement: true }),
    userId: integer('user_id').notNull().references(() => users.id, { onDelete: 'cascade' }),
    gameType: text('game_type').notNull(),
    score: integer('score').notNull(),
    durationSeconds: integer('duration_seconds'),
    completionStatus: text('completion_status'),
    stressBefore: real('stress_before'),
    stressAfter: real('stress_after'),
    effectivenessScore: real('effectiveness_score'),
    createdAt: integer('created_at', { mode: 'timestamp' }).default(new Date()),
});

export const pomodoroSessions = sqliteTable('pomodoro_sessions', {
    id: integer('id').primaryKey({ autoIncrement: true }),
    userId: integer('user_id').notNull().references(() => users.id, { onDelete: 'cascade' }),
    taskName: text('task_name'),
    durationMinutes: integer('duration_minutes').notNull(),
    sessionType: text('session_type').notNull(),
    completed: integer('completed', { mode: 'boolean' }).default(false),
    interruptions: integer('interruptions').default(0),
    ocusScore: real('focus_score'),
    createdAt: integer('created_at', { mode: 'timestamp' }).default(new Date()),
});

export const wordcloudEntries = sqliteTable('wordcloud_entries', {
    id: integer('id').primaryKey({ autoIncrement: true }),
    userId: integer('user_id').notNull().references(() => users.id, { onDelete: 'cascade' }),
    word: text('word').notNull(),
    frequency: integer('frequency').default(1).notNull(),
    weight: real('weight').notNull(),
    sentiment: text('sentiment'),
    sourceType: text('source_type'),
    sourceId: integer('source_id'),
    createdAt: integer('created_at', { mode: 'timestamp' }).default(new Date()),
}, (table) => ({
    pk: primaryKey({ columns: [table.userId, table.word, table.sourceType, table.createdAt] }),
}));


// --- 配置与日志模型 ---

export const userPreferences = sqliteTable('user_preferences', {
    id: integer('id').primaryKey({ autoIncrement: true }),
    userId: integer('user_id').notNull().unique().references(() => users.id, { onDelete: 'cascade' }),
    breathePattern: text('breathe_pattern').default('4-7-8'),
    screamMode: text('scream_mode').default('normal'),
    pomodoroDuration: integer('pomodoro_duration').default(25),
    notificationsEnabled: integer('notifications_enabled', { mode: 'boolean' }).default(true),
    soundEnabled: integer('sound_enabled', { mode: 'boolean' }).default(true),
    agentEnabled: integer('agent_enabled', { mode: 'boolean' }).default(true),
    theme: text('theme').default('calm'),
    language: text('language').default('zh'),
    privacyLevel: text('privacy_level').default('standard'),
    exportPasswordHash: text('export_password_hash'),
    createdAt: integer('created_at', { mode: 'timestamp' }).default(new Date()),
    updatedAt: integer('updated_at', { mode: 'timestamp' }).default(new Date()),
});

export const agentInterventions = sqliteTable('agent_interventions', {
    id: integer('id').primaryKey({ autoIncrement: true }),
    userId: integer('user_id').notNull().references(() => users.id, { onDelete: 'cascade' }),
    interventionType: text('intervention_type').notNull(),
    triggerEvent: text('trigger_event').notNull(),
    triggerValue: real('trigger_value'),
    predictedEffectiveness: real('predicted_effectiveness'),
    actualEffectiveness: real('actual_effectiveness'),
    timingScore: real('timing_score'),
    userFeedback: text('user_feedback'),
    createdAt: integer('created_at', { mode: 'timestamp' }).default(new Date()),
});

export const systemLogs = sqliteTable('system_logs', {
    id: integer('id').primaryKey({ autoIncrement: true }),
    level: text('level').notNull(),
    component: text('component').notNull(),
    message: text('message').notNull(),
    context: text('context'), // JSON
    errorStack: text('error_stack'),
    createdAt: integer('created_at', { mode: 'timestamp' }).default(new Date()),
});

export const exportLogs = sqliteTable('export_logs', {
    id: integer('id').primaryKey({ autoIncrement: true }),
    userId: integer('user_id').notNull().references(() => users.id, { onDelete: 'cascade' }),
    exportType: text('export_type').notNull(),
    filePath: text('file_path').notNull(),
    fileSize: integer('file_size'),
    checksum: text('checksum'),
    encryptionMethod: text('encryption_method'),
    exportDurationMs: integer('export_duration_ms'),
    createdAt: integer('created_at', { mode: 'timestamp' }).default(new Date()),
});


// --- 关系定义 ---

export const usersRelations = relations(users, ({ one, many }) => ({
    preferences: one(userPreferences),
    moodEntries: many(moodEntries),
    chatMessages: many(chatMessages),
    breatheSessions: many(breatheSessions),
    screamSessions: many(screamSessions),
    cognitiveSessions: many(cognitiveSessions),
    diaryEntries: many(diaryEntries),
    gameSessions: many(gameSessions),
    pomodoroSessions: many(pomodoroSessions),
    wordcloudEntries: many(wordcloudEntries),
    agentInterventions: many(agentInterventions),
    exportLogs: many(exportLogs),
}));

export const moodEntriesRelations = relations(moodEntries, ({ one }) => ({
    user: one(users, {
        fields: [moodEntries.userId],
        references: [users.id],
    }),
}));