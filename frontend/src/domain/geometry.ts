/** All dimensions in millimeters. Angles are never stored — derive from geometry. */

export interface Point2D {
  x: number;
  y: number;
}

export interface Point3D {
  x: number;
  y: number;
  z: number;
}

/** Plane defined by a point on the surface and its outward normal. */
export interface Plane {
  point: Point3D;
  normal: Point3D;
}

export interface Polygon2D {
  vertices: Point2D[];
}

export interface BoundingBox3D {
  min: Point3D;
  max: Point3D;
}
