import type { ThreeEvent } from '@react-three/fiber';
import { useMemo } from 'react';
import * as THREE from 'three';
import { useViewerStore } from '@/store/viewer-store';
import { roundMm } from '@/utils/round-mm';
import type { PanelMaterialProps } from './materials';
import type { PanelInspectInfo, PanelLayer } from './panel-layers';
import { boxCenter, boxSize } from './scene-utils';

const GHOST_COLOR = '#4da3ff';
const GHOST_OPACITY = 0.118;

interface InspectablePanelMeshProps {
  x: number;
  y: number;
  z: number;
  width: number;
  height: number;
  depth: number;
  material: PanelMaterialProps;
  layer: PanelLayer;
  label: string;
  materialLabel: string;
  castShadow?: boolean;
  receiveShadow?: boolean;
  emissive?: string;
  onClick?: (event: ThreeEvent<MouseEvent>) => void;
  onPointerOver?: (event: ThreeEvent<PointerEvent>) => void;
  onPointerOut?: (event: ThreeEvent<PointerEvent>) => void;
}

export function InspectablePanelMesh({
  x,
  y,
  z,
  width,
  height,
  depth,
  material,
  layer,
  label,
  materialLabel,
  castShadow = true,
  receiveShadow = true,
  emissive,
  onClick,
  onPointerOver,
  onPointerOut,
}: InspectablePanelMeshProps) {
  const ghost = useViewerStore((s) => s.ghostLayers[layer]);
  const setHoveredPanel = useViewerStore((s) => s.setHoveredPanel);

  const inspectInfo = useMemo(
    (): PanelInspectInfo => ({
      label,
      layer,
      width: roundMm(width),
      height: roundMm(height),
      depth: roundMm(depth),
      color: material.color,
      materialLabel,
    }),
    [label, layer, width, height, depth, material.color, materialLabel],
  );

  if (width <= 0 || height <= 0 || depth <= 0) {
    return null;
  }

  const handlePointerOver = (event: ThreeEvent<PointerEvent>) => {
    event.stopPropagation();
    setHoveredPanel(inspectInfo);
    document.body.style.cursor = onClick ? 'pointer' : 'help';
    onPointerOver?.(event);
  };

  const handlePointerOut = (event: ThreeEvent<PointerEvent>) => {
    event.stopPropagation();
    setHoveredPanel(null);
    document.body.style.cursor = 'default';
    onPointerOut?.(event);
  };

  return (
    <mesh
      position={boxCenter(x, y, z, width, height, depth)}
      castShadow={castShadow && !ghost}
      receiveShadow={receiveShadow && !ghost}
      onClick={onClick}
      onPointerOver={handlePointerOver}
      onPointerOut={handlePointerOut}
    >
      <boxGeometry args={boxSize(width, height, depth)} />
      <meshStandardMaterial
        color={ghost ? GHOST_COLOR : material.color}
        roughness={ghost ? 0.2 : (material.roughness ?? 0.7)}
        metalness={ghost ? 0 : (material.metalness ?? 0.05)}
        transparent={true}
        premultipliedAlpha={true}
        opacity={ghost ? GHOST_OPACITY : 1}
        depthWrite={!ghost}
        emissive={
          ghost
            ? new THREE.Color('#2563eb')
            : emissive
              ? new THREE.Color(emissive)
              : undefined
        }
        emissiveIntensity={ghost ? 0.08 : emissive ? 0.12 : 0}
      />
    </mesh>
  );
}
