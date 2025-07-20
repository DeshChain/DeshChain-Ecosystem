import { StargateClient } from '@cosmjs/stargate'
import { Tendermint34Client } from '@cosmjs/tendermint-rpc'
import { 
  LoanApplication, 
  LoanStats, 
  FarmerProfile, 
  BusinessProfile, 
  StudentProfile,
  Loan,
  InterestRateQuote,
  EligibilityCheck,
  LendingModuleType
} from '../../types/lending'

/**
 * Client for interacting with DeshChain lending modules
 * Supports Krishi Mitra, Vyavasaya Mitra, and Shiksha Mitra
 */
export class LendingClient {
  constructor(
    private readonly client: StargateClient,
    private readonly tmClient: Tendermint34Client
  ) {}

  /**
   * Krishi Mitra (Agricultural Lending) Methods
   */

  /**
   * Get Krishi Mitra statistics
   */
  async getKrishiMitraStats(): Promise<LoanStats> {
    const response = await this.client.queryContractSmart(
      '', // Contract address would be here
      { get_krishi_stats: {} }
    )
    
    return {
      totalLoans: response.total_loans,
      totalDisbursed: response.total_disbursed,
      averageRate: response.average_rate,
      activeLoans: response.active_loans,
      defaultRate: response.default_rate,
    }
  }

  /**
   * Get farmer profile
   */
  async getFarmerProfile(farmerId: string): Promise<FarmerProfile | null> {
    try {
      const response = await this.client.queryContractSmart(
        '', // Contract address
        { get_farmer_profile: { farmer_id: farmerId } }
      )
      return response
    } catch {
      return null
    }
  }

  /**
   * Get farmer loans
   */
  async getFarmerLoans(farmerId: string): Promise<Loan[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_farmer_loans: { farmer_id: farmerId } }
    )
    return response.loans || []
  }

  /**
   * Get agricultural loan by ID
   */
  async getAgriculturalLoan(loanId: string): Promise<Loan | null> {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_loan: { loan_id: loanId } }
      )
      return response
    } catch {
      return null
    }
  }

  /**
   * Check farmer loan eligibility
   */
  async checkFarmerEligibility(application: LoanApplication): Promise<EligibilityCheck> {
    const response = await this.client.queryContractSmart(
      '',
      { check_farmer_eligibility: { application } }
    )
    return {
      eligible: response.eligible,
      reason: response.reason,
      maxAmount: response.max_amount,
      recommendedRate: response.recommended_rate,
    }
  }

  /**
   * Get agricultural interest rate quote
   */
  async getAgriculturalRateQuote(
    amount: number,
    duration: number,
    cropType: string,
    landSize: number
  ): Promise<InterestRateQuote> {
    const response = await this.client.queryContractSmart(
      '',
      { 
        get_rate_quote: { 
          amount, 
          duration, 
          crop_type: cropType,
          land_size: landSize,
          module: 'krishi'
        } 
      }
    )
    return {
      baseRate: response.base_rate,
      finalRate: response.final_rate,
      factors: response.factors,
      festivalBonus: response.festival_bonus,
    }
  }

  /**
   * Vyavasaya Mitra (Business Lending) Methods
   */

  /**
   * Get Vyavasaya Mitra statistics
   */
  async getVyavasayaMitraStats(): Promise<LoanStats> {
    const response = await this.client.queryContractSmart(
      '',
      { get_vyavasaya_stats: {} }
    )
    
    return {
      totalLoans: response.total_loans,
      totalDisbursed: response.total_disbursed,
      averageRate: response.average_rate,
      activeLoans: response.active_loans,
      defaultRate: response.default_rate,
    }
  }

  /**
   * Get business profile
   */
  async getBusinessProfile(businessId: string): Promise<BusinessProfile | null> {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_business_profile: { business_id: businessId } }
      )
      return response
    } catch {
      return null
    }
  }

  /**
   * Get business loans
   */
  async getBusinessLoans(businessId: string): Promise<Loan[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_business_loans: { business_id: businessId } }
    )
    return response.loans || []
  }

  /**
   * Check business loan eligibility
   */
  async checkBusinessEligibility(application: LoanApplication): Promise<EligibilityCheck> {
    const response = await this.client.queryContractSmart(
      '',
      { check_business_eligibility: { application } }
    )
    return {
      eligible: response.eligible,
      reason: response.reason,
      maxAmount: response.max_amount,
      recommendedRate: response.recommended_rate,
    }
  }

  /**
   * Get business interest rate quote
   */
  async getBusinessRateQuote(
    amount: number,
    duration: number,
    businessType: string,
    annualRevenue: number
  ): Promise<InterestRateQuote> {
    const response = await this.client.queryContractSmart(
      '',
      { 
        get_rate_quote: { 
          amount, 
          duration, 
          business_type: businessType,
          annual_revenue: annualRevenue,
          module: 'vyavasaya'
        } 
      }
    )
    return {
      baseRate: response.base_rate,
      finalRate: response.final_rate,
      factors: response.factors,
      festivalBonus: response.festival_bonus,
    }
  }

  /**
   * Shiksha Mitra (Education Loans) Methods
   */

  /**
   * Get Shiksha Mitra statistics
   */
  async getShikshaMitraStats(): Promise<LoanStats> {
    const response = await this.client.queryContractSmart(
      '',
      { get_shiksha_stats: {} }
    )
    
    return {
      totalLoans: response.total_loans,
      totalDisbursed: response.total_disbursed,
      averageRate: response.average_rate,
      activeLoans: response.active_loans,
      defaultRate: response.default_rate,
    }
  }

  /**
   * Get student profile
   */
  async getStudentProfile(studentId: string): Promise<StudentProfile | null> {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_student_profile: { student_id: studentId } }
      )
      return response
    } catch {
      return null
    }
  }

  /**
   * Get student loans
   */
  async getStudentLoans(studentId: string): Promise<Loan[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_student_loans: { student_id: studentId } }
    )
    return response.loans || []
  }

  /**
   * Check education loan eligibility
   */
  async checkEducationEligibility(application: LoanApplication): Promise<EligibilityCheck> {
    const response = await this.client.queryContractSmart(
      '',
      { check_education_eligibility: { application } }
    )
    return {
      eligible: response.eligible,
      reason: response.reason,
      maxAmount: response.max_amount,
      recommendedRate: response.recommended_rate,
    }
  }

  /**
   * Get education interest rate quote
   */
  async getEducationRateQuote(
    amount: number,
    duration: number,
    institutionType: string,
    academicScore: number
  ): Promise<InterestRateQuote> {
    const response = await this.client.queryContractSmart(
      '',
      { 
        get_rate_quote: { 
          amount, 
          duration, 
          institution_type: institutionType,
          academic_score: academicScore,
          module: 'shiksha'
        } 
      }
    )
    return {
      baseRate: response.base_rate,
      finalRate: response.final_rate,
      factors: response.factors,
      festivalBonus: response.festival_bonus,
    }
  }

  /**
   * Get scholarships
   */
  async getScholarships(studentId?: string): Promise<any[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_scholarships: { student_id: studentId } }
    )
    return response.scholarships || []
  }

  /**
   * General Methods
   */

  /**
   * Search loans across all modules
   */
  async searchLoans(query: string): Promise<Loan[]> {
    const [krishiLoans, vyavasayaLoans, shikshaMitraLoans] = await Promise.allSettled([
      this.searchModuleLoans('krishi', query),
      this.searchModuleLoans('vyavasaya', query),
      this.searchModuleLoans('shiksha', query),
    ])

    const allLoans = [
      ...(krishiLoans.status === 'fulfilled' ? krishiLoans.value : []),
      ...(vyavasayaLoans.status === 'fulfilled' ? vyavasayaLoans.value : []),
      ...(shikshaMitraLoans.status === 'fulfilled' ? shikshaMitraLoans.value : []),
    ]

    return allLoans
  }

  private async searchModuleLoans(module: LendingModuleType, query: string): Promise<Loan[]> {
    const response = await this.client.queryContractSmart(
      '',
      { search_loans: { module, query } }
    )
    return response.loans || []
  }

  /**
   * Get lending analytics
   */
  async getLendingAnalytics(timeframe: '24h' | '7d' | '30d' = '7d') {
    const response = await this.client.queryContractSmart(
      '',
      { get_lending_analytics: { timeframe } }
    )
    return response
  }

  /**
   * Get festival impact on lending
   */
  async getFestivalLendingImpact() {
    const response = await this.client.queryContractSmart(
      '',
      { get_festival_lending_impact: {} }
    )
    return response
  }

  /**
   * Get regional lending statistics
   */
  async getRegionalStats(state?: string, district?: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_regional_stats: { state, district } }
    )
    return response
  }

  /**
   * Get loan by ID (works across all modules)
   */
  async getLoan(loanId: string): Promise<Loan | null> {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_loan: { loan_id: loanId } }
      )
      return response
    } catch {
      return null
    }
  }

  /**
   * Get active festival bonuses for lending
   */
  async getActiveLendingBonuses() {
    const response = await this.client.queryContractSmart(
      '',
      { get_active_lending_bonuses: {} }
    )
    return response.bonuses || []
  }

  /**
   * Get lending module parameters
   */
  async getLendingParams(module: LendingModuleType) {
    const response = await this.client.queryContractSmart(
      '',
      { get_lending_params: { module } }
    )
    return response
  }
}