import type { PanelMaterialProps } from './materials';
import { boxCenter, boxSize } from './scene-utils';

interface PanelMeshProps {
  x: number;
  y: number;
  z: number;
  width: number;
  height: number;
  depth: number;
  material: PanelMaterialProps;
  castShadow?: boolean;
  receiveShadow?: boolean;
}

export function PanelMesh({
  x,
  y,
  z,
  width,
  height,
  depth,
  material,
  castShadow = true,
  receiveShadow = true,
}: PanelMeshProps) {
  if (width <= 0 || height <= 0 || depth <= 0) {
    return null;
  }

  return (
    <mesh
      position={boxCenter(x, y, z, width, height, depth)}
      castShadow={castShadow}
      receiveShadow={receiveShadow}
    >
      <boxGeometry args={boxSize(width, height, depth)} />
      <meshStandardMaterial
        color={material.color}
        roughness={material.roughness ?? 0.7}
        metalness={material.metalness ?? 0.05}
      />
    </mesh>
  );
}
