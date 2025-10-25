import { defineConfig } from 'vitest/config';
import path from 'path';
import tsconfigPathsPlugin from 'vite-tsconfig-paths';

export default defineConfig({
    plugins: [tsconfigPathsPlugin()],
    test: {
        globals: true,
        environment: 'node',
        setupFiles: ['./test/setup.ts'],
    },
    resolve: {
        alias: {
            '@': path.resolve(__dirname, './src'),
        },
    },
});