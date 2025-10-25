CREATE TABLE `agent_interventions` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`user_id` integer NOT NULL,
	`intervention_type` text NOT NULL,
	`trigger_event` text NOT NULL,
	`trigger_value` real,
	`predicted_effectiveness` real,
	`actual_effectiveness` real,
	`timing_score` real,
	`user_feedback` text,
	`created_at` integer DEFAULT '"2025-10-25T07:04:01.456Z"',
	FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE TABLE `breathe_sessions` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`user_id` integer NOT NULL,
	`duration_seconds` integer NOT NULL,
	`pattern_type` text NOT NULL,
	`completed_cycles` integer,
	`avg_heart_rate` integer,
	`completion_rate` real,
	`session_quality` text,
	`notes` text,
	`created_at` integer DEFAULT '"2025-10-25T07:04:01.455Z"',
	FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE TABLE `chat_messages` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`user_id` integer NOT NULL,
	`philosopher_id` text NOT NULL,
	`role` text NOT NULL,
	`content` text NOT NULL,
	`mood_at_time` real,
	`emotion_detected` text,
	`token_count` integer,
	`model_used` text,
	`response_time_ms` integer,
	`created_at` integer DEFAULT '"2025-10-25T07:04:01.455Z"',
	FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE TABLE `cognitive_sessions` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`user_id` integer NOT NULL,
	`session_type` text NOT NULL,
	`problem_context` text,
	`worst_case` text,
	`best_case` text,
	`most_realistic` text,
	`final_thought` text,
	`anxiety_before` real,
	`anxiety_after` real,
	`helpfulness_score` real,
	`created_at` integer DEFAULT '"2025-10-25T07:04:01.455Z"',
	FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE TABLE `diary_entries` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`user_id` integer NOT NULL,
	`encrypted_content` text NOT NULL,
	`content_hash` text NOT NULL,
	`word_count` integer,
	`mood_summary` real,
	`key_topics` text,
	`is_encrypted` integer DEFAULT true,
	`created_at` integer DEFAULT '"2025-10-25T07:04:01.455Z"',
	`updated_at` integer DEFAULT '"2025-10-25T07:04:01.455Z"',
	FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE TABLE `export_logs` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`user_id` integer NOT NULL,
	`export_type` text NOT NULL,
	`file_path` text NOT NULL,
	`file_size` integer,
	`checksum` text,
	`encryption_method` text,
	`export_duration_ms` integer,
	`created_at` integer DEFAULT '"2025-10-25T07:04:01.456Z"',
	FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE TABLE `game_sessions` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`user_id` integer NOT NULL,
	`game_type` text NOT NULL,
	`score` integer NOT NULL,
	`duration_seconds` integer,
	`completion_status` text,
	`stress_before` real,
	`stress_after` real,
	`effectiveness_score` real,
	`created_at` integer DEFAULT '"2025-10-25T07:04:01.456Z"',
	FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE TABLE `mood_entries` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`user_id` integer NOT NULL,
	`mood_score` real NOT NULL,
	`mood_label` text NOT NULL,
	`confidence` real,
	`context` text,
	`trigger_event` text,
	`trigger_value` real,
	`created_at` integer DEFAULT '"2025-10-25T07:04:01.455Z"',
	`updated_at` integer DEFAULT '"2025-10-25T07:04:01.455Z"',
	FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE TABLE `pomodoro_sessions` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`user_id` integer NOT NULL,
	`task_name` text,
	`duration_minutes` integer NOT NULL,
	`session_type` text NOT NULL,
	`completed` integer DEFAULT false,
	`interruptions` integer DEFAULT 0,
	`focus_score` real,
	`created_at` integer DEFAULT '"2025-10-25T07:04:01.456Z"',
	FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE TABLE `scream_sessions` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`user_id` integer NOT NULL,
	`scream_text` text,
	`scream_intensity` real,
	`scream_mode` text NOT NULL,
	`duration_seconds` integer,
	`character_count` integer,
	`relief_score` real,
	`created_at` integer DEFAULT '"2025-10-25T07:04:01.455Z"',
	FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE TABLE `system_logs` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`level` text NOT NULL,
	`component` text NOT NULL,
	`message` text NOT NULL,
	`context` text,
	`error_stack` text,
	`created_at` integer DEFAULT '"2025-10-25T07:04:01.456Z"'
);
--> statement-breakpoint
CREATE TABLE `user_preferences` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`user_id` integer NOT NULL,
	`breathe_pattern` text DEFAULT '4-7-8',
	`scream_mode` text DEFAULT 'normal',
	`pomodoro_duration` integer DEFAULT 25,
	`notifications_enabled` integer DEFAULT true,
	`sound_enabled` integer DEFAULT true,
	`agent_enabled` integer DEFAULT true,
	`theme` text DEFAULT 'calm',
	`language` text DEFAULT 'zh',
	`privacy_level` text DEFAULT 'standard',
	`export_password_hash` text,
	`created_at` integer DEFAULT '"2025-10-25T07:04:01.456Z"',
	`updated_at` integer DEFAULT '"2025-10-25T07:04:01.456Z"',
	FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE UNIQUE INDEX `user_preferences_user_id_unique` ON `user_preferences` (`user_id`);--> statement-breakpoint
CREATE TABLE `users` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`username` text NOT NULL,
	`created_at` integer DEFAULT '"2025-10-25T07:04:01.455Z"',
	`updated_at` integer DEFAULT '"2025-10-25T07:04:01.455Z"'
);
--> statement-breakpoint
CREATE UNIQUE INDEX `users_username_unique` ON `users` (`username`);--> statement-breakpoint
CREATE TABLE `wordcloud_entries` (
	`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL,
	`user_id` integer NOT NULL,
	`word` text NOT NULL,
	`frequency` integer DEFAULT 1 NOT NULL,
	`weight` real NOT NULL,
	`sentiment` text,
	`source_type` text,
	`source_id` integer,
	`created_at` integer DEFAULT '"2025-10-25T07:04:01.456Z"',
	PRIMARY KEY(`user_id`, `word`, `source_type`, `created_at`),
	FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON UPDATE no action ON DELETE cascade
);
