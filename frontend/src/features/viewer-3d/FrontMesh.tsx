import { useMemo } from 'react';
import type { ResolvedFront, ResolvedVolume } from '@/domain/resolved-furniture';
import { frontPanel } from './geometry';
import { HANDLE } from './materials';
import { boxCenter, boxSize } from './scene-utils';
import { PanelMesh } from './PanelMesh';

interface FrontMeshProps {
  front: ResolvedFront;
  volume: ResolvedVolume;
}

export function FrontMesh({ front, volume }: FrontMeshProps) {
  const panel = useMemo(() => frontPanel(front, volume), [front, volume]);
  const isGlass = front.type === 'glass';
  const showHandle = front.type === 'door' || front.type === 'drawer_front';

  const handle = useMemo(() => {
    if (!showHandle) {
      return null;
    }
    const hinge = front.params.hinge === 'right' ? 'right' : 'left';
    const handleX =
      hinge === 'right'
        ? front.x + front.width - 120
        : front.x + 80;
    const handleY = front.y + front.height / 2 - 60;

    return {
      x: handleX,
      y: handleY,
      z: front.z + front.depth - 4,
      width: 12,
      height: 120,
      depth: 16,
    };
  }, [front, showHandle]);

  if (isGlass) {
    return (
      <mesh position={boxCenter(panel.x, panel.y, panel.z, panel.width, panel.height, panel.depth)}>
        <boxGeometry args={boxSize(panel.width, panel.height, panel.depth)} />
        <meshPhysicalMaterial
          color="#9ec5e8"
          transparent
          opacity={0.35}
          roughness={0.1}
          metalness={0.05}
          transmission={0.6}
          thickness={0.02}
        />
      </mesh>
    );
  }

  return (
    <group>
      <PanelMesh {...panel} />
      {handle && <PanelMesh {...handle} material={HANDLE} />}
    </group>
  );
}
