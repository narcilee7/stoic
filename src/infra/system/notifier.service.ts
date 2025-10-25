import notifier from 'node-notifier';
import path from 'node:path';
import { config } from '@/config';

const APP_ICON = path.join(process.cwd(), 'assets', 'icons', 'stoic-logo.png');

export interface NotificationOptions {
  title: string;
  message: string;
  subtitle?: string;
  sound?: boolean; // true | false | 'Frog' | 'Glass'
  wait?: boolean; // 等待用户交互
  timeout?: number; // 秒
  open?: string; // URL to open on click
  actions?: string[]; // e.g., ['Snooze', 'Dismiss']
}

class NotifierService {
  constructor() {
    console.log('💡 Notifier Service initialized.');
  }

  public notify(options: NotificationOptions): Promise<void> {
    return new Promise((resolve, reject) => {
      if (!config.agent.notificationsEnabled) {
        console.log('[Notifier] Notifications are disabled by config. Skipping.');
        return resolve();
      }

      notifier.notify(
        {
          ...options,
          icon: APP_ICON,
          contentImage: undefined, // 可选
        },
        (error, response, metadata) => {
          if (error) {
            console.error('❌ Notification failed:', error);
            return reject(error);
          }
          console.log('✅ Notification sent successfully. Response:', response, metadata);
          resolve();
        },
      );
    });
  }
}

// 创建并导出一个单例的 NotifierService 实例
export const notifierService = new NotifierService();