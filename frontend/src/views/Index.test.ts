import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import Index from './Index.vue';

describe('Index', () => {
  it('should render the Index view', () => {
    const wrapper = mount(Index);
    expect(wrapper.exists()).toBe(true);
  });

  it('should display a welcome message', () => {
    const wrapper = mount(Index);
    const text = wrapper.text();
    expect(text).toContain('Gestione Caselo');
  });
});
