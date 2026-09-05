import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { ChevronRight } from "lucide-react";
import { useCreateContract } from "../hooks";
import { ContractForm } from "../components/ContractForm";
import { InteractiveContractPlanner } from "../components/InteractiveContractPlanner";
import type { CreateContractFormData } from "../schema";
import { Button } from "@/components/ui/button";

export function ContractCreatePage() {
  const navigate = useNavigate();
  const createMutation = useCreateContract();
  const [viewMode, setViewMode] = useState<"basic" | "interactive">("interactive");

  const handleSubmit = (data: CreateContractFormData | any) => {
    createMutation.mutate(data, {
      onSuccess: () => navigate("/dashboard/contracts"),
    });
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Link to="/dashboard" className="hover:text-slate-600">Dashboard</Link>
        <ChevronRight size={14} />
        <Link to="/dashboard/contracts" className="hover:text-slate-600">Contracts</Link>
        <ChevronRight size={14} />
        <span className="text-slate-600">New Contract</span>
      </div>

      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">New Contract</h1>
        <div className="flex gap-2 bg-slate-100 p-1 rounded-md">
          <Button
            variant={viewMode === "basic" ? "default" : "ghost"}
            size="sm"
            onClick={() => setViewMode("basic")}
            className="text-sm"
          >
            Basic Form
          </Button>
          <Button
            variant={viewMode === "interactive" ? "default" : "ghost"}
            size="sm"
            onClick={() => setViewMode("interactive")}
            className="text-sm"
          >
            Interactive Placement
          </Button>
        </div>
      </div>
      
      <div className="min-h-[400px]">
        {viewMode === "basic" ? (
          <div className="max-w-2xl bg-white p-8 rounded-xl border border-slate-200 shadow-sm">
            <ContractForm
              onSubmit={handleSubmit as any}
              isSubmitting={createMutation.isPending}
            />
          </div>
        ) : (
          <InteractiveContractPlanner
            onSubmit={handleSubmit}
            isSubmitting={createMutation.isPending}
          />
        )}
      </div>
    </div>
  );
}
