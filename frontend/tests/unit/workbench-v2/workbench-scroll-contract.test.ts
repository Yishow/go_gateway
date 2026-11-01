import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';
const source = (path: string) => readFileSync(`${process.cwd()}/src/features/datalink/workbench-v2/${path}`, 'utf8');
describe('document scroll integration (runtime geometry is separate browser evidence)', () => {
  it('does not scroll the non-scrollable step-content or assume a fixed wrapped header height', () => {
    const shell = source('shell/WorkbenchV2Shell.tsx');
    expect(shell).not.toMatch(/getElementById\('step-content'\)\?\.scrollTo/);
    expect(shell).not.toContain('top-[68px]');
    expect(shell).toContain('useWorkbenchScroll');
  });
  it('gives document navigation targets space below the measured sticky header without changing modal body sizing', () => {
    const css = source('styles/workbench-v2.css');
    expect(css).toContain('--wbv2-header-height');
    expect(css).toContain('scroll-margin-block-start');
    expect(css).not.toMatch(/#step-content[^}]*overflow:\s*(auto|scroll)/s);
  });
});
