import { eventBus } from '@/shared/core/event-bus';
import { Intervention } from '../types';
import { notifierService, NotificationOptions } from '@/infra/system/notifier.service';
import chalk from 'chalk';
import boxen from 'boxen';

export class ExecutorService {
  constructor() {
    console.log('💡 Executor Service initialized.');
    this.subscribeToInterventions();
  }

  private subscribeToInterventions() {
    eventBus.on('agent:intervention', this.handleIntervention);
  }

  private handleIntervention = (intervention: Intervention) => {
    console.log(`[Executor] Received intervention: ${intervention.type}`);

    switch (intervention.type) {
      case 'suggest_breathing_exercise':
        this.executeBreathingSuggestion(intervention);
        break;
      default:
        console.log(`[Executor] No action defined for intervention type: ${intervention.type}`);
        break;
    }
  };

  private async executeBreathingSuggestion(intervention: Intervention) {
    const title = 'High CPU Usage Detected!';
    const message = `How about a quick ${intervention.parameters?.duration}-second breathing exercise?`;
    
    console.log(`[Executor] Sending notification: "${title} - ${message}"`);

    try {
      await notifierService.notify({
        title: title,
        message: message,
        subtitle: 'Stoic Agent Suggestion',
        sound: true,
        wait: true, // 等待用户点击
        actions: ['Start', 'Dismiss'],
        timeout: 30,
      });
      // TODO: 在这里可以处理用户的点击行为，例如 'Start'
    } catch (error) {
      console.error('❌ Failed to send breathing suggestion notification:', error);
    }
  }
}