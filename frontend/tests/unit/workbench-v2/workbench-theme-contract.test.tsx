import { readFileSync } from 'node:fs';
import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { Button } from '@/features/datalink/workbench-v2/components/Button';
import { Input, Select } from '@/features/datalink/workbench-v2/components/inputs';

const source = (path: string) => readFileSync(`${process.cwd()}/src/${path}`, 'utf8');
const css = source('features/datalink/workbench-v2/styles/workbench-v2.css');
describe('Workbench semantic presentation contracts (not pixel acceptance)', () => {
  it('inherits resolved theme instead of forcing dark controls and decorative navy surfaces', () => {
    expect(css).not.toMatch(/color-scheme:\s*dark|#[0-9a-f]{3,8}\b|(?:linear|radial)-gradient/i);
    expect(css).toContain('--wbv2-bg: var(--gw-ground)');
    expect(css).toContain('color-scheme: inherit');
    expect(css).toContain('font-family: var(--gw-font-sans)');
  });
  it('uses legible semantic field boundaries, labels and disabled states', () => {
    expect(css).toContain('border-color: var(--gw-control)');
    expect(css).toContain('color: var(--gw-muted)');
    expect(css).not.toContain('text-transform: uppercase');
    render(<div data-workbench-v2><Input aria-label="Name" /><Select aria-label="Mode"><option>Mode</option></Select><Button disabled>Save</Button></div>);
    expect(screen.getByRole('textbox')).toHaveClass('field');
    expect(screen.getByRole('combobox')).toHaveClass('field');
    expect(screen.getByRole('button', { name: 'Save' })).not.toHaveClass('disabled:opacity-50');
  });
  it('keeps form and confirmation sizing in the existing semantic modal stylesheet', () => {
    expect(css).not.toMatch(/dialog[^}]*75(?:vw|dvh)/s);
    const modal = source('styles/semantic-modal.css');
    expect(modal).toContain('[data-size="form"]');
    expect(modal).toContain('[data-size="confirm"]');
  });
  it('does not require an external font request to render the intranet interface', () => {
    expect(source('index.css')).not.toContain('fonts.googleapis.com');
  });
});
