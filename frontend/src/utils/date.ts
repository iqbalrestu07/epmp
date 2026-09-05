/**
 * Date formatting utilities for EPMP application.
 * Formats raw ISO timestamps (e.g. 2026-09-15T07:00:00+07:00) into human readable strings.
 */

export function formatDate(dateInput?: string | Date | null, options?: Intl.DateTimeFormatOptions): string {
  if (!dateInput) return '-';

  try {
    const date = typeof dateInput === 'string' ? new Date(dateInput) : dateInput;
    if (isNaN(date.getTime())) return typeof dateInput === 'string' ? dateInput : '-';

    const defaultOptions: Intl.DateTimeFormatOptions = {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      ...options,
    };

    return new Intl.DateTimeFormat('id-ID', defaultOptions).format(date);
  } catch {
    return typeof dateInput === 'string' ? dateInput : '-';
  }
}

export function formatDateTime(dateInput?: string | Date | null): string {
  if (!dateInput) return '-';

  try {
    const date = typeof dateInput === 'string' ? new Date(dateInput) : dateInput;
    if (isNaN(date.getTime())) return typeof dateInput === 'string' ? dateInput : '-';

    return new Intl.DateTimeFormat('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    }).format(date);
  } catch {
    return typeof dateInput === 'string' ? dateInput : '-';
  }
}

export function formatDayMonth(dateInput?: string | Date | null): string {
  if (!dateInput) return '-';

  try {
    const date = typeof dateInput === 'string' ? new Date(dateInput) : dateInput;
    if (isNaN(date.getTime())) return '-';

    return new Intl.DateTimeFormat('id-ID', {
      day: 'numeric',
      month: 'short',
    }).format(date);
  } catch {
    return '-';
  }
}
