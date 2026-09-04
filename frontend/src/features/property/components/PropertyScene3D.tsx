import { Suspense } from 'react';
import { Canvas } from '@react-three/fiber';
import { OrbitControls, Html } from '@react-three/drei';
import type { Building } from '../../building/types';
import type { Floor } from '../../floor/types';
import type { Room } from '../../room/types';
import { Building3D, Floor3D } from './3d';

// ─── Scene Lights & Environment ──────────────────────────────────────────────

function SceneLights() {
  return (
    <>
      <ambientLight intensity={0.5} />
      <directionalLight
        position={[10, 10, 5]}
        intensity={1}
        castShadow
        shadow-mapSize={[2048, 2048]}
        shadow-camera-far={50}
        shadow-camera-left={-10}
        shadow-camera-right={10}
        shadow-camera-top={10}
        shadow-camera-bottom={-10}
      />
      <directionalLight position={[-10, 5, -5]} intensity={0.3} />
    </>
  );
}

// ─── Ground Plane ────────────────────────────────────────────────────────────

function GroundPlane() {
  return (
    <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, -0.1, 0]} receiveShadow>
      <planeGeometry args={[30, 30]} />
      <meshStandardMaterial color="#e2e8f0" />
    </mesh>
  );
}

// ─── Main Scene Component ────────────────────────────────────────────────────

export interface PropertyScene3DProps {
  buildings: Building[];
  floors: Floor[];
  rooms: Room[];
  selectedBuildingId: string | null;
  selectedFloorId: string | null;
  onBuildingClick: (b: Building) => void;
  onFloorClick: (f: Floor) => void;
  onRoomClick: (r: Room) => void;
}

export function PropertyScene3D({
  buildings,
  floors,
  rooms,
  selectedBuildingId,
  selectedFloorId,
  onBuildingClick,
  onFloorClick,
  onRoomClick,
}: PropertyScene3DProps) {
  const showFloorStack = selectedBuildingId !== null;

  return (
    <Canvas
      shadows
      camera={{ position: [8, 6, 8], fov: 50 }}
      className="rounded-2xl"
    >
      <Suspense fallback={<Html center><div className="text-gray-500">Loading 3D model...</div></Html>}>
        <SceneLights />
        <GroundPlane />

        {!showFloorStack ? (
          <group>
            {buildings.length > 0 ? (
              buildings.map((b, idx) => {
                const cols = Math.ceil(Math.sqrt(buildings.length));
                const col = idx % cols;
                const row = Math.floor(idx / cols);
                const x = (col - (cols - 1) / 2) * 5;
                const z = (row - (cols - 1) / 2) * 5;
                return (
                  <Building3D
                    key={b.id}
                    building={b}
                    position={[x, 0, z]}
                    onClick={onBuildingClick}
                    selected={selectedBuildingId === b.id}
                  />
                );
              })
            ) : (
              <Html center>
                <div className="text-gray-400 text-center">
                  <p className="text-sm">No buildings yet.</p>
                  <p className="text-xs mt-1">Add buildings to see the 3D model.</p>
                </div>
              </Html>
            )}
          </group>
        ) : (
          <Floor3D
            floors={floors}
            rooms={rooms}
            selectedFloorId={selectedFloorId}
            onFloorClick={onFloorClick}
            onRoomClick={onRoomClick}
          />
        )}

        <OrbitControls
          enablePan
          enableZoom
          enableRotate
          minDistance={3}
          maxDistance={25}
          maxPolarAngle={Math.PI / 2.1}
        />
      </Suspense>
    </Canvas>
  );
}
