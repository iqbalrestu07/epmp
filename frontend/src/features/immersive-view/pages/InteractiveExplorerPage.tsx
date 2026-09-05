import { useState, useRef, Suspense } from "react";
import { Canvas, useFrame } from "@react-three/fiber";
import { OrbitControls, Html, Environment, ContactShadows } from "@react-three/drei";
import { useBuildings } from "../../building/hooks";
import { useFloors } from "../../floor/hooks";
import { useRooms } from "../../room/hooks";
import { Button } from "@/components/ui/button";
import { BuildingTable } from "../../building/components/BuildingTable";
import { useNavigate } from "react-router-dom";
import * as THREE from "three";

export function InteractiveExplorerPage() {
  const { data: bData, isLoading } = useBuildings({ per_page: 100 });
  const { data: fData } = useFloors({ per_page: 100 });
  const { data: rData } = useRooms({ per_page: 100 });
  const navigate = useNavigate();

  const buildings = bData?.data || [];
  const floors = fData?.data || [];
  const rooms = rData?.data || [];

  const [viewMode, setViewMode] = useState<"basic" | "interactive">("interactive");
  const [explodedView, setExplodedView] = useState(0); // 0 to 1
  const [selectedBuilding, setSelectedBuilding] = useState<string | null>(null);

  if (isLoading) return <div className="p-8 text-center text-slate-500">Loading Explorer...</div>;

  return (
    <div className="flex flex-col h-[calc(100vh-6rem)] -m-4 md:-m-8">
      {/* Top Header & Toggle — fixed height, not overlapping */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 px-6 py-4 bg-white border-b border-slate-200 z-10">
        <div className="flex items-center gap-4">
          <div>
            <h1 className="text-xl font-bold text-slate-900">
              Spatial Explorer
            </h1>
            <p className="text-sm text-slate-500">Property Digital Twin</p>
          </div>
          <div className="h-8 w-px bg-slate-200 mx-2"></div>
          <div className="flex bg-slate-100 p-1 rounded-md">
            <Button
              variant={viewMode === "basic" ? "default" : "ghost"}
              size="sm"
              onClick={() => setViewMode("basic")}
              className="text-xs"
            >
              Basic View (Data)
            </Button>
            <Button
              variant={viewMode === "interactive" ? "default" : "ghost"}
              size="sm"
              onClick={() => setViewMode("interactive")}
              className="text-xs"
            >
              Interactive 3D
            </Button>
          </div>
        </div>

        {viewMode === "interactive" && (
          <div className="flex items-center gap-4">
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
                Reset Focus
              </button>
            )}
          </div>
        )}
      </div>

      {/* Content area — fills remaining space, no overlap */}
      <div className="flex-1 w-full relative bg-slate-50 overflow-hidden">
        {viewMode === "basic" ? (
          <div className="p-6 h-full overflow-y-auto">
            <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200">
              <div className="flex justify-between items-center mb-6">
                <h2 className="text-lg font-bold text-slate-900">Buildings Data</h2>
                <Button onClick={() => navigate("/dashboard/buildings/new")}>+ Add Building</Button>
              </div>
              {buildings.length === 0 ? (
                <p className="text-center text-slate-500 py-8">No buildings found. Please add a building first.</p>
              ) : (
                <BuildingTable
                  data={buildings}
                  onRowClick={(row) => navigate(`/dashboard/buildings/${row.id}`)}
                />
              )}
            </div>
          </div>
        ) : (
          <Canvas shadows camera={{ position: [20, 15, 20], fov: 40 }}>
            <fog attach="fog" args={["#f8fafc", 20, 60]} />
            <color attach="background" args={["#f8fafc"]} /> {/* Changed to light slate background */}
            
            <ambientLight intensity={0.6} />
            <directionalLight
              castShadow
              position={[20, 30, 10]}
              intensity={1.5}
              shadow-mapSize={[2048, 2048]}
              shadow-bias={-0.0001}
            />
            <Suspense fallback={null}>
              <Environment preset="city" />
            </Suspense>

            {/* Grid and Ground */}
            <mesh rotation={[-Math.PI / 2, 0, 0]} receiveShadow>
              <planeGeometry args={[100, 100]} />
              <meshStandardMaterial color="#e2e8f0" roughness={0.8} />
            </mesh>
            <gridHelper args={[100, 20, "#cbd5e1", "#f1f5f9"]} position={[0, 0.01, 0]} />

            {buildings.length === 0 ? (
              <Html center>
                <div className="bg-white/90 p-4 rounded-xl shadow-xl text-center border border-slate-200 w-64 pointer-events-auto">
                  <h3 className="font-bold text-slate-800 mb-2">No Buildings Found</h3>
                  <p className="text-sm text-slate-500 mb-4">You need to add a building first to see it in the 3D space.</p>
                  <Button onClick={() => navigate("/dashboard/buildings/new")} size="sm" className="w-full">
                    Create Building
                  </Button>
                </div>
              </Html>
            ) : (
              <group position={[0, 0, 0]}>
                {buildings.map((b, idx) => {
                  const cols = Math.ceil(Math.sqrt(buildings.length));
                  const col = idx % cols;
                  const row = Math.floor(idx / cols);
                  const x = (col - (cols - 1) / 2) * 15;
                  const z = (row - (cols - 1) / 2) * 15;
                  
                  const bFloors = floors.filter(f => f.building_id === b.id);
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
            )}

            <OrbitControls 
              makeDefault
              enableDamping 
              dampingFactor={0.05} 
              maxPolarAngle={Math.PI / 2 - 0.05}
              minDistance={5}
              maxDistance={50}
            />
            
            <ContactShadows position={[0, 0.02, 0]} opacity={0.4} scale={100} blur={2} far={10} />
          </Canvas>
        )}
      </div>
    </div>
  );
}

function BuildingGroup({ building, floors, rooms, position, explodedView, isSelected, onClick }: any) {
  const groupRef = useRef<THREE.Group>(null);
  const floorHeight = 2;
  const explosionGap = 3; 

  return (
    <group ref={groupRef} position={position} onClick={(e) => {
      e.stopPropagation();
      onClick();
    }}>
      <mesh position={[0, 0.1, 0]} receiveShadow>
        <boxGeometry args={[12, 0.2, 12]} />
        <meshStandardMaterial color={isSelected ? "#3b82f6" : "#94a3b8"} />
      </mesh>

      <Html position={[0, (floors.length * floorHeight) + (floors.length * explosionGap * explodedView) + 2, 0]} center zIndexRange={[100, 0]}>
        <div className={`transition-all duration-300 pointer-events-none ${isSelected ? 'scale-110 opacity-100' : 'scale-90 opacity-70'}`}>
          <div className="bg-white/90 backdrop-blur-md text-slate-800 px-4 py-2 rounded-2xl border border-slate-200 whitespace-nowrap text-center shadow-xl">
            <h3 className="font-bold text-lg">{building.name}</h3>
            <p className="text-xs text-slate-500">{floors.length} Floors</p>
          </div>
          <div className="w-px h-8 bg-gradient-to-b from-slate-400 to-transparent mx-auto"></div>
        </div>
      </Html>

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
      ref.current.position.y = THREE.MathUtils.lerp(ref.current.position.y, position[1], 0.1);
    }
  });

  return (
    <group ref={ref} position={[position[0], position[1], position[2]]}>
      <mesh position={[0, height/2, 0]} castShadow receiveShadow>
        <boxGeometry args={[10, height - 0.1, 10]} />
        <meshStandardMaterial 
          color="#f1f5f9" 
          transparent 
          opacity={explodedView > 0.1 ? 0.9 : 1}
          roughness={0.2} 
        />
      </mesh>
      
      {explodedView > 0.3 && (
        <Html position={[-5.5, height/2, 0]} center>
          <div className="bg-orange/90 text-white text-xs font-bold px-2 py-1 rounded shadow-lg whitespace-nowrap">
            {floor.name}
          </div>
        </Html>
      )}

      {explodedView > 0.1 && rooms.map((r: any, rIdx: number) => {
        const angle = (rIdx / rooms.length) * Math.PI * 2;
        const radius = 4.2;
        const rx = Math.cos(angle) * radius;
        const rz = Math.sin(angle) * radius;
        const roomColor = r.status === 'Occupied'
          ? '#ef4444' // Red
          : r.status === 'Reserved'
          ? '#f59e0b' // Amber/Yellow
          : r.status === 'Maintenance'
          ? '#64748b' // Slate
          : '#22c55e'; // Green

        const badgeDotClass = r.status === 'Occupied'
          ? 'bg-red-500'
          : r.status === 'Reserved'
          ? 'bg-amber-500'
          : r.status === 'Maintenance'
          ? 'bg-slate-400'
          : 'bg-green-500';

        return (
          <group key={r.id} position={[rx, height/2, rz]}>
            <mesh castShadow>
              <boxGeometry args={[1, height - 0.2, 1]} />
              <meshStandardMaterial 
                color={roomColor} 
                emissive={roomColor}
                emissiveIntensity={isSelected ? 0.4 : 0.1}
              />
            </mesh>
            {isSelected && explodedView > 0.8 && (
              <Html position={[0, height/2 + 0.5, 0]} center zIndexRange={[90, 0]}>
                <div className="bg-white/95 backdrop-blur text-slate-800 text-[10px] px-2 py-0.5 rounded shadow border border-slate-200 whitespace-nowrap font-semibold flex items-center gap-1.5">
                  <span className={`w-2 h-2 inline-block rounded-full ${badgeDotClass}`}></span>
                  <span>{r.name}</span>
                  <span className="text-[9px] text-slate-400">({r.status || 'Available'})</span>
                </div>
              </Html>
            )}
          </group>
        );
      })}
    </group>
  );
}
