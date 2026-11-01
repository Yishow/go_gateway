import { SaveBar } from '@/features/datalink/workbench-v2/settings/SaveBar';
import { useRef } from 'react';
import { cleanup, fireEvent, render } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { installWorkbenchScroll, scrollWorkbenchContent, useWorkbenchScroll } from '@/features/datalink/workbench-v2/shell/useWorkbenchScroll';

const rectangle = (top: number, height: number) => ({ top, height, bottom: top + height, left: 0, right: 390, width: 390, x: 0, y: top, toJSON: () => ({}) });
afterEach(() => { cleanup(); vi.restoreAllMocks(); vi.unstubAllGlobals(); document.documentElement.style.removeProperty('scroll-padding-top'); document.documentElement.style.removeProperty('scroll-padding-bottom'); Reflect.deleteProperty(document, 'scrollingElement'); });
function fixture() {
  const root = document.createElement('div');
  root.innerHTML = '<header class="wbv2-topbar"></header><main id="step-content"><button>Target</button></main>';
  document.body.append(root);
  return root;
}
describe('measured document scroll adapter (mocked geometry, not browser acceptance)', () => {
  it('tracks wrapped header height, preserves existing document padding and disconnects on unmount', () => {
    const root = fixture();
    const header = root.querySelector('header')!;
    let height = 61;
    vi.spyOn(header, 'getBoundingClientRect').mockImplementation(() => rectangle(0, height));
    let resize: () => void = () => undefined;
    const disconnect = vi.fn();
    vi.stubGlobal('ResizeObserver', class { constructor(callback: () => void) { resize = callback; } observe = vi.fn(); disconnect = disconnect; });
    document.documentElement.style.setProperty('scroll-padding-top', '12px', 'important');
    const dispose = installWorkbenchScroll(root);
    expect(document.documentElement.style.scrollPaddingTop).toBe('69px');
    expect(root.style.getPropertyValue('--wbv2-header-height')).toBe('61px');
    height = 143; resize();
    expect(document.documentElement.style.scrollPaddingTop).toBe('151px');
    expect(root.style.getPropertyValue('--wbv2-header-height')).toBe('143px');
    dispose();
    expect(disconnect).toHaveBeenCalledOnce();
    expect(document.documentElement.style.scrollPaddingTop).toBe('12px');
    expect(document.documentElement.style.getPropertyPriority('scroll-padding-top')).toBe('important');
    expect(root.style.getPropertyValue('--wbv2-header-height')).toBe('');
    root.remove();
  });
  it('scrolls the document, not the unbounded main, using current header geometry', () => {
    const root = fixture();
    vi.spyOn(root.querySelector('header')!, 'getBoundingClientRect').mockReturnValue(rectangle(0, 143));
    vi.spyOn(root.querySelector('main')!, 'getBoundingClientRect').mockReturnValue(rectangle(-500, 5035));
    Object.defineProperty(document, 'scrollingElement', { configurable: true, value: document.documentElement });
    document.documentElement.scrollTop = 1000;
    const scroll = vi.spyOn(window, 'scrollTo').mockImplementation(() => undefined);
    scrollWorkbenchContent(root);
    expect(scroll).toHaveBeenCalledExactlyOnceWith({ top: 349, behavior: 'auto' });
    root.remove();
  });
  it('does not move pointer targets or restore stale scrolling on rerender; scrolls only committed navigation', () => {
    const scroll = vi.spyOn(window, 'scrollTo').mockImplementation(() => undefined);
    function Page({ navigation }: { navigation: string }) {
      const root = useRef<HTMLDivElement>(null);
      useWorkbenchScroll(root, navigation);
      return <div ref={root}><header className="wbv2-topbar" /><main id="step-content"><button>Target</button></main></div>;
    }
    const { rerender, getByRole } = render(<Page navigation="flow:4" />);
    fireEvent.pointerDown(getByRole('button')); getByRole('button').focus(); fireEvent.pointerUp(getByRole('button'));
    rerender(<Page navigation="flow:4" />);
    expect(scroll).not.toHaveBeenCalled();
    rerender(<Page navigation="settings:4" />);
    expect(scroll).toHaveBeenCalledOnce();
  });
});
it('reserves measured Settings footer space without moving focus and restores document bottom padding on exit', () => {
  let height = 144; let bottom = '16px'; let resize: () => void = () => undefined;
  const disconnect = vi.fn();
  vi.spyOn(HTMLElement.prototype, 'getBoundingClientRect').mockImplementation(() => rectangle(600, height));
  vi.spyOn(window, 'getComputedStyle').mockImplementation(() => ({ bottom }) as CSSStyleDeclaration);
  vi.stubGlobal('ResizeObserver', class { constructor(callback: () => void) { resize = callback; } observe = vi.fn(); disconnect = disconnect; });
  const scroll = vi.spyOn(window, 'scrollTo').mockImplementation(() => undefined);
  document.documentElement.style.setProperty('scroll-padding-bottom', '12px', 'important');
  document.documentElement.style.scrollPaddingTop = '169px';
  const view = render(<SaveBar onReset={vi.fn()} onSave={vi.fn()} />);
  expect(document.documentElement.style.scrollPaddingBottom).toBe('168px');
  const button = view.getAllByRole('button')[0];
  fireEvent.pointerDown(button); button.focus(); fireEvent.pointerUp(button);
  expect(document.activeElement).toBe(button); expect(scroll).not.toHaveBeenCalled();
  height = 100; resize(); expect(document.documentElement.style.scrollPaddingBottom).toBe('124px');
  bottom = '24px'; fireEvent(window, new Event('resize')); expect(document.documentElement.style.scrollPaddingBottom).toBe('132px');
  expect(document.documentElement.style.scrollPaddingTop).toBe('169px');
  bottom = 'auto'; resize(); expect(document.documentElement.style.scrollPaddingBottom).toBe('108px');
  height = -1; resize(); expect(document.documentElement.style.scrollPaddingBottom).toBe('8px');
  view.unmount(); expect(disconnect).toHaveBeenCalledOnce();
  expect(document.documentElement.style.scrollPaddingBottom).toBe('12px');
  expect(document.documentElement.style.getPropertyPriority('scroll-padding-bottom')).toBe('important');
  height = 222; fireEvent(window, new Event('resize'));
  expect(document.documentElement.style.scrollPaddingBottom).toBe('12px');
});
