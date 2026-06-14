import { useViewerStore } from '@/store/viewer-store';
import './viewer.css';

interface ViewerControlsProps {
  furnitureName: string;
  source?: 'demo' | 'pipeline' | 'override';
}

const SOURCE_LABELS: Record<NonNullable<ViewerControlsProps['source']>, string> = {
  override: 'JSON manual (sobreescribe el diseño)',
  pipeline: 'Diseño generado desde el asistente',
  demo: 'Demo (fixture de ejemplo)',
};

export function ViewerControls({ furnitureName, source = 'demo' }: ViewerControlsProps) {
  const selectedVolumeId = useViewerStore((s) => s.selectedVolumeId);
  const showVolumes = useViewerStore((s) => s.showVolumes);
  const showFeatures = useViewerStore((s) => s.showFeatures);
  const showFronts = useViewerStore((s) => s.showFronts);
  const overrideJsonText = useViewerStore((s) => s.overrideJsonText);
  const overrideError = useViewerStore((s) => s.overrideError);
  const toggleVolumes = useViewerStore((s) => s.toggleVolumes);
  const toggleFeatures = useViewerStore((s) => s.toggleFeatures);
  const toggleFronts = useViewerStore((s) => s.toggleFronts);
  const setSelectedVolumeId = useViewerStore((s) => s.setSelectedVolumeId);
  const setOverrideJsonText = useViewerStore((s) => s.setOverrideJsonText);
  const applyJsonOverride = useViewerStore((s) => s.applyJsonOverride);
  const clearJsonOverride = useViewerStore((s) => s.clearJsonOverride);
  const openDrawers = useViewerStore((s) => s.openDrawers);
  const openDrawerCount = Object.values(openDrawers).filter((state) => state === 'half').length;

  return (
    <aside className="viewer-panel">
      <div>
        <p className="eyebrow">Resolved Furniture</p>
        <h2>{furnitureName}</h2>
        <p className="viewer-source">{SOURCE_LABELS[source]}</p>
      </div>

      <div className="viewer-json-override">
        <div className="viewer-json-override-header">
          <span>JSON manual</span>
          <p>Pega un ResolvedFurniture y aplica para sobreescribir el diseño en el visor.</p>
        </div>
        <textarea
          className="viewer-json-input"
          value={overrideJsonText}
          onChange={(e) => setOverrideJsonText(e.target.value)}
          placeholder='{"id":"test","name":"Prueba","root":{...}}'
          spellCheck={false}
        />
        {overrideError && <p className="viewer-json-error">{overrideError}</p>}
        <div className="viewer-json-actions">
          <button type="button" onClick={() => applyJsonOverride()}>
            Aplicar
          </button>
          <button type="button" className="secondary" onClick={() => clearJsonOverride()}>
            Limpiar
          </button>
        </div>
      </div>

      <div className="viewer-toggles">
        <label>
          <input type="checkbox" checked={showVolumes} onChange={toggleVolumes} />
          Debug (volúmenes)
        </label>
        <label>
          <input type="checkbox" checked={showFeatures} onChange={toggleFeatures} />
          Features
        </label>
        <label>
          <input type="checkbox" checked={showFronts} onChange={toggleFronts} />
          Frentes
        </label>
      </div>

      <div className="viewer-selection">
        <span>Cajones</span>
        <strong>{openDrawerCount > 0 ? `${openDrawerCount} semi-abierto(s)` : 'Clic para abrir/cerrar'}</strong>
      </div>

      <div className="viewer-selection">
        <span>Selección</span>
        <strong>{selectedVolumeId ?? 'Ninguna'}</strong>
        {selectedVolumeId && (
          <button type="button" onClick={() => setSelectedVolumeId(null)}>
            Limpiar
          </button>
        )}
      </div>
    </aside>
  );
}
