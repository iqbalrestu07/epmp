import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  createOrganizationSchema,
  type CreateOrganizationFormData,
} from "../schema";

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
    formState: { errors },
  } = useForm<CreateOrganizationFormData>({
    resolver: zodResolver(createOrganizationSchema),
    defaultValues: {
      is_active: true,
      ...defaultValues,
    },
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="name">Organization Name</Label>
        <Input
          id="name"
          {...register("name")}
          placeholder="e.g. My Property Management Co."
        />
        {errors.name && (
          <p className="text-sm text-red-500">{errors.name.message}</p>
        )}
      </div>
      <div className="space-y-2">
        <Label htmlFor="domain">Domain</Label>
        <Input
          id="domain"
          {...register("domain")}
          placeholder="e.g. mycompany.com"
        />
        {errors.domain && (
          <p className="text-sm text-red-500">{errors.domain.message}</p>
        )}
      </div>
      <div className="flex items-center space-x-2">
        <input
          id="is_active"
          type="checkbox"
          {...register("is_active")}
          className="h-4 w-4 rounded border-gray-300"
        />
        <Label htmlFor="is_active">Active</Label>
      </div>

      <Button type="submit" disabled={isSubmitting}>
        {isSubmitting ? "Saving..." : "Save"}
      </Button>
    </form>
  );
}
