import { useState } from "react";
import { Canvas } from "@react-three/fiber";
import { OrbitControls, Html, Environment, ContactShadows } from "@react-three/drei";
import { useAssets } from "../../asset/hooks";
import { useRooms } from "../../room/hooks";
import type { Asset } from "../../asset/types";

interface Props {
  onSubmit: (data: { asset_id: string; room_id: string }) => void;
  isSubmitting: boolean;
}

export function InteractiveAssetAssignment({ onSubmit }: Props) {
  const { data: assetsData } = useAssets({ per_page: 100 });
  const { data: roomsData } = useRooms({ per_page: 100 });
  
  const rawAssets = Array.isArray(assetsData?.data)
    ? assetsData.data
    : Array.isArray(assetsData)
    ? assetsData
    : [];
  const rawRooms = Array.isArray(roomsData?.data)
    ? roomsData.data
    : Array.isArray(roomsData)
    ? roomsData
    : [];
  const assets = rawAssets.filter((a: any) => a.status === "Available");
  const rooms = rawRooms;
  
  const [selectedRoomId, setSelectedRoomId] = useState<string>("");
  const [selectedAsset, setSelectedAsset] = useState<Asset | null>(null);
  
  const handleFloorClick = (e: any) => {
    e.stopPropagation();
    if (!selectedRoomId) {
      alert("Please select a room first!");
      return;
    }
    if (!selectedAsset) {
      alert("Please select an available asset from the left sidebar to place it.");
      return;
    }
    
    // Simulate placing the asset and submitting
    onSubmit({ asset_id: selectedAsset.id, room_id: selectedRoomId });
    setSelectedAsset(null); // reset
  };

  return (
    <div className="flex h-[600px] bg-slate-50 rounded-xl overflow-hidden border border-slate-200">
      {/* Sidebar: Available Assets */}
      <div className="w-64 bg-white border-r border-slate-200 p-4 flex flex-col h-full">
        <h3 className="font-bold text-slate-800 mb-4">Available Assets</h3>
        <div className="flex-1 overflow-y-auto space-y-2 custom-scrollbar">
          {assets.length === 0 && (
            <p className="text-sm text-slate-500 text-center mt-4">No available assets.</p>
          )}
          {assets.map(asset => (
            <div 
              key={asset.id}
              onClick={() => setSelectedAsset(asset)}
              className={`p-3 rounded-lg border cursor-pointer transition-all ${
                selectedAsset?.id === asset.id 
                  ? 'border-orange bg-orange/5 shadow-sm' 
                  : 'border-slate-200 hover:border-orange/50 hover:bg-slate-50'
              }`}
            >
              <p className="text-sm font-semibold">{asset.name}</p>
              <p className="text-xs text-slate-500">{asset.category}</p>
            </div>
          ))}
        </div>
      </div>

      {/* Main 3D Canvas area */}
      <div className="flex-1 relative">
        <div className="absolute top-4 left-4 right-4 z-10 flex gap-4 pointer-events-none">
          <div className="bg-white/90 backdrop-blur p-3 rounded-lg shadow-sm border border-slate-200 pointer-events-auto">
            <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">
              Select Target Room
            </label>
            <select
              value={selectedRoomId}
              onChange={(e) => setSelectedRoomId(e.target.value)}
              className="w-48 text-sm border-none bg-slate-100 rounded px-2 py-1 focus:ring-2 focus:ring-orange/30"
            >
              <option value="">-- Choose Room --</option>
              {rooms.map(r => (
                <option key={r.id} value={r.id}>{r.name}</option>
              ))}
            </select>
          </div>
          
          <div className="bg-white/90 backdrop-blur p-3 rounded-lg shadow-sm border border-slate-200 pointer-events-auto flex-1 flex items-center">
             <p className="text-sm text-slate-600 font-medium">
               {selectedAsset 
                 ? `Now click on the 3D floor below to place "${selectedAsset.name}"` 
                 : "Select an asset from the sidebar to begin placing."}
             </p>
          </div>
        </div>

        <Canvas shadows camera={{ position: [8, 8, 8], fov: 45 }}>
          <color attach="background" args={["#f8fafc"]} />
          <ambientLight intensity={0.5} />
          <directionalLight castShadow position={[10, 20, 10]} intensity={1} shadow-mapSize={[1024, 1024]} />
          <Environment preset="city" />

          {/* 3D Room Mockup */}
          {selectedRoomId ? (
            <group>
              {/* Floor */}
              <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, 0, 0]} receiveShadow onClick={handleFloorClick}>
                <planeGeometry args={[10, 10]} />
                <meshStandardMaterial color="#e2e8f0" />
                <gridHelper args={[10, 10, "#cbd5e1", "#f1f5f9"]} rotation={[Math.PI / 2, 0, 0]} />
              </mesh>
              {/* Walls */}
              <mesh position={[0, 1.5, -5]} receiveShadow castShadow>
                <boxGeometry args={[10, 3, 0.2]} />
                <meshStandardMaterial color="#f1f5f9" />
              </mesh>
              <mesh position={[-5, 1.5, 0]} rotation={[0, Math.PI / 2, 0]} receiveShadow castShadow>
                <boxGeometry args={[10, 3, 0.2]} />
                <meshStandardMaterial color="#f1f5f9" />
              </mesh>
            </group>
          ) : (
            <Html center>
              <div className="bg-black/50 text-white px-4 py-2 rounded-lg backdrop-blur">
                Please select a room to view the 3D planner.
              </div>
            </Html>
          )}
          
          <OrbitControls 
            enableDamping 
            dampingFactor={0.05} 
            maxPolarAngle={Math.PI / 2 - 0.1}
            minDistance={3}
            maxDistance={20}
          />
          <ContactShadows position={[0, 0, 0]} opacity={0.4} scale={20} blur={2} far={10} />
        </Canvas>
      </div>
    </div>
  );
}
