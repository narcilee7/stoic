import { eventBus } from '@/shared/core/event-bus';
import { AgentEvent, Intervention, interventionSchema } from '../types';
import { v4 as uuidv4 } from 'uuid';

const SOURCE = 'simple-rules-planner';

export class PlannerService {
  constructor() {
    console.log('💡 Planner Service initialized.');
    this.subscribeToEvents();
  }

  private subscribeToEvents() {
    eventBus.on('agent:event', this.handleEvent);
  }

  private handleEvent = (event: AgentEvent) => {
    console.log(`[Planner] Received event: ${event.type}`);

    // 根据事件类型进行决策
    switch (event.type) {
      case 'cpu_usage_warning':
      case 'cpu_usage_critical':
        this.planBreathingExercise(event);
        break;
      // TODO: 在这里为其他事件类型添加决策逻辑
      default:
        // console.log(`[Planner] No action defined for event type: ${event.type}`);
        break;
    }
  };

  /**
   * 计划一个呼吸练习干预
   * @param event 触发此计划的原始事件
   */
  private planBreathingExercise(event: AgentEvent) {
    const reason = `Detected ${event.severity} CPU usage (${(event.value! * 100).toFixed(0)}%). A short break could be helpful.`;

    const intervention: Intervention = {
      id: uuidv4(),
      type: 'suggest_breathing_exercise',
      source: SOURCE,
      reason: reason,
      timestamp: new Date(),
      urgency: event.severity === 'critical' ? 0.9 : 0.6,
      parameters: {
        duration: 60, // 建议时长60秒
        pattern: '4-7-8',
      },
    };

    // 验证干预计划的结构
    const validationResult = interventionSchema.safeParse(intervention);
    if (validationResult.success) {
      console.log(`[Planner] Firing intervention: ${intervention.type}`);
      // 将干预计划发布到事件总线
      eventBus.emit('agent:intervention', validationResult.data);
    } else {
      console.error('❌ Invalid intervention schema:', validationResult.error);
    }
  }
}