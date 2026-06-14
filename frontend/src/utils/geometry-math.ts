import type { BoundingBox3D, Plane, Point2D, Point3D, Polygon2D } from '@/domain/geometry';

/** Default tolerance in millimeters. */
export const EPSILON = 0.001;

export function subtract3(a: Point3D, b: Point3D): Point3D {
  return { x: a.x - b.x, y: a.y - b.y, z: a.z - b.z };
}

export function add3(a: Point3D, b: Point3D): Point3D {
  return { x: a.x + b.x, y: a.y + b.y, z: a.z + b.z };
}

export function scale3(v: Point3D, s: number): Point3D {
  return { x: v.x * s, y: v.y * s, z: v.z * s };
}

export function dot3(a: Point3D, b: Point3D): number {
  return a.x * b.x + a.y * b.y + a.z * b.z;
}

export function cross3(a: Point3D, b: Point3D): Point3D {
  return {
    x: a.y * b.z - a.z * b.y,
    y: a.z * b.x - a.x * b.z,
    z: a.x * b.y - a.y * b.x,
  };
}

export function magnitude3(v: Point3D): number {
  return Math.hypot(v.x, v.y, v.z);
}

export function normalize3(v: Point3D): Point3D | null {
  const len = magnitude3(v);
  if (len < EPSILON) return null;
  return scale3(v, 1 / len);
}

export function distance3(a: Point3D, b: Point3D): number {
  return magnitude3(subtract3(a, b));
}

export function distance2(a: Point2D, b: Point2D): number {
  return Math.hypot(a.x - b.x, a.y - b.y);
}

/** Angle between two vectors in radians. Never stored in domain models. */
export function angleBetween3(a: Point3D, b: Point3D): number {
  const na = normalize3(a);
  const nb = normalize3(b);
  if (!na || !nb) return 0;
  const cos = Math.min(1, Math.max(-1, dot3(na, nb)));
  return Math.acos(cos);
}

export function angleBetween2(a: Point2D, b: Point2D): number {
  const lenA = Math.hypot(a.x, a.y);
  const lenB = Math.hypot(b.x, b.y);
  if (lenA < EPSILON || lenB < EPSILON) return 0;
  const cos = Math.min(
    1,
    Math.max(-1, (a.x * b.x + a.y * b.y) / (lenA * lenB)),
  );
  return Math.acos(cos);
}

export function planeFromPoints(a: Point3D, b: Point3D, c: Point3D): Plane | null {
  const ab = subtract3(b, a);
  const ac = subtract3(c, a);
  const normal = normalize3(cross3(ab, ac));
  if (!normal) return null;
  return { point: { ...a }, normal };
}

export function signedDistanceToPlane(point: Point3D, plane: Plane): number {
  return dot3(subtract3(point, plane.point), plane.normal);
}

export function projectPointOnPlane(point: Point3D, plane: Plane): Point3D {
  const dist = signedDistanceToPlane(point, plane);
  return subtract3(point, scale3(plane.normal, dist));
}

export function normalizePlane(plane: Plane): Plane | null {
  const normal = normalize3(plane.normal);
  if (!normal) return null;
  return { point: { ...plane.point }, normal };
}

/** Shoelace formula — signed area in mm². */
export function polygonArea2D(polygon: Polygon2D): number {
  const { vertices } = polygon;
  if (vertices.length < 3) return 0;

  let sum = 0;
  for (let i = 0; i < vertices.length; i++) {
    const j = (i + 1) % vertices.length;
    sum += vertices[i].x * vertices[j].y - vertices[j].x * vertices[i].y;
  }
  return sum / 2;
}

export function polygonPerimeter2D(polygon: Polygon2D): number {
  const { vertices } = polygon;
  if (vertices.length < 2) return 0;

  let perimeter = 0;
  for (let i = 0; i < vertices.length; i++) {
    const j = (i + 1) % vertices.length;
    perimeter += distance2(vertices[i], vertices[j]);
  }
  return perimeter;
}

/** Interior angle at a polygon vertex in radians. */
export function interiorAngleAtVertex2D(polygon: Polygon2D, index: number): number {
  const { vertices } = polygon;
  const n = vertices.length;
  if (n < 3 || index < 0 || index >= n) return 0;

  const prev = vertices[(index - 1 + n) % n];
  const curr = vertices[index];
  const next = vertices[(index + 1) % n];

  const incoming = { x: prev.x - curr.x, y: prev.y - curr.y };
  const outgoing = { x: next.x - curr.x, y: next.y - curr.y };

  return Math.PI - angleBetween2(incoming, outgoing);
}

export function normalizeBounds(bounds: BoundingBox3D): BoundingBox3D {
  return {
    min: {
      x: Math.min(bounds.min.x, bounds.max.x),
      y: Math.min(bounds.min.y, bounds.max.y),
      z: Math.min(bounds.min.z, bounds.max.z),
    },
    max: {
      x: Math.max(bounds.min.x, bounds.max.x),
      y: Math.max(bounds.min.y, bounds.max.y),
      z: Math.max(bounds.min.z, bounds.max.z),
    },
  };
}

export function boundsVolume(bounds: BoundingBox3D): number {
  const b = normalizeBounds(bounds);
  return (
    (b.max.x - b.min.x) * (b.max.y - b.min.y) * (b.max.z - b.min.z)
  );
}
