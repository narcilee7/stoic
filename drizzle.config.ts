import { defineConfig } from 'drizzle-kit';

export default defineConfig({
    dialect: 'sqlite',
    schema: './src/infra/database/schema.ts',
    out: './database/migrations',
    dbCredentials: {
      url: process.env.DATABASE_URL || 'database/stoic.db',
    },
    verbose: true,
    strict: true,
})