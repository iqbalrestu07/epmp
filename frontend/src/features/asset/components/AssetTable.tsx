import type { Asset } from "../types";
import { Edit2, Trash2, Tag } from "lucide-react";
import { usePropertys } from "../../property/hooks";

interface AssetTableProps {
  data: Asset[];
  onRowClick?: (row: Asset) => void;
  onEditClick?: (row: Asset, e: React.MouseEvent) => void;
  onDeleteClick?: (row: Asset, e: React.MouseEvent) => void;
}

export function AssetTable({ data, onRowClick, onEditClick, onDeleteClick }: AssetTableProps) {
  const { data: propData } = usePropertys({ per_page: 100 });
  const properties = Array.isArray(propData?.data)
    ? propData.data
    : Array.isArray(propData)
    ? propData
    : [];

  const getPropertyName = (id: string) => {
    return properties.find(p => p.id === id)?.name || id;
  };

  const getStatusColor = (status: string) => {
    switch(status) {
      case "Available": return "bg-green-50 text-green-700 border-green-200";
      case "Assigned": return "bg-blue-50 text-blue-700 border-blue-200";
      case "Maintenance": return "bg-yellow-50 text-yellow-700 border-yellow-200";
      case "Disposed": return "bg-red-50 text-red-700 border-red-200";
      default: return "bg-slate-100 text-slate-700 border-slate-200";
    }
  };

  const list: Asset[] = Array.isArray(data) ? data : (data as any)?.data ?? [];

  if (!list || list.length === 0) {
    return (
      <div className="rounded-xl border border-slate-200 bg-white overflow-hidden shadow-sm">
        <div className="py-12 text-center text-slate-400">
          <Tag size={36} className="mx-auto mb-2 opacity-30" />
          <p className="font-medium text-slate-600">No assets found.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="rounded-xl border border-slate-200 bg-white overflow-hidden shadow-sm">
      <table className="w-full text-sm text-left">
        <thead className="border-b border-slate-200 bg-slate-50/75 text-slate-600 font-semibold">
          <tr>
            <th className="px-5 py-3.5">Asset Name</th>
            <th className="px-5 py-3.5">Category</th>
            <th className="px-5 py-3.5">Property</th>
            <th className="px-5 py-3.5">Price</th>
            <th className="px-5 py-3.5">Status</th>
            <th className="px-5 py-3.5 text-right">Actions</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100">
          {list.map((row: Asset) => (
            <tr
              key={row.id}
              onClick={() => onRowClick?.(row)}
              className="hover:bg-orange/5 cursor-pointer transition-colors"
            >
              <td className="px-5 py-4">
                <div className="flex items-center gap-2">
                  <div className="w-8 h-8 rounded-lg bg-orange/10 flex items-center justify-center text-orange shrink-0">
                    <Tag size={16} />
                  </div>
                  <span className="font-semibold text-slate-800">{row.name}</span>
                </div>
              </td>
              <td className="px-5 py-4 text-slate-600">{row.category}</td>
              <td className="px-5 py-4 text-slate-600">{getPropertyName(row.property_id)}</td>
              <td className="px-5 py-4 text-slate-600">
                {row.purchase_price > 0 ? `$${row.purchase_price.toLocaleString()}` : "—"}
              </td>
              <td className="px-5 py-4">
                <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold border ${getStatusColor(row.status)}`}>
                  {row.status}
                </span>
              </td>
              <td className="px-5 py-4 text-right">
                <div className="flex items-center justify-end gap-2">
                  <button
                    onClick={(e) => onEditClick?.(row, e)}
                    className="p-1.5 text-slate-400 hover:text-orange hover:bg-orange/10 rounded-md transition-colors"
                  >
                    <Edit2 size={16} />
                  </button>
                  <button
                    onClick={(e) => onDeleteClick?.(row, e)}
                    className="p-1.5 text-slate-400 hover:text-red-500 hover:bg-red-500/10 rounded-md transition-colors"
                  >
                    <Trash2 size={16} />
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
