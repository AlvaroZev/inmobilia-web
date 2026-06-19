import { ContactShadows, Environment, Grid, OrbitControls } from '@react-three/drei';
import { Canvas } from '@react-three/fiber';
import { useMemo } from 'react';
import type { ResolvedFurniture } from '@/domain/resolved-furniture';
import { furnitureSceneCenter, furnitureSceneSize } from './scene-utils';
import { ResolvedVolumeMesh } from './ResolvedVolumeMesh';

interface ResolvedFurnitureViewerProps {
  furniture: ResolvedFurniture;
}

function FurnitureScene({ furniture }: ResolvedFurnitureViewerProps) {
  const center = useMemo(() => furnitureSceneCenter(furniture), [furniture]);
  const sceneSize = useMemo(() => furnitureSceneSize(furniture), [furniture]);

  return (
    <>
      <ambientLight intensity={0.45} />
      <directionalLight position={[5, 9, 7]} intensity={1.35} castShadow shadow-mapSize={[1024, 1024]} />
      <directionalLight position={[-5, 4, -3]} intensity={0.45} />
      <hemisphereLight args={['#f5f0e8', '#1a1f2e', 0.35]} />

      <group position={[-center[0], -center[1], -center[2]]}>
        <ResolvedVolumeMesh volume={furniture.root} depth={0} />
        <Grid
          infiniteGrid
          fadeDistance={sceneSize * 4}
          sectionColor="#4a5568"
          cellColor="#2d3748"
          position={[center[0], 0.001, center[2]]}
        />
        <ContactShadows
          position={[center[0], 0.002, center[2]]}
          opacity={0.5}
          scale={sceneSize * 2}
          blur={2.8}
          far={sceneSize}
        />
      </group>

      <OrbitControls
        makeDefault
        target={center}
        minDistance={sceneSize * 0.4}
        maxDistance={sceneSize * 4}
      />
      <Environment preset="apartment" />
    </>
  );
}

export function ResolvedFurnitureViewer({ furniture }: ResolvedFurnitureViewerProps) {
  const cameraDistance = furnitureSceneSize(furniture) * 1.8;

  return (
    <Canvas
      shadows
      camera={{
        // Frente del mueble en −Z; cámara al frente para ver apertura de cajones hacia el usuario.
        position: [cameraDistance * 0.8, cameraDistance * 0.65, -cameraDistance],
        fov: 45,
        near: 0.01,
        far: 100,
      }}
      style={{ width: '100%', height: '100%' }}
    >
      <color attach="background" args={['#0a0c10']} />
      <FurnitureScene furniture={furniture} />
    </Canvas>
  );
}
