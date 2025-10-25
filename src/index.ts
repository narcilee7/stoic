import { runMigrations } from '@/infra/database';
import { db } from '@/infra/database';
import { users } from '@/infra/database/schema';

/**
 * 主应用函数
 */
async function main() {
  console.log('🚀 Stoic Agent is starting...');

  // 运行数据库迁移
  runMigrations();

  // 2. (示例) 插入一个新用户并查询
  try {
    console.log('Inserting a new user...');
    await db.insert(users).values({ username: 'narcilee' }).onConflictDoNothing();

    console.log('Querying all users...');
    const allUsers = await db.select().from(users);
    console.log('All users:', allUsers);
  } catch (error) {
    console.error('Database operation failed:', error);
  }

  console.log('✅ Stoic Agent started successfully.');
  // 3. 在这里启动CLI或Agent主循环
  // ...
}

// 捕获未处理的异常
process.on('uncaughtException', (error) => {
  console.error('😱 Uncaught Exception:', error);
  process.exit(1);
});

// 捕获未处理的Promise拒绝
process.on('unhandledRejection', (reason, promise) => {
  console.error('😡 Unhandled Rejection at:', promise, 'reason:', reason);
  process.exit(1);
});

// 运行主函数
main();