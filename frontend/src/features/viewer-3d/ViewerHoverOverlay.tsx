import { useViewerStore } from '@/store/viewer-store';
import { PANEL_LAYER_LABELS } from './panel-layers';

export function ViewerHoverOverlay() {
  const hoveredPanel = useViewerStore((s) => s.hoveredPanel);
  if (!hoveredPanel) {
    return null;
  }

  return (
    <div className="viewer-hover-overlay" aria-live="polite">
      <strong>{hoveredPanel.label}</strong>
      <span>
        {hoveredPanel.width} × {hoveredPanel.height} × {hoveredPanel.depth} mm
      </span>
      <span>{hoveredPanel.materialLabel}</span>
      <span className="viewer-color-swatch">
        <span style={{ background: hoveredPanel.color }} aria-hidden />
        {PANEL_LAYER_LABELS[hoveredPanel.layer]}
      </span>
    </div>
  );
}
