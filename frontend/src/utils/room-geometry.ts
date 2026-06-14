import type { Plane, Point3D } from '@/domain/geometry';
import type { Obstacle, Opening, RoomGeometry, Wall } from '@/domain/room-geometry';
import type { ValidationError, ValidationResult } from './validation';
import {
  EPSILON,
  angleBetween3,
  boundsVolume,
  distance3,
  dot3,
  magnitude3,
  normalize3,
  normalizeBounds,
  normalizePlane,
  planeFromPoints,
  polygonArea2D,
  polygonPerimeter2D,
  interiorAngleAtVertex2D,
  projectPointOnPlane,
  signedDistanceToPlane,
  subtract3,
} from './geometry-math';

export type { ValidationError, ValidationResult } from './validation';

export interface WallLocalFrame {
  origin: Point3D;
  /** Unit vector along the bottom edge (left → right). */
  u: Point3D;
  /** Unit vector along the left edge (bottom → top). */
  v: Point3D;
  normal: Point3D;
  width: number;
  height: number;
}

function pushError(errors: ValidationError[], field: string, message: string): void {
  errors.push({ field, message });
}

export function findWallById(room: RoomGeometry, wallId: string): Wall | undefined {
  return room.walls.find((w) => w.id === wallId);
}

export function getWallPlane(wall: Wall): Plane | null {
  if (wall.vertices.length < 3) return null;
  return planeFromPoints(wall.vertices[0], wall.vertices[1], wall.vertices[2]);
}

/**
 * Quad wall convention: v0–v1 bottom edge, v1–v2 right edge, v2–v3 top, v3–v0 left.
 */
export function getWallLocalFrame(wall: Wall): WallLocalFrame | null {
  if (wall.vertices.length < 4) return null;

  const v0 = wall.vertices[0];
  const v1 = wall.vertices[1];
  const v3 = wall.vertices[3];

  const uRaw = subtract3(v1, v0);
  const vRaw = subtract3(v3, v0);
  const u = normalize3(uRaw);
  const v = normalize3(vRaw);
  if (!u || !v) return null;

  const normal = normalize3({
    x: u.y * v.z - u.z * v.y,
    y: u.z * v.x - u.x * v.z,
    z: u.x * v.y - u.y * v.x,
  });
  if (!normal) return null;

  return {
    origin: v0,
    u,
    v,
    normal,
    width: magnitude3(uRaw),
    height: magnitude3(vRaw),
  };
}

export function projectPointOnWallLocal(
  wall: Wall,
  point: Point3D,
): { u: number; v: number } | null {
  const frame = getWallLocalFrame(wall);
  if (!frame) return null;

  const onPlane = projectPointOnPlane(point, {
    point: frame.origin,
    normal: frame.normal,
  });
  const rel = subtract3(onPlane, frame.origin);

  return {
    u: dot3(rel, frame.u),
    v: dot3(rel, frame.v),
  };
}

/** Bottom edge length in mm. */
export function getWallBottomEdgeLength(wall: Wall): number {
  if (wall.vertices.length < 2) return 0;
  return distance3(wall.vertices[0], wall.vertices[1]);
}

/** Wall height from left edge in mm. */
export function getWallHeight(wall: Wall): number {
  if (wall.vertices.length < 4) return 0;
  return distance3(wall.vertices[0], wall.vertices[3]);
}

/**
 * Bottom edge angle in the floor plan (XZ plane) in radians.
 * Supports non-orthogonal walls.
 */
export function getWallBottomEdgeAngle(wall: Wall): number {
  if (wall.vertices.length < 2) return 0;
  const edge = subtract3(wall.vertices[1], wall.vertices[0]);
  return Math.atan2(edge.z, edge.x);
}

/** Deviation from true vertical, computed from geometry when not stored. */
export function getWallOutOfPlumbVector(wall: Wall): Point3D {
  if (wall.outOfPlumb) return wall.outOfPlumb;
  if (wall.vertices.length < 4) return { x: 0, y: 0, z: 0 };

  const vertical = subtract3(wall.vertices[3], wall.vertices[0]);
  const trueUp = { x: 0, y: 1, z: 0 };
  const projectedLength = dot3(vertical, trueUp);
  const plumb = { x: 0, y: projectedLength, z: 0 };
  return subtract3(vertical, plumb);
}

export function getWallOutOfPlumbMagnitude(wall: Wall): number {
  return magnitude3(getWallOutOfPlumbVector(wall));
}

/** Vertical distance between floor and ceiling planes in mm. */
export function getFloorCeilingHeight(room: RoomGeometry): number {
  const floor = normalizePlane(room.floor);
  const ceiling = normalizePlane(room.ceiling);
  if (!floor || !ceiling) return 0;
  return Math.abs(signedDistanceToPlane(ceiling.point, floor));
}

export function getRoomFloorArea(room: RoomGeometry): number {
  return Math.abs(polygonArea2D(room.perimeter));
}

export function getRoomPerimeterLength(room: RoomGeometry): number {
  return polygonPerimeter2D(room.perimeter);
}

/** Interior corner angles of the floor perimeter in radians. */
export function getPerimeterInteriorAngles(room: RoomGeometry): number[] {
  return room.perimeter.vertices.map((_, i) =>
    interiorAngleAtVertex2D(room.perimeter, i),
  );
}

export function getSkirtingObstacles(room: RoomGeometry): Obstacle[] {
  return room.obstacles.filter((o) => o.type === 'skirting');
}

export function getSkirtingHeight(obstacle: Obstacle): number {
  const bounds = normalizeBounds(obstacle.bounds);
  return bounds.max.y - bounds.min.y;
}

function validatePolygon(
  polygon: RoomGeometry['perimeter'],
  prefix: string,
  errors: ValidationError[],
): void {
  if (polygon.vertices.length < 3) {
    pushError(errors, `${prefix}.vertices`, 'at least 3 vertices required');
    return;
  }

  for (let i = 0; i < polygon.vertices.length; i++) {
    const j = (i + 1) % polygon.vertices.length;
    const a = polygon.vertices[i];
    const b = polygon.vertices[j];
    if (Math.hypot(a.x - b.x, a.y - b.y) < EPSILON) {
      pushError(errors, `${prefix}.vertices[${i}]`, 'zero-length edge');
    }
  }

  if (Math.abs(polygonArea2D(polygon)) < EPSILON) {
    pushError(errors, prefix, 'polygon area must be greater than zero');
  }
}

function validatePlane(plane: Plane, field: string, errors: ValidationError[]): void {
  if (magnitude3(plane.normal) < EPSILON) {
    pushError(errors, `${field}.normal`, 'normal vector must be non-zero');
  }
}

function validateWall(wall: Wall, prefix: string, errors: ValidationError[]): void {
  if (!wall.id) {
    pushError(errors, `${prefix}.id`, 'id is required');
  }
  if (wall.thickness <= EPSILON) {
    pushError(errors, `${prefix}.thickness`, 'thickness must be greater than zero');
  }
  if (wall.vertices.length < 3) {
    pushError(errors, `${prefix}.vertices`, 'at least 3 vertices required');
    return;
  }

  const facePlane = getWallPlane(wall);
  if (!facePlane) {
    pushError(errors, `${prefix}.vertices`, 'vertices must define a non-degenerate face');
  }

  if (wall.vertices.length >= 4) {
    const frame = getWallLocalFrame(wall);
    if (!frame) {
      pushError(errors, `${prefix}.vertices`, 'quad wall frame is degenerate');
    } else if (frame.width < EPSILON || frame.height < EPSILON) {
      pushError(errors, `${prefix}.vertices`, 'wall width and height must be greater than zero');
    }
  }
}

function validateOpening(
  opening: Opening,
  wall: Wall | undefined,
  prefix: string,
  errors: ValidationError[],
): void {
  if (!opening.id) {
    pushError(errors, `${prefix}.id`, 'id is required');
  }
  if (opening.width <= EPSILON) {
    pushError(errors, `${prefix}.width`, 'width must be greater than zero');
  }
  if (opening.height <= EPSILON) {
    pushError(errors, `${prefix}.height`, 'height must be greater than zero');
  }
  if (!wall) {
    pushError(errors, `${prefix}.wallId`, `wall "${opening.wallId}" not found`);
    return;
  }

  const frame = getWallLocalFrame(wall);
  if (!frame) {
    pushError(errors, `${prefix}.wallId`, 'wall geometry is not a valid quad face');
    return;
  }

  const dist = Math.abs(
    signedDistanceToPlane(opening.origin, {
      point: frame.origin,
      normal: frame.normal,
    }),
  );
  if (dist > 5) {
    pushError(errors, `${prefix}.origin`, 'opening origin is not on the wall face');
  }

  const local = projectPointOnWallLocal(wall, opening.origin);
  if (!local) return;

  if (local.u < -EPSILON || local.v < -EPSILON) {
    pushError(errors, `${prefix}.origin`, 'opening origin is outside wall bounds');
  }
  if (local.u + opening.width > frame.width + EPSILON) {
    pushError(errors, `${prefix}.width`, 'opening exceeds wall width');
  }
  if (local.v + opening.height > frame.height + EPSILON) {
    pushError(errors, `${prefix}.height`, 'opening exceeds wall height');
  }
}

function validateObstacle(obstacle: Obstacle, prefix: string, errors: ValidationError[]): void {
  if (!obstacle.id) {
    pushError(errors, `${prefix}.id`, 'id is required');
  }

  const bounds = normalizeBounds(obstacle.bounds);
  if (boundsVolume(bounds) < EPSILON) {
    pushError(errors, `${prefix}.bounds`, 'bounds must have positive volume');
  }
}

export function validateRoomGeometry(room: RoomGeometry): ValidationResult {
  const errors: ValidationError[] = [];

  if (!room.id) {
    pushError(errors, 'id', 'id is required');
  }

  validatePolygon(room.perimeter, 'perimeter', errors);
  validatePlane(room.floor, 'floor', errors);
  validatePlane(room.ceiling, 'ceiling', errors);

  const ids = new Map<string, string>();

  const trackId = (id: string, field: string) => {
    if (!id) return;
    const existing = ids.get(id);
    if (existing) {
      pushError(errors, field, `duplicate id "${id}" (also used in ${existing})`);
    } else {
      ids.set(id, field);
    }
  };

  room.walls.forEach((wall, i) => {
    trackId(wall.id, `walls[${i}].id`);
    validateWall(wall, `walls[${i}]`, errors);
  });

  room.openings.forEach((opening, i) => {
    trackId(opening.id, `openings[${i}].id`);
    const wall = findWallById(room, opening.wallId);
    validateOpening(opening, wall, `openings[${i}]`, errors);
  });

  room.obstacles.forEach((obstacle, i) => {
    trackId(obstacle.id, `obstacles[${i}].id`);
    validateObstacle(obstacle, `obstacles[${i}]`, errors);
  });

  return { valid: errors.length === 0, errors };
}

/** Angle between floor and ceiling planes in radians (for sloped ceilings). */
export function getFloorCeilingAngle(room: RoomGeometry): number {
  const floor = normalizePlane(room.floor);
  const ceiling = normalizePlane(room.ceiling);
  if (!floor || !ceiling) return 0;
  return angleBetween3(floor.normal, ceiling.normal);
}
