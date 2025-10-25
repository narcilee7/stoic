import TOML from '@ltd/j-toml';
import fs from 'node:fs';
import path from 'node:path';
import { fullConfigSchema, FullConfig } from './schema';

const CONFIG_DIR = path.join(process.cwd(), 'configs');

/**
 * 加载并解析单个 TOML 配置文件
 * @param filePath 文件的绝对路径
 * @returns 解析后的对象，如果文件不存在或解析失败则返回空对象
 */
function loadConfigFile(fileName: string): object {
  const filePath = path.join(CONFIG_DIR, fileName);
  try {
    if (!fs.existsSync(filePath)) {
      console.warn(`🤔 Config file not found: ${filePath}. Using defaults.`);
      return {};
    }
    const fileContent = fs.readFileSync(filePath, 'utf-8');
    // 使用 TOML.parse，并强制转换为 object
    return TOML.parse(fileContent) as unknown as object;
  } catch (error) {
    console.error(`❌ Error parsing config file ${filePath}:`, error);
    // 在解析失败时返回空对象，以便 Zod 可以使用默认值
    return {};
  }
}

/**
 * 加载、合并并验证所有配置文件
 * @returns 一个经过验证和类型安全的配置对象
 */
function loadAndValidateConfig(): FullConfig {
  // 1. 加载各个配置文件
  const appConfigData = loadConfigFile('config.toml');
  const agentConfigData = loadConfigFile('agent.toml');
  const databaseConfigData = loadConfigFile('database.toml');
  const widgetConfigData = loadConfigFile('widget.toml');

  // 2. 创建一个原始的、未经验证的配置对象结构
  const rawConfig = {
    app: appConfigData,
    agent: agentConfigData,
    database: databaseConfigData,
    widget: widgetConfigData,
  };

  // 3. 使用 Zod 进行解析和验证
  // .parse() 会在验证失败时抛出错误，确保我们不会在配置错误的情况下运行应用
  const validatedConfig = fullConfigSchema.parse(rawConfig);

  // 4. 返回一个深度冻结的对象，防止在运行时被意外修改
  return Object.freeze(validatedConfig);
}

// --- 导出单例配置 ---
export const config = loadAndValidateConfig();

// --- 打印加载的配置 (用于调试) ---
console.log('✅ Configuration loaded successfully:');
console.log(JSON.stringify(config, null, 2));
