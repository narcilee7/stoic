import EventEmitter from 'eventemitter3';

class TypedEventEmitter<T extends Record<string, any>> {
    private emitter = new EventEmitter();

    on<K extends keyof T>(event: K, listener: (event: T[K]) => void) {
        this.emitter.on(event as string, listener);
    }

    emit<K extends keyof T>(event: K, data: T[K]) {
        this.emitter.emit(event as string, data);
    }

    off<K extends keyof T>(event: K, listener: (event: T[K]) => void) {
        this.emitter.off(event as string, listener);
    }
}

// 定义应用中所有事件及其负载的类型
interface AppEvents {
  'agent:event': import('@/modules/agent/types').AgentEvent;
  'agent:intervention': import('@/modules/agent/types').Intervention;
  'config:reloaded': import('@/config/schema').FullConfig;
}

// 创建并导出一个单例的事件总线实例
export const eventBus = new TypedEventEmitter<AppEvents>();