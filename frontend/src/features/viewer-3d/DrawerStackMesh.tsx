import type { ThreeEvent } from '@react-three/fiber';
import { useFrame } from '@react-three/fiber';
import { useMemo, useRef } from 'react';
import type { Group } from 'three';
import * as THREE from 'three';
import type { ResolvedFeature, ResolvedVolume } from '@/domain/resolved-furniture';
import { useViewerStore } from '@/store/viewer-store';
import {
  buildDrawerStackGeometry,
  type DrawerPanelRole,
  type DrawerPanelSpec,
  type DrawerUnitSpec,
} from './drawer-geometry';
import { InspectablePanelMesh } from './InspectablePanelMesh';
import {
  HANDLE,
  GROOVE,
  METAL,
  MELAMINE,
  STRUCTURE,
  describePanelMaterial,
  resolveMaterialColor,
  resolvePanelMaterial,
  type PanelMaterialProps,
} from './materials';
import { drawerRoleLabel, drawerRoleLayer } from './panel-layers';
import { mm } from './scene-utils';

interface DrawerStackMeshProps {
  feature: ResolvedFeature;
  volume: ResolvedVolume;
  nestedInDesk?: boolean;
}

interface DrawerUnitMeshProps {
  unit: DrawerUnitSpec;
  volume: ResolvedVolume;
  feature: ResolvedFeature;
  nestedInDesk: boolean;
}

function panelMaterial(
  role: DrawerPanelRole,
  volume: ResolvedVolume,
  feature: ResolvedFeature,
  nestedInDesk: boolean,
): PanelMaterialProps {
  switch (role) {
    case 'structure':
      return STRUCTURE;
    case 'bottom':
      return resolvePanelMaterial(
        typeof feature.params.bottomMaterialId === 'string'
          ? feature.params.bottomMaterialId
          : typeof feature.params.backMaterialId === 'string'
            ? feature.params.backMaterialId
            : nestedInDesk
              ? 'nordex'
              : 'nordex',
      );
    case 'front':
      return {
        ...MELAMINE,
        color: resolveMaterialColor(volume.materialId),
      };
    case 'front_edge_top':
    case 'front_edge_bottom':
    case 'front_edge_left':
    case 'front_edge_right':
      return {
        ...MELAMINE,
        color: '#e8e0d4',
        roughness: 0.65,
      };
    case 'handle':
      return HANDLE;
    case 'runner_cabinet':
      return { ...METAL, color: '#8a9199' };
    case 'runner_drawer':
      return { ...METAL, color: '#b0b8c0' };
    case 'runner':
    case 'groove_rail':
      return METAL;
    case 'groove_cut':
      return GROOVE;
    default:
      return MELAMINE;
  }
}

function panelMaterialId(
  role: DrawerPanelRole,
  volume: ResolvedVolume,
  feature: ResolvedFeature,
  _nestedInDesk: boolean,
): string | undefined {
  if (role === 'front' || role === 'front_edge_top' || role === 'front_edge_bottom' || role === 'front_edge_left' || role === 'front_edge_right') {
    return volume.materialId;
  }
  if (role === 'bottom') {
    const id = feature.params.backMaterialId ?? feature.params.bottomMaterialId;
    if (typeof id === 'string') {
      return id;
    }
    return 'nordex';
  }
  return undefined;
}

function DrawerPanel({
  panel,
  material,
  materialId,
  emissive,
}: {
  panel: DrawerPanelSpec;
  material: PanelMaterialProps;
  materialId?: string;
  emissive?: string;
}) {
  return (
    <InspectablePanelMesh
      x={panel.x}
      y={panel.y}
      z={panel.z}
      width={panel.width}
      height={panel.height}
      depth={panel.depth}
      material={material}
      layer={drawerRoleLayer(panel.role)}
      label={drawerRoleLabel(panel.role)}
      materialLabel={describePanelMaterial(material, materialId)}
      emissive={emissive}
    />
  );
}

function DrawerUnitMesh({ unit, volume, feature, nestedInDesk }: DrawerUnitMeshProps) {
  const groupRef = useRef<Group>(null);
  const slideRef = useRef(0);
  const drawerState = useViewerStore((s) => s.openDrawers[unit.id] ?? 'closed');
  const toggleDrawer = useViewerStore((s) => s.toggleDrawer);
  const isOpen = drawerState === 'half';
  const targetSlide = isOpen ? unit.slideDistance : 0;

  useFrame(() => {
    if (!groupRef.current) {
      return;
    }
    slideRef.current = THREE.MathUtils.lerp(slideRef.current, targetSlide, 0.14);
    groupRef.current.position.z = mm(unit.worldOrigin[2] - slideRef.current);
  });

  const handleClick = (event: ThreeEvent<MouseEvent>) => {
    event.stopPropagation();
    toggleDrawer(unit.id);
  };

  return (
    <group
      ref={groupRef}
      position={[mm(unit.worldOrigin[0]), mm(unit.worldOrigin[1]), mm(unit.worldOrigin[2])]}
      onClick={handleClick}
    >
      {unit.panels.map((panel, index) => {
        const mat = panelMaterial(panel.role, volume, feature, nestedInDesk);
        return (
          <DrawerPanel
            key={`${unit.id}-${index}`}
            panel={panel}
            material={mat}
            materialId={panelMaterialId(panel.role, volume, feature, nestedInDesk)}
            emissive={isOpen && panel.role === 'front' ? '#f39c12' : undefined}
          />
        );
      })}
    </group>
  );
}

function FixedRunnerSlot({
  slot,
  volumeOrigin,
  volume,
  feature,
  nestedInDesk,
}: {
  slot: { slotY: number; panels: DrawerPanelSpec[] };
  volumeOrigin: [number, number, number];
  volume: ResolvedVolume;
  feature: ResolvedFeature;
  nestedInDesk: boolean;
}) {
  return (
    <group
      position={[
        mm(volumeOrigin[0]),
        mm(volumeOrigin[1] + slot.slotY),
        mm(volumeOrigin[2]),
      ]}
    >
      {slot.panels.map((panel, index) => (
        <DrawerPanel
          key={`cabinet-runner-${index}`}
          panel={panel}
          material={panelMaterial(panel.role, volume, feature, nestedInDesk)}
        />
      ))}
    </group>
  );
}

export function DrawerStackMesh({ feature, volume, nestedInDesk = false }: DrawerStackMeshProps) {
  const geometry = useMemo(
    () => buildDrawerStackGeometry(feature, volume, nestedInDesk),
    [feature, volume, nestedInDesk],
  );

  const volumeOrigin = useMemo((): [number, number, number] => {
    const first = geometry.units[0];
    if (first) {
      return [first.worldOrigin[0], volume.y, first.worldOrigin[2]];
    }
    return [volume.x, volume.y, volume.z];
  }, [geometry.units, volume]);

  return (
    <group>
      {geometry.fixedRunnerSlots.map((slot, index) => (
        <FixedRunnerSlot
          key={`fixed-runners-${index}`}
          slot={slot}
          volumeOrigin={volumeOrigin}
          volume={volume}
          feature={feature}
          nestedInDesk={nestedInDesk}
        />
      ))}
      {geometry.units.map((unit) => (
        <DrawerUnitMesh
          key={unit.id}
          unit={unit}
          volume={volume}
          feature={feature}
          nestedInDesk={nestedInDesk}
        />
      ))}
    </group>
  );
}
