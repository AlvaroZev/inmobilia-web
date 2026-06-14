/**
 * Layer 6 — Manufacturing Model
 * Physical parts derived from resolved furniture.
 */

export interface ManufacturingModel {
  furnitureId: string;
  parts: Part[];
  hardware: Hardware[];
  edgeBanding: EdgeBanding[];
  drilling: Drilling[];
}

export interface Material {
  id: string;
  name: string;
  type: string;
  thickness: number;
  color?: string;
}

export type PartType =
  | 'lateral'
  | 'base'
  | 'top'
  | 'shelf'
  | 'door'
  | 'divider'
  | 'back'
  | 'drawer_side'
  | 'drawer_bottom'
  | 'other';

export interface Part {
  id: string;
  name: string;
  type: PartType | string;
  volumeId: string;
  width: number;
  height: number;
  thickness: number;
  material: Material;
  grainDirection?: 'horizontal' | 'vertical';
}

export interface Hardware {
  id: string;
  type: string;
  quantity: number;
  partIds?: string[];
  params?: Record<string, unknown>;
}

export type EdgeSide = 'top' | 'bottom' | 'left' | 'right';

export interface EdgeBanding {
  partId: string;
  edge: EdgeSide;
  material: string;
  length: number;
}

export interface Drilling {
  partId: string;
  x: number;
  y: number;
  diameter: number;
  depth: number;
  type: string;
}
