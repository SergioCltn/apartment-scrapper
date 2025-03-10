export interface Apartment {
  id: string;
  title?: string;
  propertyPrice?: string;
  pricePerSqm?: string;
  communityFees?: string;
  details?: string;
  location?: string;
  description?: string;
  createdAt: Date;
  updatedAt?: Date;
}
