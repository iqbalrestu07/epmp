import { useState, useEffect, useMemo, useRef } from 'react';
import {
  ChevronLeft,
  Building2,
  ArrowLeft,
  Box,
  Maximize2,
  Minimize2,
  Sparkles,
  Layers,
  RefreshCw,
  Globe2,
  PanelRightClose,
  PanelRightOpen,
  DoorOpen,
} from 'lucide-react';
import { Link } from 'react-router-dom';
import { usePropertys } from '../hooks';
import { useBuildings } from '../../building/hooks';
import { useFloors } from '../../floor/hooks';
import { useRooms } from '../../room/hooks';
import type { Property } from '../types';
import type { Building } from '../../building/types';
import type { Floor } from '../../floor/types';
import type { Room } from '../../room/types';
import { PropertyScene3D, type PropertyGroup } from '../components/PropertyScene3D';

export default function PropertyInteractiveView() {
  const containerRef = useRef<HTMLDivElement>(null);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [isInspectorOpen, setIsInspectorOpen] = useState(true);

  // Fetch all properties
  const { data: propertyData, isLoading: isLoadingProps } = usePropertys({ page: 1, per_page: 100 });
  const properties: Property[] = propertyData?.data ?? [];

  // Fetch all buildings across properties
  const { data: buildingData, isLoading: isLoadingBuildings } = useBuildings({ per_page: 200 });
  const allBuildings: Building[] = Array.isArray(buildingData?.data)
    ? buildingData.data
    : Array.isArray(buildingData)
    ? buildingData
    : [];

  const [selectedProperty, setSelectedProperty] = useState<Property | null>(null);
  const [selectedBuilding, setSelectedBuilding] = useState<Building | null>(null);
  const [selectedFloor, setSelectedFloor] = useState<Floor | null>(null);
  const [selectedRoom, setSelectedRoom] = useState<Room | null>(null);

  // Dynamic 3D model style overrides per building (tower vs hotel)
  const [modelOverrides, setModelOverrides] = useState<Record<string, 'building' | 'hotel'>>({});

  // Group buildings by property
  const propertyGroups: PropertyGroup[] = useMemo(() => {
    return properties.map((prop) => ({
      property: prop,
      buildings: allBuildings.filter((b) => b.property_id === prop.id),
    }));
  }, [properties, allBuildings]);

  // Active buildings for the currently selected property (or all if none selected)
  const activeBuildings: Building[] = useMemo(() => {
    if (selectedProperty) {
      return allBuildings.filter((b) => b.property_id === selectedProperty.id);
    }
    return allBuildings;
  }, [selectedProperty, allBuildings]);

  // Fetch floors for selected building
  const { data: floorData } = useFloors(
    selectedBuilding ? { building_id: selectedBuilding.id, per_page: 100 } : undefined
  );
  const floors: Floor[] = Array.isArray(floorData?.data)
    ? floorData.data
    : Array.isArray(floorData)
    ? floorData
    : [];

  // Fetch rooms
  const { data: roomData } = useRooms(
    selectedProperty ? { per_page: 100 } : undefined
  );
  const allRooms: Room[] = Array.isArray(roomData?.data)
    ? roomData.data
    : Array.isArray(roomData)
    ? roomData
    : [];
  const rooms = selectedBuilding
    ? allRooms.filter((r) => floors.some((f) => f.id === r.floor_id))
    : allRooms.filter((r) => r.property_id === selectedProperty?.id);

  // Fullscreen toggle handler
  const toggleFullscreen = () => {
    if (!document.fullscreenElement) {
      containerRef.current?.requestFullscreen?.().catch(() => {
        setIsFullscreen(true);
      });
      setIsFullscreen(true);
    } else {
      document.exitFullscreen?.().catch(() => {});
      setIsFullscreen(false);
    }
  };

  useEffect(() => {
    const handleFsChange = () => {
      setIsFullscreen(!!document.fullscreenElement);
    };
    document.addEventListener('fullscreenchange', handleFsChange);
    return () => document.removeEventListener('fullscreenchange', handleFsChange);
  }, []);

  const handleToggleBuildingModel = (b: Building, e: React.MouseEvent) => {
    e.stopPropagation();
    const current = modelOverrides[b.id] || (b.name.toLowerCase().includes('hotel') ? 'hotel' : 'building');
    const next = current === 'hotel' ? 'building' : 'hotel';
    setModelOverrides((prev) => ({ ...prev, [b.id]: next }));
    try {
      localStorage.setItem(`building_model_pref_${b.name}`, next === 'hotel' ? 'hotel_resort' : 'office_tower');
    } catch {
      // ignore
    }
  };

  const handlePropertySelect = (p: Property | null) => {
    setSelectedProperty(p);
    setSelectedBuilding(null);
    setSelectedFloor(null);
    setSelectedRoom(null);
  };

  const handleBuildingClick = (b: Building) => {
    // If we're in all properties mode, also set the selected property
    const parentProp = properties.find((p) => p.id === b.property_id);
    if (parentProp && (!selectedProperty || selectedProperty.id !== parentProp.id)) {
      setSelectedProperty(parentProp);
    }
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

  if (isLoadingProps || isLoadingBuildings) {
    return (
      <div className="flex items-center justify-center min-h-[450px]">
        <div className="text-slate-400 flex items-center gap-2">
          <RefreshCw className="animate-spin" size={18} /> Loading spatial property map...
        </div>
      </div>
    );
  }

  return (
    <div
      ref={containerRef}
      className={`transition-all ${
        isFullscreen
          ? 'fixed inset-0 z-50 bg-slate-900 p-6 flex flex-col justify-between overflow-hidden'
          : 'animate-in fade-in duration-500 space-y-5'
      }`}
    >
      {/* Top Header & Property Navigation */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          {!isFullscreen && (
            <Link
              to="/dashboard/properties"
              className="p-2 hover:bg-slate-200 rounded-full transition-colors text-slate-700"
            >
              <ChevronLeft size={24} />
            </Link>
          )}
          <div>
            <div className="flex items-center gap-2">
              <h1 className={`font-display font-bold ${isFullscreen ? 'text-2xl text-white' : 'text-3xl text-slate-900'}`}>
                Spatial 3D Digital Twin
              </h1>
              <span className="px-2.5 py-0.5 rounded-full text-xs font-semibold bg-orange/10 text-orange border border-orange/20">
                WebGL 3D
              </span>
            </div>
            <p className={`text-xs mt-0.5 ${isFullscreen ? 'text-slate-400' : 'text-slate-500'}`}>
              Interactive drill-down grouped by property: Property → Building → Floor → Room
            </p>
          </div>
        </div>

        {/* Property Selector Tabs (Includes "All Properties (Grouped)" and each individual property) */}
        {properties.length > 0 && (
          <div className="flex bg-white/95 backdrop-blur p-1 rounded-xl border border-slate-200 shadow-sm w-full md:w-auto overflow-x-auto gap-1">
            {/* All Properties Grouped Button */}
            <button
              onClick={() => handlePropertySelect(null)}
              className={`flex items-center gap-1.5 px-4 py-2 rounded-lg text-xs font-semibold transition-all whitespace-nowrap ${
                selectedProperty === null
                  ? 'bg-slate-900 text-white shadow-md'
                  : 'text-slate-600 hover:bg-slate-100'
              }`}
            >
              <Globe2 size={14} className={selectedProperty === null ? 'text-orange' : ''} />
              All Properties ({properties.length})
            </button>

            {properties.map((p) => {
              const bCount = allBuildings.filter((b) => b.property_id === p.id).length;
              const isSelected = selectedProperty?.id === p.id;
              return (
                <button
                  key={p.id}
                  onClick={() => handlePropertySelect(p)}
                  className={`flex items-center gap-1.5 px-3.5 py-2 rounded-lg text-xs font-semibold transition-all whitespace-nowrap ${
                    isSelected
                      ? 'bg-slate-900 text-white shadow-md'
                      : 'text-slate-600 hover:bg-slate-100'
                  }`}
                >
                  <Building2 size={14} className={isSelected ? 'text-orange' : ''} />
                  <span>{p.name}</span>
                  <span className={`text-[10px] px-1.5 py-0.2 rounded-full ${
                    isSelected ? 'bg-orange/20 text-orange' : 'bg-slate-200 text-slate-600'
                  }`}>
                    {bCount}
                  </span>
                </button>
              );
            })}
          </div>
        )}
      </div>

      {/* Breadcrumb Navigation */}
      <div className="flex items-center justify-between gap-4 text-xs">
        <div className={`flex items-center gap-1.5 ${isFullscreen ? 'text-slate-400' : 'text-slate-500'}`}>
          <button
            onClick={() => handlePropertySelect(null)}
            className="hover:text-orange transition-colors font-medium"
          >
            All Properties (Grouped)
          </button>
          {selectedProperty && (
            <>
              <span>/</span>
              <button
                onClick={() => handlePropertySelect(selectedProperty)}
                className="hover:text-orange transition-colors font-semibold text-slate-800 dark:text-slate-200"
              >
                {selectedProperty.name}
              </button>
            </>
          )}
          {selectedBuilding && (
            <>
              <span>/</span>
              <button
                onClick={handleBackToBuildings}
                className="hover:text-orange transition-colors font-medium text-slate-700"
              >
                {selectedBuilding.name}
              </button>
            </>
          )}
          {selectedFloor && (
            <>
              <span>/</span>
              <button
                onClick={() => handleFloorClick(selectedFloor)}
                className="hover:text-orange transition-colors font-medium text-slate-700"
              >
                {selectedFloor.name}
              </button>
            </>
          )}
          {selectedRoom && (
            <>
              <span>/</span>
              <span className="font-bold text-orange">{selectedRoom.name}</span>
            </>
          )}
        </div>

        {/* Dynamic Create Button — changes based on selection level */}
        <div className="flex items-center gap-2">
          {selectedFloor ? (
            <Link
              to={`/dashboard/rooms/new?property_id=${selectedProperty?.id ?? ''}&floor_id=${selectedFloor.id}`}
              className="px-3 py-1.5 rounded-lg bg-orange text-white font-medium hover:bg-orange/90 transition-colors flex items-center gap-1"
            >
              <DoorOpen size={13} />
              + New Room
            </Link>
          ) : selectedBuilding ? (
            <Link
              to={`/dashboard/floors/new?building_id=${selectedBuilding.id}`}
              className="px-3 py-1.5 rounded-lg bg-orange text-white font-medium hover:bg-orange/90 transition-colors flex items-center gap-1"
            >
              <Layers size={13} />
              + New Floor
            </Link>
          ) : (
            <Link
              to={`/dashboard/buildings/new${selectedProperty ? `?property_id=${selectedProperty.id}` : ''}`}
              className="px-3 py-1.5 rounded-lg bg-orange text-white font-medium hover:bg-orange/90 transition-colors flex items-center gap-1"
            >
              <Box size={13} />
              + New Building
            </Link>
          )}
        </div>
      </div>

      {properties.length === 0 ? (
        <div className="bg-white rounded-3xl border shadow-sm flex flex-col items-center justify-center text-slate-400 p-12 text-center min-h-[500px]">
          <Building2 size={64} className="mb-4 opacity-20" />
          <h3 className="text-xl font-medium mb-2 text-slate-800">No Properties Found</h3>
          <p className="text-sm text-slate-500 max-w-sm mb-4">
            Create a property first to view interactive spatial 3D buildings.
          </p>
          <Link
            to="/dashboard/properties/new"
            className="px-6 py-2.5 bg-slate-900 text-white rounded-xl text-sm font-semibold hover:bg-gray-800 transition-colors"
          >
            Create First Property
          </Link>
        </div>
      ) : (
        /* Main Spatial Grid & Info Split Layout */
        <div className={`grid grid-cols-1 ${isInspectorOpen ? 'lg:grid-cols-4' : 'lg:grid-cols-1'} gap-5 ${isFullscreen ? 'flex-1 min-h-0' : ''}`}>
          {/* 3D Canvas Area (Expanded size: 800px / Fullscreen / Wide mode) */}
          <div
            className={`col-span-1 ${isInspectorOpen ? 'lg:col-span-3' : 'lg:col-span-1'} bg-slate-950 rounded-2xl border border-slate-800 shadow-xl overflow-hidden relative transition-all duration-300 ${
              isFullscreen ? 'h-full' : 'h-[800px]'
            }`}
          >
            {/* Top Toolbar Overlay inside 3D Canvas */}
            <div className="absolute top-4 left-4 right-4 z-20 flex items-center justify-between pointer-events-none">
              <div className="flex items-center gap-2 pointer-events-auto">
                {selectedBuilding ? (
                  <button
                    onClick={handleBackToBuildings}
                    className="flex items-center gap-1.5 px-3 py-1.5 bg-slate-900/90 text-white hover:bg-slate-800 border border-slate-700 rounded-xl text-xs font-semibold backdrop-blur shadow-md transition-all"
                  >
                    <ArrowLeft size={14} />
                    Back to Buildings
                  </button>
                ) : (
                  <div className="px-3.5 py-1.5 bg-slate-900/90 text-white border border-slate-700 rounded-xl text-xs font-semibold backdrop-blur shadow-md flex items-center gap-2">
                    <span className="w-2 h-2 rounded-full bg-green-400 animate-pulse"></span>
                    <span>{selectedProperty ? selectedProperty.name : 'All Properties (Grouped Clusters)'}</span>
                    <span className="text-slate-400 text-[11px]">({activeBuildings.length} buildings)</span>
                  </div>
                )}
              </div>

              {/* Action Buttons: Toggle Inspector & Fullscreen */}
              <div className="flex items-center gap-2 pointer-events-auto">
                <button
                  onClick={() => setIsInspectorOpen((prev) => !prev)}
                  className="px-3 py-1.5 bg-slate-900/90 hover:bg-slate-800 text-white border border-slate-700 rounded-xl text-xs font-semibold backdrop-blur shadow-md flex items-center gap-1.5 transition-all"
                  title={isInspectorOpen ? 'Collapse sidebar (Focus 3D Canvas)' : 'Show details inspector'}
                >
                  {isInspectorOpen ? (
                    <>
                      <PanelRightClose size={14} className="text-orange" />
                      <span className="hidden sm:inline">Focus Canvas</span>
                    </>
                  ) : (
                    <>
                      <PanelRightOpen size={14} className="text-orange" />
                      <span className="hidden sm:inline">Show Details</span>
                    </>
                  )}
                </button>

                <button
                  onClick={toggleFullscreen}
                  className="px-3 py-1.5 bg-slate-900/90 hover:bg-slate-800 text-white border border-slate-700 rounded-xl text-xs font-semibold backdrop-blur shadow-md flex items-center gap-1.5 transition-all"
                  title={isFullscreen ? 'Exit Fullscreen' : 'Enter Fullscreen'}
                >
                  {isFullscreen ? (
                    <>
                      <Minimize2 size={14} className="text-orange" />
                      <span>Exit Fullscreen</span>
                    </>
                  ) : (
                    <>
                      <Maximize2 size={14} className="text-orange" />
                      <span>Fullscreen View</span>
                    </>
                  )}
                </button>
              </div>
            </div>

            {/* 3D Scene Component */}
            <PropertyScene3D
              buildings={activeBuildings}
              floors={floors}
              rooms={rooms}
              selectedBuildingId={selectedBuilding?.id ?? null}
              selectedFloorId={selectedFloor?.id ?? null}
              selectedRoomId={selectedRoom?.id ?? null}
              onBuildingClick={handleBuildingClick}
              onFloorClick={handleFloorClick}
              onRoomClick={handleRoomClick}
              propertyGroups={propertyGroups}
              isAllPropertiesMode={selectedProperty === null}
              onPropertyClick={(p) => handlePropertySelect(p)}
              buildingModelOverrides={modelOverrides}
            />

            {/* Bottom Legend Overlay */}
            <div className="absolute bottom-4 left-4 z-20 bg-slate-900/90 text-white backdrop-blur rounded-xl border border-slate-800 px-4 py-2 shadow-lg">
              <div className="flex items-center gap-4 text-xs font-medium">
                <div className="flex items-center gap-1.5">
                  <span className="w-2.5 h-2.5 rounded bg-green-500"></span> Available
                </div>
                <div className="flex items-center gap-1.5">
                  <span className="w-2.5 h-2.5 rounded bg-red-500"></span> Occupied
                </div>
                <div className="flex items-center gap-1.5">
                  <span className="w-2.5 h-2.5 rounded bg-orange"></span> Selected
                </div>
                <div className="border-l border-slate-700 pl-3 flex items-center gap-2 text-slate-400 text-[11px]">
                  <span>🏢 = Office Tower</span>
                  <span>🏨 = Hotel/Resort</span>
                </div>
              </div>
            </div>
          </div>

          {/* RIGHT: Properties & Buildings Inspector Panel */}
          {isInspectorOpen && (
            <div className="col-span-1 space-y-4 overflow-y-auto max-h-[800px] pr-1 animate-in fade-in slide-in-from-right-2 duration-200">
            {/* GROUPED LIST OF BUILDINGS BY PROPERTY */}
            {!selectedBuilding ? (
              <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-5 space-y-4">
                <div className="flex items-center justify-between border-b pb-3">
                  <div>
                    <h3 className="font-bold text-slate-900 text-base">Buildings Grouped by Property</h3>
                    <p className="text-xs text-slate-500 mt-0.5">Click any building to inspect floors & rooms</p>
                  </div>
                </div>

                {propertyGroups.map((group) => {
                  const isCurrentProp = selectedProperty?.id === group.property.id;
                  return (
                    <div
                      key={group.property.id}
                      className={`rounded-xl border p-3 transition-all ${
                        isCurrentProp
                          ? 'border-orange/60 bg-orange/5'
                          : 'border-slate-200 bg-slate-50/50 hover:bg-slate-50'
                      }`}
                    >
                      {/* Property Header */}
                      <div className="flex items-center justify-between mb-2">
                        <button
                          onClick={() => handlePropertySelect(group.property)}
                          className="flex items-center gap-1.5 font-bold text-xs text-slate-800 hover:text-orange text-left"
                        >
                          <Building2 size={14} className="text-orange shrink-0" />
                          <span>{group.property.name}</span>
                        </button>
                        <span className="text-[10px] px-2 py-0.5 rounded-full bg-white border text-slate-600 font-medium">
                          {group.buildings.length} buildings
                        </span>
                      </div>

                      {/* Buildings under this property */}
                      {group.buildings.length === 0 ? (
                        <p className="text-[11px] text-slate-400 italic py-1 pl-1">
                          No buildings created yet for this property.
                        </p>
                      ) : (
                        <div className="space-y-1.5 mt-2">
                          {group.buildings.map((b) => {
                            const isHotel =
                              modelOverrides[b.id] === 'hotel' ||
                              (!modelOverrides[b.id] &&
                                (b.name.toLowerCase().includes('hotel') ||
                                  b.name.toLowerCase().includes('resort')));
                            const isSelected = false;

                            return (
                              <div
                                key={b.id}
                                onClick={() => handleBuildingClick(b)}
                                className={`w-full flex items-center justify-between p-2.5 rounded-lg border text-left transition-all cursor-pointer ${
                                  isSelected
                                    ? 'bg-slate-900 text-white border-slate-900 shadow-sm'
                                    : 'bg-white hover:border-slate-300 text-slate-800'
                                }`}
                              >
                                <div className="flex items-center gap-2 min-w-0">
                                  <span className="text-lg">{isHotel ? '🏨' : '🏢'}</span>
                                  <div className="min-w-0">
                                    <p className="font-semibold text-xs truncate">{b.name}</p>
                                    <p className={`text-[10px] ${isSelected ? 'text-slate-300' : 'text-slate-400'}`}>
                                      {b.total_floors} floors · {isHotel ? 'Hotel Style' : 'Tower Style'}
                                    </p>
                                  </div>
                                </div>

                                {/* 3D Model Switcher Button */}
                                <button
                                  type="button"
                                  onClick={(e) => handleToggleBuildingModel(b, e)}
                                  className={`px-2 py-1 rounded text-[10px] font-medium border flex items-center gap-1 transition-all ${
                                    isSelected
                                      ? 'bg-slate-800 border-slate-700 text-orange hover:bg-slate-700'
                                      : 'bg-slate-50 border-slate-200 text-slate-600 hover:bg-orange/10 hover:text-orange'
                                  }`}
                                  title="Toggle 3D model style between Tower and Hotel"
                                >
                                  <Sparkles size={10} />
                                  <span>{isHotel ? 'Hotel' : 'Tower'}</span>
                                </button>
                              </div>
                            );
                          })}
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            ) : (
              <>
                {/* Single Building Floors & Rooms Inspector */}
                {selectedBuilding && (
                  <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-5 space-y-4">
                    <div className="flex items-start justify-between">
                      <div>
                        <div className="flex items-center gap-2">
                          <span className="text-xl">
                            {modelOverrides[selectedBuilding.id] === 'hotel' ||
                            selectedBuilding.name.toLowerCase().includes('hotel')
                              ? '🏨'
                              : '🏢'}
                          </span>
                          <div>
                            <h3 className="font-bold text-slate-900 text-base">{selectedBuilding.name}</h3>
                            <p className="text-xs text-slate-500">{selectedBuilding.total_floors} Floors Stack</p>
                          </div>
                        </div>
                      </div>

                      <button
                        type="button"
                        onClick={(e) => handleToggleBuildingModel(selectedBuilding, e)}
                        className="px-2.5 py-1 rounded-lg text-xs font-semibold border border-orange/30 bg-orange/10 text-orange hover:bg-orange/20 transition-colors flex items-center gap-1"
                      >
                        <Sparkles size={12} />
                        Switch Model
                      </button>
                    </div>

                    <div className="border-t pt-3">
                    <h4 className="font-semibold text-xs text-slate-700 mb-2">Select a Floor</h4>
                    {floors.length === 0 ? (
                      <p className="text-xs text-slate-400">No floors added to this building yet.</p>
                    ) : (
                      <div className="space-y-1.5 max-h-56 overflow-y-auto">
                        {[...floors]
                          .sort((a, b) => b.floor_number - a.floor_number)
                          .map((f) => (
                            <button
                              key={f.id}
                              onClick={() => handleFloorClick(f)}
                              className={`w-full flex items-center justify-between p-2.5 rounded-lg border text-left text-xs transition-colors ${
                                selectedFloor?.id === f.id
                                  ? 'border-orange bg-orange/10 font-bold text-orange'
                                  : 'hover:bg-slate-50 border-slate-200 text-slate-700'
                              }`}
                            >
                              <span className="flex items-center gap-1.5">
                                <Layers size={13} />
                                {f.name} (Lvl {f.floor_number})
                              </span>
                              <ChevronLeft
                                size={14}
                                className={`rotate-90 transition-transform ${
                                  selectedFloor?.id === f.id ? 'text-orange' : 'text-slate-400'
                                }`}
                              />
                            </button>
                          ))}
                      </div>
                    )}
                  </div>
                </div>
              )}

                {/* Rooms on selected floor */}
                {selectedFloor && (
                  <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-5 space-y-3">
                    <h4 className="font-bold text-slate-800 text-sm">{selectedFloor.name} Rooms</h4>
                    {rooms.filter((r) => r.floor_id === selectedFloor.id).length === 0 ? (
                      <p className="text-xs text-slate-400">No rooms mapped on this floor.</p>
                    ) : (
                      <div className="space-y-2 max-h-60 overflow-y-auto">
                        {rooms
                          .filter((r) => r.floor_id === selectedFloor.id)
                          .map((r) => (
                            <button
                              key={r.id}
                              onClick={() => handleRoomClick(r)}
                              className={`w-full flex items-center justify-between p-2.5 rounded-lg border text-left transition-all text-xs ${
                                selectedRoom?.id === r.id
                                  ? 'border-orange bg-orange/5 ring-1 ring-orange'
                                  : 'hover:bg-slate-50 border-slate-200'
                              }`}
                            >
                              <div>
                                <p className="font-semibold text-slate-900">{r.name}</p>
                                <p className="text-[10px] text-slate-400">
                                  Cap: {r.capacity} · Rp {r.price?.toLocaleString()}
                                </p>
                              </div>
                              <span
                                className={`px-2 py-0.5 rounded-full text-[10px] font-semibold ${
                                  r.is_available
                                    ? 'bg-green-50 text-green-700 border border-green-200'
                                    : 'bg-red-50 text-red-700 border border-red-200'
                                }`}
                              >
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
          )}
        </div>
      )}
    </div>
  );
}

