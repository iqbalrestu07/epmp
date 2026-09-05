import { Suspense, useMemo } from 'react';
import { Canvas } from '@react-three/fiber';
import { OrbitControls, Html } from '@react-three/drei';
import * as THREE from 'three';
import type { Property } from '../types';
import type { Building } from '../../building/types';
import type { Floor } from '../../floor/types';
import type { Room } from '../../room/types';
import { Building3D, Floor3D } from './3d';

export interface PropertyGroup {
  property: Property;
  buildings: Building[];
}

export interface PropertyScene3DProps {
  buildings: Building[];
  floors: Floor[];
  rooms: Room[];
  selectedBuildingId: string | null;
  selectedFloorId: string | null;
  onBuildingClick: (b: Building) => void;
  onFloorClick: (f: Floor) => void;
  onRoomClick: (r: Room) => void;
  propertyGroups?: PropertyGroup[];
  isAllPropertiesMode?: boolean;
  onPropertyClick?: (p: Property) => void;
  buildingModelOverrides?: Record<string, 'building' | 'hotel'>;
}

// ─── Scene Lights & Environment ──────────────────────────────────────────────

function SceneLights() {
  return (
    <>
      <ambientLight intensity={0.65} />
      <directionalLight
        position={[15, 20, 10]}
        intensity={1.5}
        castShadow
        shadow-mapSize={[2048, 2048]}
        shadow-camera-far={60}
        shadow-camera-left={-25}
        shadow-camera-right={25}
        shadow-camera-top={25}
        shadow-camera-bottom={-25}
      />
      <directionalLight position={[-15, 10, -10]} intensity={0.4} />
    </>
  );
}

// ─── Ground Plane ────────────────────────────────────────────────────────────

function GroundPlane({ size = 50 }: { size?: number }) {
  return (
    <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, -0.1, 0]} receiveShadow>
      <planeGeometry args={[size, size]} />
      <meshStandardMaterial color="#e2e8f0" roughness={0.8} />
    </mesh>
  );
}

// ─── Property Zone Pad (in All Properties grouped view) ──────────────────────

function PropertyZonePad({
  property,
  position,
  width = 14,
  depth = 14,
  onClick,
}: {
  property: Property;
  position: [number, number, number];
  width?: number;
  depth?: number;
  onClick: (p: Property) => void;
}) {
  const edgesGeo = useMemo(() => {
    const plane = new THREE.PlaneGeometry(width, depth);
    return new THREE.EdgesGeometry(plane);
  }, [width, depth]);

  return (
    <group position={position}>
      {/* Zone Ground Plate */}
      <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, 0.01, 0]} receiveShadow>
        <planeGeometry args={[width, depth]} />
        <meshStandardMaterial color="#f1f5f9" roughness={0.6} />
      </mesh>

      {/* Zone Boundary Border */}
      <lineSegments geometry={edgesGeo} position={[0, 0.03, 0]} rotation={[-Math.PI / 2, 0, 0]}>
        <lineBasicMaterial color="#f97316" linewidth={2} />
      </lineSegments>

      {/* Property Floating Title */}
      <Html position={[0, 4.8, -depth / 2 + 1]} center distanceFactor={18}>
        <button
          onClick={(e) => {
            e.stopPropagation();
            onClick(property);
          }}
          className="px-3 py-1.5 rounded-xl bg-slate-900/90 hover:bg-orange text-white text-xs font-bold shadow-xl border border-white/20 transition-all flex items-center gap-1.5 cursor-pointer whitespace-nowrap"
        >
          <span>🏛️</span>
          <span>{property.name}</span>
          <span className="text-[10px] text-orange-200 font-normal capitalize">
            ({property.property_type ? property.property_type.replace('_', ' ') : 'Property'})
          </span>
        </button>
      </Html>
    </group>
  );
}

// ─── Main Scene Component ────────────────────────────────────────────────────

export function PropertyScene3D({
  buildings,
  floors,
  rooms,
  selectedBuildingId,
  selectedFloorId,
  onBuildingClick,
  onFloorClick,
  onRoomClick,
  propertyGroups = [],
  isAllPropertiesMode = false,
  onPropertyClick,
  buildingModelOverrides = {},
}: PropertyScene3DProps) {
  const showFloorStack = selectedBuildingId !== null;

  return (
    <Canvas
      shadows
      camera={{
        position: isAllPropertiesMode ? [22, 18, 22] : [12, 10, 12],
        fov: 42,
      }}
      className="rounded-2xl"
    >
      <Suspense fallback={<Html center><div className="text-slate-500 font-medium">Loading 3D scene & assets...</div></Html>}>
        <SceneLights />
        <GroundPlane size={isAllPropertiesMode ? 70 : 40} />
        <gridHelper
          args={[isAllPropertiesMode ? 70 : 40, isAllPropertiesMode ? 70 : 40, '#94a3b8', '#cbd5e1']}
          position={[0, 0.02, 0]}
        />

        {showFloorStack ? (
          <Floor3D
            floors={floors}
            rooms={rooms}
            selectedFloorId={selectedFloorId}
            onFloorClick={onFloorClick}
            onRoomClick={onRoomClick}
          />
        ) : isAllPropertiesMode && propertyGroups.length > 0 ? (
          // Grouped Multi-Property Spatial Clusters
          <group>
            {propertyGroups.map((group, pIdx) => {
              const count = propertyGroups.length;
              const spacing = 18;
              const cols = Math.ceil(Math.sqrt(count));
              const col = pIdx % cols;
              const row = Math.floor(pIdx / cols);
              const zoneX = (col - (cols - 1) / 2) * spacing;
              const zoneZ = (row - (cols - 1) / 2) * spacing;

              return (
                <group key={group.property.id}>
                  <PropertyZonePad
                    property={group.property}
                    position={[zoneX, 0, zoneZ]}
                    onClick={(p) => onPropertyClick?.(p)}
                  />

                  {/* Buildings inside this property's zone */}
                  {group.buildings.map((b, bIdx) => {
                    const bCols = Math.ceil(Math.sqrt(Math.max(group.buildings.length, 1)));
                    const bCol = bIdx % bCols;
                    const bRow = Math.floor(bIdx / bCols);
                    const bx = zoneX + (bCol - (bCols - 1) / 2) * 5;
                    const bz = zoneZ + (bRow - (bCols - 1) / 2) * 5;

                    return (
                      <Building3D
                        key={b.id}
                        building={b}
                        position={[bx, 0, bz]}
                        onClick={onBuildingClick}
                        selected={selectedBuildingId === b.id}
                        modelStyle={buildingModelOverrides[b.id] ?? 'auto'}
                      />
                    );
                  })}
                </group>
              );
            })}
          </group>
        ) : (
          // Single Selected Property's Buildings
          <group>
            {buildings.length > 0 ? (
              buildings.map((b, idx) => {
                const cols = Math.ceil(Math.sqrt(buildings.length));
                const col = idx % cols;
                const row = Math.floor(idx / cols);
                const x = (col - (cols - 1) / 2) * 6;
                const z = (row - (cols - 1) / 2) * 6;
                return (
                  <Building3D
                    key={b.id}
                    building={b}
                    position={[x, 0, z]}
                    onClick={onBuildingClick}
                    selected={selectedBuildingId === b.id}
                    modelStyle={buildingModelOverrides[b.id] ?? 'auto'}
                  />
                );
              })
            ) : (
              <Html center>
                <div className="text-slate-400 text-center bg-white/90 p-4 rounded-xl border shadow-sm">
                  <p className="text-sm font-semibold text-slate-700">No buildings in this property yet.</p>
                  <p className="text-xs text-slate-400 mt-1">Add buildings to view the spatial 3D model.</p>
                </div>
              </Html>
            )}
          </group>
        )}

        <OrbitControls
          enablePan
          enableZoom
          enableRotate
          minDistance={2}
          maxDistance={70}
          maxPolarAngle={Math.PI / 2.05}
          dampingFactor={0.06}
        />
      </Suspense>
    </Canvas>
  );
}
