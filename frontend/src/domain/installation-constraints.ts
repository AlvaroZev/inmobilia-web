import type { BoundingBox3D } from './geometry';

/**
 * Layer 3 — Installation Constraints
 * Connects room geometry and furniture definition.
 */

export interface InstallationConstraints {
  id: string;
  zone: InstallationZone;
  clearances: Clearances;
  tolerances: Tolerances;
  references: InstallationReferences;
}

export interface InstallationZone {
  /** Walls the furniture anchors to. */
  anchorWallIds?: string[];
  /** Optional bounding region within the room. */
  bounds?: BoundingBox3D;
}

export interface Clearances {
  top: number;
  bottom: number;
  left: number;
  right: number;
  back: number;
  front: number;
}

export interface Tolerances {
  width: number;
  height: number;
  depth: number;
}

export interface InstallationReferences {
  floorOffset: number;
  ceilingOffset?: number;
  wallOffset?: number;
  /** Reference wall for horizontal positioning. */
  referenceWallId?: string;
}
