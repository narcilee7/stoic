export interface MovingAverage {
    next: (value: number) => number;
    getAverage: () => number;
    getValues: () => number[];
}

export function movingAverage(size: number): MovingAverage {
    const values: number[] = [];
    let sum = 0;
    
    const next = (v: number): number => {
        values.push(v);
        sum += v;
        if (values.length > size) {
            sum -= values.shift()!;
        }
        return sum / values.length;
    }

    const getAverage = (): number => {
        return sum / values.length;
    }

    const getValues = (): number[] => {
        return values;
    }

    return {
        next,
        getAverage,
        getValues,
    }
}