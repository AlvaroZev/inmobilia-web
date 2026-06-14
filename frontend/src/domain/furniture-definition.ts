/**
 * Layer 2 — Furniture Definition
 * Design intent only. No final geometry, cuts, or parts.
 */

export interface FurnitureDefinition {
  id: string;
  name: string;
  description?: string;
  root: VolumeNode;
  metadata?: Record<string, unknown>;
}

export interface VolumeNode {
  id: string;
  label?: string;
  constraints: VolumeConstraints;
  split?: VolumeSplit;
  children: VolumeNode[];
  features: Feature[];
  fronts: Front[];
  adaptation?: AdaptationRules;
  manufacturing?: ManufacturingHints;
}

export interface VolumeConstraints {
  width?: DimensionConstraint;
  height?: DimensionConstraint;
  depth?: DimensionConstraint;
}

export type DimensionConstraint =
  | { mode: 'fixed'; value: number }
  | { mode: 'ratio'; value: number }
  | { mode: 'fill' }
  | { mode: 'min'; value: number }
  | { mode: 'max'; value: number };

export type SplitAxis = 'x' | 'y' | 'z';

export interface VolumeSplit {
  axis: SplitAxis;
  /** Relative proportions for child volumes (sum should equal 1). */
  ratios?: number[];
  /** Fixed sizes in mm for children that are not ratio-based. */
  fixed?: number[];
}

/**
 * Extensible feature — never use a closed enum of all furniture types.
 * Examples: shelf_set, drawer_stack, hanger_rod, divider, lighting, appliance_space
 */
export interface Feature {
  id: string;
  type: string;
  params: Record<string, unknown>;
}

/**
 * Visible front elements.
 * Examples: door, sliding_door, glass, mirror, drawer_front
 */
export interface Front {
  id: string;
  type: string;
  params: Record<string, unknown>;
}

export interface AdaptationRules {
  followFloor?: boolean;
  followCeiling?: boolean;
  followWall?: boolean;
  compensateSkirting?: boolean;
  params?: Record<string, unknown>;
}

export interface ManufacturingHints {
  materialId?: string;
  edgeBanding?: string;
  backPanel?: boolean;
  params?: Record<string, unknown>;
}
