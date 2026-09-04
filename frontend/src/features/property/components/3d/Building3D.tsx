import { useRef, useState, Suspense } from 'react';
import { Html, useGLTF, Bounds } from '@react-three/drei';
import * as THREE from 'three';
import type { Building } from '../../../building/types';

// Preload the building GLB model
useGLTF.preload('/3d-models/buildings/building.glb');

export interface Building3DProps {
  building: Building;
  position: [number, number, number];
  onClick: (b: Building) => void;
  selected: boolean;
}

// GLB Model loader component (wrapped in Suspense by parent)
function GLBBuildingModel() {
  const { scene } = useGLTF('/3d-models/buildings/building.glb');
  const cloned = scene.clone(true);

  cloned.traverse((child) => {
    if (child instanceof THREE.Mesh) {
      child.castShadow = true;
      child.receiveShadow = true;
    }
  });

  return (
    <Bounds fit clip observe margin={1.2}>
      <primitive object={cloned} />
    </Bounds>
  );
}

export function Building3D({ building, position, onClick, selected }: Building3DProps) {
  const meshRef = useRef<THREE.Group>(null);
  const [hovered, setHovered] = useState(false);

  return (
    <group
      ref={meshRef}
      position={position}
      onClick={(e) => {
        e.stopPropagation();
        onClick(building);
      }}
      onPointerOver={(e) => {
        e.stopPropagation();
        setHovered(true);
        document.body.style.cursor = 'pointer';
      }}
      onPointerOut={() => {
        setHovered(false);
        document.body.style.cursor = 'default';
      }}
    >
      <Suspense fallback={
        <mesh castShadow receiveShadow>
          <boxGeometry args={[2, 2 + building.total_floors * 0.5, 2]} />
          <meshStandardMaterial
            color={selected ? '#f97316' : hovered ? '#fbbf24' : '#94a3b8'}
            transparent
            opacity={0.85}
          />
        </mesh>
      }>
        <GLBBuildingModel />
      </Suspense>

      {selected && (
        <mesh>
          <boxGeometry args={[2.15, 2.15 + building.total_floors * 0.5, 2.15]} />
          <meshBasicMaterial color="#f97316" wireframe />
        </mesh>
      )}

      <Html position={[0, 2 + building.total_floors * 0.5 + 0.3, 0]} center distanceFactor={8}>
        <div className="px-2 py-1 bg-black/80 text-white text-xs rounded whitespace-nowrap pointer-events-none">
          {building.name}
          <span className="block text-[10px] text-orange-400">{building.total_floors} floors</span>
        </div>
      </Html>
    </group>
  );
}
