import { useRef, useState, Suspense } from 'react';
import { Html, useGLTF } from '@react-three/drei';
import * as THREE from 'three';
import type { Building } from '../../../building/types';

// Preload both 3D models
useGLTF.preload('/3d-models/buildings/building.glb');
useGLTF.preload('/3d-models/buildings/hotel.glb');

export interface Building3DProps {
  building: Building;
  position: [number, number, number];
  onClick: (b: Building) => void;
  selected: boolean;
  modelStyle?: 'building' | 'hotel' | 'auto';
}

function DynamicGLBModel({ path, scale }: { path: string; scale: [number, number, number] }) {
  const { scene } = useGLTF(path);
  const cloned = scene.clone(true);

  cloned.traverse((child) => {
    if (child instanceof THREE.Mesh) {
      child.castShadow = true;
      child.receiveShadow = true;
    }
  });

  return (
    <primitive object={cloned} scale={scale} position={[0, 0, 0]} />
  );
}

export function Building3D({ building, position, onClick, selected, modelStyle = 'auto' }: Building3DProps) {
  const meshRef = useRef<THREE.Group>(null);
  const [hovered, setHovered] = useState(false);

  // Check stored preference or name heuristics
  let isHotel = false;
  if (modelStyle === 'hotel') {
    isHotel = true;
  } else if (modelStyle === 'building') {
    isHotel = false;
  } else {
    try {
      const pref = localStorage.getItem(`building_model_pref_${building.name}`);
      if (pref === 'hotel_resort') {
        isHotel = true;
      } else {
        const lower = building.name.toLowerCase();
        isHotel = lower.includes('hotel') || lower.includes('resort') || lower.includes('suite');
      }
    } catch {
      const lower = building.name.toLowerCase();
      isHotel = lower.includes('hotel') || lower.includes('resort') || lower.includes('suite');
    }
  }

  const modelPath = isHotel ? '/3d-models/buildings/hotel.glb' : '/3d-models/buildings/building.glb';
  const modelScale: [number, number, number] = isHotel ? [0.55, 0.55, 0.55] : [0.6, 0.6, 0.6];
  const modelIcon = isHotel ? '🏨' : '🏢';

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
        <DynamicGLBModel key={modelPath} path={modelPath} scale={modelScale} />
      </Suspense>

      {/* Selected Indicator Ring */}
      {selected && (
        <mesh position={[0, 0.05, 0]} rotation={[-Math.PI / 2, 0, 0]}>
          <ringGeometry args={[1.8, 2.1, 32]} />
          <meshBasicMaterial color="#f97316" side={THREE.DoubleSide} />
        </mesh>
      )}

      {/* Building Label with Model Icon */}
      <Html position={[0, 3.8, 0]} center distanceFactor={13}>
        <div className={`px-2.5 py-1 rounded-lg text-xs font-semibold whitespace-nowrap shadow-lg border transition-all pointer-events-none flex items-center gap-1.5 ${
          selected
            ? 'bg-orange text-white border-orange shadow-orange/30'
            : 'bg-black/85 text-white border-white/10'
        }`}>
          <span>{modelIcon}</span>
          <div>
            <span>{building.name}</span>
            <span className="block text-[10px] opacity-75 font-normal">{building.total_floors} floors</span>
          </div>
        </div>
      </Html>
    </group>
  );
}
