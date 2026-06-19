import type { ResolvedFeature, ResolvedVolume } from '@/domain/resolved-furniture';
import { roundMm } from '@/utils/round-mm';
import {
  backPanelsUseGrooves,
  nordexBottomPanelYMm,
  nordexPanelSpanMm,
  resolveFeatureBackMaterialId,
} from './back-panel';
import {
  computeCarcassStructureDepth,
  computeDrawerBodyHeights,
  computeDrawerBoxDepth,
  computeDrawerBoxOriginX,
  computeDrawerBoxWidth,
  computeDrawerFrontDimensions,
  resolveDrawerSlotEnclosure,
  type DrawerSlotEnclosure,
  parseDrawerStackConfig,
  snapRunnerLengthMm,
  type DrawerStackConfig,
} from './drawer-config';
import {
  buildCabinetRunnerU,
  buildDrawerRunnerU,
  runnerZoneLeftX,
  runnerZoneRightX,
} from './runner-u-geometry';

/** Convención Z local del cajón: 0 = frente (habitación), +Z = fondo (pared). */

export type DrawerPanelRole =
  | 'structure'
  | 'bottom'
  | 'front'
  | 'front_edge_top'
  | 'front_edge_bottom'
  | 'front_edge_left'
  | 'front_edge_right'
  | 'handle'
  | 'runner'
  | 'runner_cabinet'
  | 'runner_drawer'
  | 'groove_rail'
  | 'groove_cut';

export interface DrawerPanelSpec {
  x: number;
  y: number;
  z: number;
  width: number;
  height: number;
  depth: number;
  role: DrawerPanelRole;
}

export interface DrawerUnitSpec {
  id: string;
  worldOrigin: [number, number, number];
  slideDistance: number;
  outerHeight: number;
  panels: DrawerPanelSpec[];
}

export interface DrawerStackGeometry {
  units: DrawerUnitSpec[];
  /** Correderas de carcasa (U grande), fijas — una entrada por cajón. */
  fixedRunnerSlots: { slotY: number; panels: DrawerPanelSpec[] }[];
}

function intParam(params: Record<string, unknown>, key: string, fallback: number) {
  const value = params[key];
  if (typeof value === 'number' && Number.isFinite(value)) {
    return Math.max(1, Math.round(value));
  }
  return fallback;
}

function numberParam(params: Record<string, unknown>, key: string, fallback: number) {
  const value = params[key];
  if (typeof value === 'number' && Number.isFinite(value)) {
    return value;
  }
  return fallback;
}

function stringParam(params: Record<string, unknown>, key: string): string | undefined {
  const value = params[key];
  return typeof value === 'string' ? value : undefined;
}

function addPanel(panels: DrawerPanelSpec[], panel: DrawerPanelSpec) {
  const p: DrawerPanelSpec = {
    ...panel,
    x: roundMm(panel.x),
    y: roundMm(panel.y),
    z: roundMm(panel.z),
    width: roundMm(panel.width),
    height: roundMm(panel.height),
    depth: roundMm(panel.depth),
  };
  if (p.width > 0 && p.height > 0 && p.depth > 0) {
    panels.push(p);
  }
}

interface DrawerSlotLayout {
  slotY: number;
  externalWidth: number;
  finalFrontHeight: number;
  bodySideHeight: number;
  falseFrontHeight: number;
  falseFrontY: number;
  bodyBottomY: number;
  boxX: number;
  boxW: number;
  boxDepth: number;
  runnerLength: number;
  runnerY: number;
  frontOuterX: number;
  frontOuterY: number;
  frontOuterWidth: number;
  frontOuterHeight: number;
  frontPanelX: number;
  frontPanelY: number;
  frontPanelWidth: number;
  frontPanelHeight: number;
  frontLocalZ: number;
  structureDepth: number;
}

function layoutDrawerSlot(
  slotY: number,
  finalFrontHeight: number,
  externalWidth: number,
  externalHeight: number,
  boxDepth: number,
  structureDepth: number,
  config: DrawerStackConfig,
  sharesDeskLateral: boolean,
  enclosure: DrawerSlotEnclosure,
): DrawerSlotLayout {
  const frontDims = computeDrawerFrontDimensions(externalWidth, externalHeight, config);
  const body = computeDrawerBodyHeights(finalFrontHeight, config, enclosure);

  const boxX = computeDrawerBoxOriginX(config);
  const boxW = computeDrawerBoxWidth(externalWidth, config, sharesDeskLateral);
  const runnerLength = snapRunnerLengthMm(boxDepth, config);
  const runnerY = config.runnerOffsetFromBaseMm;

  return {
    slotY,
    externalWidth,
    finalFrontHeight,
    bodySideHeight: body.bodySideHeight,
    falseFrontHeight: body.falseFrontHeight,
    falseFrontY: body.falseFrontY,
    bodyBottomY: body.bodyBottomY,
    boxX,
    boxW,
    boxDepth,
    runnerLength,
    runnerY,
    frontOuterX: frontDims.outerX,
    frontOuterY: frontDims.outerY,
    frontOuterWidth: frontDims.outerWidth,
    frontOuterHeight: frontDims.outerHeight,
    frontPanelX: frontDims.panelX,
    frontPanelY: frontDims.panelY,
    frontPanelWidth: frontDims.panelWidth,
    frontPanelHeight: frontDims.panelHeight,
    frontLocalZ: -config.frontThicknessMm - config.frontGapMm,
    structureDepth,
  };
}

function buildFrontEdgeBandPanels(
  outerX: number,
  outerY: number,
  outerW: number,
  outerH: number,
  panelY: number,
  panelH: number,
  z: number,
  depth: number,
  band: number,
): DrawerPanelSpec[] {
  if (band <= 0 || outerW <= 0 || outerH <= 0) {
    return [];
  }
  return [
    {
      x: outerX,
      y: outerY + outerH - band,
      z,
      width: outerW,
      height: band,
      depth,
      role: 'front_edge_top',
    },
    {
      x: outerX,
      y: outerY,
      z,
      width: outerW,
      height: band,
      depth,
      role: 'front_edge_bottom',
    },
    {
      x: outerX,
      y: panelY,
      z,
      width: band,
      height: panelH,
      depth,
      role: 'front_edge_left',
    },
    {
      x: outerX + outerW - band,
      y: panelY,
      z,
      width: band,
      height: panelH,
      depth,
      role: 'front_edge_right',
    },
  ];
}

function buildDrawerBottomGrooves(
  layout: Pick<DrawerSlotLayout, 'boxX' | 'boxW' | 'bodyBottomY' | 'boxDepth'>,
  innerW: number,
  config: DrawerStackConfig,
  sharesDeskLateral: boolean,
): DrawerPanelSpec[] {
  const st = config.panelThicknessMm;
  const go = config.grooveOffsetFromEdgeMm;
  const gw = config.grooveWidthMm;
  const gd = config.grooveDepthMm;
  const grooveY = layout.bodyBottomY + go;
  const panelDepth = layout.boxDepth;
  const panels: DrawerPanelSpec[] = [];

  // Ranura = todo el borde del tablero de esa pieza (no fuera de la melamina).
  addPanel(panels, {
    x: layout.boxX + st - gd,
    y: grooveY,
    z: 0,
    width: gd,
    height: gw,
    depth: panelDepth,
    role: 'groove_cut',
  });

  if (!sharesDeskLateral) {
    addPanel(panels, {
      x: layout.boxX + layout.boxW - st,
      y: grooveY,
      z: 0,
      width: gd,
      height: gw,
      depth: panelDepth,
      role: 'groove_cut',
    });
  }

  addPanel(panels, {
    x: layout.boxX + st,
    y: grooveY,
    z: st - gd,
    width: innerW,
    height: gw,
    depth: gd,
    role: 'groove_cut',
  });

  addPanel(panels, {
    x: layout.boxX + st,
    y: grooveY,
    z: layout.boxDepth - st,
    width: innerW,
    height: gw,
    depth: gd,
    role: 'groove_cut',
  });

  return panels;
}

function buildDrawerBoxPanels(
  layout: DrawerSlotLayout,
  config: DrawerStackConfig,
  sharesDeskLateral: boolean,
  backMaterialId: string,
): DrawerPanelSpec[] {
  const panels: DrawerPanelSpec[] = [];
  const st = config.panelThicknessMm;
  const bt = config.bottomThicknessMm;
  const stackW = config.lateralInsetMm;
  const {
    bodySideHeight,
    falseFrontHeight,
    bodyBottomY,
    boxX,
    boxW,
    boxDepth,
    runnerY,
    runnerLength,
    frontOuterX,
    frontOuterY,
    frontOuterWidth,
    frontOuterHeight,
    frontPanelX,
    frontPanelY,
    frontPanelWidth,
    frontPanelHeight,
    frontLocalZ,
  } = layout;

  const backZ = boxDepth - st;
  const innerW = Math.max(0, boxW - (sharesDeskLateral ? st : 2 * st));
  const innerD = Math.max(0, backZ - st);
  const sideDepth = boxDepth;
  const useGrooves = backPanelsUseGrooves(backMaterialId);
  const grooveInset = config.nordexGrooveInsetMm;
  const lateralGrooveCount = sharesDeskLateral ? 1 : 2;
  const bottomDepth = useGrooves ? nordexPanelSpanMm(innerD, 2) : Math.max(0, sideDepth - st);
  const bottomY = useGrooves
    ? nordexBottomPanelYMm(bodyBottomY, bt, config.grooveOffsetFromEdgeMm, config.grooveWidthMm)
    : bodyBottomY;
  const bottomPanelWidth = useGrooves
    ? nordexPanelSpanMm(innerW, lateralGrooveCount)
    : innerW;
  const bottomPanelX = useGrooves
    ? boxX + st - grooveInset
    : boxX + st;
  const bottomPanelZ = useGrooves ? st - grooveInset : 0;

  addPanel(panels, {
    x: boxX,
    y: bodyBottomY,
    z: 0,
    width: st,
    height: bodySideHeight,
    depth: sideDepth,
    role: 'structure',
  });

  if (!sharesDeskLateral) {
    addPanel(panels, {
      x: boxX + boxW - st,
      y: bodyBottomY,
      z: 0,
      width: st,
      height: bodySideHeight,
      depth: sideDepth,
      role: 'structure',
    });
  }

  addPanel(panels, {
    x: boxX + st,
    y: bodyBottomY,
    z: backZ,
    width: innerW,
    height: bodySideHeight,
    depth: st,
    role: 'structure',
  });

  addPanel(panels, {
    x: bottomPanelX,
    y: bottomY,
    z: bottomPanelZ,
    width: bottomPanelWidth,
    height: bt,
    depth: bottomDepth,
    role: 'bottom',
  });

  if (useGrooves) {
    panels.push(
      ...buildDrawerBottomGrooves(layout, innerW, config, sharesDeskLateral),
    );
  }

  addPanel(panels, {
    x: boxX + st,
    y: bodyBottomY,
    z: 0,
    width: innerW,
    height: falseFrontHeight,
    depth: st,
    role: 'structure',
  });

  // Frente final: fuera del volumen, altura = extremo a extremo del volumen externo.
  addPanel(panels, {
    x: frontPanelX,
    y: frontPanelY,
    z: frontLocalZ,
    width: frontPanelWidth,
    height: frontPanelHeight,
    depth: config.frontThicknessMm,
    role: 'front',
  });

  panels.push(
    ...buildFrontEdgeBandPanels(
      frontOuterX,
      frontOuterY,
      frontOuterWidth,
      frontOuterHeight,
      frontPanelY,
      frontPanelHeight,
      frontLocalZ,
      config.frontThicknessMm,
      config.thickEdgeBandingMm,
    ),
  );

  addPanel(panels, {
    x: frontOuterX + frontOuterWidth / 2 - 40,
    y: frontPanelY + frontPanelHeight / 2 - 6,
    z: frontLocalZ - 12,
    width: 80,
    height: 12,
    depth: 12,
    role: 'handle',
  });

  const leftRunnerX = runnerZoneLeftX(boxX, stackW);
  const rightRunnerX = runnerZoneRightX(boxX + boxW);
  const runnerZ = 0;
  const runnerDepth = runnerLength;

  panels.push(
    ...buildDrawerRunnerU(
      'left',
      leftRunnerX,
      stackW,
      runnerY,
      config.runnerHeightMm,
      runnerDepth,
      runnerZ,
    ),
  );

  if (!sharesDeskLateral) {
    panels.push(
      ...buildDrawerRunnerU(
        'right',
        rightRunnerX,
        stackW,
        runnerY,
        config.runnerHeightMm,
        runnerDepth,
        runnerZ,
      ),
    );
  }

  return panels;
}

function buildFixedCabinetRunners(
  layout: DrawerSlotLayout,
  config: DrawerStackConfig,
  sharesDeskLateral: boolean,
): DrawerPanelSpec[] {
  const panels: DrawerPanelSpec[] = [];
  const stackW = config.lateralInsetMm;
  const runnerZ = 0;
  const leftRunnerX = runnerZoneLeftX(layout.boxX, stackW);
  const rightRunnerX = runnerZoneRightX(layout.boxX + layout.boxW);

  panels.push(
    ...buildCabinetRunnerU(
      'left',
      leftRunnerX,
      stackW,
      layout.runnerY,
      config.runnerHeightMm,
      layout.runnerLength,
      runnerZ,
    ),
  );

  if (!sharesDeskLateral) {
    panels.push(
      ...buildCabinetRunnerU(
        'right',
        rightRunnerX,
        stackW,
        layout.runnerY,
        config.runnerHeightMm,
        layout.runnerLength,
        runnerZ,
      ),
    );
  }

  return panels;
}

function buildDeskGrooveRails(
  layout: DrawerSlotLayout,
  config: DrawerStackConfig,
  sharesDeskLateral: boolean,
): DrawerPanelSpec[] {
  const panels: DrawerPanelSpec[] = [];
  const st = config.panelThicknessMm;
  const stackW = config.lateralInsetMm;
  const grooveX = layout.frontOuterX - st;
  const channelW = config.deskGrooveChannelMm;
  const railX = grooveX + (channelW - config.grooveRailThicknessMm) / 2;
  const grooveZ = layout.boxDepth - config.grooveDepthMm;

  addPanel(panels, {
    x: railX,
    y: layout.runnerY,
    z: grooveZ,
    width: config.grooveRailThicknessMm,
    height: config.runnerHeightMm,
    depth: config.grooveDepthMm,
    role: 'groove_rail',
  });

  const cabinetRunnerZ = grooveZ + config.grooveRailThicknessMm;
  addPanel(panels, {
    x: grooveX + 1,
    y: layout.runnerY,
    z: cabinetRunnerZ,
    width: stackW,
    height: config.runnerHeightMm,
    depth: layout.runnerLength,
    role: 'runner_cabinet',
  });

  if (sharesDeskLateral) {
    const rightRunnerX = layout.frontOuterX + layout.frontOuterWidth - st - stackW - 1;
    addPanel(panels, {
      x: rightRunnerX,
      y: layout.runnerY,
      z: cabinetRunnerZ,
      width: stackW,
      height: config.runnerHeightMm,
      depth: layout.runnerLength,
      role: 'runner_cabinet',
    });
  }

  return panels;
}

export function buildDrawerStackGeometry(
  feature: ResolvedFeature,
  volume: ResolvedVolume,
  nestedInDesk: boolean,
): DrawerStackGeometry {
  const params = feature.params;
  const config = parseDrawerStackConfig(params);
  const count = intParam(params, 'count', 1);
  const sharedLateral = stringParam(params, 'sharedLateral');
  const sharesDeskLateral = nestedInDesk && sharedLateral === 'right';
  const enclosure = resolveDrawerSlotEnclosure(params, nestedInDesk, count);

  const backMaterialId = resolveFeatureBackMaterialId(feature, volume);
  const structureDepth = computeCarcassStructureDepth(volume.depth, config, params, backMaterialId);
  const boxDepth = computeDrawerBoxDepth(volume.depth, config, params);

  const slotHeight = nestedInDesk
    ? numberParam(params, 'drawerHeightMm', 175)
    : volume.height / count;
  const gap = nestedInDesk ? 2 : 0;
  const stackTotal = count * slotHeight + Math.max(0, count - 1) * gap;
  const stackBottomY = nestedInDesk
    ? volume.height - stackTotal
    : 0;

  const originX = volume.x;
  const originY = volume.y;
  const originZ = volume.z;

  const units: DrawerUnitSpec[] = [];
  const fixedRunnerSlots: { slotY: number; panels: DrawerPanelSpec[] }[] = [];

  for (let i = 0; i < count; i++) {
    const slotY = stackBottomY + i * (slotHeight + gap);
    const finalFrontHeight = nestedInDesk ? slotHeight : volume.height / count;

    const layout = layoutDrawerSlot(
      slotY,
      finalFrontHeight,
      volume.width,
      volume.height,
      boxDepth,
      structureDepth,
      config,
      sharesDeskLateral,
      enclosure,
    );

    const unitPanels = [...buildDrawerBoxPanels(layout, config, sharesDeskLateral, backMaterialId)];
    if (nestedInDesk) {
      unitPanels.push(...buildDeskGrooveRails(layout, config, sharesDeskLateral));
    }

    fixedRunnerSlots.push({
      slotY,
      panels: buildFixedCabinetRunners(layout, config, sharesDeskLateral),
    });

    units.push({
      id: `${feature.id}-drawer-${i + 1}`,
      worldOrigin: [originX, originY + slotY, originZ],
      slideDistance: layout.runnerLength,
      outerHeight: finalFrontHeight,
      panels: unitPanels,
    });
  }

  return { units, fixedRunnerSlots };
}

export function buildDrawerStackUnits(
  feature: ResolvedFeature,
  volume: ResolvedVolume,
  nestedInDesk: boolean,
): DrawerUnitSpec[] {
  return buildDrawerStackGeometry(feature, volume, nestedInDesk).units;
}
