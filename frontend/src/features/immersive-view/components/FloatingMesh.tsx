import { useRef } from 'react';
import { useFrame } from '@react-three/fiber';
import * as THREE from 'three';
import { Model as CityNightModel } from './CityNight';

export default function FloatingMesh({ scrollProgress }: { scrollProgress: React.MutableRefObject<number> }) {
  const groupRef = useRef<THREE.Group>(null);

  // Smooth mouse tracking
  const targetRotation = useRef({ x: 0, y: 0 });

  useFrame((state) => {
    if (!groupRef.current) return;

    // 1. Get mouse coords (-1 to 1)
    const { x, y } = state.pointer;
    
    // 2. Calculate target rotation (subtle for the city)
    targetRotation.current.y = (x * Math.PI) / 8;
    targetRotation.current.x = (y * Math.PI) / 16;

    // 3. Smooth lerp to target
    groupRef.current.rotation.x = THREE.MathUtils.lerp(groupRef.current.rotation.x, targetRotation.current.x, 0.05);
    groupRef.current.rotation.y = THREE.MathUtils.lerp(groupRef.current.rotation.y, targetRotation.current.y, 0.05);

    // 4. Scrollytelling effect: zoom out slightly and move down
    const scale = 1 - scrollProgress.current * 0.2;
    groupRef.current.scale.set(scale, scale, scale);
    groupRef.current.position.y = -2 - scrollProgress.current * 4;
  });

  return (
    <group ref={groupRef} position={[0, -2, 0]}>
      {/* Add a dramatic spotlight pointing at the city */}
      <spotLight position={[0, 10, 0]} intensity={2} color="#ffaa00" angle={0.5} penumbra={1} />
      <CityNightModel scale={1.7} />
    </group>
  );
}
