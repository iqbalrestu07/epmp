import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  createOrganizationSchema,
  type CreateOrganizationFormData,
} from "../schema";
import { Globe2, ShieldCheck, Sparkles, Building } from "lucide-react";

interface OrganizationFormProps {
  onSubmit: (data: CreateOrganizationFormData) => void;
  defaultValues?: Partial<CreateOrganizationFormData>;
  isSubmitting?: boolean;
}

export function OrganizationForm({
  onSubmit,
  defaultValues,
  isSubmitting,
}: OrganizationFormProps) {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<CreateOrganizationFormData>({
    resolver: zodResolver(createOrganizationSchema),
    defaultValues: {
      is_active: true,
      ...defaultValues,
    },
  });

  const nameVal = watch("name");
  const domainVal = watch("domain");
  const isActiveVal = watch("is_active");

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      {/* Interactive Preview Card */}
      <div className="bg-gradient-to-br from-slate-900 via-slate-800 to-black text-white p-5 rounded-2xl shadow-lg border border-slate-700/50 relative overflow-hidden">
        <div className="absolute top-0 right-0 p-6 opacity-10 pointer-events-none">
          <Building size={120} />
        </div>
        <div className="flex items-center justify-between mb-3 relative z-10">
          <span className="text-[11px] font-bold uppercase tracking-wider text-orange flex items-center gap-1.5">
            <Sparkles size={13} />
            Live Organization Preview
          </span>
          <span className={`text-[11px] font-semibold px-2.5 py-0.5 rounded-full border ${
            isActiveVal
              ? "bg-green-500/20 text-green-400 border-green-500/30"
              : "bg-red-500/20 text-red-400 border-red-500/30"
          }`}>
            {isActiveVal ? "Active Status" : "Inactive"}
          </span>
        </div>
        <h4 className="text-xl font-bold text-white tracking-tight relative z-10">
          {nameVal || "Your Organization Name"}
        </h4>
        <p className="text-xs text-slate-400 mt-1 flex items-center gap-1.5 relative z-10">
          <Globe2 size={13} className="text-orange" />
          {domainVal ? `https://${domainVal}` : "domain.company.com"}
        </p>
      </div>

      {/* Organization Name Field */}
      <div className="space-y-2">
        <Label htmlFor="name" className="text-sm font-semibold text-slate-800 flex items-center gap-1.5">
          <Building size={15} className="text-orange" />
          Organization Name
        </Label>
        <Input
          id="name"
          {...register("name")}
          placeholder="e.g. Acme Property Group"
          className="rounded-xl border-slate-300 py-2.5 px-3.5 focus:ring-orange/20"
        />
        <p className="text-[11px] text-slate-500">The primary legal or operating entity name for your tenant spaces.</p>
        {errors.name && (
          <p className="text-xs text-red-600 mt-1 font-medium">{errors.name.message}</p>
        )}
      </div>

      {/* Domain / Subdomain Field */}
      <div className="space-y-2">
        <Label htmlFor="domain" className="text-sm font-semibold text-slate-800 flex items-center gap-1.5">
          <Globe2 size={15} className="text-orange" />
          Custom Domain / Slug
        </Label>
        <div className="relative">
          <Input
            id="domain"
            {...register("domain")}
            placeholder="acmeproperty.com"
            className="rounded-xl border-slate-300 py-2.5 px-3.5 pl-9 focus:ring-orange/20"
          />
          <span className="absolute left-3 top-3 text-slate-400">
            @
          </span>
        </div>
        <p className="text-[11px] text-slate-500">Domain or unique slug used for branding and member invitations.</p>
        {errors.domain && (
          <p className="text-xs text-red-600 mt-1 font-medium">{errors.domain.message}</p>
        )}
      </div>

      {/* Active Checkbox Card */}
      <div className="flex items-center justify-between p-4 bg-slate-50 rounded-xl border border-slate-200 hover:border-slate-300 transition-colors cursor-pointer">
        <div className="flex items-start gap-3">
          <div className="w-8 h-8 rounded-lg bg-green-100 flex items-center justify-center text-green-700 shrink-0 mt-0.5">
            <ShieldCheck size={18} />
          </div>
          <div>
            <label htmlFor="is_active" className="text-sm font-semibold text-slate-800 cursor-pointer">
              Set as Active Immediately
            </label>
            <p className="text-xs text-slate-500">
              Active organizations can immediately register properties, buildings, and invite team members.
            </p>
          </div>
        </div>
        <input
          id="is_active"
          type="checkbox"
          {...register("is_active")}
          className="h-5 w-5 rounded border-slate-300 text-orange focus:ring-orange cursor-pointer accent-orange"
        />
      </div>

      <Button
        type="submit"
        disabled={isSubmitting}
        className="w-full py-3 bg-orange hover:bg-orange/90 text-white font-semibold rounded-xl shadow-md shadow-orange/20 transition-all text-sm"
      >
        {isSubmitting ? "Creating Organization..." : "Create Organization"}
      </Button>
    </form>
  );
}
