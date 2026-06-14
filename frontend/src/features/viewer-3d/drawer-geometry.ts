import type { ResolvedFeature, ResolvedVolume } from '@/domain/resolved-furniture';
import {
  parseDrawerStackConfig,
  snapRunnerLengthMm,
  type DrawerStackConfig,
} from './drawer-config';

export type DrawerPanelRole =
  | 'structure'
  | 'bottom'
  | 'front'
  | 'handle'
  | 'runner'
  | 'groove_rail';

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
  framePanels: DrawerPanelSpec[];
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
  if (panel.width > 0 && panel.height > 0 && panel.depth > 0) {
    panels.push(panel);
  }
}

interface DrawerSlotLayout {
  drawerHeight: number;
  boxX: number;
  boxW: number;
  boxDepth: number;
  innerSideHeight: number;
  falseFrontHeight: number;
  falseFrontY: number;
  runnerLength: number;
  runnerY: number;
  frontLocalX: number;
  frontWidth: number;
}

function layoutDrawerSlot(
  drawerHeight: number,
  frontLocalX: number,
  frontWidth: number,
  boxDepth: number,
  config: DrawerStackConfig,
  sharesDeskLateral: boolean,
): DrawerSlotLayout {
  const bt = config.bottomThicknessMm;
  const inset = config.boxInsetSideMm;
  const insetTop = config.boxInsetTopMm;

  const boxX = frontLocalX + inset;
  const boxW = Math.max(0, frontWidth - inset - (sharesDeskLateral ? inset : inset * 2));
  const innerSideHeight = Math.max(0, drawerHeight - insetTop - bt);
  const falseFrontHeight = innerSideHeight * config.falseFrontHeightRatio;
  const falseFrontY = bt + (innerSideHeight - falseFrontHeight) / 2;
  const runnerLength = snapRunnerLengthMm(boxDepth, config);
  const runnerY = bt + 4;

  return {
    drawerHeight,
    boxX,
    boxW,
    boxDepth,
    innerSideHeight,
    falseFrontHeight,
    falseFrontY,
    runnerLength,
    runnerY,
    frontLocalX,
    frontWidth,
  };
}

function buildDrawerBoxPanels(
  layout: DrawerSlotLayout,
  config: DrawerStackConfig,
  sharesDeskLateral: boolean,
): DrawerPanelSpec[] {
  const panels: DrawerPanelSpec[] = [];
  const st = config.panelThicknessMm;
  const bt = config.bottomThicknessMm;
  const { boxX, boxW, boxDepth, innerSideHeight, falseFrontHeight, falseFrontY, runnerY, runnerLength } =
    layout;

  const backZ = 0;
  const innerW = Math.max(0, boxW - st - (sharesDeskLateral ? 0 : st));
  const gapBeforeOuterFront = 2;
  const sideZ = st;
  const sideDepth = Math.max(0, boxDepth - sideZ - gapBeforeOuterFront);
  const bottomDepth = Math.max(0, sideDepth - st);

  addPanel(panels, {
    x: boxX,
    y: bt,
    z: sideZ,
    width: st,
    height: innerSideHeight,
    depth: sideDepth,
    role: 'structure',
  });

  if (!sharesDeskLateral) {
    addPanel(panels, {
      x: boxX + boxW - st,
      y: bt,
      z: sideZ,
      width: st,
      height: innerSideHeight,
      depth: sideDepth,
      role: 'structure',
    });
  }

  addPanel(panels, {
    x: boxX + st,
    y: bt,
    z: backZ,
    width: innerW,
    height: innerSideHeight,
    depth: st,
    role: 'structure',
  });

  addPanel(panels, {
    x: boxX + st,
    y: 0,
    z: sideZ,
    width: innerW,
    height: bt,
    depth: bottomDepth,
    role: 'bottom',
  });

  const falseFrontZ = sideZ + sideDepth - st;
  addPanel(panels, {
    x: boxX + st,
    y: falseFrontY,
    z: falseFrontZ,
    width: innerW,
    height: falseFrontHeight,
    depth: st,
    role: 'structure',
  });

  addPanel(panels, {
    x: layout.frontLocalX,
    y: 0,
    z: boxDepth,
    width: layout.frontWidth,
    height: layout.drawerHeight,
    depth: config.frontThicknessMm,
    role: 'front',
  });

  addPanel(panels, {
    x: layout.frontLocalX + layout.frontWidth / 2 - 40,
    y: layout.drawerHeight / 2 - 6,
    z: boxDepth + config.frontThicknessMm - 6,
    width: 80,
    height: 12,
    depth: 12,
    role: 'handle',
  });

  addPanel(panels, {
    x: boxX + 1,
    y: runnerY,
    z: sideZ,
    width: config.runnerWidthMm,
    height: config.runnerHeightMm,
    depth: Math.min(runnerLength, sideDepth),
    role: 'runner',
  });

  if (!sharesDeskLateral) {
    addPanel(panels, {
      x: boxX + boxW - st - config.runnerWidthMm - 1,
      y: runnerY,
      z: sideZ,
      width: config.runnerWidthMm,
      height: config.runnerHeightMm,
      depth: Math.min(runnerLength, sideDepth),
      role: 'runner',
    });
  }

  return panels;
}

function buildFrameGrooveAndRunners(
  slotLayouts: DrawerSlotLayout[],
  config: DrawerStackConfig,
  sharesDeskLateral: boolean,
  frontLocalX: number,
  frontWidth: number,
): DrawerPanelSpec[] {
  const panels: DrawerPanelSpec[] = [];
  const grooveX = frontLocalX - config.panelThicknessMm;
  const railX = grooveX + (config.grooveWidthMm - config.grooveRailThicknessMm) / 2;

  for (const layout of slotLayouts) {
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
      width: config.runnerWidthMm,
      height: config.runnerHeightMm,
      depth: layout.runnerLength,
      role: 'runner',
    });

    if (sharesDeskLateral) {
      const rightRunnerX = frontLocalX + frontWidth - config.panelThicknessMm - config.runnerWidthMm - 1;
      addPanel(panels, {
        x: rightRunnerX,
        y: layout.runnerY,
        z: cabinetRunnerZ,
        width: config.runnerWidthMm,
        height: config.runnerHeightMm,
        depth: layout.runnerLength,
        role: 'runner',
      });
    }
  }

  return panels;
}

export function buildDrawerStackGeometry(
  feature: ResolvedFeature,
  volume: ResolvedVolume,
  nestedInDesk: boolean,
): DrawerStackGeometry {
  const config = parseDrawerStackConfig(feature.params);
  const st = config.panelThicknessMm;
  const count = intParam(feature.params, 'count', 1);
  const sharedLateral = stringParam(feature.params, 'sharedLateral');
  const sharesDeskLateral = nestedInDesk && sharedLateral === 'right';

  const drawerHeight = nestedInDesk
    ? numberParam(feature.params, 'drawerHeightMm', 175)
    : Math.max(0, volume.height - 2 * st) / count;
  const gap = nestedInDesk ? 2 : 0;
  const stackTotal = count * drawerHeight + Math.max(0, count - 1) * gap;
  const stackBottomY = nestedInDesk
    ? volume.y + volume.height - stackTotal
    : volume.y + st;

  const originX = volume.x + st;
  const originZ = volume.z + config.backClearanceMm;
  const drawerDepth = Math.max(0, volume.depth - config.backClearanceMm - config.frontGapMm);
  const boxDepth = Math.max(0, drawerDepth - config.frontThicknessMm - 2);

  const frontWidth = sharesDeskLateral ? volume.width - st : volume.width;
  const frontLocalX = volume.x - originX;

  const units: DrawerUnitSpec[] = [];
  const slotLayouts: DrawerSlotLayout[] = [];

  for (let i = 0; i < count; i++) {
    const slotY = nestedInDesk ? stackBottomY + i * (drawerHeight + gap) : stackBottomY + i * drawerHeight;
    const layout = layoutDrawerSlot(
      drawerHeight,
      frontLocalX,
      frontWidth,
      boxDepth,
      config,
      sharesDeskLateral,
    );
    slotLayouts.push(layout);

    units.push({
      id: `${feature.id}-drawer-${i + 1}`,
      worldOrigin: [originX, slotY, originZ],
      slideDistance: boxDepth * 0.5,
      outerHeight: drawerHeight,
      panels: buildDrawerBoxPanels(layout, config, sharesDeskLateral),
    });
  }

  const framePanels = nestedInDesk
    ? buildFrameGrooveAndRunners(slotLayouts, config, sharesDeskLateral, frontLocalX, frontWidth)
    : [];

  return { units, framePanels };
}

export function buildDrawerStackUnits(
  feature: ResolvedFeature,
  volume: ResolvedVolume,
  nestedInDesk: boolean,
): DrawerUnitSpec[] {
  return buildDrawerStackGeometry(feature, volume, nestedInDesk).units;
}
