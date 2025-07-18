import 'dart:typed_data';
import 'package:flutter/material.dart';

import '../blockchain/deshchain_client.dart';
import '../wallet/hd_wallet.dart';
import '../../utils/constants.dart';
import '../../utils/logger.dart';

/// NAMO Token - Native token of DeshChain
class NAMOToken {
  static const String symbol = 'NAMO';
  static const String name = 'NAMO Token';
  static const String denom = 'namo';
  static const int decimals = 6;
  static const String description = 'Native token of DeshChain - India\'s Cultural Blockchain';
  static const String website = 'https://deshchain.org';
  static const String whitepaper = 'https://docs.deshchain.org/namo-tokenomics';
  
  // Cultural quotes related to NAMO
  static const List<String> culturalQuotes = [
    'नमो नमो भारत की शक्ति को',
    'NAMO - Nation's Asset, Mutual Ownership',
    'भारत की आत्मा, डिजिटल भविष्य',
    'Unity in Diversity, Strength in NAMO',
    'सत्यमेव जयते - Truth Prevails with NAMO',
  ];
  
  // Token icon (placeholder - replace with actual NAMO icon)
  static const IconData icon = Icons.account_balance;
  
  // Token colors matching Indian flag
  static const Color primaryColor = AppColors.saffron;
  static const Color secondaryColor = AppColors.green;
  static const Color accentColor = AppColors.white;
  
  /// Format NAMO amount with proper decimals
  static String formatAmount(BigInt amount) {
    final formatted = (amount.toDouble() / 1000000).toStringAsFixed(6);
    return formatted.replaceAll(RegExp(r'\.?0+$'), '');
  }
  
  /// Parse NAMO amount from string
  static BigInt parseAmount(String amount) {
    final double value = double.parse(amount);
    return BigInt.from(value * 1000000);
  }
  
  /// Get NAMO balance for address
  static Future<BigInt> getBalance(String address) async {
    try {
      final client = DeshChainClient();
      final balance = await client.getBalance(address);
      return balance.namoBalance;
    } catch (e) {
      AppLogger.error('Error getting NAMO balance: $e');
      return BigInt.zero;
    }
  }
  
  /// Send NAMO tokens
  static Future<String> sendTokens({
    required String fromAddress,
    required String toAddress,
    required BigInt amount,
    required String memo,
    required HDWallet wallet,
    required int accountIndex,
  }) async {
    try {
      final client = DeshChainClient();
      return await client.sendNAMO(
        fromAddress: fromAddress,
        toAddress: toAddress,
        amount: amount,
        memo: memo,
        wallet: wallet,
        accountIndex: accountIndex,
      );
    } catch (e) {
      AppLogger.error('Error sending NAMO tokens: $e');
      throw NAMOException('Failed to send NAMO tokens: $e');
    }
  }
  
  /// Get transaction history for NAMO
  static Future<List<NAMOTransaction>> getTransactionHistory(
    String address, {
    int limit = 50,
    int offset = 0,
  }) async {
    try {
      final client = DeshChainClient();
      final transactions = await client.getTransactionHistory(
        address,
        limit: limit,
        offset: offset,
      );
      
      return transactions
          .where((tx) => _isNAMOTransaction(tx))
          .map((tx) => NAMOTransaction.fromDeshChainTransaction(tx))
          .toList();
    } catch (e) {
      AppLogger.error('Error getting NAMO transaction history: $e');
      return [];
    }
  }
  
  /// Check if transaction involves NAMO
  static bool _isNAMOTransaction(DeshChainTransaction tx) {
    return tx.messages.any((msg) => 
      msg.type == '/cosmos.bank.v1beta1.MsgSend' &&
      msg.data['amount'] != null &&
      (msg.data['amount'] as List).any((amount) => 
        amount['denom'] == 'namo' || amount['denom'] == 'unamo'
      )
    );
  }
  
  /// Get current NAMO price in INR (placeholder)
  static Future<double> getCurrentPriceINR() async {
    // TODO: Implement actual price fetching from DeshChain API
    return 1.0; // Placeholder price
  }
  
  /// Get NAMO market stats
  static Future<NAMOMarketStats> getMarketStats() async {
    // TODO: Implement actual market stats fetching
    return NAMOMarketStats(
      priceINR: 1.0,
      priceUSD: 0.012,
      marketCap: 1428627666.0,
      volume24h: 50000000.0,
      change24h: 5.2,
      circulatingSupply: 1428627666.0,
      totalSupply: 1428627666.0,
    );
  }
  
  /// Get cultural quote for transaction
  static String getCulturalQuote() {
    final random = DateTime.now().millisecondsSinceEpoch % culturalQuotes.length;
    return culturalQuotes[random];
  }
  
  /// Validate NAMO address
  static bool isValidAddress(String address) {
    return address.startsWith('desh') && address.length == 45;
  }
  
  /// Get minimum transfer amount
  static BigInt getMinimumTransferAmount() {
    return BigInt.from(1000); // 0.001 NAMO
  }
  
  /// Get transaction fee
  static BigInt getTransactionFee() {
    return BigInt.from(5000); // 0.005 NAMO
  }
}

/// NAMO Transaction model
class NAMOTransaction {
  final String txHash;
  final String fromAddress;
  final String toAddress;
  final BigInt amount;
  final String memo;
  final DateTime timestamp;
  final bool isSuccess;
  final String? errorMessage;
  final TransactionType type;
  
  NAMOTransaction({
    required this.txHash,
    required this.fromAddress,
    required this.toAddress,
    required this.amount,
    required this.memo,
    required this.timestamp,
    required this.isSuccess,
    this.errorMessage,
    required this.type,
  });
  
  factory NAMOTransaction.fromDeshChainTransaction(DeshChainTransaction tx) {
    // Parse transaction data
    final sendMsg = tx.messages.firstWhere(
      (msg) => msg.type == '/cosmos.bank.v1beta1.MsgSend',
      orElse: () => throw NAMOException('No send message found'),
    );
    
    final fromAddress = sendMsg.data['from_address'] as String;
    final toAddress = sendMsg.data['to_address'] as String;
    final amounts = sendMsg.data['amount'] as List;
    final namoAmount = amounts.firstWhere(
      (amount) => amount['denom'] == 'namo' || amount['denom'] == 'unamo',
      orElse: () => {'amount': '0'},
    );
    
    return NAMOTransaction(
      txHash: tx.txHash,
      fromAddress: fromAddress,
      toAddress: toAddress,
      amount: BigInt.parse(namoAmount['amount']),
      memo: tx.memo ?? '',
      timestamp: DateTime.parse(tx.timestamp),
      isSuccess: tx.success,
      errorMessage: tx.success ? null : 'Transaction failed',
      type: TransactionType.transfer,
    );
  }
  
  /// Get formatted amount string
  String get formattedAmount => NAMOToken.formatAmount(amount);
  
  /// Get transaction direction
  TransactionDirection getDirection(String userAddress) {
    if (fromAddress == userAddress) {
      return TransactionDirection.sent;
    } else if (toAddress == userAddress) {
      return TransactionDirection.received;
    } else {
      return TransactionDirection.unknown;
    }
  }
  
  /// Get display address based on direction
  String getDisplayAddress(String userAddress) {
    final direction = getDirection(userAddress);
    switch (direction) {
      case TransactionDirection.sent:
        return toAddress;
      case TransactionDirection.received:
        return fromAddress;
      case TransactionDirection.unknown:
        return fromAddress;
    }
  }
}

/// NAMO Market Statistics
class NAMOMarketStats {
  final double priceINR;
  final double priceUSD;
  final double marketCap;
  final double volume24h;
  final double change24h;
  final double circulatingSupply;
  final double totalSupply;
  
  NAMOMarketStats({
    required this.priceINR,
    required this.priceUSD,
    required this.marketCap,
    required this.volume24h,
    required this.change24h,
    required this.circulatingSupply,
    required this.totalSupply,
  });
  
  /// Get formatted market cap
  String get formattedMarketCap {
    if (marketCap >= 1000000000) {
      return '₹${(marketCap / 1000000000).toStringAsFixed(2)}B';
    } else if (marketCap >= 1000000) {
      return '₹${(marketCap / 1000000).toStringAsFixed(2)}M';
    } else {
      return '₹${marketCap.toStringAsFixed(2)}';
    }
  }
  
  /// Get formatted volume
  String get formattedVolume {
    if (volume24h >= 1000000000) {
      return '₹${(volume24h / 1000000000).toStringAsFixed(2)}B';
    } else if (volume24h >= 1000000) {
      return '₹${(volume24h / 1000000).toStringAsFixed(2)}M';
    } else {
      return '₹${volume24h.toStringAsFixed(2)}';
    }
  }
  
  /// Get change color
  Color get changeColor {
    return change24h >= 0 ? Colors.green : Colors.red;
  }
  
  /// Get change icon
  IconData get changeIcon {
    return change24h >= 0 ? Icons.arrow_upward : Icons.arrow_downward;
  }
}

/// Transaction types
enum TransactionType {
  transfer,
  stake,
  unstake,
  reward,
  fee,
}

/// Transaction direction
enum TransactionDirection {
  sent,
  received,
  unknown,
}

/// NAMO Exception
class NAMOException implements Exception {
  final String message;
  
  NAMOException(this.message);
  
  @override
  String toString() => 'NAMOException: $message';
}