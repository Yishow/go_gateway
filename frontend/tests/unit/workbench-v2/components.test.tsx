import * as React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import {
  Icon,
  Button,
  Input,
  Select,
  Textarea,
  Field,
  Toggle,
  StatusChip,
  SectionCard,
} from '../../../src/features/datalink/workbench-v2/components';

describe('Icon Component', () => {
  it('should render SVG element with aria-hidden="true"', () => {
    const { container } = render(<Icon name="device" />);
    const svg = container.querySelector('svg');
    expect(svg).toBeInTheDocument();
    expect(svg).toHaveAttribute('aria-hidden', 'true');
  });

  it('should render specific paths for existing icon names', () => {
    const { container } = render(<Icon name="check" />);
    const path = container.querySelector('path');
    expect(path).toBeInTheDocument();
    expect(path).toHaveAttribute('d', 'M4 12l5 5L20 6');
  });
});

describe('Button Component', () => {
  it('should render children and trigger onClick on click', () => {
    const handleClick = vi.fn();
    render(<Button onClick={handleClick}>Click Me</Button>);
    const button = screen.getByRole('button', { name: /click me/i });
    expect(button).toBeInTheDocument();
    fireEvent.click(button);
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('should not trigger onClick when disabled', () => {
    const handleClick = vi.fn();
    render(<Button disabled onClick={handleClick}>Click Me</Button>);
    const button = screen.getByRole('button', { name: /click me/i });
    expect(button).toBeDisabled();
    fireEvent.click(button);
    expect(handleClick).not.toHaveBeenCalled();
  });
});

describe('Inputs Components', () => {
  it('should forward ref for Input', () => {
    const ref = React.createRef<HTMLInputElement>();
    render(<Input ref={ref} defaultValue="test" />);
    expect(ref.current).toBeInstanceOf(HTMLInputElement);
    expect(ref.current?.value).toBe('test');
  });

  it('should forward ref for Select', () => {
    const ref = React.createRef<HTMLSelectElement>();
    render(
      <Select ref={ref} defaultValue="opt1">
        <option value="opt1">Option 1</option>
      </Select>
    );
    expect(ref.current).toBeInstanceOf(HTMLSelectElement);
    expect(ref.current?.value).toBe('opt1');
  });

  it('should forward ref for Textarea', () => {
    const ref = React.createRef<HTMLTextAreaElement>();
    render(<Textarea ref={ref} defaultValue="textarea-test" />);
    expect(ref.current).toBeInstanceOf(HTMLTextAreaElement);
    expect(ref.current?.value).toBe('textarea-test');
  });
});

describe('Field Component', () => {
  it('should show label and asterisk when required', () => {
    render(
      <Field label="Name" required>
        <Input />
      </Field>
    );
    const labelSpan = screen.getByText('Name');
    expect(labelSpan).toHaveClass('after:content-["*"]');
  });

  it('should show error and hide hint when error is present', () => {
    render(
      <Field label="Name" hint="Enter name" error="Name is required">
        <Input />
      </Field>
    );
    expect(screen.getByText('Name is required')).toBeInTheDocument();
    expect(screen.queryByText('Enter name')).not.toBeInTheDocument();
  });

  it('should show hint when error is not present', () => {
    render(
      <Field label="Name" hint="Enter name">
        <Input />
      </Field>
    );
    expect(screen.getByText('Enter name')).toBeInTheDocument();
  });
});

describe('Toggle Component', () => {
  it('should render with switch role and correct aria-checked state', () => {
    const handleChange = vi.fn();
    render(<Toggle checked={false} onChange={handleChange} label="Active" />);
    const toggle = screen.getByRole('switch', { name: /active/i });
    expect(toggle).toBeInTheDocument();
    expect(toggle).toHaveAttribute('aria-checked', 'false');

    fireEvent.click(toggle);
    expect(handleChange).toHaveBeenCalledWith(true);
  });
});

describe('StatusChip Component', () => {
  it('should render correct tone classes and text', () => {
    render(<StatusChip tone="success">Active</StatusChip>);
    const chip = screen.getByText('Active');
    expect(chip).toHaveClass('border-emerald-500/30');
  });

  it('should render pulse-dot for info and success tones', () => {
    const { container: successContainer } = render(
      <StatusChip tone="success">Success</StatusChip>
    );
    const successDot = successContainer.querySelector('.pulse-dot');
    expect(successDot).toBeInTheDocument();

    const { container: neutralContainer } = render(
      <StatusChip tone="neutral">Neutral</StatusChip>
    );
    const neutralDot = neutralContainer.querySelector('.pulse-dot');
    expect(neutralDot).not.toBeInTheDocument();
  });
});

describe('SectionCard Component', () => {
  it('should render header only if title or aside is provided', () => {
    const { container: emptyContainer } = render(
      <SectionCard>Content</SectionCard>
    );
    expect(emptyContainer.querySelector('header')).not.toBeInTheDocument();

    const { container: withTitleContainer } = render(
      <SectionCard title="My Title">Content</SectionCard>
    );
    expect(withTitleContainer.querySelector('header')).toBeInTheDocument();
    expect(screen.getByText('My Title')).toBeInTheDocument();
  });
});
