import { describe, it, expect } from 'vitest';
import { createMemoryHistory } from 'vue-router';
import { createRouter } from './index';

describe('Router', () => {
  it('should have a root route that maps to Index view', async () => {
    const router = createRouter(createMemoryHistory());
    await router.push('/');
    await router.isReady();

    expect(router.currentRoute.value.name).toBe('index');
  });

  it('should export a router instance with browser history for main app', () => {
    const { router } = require('./index');
    expect(router).toBeDefined();
    expect(router.options.history).toBeDefined();
  });
});
