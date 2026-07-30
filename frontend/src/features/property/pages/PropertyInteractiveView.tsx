import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { ChevronLeft, DoorOpen, Users, Settings, Building2 } from 'lucide-react';
import { Link } from 'react-router-dom';

// MOCK DATA: Dinamis dari API nantinya
const buildingData = {
  name: "Tower Alpha - Student Housing",
  floors: [
    {
      id: 'f3', level: 3, name: 'Floor 3',
      rooms: [
        { id: '301', status: 'occupied', tenant: 'Budi' },
        { id: '302', status: 'available', tenant: null },
        { id: '303', status: 'maintenance', tenant: null },
        { id: '304', status: 'occupied', tenant: 'Andi' },
      ]
    },
    {
      id: 'f2', level: 2, name: 'Floor 2',
      rooms: [
        { id: '201', status: 'occupied', tenant: 'Citra' },
        { id: '202', status: 'occupied', tenant: 'Dewi' },
        { id: '203', status: 'available', tenant: null },
        { id: '204', status: 'available', tenant: null },
      ]
    },
    {
      id: 'f1', level: 1, name: 'Floor 1 (Lobby)',
      rooms: [
        { id: '101', status: 'occupied', tenant: 'Eka' },
        { id: '102', status: 'available', tenant: null },
      ]
    }
  ]
};

export default function PropertyInteractiveView() {
  const [selectedFloor, setSelectedFloor] = useState<any | null>(null);

  const getStatusColor = (status: string) => {
    switch(status) {
      case 'occupied': return 'bg-red-500 text-white shadow-red-500/40';
      case 'available': return 'bg-green-500 text-white shadow-green-500/40';
      case 'maintenance': return 'bg-yellow-500 text-white shadow-yellow-500/40';
      default: return 'bg-gray-200 text-gray-800';
    }
  };

  return (
    <div className="animate-in fade-in duration-500">
      <div className="flex items-center gap-4 mb-8">
        <Link to="/dashboard/properties" className="p-2 hover:bg-gray-200 rounded-full transition-colors">
          <ChevronLeft size={24} />
        </Link>
        <div>
          <h1 className="text-3xl font-display">{buildingData.name}</h1>
          <p className="text-gray-500">Interactive Building Map</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* LEFT: Building Stacking Plan */}
        <div className="col-span-1 bg-white p-8 rounded-3xl border shadow-sm flex flex-col items-center justify-end min-h-[600px] relative">
          <h2 className="absolute top-8 left-8 text-xl font-bold">Stacking Plan</h2>
          
          <div className="w-full max-w-[250px] flex flex-col gap-3">
            {buildingData.floors.map((floor) => {
              const isSelected = selectedFloor?.id === floor.id;
              
              // Calculate occupancy bar
              const total = floor.rooms.length;
              const occupied = floor.rooms.filter(r => r.status === 'occupied').length;
              const pct = (occupied / total) * 100;

              return (
                <motion.div
                  key={floor.id}
                  whileHover={{ scale: 1.05, x: 10 }}
                  onClick={() => setSelectedFloor(floor)}
                  className={`relative p-4 rounded-xl cursor-pointer border-2 transition-all duration-300 shadow-lg ${
                    isSelected ? 'border-orange bg-orange/5' : 'border-gray-200 bg-white hover:border-orange/50'
                  }`}
                  style={{
                    transformStyle: "preserve-3d",
                    transform: isSelected ? "perspective(1000px) rotateX(10deg)" : "perspective(1000px) rotateX(0deg)"
                  }}
                >
                  <div className="flex justify-between items-center mb-2">
                    <span className="font-bold">{floor.name}</span>
                    <span className="text-xs font-medium text-gray-500">{occupied}/{total} Rooms</span>
                  </div>
                  {/* Mini Progress Bar */}
                  <div className="w-full h-1.5 bg-gray-100 rounded-full overflow-hidden">
                    <div className="h-full bg-orange transition-all duration-1000" style={{ width: `${pct}%` }}></div>
                  </div>
                </motion.div>
              )
            })}
            
            {/* Ground / Foundation indicator */}
            <div className="w-[120%] h-4 bg-gray-300 rounded-full mx-auto mt-4 shadow-inner"></div>
          </div>
        </div>

        {/* RIGHT: Floor Details (Rooms Grid) */}
        <div className="col-span-1 lg:col-span-2">
          <AnimatePresence mode="wait">
            {!selectedFloor ? (
              <motion.div 
                key="empty"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                className="h-full bg-white rounded-3xl border shadow-sm flex flex-col items-center justify-center text-gray-400 p-10 text-center"
              >
                <Building2 size={64} className="mb-4 opacity-20" />
                <h3 className="text-xl font-medium mb-2">Select a Floor</h3>
                <p>Click on any floor from the stacking plan to view its detailed room layout.</p>
              </motion.div>
            ) : (
              <motion.div 
                key={selectedFloor.id}
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                transition={{ type: "spring", stiffness: 300, damping: 30 }}
                className="bg-white p-8 rounded-3xl border shadow-sm h-full"
              >
                <div className="flex justify-between items-center mb-8 pb-4 border-b">
                  <div>
                    <h2 className="text-2xl font-bold">{selectedFloor.name} Layout</h2>
                    <p className="text-gray-500">Live occupancy status</p>
                  </div>
                  
                  {/* Legends */}
                  <div className="flex gap-4 text-sm font-medium">
                    <div className="flex items-center gap-2"><span className="w-3 h-3 rounded-full bg-green-500"></span> Available</div>
                    <div className="flex items-center gap-2"><span className="w-3 h-3 rounded-full bg-red-500"></span> Occupied</div>
                    <div className="flex items-center gap-2"><span className="w-3 h-3 rounded-full bg-yellow-500"></span> Maintenance</div>
                  </div>
                </div>

                <div className="grid grid-cols-2 md:grid-cols-3 gap-6">
                  {selectedFloor.rooms.map((room: any) => (
                    <motion.div
                      key={room.id}
                      whileHover={{ scale: 1.05 }}
                      whileTap={{ scale: 0.95 }}
                      className={`p-6 rounded-2xl shadow-lg cursor-pointer ${getStatusColor(room.status)}`}
                    >
                      <div className="flex justify-between items-start mb-6">
                        <span className="text-2xl font-bold">Room {room.id}</span>
                        {room.status === 'occupied' && <Users size={20} className="opacity-80" />}
                        {room.status === 'available' && <DoorOpen size={20} className="opacity-80" />}
                        {room.status === 'maintenance' && <Settings size={20} className="opacity-80" />}
                      </div>
                      <div className="mt-4">
                        <p className="text-sm opacity-80 uppercase tracking-wider font-semibold">
                          {room.status}
                        </p>
                        {room.tenant && (
                          <p className="font-medium mt-1">{room.tenant}</p>
                        )}
                      </div>
                    </motion.div>
                  ))}
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </div>
    </div>
  );
}
