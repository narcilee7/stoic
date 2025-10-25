import { config } from '@/config';
import { eventBus } from '@/shared/core/event-bus';
import { AgentEvent, Intervention } from './types';
import { CpuListener } from './listeners/cpu.listener';
import { PlannerService } from './planner/planner.service';   // 导入 Planner
import { ExecutorService } from './executors/executor.service'; // 导入 Executor

class AgentService {
  private isRunning = false;
  private listeners: { start: () => void; stop: () => void }[] = [];
  private planner?: PlannerService;
  private executor?: ExecutorService;

  constructor() {
    if (config.agent.enabled) {
      this.initialize();
    } else {
      console.log('Agent is disabled by config.');
    }
  }

  private initialize() {
    console.log('🚀 Initializing Agent Service...');

    // 订阅事件 (Agent Service 自身也可以订阅)
    eventBus.on('agent:event', this.handleEvent);
    eventBus.on('agent:intervention', this.handleIntervention);

    // 初始化所有子模块
    this.initListeners();
    this.initPlanner();
    this.initExecutors();

    console.log('✅ Agent Service initialized.');
  }

  private initListeners() {
    console.log('👂 Initializing listeners...');
    this.listeners.push(new CpuListener());
  }

  private initPlanner() {
    console.log('🧠 Initializing planner...');
    this.planner = new PlannerService();
  }

  private initExecutors() {
    console.log('💪 Initializing executors...');
    this.executor = new ExecutorService();
  }

  public start() {
    if (!config.agent.enabled || this.isRunning) return;
    console.log('▶️ Starting Agent...');
    this.isRunning = true;
    this.listeners.forEach(listener => listener.start());
  }

  public stop() {
    if (!this.isRunning) return;
    console.log('⏹️ Stopping Agent...');
    this.isRunning = false;
    this.listeners.forEach(listener => listener.stop());
  }

  private handleEvent = (event: AgentEvent) => {
    if (!this.isRunning) return;
    // Agent Service 自身也可以对事件做一些通用处理，比如记录日志
    // console.log(`[Agent Service] Logging event: ${event.type}`);
  };

  private handleIntervention = (intervention: Intervention) => {
    if (!this.isRunning) return;
    // Agent Service 自身也可以对干预做一些通用处理
  };
}

export const agentService = new AgentService();