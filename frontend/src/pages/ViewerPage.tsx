import { useMemo, useState } from 'react';
import resolvedCloset from '@/domain/fixtures/example-resolved-closet.json';
import resolvedSingleDrawer from '@/domain/fixtures/example-resolved-single-drawer.json';
import { ResolvedFurnitureViewer, ViewerControls, ViewerHoverOverlay } from '@/features/viewer-3d';
import { useViewerStore } from '@/store/viewer-store';
import { normalizeResolvedFurniture } from '@/utils/load-resolved';
import '@/features/viewer-3d/viewer.css';

const VIEWER_DEMO_CASES = {
  closet: { label: 'Caso 00 — Closet', data: resolvedCloset },
  'single-drawer': { label: 'Caso 01 — Cajonera 50×35×40', data: resolvedSingleDrawer },
} as const;

type DemoCaseId = keyof typeof VIEWER_DEMO_CASES;

export function ViewerPage() {
  const [demoCaseId, setDemoCaseId] = useState<DemoCaseId>('single-drawer');

  const overrideFurniture = useViewerStore((s) => s.overrideFurniture);
  const resolvedFurniture = useViewerStore((s) => s.resolvedFurniture);
  const demoFurniture = useMemo(
    () => normalizeResolvedFurniture(VIEWER_DEMO_CASES[demoCaseId].data),
    [demoCaseId],
  );
  const furniture = overrideFurniture ?? resolvedFurniture ?? demoFurniture;
  const source = overrideFurniture ? 'override' : resolvedFurniture ? 'pipeline' : 'demo';

  return (
    <div className="viewer-page">
      <div className="viewer-layout">
        <ViewerControls
          furnitureName={furniture.name}
          source={source}
          demoCaseId={source === 'demo' ? demoCaseId : undefined}
          demoCases={VIEWER_DEMO_CASES}
          onDemoCaseChange={(id) => setDemoCaseId(id as DemoCaseId)}
        />
        <div className="viewer-canvas">
          <ResolvedFurnitureViewer furniture={furniture} />
          <ViewerHoverOverlay />
        </div>
      </div>
    </div>
  );
}
