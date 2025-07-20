/**
 * Lending module types for DeshChain
 * Includes Krishi Mitra, Vyavasaya Mitra, and Shiksha Mitra
 */

export type LendingModuleType = 'krishi' | 'vyavasaya' | 'shiksha'

export interface LoanApplication {
  applicantId: string
  amount: number
  duration: number
  purpose: string
  collateral?: any
  documents: string[]
  metadata: Record<string, any>
}

export interface Loan {
  loanId: string
  applicantId: string
  amount: number
  interestRate: number
  duration: number
  status: LoanStatus
  disbursedAmount?: number
  repaidAmount?: number
  collateral?: any
  createdAt: string
  dueDate: string
  module: LendingModuleType
}

export type LoanStatus = 
  | 'pending'
  | 'approved' 
  | 'disbursed'
  | 'active'
  | 'completed'
  | 'defaulted'
  | 'rejected'

export interface LoanStats {
  totalLoans: number
  totalDisbursed: string
  averageRate: number
  activeLoans: number
  defaultRate: number
}

export interface FarmerProfile {
  farmerId: string
  name: string
  aadharNumber: string
  phoneNumber: string
  address: {
    village: string
    district: string
    state: string
    pincode: string
  }
  landDetails: {
    totalArea: number
    irrigatedArea: number
    landType: string
    soilType: string
  }
  cropHistory: CropInfo[]
  bankAccount: string
  kycStatus: string
  creditScore: number
  registrationDate: string
}

export interface CropInfo {
  season: string
  year: number
  cropType: string
  area: number
  expectedYield: number
  actualYield?: number
  marketPrice?: number
  expenses: number
  revenue?: number
}

export interface BusinessProfile {
  businessId: string
  businessName: string
  ownerName: string
  registrationNumber: string
  businessType: string
  industry: string
  address: {
    street: string
    city: string
    district: string
    state: string
    pincode: string
  }
  financials: {
    annualRevenue: number
    monthlyRevenue: number
    expenses: number
    assets: number
    liabilities: number
  }
  bankAccount: string
  gstNumber?: string
  panNumber: string
  kycStatus: string
  creditScore: number
  registrationDate: string
}

export interface StudentProfile {
  studentId: string
  name: string
  aadharNumber: string
  phoneNumber: string
  email: string
  address: {
    street: string
    city: string
    district: string
    state: string
    pincode: string
  }
  academicDetails: {
    currentInstitution: string
    course: string
    year: number
    marks: number
    percentage: number
    institutionType: 'government' | 'private'
  }
  familyDetails: {
    fatherName: string
    motherName: string
    fatherOccupation: string
    motherOccupation: string
    annualIncome: number
  }
  bankAccount: string
  kycStatus: string
  creditScore: number
  registrationDate: string
}

export interface InterestRateQuote {
  baseRate: number
  finalRate: number
  factors: {
    riskFactor: number
    culturalBonus: number
    festivalBonus: number
    regionBonus: number
    profileBonus: number
  }
  festivalBonus?: {
    festival: string
    bonus: number
    validUntil: string
  }
}

export interface EligibilityCheck {
  eligible: boolean
  reason?: string
  maxAmount: number
  recommendedRate: number
  requirements?: string[]
  documents?: string[]
}

export interface WeatherData {
  location: string
  temperature: number
  humidity: number
  rainfall: number
  windSpeed: number
  forecast: string
  lastUpdated: string
}

export interface MarketPrice {
  commodity: string
  market: string
  price: number
  unit: string
  date: string
  trend: 'up' | 'down' | 'stable'
}

export interface Scholarship {
  scholarshipId: string
  name: string
  provider: string
  amount: number
  eligibility: string[]
  deadline: string
  status: 'active' | 'expired' | 'closed'
  category: string
}

export interface LendingAnalytics {
  totalVolume: number
  averageTicketSize: number
  approvalRate: number
  disbursalTime: number
  repaymentRate: number
  regionalBreakdown: Record<string, number>
  monthlyTrends: Array<{
    month: string
    volume: number
    count: number
  }>
}

export interface FestivalOffer {
  festivalId: string
  festivalName: string
  discountPercent: number
  bonusAmount: number
  validFrom: string
  validTo: string
  eligibility: string[]
  terms: string[]
  maxBenefit: number
}

export interface RegionalStats {
  state: string
  district?: string
  totalLoans: number
  totalDisbursed: number
  averageTicketSize: number
  defaultRate: number
  popularCrops?: string[]
  popularBusinessTypes?: string[]
  popularCourses?: string[]
}

export interface LoanRepayment {
  loanId: string
  installmentNumber: number
  dueDate: string
  amount: number
  principal: number
  interest: number
  status: 'pending' | 'paid' | 'overdue'
  paidDate?: string
  lateFee?: number
}

export interface CreditHistory {
  applicantId: string
  totalLoans: number
  repaidLoans: number
  defaultedLoans: number
  totalAmountBorrowed: number
  totalAmountRepaid: number
  averageRepaymentTime: number
  creditScore: number
  lastUpdated: string
}