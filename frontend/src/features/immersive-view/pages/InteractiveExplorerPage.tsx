import { useState, useRef } from "react";
import { Canvas, useFrame } from "@react-three/fiber";
import { OrbitControls, Html, Environment, ContactShadows } from "@react-three/drei";
import { useBuildings } from "../../building/hooks";
import { useFloors } from "../../floor/hooks";
import { useRooms } from "../../room/hooks";
import * as THREE from "three";

export function InteractiveExplorerPage() {
  const { data: bData } = useBuildings({ per_page: 100 });
  const { data: fData } = useFloors({ per_page: 100 });
  const { data: rData } = useRooms({ per_page: 100 });

  const buildings = bData?.data || [];
  const floors = fData?.data || [];
  const rooms = rData?.data || [];

  const [explodedView, setExplodedView] = useState(0); // 0 to 1
  const [selectedBuilding, setSelectedBuilding] = useState<string | null>(null);

  return (
    <div className="flex flex-col h-[calc(100vh-6rem)] -m-4 md:-m-8">
      {/* Top HUD overlay */}
      <div className="absolute z-10 top-4 left-4 right-4 flex justify-between pointer-events-none">
        <div className="bg-white/80 backdrop-blur p-4 rounded-xl shadow-lg border border-white/40 pointer-events-auto">
          <h1 className="text-xl font-bold bg-gradient-to-r from-orange to-red-500 bg-clip-text text-transparent">
            Spatial Explorer
          </h1>
          <p className="text-sm text-slate-600">Immersive Digital Twin</p>
        </div>
        
        <div className="bg-white/80 backdrop-blur p-4 rounded-xl shadow-lg border border-white/40 pointer-events-auto flex items-center gap-4">
          <div className="flex flex-col">
            <label className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">
              Exploded View
            </label>
            <input 
              type="range" 
              min="0" max="1" step="0.01" 
              value={explodedView}
              onChange={(e) => setExplodedView(parseFloat(e.target.value))}
              className="accent-orange w-32"
            />
          </div>
          {selectedBuilding && (
            <button 
              onClick={() => setSelectedBuilding(null)}
              className="ml-2 text-xs bg-slate-200 hover:bg-slate-300 px-3 py-1.5 rounded-full font-medium transition-colors"
            >
              Reset Camera
            </button>
          )}
        </div>
      </div>

      <div className="flex-1 w-full bg-slate-900 relative">
        <Canvas shadows camera={{ position: [20, 15, 20], fov: 40 }}>
          <fog attach="fog" args={["#0f172a", 20, 60]} />
          <color attach="background" args={["#0f172a"]} />
          
          <ambientLight intensity={0.4} />
          <directionalLight
            castShadow
            position={[20, 30, 10]}
            intensity={1.5}
            shadow-mapSize={[2048, 2048]}
            shadow-bias={-0.0001}
          />
          <Environment preset="city" />

          {/* Grid and Ground */}
          <mesh rotation={[-Math.PI / 2, 0, 0]} receiveShadow>
            <planeGeometry args={[100, 100]} />
            <meshStandardMaterial color="#1e293b" roughness={0.8} />
          </mesh>
          <gridHelper args={[100, 20, "#334155", "#1e293b"]} position={[0, 0.01, 0]} />

          <group position={[0, 0, 0]}>
            {buildings.map((b, idx) => {
              // Layout buildings in a grid
              const cols = Math.ceil(Math.sqrt(buildings.length));
              const col = idx % cols;
              const row = Math.floor(idx / cols);
              const x = (col - (cols - 1) / 2) * 15;
              const z = (row - (cols - 1) / 2) * 15;
              
              const bFloors = floors.filter(f => f.building_id === b.id);
              // sort floors by floor_number ascending
              bFloors.sort((a, b) => a.floor_number - b.floor_number);

              return (
                <BuildingGroup 
                  key={b.id}
                  building={b}
                  floors={bFloors}
                  rooms={rooms}
                  position={[x, 0, z]}
                  explodedView={explodedView}
                  isSelected={selectedBuilding === b.id}
                  onClick={() => setSelectedBuilding(b.id === selectedBuilding ? null : b.id)}
                />
              );
            })}
          </group>

          <OrbitControls 
            makeDefault
            enableDamping 
            dampingFactor={0.05} 
            maxPolarAngle={Math.PI / 2 - 0.05}
            minDistance={5}
            maxDistance={50}
          />
          
          <ContactShadows position={[0, 0.02, 0]} opacity={0.5} scale={100} blur={2} far={10} />
        </Canvas>
      </div>
    </div>
  );
}

function BuildingGroup({ building, floors, rooms, position, explodedView, isSelected, onClick }: any) {
  const groupRef = useRef<THREE.Group>(null);
  
  // Floor height constant
  const floorHeight = 2;
  const explosionGap = 3; // Max gap between floors when slider is at 1

  return (
    <group ref={groupRef} position={position} onClick={(e) => {
      e.stopPropagation();
      onClick();
    }}>
      {/* Foundation / Land pad */}
      <mesh position={[0, 0.1, 0]} receiveShadow>
        <boxGeometry args={[12, 0.2, 12]} />
        <meshStandardMaterial color={isSelected ? "#3b82f6" : "#475569"} />
      </mesh>

      {/* Floating Spatial UI Label */}
      <Html position={[0, (floors.length * floorHeight) + (floors.length * explosionGap * explodedView) + 2, 0]} center zIndexRange={[100, 0]}>
        <div className={`transition-all duration-300 pointer-events-none ${isSelected ? 'scale-110 opacity-100' : 'scale-90 opacity-70'}`}>
          <div className="bg-black/40 backdrop-blur-md text-white px-4 py-2 rounded-2xl border border-white/20 whitespace-nowrap text-center shadow-2xl">
            <h3 className="font-bold text-lg">{building.name}</h3>
            <p className="text-xs text-white/70">{floors.length} Floors</p>
          </div>
          {/* Stem pointing down */}
          <div className="w-px h-8 bg-gradient-to-b from-white/50 to-transparent mx-auto"></div>
        </div>
      </Html>

      {/* Render Floors */}
      {floors.map((f: any, idx: number) => {
        const floorRooms = rooms.filter((r: any) => r.floor_id === f.id);
        const yPos = 0.2 + (idx * floorHeight) + (idx * explosionGap * explodedView);
        
        return (
          <FloorGroup 
            key={f.id} 
            floor={f} 
            rooms={floorRooms} 
            position={[0, yPos, 0]} 
            height={floorHeight}
            explodedView={explodedView}
            isSelected={isSelected}
          />
        );
      })}
    </group>
  );
}

function FloorGroup({ floor, rooms, position, height, explodedView, isSelected }: any) {
  const ref = useRef<THREE.Group>(null);

  useFrame(() => {
    if (ref.current) {
      // Smoothly animate to the target Y position
      ref.current.position.y = THREE.MathUtils.lerp(ref.current.position.y, position[1], 0.1);
    }
  });

  return (
    <group ref={ref} position={[position[0], position[1], position[2]]}>
      {/* Floor Base Slab */}
      <mesh position={[0, height/2, 0]} castShadow receiveShadow>
        <boxGeometry args={[10, height - 0.1, 10]} />
        <meshStandardMaterial 
          color="#f8fafc" 
          transparent 
          opacity={explodedView > 0.1 ? 0.9 : 1}
          roughness={0.2} 
        />
      </mesh>
      
      {/* Floor Label (Only visible when exploded) */}
      {explodedView > 0.3 && (
        <Html position={[-5.5, height/2, 0]} center>
          <div className="bg-orange/90 text-white text-xs font-bold px-2 py-1 rounded shadow-lg whitespace-nowrap">
            {floor.name}
          </div>
        </Html>
      )}

      {/* Render Rooms as little colored blocks on the edges if exploded */}
      {explodedView > 0.1 && rooms.map((r: any, rIdx: number) => {
        // Simple layout around the perimeter for visualization
        const angle = (rIdx / rooms.length) * Math.PI * 2;
        const radius = 4.2;
        const rx = Math.cos(angle) * radius;
        const rz = Math.sin(angle) * radius;
        const isAvailable = r.is_available;

        return (
          <group key={r.id} position={[rx, height/2, rz]}>
            <mesh castShadow>
              <boxGeometry args={[1, height - 0.2, 1]} />
              <meshStandardMaterial 
                color={isAvailable ? "#22c55e" : "#ef4444"} 
                emissive={isAvailable ? "#22c55e" : "#ef4444"}
                emissiveIntensity={isSelected ? 0.4 : 0.1}
              />
            </mesh>
            {/* Room HUD label */}
            {isSelected && explodedView > 0.8 && (
              <Html position={[0, height/2 + 0.5, 0]} center zIndexRange={[90, 0]}>
                <div className="bg-slate-900/80 backdrop-blur text-white text-[10px] px-1.5 py-0.5 rounded border border-white/10 whitespace-nowrap">
                  {r.name}
                  <span className={`ml-1 w-2 h-2 inline-block rounded-full ${isAvailable ? 'bg-green-500' : 'bg-red-500'}`}></span>
                </div>
              </Html>
            )}
          </group>
        );
      })}
    </group>
  );
}
