/**
 * Validation utilities for DeshChain SDK
 */

export class ValidationUtils {
  /**
   * Validate DeshChain address
   */
  static validateAddress(address: string): { valid: boolean; error?: string } {
    if (!address) {
      return { valid: false, error: 'Address is required' }
    }

    if (!address.startsWith('deshchain')) {
      return { valid: false, error: 'Address must start with "deshchain"' }
    }

    if (address.length !== 45) {
      return { valid: false, error: 'Address must be 45 characters long' }
    }

    const addressPart = address.slice(9) // Remove 'deshchain' prefix
    if (!/^[a-z0-9]+$/.test(addressPart)) {
      return { valid: false, error: 'Address contains invalid characters' }
    }

    return { valid: true }
  }

  /**
   * Validate amount
   */
  static validateAmount(amount: number | string): { valid: boolean; error?: string } {
    const num = typeof amount === 'string' ? parseFloat(amount) : amount

    if (isNaN(num)) {
      return { valid: false, error: 'Amount must be a valid number' }
    }

    if (num <= 0) {
      return { valid: false, error: 'Amount must be greater than 0' }
    }

    if (num > Number.MAX_SAFE_INTEGER) {
      return { valid: false, error: 'Amount is too large' }
    }

    return { valid: true }
  }

  /**
   * Validate denomination
   */
  static validateDenom(denom: string): { valid: boolean; error?: string } {
    if (!denom) {
      return { valid: false, error: 'Denomination is required' }
    }

    if (!/^[a-z][a-z0-9]*$/.test(denom)) {
      return { valid: false, error: 'Invalid denomination format' }
    }

    return { valid: true }
  }

  /**
   * Validate memo
   */
  static validateMemo(memo: string): { valid: boolean; error?: string } {
    if (memo.length > 512) {
      return { valid: false, error: 'Memo too long (max 512 characters)' }
    }

    return { valid: true }
  }

  /**
   * Validate gas limit
   */
  static validateGasLimit(gasLimit: number): { valid: boolean; error?: string } {
    if (!Number.isInteger(gasLimit)) {
      return { valid: false, error: 'Gas limit must be an integer' }
    }

    if (gasLimit <= 0) {
      return { valid: false, error: 'Gas limit must be positive' }
    }

    if (gasLimit > 10000000) {
      return { valid: false, error: 'Gas limit too high' }
    }

    return { valid: true }
  }

  /**
   * Validate transaction hash
   */
  static validateTxHash(hash: string): { valid: boolean; error?: string } {
    if (!hash) {
      return { valid: false, error: 'Transaction hash is required' }
    }

    if (!/^[A-F0-9]{64}$/i.test(hash)) {
      return { valid: false, error: 'Invalid transaction hash format' }
    }

    return { valid: true }
  }

  /**
   * Validate block height
   */
  static validateBlockHeight(height: number): { valid: boolean; error?: string } {
    if (!Number.isInteger(height)) {
      return { valid: false, error: 'Block height must be an integer' }
    }

    if (height < 0) {
      return { valid: false, error: 'Block height must be non-negative' }
    }

    return { valid: true }
  }

  /**
   * Validate URL
   */
  static validateUrl(url: string): { valid: boolean; error?: string } {
    try {
      new URL(url)
      return { valid: true }
    } catch {
      return { valid: false, error: 'Invalid URL format' }
    }
  }

  /**
   * Validate email
   */
  static validateEmail(email: string): { valid: boolean; error?: string } {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    
    if (!emailRegex.test(email)) {
      return { valid: false, error: 'Invalid email format' }
    }

    return { valid: true }
  }

  /**
   * Validate phone number (Indian format)
   */
  static validatePhoneNumber(phone: string): { valid: boolean; error?: string } {
    // Indian phone number validation
    const phoneRegex = /^[6-9]\d{9}$/
    
    if (!phoneRegex.test(phone)) {
      return { valid: false, error: 'Invalid phone number format' }
    }

    return { valid: true }
  }

  /**
   * Validate pincode (Indian format)
   */
  static validatePincode(pincode: string): { valid: boolean; error?: string } {
    const pincodeRegex = /^[1-9][0-9]{5}$/
    
    if (!pincodeRegex.test(pincode)) {
      return { valid: false, error: 'Invalid pincode format' }
    }

    return { valid: true }
  }

  /**
   * Validate Aadhar number
   */
  static validateAadhar(aadhar: string): { valid: boolean; error?: string } {
    // Remove spaces and hyphens
    const cleanAadhar = aadhar.replace(/[\s-]/g, '')
    
    if (!/^\d{12}$/.test(cleanAadhar)) {
      return { valid: false, error: 'Aadhar must be 12 digits' }
    }

    // Validate using Luhn-like algorithm for Aadhar
    const digits = cleanAadhar.split('').map(Number)
    const multiplicands = [2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1]
    
    let sum = 0
    for (let i = 0; i < 12; i++) {
      let product = digits[i] * multiplicands[i]
      sum += Math.floor(product / 10) + (product % 10)
    }

    if (sum % 10 !== 0) {
      return { valid: false, error: 'Invalid Aadhar number' }
    }

    return { valid: true }
  }

  /**
   * Validate PAN number
   */
  static validatePAN(pan: string): { valid: boolean; error?: string } {
    const panRegex = /^[A-Z]{5}[0-9]{4}[A-Z]{1}$/
    
    if (!panRegex.test(pan.toUpperCase())) {
      return { valid: false, error: 'Invalid PAN format' }
    }

    return { valid: true }
  }

  /**
   * Validate GST number
   */
  static validateGST(gst: string): { valid: boolean; error?: string } {
    const gstRegex = /^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z]{1}[1-9A-Z]{1}Z[0-9A-Z]{1}$/
    
    if (!gstRegex.test(gst.toUpperCase())) {
      return { valid: false, error: 'Invalid GST format' }
    }

    return { valid: true }
  }

  /**
   * Validate token symbol
   */
  static validateTokenSymbol(symbol: string): { valid: boolean; error?: string } {
    if (!symbol) {
      return { valid: false, error: 'Token symbol is required' }
    }

    if (symbol.length < 2 || symbol.length > 8) {
      return { valid: false, error: 'Token symbol must be 2-8 characters' }
    }

    if (!/^[A-Z][A-Z0-9]*$/.test(symbol)) {
      return { valid: false, error: 'Token symbol must start with letter and contain only uppercase letters and numbers' }
    }

    return { valid: true }
  }

  /**
   * Validate percentage (0-100)
   */
  static validatePercentage(percentage: number): { valid: boolean; error?: string } {
    if (isNaN(percentage)) {
      return { valid: false, error: 'Percentage must be a number' }
    }

    if (percentage < 0 || percentage > 100) {
      return { valid: false, error: 'Percentage must be between 0 and 100' }
    }

    return { valid: true }
  }

  /**
   * Validate date string (ISO format)
   */
  static validateDate(dateString: string): { valid: boolean; error?: string } {
    const date = new Date(dateString)
    
    if (isNaN(date.getTime())) {
      return { valid: false, error: 'Invalid date format' }
    }

    return { valid: true }
  }

  /**
   * Validate future date
   */
  static validateFutureDate(dateString: string): { valid: boolean; error?: string } {
    const dateValidation = ValidationUtils.validateDate(dateString)
    if (!dateValidation.valid) {
      return dateValidation
    }

    const date = new Date(dateString)
    const now = new Date()
    
    if (date <= now) {
      return { valid: false, error: 'Date must be in the future' }
    }

    return { valid: true }
  }

  /**
   * Validate object with multiple fields
   */
  static validateObject(
    obj: Record<string, any>,
    validators: Record<string, (value: any) => { valid: boolean; error?: string }>
  ): { valid: boolean; errors: Record<string, string> } {
    const errors: Record<string, string> = {}
    let valid = true

    for (const [field, validator] of Object.entries(validators)) {
      const result = validator(obj[field])
      if (!result.valid) {
        errors[field] = result.error || 'Invalid value'
        valid = false
      }
    }

    return { valid, errors }
  }
}