# Recording style preview (development only)

Start the existing Vite dev server (`cd frontend && npm run dev`) and open
`/__dev/recording-style`. There is deliberately no operator navigation item.
The route loads `recording-style.html` in its own iframe document. That document
creates its own React root, so its theme classes and Radix portals belong to the
preview document, not the application's theme provider or portal root.

The preview reuses `RecordingHome`, `SavedRecordingRemoveDialog`, shared
`Button`/`Alert`, and the existing semantic CSS/tokens. Labels explicitly identify
simulated props; no production bootstrap or resource identity is created.
Navigation/retry/confirmation callbacks only update a local simulation notice.
There are no service/query/autosave/command owners, credentials or persistence.
Its disposable i18n instance loads only shared JSON resources, with no detector,
cache, backend or production i18n/config import. Bootstrap is included in the
storage/network/command-spy test, not merely component interactions.
Light/dark switching changes only this disposable document's root class.

Use Tab/Shift+Tab to inspect focus, select guidance states, edit the sample input,
and open/cancel/confirm the actual remove dialog with simulated callbacks. Resize
the browser to inspect mobile/short viewport behavior. Component tests do not
prove browser geometry, visual contrast or native iframe focus behavior.

Related automated checks (from `frontend/`):

```sh
node_modules/.bin/vitest run tests/unit/dev
```

The production test executes a real Vite production build in a temporary output
folder and rejects the preview entry, route, fixture marker and chunks. It sets
`NODE_ENV=production` explicitly because Vitest inherits `NODE_ENV=test`.
This folder is not under `public`; none of its HTML is a production build input.
Browser/visual verification remains a separate gate, not implied by these tests.
