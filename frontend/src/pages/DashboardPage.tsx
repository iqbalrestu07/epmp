import { Building2, Users, DoorOpen, Banknote } from 'lucide-react';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { motion } from 'framer-motion';

const data = [
  { name: 'Jan', revenue: 40000, occupancy: 80 },
  { name: 'Feb', revenue: 45000, occupancy: 82 },
  { name: 'Mar', revenue: 48000, occupancy: 85 },
  { name: 'Apr', revenue: 52000, occupancy: 89 },
  { name: 'May', revenue: 58000, occupancy: 91 },
  { name: 'Jun', revenue: 64000, occupancy: 94 },
  { name: 'Jul', revenue: 70000, occupancy: 95 },
];

export default function DashboardPage() {
  const metrics = [
    { title: 'Total Properties', value: '12', icon: Building2, trend: '+2 this month', color: 'text-blue-500', bg: 'bg-blue-100' },
    { title: 'Total Rooms', value: '1,204', icon: DoorOpen, trend: '+40 this month', color: 'text-purple-500', bg: 'bg-purple-100' },
    { title: 'Active Tenants', value: '984', icon: Users, trend: '+12% vs last month', color: 'text-green-500', bg: 'bg-green-100' },
    { title: 'Monthly Revenue', value: '$124,500', icon: Banknote, trend: '+8% vs last month', color: 'text-orange', bg: 'bg-orange/20' },
  ];

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <h1 className="text-3xl font-display mb-2">Overview</h1>
      <p className="text-gray-500 mb-8">Welcome back to EPMP. Here's what's happening today.</p>

      {/* Metrics Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {metrics.map((m, i) => {
          const Icon = m.icon;
          return (
            <motion.div 
              key={i} 
              whileHover={{ y: -5, boxShadow: "0px 10px 30px rgba(0,0,0,0.1)" }}
              className="bg-white p-6 rounded-2xl border shadow-sm transition-all cursor-default"
            >
              <div className="flex justify-between items-start mb-4">
                <div className={`w-10 h-10 rounded-full ${m.bg} flex items-center justify-center ${m.color}`}>
                  <Icon size={20} />
                </div>
              </div>
              <h3 className="text-3xl font-bold mb-1">{m.value}</h3>
              <p className="text-sm font-medium text-gray-900 mb-1">{m.title}</p>
              <p className="text-xs text-gray-500">{m.trend}</p>
            </motion.div>
          )
        })}
      </div>

      {/* Interactive Chart */}
      <motion.div 
        initial={{ opacity: 0, y: 20 }} 
        animate={{ opacity: 1, y: 0 }} 
        transition={{ delay: 0.2 }}
        className="bg-white p-8 rounded-2xl border shadow-sm mb-8"
      >
        <div className="mb-6">
          <h2 className="text-xl font-bold text-gray-800">Revenue & Occupancy Trend</h2>
          <p className="text-sm text-gray-500">Interactive overview of the last 7 months</p>
        </div>
        <div className="h-80 w-full">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={data} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
              <defs>
                <linearGradient id="colorRevenue" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#ffaa00" stopOpacity={0.8}/>
                  <stop offset="95%" stopColor="#ffaa00" stopOpacity={0}/>
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#f0f0f0" />
              <XAxis dataKey="name" axisLine={false} tickLine={false} tick={{ fill: '#888' }} dy={10} />
              <YAxis axisLine={false} tickLine={false} tick={{ fill: '#888' }} />
              <Tooltip 
                contentStyle={{ borderRadius: '12px', border: 'none', boxShadow: '0 4px 20px rgba(0,0,0,0.1)' }}
              />
              <Area 
                type="monotone" 
                dataKey="revenue" 
                stroke="#ffaa00" 
                strokeWidth={3}
                fillOpacity={1} 
                fill="url(#colorRevenue)" 
                activeDot={{ r: 8, fill: '#ffaa00', stroke: '#fff', strokeWidth: 2 }}
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </motion.div>
    </div>
  );
}
