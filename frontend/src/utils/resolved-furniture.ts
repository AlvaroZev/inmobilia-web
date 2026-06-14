import type { BoundingBox3D } from '@/domain/geometry';
import type {
  ResolvedFeature,
  ResolvedFront,
  ResolvedFurniture,
  ResolvedVolume,
} from '@/domain/resolved-furniture';
import { EPSILON, normalizeBounds } from './geometry-math';
import type { ValidationError, ValidationResult } from './validation';

export interface ResolvedVolumeRef {
  volume: ResolvedVolume;
  path: string[];
  depth: number;
  parent: ResolvedVolume | null;
}

function pushError(errors: ValidationError[], field: string, message: string): void {
  errors.push({ field, message });
}

export function volumeBounds(volume: ResolvedVolume): BoundingBox3D {
  return {
    min: { x: volume.x, y: volume.y, z: volume.z },
    max: {
      x: volume.x + volume.width,
      y: volume.y + volume.height,
      z: volume.z + volume.depth,
    },
  };
}

export function furnitureBounds(furniture: ResolvedFurniture): BoundingBox3D {
  return volumeBounds(furniture.root);
}

export function volumeContains(outer: ResolvedVolume, inner: ResolvedVolume): boolean {
  const ob = volumeBounds(outer);
  const ib = volumeBounds(inner);

  return (
    ib.min.x >= ob.min.x - EPSILON &&
    ib.min.y >= ob.min.y - EPSILON &&
    ib.min.z >= ob.min.z - EPSILON &&
    ib.max.x <= ob.max.x + EPSILON &&
    ib.max.y <= ob.max.y + EPSILON &&
    ib.max.z <= ob.max.z + EPSILON
  );
}

export function boundsOverlap(a: BoundingBox3D, b: BoundingBox3D): boolean {
  const na = normalizeBounds(a);
  const nb = normalizeBounds(b);

  return (
    na.min.x < nb.max.x - EPSILON &&
    na.max.x > nb.min.x + EPSILON &&
    na.min.y < nb.max.y - EPSILON &&
    na.max.y > nb.min.y + EPSILON &&
    na.min.z < nb.max.z - EPSILON &&
    na.max.z > nb.min.z + EPSILON
  );
}

export function volumeOverlap(a: ResolvedVolume, b: ResolvedVolume): boolean {
  return boundsOverlap(volumeBounds(a), volumeBounds(b));
}

export function compartmentVolumeMm3(volume: ResolvedVolume): number {
  return volume.width * volume.height * volume.depth;
}

export function externalDimensions(furniture: ResolvedFurniture): {
  width: number;
  height: number;
  depth: number;
} {
  return {
    width: furniture.root.width,
    height: furniture.root.height,
    depth: furniture.root.depth,
  };
}

export function walkResolvedTree(
  root: ResolvedVolume,
  callback: (ref: ResolvedVolumeRef) => boolean | void,
): void {
  const visit = (
    volume: ResolvedVolume,
    path: string[],
    depth: number,
    parent: ResolvedVolume | null,
  ) => {
    const ref: ResolvedVolumeRef = { volume, path, depth, parent };
    if (callback(ref) === false) return;

    volume.children.forEach((child) => {
      visit(child, [...path, child.id], depth + 1, volume);
    });
  };

  visit(root, [root.id], 0, null);
}

export function findResolvedVolumeById(
  root: ResolvedVolume,
  id: string,
): ResolvedVolumeRef | null {
  let found: ResolvedVolumeRef | null = null;

  walkResolvedTree(root, (ref) => {
    if (ref.volume.id === id) {
      found = ref;
      return false;
    }
  });

  return found;
}

export function flattenResolvedTree(root: ResolvedVolume): ResolvedVolumeRef[] {
  const nodes: ResolvedVolumeRef[] = [];
  walkResolvedTree(root, (ref) => {
    nodes.push(ref);
  });
  return nodes;
}

export function getResolvedTreeDepth(root: ResolvedVolume): number {
  let maxDepth = 0;
  walkResolvedTree(root, (ref) => {
    if (ref.depth > maxDepth) maxDepth = ref.depth;
  });
  return maxDepth;
}

export function getResolvedNodeCount(root: ResolvedVolume): number {
  let count = 0;
  walkResolvedTree(root, () => {
    count++;
  });
  return count;
}

export function getResolvedLeafVolumes(root: ResolvedVolume): ResolvedVolume[] {
  const leaves: ResolvedVolume[] = [];
  walkResolvedTree(root, (ref) => {
    if (ref.volume.children.length === 0) leaves.push(ref.volume);
  });
  return leaves;
}

export function collectResolvedFeatures(root: ResolvedVolume): ResolvedFeature[] {
  const features: ResolvedFeature[] = [];
  walkResolvedTree(root, (ref) => {
    features.push(...ref.volume.features);
  });
  return features;
}

export function collectResolvedFronts(root: ResolvedVolume): ResolvedFront[] {
  const fronts: ResolvedFront[] = [];
  walkResolvedTree(root, (ref) => {
    fronts.push(...ref.volume.fronts);
  });
  return fronts;
}

export function totalLeafVolumeMm3(root: ResolvedVolume): number {
  return getResolvedLeafVolumes(root).reduce((sum, leaf) => sum + compartmentVolumeMm3(leaf), 0);
}

function elementBounds(element: {
  x: number;
  y: number;
  z: number;
  width: number;
  height: number;
  depth: number;
}): BoundingBox3D {
  return {
    min: { x: element.x, y: element.y, z: element.z },
    max: {
      x: element.x + element.width,
      y: element.y + element.height,
      z: element.z + element.depth,
    },
  };
}

function elementContainedInVolume(
  bounds: BoundingBox3D,
  volume: ResolvedVolume,
): boolean {
  const vb = volumeBounds(volume);
  const b = normalizeBounds(bounds);

  return (
    b.min.x >= vb.min.x - EPSILON &&
    b.min.y >= vb.min.y - EPSILON &&
    b.min.z >= vb.min.z - EPSILON &&
    b.max.x <= vb.max.x + EPSILON &&
    b.max.y <= vb.max.y + EPSILON &&
    b.max.z <= vb.max.z + EPSILON
  );
}

function validateDimensions(volume: ResolvedVolume, prefix: string, errors: ValidationError[]): void {
  if (volume.width <= EPSILON) pushError(errors, `${prefix}.width`, 'width must be greater than zero');
  if (volume.height <= EPSILON) pushError(errors, `${prefix}.height`, 'height must be greater than zero');
  if (volume.depth <= EPSILON) pushError(errors, `${prefix}.depth`, 'depth must be greater than zero');
}

function validateResolvedVolume(
  volume: ResolvedVolume,
  prefix: string,
  errors: ValidationError[],
  ids: Map<string, string>,
  parent: ResolvedVolume | null,
): void {
  if (!volume.id) {
    pushError(errors, `${prefix}.id`, 'id is required');
  } else {
    const existing = ids.get(volume.id);
    if (existing) {
      pushError(errors, `${prefix}.id`, `duplicate id "${volume.id}" (also used in ${existing})`);
    } else {
      ids.set(volume.id, prefix);
    }
  }

  validateDimensions(volume, prefix, errors);

  if (parent && !volumeContains(parent, volume)) {
    pushError(errors, prefix, 'volume exceeds parent bounds');
  }

  volume.children.forEach((child, i) => {
    validateResolvedVolume(child, `${prefix}.children[${i}]`, errors, ids, volume);
  });

  for (let i = 0; i < volume.children.length; i++) {
    for (let j = i + 1; j < volume.children.length; j++) {
      if (volumeOverlap(volume.children[i], volume.children[j])) {
        pushError(
          errors,
          `${prefix}.children[${i}]`,
          `overlaps with sibling "${volume.children[j].id}"`,
        );
      }
    }
  }

  volume.features.forEach((feature, i) => {
    const field = `${prefix}.features[${i}]`;
    if (!feature.id) pushError(errors, `${field}.id`, 'id is required');
    if (!feature.type) pushError(errors, `${field}.type`, 'type is required');
    if (feature.id) {
      const existing = ids.get(feature.id);
      if (existing) {
        pushError(errors, `${field}.id`, `duplicate id "${feature.id}" (also used in ${existing})`);
      } else {
        ids.set(feature.id, field);
      }
    }
    if (!elementContainedInVolume(elementBounds(feature), volume)) {
      pushError(errors, field, 'feature exceeds volume bounds');
    }
  });

  volume.fronts.forEach((front, i) => {
    const field = `${prefix}.fronts[${i}]`;
    if (!front.id) pushError(errors, `${field}.id`, 'id is required');
    if (!front.type) pushError(errors, `${field}.type`, 'type is required');
    if (front.id) {
      const existing = ids.get(front.id);
      if (existing) {
        pushError(errors, `${field}.id`, `duplicate id "${front.id}" (also used in ${existing})`);
      } else {
        ids.set(front.id, field);
      }
    }
    if (!elementContainedInVolume(elementBounds(front), volume)) {
      pushError(errors, field, 'front exceeds volume bounds');
    }
  });
}

/** Validates that resolved furniture contains only real dimensions and valid geometry. */
export function validateResolvedFurniture(furniture: ResolvedFurniture): ValidationResult {
  const errors: ValidationError[] = [];

  if (!furniture.id) pushError(errors, 'id', 'id is required');
  if (!furniture.name) pushError(errors, 'name', 'name is required');

  const ids = new Map<string, string>();
  validateResolvedVolume(furniture.root, 'root', errors, ids, null);

  return { valid: errors.length === 0, errors };
}
