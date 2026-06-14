import { useState } from 'react';
import type { CostResult } from '@/domain/cost';
import type { ManufacturingModel } from '@/domain/manufacturing';
import { downloadBOMJSON, downloadCutPlanJSON, downloadPDF } from '@/services/api';
import './export-panel.css';

interface ExportPanelProps {
  furnitureName: string;
  model: ManufacturingModel;
  cost?: CostResult;
}

export function ExportPanel({ furnitureName, model, cost }: ExportPanelProps) {
  const [loading, setLoading] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const payload = { furnitureName, model, cost };

  const run = async (label: string, action: () => Promise<void>) => {
    setLoading(label);
    setError(null);
    try {
      await action();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al exportar');
    } finally {
      setLoading(null);
    }
  };

  return (
    <div className="export-panel">
      <h3>Exportar</h3>
      <div className="export-actions">
        <button
          type="button"
          disabled={Boolean(loading)}
          onClick={() => run('pdf', () => downloadPDF(payload))}
        >
          {loading === 'pdf' ? 'Generando...' : 'PDF (BOM + cortes)'}
        </button>
        <button
          type="button"
          disabled={Boolean(loading)}
          onClick={() => run('bom', () => downloadBOMJSON(payload))}
        >
          BOM JSON
        </button>
        <button
          type="button"
          disabled={Boolean(loading)}
          onClick={() => run('cuts', () => downloadCutPlanJSON(payload))}
        >
          Planos de corte JSON
        </button>
      </div>
      {error && <p className="export-error">{error}</p>}
    </div>
  );
}
