import { fireEvent, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { mockMappings, mockPoints, mockTags, renderPage, resetTagStepMocks } from './DatalinkWorkbenchTagStep.testHarness';

describe('DatalinkWorkbench tag step', () => {
  beforeEach(() => {
    resetTagStepMocks();
  });

  it('shows an empty state when the selected device has no points', () => {
    mockPoints.splice(0, mockPoints.length);

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByText('workbench.tag.empty.title')).toBeInTheDocument();
    expect(screen.getByTestId('tag-empty-eligible-spans')).toBeInTheDocument();
    expect(screen.getByText('workbench.tag.empty.eligibleSpansNone')).toBeInTheDocument();
  });

  it('shows a go-to-source action in the empty state that navigates to step 2', () => {
    mockPoints.splice(0, mockPoints.length);

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const goToSourceButton = screen.getByTestId('tag-empty-goto-source');
    expect(goToSourceButton).toBeInTheDocument();

    fireEvent.click(goToSourceButton);

    // After clicking, active step should switch to source (step 2)
    const sourceStepButton = screen.getByRole('tab', { name: /workbench.steps.source/ });
    expect(sourceStepButton).toHaveAttribute('aria-selected', 'true');
  });

  it('previews tag keys for selected points and updates the preview when the prefix changes', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByTestId('tag-preview-point-1')).toHaveTextContent('TAG_40001');

    fireEvent.click(screen.getByTestId('tag-template-toggle'));
    fireEvent.change(screen.getByLabelText('workbench.tag.template.prefix'), {
      target: { value: 'linea' },
    });

    expect(screen.getByTestId('tag-preview-point-1')).toHaveTextContent('LINEA_40001');
  });

  it('shows dense source metadata for each candidate row', () => {
    mockPoints[0].last_value = 123;

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByTestId('tag-candidate-board')).toHaveAttribute(
      'data-layout',
      'two-line-board',
    );
    expect(screen.getByTestId('tag-candidate-point-1')).toHaveAttribute('data-layout', 'two-line');
    expect(screen.getByTestId('tag-candidate-point-1')).toHaveTextContent('40001');
    expect(screen.getByTestId('tag-raw-point-1')).toHaveTextContent('123');
    expect(screen.getByTestId('tag-status-point-1')).toHaveTextContent(
      'workbench.tag.board.status.unbound',
    );
  });

  it('shows review-first summary metrics for generated and pending rows', () => {
    mockTags.push({
      id: 'tag-bound',
      key: 'TAG_40001',
      display_name: 'Flow Tag',
      description: '',
      data_type: 'int16',
      unit: '',
      labels: null,
      status: 'active',
      created_at: '',
      updated_at: '',
    });
    mockMappings.push({
      id: 'mapping-bound',
      point_id: 'point-1',
      tag_id: 'tag-bound',
      enabled: true,
      transform_pipeline: '',
      created_at: '',
      updated_at: '',
    });

    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByTestId('tag-review-summary')).toBeInTheDocument();
    expect(screen.getByTestId('tag-review-generated')).toHaveTextContent('1');
    expect(screen.getByTestId('tag-review-needs-review')).toHaveTextContent('1');
    expect(screen.getByTestId('tag-review-selected-exceptions')).toHaveTextContent('0');
  });

  it('shows adaptive board guidance and flow hints as selection state changes', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByText('workbench.tag.board.selectionHint.none')).toBeInTheDocument();
    expect(screen.getByText('workbench.tag.board.flowModeHint.create')).toBeInTheDocument();

    fireEvent.click(screen.getByLabelText('Flow Sensor'));

    expect(screen.getByText('workbench.tag.board.selectionHint.ready')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.tag.board.flow.existing' }));

    expect(screen.getByText('workbench.tag.board.flowModeHint.existing')).toBeInTheDocument();
  });

  it('groups manual correction tools under an exception-handling section', () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByTestId('tag-exception-tools')).toBeInTheDocument();
    expect(screen.getByText('workbench.tag.exception.title')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'workbench.tag.board.flow.create' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'workbench.tag.board.flow.existing' })).toBeInTheDocument();
  });
});
