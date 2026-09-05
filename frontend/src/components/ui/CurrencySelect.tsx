import React, { forwardRef } from 'react';
import { SUPPORTED_CURRENCIES } from '@/utils/currency';

export interface CurrencySelectProps extends React.SelectHTMLAttributes<HTMLSelectElement> {
  value?: string;
  onValueChange?: (code: string) => void;
  className?: string;
}

export const CurrencySelect = forwardRef<HTMLSelectElement, CurrencySelectProps>(
  ({ value, defaultValue, onValueChange, onChange, className = '', ...props }, ref) => {
    const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
      if (onValueChange) {
        onValueChange(e.target.value);
      }
      if (onChange) {
        onChange(e);
      }
    };

    return (
      <select
        ref={ref}
        value={value}
        defaultValue={value !== undefined ? undefined : defaultValue ?? 'IDR'}
        onChange={handleChange}
        className={`rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm font-semibold focus:outline-none focus:ring-2 focus:ring-orange/30 text-slate-800 ${className}`}
        {...props}
      >
        {SUPPORTED_CURRENCIES.map((c) => (
          <option key={c.code} value={c.code}>
            {c.code} ({c.symbol}) - {c.country}
          </option>
        ))}
      </select>
    );
  }
);

CurrencySelect.displayName = 'CurrencySelect';
