import type {
  DimensionConstraint,
  Feature,
  Front,
  FurnitureDefinition,
  SplitAxis,
  VolumeConstraints,
  VolumeNode,
  VolumeSplit,
} from '@/domain/furniture-definition';
import { EPSILON } from './geometry-math';
import type { ValidationError, ValidationResult } from './validation';

export type { ValidationError, ValidationResult } from './validation';

export interface VolumeNodeRef {
  node: VolumeNode;
  path: string[];
  depth: number;
  parent: VolumeNode | null;
}

export interface ConstraintSummary {
  width: DimensionConstraint | null;
  height: DimensionConstraint | null;
  depth: DimensionConstraint | null;
  hasFill: boolean;
  hasRatio: boolean;
  hasFixed: boolean;
}

const SPLIT_AXES: SplitAxis[] = ['x', 'y', 'z'];

function pushError(errors: ValidationError[], field: string, message: string): void {
  errors.push({ field, message });
}

export function axisToDimension(axis: SplitAxis): keyof VolumeConstraints {
  switch (axis) {
    case 'x':
      return 'width';
    case 'y':
      return 'height';
    case 'z':
      return 'depth';
  }
}

export function sumSplitRatios(split: VolumeSplit): number {
  if (!split.ratios?.length) return 0;
  return split.ratios.reduce((sum, r) => sum + r, 0);
}

export function summarizeConstraints(node: VolumeNode): ConstraintSummary {
  const { constraints } = node;
  const all = [constraints.width, constraints.height, constraints.depth];

  return {
    width: constraints.width ?? null,
    height: constraints.height ?? null,
    depth: constraints.depth ?? null,
    hasFill: all.some((c) => c?.mode === 'fill'),
    hasRatio: all.some((c) => c?.mode === 'ratio'),
    hasFixed: all.some((c) => c?.mode === 'fixed'),
  };
}

export function walkVolumeTree(
  root: VolumeNode,
  callback: (ref: VolumeNodeRef) => boolean | void,
): void {
  const visit = (node: VolumeNode, path: string[], depth: number, parent: VolumeNode | null) => {
    const ref: VolumeNodeRef = { node, path, depth, parent };
    if (callback(ref) === false) return;

    (node.children ?? []).forEach((child) => {
      visit(child, [...path, child.id], depth + 1, node);
    });
  };

  visit(root, [root.id], 0, null);
}

export function findVolumeNodeById(root: VolumeNode, id: string): VolumeNodeRef | null {
  let found: VolumeNodeRef | null = null;

  walkVolumeTree(root, (ref) => {
    if (ref.node.id === id) {
      found = ref;
      return false;
    }
  });

  return found;
}

export function flattenVolumeTree(root: VolumeNode): VolumeNodeRef[] {
  const nodes: VolumeNodeRef[] = [];
  walkVolumeTree(root, (ref) => {
    nodes.push(ref);
  });
  return nodes;
}

export function getTreeDepth(root: VolumeNode): number {
  let maxDepth = 0;
  walkVolumeTree(root, (ref) => {
    if (ref.depth > maxDepth) maxDepth = ref.depth;
  });
  return maxDepth;
}

export function getNodeCount(root: VolumeNode): number {
  let count = 0;
  walkVolumeTree(root, () => {
    count++;
  });
  return count;
}

export function getLeafNodes(root: VolumeNode): VolumeNode[] {
  const leaves: VolumeNode[] = [];
  walkVolumeTree(root, (ref) => {
    if (ref.node.children.length === 0) leaves.push(ref.node);
  });
  return leaves;
}

export function collectFeatures(root: VolumeNode): Feature[] {
  const features: Feature[] = [];
  walkVolumeTree(root, (ref) => {
    features.push(...ref.node.features);
  });
  return features;
}

export function collectFronts(root: VolumeNode): Front[] {
  const fronts: Front[] = [];
  walkVolumeTree(root, (ref) => {
    fronts.push(...ref.node.fronts);
  });
  return fronts;
}

export function countFeaturesByType(root: VolumeNode): Record<string, number> {
  const counts: Record<string, number> = {};
  for (const feature of collectFeatures(root)) {
    counts[feature.type] = (counts[feature.type] ?? 0) + 1;
  }
  return counts;
}

function validateDimensionConstraint(
  constraint: DimensionConstraint,
  field: string,
  errors: ValidationError[],
): void {
  switch (constraint.mode) {
    case 'fixed':
    case 'min':
    case 'max':
      if (constraint.value === undefined || constraint.value <= EPSILON) {
        pushError(errors, field, `mode "${constraint.mode}" requires a positive value`);
      }
      break;
    case 'ratio':
      if (constraint.value === undefined || constraint.value <= EPSILON || constraint.value > 1 + EPSILON) {
        pushError(errors, field, 'mode "ratio" requires a value between 0 and 1');
      }
      break;
    case 'fill':
      break;
    default:
      pushError(errors, field, `unknown constraint mode "${(constraint as DimensionConstraint).mode}"`);
  }
}

function validateConstraints(
  constraints: VolumeConstraints,
  prefix: string,
  errors: ValidationError[],
): void {
  const dims: (keyof VolumeConstraints)[] = ['width', 'height', 'depth'];
  for (const dim of dims) {
    const c = constraints[dim];
    if (c) validateDimensionConstraint(c, `${prefix}.constraints.${dim}`, errors);
  }
}

function validateSplit(split: VolumeSplit, childCount: number, prefix: string, errors: ValidationError[]): void {
  if (!SPLIT_AXES.includes(split.axis)) {
    pushError(errors, `${prefix}.split.axis`, `axis must be one of: ${SPLIT_AXES.join(', ')}`);
  }

  const hasRatios = Boolean(split.ratios?.length);
  const hasFixed = Boolean(split.fixed?.length);

  if (!hasRatios && !hasFixed) {
    pushError(errors, `${prefix}.split`, 'split requires ratios or fixed sizes');
    return;
  }

  if (hasRatios && hasFixed) {
    pushError(errors, `${prefix}.split`, 'split cannot define both ratios and fixed sizes');
  }

  if (hasRatios) {
    if (split.ratios!.length !== childCount) {
      pushError(
        errors,
        `${prefix}.split.ratios`,
        `ratio count (${split.ratios!.length}) must match children count (${childCount})`,
      );
    }
    for (let i = 0; i < split.ratios!.length; i++) {
      if (split.ratios![i] <= EPSILON) {
        pushError(errors, `${prefix}.split.ratios[${i}]`, 'ratio must be greater than zero');
      }
    }
    const total = sumSplitRatios(split);
    if (Math.abs(total - 1) > 0.01) {
      pushError(errors, `${prefix}.split.ratios`, `ratios must sum to 1 (got ${total})`);
    }
  }

  if (hasFixed) {
    if (split.fixed!.length !== childCount) {
      pushError(
        errors,
        `${prefix}.split.fixed`,
        `fixed count (${split.fixed!.length}) must match children count (${childCount})`,
      );
    }
    for (let i = 0; i < split.fixed!.length; i++) {
      if (split.fixed![i] <= EPSILON) {
        pushError(errors, `${prefix}.split.fixed[${i}]`, 'fixed size must be greater than zero');
      }
    }
  }
}

function validateChildSplitAlignment(
  parent: VolumeNode,
  child: VolumeNode,
  childIndex: number,
  prefix: string,
  errors: ValidationError[],
): void {
  if (!parent.split) return;

  const dim = axisToDimension(parent.split.axis);
  const childConstraint = child.constraints[dim];

  if (!childConstraint) {
    pushError(
      errors,
      `${prefix}.children[${childIndex}].constraints.${dim}`,
      `child must define a ${dim} constraint matching parent split on axis "${parent.split.axis}"`,
    );
    return;
  }

  if (parent.split.ratios?.length) {
    const expected = parent.split.ratios[childIndex];
    if (childConstraint.mode === 'ratio' && Math.abs(childConstraint.value! - expected) > 0.01) {
      pushError(
        errors,
        `${prefix}.children[${childIndex}].constraints.${dim}`,
        `ratio value ${childConstraint.value} does not match parent split ratio ${expected}`,
      );
    }
    if (childConstraint.mode === 'fixed') {
      pushError(
        errors,
        `${prefix}.children[${childIndex}].constraints.${dim}`,
        'child cannot use fixed constraint when parent splits by ratios',
      );
    }
  }

  if (parent.split.fixed?.length) {
    const expected = parent.split.fixed[childIndex];
    if (childConstraint.mode === 'fixed' && Math.abs(childConstraint.value! - expected) > EPSILON) {
      pushError(
        errors,
        `${prefix}.children[${childIndex}].constraints.${dim}`,
        `fixed value ${childConstraint.value} does not match parent split fixed ${expected}`,
      );
    }
    if (childConstraint.mode === 'ratio') {
      pushError(
        errors,
        `${prefix}.children[${childIndex}].constraints.${dim}`,
        'child cannot use ratio constraint when parent splits by fixed sizes',
      );
    }
  }
}

function validateFeature(feature: Feature, prefix: string, errors: ValidationError[]): void {
  if (!feature.id) pushError(errors, `${prefix}.id`, 'id is required');
  if (!feature.type) pushError(errors, `${prefix}.type`, 'type is required');
}

function validateFront(front: Front, prefix: string, errors: ValidationError[]): void {
  if (!front.id) pushError(errors, `${prefix}.id`, 'id is required');
  if (!front.type) pushError(errors, `${prefix}.type`, 'type is required');
}

function validateVolumeNode(
  node: VolumeNode,
  prefix: string,
  errors: ValidationError[],
  ids: Map<string, string>,
): void {
  if (!node.id) {
    pushError(errors, `${prefix}.id`, 'id is required');
  } else {
    const existing = ids.get(node.id);
    if (existing) {
      pushError(errors, `${prefix}.id`, `duplicate id "${node.id}" (also used in ${existing})`);
    } else {
      ids.set(node.id, prefix);
    }
  }

  validateConstraints(node.constraints, prefix, errors);

  const children = node.children ?? [];
  const hasChildren = children.length > 0;
  const hasSplit = Boolean(node.split);

  if (hasChildren && !hasSplit) {
    pushError(errors, `${prefix}.children`, 'nodes with children must define a split');
  }
  if (hasSplit && !hasChildren) {
    pushError(errors, `${prefix}.split`, 'split requires at least one child');
  }
  if (hasSplit && hasChildren) {
    validateSplit(node.split!, children.length, prefix, errors);
  }

  (node.features ?? []).forEach((feature, i) => {
    const field = `${prefix}.features[${i}].id`;
    validateFeature(feature, `${prefix}.features[${i}]`, errors);
    if (feature.id) {
      const existing = ids.get(feature.id);
      if (existing) {
        pushError(errors, field, `duplicate id "${feature.id}" (also used in ${existing})`);
      } else {
        ids.set(feature.id, field);
      }
    }
  });

  (node.fronts ?? []).forEach((front, i) => {
    const field = `${prefix}.fronts[${i}].id`;
    validateFront(front, `${prefix}.fronts[${i}]`, errors);
    if (front.id) {
      const existing = ids.get(front.id);
      if (existing) {
        pushError(errors, field, `duplicate id "${front.id}" (also used in ${existing})`);
      } else {
        ids.set(front.id, field);
      }
    }
  });

  children.forEach((child, i) => {
    validateVolumeNode(child, `${prefix}.children[${i}]`, errors, ids);
    validateChildSplitAlignment(node, child, i, prefix, errors);
  });
}

export function validateFurnitureDefinition(furniture: FurnitureDefinition): ValidationResult {
  const errors: ValidationError[] = [];

  if (!furniture.id) pushError(errors, 'id', 'id is required');
  if (!furniture.name) pushError(errors, 'name', 'name is required');

  const ids = new Map<string, string>();
  validateVolumeNode(furniture.root, 'root', errors, ids);

  return { valid: errors.length === 0, errors };
}
