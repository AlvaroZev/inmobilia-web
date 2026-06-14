import resolvedCloset from '@/domain/fixtures/example-resolved-closet.json';
import { ResolvedFurnitureViewer, ViewerControls } from '@/features/viewer-3d';
import { useViewerStore } from '@/store/viewer-store';
import { normalizeResolvedFurniture } from '@/utils/load-resolved';
import '@/features/viewer-3d/viewer.css';

const demoFurniture = normalizeResolvedFurniture(resolvedCloset);

export function ViewerPage() {
  const overrideFurniture = useViewerStore((s) => s.overrideFurniture);
  const resolvedFurniture = useViewerStore((s) => s.resolvedFurniture);
  const furniture = overrideFurniture ?? resolvedFurniture ?? demoFurniture;
  const source = overrideFurniture ? 'override' : resolvedFurniture ? 'pipeline' : 'demo';

  return (
    <div className="viewer-page">
      <div className="viewer-layout">
        <ViewerControls furnitureName={furniture.name} source={source} />
        <div className="viewer-canvas">
          <ResolvedFurnitureViewer furniture={furniture} />
        </div>
      </div>
    </div>
  );
}
