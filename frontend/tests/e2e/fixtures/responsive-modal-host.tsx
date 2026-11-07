import { useState } from 'react';
import { createRoot } from 'react-dom/client';
import * as Radix from '@radix-ui/react-dialog';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '../../../src/components/ui/dialog';
import ConfirmDialog from '../../../src/components/datalink/ConfirmDialog';
import { SlidePanel } from '../../../src/components/datalink/SlidePanel';
import '../../../src/index.css';
import '../../../src/styles/semantic-modal.css';
// Import consumer overrides last to test the actual cascade, including memory's 85% size.
import '../../../src/features/datalink/workbench-v2/recording-home/recording-home.css';

const kind = new URLSearchParams(location.search).get('kind') ?? 'form';
const longText = Array.from({ length: 240 }, (_, i) => `設定說明 ${i + 1}：保留驗證及操作流程。`).join('\n');
function Content() { return <><label>名稱<input defaultValue="測試資料" /></label><pre style={{ whiteSpace: 'pre-wrap' }}>{longText}</pre><button>最後內容按鈕</button></>; }
function Fixture() {
  const [open, setOpen] = useState(false);
  return <><button onClick={() => setOpen(true)}>開啟測試</button><div role="dialog" data-testid="non-modal" style={{ width: 200, height: 60 }}>非 modal 資訊</div>
    {kind === 'legacy' ? <Dialog open={open} onOpenChange={setOpen}><DialogContent><DialogHeader><DialogTitle>標題</DialogTitle></DialogHeader><div className="gw-dialog-body"><Content /></div><DialogFooter><button onClick={() => setOpen(false)}>完成</button></DialogFooter></DialogContent></Dialog>
    : kind === 'confirm' ? <ConfirmDialog isOpen={open} title="標題" message={longText} onCancel={() => setOpen(false)} onConfirm={() => setOpen(false)} />
    : kind === 'slide' ? <SlidePanel isOpen={open} title="標題" onClose={() => setOpen(false)}><Content /></SlidePanel>
    : <Radix.Root open={open} onOpenChange={setOpen}><Radix.Portal><Radix.Overlay className="gw-modal-overlay" /><Radix.Content className={`gw-modal ${kind === 'memory' ? 'gw-home-source-modal' : ''}`} data-size={kind === 'radix-confirm' ? 'confirm' : 'form'} data-source-view={kind === 'memory' ? 'memory' : undefined} aria-describedby={undefined}>
      <header className="gw-modal-header"><Radix.Title>標題</Radix.Title></header><div className="gw-modal-body"><Content /></div><footer className="gw-modal-footer"><Radix.Close>完成</Radix.Close></footer>
    </Radix.Content></Radix.Portal></Radix.Root>}
  </>;
}
createRoot(document.getElementById('root')!).render(<Fixture />);
