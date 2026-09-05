import { useState } from 'react';
import { Html } from '@react-three/drei';
import type { Floor } from '../../../floor/types';
import type { Room } from '../../../room/types';
import { Room3D } from './Room3D';

export interface Floor3DProps {
  floors: Floor[];
  rooms: Room[];
  selectedFloorId: string | null;
  onFloorClick: (f: Floor) => void;
  onRoomClick: (r: Room) => void;
}

export function Floor3D({ floors, rooms, selectedFloorId, onFloorClick, onRoomClick }: Floor3DProps) {
  const [hoveredFloor, setHoveredFloor] = useState<string | null>(null);

  const sortedFloors = [...floors].sort((a, b) => a.floor_number - b.floor_number);

  return (
    <group>
      {sortedFloors.map((floor, idx) => {
        const y = idx * 1.2;
        const floorRooms = rooms.filter((r) => r.floor_id === floor.id);
        const isSelected = selectedFloorId === floor.id;
        const isHovered = hoveredFloor === floor.id;

        return (
          <group key={floor.id} position={[0, y, 0]}>
            {/* Floor slab */}
            <mesh
              castShadow
              receiveShadow
              onClick={(e) => {
                e.stopPropagation();
                onFloorClick(floor);
              }}
              onPointerOver={(e) => {
                e.stopPropagation();
                setHoveredFloor(floor.id);
                document.body.style.cursor = 'pointer';
              }}
              onPointerOut={() => {
                setHoveredFloor(null);
                document.body.style.cursor = 'default';
              }}
            >
              <boxGeometry args={[3, 0.1, 3]} />
              <meshStandardMaterial
                color={isSelected ? '#f97316' : isHovered ? '#fbbf24' : '#cbd5e1'}
                transparent
                opacity={0.7}
              />
            </mesh>

            {/* Floor label */}
            <Html position={[1.7, 0, 0]} center distanceFactor={6}>
              <div className={`px-2 py-0.5 text-xs rounded whitespace-nowrap pointer-events-none ${
                isSelected ? 'bg-orange-500 text-white' : 'bg-slate-800 text-white'
              }`}>
                {floor.name}
              </div>
            </Html>

            {/* Rooms on this floor */}
            {isSelected && floorRooms.length > 0 && (
              <group position={[0, 0.15, 0]}>
                {floorRooms.map((room, rIdx) => {
                  const cols = Math.ceil(Math.sqrt(floorRooms.length));
                  const col = rIdx % cols;
                  const row = Math.floor(rIdx / cols);
                  const x = (col - (cols - 1) / 2) * 0.8;
                  const z = (row - (cols - 1) / 2) * 0.8;

                  return (
                    <Room3D
                      key={room.id}
                      room={room}
                      position={[x, 0.3, z]}
                      onClick={onRoomClick}
                    />
                  );
                })}
              </group>
            )}
          </group>
        );
      })}
    </group>
  );
}
