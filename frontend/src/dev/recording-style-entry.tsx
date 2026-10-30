import { I18nextProvider } from 'react-i18next';
// Not a production entry. No application, query, settings or command owner is mounted.
export const previewReady = import.meta.env.DEV
  ? Promise.all([import('react-dom/client'), import('./RecordingStylePreview'), import('./previewI18n'),
    import('../index.css'), import('@fontsource/inter/400.css'), import('@fontsource/inter/500.css'), import('@fontsource/inter/600.css')])
    .then(async ([{ createRoot }, { default: Preview }, { previewI18n, previewTranslationsReady }]) => {
      await previewTranslationsReady;
      const node = document.getElementById('root');
      if (!node) return;
      const root = createRoot(node);
      root.render(<I18nextProvider i18n={previewI18n}><Preview /></I18nextProvider>);
      return root;
    }).catch(() => {
      const root = document.getElementById('root');
      if (root) root.textContent = 'Preview could not load. Check the development server and reload; no data was written.';
      return undefined;
    })
  : Promise.resolve(undefined);
