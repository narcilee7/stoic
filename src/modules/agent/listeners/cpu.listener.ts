import { config } from '@/config';
import { eventBus } from '@/shared/core/event-bus';
import { agentEventSchema, AgentEvent } from '../types';
import si from 'systeminformation';
import { v4 as uuidv4 } from 'uuid';
import { movingAverage, MovingAverage } from '@/shared/utils/moving-average';

const SOURCE = 'cpu-listener';

type CpuState = 'normal' | 'warning' | 'critical';

export class CpuListener {
  private intervalId: NodeJS.Timeout | null = null;
  private history: number[] = [];
  private state: CpuState = 'normal';
  private movingAverage: MovingAverage;

  constructor() {
    const historySize = Math.round(config.agent.cooldownPeriod / config.agent.processInterval);
    this.movingAverage = movingAverage(historySize);
    console.log(`💡 CPU Listener initialized (history size: ${historySize}).`);
  }

  public start() {
    if (this.intervalId) return;
    console.log('▶️ Starting CPU Listener...');
    this.checkCpuUsage();
    this.intervalId = setInterval(this.checkCpuUsage, config.agent.processInterval);
  }

  public stop() {
    if (this.intervalId) {
      console.log('⏹️ Stopping CPU Listener...');
      clearInterval(this.intervalId);
      this.intervalId = null;
    }
  }

  private checkCpuUsage = async () => {
    try {
      const currentLoad = await si.currentLoad();
      const usage = currentLoad.currentLoad;
      const avg = this.movingAverage.next(usage);

      console.log(`[CPU] Current: ${usage.toFixed(2)}%, Moving Avg: ${avg.toFixed(2)}%`);

      const previousState = this.state;
      let newState: CpuState = 'normal';
      let eventType: AgentEvent['type'] | null = null;
      let severity: AgentEvent['severity'] = 'low';

      if (usage > config.agent.cpuCriticalThreshold) {
        newState = 'critical';
        eventType = 'cpu_usage_critical';
        severity = 'critical';
      } else if (usage > config.agent.cpuWarningThreshold) {
        newState = 'warning';
        eventType = 'cpu_usage_warning';
        severity = 'high';
      }

      // 只有在状态发生变化时才触发事件，避免事件风暴
      if (eventType && newState !== previousState) {
        const event: AgentEvent = {
          id: uuidv4(),
          type: eventType,
          source: SOURCE,
          severity: severity,
          timestamp: new Date(),
          value: usage / 100,
          metadata: {
            average: avg,
            stateChangedFrom: previousState,
            stateChangedTo: newState,
            cores: currentLoad.cpus.map(c => c.load),
          },
        };
        
        this.fireEvent(event);
      }
      
      this.state = newState;

    } catch (error) {
      console.error('❌ Error checking CPU usage:', error);
    }
  };

  private fireEvent(event: AgentEvent) {
    const validationResult = agentEventSchema.safeParse(event);
    if (validationResult.success) {
      console.log(`🔥 Firing event: ${event.type}, usage: ${(event.value! * 100).toFixed(2)}%`);
      eventBus.emit('agent:event', validationResult.data);
    } else {
      console.error('❌ Invalid CPU event schema:', validationResult.error);
    }
  }
}