import type { ResolvedFeature, ResolvedFront, ResolvedVolume } from '@/domain/resolved-furniture';
import {
  backPanelThicknessMm,
  backPanelsUseGrooves,
  carcassStructureDepthMm,
  nordexBackPanelZMm,
  nordexPanelSpanMm,
  NORDEX_GROOVE_INSET_MM,
  resolveVolumeBackMaterialId,
} from './back-panel';
import { DEFAULT_DRAWER_STACK_CONFIG } from './drawer-config';
import {
  FRONT_THICKNESS_MM,
  GROOVE,
  MELAMINE,
  METAL,
  PANEL_THICKNESS_MM,
  THICK_EDGE_BANDING_MM,
  resolveMaterialColor,
  resolvePanelMaterial,
  type PanelMaterialProps,
} from './materials';

export interface PanelSpec {
  x: number;
  y: number;
  z: number;
  width: number;
  height: number;
  depth: number;
  material: PanelMaterialProps;
  label?: string;
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

export function innerWidth(volume: ResolvedVolume) {
  return Math.max(0, volume.width - 2 * PANEL_THICKNESS_MM);
}

export function innerHeight(volume: ResolvedVolume) {
  return Math.max(0, volume.height - 2 * PANEL_THICKNESS_MM);
}

export function innerDepth(volume: ResolvedVolume, backMaterialId?: string) {
  const backT = backPanelThicknessMm(backMaterialId ?? resolveVolumeBackMaterialId(volume));
  return Math.max(0, volume.depth - backT);
}

export function hasDeskFrame(volume: ResolvedVolume): boolean {
  return volume.features.some((feature) => feature.type === 'desk_frame');
}

export function isNestedDrawerTower(volume: ResolvedVolume, parent?: ResolvedVolume): boolean {
  if (!parent || !hasDeskFrame(parent)) {
    return false;
  }
  return volume.features.some((feature) => feature.type === 'drawer_stack');
}

function boolParam(params: Record<string, unknown>, key: string, fallback: boolean) {
  const value = params[key];
  return typeof value === 'boolean' ? value : fallback;
}

function drawerStackFeature(volume: ResolvedVolume): ResolvedFeature | undefined {
  return volume.features.find((feature) => feature.type === 'drawer_stack');
}

export function nestedDrawerTowerPanels(volume: ResolvedVolume): PanelSpec[] {
  const feature = drawerStackFeature(volume);
  if (!feature || !boolParam(feature.params, 'hasBase', true)) {
    return [];
  }

  const t = PANEL_THICKNESS_MM;
  const material = {
    ...MELAMINE,
    color: resolveMaterialColor(volume.materialId),
  };
  const innerW = Math.max(0, volume.width - 2 * t);

  return [
    {
      x: volume.x + t,
      y: volume.y,
      z: volume.z,
      width: innerW,
      height: t,
      depth: volume.depth,
      material,
    },
  ];
}

export function deskFramePanels(feature: ResolvedFeature, volume: ResolvedVolume): PanelSpec[] {
  const t = PANEL_THICKNESS_MM;
  const material = {
    ...MELAMINE,
    color: resolveMaterialColor(volume.materialId),
  };
  const overhang = numberParam(feature.params, 'topOverhangMm', 25);
  const braceRatio = numberParam(feature.params, 'braceHeightRatio', 0.5);

  const legHeight = volume.height;
  const braceHeight = legHeight * braceRatio;
  const braceY = volume.y + legHeight - braceHeight;

  const leftLateral: PanelSpec = {
    x: volume.x,
    y: volume.y,
    z: volume.z,
    width: t,
    height: legHeight,
    depth: volume.depth,
    material,
  };

  const rightLateral: PanelSpec = {
    x: volume.x + volume.width - t,
    y: volume.y,
    z: volume.z,
    width: t,
    height: legHeight,
    depth: volume.depth,
    material,
  };

  const backBrace: PanelSpec = {
    x: volume.x + t,
    y: braceY,
    z: volume.z,
    width: volume.width - 2 * t,
    height: braceHeight,
    depth: t,
    material,
  };

  const desktop: PanelSpec = {
    x: volume.x - overhang,
    y: volume.y + legHeight,
    z: volume.z,
    width: volume.width + 2 * overhang,
    height: t,
    depth: volume.depth + overhang,
    material,
  };

  return [leftLateral, rightLateral, backBrace, desktop];
}

export function outerCarcassPanels(volume: ResolvedVolume, backMaterialId?: string): PanelSpec[] {
  const t = PANEL_THICKNESS_MM;
  const backId = backMaterialId ?? resolveVolumeBackMaterialId(volume);
  const backT = backPanelThicknessMm(backId);
  const useGrooves = backPanelsUseGrooves(backId);
  const groove = DEFAULT_DRAWER_STACK_CONFIG;
  const structureDepth = carcassStructureDepthMm(volume.depth, backId);
  const innerW = Math.max(0, volume.width - 2 * t);
  const innerH = Math.max(0, volume.height - 2 * t);
  const grooveZ = volume.z + structureDepth - groove.grooveOffsetFromEdgeMm - groove.grooveWidthMm;
  const backPanelZ = useGrooves
    ? nordexBackPanelZMm(
        volume.z,
        structureDepth,
        backT,
        groove.grooveOffsetFromEdgeMm,
        groove.grooveWidthMm,
      )
    : volume.z + structureDepth - backT;
  const material = {
    ...MELAMINE,
    color: resolveMaterialColor(volume.materialId),
  };
  const backMaterial = resolvePanelMaterial(backId === 'nordex' ? 'nordex' : volume.materialId);

  const panels: PanelSpec[] = [
    {
      x: volume.x,
      y: volume.y,
      z: volume.z,
      width: t,
      height: volume.height,
      depth: structureDepth,
      material,
      label: 'Lateral izquierdo',
    },
    {
      x: volume.x + volume.width - t,
      y: volume.y,
      z: volume.z,
      width: t,
      height: volume.height,
      depth: structureDepth,
      material,
      label: 'Lateral derecho',
    },
    {
      x: volume.x + t,
      y: volume.y,
      z: volume.z,
      width: innerW,
      height: t,
      depth: structureDepth,
      material,
      label: 'Piso',
    },
    {
      x: volume.x + t,
      y: volume.y + volume.height - t,
      z: volume.z,
      width: innerW,
      height: t,
      depth: structureDepth,
      material,
      label: 'Techo',
    },
    {
      x: volume.x + t - (useGrooves ? NORDEX_GROOVE_INSET_MM : 0),
      y: volume.y + t - (useGrooves ? NORDEX_GROOVE_INSET_MM : 0),
      z: backPanelZ,
      width: useGrooves ? nordexPanelSpanMm(innerW, 2) : innerW,
      height: useGrooves ? nordexPanelSpanMm(innerH, 2) : innerH,
      depth: backT,
      material: backMaterial,
      label: backId === 'nordex' ? 'Trasera nordex' : 'Trasera melamina',
    },
  ];

  if (useGrooves) {
    const gw = groove.grooveWidthMm;
    const gd = groove.grooveDepthMm;
    // Ranura = todo el borde del tablero de esa pieza (no el volumen exterior).
    panels.push(
      {
        x: volume.x + t - gd,
        y: volume.y,
        z: grooveZ,
        width: gd,
        height: volume.height,
        depth: gw,
        material: GROOVE,
        label: 'Ranura trasera (izq)',
      },
      {
        x: volume.x + volume.width - t,
        y: volume.y,
        z: grooveZ,
        width: gd,
        height: volume.height,
        depth: gw,
        material: GROOVE,
        label: 'Ranura trasera (der)',
      },
      {
        x: volume.x + t,
        y: volume.y + t - gd,
        z: grooveZ,
        width: innerW,
        height: gd,
        depth: gw,
        material: GROOVE,
        label: 'Ranura trasera (piso)',
      },
      {
        x: volume.x + t,
        y: volume.y + volume.height - t,
        z: grooveZ,
        width: innerW,
        height: gd,
        depth: gw,
        material: GROOVE,
        label: 'Ranura trasera (techo)',
      },
    );
  }

  return panels;
}

export function dividerPanel(left: ResolvedVolume, right: ResolvedVolume, parent: ResolvedVolume): PanelSpec | null {
  const t = PANEL_THICKNESS_MM;
  const material = {
    ...MELAMINE,
    color: resolveMaterialColor(parent.materialId ?? left.materialId),
  };

  if (left.x !== right.x) {
    const x = left.x + left.width;
    return {
      x,
      y: parent.y,
      z: parent.z,
      width: t,
      height: parent.height,
      depth: parent.depth,
      material,
    };
  }

  if (left.y !== right.y) {
    const y = left.y + left.height;
    return {
      x: parent.x + t,
      y,
      z: parent.z,
      width: innerWidth(parent),
      height: t,
      depth: parent.depth,
      material,
    };
  }

  if (left.z !== right.z) {
    const z = left.z + left.depth;
    return {
      x: parent.x + t,
      y: parent.y + t,
      z,
      width: innerWidth(parent),
      height: innerHeight(parent),
      depth: t,
      material,
    };
  }

  return null;
}

export function featurePanels(
  feature: ResolvedFeature,
  volume: ResolvedVolume,
  options?: { nestedInDesk?: boolean },
): PanelSpec[] {
  const material = {
    ...MELAMINE,
    color: resolveMaterialColor(volume.materialId),
  };
  const nestedInDesk = options?.nestedInDesk ?? false;
  const t = PANEL_THICKNESS_MM;
  const iw = nestedInDesk ? Math.max(0, volume.width - 2 * t) : innerWidth(volume);
  const ih = nestedInDesk ? volume.height : innerHeight(volume);
  const id = innerDepth(volume);
  const innerX = nestedInDesk ? volume.x + t : volume.x + t;
  const innerY = nestedInDesk ? volume.y : volume.y + t;
  const innerZ = volume.z;

  switch (feature.type) {
    case 'desk_frame':
      return deskFramePanels(feature, volume);
    case 'shelf_set':
      return shelfPanels(feature, innerX, innerY, innerZ, iw, ih, id, t, material);
    case 'drawer_stack':
      return [];
    case 'hanger_rod':
      return rodPanels(feature, volume, innerX, innerZ, iw, id);
    case 'appliance_space':
      return appliancePanels(innerX, innerY, innerZ, iw, ih, id);
    default:
      return [];
  }
}

function shelfPanels(
  feature: ResolvedFeature,
  innerX: number,
  innerY: number,
  innerZ: number,
  iw: number,
  ih: number,
  id: number,
  t: number,
  material: PanelMaterialProps,
): PanelSpec[] {
  const count = intParam(feature.params, 'count', 1);
  const panels: PanelSpec[] = [];

  for (let i = 0; i < count; i++) {
    const y = innerY + ((i + 1) * ih) / (count + 1) - t / 2;
    panels.push({
      x: innerX,
      y,
      z: innerZ,
      width: iw,
      height: t,
      depth: id,
      material,
    });
  }

  return panels;
}

function appliancePanels(
  innerX: number,
  innerY: number,
  innerZ: number,
  iw: number,
  ih: number,
  id: number,
): PanelSpec[] {
  const inset = 30;
  return [
    {
      x: innerX + inset,
      y: innerY + inset,
      z: innerZ + inset,
      width: Math.max(0, iw - inset * 2),
      height: Math.max(0, ih - inset * 2),
      depth: Math.max(0, id - inset * 2),
      material: { color: '#1a1f2e', roughness: 0.95, metalness: 0 },
    },
  ];
}

function rodPanels(
  feature: ResolvedFeature,
  volume: ResolvedVolume,
  innerX: number,
  innerZ: number,
  iw: number,
  id: number,
): PanelSpec[] {
  const fromTop = numberParam(feature.params, 'heightFromTop', 1800);
  const y = volume.y + volume.height - fromTop;
  const z = innerZ + id / 2;

  return [
    {
      x: innerX + 20,
      y: y - 10,
      z: z - 10,
      width: iw - 40,
      height: 20,
      depth: 20,
      material: METAL,
    },
    {
      x: innerX,
      y: y - 25,
      z: z - 8,
      width: 16,
      height: 50,
      depth: 16,
      material: METAL,
    },
    {
      x: innerX + iw - 16,
      y: y - 25,
      z: z - 8,
      width: 16,
      height: 50,
      depth: 16,
      material: METAL,
    },
  ];
}

export function frontPanel(front: ResolvedFront, volume: ResolvedVolume): PanelSpec {
  const material = {
    ...MELAMINE,
    color: resolveMaterialColor(
      typeof front.params.materialId === 'string' ? front.params.materialId : volume.materialId,
    ),
  };

  const z = front.z;

  if (front.type === 'door' || front.type === 'drawer_front') {
    const tol = 1;
    const outerWidth = Math.max(0, front.width - 2 * tol);
    const panelWidth = Math.max(0, outerWidth - 2 * THICK_EDGE_BANDING_MM);
    const outerX = front.x + tol;
    const panelX = outerX + THICK_EDGE_BANDING_MM;
    return {
      x: panelX,
      y: front.y,
      z,
      width: panelWidth,
      height: front.height,
      depth: FRONT_THICKNESS_MM,
      material,
    };
  }

  return {
    x: front.x,
    y: front.y,
    z,
    width: front.width,
    height: front.height,
    depth: FRONT_THICKNESS_MM,
    material,
  };
}
