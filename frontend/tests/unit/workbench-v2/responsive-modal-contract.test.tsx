import { readFileSync } from 'node:fs';
import { cleanup, render, screen } from '@testing-library/react';
import { afterEach, expect, it } from 'vitest';
import { Dialog, DialogContent, DialogHeader, DialogFooter } from '@/components/ui/dialog';
import ConfirmDialog from '@/components/datalink/ConfirmDialog';
import { SlidePanel } from '@/components/datalink/SlidePanel';

afterEach(cleanup);
it('legacy primitives explicitly opt in while unrelated popovers do not', () => {
  const { unmount } = render(<Dialog open><DialogContent><DialogHeader>Header</DialogHeader><div>Body</div><DialogFooter>Footer</DialogFooter></DialogContent></Dialog>);
  expect(screen.getByRole('dialog')).toHaveClass('gw-legacy-dialog');
  unmount();
  const confirm = render(<ConfirmDialog isOpen title="Confirm" message="Message" onConfirm={() => {}} onCancel={() => {}} />);
  expect(screen.getByRole('dialog')).toHaveClass('gw-legacy-confirm');
  expect(screen.getByText('Message')).toHaveAttribute('tabindex', '0');
  confirm.unmount();
  render(<SlidePanel isOpen title="Panel" onClose={() => {}}>Body</SlidePanel>);
  expect(screen.getByRole('dialog')).toHaveClass('gw-slide-dialog');
});
it('fullscreen contract covers tablets, short viewports, safe areas and classified confirmations', () => {
  const css = readFileSync('src/styles/responsive-modal.css', 'utf8');
  expect(css).toContain('(max-width: 1024px)');
  expect(css).toContain('(pointer: coarse) and (max-width: 1366px)');
  expect(css).toContain('(max-height: 479px)');
  expect(css).toContain('.gw-modal.gw-modal.gw-modal[data-size]');
  expect(css).toContain('height: 100dvh');
  expect(css).toContain('--tw-translate-x: 0; --tw-translate-y: 0;');
  expect(css).toContain('font-size: max(16px, 1rem)');
  expect(css).toContain('input:not([type=checkbox], [type=radio], [type=range])');
  expect(css).toContain('safe-area-inset-bottom');
  expect(css).not.toContain('[role=');
  expect(readFileSync('src/index.css', 'utf8')).toContain("@import './styles/responsive-modal.css'");
});
