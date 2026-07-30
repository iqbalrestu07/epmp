import { Suspense } from 'react';
import { Canvas } from '@react-three/fiber';
import { Environment, Float, Html, OrbitControls } from '@react-three/drei';
import FloatingMesh from './FloatingMesh';

function Loader() {
  return <Html center><div className="text-orange text-xs tracking-widest uppercase">Loading 3D...</div></Html>
}

export default function Scene3D({ scrollProgress }: { scrollProgress: React.MutableRefObject<number> }) {
  return (
    <div className="fixed inset-0 z-0 pointer-events-none">
      {/* We make the Canvas itself accept pointer events so OrbitControls can be dragged, but keep it out of the way for HTML clicks using z-index */}
      <Canvas camera={{ position: [0, 2, 8], fov: 45 }} className="pointer-events-auto cursor-grab active:cursor-grabbing">
        <fog attach="fog" args={['#0b0b0c', 5, 20]} />
        <ambientLight intensity={0.2} />
        <directionalLight position={[10, 10, 5]} intensity={1} color="#ffaa00" />
        
        {/* Adds continuous slow rotation and allows the user to click-drag the city with smooth damping! */}
        <OrbitControls 
          autoRotate 
          autoRotateSpeed={0.2} 
          enableZoom={false} 
          enablePan={false} 
          maxPolarAngle={Math.PI / 2} 
          minPolarAngle={Math.PI / 4} 
          enableDamping={true}
          dampingFactor={0.05}
        />

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
