import type { ResolvedFeature, ResolvedVolume } from '@/domain/resolved-furniture';
import { roundMm } from '@/utils/round-mm';

export const NORDEX_THICKNESS_MM = 3;
export const MELAMINE_BACK_THICKNESS_MM = 18;
/** Cuánto entra el nordex en la ranura (ranura 7 mm de profundidad). */
export const NORDEX_GROOVE_INSET_MM = 6;

/** Las ranuras se fresan a lo largo de todo el borde del tablero (limitación de taller),
 *  pero solo dentro del contorno de esa pieza — no se extienden fuera de la melamina. */

export type BackMaterialKind = 'nordex' | 'melamine';

export function isNordexMaterialId(materialId?: string): boolean {
  return materialId === 'nordex';
}

export function backMaterialKind(materialId?: string): BackMaterialKind {
  return isNordexMaterialId(materialId) ? 'nordex' : 'melamine';
}

export function backPanelThicknessMm(materialId?: string): number {
  return isNordexMaterialId(materialId) ? NORDEX_THICKNESS_MM : MELAMINE_BACK_THICKNESS_MM;
}

export function backPanelsUseGrooves(materialId?: string): boolean {
  return isNordexMaterialId(materialId);
}

/** Profundidad estructural de carcasa: nordex en ranura → laterales a todo el volumen. */
export function carcassStructureDepthMm(volumeDepthMm: number, backMaterialId?: string): number {
  if (backPanelsUseGrooves(backMaterialId)) {
    return volumeDepthMm;
  }
  return volumeDepthMm - backPanelThicknessMm(backMaterialId);
}

/** Z del panel trasero nordex centrado en la ranura del borde posterior. */
export function nordexBackPanelZMm(
  volumeZMm: number,
  structureDepthMm: number,
  backThicknessMm: number,
  grooveOffsetFromEdgeMm: number,
  grooveWidthMm: number,
): number {
  const grooveZ = volumeZMm + structureDepthMm - grooveOffsetFromEdgeMm - grooveWidthMm;
  return roundMm(grooveZ + (grooveWidthMm - backThicknessMm) / 2);
}

/** Y del fondo nordex del cajón centrado en la ranura del borde inferior. */
export function nordexBottomPanelYMm(
  bodyBottomYMm: number,
  bottomThicknessMm: number,
  grooveOffsetFromEdgeMm: number,
  grooveWidthMm: number,
): number {
  const grooveY = bodyBottomYMm + grooveOffsetFromEdgeMm;
  return roundMm(grooveY + (grooveWidthMm - bottomThicknessMm) / 2);
}

/** Ancho/largo del nordex = hueco interno + 6 mm por cada ranura donde encaja. */
export function nordexPanelSpanMm(internalSpanMm: number, engagingGrooveCount: number): number {
  return roundMm(internalSpanMm + engagingGrooveCount * NORDEX_GROOVE_INSET_MM);
}

/** Lee backMaterialId o bottomMaterialId del drawer_stack del volumen. */
export function resolveVolumeBackMaterialId(volume: ResolvedVolume): string {
  for (const feature of volume.features) {
    if (feature.type !== 'drawer_stack') {
      continue;
    }
    const fromBack = feature.params.backMaterialId;
    if (typeof fromBack === 'string' && fromBack.length > 0) {
      return fromBack;
    }
    const fromBottom = feature.params.bottomMaterialId;
    if (typeof fromBottom === 'string' && fromBottom.length > 0) {
      return fromBottom;
    }
  }
  return 'melamine-white-18';
}

export function resolveFeatureBackMaterialId(
  feature: ResolvedFeature,
  volume: ResolvedVolume,
): string {
  const fromBack = feature.params.backMaterialId;
  if (typeof fromBack === 'string' && fromBack.length > 0) {
    return fromBack;
  }
  const fromBottom = feature.params.bottomMaterialId;
  if (typeof fromBottom === 'string' && fromBottom.length > 0) {
    return fromBottom;
  }
  return volume.materialId ?? 'melamine-white-18';
}
