-- 用户基础信息表
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 情绪记录表（核心表）
CREATE TABLE IF NOT EXISTS mood_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    mood_score REAL NOT NULL CHECK (mood_score >= 0 AND mood_score <= 1), -- 0=很差, 1=很好
    mood_label TEXT NOT NULL, -- happy, sad, angry, anxious, calm, confused
    confidence REAL CHECK (confidence >= 0 AND confidence <= 1), -- 情绪检测置信度
    context TEXT, -- 情绪上下文描述
    trigger_event TEXT, -- 触发事件类型: keyboard_burst, git_reset, idle_timeout, etc.
    trigger_value REAL, -- 触发事件的数值
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 聊天对话表（哲学家对话）
CREATE TABLE IF NOT EXISTS chat_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('user', 'assistant', 'system')), -- 消息角色
    content TEXT NOT NULL, -- 消息内容
    mood_at_time REAL, -- 当时情绪分数
    emotion_detected TEXT, -- 检测到的情绪
    token_count INTEGER, -- token数量（用于统计）
    model_used TEXT, -- 使用的AI模型
    response_time_ms INTEGER, -- 响应时间（毫秒）
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 呼吸训练记录表
CREATE TABLE IF NOT EXISTS breathe_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    duration_seconds INTEGER NOT NULL, -- 训练时长（秒）
    pattern_type TEXT NOT NULL, -- 呼吸模式: 4-7-8, box_breathing, etc.
    completed_cycles INTEGER, -- 完成呼吸周期数
    avg_heart_rate INTEGER, -- 平均心率（如果有设备）
    completion_rate REAL CHECK (completion_rate >= 0 AND completion_rate <= 1), -- 完成率
    session_quality TEXT CHECK (session_quality IN ('excellent', 'good', 'fair', 'poor')), -- 训练质量
    notes TEXT, -- 用户备注
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 尖叫室泄压记录表
CREATE TABLE IF NOT EXISTS scream_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    scream_text TEXT, -- 用户输入的泄压文本
    scream_intensity REAL CHECK (scream_intensity >= 0 AND scream_intensity <= 1), -- 泄压强度
    scream_mode TEXT CHECK (scream_mode IN ('normal', 'mute', 'visual')), -- 泄压模式
    duration_seconds INTEGER, -- 泄压时长
    character_count INTEGER, -- 输入字符数
    relief_score REAL CHECK (relief_score >= 0 AND relief_score <= 1), -- 泄压效果评分
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 认知训练记录表（最坏-最好-最真）
CREATE TABLE IF NOT EXISTS cognitive_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_type TEXT NOT NULL CHECK (session_type IN ('wcg', 'angle', 'wordcloud')), -- 训练类型
    problem_context TEXT, -- 问题背景
    worst_case TEXT, -- 最坏情况
    best_case TEXT, -- 最好情况
    most_realistic TEXT, -- 最真实情况
    final_thought TEXT, -- 最终想法
    anxiety_before REAL CHECK (anxiety_before >= 0 AND anxiety_before <= 1), -- 训练前焦虑
    anxiety_after REAL CHECK (anxiety_after >= 0 AND anxiety_after <= 1), -- 训练后焦虑
    helpfulness_score REAL CHECK (helpfulness_score >= 0 AND helpfulness_score <= 1), -- 有用性评分
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 日记记录表（加密存储）
CREATE TABLE IF NOT EXISTS diary_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    encrypted_content TEXT NOT NULL, -- AES-256加密的内容
    content_hash TEXT NOT NULL, -- 内容哈希（完整性验证）
    word_count INTEGER, -- 字数统计
    mood_summary REAL, -- 情绪摘要分数
    key_topics TEXT, -- 关键词主题（JSON格式）
    is_encrypted BOOLEAN DEFAULT 1, -- 是否加密
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 游戏记录表
CREATE TABLE IF NOT EXISTS game_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    game_type TEXT NOT NULL CHECK (game_type IN ('2048', 'snake', 'tictactoe', 'beach')), -- 游戏类型
    score INTEGER NOT NULL, -- 游戏分数
    duration_seconds INTEGER, -- 游戏时长
    completion_status TEXT CHECK (completion_status IN ('won', 'lost', 'quit', 'timeout')), -- 完成状态
    stress_before REAL CHECK (stress_before >= 0 AND stress_before <= 1), -- 游戏前压力
    stress_after REAL CHECK (stress_after >= 0 AND stress_after <= 1), -- 游戏后压力
    effectiveness_score REAL CHECK (effectiveness_score >= 0 AND effectiveness_score <= 1), -- 减压效果
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 番茄工作法记录表
CREATE TABLE IF NOT EXISTS pomodoro_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    task_name TEXT, -- 任务名称
    duration_minutes INTEGER NOT NULL CHECK (duration_minutes IN (25, 15, 5)), -- 时长（标准番茄）
    session_type TEXT CHECK (session_type IN ('work', 'short_break', 'long_break')), -- 会话类型
    completed BOOLEAN DEFAULT 0, -- 是否完成
    interruptions INTEGER DEFAULT 0, -- 中断次数
    focus_score REAL CHECK (focus_score >= 0 AND focus_score <= 1), -- 专注度评分
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 词云数据表
CREATE TABLE IF NOT EXISTS wordcloud_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    word TEXT NOT NULL, -- 词汇
    frequency INTEGER NOT NULL DEFAULT 1, -- 出现频率
    weight REAL NOT NULL, -- 权重（情绪强度）
    sentiment TEXT CHECK (sentiment IN ('positive', 'negative', 'neutral')), -- 情感倾向
    source_type TEXT CHECK (source_type IN ('chat', 'diary', 'scream', 'cognitive')), -- 来源类型
    source_id INTEGER, -- 来源记录ID
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, word, source_type, DATE(created_at)) -- 每天每词每来源唯一
);

-- 用户偏好设置表
CREATE TABLE IF NOT EXISTS user_preferences (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL UNIQUE,
    breathe_pattern TEXT DEFAULT '4-7-8', -- 默认呼吸模式
    scream_mode TEXT DEFAULT 'normal', -- 默认泄压模式
    pomodoro_duration INTEGER DEFAULT 25, -- 默认番茄时长
    notifications_enabled BOOLEAN DEFAULT 1, -- 通知开关
    sound_enabled BOOLEAN DEFAULT 1, -- 声音开关
    agent_enabled BOOLEAN DEFAULT 1, -- Agent开关
    theme TEXT DEFAULT 'calm', -- 主题
    language TEXT DEFAULT 'zh', -- 语言
    privacy_level TEXT DEFAULT 'standard', -- 隐私级别
    export_password_hash TEXT, -- 导出密码哈希
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Agent干预记录表
CREATE TABLE IF NOT EXISTS agent_interventions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    intervention_type TEXT NOT NULL, -- 干预类型: breathe, scream, q, etc.
    trigger_event TEXT NOT NULL, -- 触发事件
    trigger_value REAL, -- 触发值
    predicted_effectiveness REAL CHECK (predicted_effectiveness >= 0 AND predicted_effectiveness <= 1), -- 预测效果
    actual_effectiveness REAL CHECK (actual_effectiveness >= 0 AND actual_effectiveness <= 1), -- 实际效果
    timing_score REAL CHECK (timing_score >= 0 AND timing_score <= 1), -- 时机评分
    user_feedback TEXT CHECK (user_feedback IN ('helpful', 'neutral', 'annoying')), -- 用户反馈
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 系统日志表（用于调试和审计）
CREATE TABLE IF NOT EXISTS system_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    level TEXT NOT NULL CHECK (level IN ('debug', 'info', 'warning', 'error', 'fatal')), -- 日志级别
    component TEXT NOT NULL, -- 组件: agent, breathe, cognitive, etc.
    message TEXT NOT NULL, -- 日志消息
    context TEXT, -- 上下文信息（JSON格式）
    error_stack TEXT, -- 错误堆栈
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 数据导出记录表（审计用途）
CREATE TABLE IF NOT EXISTS export_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    export_type TEXT NOT NULL, -- 导出类型: full, partial, encrypted
    file_path TEXT NOT NULL, -- 导出文件路径
    file_size INTEGER, -- 文件大小（字节）
    checksum TEXT, -- 文件校验和
    encryption_method TEXT, -- 加密方法
    export_duration_ms INTEGER, -- 导出耗时（毫秒）
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 创建索引优化查询性能
CREATE INDEX IF NOT EXISTS idx_mood_entries_user_date ON mood_entries(user_id, DATE(created_at));
CREATE INDEX IF NOT EXISTS idx_mood_entries_score ON mood_entries(mood_score);
CREATE INDEX IF NOT EXISTS idx_chat_messages_user_date ON chat_messages(user_id, DATE(created_at));
CREATE INDEX IF NOT EXISTS idx_breathe_sessions_user_date ON breathe_sessions(user_id, DATE(created_at));
CREATE INDEX IF NOT EXISTS idx_diary_entries_user_date ON diary_entries(user_id, DATE(created_at));
CREATE INDEX IF NOT EXISTS idx_wordcloud_user_word ON wordcloud_entries(user_id, word);
CREATE INDEX IF NOT EXISTS idx_agent_interventions_user_date ON agent_interventions(user_id, DATE(created_at));

-- 创建视图：每日情绪摘要
CREATE VIEW IF NOT EXISTS daily_mood_summary AS
SELECT 
    user_id,
    DATE(created_at) as date,
    AVG(mood_score) as avg_mood,
    MIN(mood_score) as min_mood,
    MAX(mood_score) as max_mood,
    COUNT(*) as entry_count,
    GROUP_CONCAT(DISTINCT mood_label) as mood_labels
FROM mood_entries
GROUP BY user_id, DATE(created_at);

-- 创建视图：每周活动摘要
CREATE VIEW IF NOT EXISTS weekly_activity_summary AS
SELECT 
    user_id,
    strftime('%Y-W%W', created_at) as week,
    COUNT(DISTINCT CASE WHEN source_type = 'breathe' THEN source_id END) as breathe_sessions,
    COUNT(DISTINCT CASE WHEN source_type = 'scream' THEN source_id END) as scream_sessions,
    COUNT(DISTINCT CASE WHEN source_type = 'cognitive' THEN source_id END) as cognitive_sessions,
    COUNT(DISTINCT CASE WHEN source_type = 'game' THEN source_id END) as game_sessions,
    COUNT(DISTINCT CASE WHEN source_type = 'diary' THEN source_id END) as diary_entries,
    AVG(CASE WHEN source_type = 'breathe' THEN weight END) as avg_breathe_effect
FROM (
    SELECT user_id, created_at, 'breathe' as source_type, id as source_id, completion_rate as weight FROM breathe_sessions
    UNION ALL
    SELECT user_id, created_at, 'scream' as source_type, id as source_id, relief_score as weight FROM scream_sessions
    UNION ALL
    SELECT user_id, created_at, 'cognitive' as source_type, id as source_id, helpfulness_score as weight FROM cognitive_sessions
    UNION ALL
    SELECT user_id, created_at, 'game' as source_type, id as source_id, effectiveness_score as weight FROM game_sessions
    UNION ALL
    SELECT user_id, created_at, 'diary' as source_type, id as source_id, NULL as weight FROM diary_entries
)
GROUP BY user_id, strftime('%Y-W%W', created_at);

-- 插入默认用户（如果不存在）
INSERT OR IGNORE INTO users (username) VALUES ('default_user');

-- 插入默认用户偏好设置
INSERT OR IGNORE INTO user_preferences (user_id) 
SELECT id FROM users WHERE username = 'default_user';