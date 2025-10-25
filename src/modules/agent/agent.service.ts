import { config } from '@/config';
import { eventBus } from '@/shared/core/event-bus';
import { AgentEvent, Intervention } from './types';

class AgentService {
  private isRunning = false;

  constructor() {
    if (config.agent.enabled) {
      this.initialize();
    } else {
      console.log('Agent is disabled by config.');
    }
  }

  private initialize() {
    console.log('🚀 Initializing Agent Service...');

    // 订阅事件
    eventBus.on('agent:event', this.handleEvent);
    eventBus.on('agent:intervention', this.handleIntervention);

    // TODO: 初始化 Listeners, Planner, Executors
    // this.initListeners();
    // this.initPlanner();
    // this.initExecutors();

    console.log('✅ Agent Service initialized.');
  }

  public start() {
    if (!config.agent.enabled || this.isRunning) {
      return;
    }
    console.log('▶️ Starting Agent...');
    this.isRunning = true;
    // TODO: 启动所有监听器
  }

  public stop() {
    if (!this.isRunning) {
      return;
    }
    console.log('⏹️ Stopping Agent...');
    this.isRunning = false;
    // TODO: 停止所有监听器
  }

  /**
   * 处理由监听器发布的事件
   * @param event AgentEvent
   */
  private handleEvent = (event: AgentEvent) => {
    if (!this.isRunning) return;

    console.log(`🧠 Received event: ${event.type} from ${event.source}`, event);
    // TODO: 将事件转发给 Planner 进行决策
  };

  /**
   * 处理由规划器制定的干预计划
   * @param intervention Intervention
   */
  private handleIntervention = (intervention: Intervention) => {
    if (!this.isRunning) return;

    console.log(`💪 Received intervention: ${intervention.type} from ${intervention.source}`, intervention);
    // TODO: 将干预计划转发给 Executor 执行
  };
}

// 创建并导出一个单例的 AgentService 实例
export const agentService = new AgentService();