import 'dart:typed_data';
import 'package:flutter/material.dart';

import '../blockchain/deshchain_client.dart';
import '../wallet/hd_wallet.dart';
import '../../utils/constants.dart';
import '../../utils/logger.dart';

/// DINR Token - Algorithmic INR Stablecoin on DeshChain
class DINRToken {
  static const String symbol = 'DINR';
  static const String name = 'Digital INR';
  static const String denom = 'dinr';
  static const int decimals = 6;
  static const String description = 'Algorithmic stablecoin pegged 1:1 to Indian Rupee on DeshChain';
  static const String website = 'https://deshchain.org/dinr';
  static const String whitepaper = 'https://docs.deshchain.org/dinr-stablecoin';
  
  // Cultural quotes related to DINR stability
  static const List<String> culturalQuotes = [
    'स्थिरता में शक्ति - Stability in Strength',
    'DINR - Digital India's New Rupee',
    'भारत की डिजिटल मुद्रा, वैश्विक स्तर पर',
    'Stability meets Innovation with DINR',
    'एक रुपया डिजिटल, अनंत संभावनाएं',
  ];
  
  // Token icon
  static const IconData icon = Icons.currency_rupee;
  
  // DINR specific colors (Navy Blue representing stability)
  static const Color primaryColor = AppColors.navy;
  static const Color secondaryColor = AppColors.saffron;
  static const Color accentColor = AppColors.white;
  
  /// Format DINR amount with proper decimals (always shows 2 decimal places for currency)
  static String formatAmount(BigInt amount) {
    final value = amount.toDouble() / 1000000;
    return value.toStringAsFixed(2);
  }
  
  /// Parse DINR amount from string
  static BigInt parseAmount(String amount) {
    final double value = double.parse(amount);
    return BigInt.from(value * 1000000);
  }
  
  /// Get DINR balance for address
  static Future<BigInt> getBalance(String address) async {
    try {
      final client = DeshChainClient();
      final balance = await client.getDINRBalance(address);
      return balance;
    } catch (e) {
      AppLogger.error('Error getting DINR balance: $e');
      return BigInt.zero;
    }
  }
  
  /// Get collateral positions for address
  static Future<List<CollateralPosition>> getCollateralPositions(String address) async {
    try {
      final client = DeshChainClient();
      return await client.getCollateralPositions(address);
    } catch (e) {
      AppLogger.error('Error getting collateral positions: $e');
      return [];
    }
  }
  
  /// Mint DINR with collateral
  static Future<String> mintDINR({
    required String userAddress,
    required String collateralDenom,
    required BigInt collateralAmount,
    required BigInt dinrToMint,
    required HDWallet wallet,
    required int accountIndex,
  }) async {
    try {
      final client = DeshChainClient();
      return await client.mintDINR(
        userAddress: userAddress,
        collateralDenom: collateralDenom,
        collateralAmount: collateralAmount,
        dinrToMint: dinrToMint,
        wallet: wallet,
        accountIndex: accountIndex,
      );
    } catch (e) {
      AppLogger.error('Error minting DINR: $e');
      throw DINRException('Failed to mint DINR: $e');
    }
  }
  
  /// Burn DINR to retrieve collateral
  static Future<String> burnDINR({
    required String userAddress,
    required BigInt dinrToBurn,
    required HDWallet wallet,
    required int accountIndex,
  }) async {
    try {
      final client = DeshChainClient();
      return await client.burnDINR(
        userAddress: userAddress,
        dinrToBurn: dinrToBurn,
        wallet: wallet,
        accountIndex: accountIndex,
      );
    } catch (e) {
      AppLogger.error('Error burning DINR: $e');
      throw DINRException('Failed to burn DINR: $e');
    }
  }
  
  /// Send DINR tokens
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
      return await client.sendDINR(
        fromAddress: fromAddress,
        toAddress: toAddress,
        amount: amount,
        memo: memo,
        wallet: wallet,
        accountIndex: accountIndex,
      );
    } catch (e) {
      AppLogger.error('Error sending DINR tokens: $e');
      throw DINRException('Failed to send DINR tokens: $e');
    }
  }
  
  /// Get transaction history for DINR
  static Future<List<DINRTransaction>> getTransactionHistory(
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
          .where((tx) => _isDINRTransaction(tx))
          .map((tx) => DINRTransaction.fromDeshChainTransaction(tx))
          .toList();
    } catch (e) {
      AppLogger.error('Error getting DINR transaction history: $e');
      return [];
    }
  }
  
  /// Check if transaction involves DINR
  static bool _isDINRTransaction(DeshChainTransaction tx) {
    return tx.messages.any((msg) => 
      (msg.type == '/cosmos.bank.v1beta1.MsgSend' &&
       msg.data['amount'] != null &&
       (msg.data['amount'] as List).any((amount) => 
         amount['denom'] == 'dinr' || amount['denom'] == 'udinr'
       )) ||
      msg.type == '/deshchain.dinr.v1.MsgMintDINR' ||
      msg.type == '/deshchain.dinr.v1.MsgBurnDINR' ||
      msg.type == '/deshchain.dinr.v1.MsgDepositCollateral' ||
      msg.type == '/deshchain.dinr.v1.MsgWithdrawCollateral'
    );
  }
  
  /// Get current DINR price (should be 1.00 INR)
  static Future<double> getCurrentPriceINR() async {
    try {
      final client = DeshChainClient();
      final priceData = await client.getDINRPriceData();
      return priceData.priceINR;
    } catch (e) {
      AppLogger.error('Error getting DINR price: $e');
      return 1.0; // Default to 1 INR
    }
  }
  
  /// Get DINR stability metrics
  static Future<DINRStabilityMetrics> getStabilityMetrics() async {
    try {
      final client = DeshChainClient();
      return await client.getDINRStabilityMetrics();
    } catch (e) {
      AppLogger.error('Error getting DINR stability metrics: $e');
      return DINRStabilityMetrics.defaultMetrics();
    }
  }
  
  /// Get supported collateral assets
  static Future<List<CollateralAsset>> getSupportedCollaterals() async {
    try {
      final client = DeshChainClient();
      return await client.getSupportedCollaterals();
    } catch (e) {
      AppLogger.error('Error getting supported collaterals: $e');
      return _getDefaultCollaterals();
    }
  }
  
  /// Default collateral assets
  static List<CollateralAsset> _getDefaultCollaterals() {
    return [
      CollateralAsset(
        denom: 'btc',
        name: 'Bitcoin',
        symbol: 'BTC',
        decimals: 8,
        collateralRatio: 150,
        liquidationRatio: 120,
        minCollateralAmount: BigInt.from(1000), // 0.00001 BTC
        color: const Color(0xFFF7931A),
      ),
      CollateralAsset(
        denom: 'eth',
        name: 'Ethereum',
        symbol: 'ETH',
        decimals: 18,
        collateralRatio: 150,
        liquidationRatio: 120,
        minCollateralAmount: BigInt.from(1000000000000000), // 0.001 ETH
        color: const Color(0xFF627EEA),
      ),
      CollateralAsset(
        denom: 'usdt',
        name: 'Tether USD',
        symbol: 'USDT',
        decimals: 6,
        collateralRatio: 110,
        liquidationRatio: 105,
        minCollateralAmount: BigInt.from(1000000), // 1 USDT
        color: const Color(0xFF26A17B),
      ),
      CollateralAsset(
        denom: 'usdc',
        name: 'USD Coin',
        symbol: 'USDC',
        decimals: 6,
        collateralRatio: 110,
        liquidationRatio: 105,
        minCollateralAmount: BigInt.from(1000000), // 1 USDC
        color: const Color(0xFF2775CA),
      ),
    ];
  }
  
  /// Calculate maximum DINR that can be minted with given collateral
  static BigInt calculateMaxMintable(CollateralAsset asset, BigInt collateralAmount, double collateralPriceINR) {
    final collateralValueINR = collateralAmount.toDouble() / 
        (BigInt.from(10).pow(asset.decimals).toDouble()) * collateralPriceINR;
    final maxDINRValue = collateralValueINR / (asset.collateralRatio / 100);
    return BigInt.from(maxDINRValue * 1000000); // Convert to micro-DINR
  }
  
  /// Calculate required collateral for minting DINR
  static BigInt calculateRequiredCollateral(CollateralAsset asset, BigInt dinrAmount, double collateralPriceINR) {
    final dinrValueINR = dinrAmount.toDouble() / 1000000; // Convert from micro-DINR
    final requiredCollateralValueINR = dinrValueINR * (asset.collateralRatio / 100);
    final requiredCollateralAmount = requiredCollateralValueINR / collateralPriceINR;
    return BigInt.from(requiredCollateralAmount * BigInt.from(10).pow(asset.decimals).toDouble());
  }
  
  /// Get cultural quote for DINR operations
  static String getCulturalQuote() {
    final random = DateTime.now().millisecondsSinceEpoch % culturalQuotes.length;
    return culturalQuotes[random];
  }
  
  /// Validate DINR address (same as DeshChain addresses)
  static bool isValidAddress(String address) {
    return address.startsWith('desh') && address.length == 45;
  }
  
  /// Get minimum transfer amount
  static BigInt getMinimumTransferAmount() {
    return BigInt.from(1000000); // 1.0 DINR
  }
  
  /// Get transaction fee (0.1% capped at ₹100)
  static BigInt getTransactionFee(BigInt amount) {
    final feeAmount = (amount.toDouble() * 0.001).round(); // 0.1%
    final maxFee = BigInt.from(100000000); // ₹100 in micro-DINR
    return BigInt.from(feeAmount) > maxFee ? maxFee : BigInt.from(feeAmount);
  }
}

/// DINR Transaction model
class DINRTransaction {
  final String txHash;
  final String fromAddress;
  final String toAddress;
  final BigInt amount;
  final String memo;
  final DateTime timestamp;
  final bool isSuccess;
  final String? errorMessage;
  final DINRTransactionType type;
  final Map<String, dynamic>? additionalData;
  
  DINRTransaction({
    required this.txHash,
    required this.fromAddress,
    required this.toAddress,
    required this.amount,
    required this.memo,
    required this.timestamp,
    required this.isSuccess,
    this.errorMessage,
    required this.type,
    this.additionalData,
  });
  
  factory DINRTransaction.fromDeshChainTransaction(DeshChainTransaction tx) {
    // Determine transaction type and extract relevant data
    DINRTransactionType type;
    String fromAddress = '';
    String toAddress = '';
    BigInt amount = BigInt.zero;
    Map<String, dynamic>? additionalData;
    
    final firstMsg = tx.messages.first;
    switch (firstMsg.type) {
      case '/cosmos.bank.v1beta1.MsgSend':
        type = DINRTransactionType.transfer;
        fromAddress = firstMsg.data['from_address'] as String;
        toAddress = firstMsg.data['to_address'] as String;
        final amounts = firstMsg.data['amount'] as List;
        final dinrAmount = amounts.firstWhere(
          (amount) => amount['denom'] == 'dinr' || amount['denom'] == 'udinr',
          orElse: () => {'amount': '0'},
        );
        amount = BigInt.parse(dinrAmount['amount']);
        break;
        
      case '/deshchain.dinr.v1.MsgMintDINR':
        type = DINRTransactionType.mint;
        fromAddress = firstMsg.data['minter'] as String;
        toAddress = fromAddress; // Mint to self
        amount = BigInt.parse(firstMsg.data['dinr_to_mint']['amount']);
        additionalData = {
          'collateral': firstMsg.data['collateral'],
        };
        break;
        
      case '/deshchain.dinr.v1.MsgBurnDINR':
        type = DINRTransactionType.burn;
        fromAddress = firstMsg.data['burner'] as String;
        toAddress = fromAddress; // Burn from self
        amount = BigInt.parse(firstMsg.data['dinr_to_burn']['amount']);
        break;
        
      case '/deshchain.dinr.v1.MsgDepositCollateral':
        type = DINRTransactionType.depositCollateral;
        fromAddress = firstMsg.data['depositor'] as String;
        toAddress = 'DINR Module';
        amount = BigInt.zero; // No DINR involved
        additionalData = {
          'collateral': firstMsg.data['collateral'],
        };
        break;
        
      case '/deshchain.dinr.v1.MsgWithdrawCollateral':
        type = DINRTransactionType.withdrawCollateral;
        fromAddress = 'DINR Module';
        toAddress = firstMsg.data['withdrawer'] as String;
        amount = BigInt.zero; // No DINR involved
        additionalData = {
          'collateral': firstMsg.data['collateral'],
        };
        break;
        
      default:
        type = DINRTransactionType.other;
        fromAddress = 'Unknown';
        toAddress = 'Unknown';
        amount = BigInt.zero;
    }
    
    return DINRTransaction(
      txHash: tx.txHash,
      fromAddress: fromAddress,
      toAddress: toAddress,
      amount: amount,
      memo: tx.memo ?? '',
      timestamp: DateTime.parse(tx.timestamp),
      isSuccess: tx.success,
      errorMessage: tx.success ? null : 'Transaction failed',
      type: type,
      additionalData: additionalData,
    );
  }
  
  /// Get formatted amount string
  String get formattedAmount => DINRToken.formatAmount(amount);
  
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
  
  /// Get display description based on transaction type
  String getDisplayDescription() {
    switch (type) {
      case DINRTransactionType.transfer:
        return 'DINR Transfer';
      case DINRTransactionType.mint:
        return 'Mint DINR';
      case DINRTransactionType.burn:
        return 'Burn DINR';
      case DINRTransactionType.depositCollateral:
        return 'Deposit Collateral';
      case DINRTransactionType.withdrawCollateral:
        return 'Withdraw Collateral';
      case DINRTransactionType.liquidation:
        return 'Position Liquidated';
      case DINRTransactionType.yieldEarned:
        return 'Yield Earned';
      case DINRTransactionType.other:
        return 'Other Transaction';
    }
  }
}

/// DINR Stability Metrics
class DINRStabilityMetrics {
  final double currentPrice;
  final BigInt totalSupply;
  final double collateralRatio;
  final double stabilityFee;
  final double yieldAPY;
  final bool oracleStatus;
  final DateTime lastUpdated;
  
  DINRStabilityMetrics({
    required this.currentPrice,
    required this.totalSupply,
    required this.collateralRatio,
    required this.stabilityFee,
    required this.yieldAPY,
    required this.oracleStatus,
    required this.lastUpdated,
  });
  
  factory DINRStabilityMetrics.defaultMetrics() {
    return DINRStabilityMetrics(
      currentPrice: 1.00,
      totalSupply: BigInt.from(12500000000000), // 12.5M DINR
      collateralRatio: 156.0,
      stabilityFee: 0.1,
      yieldAPY: 5.2,
      oracleStatus: true,
      lastUpdated: DateTime.now(),
    );
  }
  
  /// Get formatted total supply
  String get formattedTotalSupply {
    final value = totalSupply.toDouble() / 1000000;
    if (value >= 1000000) {
      return '₹${(value / 1000000).toStringAsFixed(1)}M';
    } else if (value >= 1000) {
      return '₹${(value / 1000).toStringAsFixed(1)}K';
    } else {
      return '₹${value.toStringAsFixed(0)}';
    }
  }
}

/// Collateral Asset model
class CollateralAsset {
  final String denom;
  final String name;
  final String symbol;
  final int decimals;
  final int collateralRatio;
  final int liquidationRatio;
  final BigInt minCollateralAmount;
  final Color color;
  
  CollateralAsset({
    required this.denom,
    required this.name,
    required this.symbol,
    required this.decimals,
    required this.collateralRatio,
    required this.liquidationRatio,
    required this.minCollateralAmount,
    required this.color,
  });
  
  /// Format collateral amount
  String formatAmount(BigInt amount) {
    final value = amount.toDouble() / BigInt.from(10).pow(decimals).toDouble();
    return value.toStringAsFixed(decimals <= 6 ? decimals : 6);
  }
}

/// Collateral Position model
class CollateralPosition {
  final String userAddress;
  final String collateralDenom;
  final BigInt collateralAmount;
  final BigInt dinrMinted;
  final double collateralRatio;
  final bool isHealthy;
  final DateTime lastUpdated;
  
  CollateralPosition({
    required this.userAddress,
    required this.collateralDenom,
    required this.collateralAmount,
    required this.dinrMinted,
    required this.collateralRatio,
    required this.isHealthy,
    required this.lastUpdated,
  });
  
  /// Get health status color
  Color get healthColor {
    if (collateralRatio >= 150) {
      return Colors.green;
    } else if (collateralRatio >= 130) {
      return Colors.orange;
    } else {
      return Colors.red;
    }
  }
  
  /// Get health status text
  String get healthStatus {
    if (collateralRatio >= 150) {
      return 'Healthy';
    } else if (collateralRatio >= 130) {
      return 'Warning';
    } else {
      return 'At Risk';
    }
  }
}

/// DINR Transaction types
enum DINRTransactionType {
  transfer,
  mint,
  burn,
  depositCollateral,
  withdrawCollateral,
  liquidation,
  yieldEarned,
  other,
}

/// Transaction direction (reused from NAMO)
enum TransactionDirection {
  sent,
  received,
  unknown,
}

/// DINR Exception
class DINRException implements Exception {
  final String message;
  
  DINRException(this.message);
  
  @override
  String toString() => 'DINRException: $message';
}

// Placeholder classes for DeshChain client integration
class DeshChainTransaction {
  final String txHash;
  final List<TransactionMessage> messages;
  final String? memo;
  final String timestamp;
  final bool success;
  
  DeshChainTransaction({
    required this.txHash,
    required this.messages,
    this.memo,
    required this.timestamp,
    required this.success,
  });
}

class TransactionMessage {
  final String type;
  final Map<String, dynamic> data;
  
  TransactionMessage({
    required this.type,
    required this.data,
  });
}

class DINRPriceData {
  final double priceINR;
  final double priceUSD;
  final DateTime timestamp;
  
  DINRPriceData({
    required this.priceINR,
    required this.priceUSD,
    required this.timestamp,
  });
}