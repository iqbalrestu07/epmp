import { useRef, useState, Suspense } from 'react';
import { Html, useGLTF } from '@react-three/drei';
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

// GLB Model loader component
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
    <primitive object={cloned} scale={[0.6, 0.6, 0.6]} position={[0, 0, 0]} />
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
      <Suspense
        fallback={
          <mesh castShadow receiveShadow position={[0, (2 + building.total_floors * 0.5) / 2, 0]}>
            <boxGeometry args={[2, 2 + building.total_floors * 0.5, 2]} />
            <meshStandardMaterial
              color={selected ? '#f97316' : hovered ? '#fbbf24' : '#94a3b8'}
              transparent
              opacity={0.85}
            />
          </mesh>
        }
      >
        <GLBBuildingModel />
      </Suspense>

      {/* Selected Indicator Ring */}
      {selected && (
        <mesh position={[0, 0.05, 0]} rotation={[-Math.PI / 2, 0, 0]}>
          <ringGeometry args={[1.8, 2.1, 32]} />
          <meshBasicMaterial color="#f97316" side={THREE.DoubleSide} />
        </mesh>
      )}

      {/* Building Label */}
      <Html position={[0, 3.5, 0]} center distanceFactor={12}>
        <div className={`px-2.5 py-1 rounded-lg text-xs font-semibold whitespace-nowrap shadow-lg border transition-all pointer-events-none ${
          selected
            ? 'bg-orange text-white border-orange shadow-orange/30'
            : 'bg-black/80 text-white border-white/10'
        }`}>
          {building.name}
          <span className="block text-[10px] opacity-75 font-normal">{building.total_floors} floors</span>
        </div>
      </Html>
    </group>
  );
}
