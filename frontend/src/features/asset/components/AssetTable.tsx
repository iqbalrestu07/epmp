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
      case "Available": return "bg-green-100 text-green-700 border-green-200";
      case "Assigned": return "bg-blue-100 text-blue-700 border-blue-200";
      case "Maintenance": return "bg-yellow-100 text-yellow-700 border-yellow-200";
      case "Disposed": return "bg-red-100 text-red-700 border-red-200";
      default: return "bg-gray-100 text-gray-700 border-gray-200";
    }
  };

  const list: Asset[] = Array.isArray(data) ? data : (data as any)?.data ?? [];

  if (!list || list.length === 0) {
    return <div className="text-center py-12 text-black/50 text-sm">No assets found.</div>;
  }

  return (
    <div className="w-full overflow-auto">
      <table className="w-full text-left text-sm">
        <thead className="bg-black/[0.02] text-black/60 font-medium">
          <tr>
            <th className="px-4 py-3">Asset Name</th>
            <th className="px-4 py-3">Category</th>
            <th className="px-4 py-3">Property</th>
            <th className="px-4 py-3">Price</th>
            <th className="px-4 py-3">Status</th>
            <th className="px-4 py-3 text-right">Actions</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-black/5">
          {list.map((row: Asset) => (
            <tr 
              key={row.id} 
              onClick={() => onRowClick?.(row)}
              className="hover:bg-black/[0.02] cursor-pointer transition-colors"
            >
              <td className="px-4 py-3">
                <div className="flex items-center gap-2">
                  <Tag size={16} className="text-orange" />
                  <span className="font-medium text-black/80">{row.name}</span>
                </div>
              </td>
              <td className="px-4 py-3 text-black/60">{row.category}</td>
              <td className="px-4 py-3 text-black/60">{getPropertyName(row.property_id)}</td>
              <td className="px-4 py-3 text-black/60">
                {row.purchase_price > 0 ? `$${row.purchase_price.toLocaleString()}` : "-"}
              </td>
              <td className="px-4 py-3">
                <span className={`px-2 py-1 rounded-full text-xs border ${getStatusColor(row.status)}`}>
                  {row.status}
                </span>
              </td>
              <td className="px-4 py-3 text-right">
                <div className="flex items-center justify-end gap-2">
                  <button 
                    onClick={(e) => onEditClick?.(row, e)}
                    className="p-1.5 text-black/40 hover:text-orange hover:bg-orange/10 rounded-md transition-colors"
                  >
                    <Edit2 size={16} />
                  </button>
                  <button 
                    onClick={(e) => onDeleteClick?.(row, e)}
                    className="p-1.5 text-black/40 hover:text-red-500 hover:bg-red-500/10 rounded-md transition-colors"
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
