import type { ResolvedFurniture, ResolvedVolume } from '@/domain/resolved-furniture';
import { furnitureBounds } from '@/utils/resolved-furniture';

/** Converts millimeters to Three.js scene units (meters). */
export const MM_TO_M = 0.001;

export function mm(value: number): number {
  return value * MM_TO_M;
}

export function volumeCenter(volume: ResolvedVolume): [number, number, number] {
  return [
    mm(volume.x + volume.width / 2),
    mm(volume.y + volume.height / 2),
    mm(volume.z + volume.depth / 2),
  ];
}

export function volumeSize(volume: ResolvedVolume): [number, number, number] {
  return [mm(volume.width), mm(volume.height), mm(volume.depth)];
}

/** Caja debug alineada con profundidad estructural (sin el vuelo de la trasera). */
export function volumeStructureDebugBox(
  volume: ResolvedVolume,
  structureDepthMm: number,
): { position: [number, number, number]; size: [number, number, number] } {
  const depth = Math.max(0, structureDepthMm);
  return {
    position: [
      mm(volume.x + volume.width / 2),
      mm(volume.y + volume.height / 2),
      mm(volume.z + depth / 2),
    ],
    size: [mm(volume.width), mm(volume.height), mm(depth)],
  };
}

export function furnitureSceneCenter(furniture: ResolvedFurniture): [number, number, number] {
  const bounds = furnitureBounds(furniture);
  return [
    mm((bounds.min.x + bounds.max.x) / 2),
    mm((bounds.min.y + bounds.max.y) / 2),
    mm((bounds.min.z + bounds.max.z) / 2),
  ];
}

export function furnitureSceneSize(furniture: ResolvedFurniture): number {
  const bounds = furnitureBounds(furniture);
  const dx = bounds.max.x - bounds.min.x;
  const dy = bounds.max.y - bounds.min.y;
  const dz = bounds.max.z - bounds.min.z;
  return mm(Math.max(dx, dy, dz));
}

const VOLUME_PALETTE = [
  '#c8d6e5',
  '#82ccdd',
  '#60a3bc',
  '#3c6382',
  '#0a3d62',
  '#b8e994',
  '#78e08f',
  '#38ada9',
];

export function boxCenter(
  x: number,
  y: number,
  z: number,
  width: number,
  height: number,
  depth: number,
): [number, number, number] {
  return [mm(x + width / 2), mm(y + height / 2), mm(z + depth / 2)];
}

export function boxSize(
  width: number,
  height: number,
  depth: number,
): [number, number, number] {
  return [mm(width), mm(height), mm(depth)];
}

export function volumeColor(depth: number, id: string): string {
  let hash = 0;
  for (let i = 0; i < id.length; i++) {
    hash = (hash + id.charCodeAt(i) * (depth + 1)) % VOLUME_PALETTE.length;
  }
  return VOLUME_PALETTE[hash];
}
