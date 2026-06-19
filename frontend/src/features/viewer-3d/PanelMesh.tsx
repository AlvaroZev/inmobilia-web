import type { PanelMaterialProps } from './materials';
import { describePanelMaterial } from './materials';
import type { PanelLayer } from './panel-layers';
import { InspectablePanelMesh } from './InspectablePanelMesh';

interface PanelMeshProps {
  x: number;
  y: number;
  z: number;
  width: number;
  height: number;
  depth: number;
  material: PanelMaterialProps;
  layer?: PanelLayer;
  label?: string;
  materialLabel?: string;
  materialId?: string;
  castShadow?: boolean;
  receiveShadow?: boolean;
  emissive?: string;
  onClick?: (event: import('@react-three/fiber').ThreeEvent<MouseEvent>) => void;
}

export function PanelMesh({
  x,
  y,
  z,
  width,
  height,
  depth,
  material,
  layer = 'carcass',
  label = 'Panel',
  materialLabel,
  materialId,
  castShadow = true,
  receiveShadow = true,
  emissive,
  onClick,
}: PanelMeshProps) {
  return (
    <InspectablePanelMesh
      x={x}
      y={y}
      z={z}
      width={width}
      height={height}
      depth={depth}
      material={material}
      layer={layer}
      label={label}
      materialLabel={materialLabel ?? describePanelMaterial(material, materialId)}
      castShadow={castShadow}
      receiveShadow={receiveShadow}
      emissive={emissive}
      onClick={onClick}
    />
  );
}
