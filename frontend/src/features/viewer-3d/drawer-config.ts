/** Configurable drawer / runner / groove defaults (override via feature.params). */
export interface DrawerStackConfig {
  panelThicknessMm: number;
  bottomThicknessMm: number;
  frontThicknessMm: number;
  backClearanceMm: number;
  frontGapMm: number;
  lateralGapMm: number;
  boxInsetSideMm: number;
  boxInsetTopMm: number;
  grooveWidthMm: number;
  grooveDepthMm: number;
  grooveRailThicknessMm: number;
  runnerHeightMm: number;
  runnerWidthMm: number;
  runnerLengthStepMm: number;
  runnerLengthMinMm: number;
  falseFrontHeightRatio: number;
}

export const DEFAULT_DRAWER_STACK_CONFIG: DrawerStackConfig = {
  panelThicknessMm: 18,
  bottomThicknessMm: 3,
  frontThicknessMm: 18,
  backClearanceMm: 40,
  frontGapMm: 8,
  lateralGapMm: 3,
  boxInsetSideMm: 2,
  boxInsetTopMm: 2,
  grooveWidthMm: 18,
  grooveDepthMm: 7,
  grooveRailThicknessMm: 4,
  runnerHeightMm: 40,
  runnerWidthMm: 8,
  runnerLengthStepMm: 50,
  runnerLengthMinMm: 200,
  falseFrontHeightRatio: 0.8,
};

function numberParam(params: Record<string, unknown>, key: string, fallback: number) {
  const value = params[key];
  if (typeof value === 'number' && Number.isFinite(value)) {
    return value;
  }
  return fallback;
}

export function parseDrawerStackConfig(params: Record<string, unknown>): DrawerStackConfig {
  const d = DEFAULT_DRAWER_STACK_CONFIG;
  return {
    panelThicknessMm: numberParam(params, 'panelThicknessMm', d.panelThicknessMm),
    bottomThicknessMm: numberParam(params, 'bottomThicknessMm', d.bottomThicknessMm),
    frontThicknessMm: numberParam(params, 'frontThicknessMm', d.frontThicknessMm),
    backClearanceMm: numberParam(params, 'backClearanceMm', d.backClearanceMm),
    frontGapMm: numberParam(params, 'frontGapMm', d.frontGapMm),
    lateralGapMm: numberParam(params, 'lateralGapMm', d.lateralGapMm),
    boxInsetSideMm: numberParam(params, 'boxInsetSideMm', d.boxInsetSideMm),
    boxInsetTopMm: numberParam(params, 'boxInsetTopMm', d.boxInsetTopMm),
    grooveWidthMm: numberParam(params, 'grooveWidthMm', d.grooveWidthMm),
    grooveDepthMm: numberParam(params, 'grooveDepthMm', d.grooveDepthMm),
    grooveRailThicknessMm: numberParam(params, 'grooveRailThicknessMm', d.grooveRailThicknessMm),
    runnerHeightMm: numberParam(params, 'runnerHeightMm', d.runnerHeightMm),
    runnerWidthMm: numberParam(params, 'runnerWidthMm', d.runnerWidthMm),
    runnerLengthStepMm: numberParam(params, 'runnerLengthStepMm', d.runnerLengthStepMm),
    runnerLengthMinMm: numberParam(params, 'runnerLengthMinMm', d.runnerLengthMinMm),
    falseFrontHeightRatio: numberParam(params, 'falseFrontHeightRatio', d.falseFrontHeightRatio),
  };
}

/** Runner length snapped to step (5 cm) from minimum 20 cm. */
export function snapRunnerLengthMm(depthMm: number, config: DrawerStackConfig): number {
  const raw = Math.max(config.runnerLengthMinMm, depthMm);
  const steps = Math.ceil((raw - config.runnerLengthMinMm) / config.runnerLengthStepMm);
  return config.runnerLengthMinMm + steps * config.runnerLengthStepMm;
}
