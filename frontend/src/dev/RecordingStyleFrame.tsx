/** Development route only; the document isolates preview theme/portals from application owners. */
export default function RecordingStyleFrame() {
  return <iframe title="Recording style preview — simulated" src="/src/dev/recording-style.html"
    style={{ width: '100%', height: '100dvh', border: 0 }} />;
}
