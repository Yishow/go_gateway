import { expect, test } from '@playwright/test';

for (const viewport of [{ width: 1440, height: 900 }, { width: 390, height: 844 }]) {
  test(`saved MC range loads scoped evidence at ${viewport.width}px`, async ({ page }, testInfo) => {
    const apiRequests: string[] = [];
    await page.setViewportSize(viewport);
    await page.route('**/api/**', async route => {
      const request = route.request(), url = new URL(request.url());
      apiRequests.push(`${request.method()} ${url.pathname}${url.search}`);
      if (request.method() !== 'GET' || url.pathname !== '/api/v1/datalink/source-rules/R/candidates' ||
          url.searchParams.get('view') !== 'tags' || url.searchParams.get('workspace_id') !== 'W' ||
          url.searchParams.get('expected_workspace_revision') !== 'W2' || url.searchParams.get('revision_id') !== 'R2') {
        await route.abort();
        return;
      }
      await new Promise(resolve => setTimeout(resolve, 500));
      await route.fulfill({ json: { success: true, data: { source_rule_id: 'R', workspace_id: 'W', workspace_revision: 'W2', revision_id: 'R2',
        tags: { status: 'ready', candidates: Array.from({ length: 32 }, (_, index) => ({
          id: `C${index}`, identity: { source_rule_id: 'R', candidate_type: 'tags', candidate_kind: 'tag', derived_from_rule_address: 'D0' },
          proposed_signature: `sig${index}`, address: `D${index}`, point_id: `P${index}`, tag_key: `tag${index}`,
          display_name: `TAG_${index}`, data_type: 'int16', status: 'draft',
        })) } } } });
    });
    await page.goto('/tests/e2e/fixtures/recording-home-mc-host.html');
    const dialog = page.getByRole('dialog');
    const prefix = dialog.getByTestId('rule-prefix-input');
    await prefix.fill('MC_NEW_');
    await expect(dialog.getByTestId('rule-mc-persisted-status')).toHaveCount(0);
    await expect(dialog.getByTestId('range-summary')).toContainText('D0');
    await expect(dialog.getByTestId('range-summary')).toContainText('D31');
    await expect(prefix).toHaveValue('MC_NEW_');
    await expect(dialog.getByTestId('home-range-save')).toBeEnabled();
    await page.screenshot({ path: testInfo.outputPath(`home-mc-${viewport.width}.png`) });
    if (viewport.width === 390) {
      const cue = dialog.getByText('向下捲動查看其他設定');
      await expect(cue).toBeVisible();
      const formBox = await dialog.getByTestId('home-saved-range-form').boundingBox();
      const cueBox = await cue.boundingBox();
      expect(formBox).not.toBeNull(); expect(cueBox).not.toBeNull();
      expect(cueBox!.y).toBeGreaterThanOrEqual(formBox!.y + formBox!.height);
      await dialog.getByTestId('home-saved-range-form').evaluate(element => { element.scrollTop = element.scrollHeight; });
      await expect(cue).toHaveCount(0);
      await expect(dialog.getByTestId('share-section')).toBeVisible();
      await expect(dialog.getByTestId('home-range-save')).toBeVisible();
    }
    expect(apiRequests).toHaveLength(1);
  });
}
