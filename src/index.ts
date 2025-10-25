import { config } from '@/config';
import { runMigrations } from '@/infra/database';
import { agentService } from '@/modules/agent/agent.service'; // 导入 Agent 服务

/**
 * 主应用函数
 */
async function main() {
  console.log(`🚀 Stoic Agent v${process.env.npm_package_version} is starting...`);
  console.log(`Log level set to: ${config.app.logLevel}`);

  runMigrations();

  // 启动 Agent 服务
  agentService.start();

  console.log('✅ Stoic Agent started successfully.');

  // TODO: 在这里启动CLI或Agent主循环
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