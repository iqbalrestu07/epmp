import {
  useReactTable,
  getCoreRowModel,
  flexRender,
  createColumnHelper,
} from "@tanstack/react-table";
import { DoorOpen, Eye, Pencil, Trash2, Users } from "lucide-react";
import type { Room } from "../types";
import type { Floor } from "../../floor/types";
import type { Building } from "../../building/types";
import type { Property } from "../../property/types";
import { formatCurrency } from "@/utils/currency";

const columnHelper = createColumnHelper<Room>();

interface RoomTableProps {
  data: Room[];
  floors?: Floor[];
  buildings?: Building[];
  properties?: Property[];
  onRowClick?: (row: Room) => void;
  onEdit?: (row: Room) => void;
  onDelete?: (row: Room) => void;
}

export function RoomTable({ data, floors = [], buildings = [], properties = [], onRowClick, onEdit, onDelete }: RoomTableProps) {
  const floorNameById = new Map(floors.map((f) => [f.id, f.name]));
  const floorBuildingById = new Map(floors.map((f) => [f.id, f.building_id]));
  const buildingNameById = new Map(buildings.map((b) => [b.id, b.name]));
  const propertyNameById = new Map(properties.map((p) => [p.id, p.name]));

  const columns = [
    columnHelper.accessor("name", {
      header: "Name",
      cell: (info) => (
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded-lg bg-orange/10 flex items-center justify-center flex-shrink-0">
            <DoorOpen size={16} className="text-orange" />
          </div>
          <span className="font-medium">{info.getValue()}</span>
        </div>
      ),
    }),
    columnHelper.accessor("property_id", {
      header: "Property",
      cell: (info) => {
        const id = info.getValue() as string;
        if (!id) return <span className="text-slate-400">—</span>;
        return <span className="text-sm text-slate-600">{propertyNameById.get(id) || <span className="text-slate-400 font-mono text-xs">#{id.slice(0, 8)}</span>}</span>;
      },
    }),
    columnHelper.accessor("floor_id", {
      id: "building",
      header: "Building",
      cell: (info) => {
        const floorId = info.getValue() as string;
        if (!floorId) return <span className="text-slate-400">—</span>;
        const buildingId = floorBuildingById.get(floorId);
        if (!buildingId) return <span className="text-slate-400">—</span>;
        return <span className="text-sm text-slate-600">{buildingNameById.get(buildingId) || <span className="text-slate-400 font-mono text-xs">#{buildingId.slice(0, 8)}</span>}</span>;
      },
    }),
    columnHelper.accessor("floor_id", {
      id: "floor",
      header: "Floor",
      cell: (info) => {
        const id = info.getValue() as string;
        if (!id) return <span className="text-slate-400">—</span>;
        return <span className="text-sm text-slate-600">{floorNameById.get(id) || <span className="text-slate-400 font-mono text-xs">#{id.slice(0, 8)}</span>}</span>;
      },
    }),
    columnHelper.accessor("capacity", {
      header: "Capacity",
      cell: (info) => (
        <span className="inline-flex items-center gap-1 text-sm text-slate-600">
          <Users size={14} className="text-slate-400" />
          {info.getValue()}
        </span>
      ),
    }),
    columnHelper.accessor("price", {
      header: "Price",
      cell: (info) => {
        const row = info.row.original;
        return (
          <span className="text-sm font-semibold text-slate-800">
            {formatCurrency(info.getValue(), (row as any).currency || "IDR")}
          </span>
        );
      },
    }),
    columnHelper.accessor("status", {
      header: "Status Kamar",
      cell: (info) => {
        const row = info.row.original;
        const status = (info.getValue() || (row.is_available ? "Available" : "Occupied")) as string;

        if (status === "Occupied") {
          return (
            <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-red-50 text-red-700 border border-red-200">
              <span className="w-1.5 h-1.5 rounded-full bg-red-500 animate-pulse"></span>
              Terisi (Occupied)
            </span>
          );
        }
        if (status === "Reserved") {
          return (
            <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-amber-50 text-amber-700 border border-amber-200">
              <span className="w-1.5 h-1.5 rounded-full bg-amber-500"></span>
              Dipesan (Belum Masuk)
            </span>
          );
        }
        if (status === "Maintenance") {
          return (
            <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-slate-100 text-slate-700 border border-slate-200">
              <span className="w-1.5 h-1.5 rounded-full bg-slate-500"></span>
              Perbaikan
            </span>
          );
        }
        return (
          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200">
            <span className="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
            Tersedia
          </span>
        );
      },
    }),
    columnHelper.display({
      id: "actions",
      header: "Actions",
      cell: (info) => (
        <div className="flex items-center gap-1">
          <button
            onClick={(e) => { e.stopPropagation(); onRowClick?.(info.row.original); }}
            className="p-1.5 rounded-lg hover:bg-slate-100 text-slate-500 hover:text-slate-900 transition-colors"
            title="View"
          >
            <Eye size={16} />
          </button>
          <button
            onClick={(e) => { e.stopPropagation(); onEdit?.(info.row.original); }}
            className="p-1.5 rounded-lg hover:bg-slate-100 text-slate-500 hover:text-orange transition-colors"
            title="Edit"
          >
            <Pencil size={16} />
          </button>
          <button
            onClick={(e) => { e.stopPropagation(); onDelete?.(info.row.original); }}
            className="p-1.5 rounded-lg hover:bg-slate-100 text-slate-500 hover:text-red-500 transition-colors"
            title="Delete"
          >
            <Trash2 size={16} />
          </button>
        </div>
      ),
    }),
  ];

  const table = useReactTable({
    data: data || [],
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
                No data found.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}
