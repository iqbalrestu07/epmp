import { Suspense } from 'react';
import { Canvas } from '@react-three/fiber';
import { Environment, Float, Html } from '@react-three/drei';
import FloatingMesh from './FloatingMesh';

function Loader() {
  return <Html center><div className="text-orange text-xs tracking-widest uppercase">Loading 3D...</div></Html>
}

export default function Scene3D({ scrollProgress }: { scrollProgress: React.MutableRefObject<number> }) {
  return (
    <div className="fixed inset-0 z-0 pointer-events-none">
      <Canvas camera={{ position: [1, 4, 1], fov: 85 }}>
        <fog attach="fog" args={['#0b0b0c', 5, 20]} />
        <ambientLight intensity={0.2} />
        <directionalLight position={[10, 10, 5]} intensity={1} color="#ffaa00" />
        
        <Suspense fallback={<Loader />}>
          <Float speed={1.5} rotationIntensity={0.2} floatIntensity={0.5}>
            <FloatingMesh scrollProgress={scrollProgress} />
          </Float>
          <Environment preset="city" />
        </Suspense>
        
      </Canvas>
    </div>
  );
}
