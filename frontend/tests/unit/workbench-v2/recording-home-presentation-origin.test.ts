import { expect, it } from 'vitest';
import axios from 'axios';
import { recordingHomeSourceContext } from '@/hooks/datalink/useRecordingHomeSources';
import { readHomePresentation, rememberHomePresentation } from '@/features/datalink/workbench-v2/recording-home/recordingHomePresentation';
it('late old-origin presentation cleanup cannot save into the newly selected API context', () => {
  const client = {}, backend = axios.create({ baseURL: 'https://same/api/old' });
  const old = recordingHomeSourceContext({ backend });
  const presentation = readHomePresentation(client, 'W', '', old);
  backend.defaults.baseURL = 'https://same/api/new';
  const current = recordingHomeSourceContext({ backend });
  rememberHomePresentation(client, 'W', { ...presentation, state: { ...presentation.state, keyword: 'old origin' }, scrollTop: 200 }, old);
  expect(readHomePresentation(client, 'W', '', current).state.keyword).toBe('');
  expect(readHomePresentation(client, 'W', '', old).state.keyword).toBe('old origin');
});
