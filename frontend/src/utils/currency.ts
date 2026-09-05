export interface CurrencyInfo {
  code: string;
  name: string;
  symbol: string;
  country: string;
  decimals: number;
}

export const SUPPORTED_CURRENCIES: CurrencyInfo[] = [
  { code: 'IDR', name: 'Indonesian Rupiah', symbol: 'Rp', country: 'Indonesia', decimals: 0 },
  { code: 'USD', name: 'US Dollar', symbol: '$', country: 'United States / Global', decimals: 2 },
  { code: 'SGD', name: 'Singapore Dollar', symbol: 'S$', country: 'Singapore', decimals: 2 },
  { code: 'MYR', name: 'Malaysian Ringgit', symbol: 'RM', country: 'Malaysia', decimals: 2 },
  { code: 'THB', name: 'Thai Baht', symbol: '฿', country: 'Thailand', decimals: 2 },
  { code: 'PHP', name: 'Philippine Peso', symbol: '₱', country: 'Philippines', decimals: 2 },
  { code: 'VND', name: 'Vietnamese Dong', symbol: '₫', country: 'Vietnam', decimals: 0 },
  { code: 'EUR', name: 'Euro', symbol: '€', country: 'European Union', decimals: 2 },
  { code: 'GBP', name: 'British Pound', symbol: '£', country: 'United Kingdom', decimals: 2 },
  { code: 'AUD', name: 'Australian Dollar', symbol: 'A$', country: 'Australia', decimals: 2 },
  { code: 'JPY', name: 'Japanese Yen', symbol: '¥', country: 'Japan', decimals: 0 },
];

export function getCurrencyInfo(code?: string): CurrencyInfo {
  const normalized = (code || 'IDR').toUpperCase();
  return SUPPORTED_CURRENCIES.find((c) => c.code === normalized) || SUPPORTED_CURRENCIES[0];
}

export function formatCurrency(amount: number | string | null | undefined, currencyCode = 'IDR'): string {
  if (amount === null || amount === undefined || amount === '') return '—';
  const num = typeof amount === 'string' ? parseFloat(amount) : amount;
  if (isNaN(num)) return '—';

  const info = getCurrencyInfo(currencyCode);

  try {
    const formattedNum = new Intl.NumberFormat('id-ID', {
      minimumFractionDigits: info.decimals,
      maximumFractionDigits: info.decimals,
    }).format(num);

    return `${info.symbol} ${formattedNum}`;
  } catch {
    return `${info.symbol} ${num.toLocaleString()}`;
  }
}
