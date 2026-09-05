import {
  useReactTable,
  getCoreRowModel,
  flexRender,
  createColumnHelper,
} from "@tanstack/react-table";
import { Link } from "react-router-dom";
import { Bed as BedIcon, DoorOpen, Layers, Building2, User } from "lucide-react";
import type { Bed } from "../types";
import { useRooms } from "../../room/hooks";
import { useFloors } from "../../floor/hooks";
import { useBuildings } from "../../building/hooks";
import { usePropertys } from "../../property/hooks";
import { useOccupancys } from "../../occupancy/hooks";
import { useContracts } from "../../contract/hooks";
import { useTenants } from "../../tenant/hooks";

const columnHelper = createColumnHelper<Bed>();

interface BedTableProps {
  data: Bed[];
  onRowClick?: (row: Bed) => void;
}

export function BedTable({ data, onRowClick }: BedTableProps) {
  const { data: roomsData } = useRooms({ per_page: 100 });
  const { data: floorsData } = useFloors({ per_page: 100 });
  const { data: buildingsData } = useBuildings({ per_page: 100 });
  const { data: propertiesData } = usePropertys({ per_page: 100 });
  const { data: occupanciesData } = useOccupancys({ per_page: 100 });
  const { data: contractsData } = useContracts({ per_page: 100 });
  const { data: tenantsData } = useTenants({ per_page: 100 });

  const rooms = Array.isArray(roomsData?.data) ? roomsData.data : Array.isArray(roomsData) ? roomsData : [];
  const floors = Array.isArray(floorsData?.data) ? floorsData.data : Array.isArray(floorsData) ? floorsData : [];
  const buildings = Array.isArray(buildingsData?.data) ? buildingsData.data : Array.isArray(buildingsData) ? buildingsData : [];
  const properties = Array.isArray(propertiesData?.data) ? propertiesData.data : Array.isArray(propertiesData) ? propertiesData : [];
  const occupancies = Array.isArray(occupanciesData?.data) ? occupanciesData.data : Array.isArray(occupanciesData) ? occupanciesData : [];
  const contracts = Array.isArray(contractsData?.data) ? contractsData.data : Array.isArray(contractsData) ? contractsData : [];
  const tenants = Array.isArray(tenantsData?.data) ? tenantsData.data : Array.isArray(tenantsData) ? tenantsData : [];

  const roomsMap = new Map(rooms.map((r: any) => [r.id, r]));
  const floorsMap = new Map(floors.map((f: any) => [f.id, f]));
  const buildingsMap = new Map(buildings.map((b: any) => [b.id, b]));
  const propertiesMap = new Map(properties.map((p: any) => [p.id, p]));
  const tenantsMap = new Map(tenants.map((t: any) => [t.id, t]));

  const columns = [
    columnHelper.accessor("name", {
      header: "Tempat Tidur (Bed)",
      cell: (info) => (
        <div className="flex items-center gap-2.5">
          <div className="w-8 h-8 rounded-lg bg-orange/10 flex items-center justify-center text-orange shrink-0">
            <BedIcon size={16} />
          </div>
          <div>
            <p className="font-semibold text-slate-900">{info.getValue()}</p>
            <p className="text-xs text-slate-400 font-mono">ID: #{info.row.original.id.slice(0, 8)}</p>
          </div>
        </div>
      ),
    }),
    columnHelper.accessor("room_id", {
      header: "Kamar / Unit",
      cell: (info) => {
        const roomId = info.getValue() as string;
        const room = roomsMap.get(roomId);
        if (!room) return <span className="text-slate-400 font-mono text-xs">#{roomId?.slice(0, 8) || "—"}</span>;
        return (
          <Link
            to={`/dashboard/rooms/${roomId}`}
            onClick={(e) => e.stopPropagation()}
            className="group flex items-center gap-1.5 font-medium text-slate-800 hover:text-orange"
          >
            <DoorOpen size={15} className="text-slate-400 group-hover:text-orange transition-colors" />
            <span className="group-hover:underline">{room.name}</span>
          </Link>
        );
      },
    }),
    columnHelper.display({
      id: "hierarchy",
      header: "Gedung & Lantai",
      cell: ({ row }) => {
        const room = roomsMap.get(row.original.room_id);
        const floor = room?.floor_id ? floorsMap.get(room.floor_id) : null;
        const building = floor?.building_id ? buildingsMap.get(floor.building_id) : null;

        if (!building && !floor) {
          return <span className="text-slate-400 text-xs">—</span>;
        }

        return (
          <div className="flex flex-col gap-0.5">
            <div className="flex items-center gap-1 text-slate-800 font-medium text-xs">
              <Building2 size={13} className="text-slate-400" />
              <span>{building?.name || "Gedung"}</span>
            </div>
            <div className="flex items-center gap-1 text-slate-500 text-xs">
              <Layers size={13} className="text-slate-400" />
              <span>{floor ? `${floor.name} (Lt ${floor.floor_number})` : "Lantai"}</span>
            </div>
          </div>
        );
      },
    }),
    columnHelper.display({
      id: "property",
      header: "Properti",
      cell: ({ row }) => {
        const room = roomsMap.get(row.original.room_id);
        const prop = room?.property_id ? propertiesMap.get(room.property_id) : null;
        return (
          <span className="text-xs font-medium text-slate-700">
            {prop ? prop.name : "—"}
          </span>
        );
      },
    }),
    columnHelper.accessor("status", {
      header: "Status & Penghuni",
      cell: (info) => {
        const bed = info.row.original;
        const room = roomsMap.get(bed.room_id);

        // Cari occupancy aktif atau kontrak aktif di kamar ini
        const activeOccupancy = occupancies.find(
          (o: any) => o.room_id === bed.room_id && (o.status === "CheckedIn" || o.status === "Active")
        );
        const activeContract = contracts.find(
          (c: any) => c.room_id === bed.room_id && (c.status === "Active" || c.status === "Signed")
        );
        const tenantId = activeOccupancy?.tenant_id || activeContract?.tenant_id;
        const tenant = tenantId ? tenantsMap.get(tenantId) : null;

        // Tentukan operational status:
        // Jika bed berstatus Maintenance -> Maintenance
        // Jika kamar terisi atau bed status Occupied -> Occupied
        // Jika kamar kosong dan bed status Available -> Available
        let displayStatus = bed.status || "Available";
        if (displayStatus !== "Maintenance") {
          if (activeOccupancy || activeContract || (room && !room.is_available) || room?.status === "Occupied") {
            displayStatus = "Occupied";
          } else {
            displayStatus = "Available";
          }
        }

        const isOccupied = displayStatus === "Occupied";
        const isMaintenance = displayStatus === "Maintenance";

        return (
          <div className="flex flex-col gap-1">
            <span
              className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-semibold w-fit ${
                isMaintenance
                  ? "bg-amber-50 text-amber-700 border border-amber-200"
                  : isOccupied
                  ? "bg-rose-50 text-rose-700 border border-rose-200"
                  : "bg-emerald-50 text-emerald-700 border border-emerald-200"
              }`}
            >
              {isMaintenance ? "Maintenance" : isOccupied ? "Occupied" : "Available"}
            </span>

            {isOccupied && tenant && (
              <div className="flex items-center gap-1 text-xs text-slate-600 font-medium">
                <User size={12} className="text-slate-400" />
                <span className="truncate max-w-[140px]">{tenant.full_name}</span>
              </div>
            )}
          </div>
        );
      },
    }),
  ];

  const table = useReactTable({
    data: Array.isArray(data) ? data : [],
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="rounded-xl border border-slate-200 bg-white overflow-hidden shadow-sm">
      <table className="w-full text-sm text-left">
        <thead className="border-b border-slate-200 bg-slate-50/75 text-slate-600 font-semibold">
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <th key={header.id} className="px-5 py-3.5">
                  {flexRender(
                    header.column.columnDef.header,
                    header.getContext()
                  )}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody className="divide-y divide-slate-100">
          {table.getRowModel().rows.map((row) => (
            <tr
              key={row.id}
              className="hover:bg-orange/5 cursor-pointer transition-colors"
              onClick={() => onRowClick?.(row.original)}
            >
              {row.getVisibleCells().map((cell) => (
                <td key={cell.id} className="px-5 py-4">
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              ))}
            </tr>
          ))}
          {table.getRowModel().rows.length === 0 && (
            <tr>
              <td colSpan={columns.length} className="px-5 py-12 text-center text-slate-400">
                Belum ada data tempat tidur (Bed).
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}
