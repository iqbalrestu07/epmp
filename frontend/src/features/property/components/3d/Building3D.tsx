import { useRef, useState, useMemo, Suspense } from 'react';
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

function DynamicGLBModel({ path, targetHeight = 3.2 }: { path: string; targetHeight?: number }) {
  const { scene } = useGLTF(path);

  const normalizedScene = useMemo(() => {
    const cloned = scene.clone(true);

    // Hide giant landscape/terrain plate mesh in building.glb if present (571k units)
    cloned.traverse((c) => {
      if (
        c instanceof THREE.Mesh &&
        (c.name === 'Object_181' ||
          (c.geometry?.boundingBox &&
            c.geometry.boundingBox.getSize(new THREE.Vector3()).x > 200000))
      ) {
        c.visible = false;
      }
    });

    // Compute bounding box strictly from visible meshes
    const box = new THREE.Box3();
    cloned.traverse((c) => {
      if (c instanceof THREE.Mesh && c.visible) {
        c.castShadow = true;
        c.receiveShadow = true;
        box.expandByObject(c);
      }
    });

    const size = new THREE.Vector3();
    box.getSize(size);
    const center = new THREE.Vector3();
    box.getCenter(center);

    const safeHeight = size.y > 0 ? size.y : 1;
    const scaleFactor = targetHeight / safeHeight;

    const wrapper = new THREE.Group();
    // Center horizontally and set bottom on ground (y = 0)
    cloned.position.set(-center.x, -box.min.y, -center.z);
    wrapper.add(cloned);
    wrapper.scale.set(scaleFactor, scaleFactor, scaleFactor);

    return wrapper;
  }, [scene, targetHeight]);

  return <primitive object={normalizedScene} position={[0, 0, 0]} />;
}

export function Building3D({ building, position, onClick, selected, modelStyle = 'auto' }: Building3DProps) {
  const meshRef = useRef<THREE.Group>(null);
  const pointerDownPos = useRef<{ x: number; y: number }>({ x: 0, y: 0 });
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
  const modelIcon = isHotel ? '🏨' : '🏢';

  return (
    <group
      ref={meshRef}
      position={position}
      onPointerDown={(e) => {
        // Record starting pointer position to differentiate drag/orbit from click
        pointerDownPos.current = { x: e.clientX, y: e.clientY };
      }}
      onClick={(e) => {
        e.stopPropagation();
        // If pointer moved more than 6px, it was a camera orbit/drag gesture, not a click!
        const dx = e.clientX - pointerDownPos.current.x;
        const dy = e.clientY - pointerDownPos.current.y;
        const dist = Math.sqrt(dx * dx + dy * dy);
        if (dist > 6) {
          return;
        }
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
          <mesh castShadow receiveShadow position={[0, 1.6, 0]}>
            <boxGeometry args={[2, 3.2, 2]} />
            <meshStandardMaterial
              color={selected ? '#f97316' : hovered ? '#fbbf24' : '#94a3b8'}
              transparent
              opacity={0.85}
            />
          </mesh>
        }
      >
        <DynamicGLBModel key={modelPath} path={modelPath} targetHeight={3.2} />
      </Suspense>

      {/* Selected Indicator Ring */}
      {selected && (
        <mesh position={[0, 0.05, 0]} rotation={[-Math.PI / 2, 0, 0]}>
          <ringGeometry args={[1.9, 2.25, 32]} />
          <meshBasicMaterial color="#f97316" side={THREE.DoubleSide} />
        </mesh>
      )}

      {/* Building Label with Model Icon */}
      <Html position={[0, 3.6, 0]} center distanceFactor={14}>
        <div className={`px-2.5 py-1 rounded-lg text-xs font-semibold whitespace-nowrap shadow-lg border transition-all pointer-events-none flex items-center gap-1.5 ${
          selected
            ? 'bg-orange text-white border-orange shadow-orange/30'
            : 'bg-slate-900/85 text-white border-white/10'
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
