import { useState } from "react";
import { Canvas } from "@react-three/fiber";
import { OrbitControls, Html, Environment } from "@react-three/drei";
import { useTenants } from "../../tenant/hooks";
import { useRooms } from "../../room/hooks";
import type { Tenant } from "../../tenant/types";
import { Button } from "@/components/ui/button";

interface Props {
  onSubmit: (data: any) => void;
  isSubmitting: boolean;
}

export function InteractiveContractPlanner({ onSubmit, isSubmitting }: Props) {
  const { data: tenantsData } = useTenants({ per_page: 100 });
  const { data: roomsData } = useRooms({ per_page: 100 });
  
  const tenantsList = Array.isArray(tenantsData?.data)
    ? tenantsData.data
    : Array.isArray(tenantsData)
    ? tenantsData
    : [];
  const roomsList = Array.isArray(roomsData?.data)
    ? roomsData.data
    : Array.isArray(roomsData)
    ? roomsData
    : [];
  const tenants = tenantsList;
  const rooms = roomsList.filter((r: any) => r.is_available);
  
  const [selectedTenant, setSelectedTenant] = useState<Tenant | null>(null);
  const [selectedRoomId, setSelectedRoomId] = useState<string | null>(null);
  
  const handleRoomClick = (roomId: string, e: any) => {
    e.stopPropagation();
    if (!selectedTenant) {
      alert("Please select a prospective tenant from the sidebar first!");
      return;
    }
    
    setSelectedRoomId(roomId);
  };

  const handleConfirmPlacement = () => {
    if (!selectedTenant || !selectedRoomId) return;
    
    // Default contract data based on room
    const room = rooms.find(r => r.id === selectedRoomId);
    
    onSubmit({ 
      tenant_id: selectedTenant.id, 
      room_id: selectedRoomId,
      property_id: room?.property_id || "",
      status: "Active",
      start_date: new Date().toISOString().split('T')[0],
      end_date: new Date(new Date().setFullYear(new Date().getFullYear() + 1)).toISOString().split('T')[0], // 1 year
      monthly_rent: room?.price || 0,
      deposit_amount: (room?.price || 0) * 2, // 2 months deposit
      terms: "Standard lease agreement."
    });
  };

  return (
    <div className="flex h-[600px] bg-slate-50 rounded-xl overflow-hidden border border-slate-200">
      {/* Sidebar: Prospective Tenants */}
      <div className="w-72 bg-white border-r border-slate-200 p-4 flex flex-col h-full z-10">
        <h3 className="font-bold text-slate-800 mb-4">Waitlist / Prospective Tenants</h3>
        <p className="text-xs text-slate-500 mb-4">Select a tenant and drop them into an available room to auto-generate a contract.</p>
        
        <div className="flex-1 overflow-y-auto space-y-2 custom-scrollbar pr-2">
          {tenants.length === 0 && (
            <p className="text-sm text-slate-500 text-center mt-4">No tenants found.</p>
          )}
          {tenants.map((tenant: any) => (
            <div 
              key={tenant.id}
              onClick={() => setSelectedTenant(tenant)}
              className={`p-3 rounded-xl border flex items-center gap-3 cursor-pointer transition-all ${
                selectedTenant?.id === tenant.id 
                  ? 'border-green-500 bg-green-50 shadow-sm' 
                  : 'border-slate-200 hover:border-green-300 hover:bg-slate-50'
              }`}
            >
              <div className="w-10 h-10 rounded-full bg-slate-200 flex items-center justify-center font-bold text-slate-600 shrink-0">
                {tenant.full_name.charAt(0)}
              </div>
              <div className="min-w-0">
                <p className="text-sm font-semibold truncate">{tenant.full_name}</p>
                <p className="text-xs text-slate-500 truncate">{tenant.phone || tenant.email}</p>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Main 3D Canvas area */}
      <div className="flex-1 relative">
        <div className="absolute top-4 left-4 right-4 z-10 flex flex-col gap-2 pointer-events-none">
          <div className="bg-white/90 backdrop-blur p-3 rounded-lg shadow-sm border border-slate-200 pointer-events-auto self-start">
             <p className="text-sm text-slate-700 font-medium">
               {selectedTenant 
                 ? `Click on an available green room below to place "${selectedTenant.full_name}"` 
                 : "1. Select a tenant from the sidebar."}
             </p>
          </div>
          
          {selectedRoomId && selectedTenant && (
            <div className="bg-white p-5 rounded-xl shadow-2xl border border-green-500 pointer-events-auto self-start max-w-sm">
              <h4 className="font-bold text-lg mb-2 text-slate-800">Confirm Placement</h4>
              <p className="text-sm text-slate-600 mb-4">
                You are about to draft a 1-year contract for <strong>{selectedTenant.full_name}</strong> in Room <strong>{rooms.find(r => r.id === selectedRoomId)?.name}</strong>.
              </p>
              <div className="flex gap-2">
                <Button onClick={handleConfirmPlacement} disabled={isSubmitting} className="bg-green-600 hover:bg-green-700 text-white flex-1">
                  {isSubmitting ? "Generating..." : "Generate Contract"}
                </Button>
                <Button onClick={() => setSelectedRoomId(null)} variant="outline">
                  Cancel
                </Button>
              </div>
            </div>
          )}
        </div>

        <Canvas shadows camera={{ position: [15, 15, 15], fov: 40 }}>
          <color attach="background" args={["#f1f5f9"]} />
          <ambientLight intensity={0.6} />
          <directionalLight castShadow position={[10, 20, 10]} intensity={1.5} shadow-mapSize={[1024, 1024]} />
          <Environment preset="city" />

          {/* Simple Floor Grid of Available Rooms */}
          <group position={[0, -0.5, 0]}>
            <mesh rotation={[-Math.PI / 2, 0, 0]} receiveShadow>
              <planeGeometry args={[30, 30]} />
              <meshStandardMaterial color="#e2e8f0" />
            </mesh>
            <gridHelper args={[30, 30, "#cbd5e1", "#e2e8f0"]} position={[0, 0.01, 0]} />
            
            {rooms.length === 0 ? (
              <Html center>
                <div className="bg-white/80 px-4 py-2 rounded shadow text-sm">No available rooms.</div>
              </Html>
            ) : (
              rooms.map((room, idx) => {
                const cols = Math.ceil(Math.sqrt(rooms.length));
                const col = idx % cols;
                const row = Math.floor(idx / cols);
                const x = (col - (cols - 1) / 2) * 4;
                const z = (row - (cols - 1) / 2) * 4;
                const isSelected = selectedRoomId === room.id;

                return (
                  <group key={room.id} position={[x, 0, z]}>
                    <mesh 
                      position={[0, 1, 0]} 
                      castShadow 
                      receiveShadow
                      onClick={(e) => handleRoomClick(room.id, e)}
                      onPointerOver={() => document.body.style.cursor = 'pointer'}
                      onPointerOut={() => document.body.style.cursor = 'auto'}
                    >
                      <boxGeometry args={[3.5, 2, 3.5]} />
                      <meshStandardMaterial 
                        color={isSelected ? "#22c55e" : "#4ade80"} 
                        roughness={0.2} 
                        emissive={isSelected ? "#22c55e" : "#000000"}
                        emissiveIntensity={0.2}
                      />
                    </mesh>
                    <Html position={[0, 2.5, 0]} center zIndexRange={[100, 0]}>
                      <div className="bg-slate-1000 backdrop-blur text-white text-xs px-2 py-1 rounded whitespace-nowrap pointer-events-none">
                        {room.name}
                      </div>
                    </Html>
                  </group>
                );
              })
            )}
          </group>
          
          <OrbitControls 
            enableDamping 
            dampingFactor={0.05} 
            maxPolarAngle={Math.PI / 2 - 0.1}
            minDistance={5}
            maxDistance={30}
          />
        </Canvas>
      </div>
    </div>
  );
}
