import React from 'react';
import { AlertCircle, Clock, AlertTriangle, CheckCircle2, XCircle } from 'lucide-react';

export type AlertBadgeType = 'overdue' | 'due_soon' | 'expiring_soon' | 'expired' | 'active' | 'paid' | 'pending';

interface AlertBadgeProps {
  type?: AlertBadgeType;
  variant?: AlertBadgeType;
  label?: string;
  sublabel?: string;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export const AlertBadge: React.FC<AlertBadgeProps> = ({
  type,
  variant,
  label,
  sublabel,
  size = 'md',
  className = '',
}) => {
  const resolvedType = (variant || type || 'pending') as AlertBadgeType;
  const configs: Record<AlertBadgeType, { bg: string; text: string; border: string; icon: React.ReactNode; defaultLabel: string }> = {
    overdue: {
      bg: 'bg-red-50',
      text: 'text-red-700',
      border: 'border-red-200',
      icon: <AlertCircle className="w-3.5 h-3.5 text-red-600 animate-pulse" />,
      defaultLabel: 'Jatuh Tempo',
    },
    due_soon: {
      bg: 'bg-amber-50',
      text: 'text-amber-700',
      border: 'border-amber-200',
      icon: <Clock className="w-3.5 h-3.5 text-amber-600" />,
      defaultLabel: 'Segera Jatuh Tempo',
    },
    expiring_soon: {
      bg: 'bg-orange-50',
      text: 'text-orange-700',
      border: 'border-orange-200',
      icon: <AlertTriangle className="w-3.5 h-3.5 text-orange-600" />,
      defaultLabel: 'Akan Habis Kontrak',
    },
    expired: {
      bg: 'bg-rose-100',
      text: 'text-rose-800',
      border: 'border-rose-300',
      icon: <XCircle className="w-3.5 h-3.5 text-rose-700" />,
      defaultLabel: 'Kontrak Habis',
    },
    active: {
      bg: 'bg-emerald-50',
      text: 'text-emerald-700',
      border: 'border-emerald-200',
      icon: <CheckCircle2 className="w-3.5 h-3.5 text-emerald-600" />,
      defaultLabel: 'Aktif',
    },
    paid: {
      bg: 'bg-teal-50',
      text: 'text-teal-700',
      border: 'border-teal-200',
      icon: <CheckCircle2 className="w-3.5 h-3.5 text-teal-600" />,
      defaultLabel: 'Lunas',
    },
    pending: {
      bg: 'bg-slate-50',
      text: 'text-slate-700',
      border: 'border-slate-200',
      icon: <Clock className="w-3.5 h-3.5 text-slate-500" />,
      defaultLabel: 'Menunggu',
    },
  };

  const config = configs[resolvedType] || configs.pending;
  const displayLabel = label || config.defaultLabel;

  const sizeClasses = {
    sm: 'text-xs px-2 py-0.5 gap-1',
    md: 'text-xs font-semibold px-2.5 py-1 gap-1.5',
    lg: 'text-sm font-semibold px-3 py-1.5 gap-2',
  }[size];

  return (
    <span
      className={`inline-flex items-center rounded-full border shadow-2xs font-medium transition-colors ${config.bg} ${config.text} ${config.border} ${sizeClasses} ${className}`}
    >
      {config.icon}
      <span>{displayLabel}</span>
      {sublabel && (
        <span className="opacity-75 font-normal ml-0.5">({sublabel})</span>
      )}
    </span>
  );
};
