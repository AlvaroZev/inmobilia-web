import type { FurnitureDefinition, VolumeNode } from '@/domain/furniture-definition';

function normalizeVolumeNode(raw: VolumeNode): VolumeNode {
  return {
    ...raw,
    children: (raw.children ?? []).map(normalizeVolumeNode),
    features: raw.features ?? [],
    fronts: raw.fronts ?? [],
  };
}

/** Normalizes API JSON where Go encodes empty slices as null. */
export function normalizeFurnitureDefinition(raw: FurnitureDefinition): FurnitureDefinition {
  return {
    ...raw,
    root: normalizeVolumeNode(raw.root),
  };
}
