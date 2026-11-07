> 歷史報告閱覽副本：正規化本機絕對路徑與附件連結；保留當時結論與未驗項目。原始位元組及完整附件見 `local-only/original-evidence.tar.gz`，不要將歷史 PASS 當成新版重驗。

# Visual comparison recheck, 2026-11-05

Baseline product: 1dea0af76b8f3e10dfbf4931aa1dfce5adb4fbfe. Docs-only HEAD: 2d437234a9ea76a20f254805cc8bdc486ccc31c6.

Method: Opened all five actual original generated images and all ten final 1440x1000 production zh-TW light/dark screenshots with the image viewer. This is a manual layout/style comparison, not a pixel-diff test. Each generated image contains both theme concepts side by side; each real screenshot contains one theme in the application modal. All files remained private and local. Original images have no mobile/tablet reference variants, so exact visual fidelity at those widths cannot be claimed.

Source pairing:
- grid: references/exec-51e970e9-3551-4cc5-9bf1-465239b04d0d.png vs live-hsl/production-zh-TW/actual-1440-{light,dark}-grid.png
- navigator: references/exec-86b48a15-0a46-4398-83f2-6b5e8d83da94.png vs live-hsl/production-zh-TW/actual-1440-{light,dark}-navigator.png
- map: references/exec-c91df0e8-3161-457e-9e8c-7c6146ddbf4b.png vs live-hsl/production-zh-TW/actual-1440-{light,dark}-map.png
- chapters: references/exec-943be82c-cc00-40cf-8873-bf81d82405f2.png vs live-hsl/production-zh-TW/actual-1440-{light,dark}-chapters.png
- detail: references/exec-201f2c58-3543-4c36-b507-c703f07757d4.png vs live-hsl/production-zh-TW/actual-1440-{light,dark}-detail.png

Findings:
- Shared: The implementation preserves the application theme rather than reproducing the reference's blue/teal type gradients, saturated teal primary actions, bolder address labels, larger value typography and soft gradients/shadows. Current cells and labels are muted gray/green and small monospaced text, with outlined buttons. The extra five-layout selector and help text take vertical space not present in the individual generated concepts. Current 1000x750 modal at 1440x1000 exposes fewer initial rows than the concept. These are visual differences, not failures of the separately verified interaction tests.
- Grid: The address grid and multi-register joined cells match the concept's structure; visual styling/density do not match exactly.
- Navigator: The left-side range navigation matches the concept, but current cards expand per-range metadata, while the concept uses a compact range list. The concept places locate and paging controls in the left sidebar; current controls sit across the shared top toolbar. Current eight grid columns differ from concept ten.
- Map: A horizontal occupancy map above the grid is present. Current rendering uses a tall multi-line map section with slider and endpoint labels, rather than the concept's compact one-line colored strip. The concept's selected-address summary is a compact bottom bar, not reproduced in the same form.
- Chapters: Most material layout difference: the concept has a full-width chapter disclosure section above a full-width grid. Current desktop implementation puts chapter disclosures to the left of the grid to preserve visible rows. This is a deliberate density adjustment, not faithful reproduction of the original arrangement.
- Detail: Left grid/right detail structure is present. Concept has five grid columns, current desktop has eight. Detail typography and field presentation differ, including long ISO timestamps in the current panel. The concept's compact range-navigation cards and current panel arrangement are not identical.

Data differences excluded from fidelity assessment: real isolated HSL data has six sparse configured ranges, additional numeric types and real readings; concept has two dense ranges and mock zeros/dashes. Those differences do not alone establish styling defects.

Conclusion: Original images were used and manually compared; the five functional layouts and interaction behavior are validated. The current product is NOT a pixel-exact or visually faithful reproduction. Existing passes must not be described as visual equivalence. The requested modal sizing adjustment can improve available space, but does not itself fix these style and arrangement gaps.
