/**
 * Layer 5 — Resolved Furniture
 * Computed real dimensions. No ratios, fills, or abstract constraints.
 */

export interface ResolvedFurniture {
  id: string;
  name: string;
  root: ResolvedVolume;
}

export interface ResolvedVolume {
  id: string;
  label?: string;
  x: number;
  y: number;
  z: number;
  width: number;
  height: number;
  depth: number;
  children: ResolvedVolume[];
  features: ResolvedFeature[];
  fronts: ResolvedFront[];
  materialId?: string;
}

export interface ResolvedFeature {
  id: string;
  type: string;
  x: number;
  y: number;
  z: number;
  width: number;
  height: number;
  depth: number;
  params: Record<string, unknown>;
}

export interface ResolvedFront {
  id: string;
  type: string;
  x: number;
  y: number;
  z: number;
  width: number;
  height: number;
  depth: number;
  params: Record<string, unknown>;
}
