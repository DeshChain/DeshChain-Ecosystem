/**
 * Encoding utilities for DeshChain SDK
 */

import { toBase64, fromBase64, toHex, fromHex } from '@cosmjs/encoding'

export class EncodingUtils {
  /**
   * Convert string to base64
   */
  static toBase64(data: string): string {
    return toBase64(new TextEncoder().encode(data))
  }

  /**
   * Convert base64 to string
   */
  static fromBase64(data: string): string {
    return new TextDecoder().decode(fromBase64(data))
  }

  /**
   * Convert bytes to hex string
   */
  static toHex(data: Uint8Array): string {
    return toHex(data)
  }

  /**
   * Convert hex string to bytes
   */
  static fromHex(data: string): Uint8Array {
    return fromHex(data)
  }

  /**
   * Convert string to bytes
   */
  static stringToBytes(data: string): Uint8Array {
    return new TextEncoder().encode(data)
  }

  /**
   * Convert bytes to string
   */
  static bytesToString(data: Uint8Array): string {
    return new TextDecoder().decode(data)
  }

  /**
   * Encode JSON to base64
   */
  static jsonToBase64(data: any): string {
    return EncodingUtils.toBase64(JSON.stringify(data))
  }

  /**
   * Decode JSON from base64
   */
  static base64ToJson<T = any>(data: string): T {
    return JSON.parse(EncodingUtils.fromBase64(data))
  }

  /**
   * Create message hash
   */
  static hashMessage(message: string): string {
    // Simple hash implementation
    let hash = 0
    for (let i = 0; i < message.length; i++) {
      const char = message.charCodeAt(i)
      hash = ((hash << 5) - hash) + char
      hash = hash & hash // Convert to 32-bit integer
    }
    return Math.abs(hash).toString(16)
  }

  /**
   * Validate address format
   */
  static isValidAddress(address: string, prefix: string = 'deshchain'): boolean {
    if (!address.startsWith(prefix)) return false
    if (address.length !== prefix.length + 39) return false
    
    const addressPart = address.slice(prefix.length)
    return /^[a-z0-9]+$/.test(addressPart)
  }

  /**
   * Validate transaction hash
   */
  static isValidTxHash(hash: string): boolean {
    return /^[A-F0-9]{64}$/i.test(hash)
  }

  /**
   * Format amount with decimals
   */
  static formatAmount(amount: number | string, decimals: number = 6): string {
    const num = typeof amount === 'string' ? parseFloat(amount) : amount
    return (num / Math.pow(10, decimals)).toFixed(decimals)
  }

  /**
   * Parse amount to micro units
   */
  static parseAmount(amount: number | string, decimals: number = 6): string {
    const num = typeof amount === 'string' ? parseFloat(amount) : amount
    return Math.floor(num * Math.pow(10, decimals)).toString()
  }

  /**
   * Normalize denom
   */
  static normalizeDenom(denom: string): string {
    return denom.toLowerCase()
  }

  /**
   * Create random ID
   */
  static generateId(length: number = 16): string {
    const chars = 'abcdefghijklmnopqrstuvwxyz0123456789'
    let result = ''
    for (let i = 0; i < length; i++) {
      result += chars.charAt(Math.floor(Math.random() * chars.length))
    }
    return result
  }
}