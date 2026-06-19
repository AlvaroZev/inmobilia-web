import { Edges } from '@react-three/drei';
import type { ThreeEvent } from '@react-three/fiber';
import { useMemo } from 'react';
import type { ResolvedVolume } from '@/domain/resolved-furniture';
import { useViewerStore } from '@/store/viewer-store';
import {
  dividerPanel,
  hasDeskFrame,
  isNestedDrawerTower,
  nestedDrawerTowerPanels,
  outerCarcassPanels,
} from './geometry';
import {
  carcassStructureDepthMm,
  resolveVolumeBackMaterialId,
} from './back-panel';
import { GROOVE } from './materials';
import { volumeColor, volumeStructureDebugBox } from './scene-utils';
import { FeatureMesh } from './FeatureMesh';
import { FrontMesh } from './FrontMesh';
import { PanelMesh } from './PanelMesh';

interface ResolvedVolumeMeshProps {
  volume: ResolvedVolume;
  depth: number;
  parent?: ResolvedVolume;
}

export function ResolvedVolumeMesh({ volume, depth, parent }: ResolvedVolumeMeshProps) {
  const selectedVolumeId = useViewerStore((s) => s.selectedVolumeId);
  const showVolumes = useViewerStore((s) => s.showVolumes);
  const showFeatures = useViewerStore((s) => s.showFeatures);
  const showFronts = useViewerStore((s) => s.showFronts);
  const setSelectedVolumeId = useViewerStore((s) => s.setSelectedVolumeId);

  const isRoot = depth === 0;
  const isLeaf = volume.children.length === 0;
  const color = volumeColor(depth, volume.id);
  const isSelected = selectedVolumeId === volume.id;

  const nestedDrawer = isNestedDrawerTower(volume, parent);
  const isStandaloneDrawerCabinet =
    isLeaf && volume.features.some((f) => f.type === 'drawer_stack') && !hasDeskFrame(volume) && !nestedDrawer;

  const backMaterialId = useMemo(() => resolveVolumeBackMaterialId(volume), [volume]);

  const carcassPanels = useMemo(() => {
    if (isRoot && !hasDeskFrame(volume) && !volume.children.some((child) => hasDeskFrame(child))) {
      return outerCarcassPanels(volume, backMaterialId);
    }
    if (nestedDrawer) {
      return nestedDrawerTowerPanels(volume);
    }
    if (isStandaloneDrawerCabinet) {
      return outerCarcassPanels(volume, backMaterialId);
    }
    return [];
  }, [isRoot, nestedDrawer, isStandaloneDrawerCabinet, volume, backMaterialId]);

  const debugBox = useMemo(() => {
    if (carcassPanels.length === 0) {
      return volumeStructureDebugBox(volume, volume.depth);
    }
    const structureDepth = carcassStructureDepthMm(volume.depth, backMaterialId);
    return volumeStructureDebugBox(volume, structureDepth);
  }, [backMaterialId, carcassPanels.length, volume]);

  const structuralFeatures = useMemo(
    () => volume.features.filter((feature) => feature.type === 'desk_frame'),
    [volume.features],
  );

  const handleClick = (event: ThreeEvent<MouseEvent>) => {
    event.stopPropagation();
    setSelectedVolumeId(volume.id);
  };

  return (
    <group onClick={handleClick}>
      {carcassPanels.map((panel, index) => (
        <PanelMesh
          key={`carcass-${volume.id}-${index}`}
          {...panel}
          layer="carcass"
          label={panel.label ?? `Carcasa ${index + 1}`}
          materialId={
            panel.material === GROOVE
              ? undefined
              : panel.label?.includes('nordex')
                ? 'nordex'
                : volume.materialId
          }
        />
      ))}

      {volume.children.map((child, index) => (
        <group key={child.id}>
          {index > 0 && (() => {
            const divider = dividerPanel(volume.children[index - 1], child, volume);
            return divider ? (
              <PanelMesh
                key={`divider-${child.id}`}
                {...divider}
                layer="carcass"
                label="División"
                materialId={volume.materialId}
              />
            ) : null;
          })()}
          <ResolvedVolumeMesh volume={child} depth={depth + 1} parent={volume} />
        </group>
      ))}

      {showFeatures &&
        structuralFeatures.map((feature) => (
          <FeatureMesh key={feature.id} feature={feature} volume={volume} />
        ))}

      {isLeaf &&
        showFeatures &&
        volume.features
          .filter((feature) => feature.type !== 'desk_frame')
          .map((feature) => (
            <FeatureMesh
              key={feature.id}
              feature={feature}
              volume={volume}
              nestedInDesk={nestedDrawer}
            />
          ))}

      {isLeaf &&
        showFronts &&
        !(parent && hasDeskFrame(parent)) &&
        volume.fronts.map((front) => (
          <FrontMesh key={front.id} front={front} volume={volume} />
        ))}

      {showVolumes && (
        <mesh position={debugBox.position}>
          <boxGeometry args={debugBox.size} />
          <meshStandardMaterial
            color={color}
            transparent
            opacity={isSelected ? 0.18 : 0.08}
            depthWrite={false}
          />
          <Edges
            color={isSelected ? '#f39c12' : '#5d6d7e'}
            lineWidth={isSelected ? 2 : 1}
            threshold={15}
          />
        </mesh>
      )}
    </group>
  );
}
