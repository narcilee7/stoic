import { drizzle } from 'drizzle-orm/better-sqlite3';
import Database from 'better-sqlite3';
import { migrate } from 'drizzle-orm/better-sqlite3/migrator';
import * as schema from './schema';
import path from 'node:path';
import fs from 'node:fs';

const dbPath = process.env.DATABASE_URL || 'database/stoic.db';

const dbDir = path.dirname(dbPath);
if (!fs.existsSync(dbDir)) {
  fs.mkdirSync(dbDir, { recursive: true });
  console.log(`Created database directory: ${dbDir}`);
}

const sqlite = new Database(dbPath, {
    verbose: console.log, // 开启 verbose 模式，打印所有 SQL 语句
})

export const db = drizzle(sqlite, { schema, logger: true });

export const runMigrations = () => {
    console.log('Running migrations...');
    try {
        migrate(db, { migrationsFolder: './database/migrations' });
    } catch (error) {
        console.error('❌ Error running migrations:', error);
        process.exit(1);
    }
    console.log('Migrations completed.');
}