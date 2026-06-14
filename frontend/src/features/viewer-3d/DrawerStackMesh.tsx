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
import {
  HANDLE,
  METAL,
  MELAMINE,
  STRUCTURE,
  resolveMaterialColor,
  resolvePanelMaterial,
  type PanelMaterialProps,
} from './materials';
import { boxCenter, boxSize, mm } from './scene-utils';

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
          : nestedInDesk
            ? 'nordex'
            : undefined,
      );
    case 'front':
      return {
        ...MELAMINE,
        color: resolveMaterialColor(volume.materialId),
      };
    case 'handle':
      return HANDLE;
    case 'runner':
    case 'groove_rail':
      return METAL;
    default:
      return MELAMINE;
  }
}

function LocalPanel({
  panel,
  material,
  emissive,
}: {
  panel: DrawerPanelSpec;
  material: PanelMaterialProps;
  emissive?: string;
}) {
  if (panel.width <= 0 || panel.height <= 0 || panel.depth <= 0) {
    return null;
  }

  return (
    <mesh position={boxCenter(panel.x, panel.y, panel.z, panel.width, panel.height, panel.depth)} castShadow receiveShadow>
      <boxGeometry args={boxSize(panel.width, panel.height, panel.depth)} />
      <meshStandardMaterial
        color={material.color}
        roughness={material.roughness ?? 0.7}
        metalness={material.metalness ?? 0.05}
        emissive={emissive ? new THREE.Color(emissive) : undefined}
        emissiveIntensity={emissive ? 0.12 : 0}
      />
    </mesh>
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
    groupRef.current.position.z = mm(slideRef.current);
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
      onPointerOver={(event) => {
        event.stopPropagation();
        document.body.style.cursor = 'pointer';
      }}
      onPointerOut={() => {
        document.body.style.cursor = 'default';
      }}
    >
      {unit.panels.map((panel, index) => (
        <LocalPanel
          key={`${unit.id}-${index}`}
          panel={panel}
          material={panelMaterial(panel.role, volume, feature, nestedInDesk)}
          emissive={isOpen && panel.role === 'front' ? '#f39c12' : undefined}
        />
      ))}
    </group>
  );
}

function FramePanels({
  panels,
  worldOrigin,
  volume,
  feature,
  nestedInDesk,
}: {
  panels: DrawerPanelSpec[];
  worldOrigin: [number, number, number];
  volume: ResolvedVolume;
  feature: ResolvedFeature;
  nestedInDesk: boolean;
}) {
  if (panels.length === 0) {
    return null;
  }

  return (
    <group position={[mm(worldOrigin[0]), mm(worldOrigin[1]), mm(worldOrigin[2])]}>
      {panels.map((panel, index) => (
        <LocalPanel
          key={`frame-${index}`}
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

  const frameOrigin = geometry.units[0]?.worldOrigin ?? [volume.x, volume.y, volume.z];

  return (
    <group>
      <FramePanels
        panels={geometry.framePanels}
        worldOrigin={frameOrigin as [number, number, number]}
        volume={volume}
        feature={feature}
        nestedInDesk={nestedInDesk}
      />
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
