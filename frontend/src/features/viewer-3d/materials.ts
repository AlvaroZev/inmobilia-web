export const PANEL_THICKNESS_MM = 18;
export const BACK_THICKNESS_MM = 6;
export const FRONT_THICKNESS_MM = 18;
export const BACK_SETBACK_MM = 8;

export interface PanelMaterialProps {
  color: string;
  roughness?: number;
  metalness?: number;
}

const MATERIAL_COLORS: Record<string, string> = {
  'melamine-white-18': '#f2efe8',
  'melamine-white': '#f2efe8',
  'pvc-white-1mm': '#faf8f4',
  nordex: '#c4b8a8',
  default: '#e8e4dc',
};

export function resolveMaterialColor(materialId?: string): string {
  if (!materialId) {
    return MATERIAL_COLORS.default;
  }
  return MATERIAL_COLORS[materialId] ?? MATERIAL_COLORS.default;
}

export const MELAMINE: PanelMaterialProps = {
  color: MATERIAL_COLORS['melamine-white-18'],
  roughness: 0.72,
  metalness: 0.02,
};

export const BACK_PANEL: PanelMaterialProps = {
  color: '#d8d0c4',
  roughness: 0.85,
  metalness: 0,
};

export const STRUCTURE: PanelMaterialProps = {
  color: '#f0f0eb',
  roughness: 0.82,
  metalness: 0.01,
};

export const NORDEX: PanelMaterialProps = {
  color: MATERIAL_COLORS.nordex,
  roughness: 0.92,
  metalness: 0,
};

export function resolvePanelMaterial(materialId?: string): PanelMaterialProps {
  if (materialId === 'nordex') {
    return NORDEX;
  }
  return {
    ...MELAMINE,
    color: resolveMaterialColor(materialId),
  };
}

export const METAL: PanelMaterialProps = {
  color: '#b8bcc4',
  roughness: 0.35,
  metalness: 0.85,
};

export const HANDLE: PanelMaterialProps = {
  color: '#4a5568',
  roughness: 0.4,
  metalness: 0.6,
};
