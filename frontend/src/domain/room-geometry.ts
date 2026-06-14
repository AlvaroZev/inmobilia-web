import type { Plane, Point3D, Polygon2D } from './geometry';

/**
 * Layer 1 — Room Geometry
 * Represents the physical environment. Independent from furniture definition.
 */

export interface RoomGeometry {
  id: string;
  name?: string;
  perimeter: Polygon2D;
  floor: Plane;
  ceiling: Plane;
  walls: Wall[];
  openings: Opening[];
  obstacles: Obstacle[];
}

export interface Wall {
  id: string;
  /** Vertices along the wall face (typically bottom then top edge). */
  vertices: Point3D[];
  thickness: number;
  /** Deviation from vertical expressed as a vector, not an angle. */
  outOfPlumb?: Point3D;
}

export type OpeningType = 'door' | 'window';

export interface Opening {
  id: string;
  type: OpeningType;
  wallId: string;
  /** Bottom-left corner of the opening on the wall face. */
  origin: Point3D;
  width: number;
  height: number;
}

export type ObstacleType =
  | 'column'
  | 'beam'
  | 'pipe'
  | 'skirting'
  | 'other';

export interface Obstacle {
  id: string;
  type: ObstacleType;
  label?: string;
  /** Obstacle bounds in room coordinates. */
  bounds: { min: Point3D; max: Point3D };
  /** Optional profile vertices for non-rectangular obstacles. */
  profile?: Point3D[];
}
