import type { DrawerPanelRole } from './drawer-geometry';

export type PanelLayer = 'carcass' | 'drawerFront' | 'drawerInterior' | 'runners' | 'handle' | 'fronts';

export const PANEL_LAYER_LABELS: Record<PanelLayer, string> = {
  carcass: 'Carcasa',
  drawerFront: 'Frente cajón',
  drawerInterior: 'Interior cajón',
  runners: 'Correderas',
  handle: 'Tirador',
  fronts: 'Frentes (puertas)',
};

export const DEFAULT_GHOST_LAYERS: Record<PanelLayer, boolean> = {
  carcass: false,
  drawerFront: false,
  drawerInterior: false,
  runners: false,
  handle: false,
  fronts: false,
};

export interface PanelInspectInfo {
  label: string;
  layer: PanelLayer;
  width: number;
  height: number;
  depth: number;
  color: string;
  materialLabel: string;
}

export function drawerRoleLayer(role: DrawerPanelRole): PanelLayer {
  switch (role) {
    case 'front':
    case 'front_edge_top':
    case 'front_edge_bottom':
    case 'front_edge_left':
    case 'front_edge_right':
      return 'drawerFront';
    case 'handle':
      return 'handle';
    case 'runner':
    case 'runner_cabinet':
    case 'runner_drawer':
    case 'groove_rail':
      return 'runners';
    case 'groove_cut':
      return 'drawerInterior';
    default:
      return 'drawerInterior';
  }
}

export function drawerRoleLabel(role: DrawerPanelRole): string {
  switch (role) {
    case 'structure':
      return 'Estructura cajón';
    case 'bottom':
      return 'Fondo cajón';
    case 'front':
      return 'Frente final';
    case 'front_edge_top':
    case 'front_edge_bottom':
    case 'front_edge_left':
    case 'front_edge_right':
      return 'Canto grueso frente';
    case 'handle':
      return 'Tirador';
    case 'runner':
    case 'runner_cabinet':
      return 'Corredera carcasa';
    case 'runner_drawer':
      return 'Corredera cajón';
    case 'groove_rail':
      return 'Riel ranura';
    case 'groove_cut':
      return 'Ranura fondo nordex';
    default:
      return role;
  }
}
