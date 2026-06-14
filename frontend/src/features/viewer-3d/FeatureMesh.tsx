import type { ResolvedFeature, ResolvedVolume } from '@/domain/resolved-furniture';
import { DrawerStackMesh } from './DrawerStackMesh';
import { featurePanels } from './geometry';
import { PanelMesh } from './PanelMesh';

interface FeatureMeshProps {
  feature: ResolvedFeature;
  volume: ResolvedVolume;
  nestedInDesk?: boolean;
}

export function FeatureMesh({ feature, volume, nestedInDesk }: FeatureMeshProps) {
  if (feature.type === 'drawer_stack') {
    return <DrawerStackMesh feature={feature} volume={volume} nestedInDesk={nestedInDesk} />;
  }

  const panels = featurePanels(feature, volume, { nestedInDesk });

  return (
    <group>
      {panels.map((panel, index) => (
        <PanelMesh key={`${feature.id}-${index}`} {...panel} />
      ))}
    </group>
  );
}
