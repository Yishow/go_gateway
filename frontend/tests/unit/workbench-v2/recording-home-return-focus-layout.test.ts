import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { watchHomeReturnFocus } from '@/features/datalink/workbench-v2/recording-home/recordingHomeReturnFocus';
const observers: { callback: () => void; observe: ReturnType<typeof vi.fn>; disconnect: ReturnType<typeof vi.fn>; options: { threshold: number } }[] = [];
beforeEach(() => {
  observers.length = 0;
  vi.stubGlobal('innerHeight', 720);
  // The observed regression changes position without resizing the region or document.
  vi.stubGlobal('ResizeObserver', class { observe = vi.fn(); disconnect = vi.fn(); unobserve = vi.fn(); });
  vi.stubGlobal('IntersectionObserver', class {
    constructor(callback: () => void, options: { threshold: number }) { observers.push({ callback, options, observe: this.observe, disconnect: this.disconnect }); }
    observe = vi.fn(); disconnect = vi.fn(); unobserve = vi.fn();
  });
});
afterEach(() => { document.body.replaceChildren(); vi.restoreAllMocks(); vi.unstubAllGlobals(); });
function start() {
  const region = document.createElement('main'), target = document.createElement('button');
  region.append(target); document.body.append(region); target.focus();
  let documentTop = 735, documentScroll = 0;
  vi.spyOn(target, 'getBoundingClientRect').mockImplementation(() => ({ top: documentTop - documentScroll,
    bottom: documentTop + 44 - documentScroll, height: 44, width: 80, left: 0, right: 80, x: 0, y: 0, toJSON: () => ({}) }));
  const scroll = vi.spyOn(window, 'scrollBy').mockImplementation((...args: unknown[]) => { const options = args[0] as { top?: number }; documentScroll += options.top ?? 0; });
  const ticket = watchHomeReturnFocus(target);
  expect(scroll).toHaveBeenCalledWith({ top: 67, behavior: 'instant' });
  const shift = () => { documentTop += 48; observers[0]?.callback(); };
  return { ticket, target, region, scroll, shift };
}
describe('bounded return-focus layout correction', () => {
  it('corrects a 48px target translation without any region or document resize', () => {
    const { ticket, scroll, shift, target } = start();
    shift();
    expect(scroll).toHaveBeenLastCalledWith({ top: 48, behavior: 'instant' });
    expect(target.getBoundingClientRect().bottom).toBe(712);
    expect(observers).toHaveLength(1);
    expect(observers[0].observe.mock.calls).toEqual([[target]]);
    expect(observers[0].options).toEqual({ threshold: 1 });
    ticket.cancel();
  });
  it.each(['wheel', 'pointerdown', 'keydown'])('stops permanently on %s without pulling the user back', type => {
    const { ticket, scroll, shift } = start();
    document.dispatchEvent(new Event(type, { bubbles: true })); shift(); ticket.finish();
    expect(scroll).toHaveBeenCalledTimes(1);
    expect(observers[0]?.disconnect).toHaveBeenCalledOnce();
  });
  it('cancels when focus leaves or the owner is disposed', () => {
    const { ticket, scroll, shift } = start();
    const other = document.createElement('button'); document.body.append(other); other.focus();
    shift(); ticket.cancel(); expect(scroll).toHaveBeenCalledTimes(1);
    expect(observers[0]?.disconnect).toHaveBeenCalledOnce();
  });
  it('finishes with the last layout and ignores later observer deliveries', () => {
    const { ticket, scroll, shift } = start();
    shift(); ticket.finish(); shift();
    expect(scroll).toHaveBeenCalledTimes(2);
    expect(observers[0]?.disconnect).toHaveBeenCalledOnce();
  });
});
