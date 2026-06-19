import { THICK_EDGE_BANDING_MM } from './materials';
import { backPanelsUseGrooves, NORDEX_THICKNESS_MM } from './back-panel';
import { roundMm } from '@/utils/round-mm';

/** Espesor trasera de carcasa (melamina 18 mm en cajonera estándar). */
export const CARCASS_BACK_THICKNESS_MM = 18;

/** Configurable drawer / runner / groove defaults (override via feature.params). */
export interface DrawerStackConfig {
  panelThicknessMm: number;
  bottomThicknessMm: number;
  frontThicknessMm: number;
  carcassBackThicknessMm: number;
  /** Holgura cajón respecto al fondo interior de carcasa (382 → 372 mm). */
  drawerRearClearanceMm: number;
  backClearanceMm: number;
  backClearanceRatio: number;
  frontGapMm: number;
  lateralInsetMm: number;
  boxInsetSideMm: number;
  boxInsetTopMm: number;
  /** Ranura para fondo nordex: inicio a 18 mm del borde, ancho 4 mm (18→22), profundidad 7 mm. */
  grooveOffsetFromEdgeMm: number;
  grooveWidthMm: number;
  grooveDepthMm: number;
  /** Canal ranura escritorio (U corredera), distinto de la ranura de panel. */
  deskGrooveChannelMm: number;
  /** Cuánto entra el nordex en la ranura (7 mm prof. → 6 mm de penetración). */
  nordexGrooveInsetMm: number;
  grooveRailThicknessMm: number;
  runnerHeightMm: number;
  runnerWidthMm: number;
  runnerLengthStepMm: number;
  runnerLengthMinMm: number;
  runnerLengthMaxMm: number;
  /** Corredera a 5 cm de la base del cajón. */
  runnerOffsetFromBaseMm: number;
  thickEdgeBandingMm: number;
  /** Tolerancia estética por lado en frente (1 mm). */
  frontToleranceInsetMm: number;
  /** @deprecated use bodySideHeightRatio + bodyHeightInsetMm */
  bodyHeightInsetMm: number;
  /** @deprecated use falseFrontHeightRatio */
  falseFrontHeightInsetMm: number;
  /** Altura laterales = min(ratio × tapa, tapa − inset mm). */
  bodySideHeightRatio: number;
  /** Altura frente falso = ratio × altura laterales. */
  falseFrontHeightRatio: number;
  /** Holgura bajo el techo interior del volumen (mm). */
  bodyTopClearanceMm: number;
  /** Holgura extra sobre el piso de carcasa (mm). */
  bodyBottomClearanceMm: number;
  /** @deprecated use frontToleranceInsetMm × 2 */
  frontOuterWidthInsetMm: number;
}

export const DEFAULT_DRAWER_STACK_CONFIG: DrawerStackConfig = {
  panelThicknessMm: 18,
  bottomThicknessMm: 3,
  frontThicknessMm: 18,
  carcassBackThicknessMm: CARCASS_BACK_THICKNESS_MM,
  drawerRearClearanceMm: 10,
  backClearanceMm: 10,
  backClearanceRatio: 0,
  frontGapMm: 0,
  lateralInsetMm: 13.5,
  boxInsetSideMm: 2,
  boxInsetTopMm: 2,
  grooveOffsetFromEdgeMm: 18,
  grooveWidthMm: 4,
  grooveDepthMm: 7,
  nordexGrooveInsetMm: 6,
  deskGrooveChannelMm: 18,
  grooveRailThicknessMm: 4,
  runnerHeightMm: 40,
  runnerWidthMm: 13.5,
  runnerLengthStepMm: 50,
  runnerLengthMinMm: 250,
  runnerLengthMaxMm: 650,
  runnerOffsetFromBaseMm: 50,
  thickEdgeBandingMm: THICK_EDGE_BANDING_MM,
  frontToleranceInsetMm: 1,
  bodyHeightInsetMm: 20,
  falseFrontHeightInsetMm: 20,
  bodySideHeightRatio: 0.85,
  falseFrontHeightRatio: 0.85,
  bodyTopClearanceMm: 24,
  bodyBottomClearanceMm: 4,
  frontOuterWidthInsetMm: 2,
};

/** Longitudes estándar de corredera (mm): 25, 30, 35 … 65 cm. */
export const RUNNER_LENGTH_OPTIONS_MM = [250, 300, 350, 400, 450, 500, 550, 600, 650] as const;

function numberParam(params: Record<string, unknown>, key: string, fallback: number) {
  const value = params[key];
  if (typeof value === 'number' && Number.isFinite(value)) {
    return value;
  }
  return fallback;
}

export function parseDrawerStackConfig(params: Record<string, unknown>): DrawerStackConfig {
  const d = DEFAULT_DRAWER_STACK_CONFIG;
  const lateralFallback = numberParam(params, 'runnerStackWidthMm', d.lateralInsetMm);
  return {
    panelThicknessMm: numberParam(params, 'panelThicknessMm', d.panelThicknessMm),
    bottomThicknessMm: numberParam(params, 'bottomThicknessMm', d.bottomThicknessMm),
    frontThicknessMm: numberParam(params, 'frontThicknessMm', d.frontThicknessMm),
    carcassBackThicknessMm: numberParam(params, 'carcassBackThicknessMm', d.carcassBackThicknessMm),
    drawerRearClearanceMm: numberParam(params, 'drawerRearClearanceMm', d.drawerRearClearanceMm),
    backClearanceMm: numberParam(params, 'backClearanceMm', d.backClearanceMm),
    backClearanceRatio: numberParam(params, 'backClearanceRatio', d.backClearanceRatio),
    frontGapMm: numberParam(params, 'frontGapMm', d.frontGapMm),
    lateralInsetMm: numberParam(params, 'lateralInsetMm', lateralFallback),
    boxInsetSideMm: numberParam(params, 'boxInsetSideMm', d.boxInsetSideMm),
    boxInsetTopMm: numberParam(params, 'boxInsetTopMm', d.boxInsetTopMm),
    grooveOffsetFromEdgeMm: numberParam(params, 'grooveOffsetFromEdgeMm', d.grooveOffsetFromEdgeMm),
    grooveWidthMm: numberParam(params, 'grooveWidthMm', d.grooveWidthMm),
    grooveDepthMm: numberParam(params, 'grooveDepthMm', d.grooveDepthMm),
    nordexGrooveInsetMm: numberParam(params, 'nordexGrooveInsetMm', d.nordexGrooveInsetMm),
    deskGrooveChannelMm: numberParam(params, 'deskGrooveChannelMm', d.deskGrooveChannelMm),
    grooveRailThicknessMm: numberParam(params, 'grooveRailThicknessMm', d.grooveRailThicknessMm),
    runnerHeightMm: numberParam(params, 'runnerHeightMm', d.runnerHeightMm),
    runnerWidthMm: numberParam(params, 'runnerWidthMm', d.runnerWidthMm),
    runnerLengthStepMm: numberParam(params, 'runnerLengthStepMm', d.runnerLengthStepMm),
    runnerLengthMinMm: numberParam(params, 'runnerLengthMinMm', d.runnerLengthMinMm),
    runnerLengthMaxMm: numberParam(params, 'runnerLengthMaxMm', d.runnerLengthMaxMm),
    runnerOffsetFromBaseMm: numberParam(params, 'runnerOffsetFromBaseMm', d.runnerOffsetFromBaseMm),
    thickEdgeBandingMm: numberParam(params, 'thickEdgeBandingMm', d.thickEdgeBandingMm),
    frontToleranceInsetMm: numberParam(params, 'frontToleranceInsetMm', d.frontToleranceInsetMm),
    bodyHeightInsetMm: numberParam(params, 'bodyHeightInsetMm', d.bodyHeightInsetMm),
    falseFrontHeightInsetMm: numberParam(params, 'falseFrontHeightInsetMm', d.falseFrontHeightInsetMm),
    bodySideHeightRatio: numberParam(params, 'bodySideHeightRatio', d.bodySideHeightRatio),
    falseFrontHeightRatio: numberParam(params, 'falseFrontHeightRatio', d.falseFrontHeightRatio),
    bodyTopClearanceMm: numberParam(params, 'bodyTopClearanceMm', d.bodyTopClearanceMm),
    bodyBottomClearanceMm: numberParam(params, 'bodyBottomClearanceMm', d.bodyBottomClearanceMm),
    frontOuterWidthInsetMm: numberParam(params, 'frontOuterWidthInsetMm', d.frontOuterWidthInsetMm),
  };
}

/** Profundidad útil de carcasa. Nordex en ranura: laterales a profundidad total del volumen. */
export function computeCarcassStructureDepth(
  volumeDepthMm: number,
  config: DrawerStackConfig,
  params: Record<string, unknown>,
  backMaterialId?: string,
): number {
  const backFromParams = params.carcassBackMaterialId ?? params.backMaterialId ?? params.bottomMaterialId;
  const backId = typeof backFromParams === 'string' ? backFromParams : backMaterialId;
  if (typeof backId === 'string' && backPanelsUseGrooves(backId)) {
    return volumeDepthMm;
  }
  const backT =
    typeof backId === 'string' && backId === 'nordex'
      ? NORDEX_THICKNESS_MM
      : numberParam(params, 'carcassBackThicknessMm', config.carcassBackThicknessMm);
  return Math.max(0, volumeDepthMm - backT);
}

/** Profundidad del cajón = interior carcasa − 10 mm (382 − 10 = 372 mm). */
export function computeDrawerBoxDepth(
  volumeDepthMm: number,
  config: DrawerStackConfig,
  params: Record<string, unknown>,
): number {
  const structure = computeCarcassStructureDepth(volumeDepthMm, config, params);
  const rear = numberParam(params, 'drawerRearClearanceMm', config.drawerRearClearanceMm);
  return Math.max(0, structure - rear);
}

export interface DrawerFrontDimensions {
  outerWidth: number;
  outerHeight: number;
  panelWidth: number;
  panelHeight: number;
  outerX: number;
  outerY: number;
  panelX: number;
  panelY: number;
}

/** Frente: tolerancia 1 mm/lado + canto grueso 3 mm/lado → tablero 492×342 en volumen 500×350. */
export function computeDrawerFrontDimensions(
  externalWidthMm: number,
  externalHeightMm: number,
  config: DrawerStackConfig,
): DrawerFrontDimensions {
  const tol = config.frontToleranceInsetMm;
  const thick = config.thickEdgeBandingMm;
  const outerWidth = Math.max(0, externalWidthMm - 2 * tol);
  const outerHeight = Math.max(0, externalHeightMm - 2 * tol);
  const panelWidth = Math.max(0, outerWidth - 2 * thick);
  const panelHeight = Math.max(0, outerHeight - 2 * thick);
  return {
    outerWidth,
    outerHeight,
    panelWidth,
    panelHeight,
    outerX: tol,
    outerY: tol,
    panelX: tol + thick,
    panelY: tol + thick,
  };
}

export interface DrawerBodyHeights {
  bodySideHeight: number;
  falseFrontHeight: number;
  bodyBottomY: number;
  falseFrontY: number;
}

export interface DrawerInteriorBounds {
  floorY: number;
  ceilingY: number;
}

/** Interior útil del slot: sobre piso estructural y bajo techo estructural. */
export function computeDrawerInteriorBounds(
  slotHeight: number,
  config: DrawerStackConfig,
): DrawerInteriorBounds {
  const st = config.panelThicknessMm;
  return {
    floorY: st,
    ceilingY: slotHeight - st,
  };
}

/** Si el slot está entre piso y techo de carcasa, restar 2× espesor (36 mm usual). */
export interface DrawerSlotEnclosure {
  subtractCarcassFloorCeiling: boolean;
}

export function resolveDrawerSlotEnclosure(
  params: Record<string, unknown>,
  nestedInDesk: boolean,
  drawerCount: number,
): DrawerSlotEnclosure {
  const explicit = params.carcassFloorCeiling;
  if (typeof explicit === 'boolean') {
    return { subtractCarcassFloorCeiling: explicit };
  }
  // Cajonera cerrada (demo): piso + techo de carcasa. Cajones continuos apilados: solo 85 % del slot.
  return { subtractCarcassFloorCeiling: !nestedInDesk && drawerCount === 1 };
}

export function computeDrawerInteriorAvailableHeight(
  slotHeight: number,
  config: DrawerStackConfig,
  enclosure: DrawerSlotEnclosure,
): number {
  const carcassInset = enclosure.subtractCarcassFloorCeiling ? 2 * config.panelThicknessMm : 0;
  return Math.max(0, slotHeight - carcassInset);
}

/** Laterales = 85 % del alto interior disponible; frente falso = 85 % de laterales. */
export function computeDrawerBodyHeights(
  slotHeight: number,
  config: DrawerStackConfig,
  enclosure: DrawerSlotEnclosure,
): DrawerBodyHeights {
  const available = computeDrawerInteriorAvailableHeight(slotHeight, config, enclosure);
  let bodySideHeight = available * config.bodySideHeightRatio;

  if (enclosure.subtractCarcassFloorCeiling) {
    const interior = computeDrawerInteriorBounds(slotHeight, config);
    const floorY = config.panelThicknessMm + config.bodyBottomClearanceMm;
    const maxTop = interior.ceilingY - config.bodyTopClearanceMm;
    const maxHeight = Math.max(0, maxTop - floorY);
    bodySideHeight = Math.min(bodySideHeight, maxHeight);
  }

  bodySideHeight = roundMm(Math.max(0, bodySideHeight));
  const falseFrontHeight = roundMm(Math.max(0, bodySideHeight * config.falseFrontHeightRatio));
  const bodyBottomY = roundMm(
    enclosure.subtractCarcassFloorCeiling
      ? config.panelThicknessMm + config.bodyBottomClearanceMm
      : slotHeight - bodySideHeight,
  );
  const falseFrontY = roundMm(bodyBottomY + (bodySideHeight - falseFrontHeight) / 2);

  return { bodySideHeight, falseFrontHeight, bodyBottomY, falseFrontY };
}

export function computeDrawerBoxOriginX(config: DrawerStackConfig): number {
  return roundMm(config.panelThicknessMm + config.lateralInsetMm);
}

export function computeDrawerBoxWidth(
  externalWidthMm: number,
  config: DrawerStackConfig,
  sharesDeskLateral: boolean,
): number {
  const st = config.panelThicknessMm;
  const inset = config.lateralInsetMm;
  const boxX = computeDrawerBoxOriginX(config);
  if (sharesDeskLateral) {
    return roundMm(Math.max(0, externalWidthMm - st - boxX));
  }
  const innerW = externalWidthMm - 2 * st;
  return roundMm(Math.max(0, innerW - 2 * inset));
}

export function snapRunnerLengthMm(availableDepthMm: number, config: DrawerStackConfig): number {
  const { runnerLengthMinMm: min, runnerLengthMaxMm: max, runnerLengthStepMm: step } = config;
  let best = min;
  for (let len = min; len <= max; len += step) {
    if (len <= availableDepthMm) {
      best = len;
    } else {
      break;
    }
  }
  return best;
}
