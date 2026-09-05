export interface BedTypeTemplate {
  id: string;
  name: string;
  dimensions: string;
  description: string;
  recommendedCapacity: number;
}

export const STANDARD_BED_TYPES: BedTypeTemplate[] = [
  {
    id: 'single_90',
    name: 'Single Bed (Standard)',
    dimensions: '90 x 200 cm',
    description: 'Cocok untuk kos kamar privat single / kamar minimalis',
    recommendedCapacity: 1,
  },
  {
    id: 'single_100',
    name: 'Single Bed (Super Single)',
    dimensions: '100 x 200 cm',
    description: 'Ranjang single lega untuk 1 orang',
    recommendedCapacity: 1,
  },
  {
    id: 'twin',
    name: 'Twin Bed (2x Single)',
    dimensions: '2 x (90 x 200 cm)',
    description: 'Dua ranjang terpisah dalam satu kamar',
    recommendedCapacity: 2,
  },
  {
    id: 'double',
    name: 'Double / Full Bed',
    dimensions: '140 x 200 cm',
    description: 'Ranjang nyaman untuk 1-2 orang',
    recommendedCapacity: 2,
  },
  {
    id: 'queen',
    name: 'Queen Bed',
    dimensions: '160 x 200 cm',
    description: 'Paling populer untuk kos eksklusif dan apartemen 1 BR',
    recommendedCapacity: 2,
  },
  {
    id: 'king',
    name: 'King Bed',
    dimensions: '180 x 200 cm',
    description: 'Sangat luas untuk tipe kamar Suite / Master Bedroom',
    recommendedCapacity: 2,
  },
  {
    id: 'super_king',
    name: 'Super King Bed',
    dimensions: '200 x 200 cm',
    description: 'Ukuran mewah maksimal untuk villa dan penthouse',
    recommendedCapacity: 2,
  },
  {
    id: 'bunk_bed',
    name: 'Bunk Bed (Tingkat)',
    dimensions: '90 x 200 cm (2 Tingkat)',
    description: 'Ranjang susun atas-bawah untuk efisiensi ruang',
    recommendedCapacity: 2,
  },
];
