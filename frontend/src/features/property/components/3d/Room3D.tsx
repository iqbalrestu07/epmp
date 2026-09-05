import { useState } from 'react';
import { Html } from '@react-three/drei';
import type { Room } from '../../../room/types';

export interface Room3DProps {
  room: Room;
  position: [number, number, number];
  selected?: boolean;
  onClick: (r: Room) => void;
}

export function Room3D({ room, position, selected = false, onClick }: Room3DProps) {
  const [hovered, setHovered] = useState(false);

  return (
    <mesh
      position={position}
      castShadow
      onClick={(e) => {
        e.stopPropagation();
        onClick(room);
      }}
      onPointerOver={(e) => {
        e.stopPropagation();
        setHovered(true);
        document.body.style.cursor = 'pointer';
      }}
      onPointerOut={() => {
        setHovered(false);
        document.body.style.cursor = 'default';
      }}
    >
      <boxGeometry args={[0.6, 0.5, 0.6]} />
      <meshStandardMaterial
        color={selected ? '#f97316' : hovered ? '#fbbf24' : room.is_available ? '#22c55e' : '#ef4444'}
        transparent
        opacity={selected ? 1 : 0.8}
        emissive={selected ? '#f97316' : '#000000'}
        emissiveIntensity={selected ? 0.3 : 0}
      />
      <Html position={[0, 0.4, 0]} center distanceFactor={4}>
        <div className={`px-1 py-0.5 text-[10px] rounded whitespace-nowrap pointer-events-none ${
          selected ? 'bg-orange-500 text-white' : 'bg-slate-900/70 text-white'
        }`}>
          {room.name}
        </div>
      </Html>
    </mesh>
  );
}
