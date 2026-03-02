import { Link } from 'react-router-dom';

function EntryCard({
  title,
  description,
  to,
}: {
  title: string;
  description: string;
  to: string;
}) {
  return (
    <Link
      to={to}
      className="block rounded-xl border border-slate-700 bg-slate-900/80 p-6 transition hover:border-blue-400 hover:bg-slate-900"
    >
      <h2 className="text-lg font-semibold text-slate-100">{title}</h2>
      <p className="mt-2 text-sm text-slate-300">{description}</p>
    </Link>
  );
}

export default function GatewayEntryPage() {
  return (
    <main className="mx-auto max-w-4xl px-6 py-10 text-slate-100">
      <h1 className="text-2xl font-bold">Gateway 雙入口</h1>
      <p className="mt-2 text-sm text-slate-300">
        請選擇你的工作模式：快速設定（新手友善）或專家工作台（進階控制）。
      </p>

      <section className="mt-6 grid gap-4 md:grid-cols-2">
        <EntryCard
          title="Quick Setup"
          description="適合常見代理場景，用引導式流程快速完成設定。"
          to="/gateway/quick-setup"
        />
        <EntryCard
          title="Expert Workbench"
          description="適合進階需求，可配置高密度路由與進階欄位。"
          to="/gateway/expert-workbench"
        />
      </section>
    </main>
  );
}
