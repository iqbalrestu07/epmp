import { useState, useRef } from "react";
import { Canvas, useFrame } from "@react-three/fiber";
import { OrbitControls, Html, Environment, ContactShadows } from "@react-three/drei";
import { BuildingForm } from "./BuildingForm";
import type { CreateBuildingFormData } from "../schema";

interface Props {
  onSubmit: (data: CreateBuildingFormData) => void;
  isSubmitting: boolean;
}

export function InteractiveBuildingCreator({ onSubmit, isSubmitting }: Props) {
  const [buildingPlaced, setBuildingPlaced] = useState(false);
  const [position, setPosition] = useState<[number, number, number]>([0, 0, 0]);

  const handleGroundClick = (e: any) => {
    if (buildingPlaced) return;
    e.stopPropagation();
    setPosition([e.point.x, 0, e.point.z]);
    setBuildingPlaced(true);
  };

  return (
    <div className="relative w-full h-[600px] bg-slate-50 rounded-xl overflow-hidden border border-slate-200">
      <div className="absolute top-4 left-4 z-10 bg-white/90 backdrop-blur p-4 rounded-lg shadow-sm border border-slate-200 text-sm max-w-xs pointer-events-none">
        <h3 className="font-bold text-slate-800 text-base mb-1">Interactive 3D Creator</h3>
        <p className="text-slate-600">
          {buildingPlaced 
            ? "Great! Now fill out the building details in the popup." 
            : "Click anywhere on the grid to drop a new building plot."}
        </p>
      </div>

      <Canvas shadows camera={{ position: [15, 12, 15], fov: 45 }}>
        <color attach="background" args={["#f8fafc"]} />
        <ambientLight intensity={0.5} />
        <directionalLight
          castShadow
          position={[10, 20, 10]}
          intensity={1}
          shadow-mapSize={[1024, 1024]}
        />
        
        <Environment preset="city" />

        {/* Ground Plane */}
        <mesh 
          rotation={[-Math.PI / 2, 0, 0]} 
          position={[0, -0.1, 0]} 
          receiveShadow 
          onClick={handleGroundClick}
        >
          <planeGeometry args={[50, 50]} />
          <meshStandardMaterial color="#e2e8f0" />
        </mesh>
        
        <gridHelper args={[50, 50, "#cbd5e1", "#f1f5f9"]} position={[0, -0.09, 0]} />

        {buildingPlaced && (
          <group position={position}>
            {/* Animated dropping building */}
            <DroppingBuilding />
            
            {/* Form attached to the building */}
            <Html position={[0, 4.5, 0]} center zIndexRange={[100, 0]}>
              <div className="bg-white p-5 rounded-xl shadow-2xl border border-slate-200 w-[350px] pointer-events-auto transition-transform">
                <div className="flex justify-between items-center mb-4">
                  <h3 className="font-semibold text-lg text-slate-800">Building Setup</h3>
                  <button 
                    onClick={(e) => {
                      e.stopPropagation();
                      setBuildingPlaced(false);
                    }}
                    className="text-slate-400 hover:text-red-500 transition-colors p-1"
                    title="Cancel"
                  >
                    ✕
                  </button>
                </div>
                <div className="max-h-[400px] overflow-y-auto pr-2 custom-scrollbar">
                  <BuildingForm 
                    onSubmit={onSubmit} 
                    isSubmitting={isSubmitting} 
                  />
                </div>
              </div>
            </Html>
          </group>
        )}

        <OrbitControls 
          enableDamping 
          dampingFactor={0.05} 
          maxPolarAngle={Math.PI / 2 - 0.1}
          minDistance={5}
          maxDistance={35}
        />
        <ContactShadows position={[0, 0, 0]} opacity={0.4} scale={50} blur={2} far={10} />
      </Canvas>
    </div>
  );
}

function DroppingBuilding() {
  const meshRef = useRef<any>();
  
  useFrame(() => {
    if (meshRef.current) {
      // Simple drop and bounce animation
      meshRef.current.position.y = Math.max(1.5, meshRef.current.position.y - 0.3);
    }
  });

  return (
    <group ref={meshRef} position={[0, 15, 0]}>
      {/* Main Building Block */}
      <mesh castShadow receiveShadow>
        <boxGeometry args={[3, 3, 3]} />
        <meshStandardMaterial color="#3b82f6" roughness={0.3} metalness={0.1} />
      </mesh>
      {/* Roof accent */}
      <mesh position={[0, 1.55, 0]} castShadow>
        <boxGeometry args={[2.8, 0.1, 2.8]} />
        <meshStandardMaterial color="#1e3a8a" />
      </mesh>
    </group>
  );
}
