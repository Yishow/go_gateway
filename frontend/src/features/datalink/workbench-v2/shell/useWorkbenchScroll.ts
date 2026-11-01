import { useLayoutEffect, useRef, type RefObject } from 'react';

const GAP = 8;

/** The V2 workspace uses the document scroller; dialogs retain their own body scroller. */
export function installWorkbenchScroll(root: HTMLElement): () => void {
  const document = root.ownerDocument;
  const window = document.defaultView;
  if (!window) return () => undefined;
  const header = root.querySelector<HTMLElement>('.wbv2-topbar');
  const style = document.documentElement.style;
  const previousPadding = style.getPropertyValue('scroll-padding-top');
  const previousPriority = style.getPropertyPriority('scroll-padding-top');
  const previousHeight = root.style.getPropertyValue('--wbv2-header-height');
  const update = () => {
    const height = Math.ceil(header?.getBoundingClientRect().height ?? 0);
    root.style.setProperty('--wbv2-header-height', `${height}px`);
    style.setProperty('scroll-padding-top', `${height + GAP}px`);
  };
  update();
  const observer = typeof ResizeObserver === 'undefined' ? undefined : new ResizeObserver(update);
  if (header) observer?.observe(header);
  window.addEventListener('resize', update);
  return () => {
    observer?.disconnect();
    window.removeEventListener('resize', update);
    if (previousPadding) style.setProperty('scroll-padding-top', previousPadding, previousPriority);
    else style.removeProperty('scroll-padding-top');
    if (previousHeight) root.style.setProperty('--wbv2-header-height', previousHeight);
    else root.style.removeProperty('--wbv2-header-height');
  };
}

/** Run after a committed navigation, never on pointer-down or focus events. */
export function scrollWorkbenchContent(root: HTMLElement): void {
  const document = root.ownerDocument;
  const content = root.querySelector<HTMLElement>('#step-content');
  const header = root.querySelector<HTMLElement>('.wbv2-topbar');
  if (!content) return;
  const scrollTop = document.scrollingElement?.scrollTop ?? document.defaultView?.scrollY ?? 0;
  const top = Math.max(0, scrollTop + content.getBoundingClientRect().top - (header?.getBoundingClientRect().height ?? 0) - GAP);
  document.defaultView?.scrollTo({ top, behavior: 'auto' });
}

export function useWorkbenchScroll(root: RefObject<HTMLDivElement | null>, navigation: string): void {
  const previous = useRef(navigation);
  useLayoutEffect(() => root.current ? installWorkbenchScroll(root.current) : undefined, [root]);
  useLayoutEffect(() => {
    if (previous.current === navigation) return;
    previous.current = navigation;
    if (root.current) scrollWorkbenchContent(root.current);
  }, [root, navigation]);
}
