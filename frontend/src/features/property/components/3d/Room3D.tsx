import { Html } from '@react-three/drei';
import type { Room } from '../../../room/types';

export interface Room3DProps {
  room: Room;
  position: [number, number, number];
  onClick: (r: Room) => void;
}

export function Room3D({ room, position, onClick }: Room3DProps) {
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
        document.body.style.cursor = 'pointer';
      }}
      onPointerOut={() => {
        document.body.style.cursor = 'default';
      }}
    >
      <boxGeometry args={[0.6, 0.5, 0.6]} />
      <meshStandardMaterial
        color={room.is_available ? '#22c55e' : '#ef4444'}
        transparent
        opacity={0.8}
      />
      <Html position={[0, 0.4, 0]} center distanceFactor={4}>
        <div className="px-1 py-0.5 bg-black/70 text-white text-[10px] rounded whitespace-nowrap pointer-events-none">
          {room.name}
        </div>
      </Html>
    </mesh>
  );
}
