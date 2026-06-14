import type { BillOfMaterials, CutPlan } from '@/domain/export';
import type { CostResult } from '@/domain/cost';
import type { FurnitureDefinition } from '@/domain/furniture-definition';
import type { InstallationConstraints } from '@/domain/installation-constraints';
import type { ManufacturingModel } from '@/domain/manufacturing';
import type { ResolvedFurniture } from '@/domain/resolved-furniture';
import type { RoomGeometry } from '@/domain/room-geometry';

const API_BASE = import.meta.env.VITE_API_BASE ?? '/api';

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: { 'Content-Type': 'application/json', ...(init?.headers ?? {}) },
    ...init,
  });

  if (!response.ok) {
    const body = (await response.json().catch(() => null)) as { error?: string } | null;
    throw new Error(body?.error ?? `Request failed (${response.status})`);
  }

  return response.json() as Promise<T>;
}

export interface ParseAIRequest {
  description: string;
  name?: string;
}

export interface SolveRequest {
  room: RoomGeometry;
  furniture: FurnitureDefinition;
  installation: InstallationConstraints;
}

export function parseFurnitureDescription(payload: ParseAIRequest) {
  return request<FurnitureDefinition>('/ai/parse', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export function solveFurniture(payload: SolveRequest) {
  return request<ResolvedFurniture>('/solver', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export function compileManufacturing(resolved: ResolvedFurniture) {
  return request<ManufacturingModel>('/manufacturing', {
    method: 'POST',
    body: JSON.stringify({ resolved }),
  });
}

export function calculateCost(model: ManufacturingModel) {
  return request<CostResult>('/cost', {
    method: 'POST',
    body: JSON.stringify({ model }),
  });
}

export interface ExportRequest {
  furnitureName?: string;
  model: ManufacturingModel;
  cost?: CostResult;
}

export function exportBOM(payload: ExportRequest) {
  return request<BillOfMaterials>('/export/bom', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export function exportCutPlan(payload: ExportRequest) {
  return request<CutPlan>('/export/cut-plans', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function downloadPDF(payload: ExportRequest): Promise<void> {
  const response = await fetch(`${API_BASE}/export/pdf`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
  if (!response.ok) {
    const body = (await response.json().catch(() => null)) as { error?: string } | null;
    throw new Error(body?.error ?? `PDF download failed (${response.status})`);
  }
  const blob = await response.blob();
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement('a');
  anchor.href = url;
  anchor.download = `inmobilia-${payload.model.furnitureId}.pdf`;
  anchor.click();
  URL.revokeObjectURL(url);
}

function downloadJSON(filename: string, data: unknown) {
  const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement('a');
  anchor.href = url;
  anchor.download = filename;
  anchor.click();
  URL.revokeObjectURL(url);
}

export async function downloadBOMJSON(payload: ExportRequest) {
  const bom = await exportBOM(payload);
  downloadJSON(`bom-${payload.model.furnitureId}.json`, bom);
}

export async function downloadCutPlanJSON(payload: ExportRequest) {
  const plan = await exportCutPlan(payload);
  downloadJSON(`cut-plans-${payload.model.furnitureId}.json`, plan);
}
