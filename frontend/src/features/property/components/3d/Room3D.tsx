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

  const getBaseColor = () => {
    if (room.status === 'Occupied') return '#ef4444'; // Red - Sedang dihuni
    if (room.status === 'Reserved') return '#f59e0b'; // Amber - Sudah disewa, belum check-in
    if (room.status === 'Maintenance') return '#64748b'; // Slate - Perbaikan
    if (room.status === 'Available' || room.is_available) return '#22c55e'; // Green - Kosong
    return '#ef4444';
  };

  const getStatusBadge = () => {
    if (room.status === 'Occupied') return '🔴';
    if (room.status === 'Reserved') return '🟡';
    if (room.status === 'Maintenance') return '⚪';
    return '🟢';
  };

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
        color={selected ? '#f97316' : hovered ? '#fbbf24' : getBaseColor()}
        transparent
        opacity={selected ? 1 : 0.85}
        emissive={selected ? '#f97316' : '#000000'}
        emissiveIntensity={selected ? 0.3 : 0}
      />
      <Html position={[0, 0.4, 0]} center distanceFactor={4}>
        <div className={`px-1.5 py-0.5 text-[10px] font-semibold rounded whitespace-nowrap pointer-events-none flex items-center gap-1 shadow-sm ${
          selected ? 'bg-orange-500 text-white' : 'bg-slate-900/80 text-white'
        }`}>
          <span>{getStatusBadge()}</span>
          <span>{room.name}</span>
        </div>
      </Html>
    </mesh>
  );
}
