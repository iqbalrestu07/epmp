import { useState } from 'react';
import { ChevronLeft, Building2, ArrowLeft, Box } from 'lucide-react';
import { Link } from 'react-router-dom';
import { usePropertys } from '../hooks';
import { useBuildings } from '../../building/hooks';
import { useFloors } from '../../floor/hooks';
import { useRooms } from '../../room/hooks';
import type { Property } from '../types';
import type { Building } from '../../building/types';
import type { Floor } from '../../floor/types';
import type { Room } from '../../room/types';
import { PropertyScene3D } from '../components/PropertyScene3D';

export default function PropertyInteractiveView() {
  const { data: propertyData, isLoading } = usePropertys({ page: 1, per_page: 100 });
  const properties: Property[] = propertyData?.data ?? [];

  const [selectedProperty, setSelectedProperty] = useState<Property | null>(null);
  const [selectedBuilding, setSelectedBuilding] = useState<Building | null>(null);
  const [selectedFloor, setSelectedFloor] = useState<Floor | null>(null);
  const [selectedRoom, setSelectedRoom] = useState<Room | null>(null);

  // Fetch buildings for selected property
  const { data: buildingData } = useBuildings(
    selectedProperty ? { property_id: selectedProperty.id, per_page: 100 } : undefined
  );
  const buildings: Building[] = buildingData?.data ?? [];

  // Fetch floors for selected building
  const { data: floorData } = useFloors(
    selectedBuilding ? { building_id: selectedBuilding.id, per_page: 100 } : undefined
  );
  const floors: Floor[] = floorData?.data ?? [];

  // Fetch rooms (filter by property_id)
  const { data: roomData } = useRooms(
    selectedProperty ? { per_page: 100 } : undefined
  );
  const allRooms: Room[] = roomData?.data ?? [];
  const rooms = allRooms.filter((r) => r.property_id === selectedProperty?.id);

  const handlePropertySelect = (p: Property) => {
    setSelectedProperty(p);
    setSelectedBuilding(null);
    setSelectedFloor(null);
    setSelectedRoom(null);
  };

  const handleBuildingClick = (b: Building) => {
    setSelectedBuilding(b);
    setSelectedFloor(null);
    setSelectedRoom(null);
  };

  const handleFloorClick = (f: Floor) => {
    setSelectedFloor((prev) => (prev?.id === f.id ? null : f));
    setSelectedRoom(null);
  };

  const handleRoomClick = (r: Room) => {
    setSelectedRoom(r);
  };

  const handleBackToBuildings = () => {
    setSelectedBuilding(null);
    setSelectedFloor(null);
    setSelectedRoom(null);
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-gray-400">Loading properties...</div>
      </div>
    );
  }

  return (
    <div className="animate-in fade-in duration-500">
      {/* Header & Property Selector */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-6 mb-6">
        <div className="flex items-center gap-4">
          <Link to="/dashboard/properties" className="p-2 hover:bg-gray-200 rounded-full transition-colors">
            <ChevronLeft size={24} />
          </Link>
          <div>
            <h1 className="text-3xl font-display">Interactive 3D Map</h1>
            <p className="text-gray-500">Drill-down: Property → Building → Floor → Room</p>
          </div>
        </div>

        {/* Property Tabs */}
        {properties.length > 0 && (
          <div className="flex bg-white p-1.5 rounded-xl border shadow-sm w-full md:w-auto overflow-x-auto">
            {properties.map((p) => (
              <button
                key={p.id}
                onClick={() => handlePropertySelect(p)}
                className={`flex items-center gap-2 px-6 py-2.5 rounded-lg text-sm font-medium transition-all whitespace-nowrap ${
                  selectedProperty?.id === p.id
                  ? 'bg-black text-white shadow-md'
                  : 'text-gray-600 hover:bg-gray-100'
                }`}
              >
                <Building2 size={16} className={selectedProperty?.id === p.id ? 'text-orange' : ''} />
                {p.name}
              </button>
            ))}
          </div>
        )}
      </div>

      {/* Breadcrumb / Drill-down navigation */}
      {selectedProperty && (
        <div className="flex items-center gap-2 mb-4 text-sm">
          <button
            onClick={() => { setSelectedProperty(null); handleBackToBuildings(); }}
            className="text-gray-500 hover:text-black"
          >
            All Properties
          </button>
          <span className="text-gray-300">/</span>
          <span className="font-medium">{selectedProperty.name}</span>
          {selectedBuilding && (
            <>
              <span className="text-gray-300">/</span>
              <button
                onClick={handleBackToBuildings}
                className="text-gray-500 hover:text-black"
              >
                {selectedBuilding.name}
              </button>
            </>
          )}
          {selectedFloor && (
            <>
              <span className="text-gray-300">/</span>
              <span className="font-medium text-orange-500">{selectedFloor.name}</span>
            </>
          )}
        </div>
      )}

      {properties.length === 0 ? (
        <div className="bg-white rounded-3xl border shadow-sm flex flex-col items-center justify-center text-gray-400 p-10 text-center min-h-[400px]">
          <Building2 size={64} className="mb-4 opacity-20" />
          <h3 className="text-xl font-medium mb-2">No Properties Yet</h3>
          <p>Create a property first to see the interactive 3D map.</p>
          <Link
            to="/dashboard/properties/new"
            className="mt-4 px-6 py-2 bg-black text-white rounded-lg text-sm font-medium hover:bg-gray-800 transition-colors"
          >
            Create Property
          </Link>
        </div>
      ) : !selectedProperty ? (
        <div className="bg-white rounded-3xl border shadow-sm flex flex-col items-center justify-center text-gray-400 p-10 text-center min-h-[400px]">
          <Building2 size={64} className="mb-4 opacity-20" />
          <h3 className="text-xl font-medium mb-2">Select a Property</h3>
          <p>Click on any property from the tabs above to view its 3D building map.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* LEFT: 3D Canvas */}
          <div className="col-span-1 lg:col-span-3 bg-white rounded-2xl border shadow-sm overflow-hidden relative" style={{ height: '600px' }}>
            {selectedBuilding && (
              <button
                onClick={handleBackToBuildings}
                className="absolute top-4 left-4 z-10 flex items-center gap-1 px-3 py-1.5 bg-white border rounded-lg text-sm font-medium hover:bg-gray-50 shadow-sm"
              >
                <ArrowLeft size={16} />
                Back to Buildings
              </button>
            )}

            {/* 3D Scene */}
            <PropertyScene3D
              buildings={buildings}
              floors={floors}
              rooms={rooms}
              selectedBuildingId={selectedBuilding?.id ?? null}
              selectedFloorId={selectedFloor?.id ?? null}
              onBuildingClick={handleBuildingClick}
              onFloorClick={handleFloorClick}
              onRoomClick={handleRoomClick}
            />

            {/* Legend */}
            <div className="absolute bottom-4 left-4 z-10 bg-white/90 backdrop-blur rounded-lg border px-4 py-2 shadow-sm">
              <div className="flex gap-4 text-xs font-medium">
                <div className="flex items-center gap-1.5">
                  <span className="w-3 h-3 rounded bg-green-500"></span> Available
                </div>
                <div className="flex items-center gap-1.5">
                  <span className="w-3 h-3 rounded bg-red-500"></span> Occupied
                </div>
                <div className="flex items-center gap-1.5">
                  <span className="w-3 h-3 rounded bg-orange-500"></span> Selected
                </div>
              </div>
            </div>
          </div>

          {/* RIGHT: Info Panel */}
          <div className="col-span-1 space-y-4">
            {/* Building info / Floor list */}
            {!selectedBuilding ? (
              <div className="bg-white rounded-2xl border shadow-sm p-5">
                <h3 className="font-bold text-lg mb-3">Buildings</h3>
                {buildings.length === 0 ? (
                  <p className="text-sm text-gray-400">No buildings yet. Add buildings to this property to see them in 3D.</p>
                ) : (
                  <div className="space-y-2">
                    {buildings.map((b) => (
                      <button
                        key={b.id}
                        onClick={() => handleBuildingClick(b)}
                        className="w-full flex items-center justify-between p-3 rounded-lg border hover:bg-gray-50 transition-colors text-left"
                      >
                        <div className="flex items-center gap-2">
                          <Box size={18} className="text-gray-400" />
                          <div>
                            <p className="font-medium text-sm">{b.name}</p>
                            <p className="text-xs text-gray-400">{b.total_floors} floors</p>
                          </div>
                        </div>
                        <ChevronLeft size={16} className="rotate-180 text-gray-300" />
                      </button>
                    ))}
                  </div>
                )}
              </div>
            ) : (
              <>
                {/* Floor list */}
                <div className="bg-white rounded-2xl border shadow-sm p-5">
                  <h3 className="font-bold text-lg mb-1">{selectedBuilding.name}</h3>
                  <p className="text-xs text-gray-400 mb-3">Select a floor to view rooms</p>
                  {floors.length === 0 ? (
                    <p className="text-sm text-gray-400">No floors yet for this building.</p>
                  ) : (
                    <div className="space-y-2">
                      {[...floors].sort((a, b) => b.floor_number - a.floor_number).map((f) => (
                        <button
                          key={f.id}
                          onClick={() => handleFloorClick(f)}
                          className={`w-full flex items-center justify-between p-3 rounded-lg border transition-colors text-left ${
                            selectedFloor?.id === f.id
                              ? 'border-orange-500 bg-orange-50'
                              : 'hover:bg-gray-50'
                          }`}
                        >
                          <div>
                            <p className="font-medium text-sm">{f.name}</p>
                            <p className="text-xs text-gray-400">Level {f.floor_number}</p>
                          </div>
                          <ChevronLeft size={16} className={`rotate-90 ${selectedFloor?.id === f.id ? 'text-orange-500' : 'text-gray-300'}`} />
                        </button>
                      ))}
                    </div>
                  )}
                </div>

                {/* Room detail */}
                {selectedFloor && (
                  <div className="bg-white rounded-2xl border shadow-sm p-5">
                    <h3 className="font-bold text-lg mb-1">{selectedFloor.name} Rooms</h3>
                    {rooms.filter((r) => r.floor === selectedFloor.floor_number).length === 0 ? (
                      <p className="text-sm text-gray-400">No rooms on this floor.</p>
                    ) : (
                      <div className="space-y-2">
                        {rooms
                          .filter((r) => r.floor === selectedFloor.floor_number)
                          .map((r) => (
                            <button
                              key={r.id}
                              onClick={() => handleRoomClick(r)}
                              className={`w-full flex items-center justify-between p-3 rounded-lg border transition-colors text-left ${
                                selectedRoom?.id === r.id
                                  ? 'border-orange-500 bg-orange-50'
                                  : 'hover:bg-gray-50'
                              }`}
                            >
                              <div>
                                <p className="font-medium text-sm">{r.name}</p>
                                <p className="text-xs text-gray-400">
                                  Cap: {r.capacity} · Rp {r.price.toLocaleString()}
                                </p>
                              </div>
                              <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${
                                r.is_available
                                  ? 'bg-green-100 text-green-700'
                                  : 'bg-red-100 text-red-700'
                              }`}>
                                {r.is_available ? 'Available' : 'Occupied'}
                              </span>
                            </button>
                          ))}
                      </div>
                    )}
                  </div>
                )}
              </>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
