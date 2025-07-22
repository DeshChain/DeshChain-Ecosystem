import 'dart:convert';
import 'dart:typed_data';
import 'package:dio/dio.dart';
import 'package:web3dart/web3dart.dart';
import 'package:flutter/material.dart';

import '../../utils/logger.dart';
import '../../utils/constants.dart';
import '../wallet/hd_wallet.dart';

/// DeshChain client for blockchain interactions
class DeshChainClient {
  static const String _mainnetRPC = 'https://rpc.deshchain.org';
  static const String _testnetRPC = 'https://testnet-rpc.deshchain.org';
  
  final Dio _dio;
  final String _rpcUrl;
  final bool _isTestnet;
  
  DeshChainClient({
    bool isTestnet = false,
    String? customRPC,
  }) : _isTestnet = isTestnet,
       _rpcUrl = customRPC ?? (isTestnet ? _testnetRPC : _mainnetRPC),
       _dio = Dio() {
    _dio.options.baseUrl = _rpcUrl;
    _dio.options.connectTimeout = const Duration(seconds: 30);
    _dio.options.receiveTimeout = const Duration(seconds: 30);
    _dio.options.headers = {
      'Content-Type': 'application/json',
      'User-Agent': 'Batua-Wallet/1.0.0',
    };
  }
  
  /// Get account balance
  Future<DeshChainBalance> getBalance(String address) async {
    try {
      final response = await _dio.get('/cosmos/bank/v1beta1/balances/$address');
      
      if (response.statusCode == 200) {
        final data = response.data;
        final balances = data['balances'] as List<dynamic>;
        
        BigInt namoBalance = BigInt.zero;
        List<TokenBalance> otherTokens = [];
        
        for (final balance in balances) {
          final denom = balance['denom'] as String;
          final amount = BigInt.parse(balance['amount'] as String);
          
          if (denom == 'namo' || denom == 'unamo') {
            namoBalance = amount;
          } else {
            otherTokens.add(TokenBalance(
              denom: denom,
              amount: amount,
              decimals: _getTokenDecimals(denom),
            ));
          }
        }
        
        return DeshChainBalance(
          address: address,
          namoBalance: namoBalance,
          otherTokens: otherTokens,
        );
      } else {
        throw DeshChainException('Failed to get balance: ${response.statusCode}');
      }
    } catch (e) {
      AppLogger.error('Error getting balance: $e');
      throw DeshChainException('Failed to get balance: $e');
    }
  }
  
  /// Get transaction history
  Future<List<DeshChainTransaction>> getTransactionHistory(
    String address, {
    int limit = 50,
    int offset = 0,
  }) async {
    try {
      final response = await _dio.get(
        '/cosmos/tx/v1beta1/txs',
        queryParameters: {
          'events': 'transfer.recipient=$address',
          'pagination.limit': limit,
          'pagination.offset': offset,
          'order_by': 'ORDER_BY_DESC',
        },
      );
      
      if (response.statusCode == 200) {
        final data = response.data;
        final txs = data['tx_responses'] as List<dynamic>;
        
        return txs.map((tx) => DeshChainTransaction.fromJson(tx)).toList();
      } else {
        throw DeshChainException('Failed to get transactions: ${response.statusCode}');
      }
    } catch (e) {
      AppLogger.error('Error getting transactions: $e');
      throw DeshChainException('Failed to get transactions: $e');
    }
  }
  
  /// Send NAMO tokens
  Future<String> sendNAMO({
    required String fromAddress,
    required String toAddress,
    required BigInt amount,
    required String memo,
    required HDWallet wallet,
    required int accountIndex,
  }) async {
    try {
      // Get account info
      final accountInfo = await _getAccountInfo(fromAddress);
      
      // Prepare transaction
      final tx = await _prepareSendTransaction(
        fromAddress: fromAddress,
        toAddress: toAddress,
        amount: amount,
        memo: memo,
        accountNumber: accountInfo.accountNumber,
        sequence: accountInfo.sequence,
      );
      
      // Sign transaction
      final signedTx = await _signTransaction(tx, wallet, accountIndex);
      
      // Broadcast transaction
      final txHash = await _broadcastTransaction(signedTx);
      
      AppLogger.info('Transaction sent successfully: $txHash');
      return txHash;
    } catch (e) {
      AppLogger.error('Error sending NAMO: $e');
      throw DeshChainException('Failed to send NAMO: $e');
    }
  }
  
  /// Get account info
  Future<AccountInfo> _getAccountInfo(String address) async {
    final response = await _dio.get('/cosmos/auth/v1beta1/accounts/$address');
    
    if (response.statusCode == 200) {
      final data = response.data['account'];
      return AccountInfo(
        address: data['address'],
        accountNumber: int.parse(data['account_number']),
        sequence: int.parse(data['sequence']),
      );
    } else {
      throw DeshChainException('Failed to get account info');
    }
  }
  
  /// Prepare send transaction
  Future<Map<String, dynamic>> _prepareSendTransaction({
    required String fromAddress,
    required String toAddress,
    required BigInt amount,
    required String memo,
    required int accountNumber,
    required int sequence,
  }) async {
    final chainId = _isTestnet ? 'deshchain-testnet-1' : 'deshchain-1';
    
    return {
      'body': {
        'messages': [
          {
            '@type': '/cosmos.bank.v1beta1.MsgSend',
            'from_address': fromAddress,
            'to_address': toAddress,
            'amount': [
              {
                'denom': 'namo',
                'amount': amount.toString(),
              }
            ],
          }
        ],
        'memo': memo,
        'timeout_height': '0',
        'extension_options': [],
        'non_critical_extension_options': [],
      },
      'auth_info': {
        'signer_infos': [],
        'fee': {
          'amount': [
            {
              'denom': 'namo',
              'amount': '5000', // 0.005 NAMO
            }
          ],
          'gas_limit': '200000',
          'payer': '',
          'granter': '',
        },
      },
      'signatures': [],
    };
  }
  
  /// Sign transaction
  Future<Map<String, dynamic>> _signTransaction(
    Map<String, dynamic> tx,
    HDWallet wallet,
    int accountIndex,
  ) async {
    // Create transaction bytes for signing
    final txBytes = _createTransactionBytes(tx);
    
    // Sign with wallet
    final signature = await wallet.signDeshChainTransaction(txBytes, accountIndex);
    
    // Add signature to transaction
    tx['signatures'] = [base64Encode(signature)];
    
    return tx;
  }
  
  /// Create transaction bytes for signing
  Uint8List _createTransactionBytes(Map<String, dynamic> tx) {
    // This is a simplified implementation
    // In production, use proper protobuf encoding
    final txString = jsonEncode(tx);
    return Uint8List.fromList(utf8.encode(txString));
  }
  
  /// Broadcast transaction
  Future<String> _broadcastTransaction(Map<String, dynamic> signedTx) async {
    final response = await _dio.post(
      '/cosmos/tx/v1beta1/txs',
      data: {
        'tx_bytes': base64Encode(
          Uint8List.fromList(utf8.encode(jsonEncode(signedTx))),
        ),
        'mode': 'BROADCAST_MODE_SYNC',
      },
    );
    
    if (response.statusCode == 200) {
      final data = response.data;
      final txResponse = data['tx_response'];
      
      if (txResponse['code'] == 0) {
        return txResponse['txhash'];
      } else {
        throw DeshChainException('Transaction failed: ${txResponse['raw_log']}');
      }
    } else {
      throw DeshChainException('Failed to broadcast transaction');
    }
  }
  
  /// Get DeshPay quote
  Future<DeshPayQuote> getDeshPayQuote() async {
    try {
      final response = await _dio.get('/deshpay/quote');
      
      if (response.statusCode == 200) {
        final data = response.data;
        return DeshPayQuote.fromJson(data);
      } else {
        throw DeshChainException('Failed to get DeshPay quote');
      }
    } catch (e) {
      AppLogger.error('Error getting DeshPay quote: $e');
      throw DeshChainException('Failed to get DeshPay quote: $e');
    }
  }
  
  /// Get staking information
  Future<List<Validator>> getValidators() async {
    try {
      final response = await _dio.get('/cosmos/staking/v1beta1/validators');
      
      if (response.statusCode == 200) {
        final data = response.data;
        final validators = data['validators'] as List<dynamic>;
        
        return validators.map((v) => Validator.fromJson(v)).toList();
      } else {
        throw DeshChainException('Failed to get validators');
      }
    } catch (e) {
      AppLogger.error('Error getting validators: $e');
      throw DeshChainException('Failed to get validators: $e');
    }
  }
  
  /// Delegate tokens
  Future<String> delegateTokens({
    required String delegatorAddress,
    required String validatorAddress,
    required BigInt amount,
    required HDWallet wallet,
    required int accountIndex,
  }) async {
    try {
      // Implementation similar to sendNAMO but with MsgDelegate
      // This is a placeholder implementation
      throw UnimplementedError('Delegation not yet implemented');
    } catch (e) {
      AppLogger.error('Error delegating tokens: $e');
      throw DeshChainException('Failed to delegate tokens: $e');
    }
  }
  
  /// Get token decimals
  int _getTokenDecimals(String denom) {
    switch (denom) {
      case 'namo':
      case 'unamo':
        return 6;
      case 'dinr':
      case 'udinr':
        return 6;
      default:
        return 6; // Default decimals
    }
  }
  
  // DINR-specific methods
  
  /// Get DINR balance for address
  Future<BigInt> getDINRBalance(String address) async {
    try {
      final balance = await getBalance(address);
      // Find DINR balance in other tokens
      final dinrToken = balance.otherTokens.firstWhere(
        (token) => token.denom == 'dinr' || token.denom == 'udinr',
        orElse: () => TokenBalance(denom: 'dinr', amount: BigInt.zero, decimals: 6),
      );
      return dinrToken.amount;
    } catch (e) {
      AppLogger.error('Error getting DINR balance: $e');
      return BigInt.zero;
    }
  }
  
  /// Get collateral positions for address
  Future<List<CollateralPosition>> getCollateralPositions(String address) async {
    try {
      final response = await _dio.get('/deshchain/dinr/v1/positions/$address');
      
      if (response.statusCode == 200) {
        final data = response.data;
        final positions = data['positions'] as List<dynamic>;
        
        return positions.map((p) => CollateralPosition.fromJson(p)).toList();
      } else {
        throw DeshChainException('Failed to get collateral positions');
      }
    } catch (e) {
      AppLogger.error('Error getting collateral positions: $e');
      return [];
    }
  }
  
  /// Mint DINR with collateral
  Future<String> mintDINR({
    required String userAddress,
    required String collateralDenom,
    required BigInt collateralAmount,
    required BigInt dinrToMint,
    required HDWallet wallet,
    required int accountIndex,
  }) async {
    try {
      // Get account info
      final accountInfo = await _getAccountInfo(userAddress);
      
      // Prepare mint transaction
      final tx = await _prepareMintDINRTransaction(
        userAddress: userAddress,
        collateralDenom: collateralDenom,
        collateralAmount: collateralAmount,
        dinrToMint: dinrToMint,
        accountNumber: accountInfo.accountNumber,
        sequence: accountInfo.sequence,
      );
      
      // Sign transaction
      final signedTx = await _signTransaction(tx, wallet, accountIndex);
      
      // Broadcast transaction
      final txHash = await _broadcastTransaction(signedTx);
      
      AppLogger.info('DINR minted successfully: $txHash');
      return txHash;
    } catch (e) {
      AppLogger.error('Error minting DINR: $e');
      throw DeshChainException('Failed to mint DINR: $e');
    }
  }
  
  /// Burn DINR to retrieve collateral
  Future<String> burnDINR({
    required String userAddress,
    required BigInt dinrToBurn,
    required HDWallet wallet,
    required int accountIndex,
  }) async {
    try {
      // Get account info
      final accountInfo = await _getAccountInfo(userAddress);
      
      // Prepare burn transaction
      final tx = await _prepareBurnDINRTransaction(
        userAddress: userAddress,
        dinrToBurn: dinrToBurn,
        accountNumber: accountInfo.accountNumber,
        sequence: accountInfo.sequence,
      );
      
      // Sign transaction
      final signedTx = await _signTransaction(tx, wallet, accountIndex);
      
      // Broadcast transaction
      final txHash = await _broadcastTransaction(signedTx);
      
      AppLogger.info('DINR burned successfully: $txHash');
      return txHash;
    } catch (e) {
      AppLogger.error('Error burning DINR: $e');
      throw DeshChainException('Failed to burn DINR: $e');
    }
  }
  
  /// Send DINR tokens
  Future<String> sendDINR({
    required String fromAddress,
    required String toAddress,
    required BigInt amount,
    required String memo,
    required HDWallet wallet,
    required int accountIndex,
  }) async {
    try {
      // Get account info
      final accountInfo = await _getAccountInfo(fromAddress);
      
      // Prepare transaction (similar to NAMO but with DINR denom)
      final tx = await _prepareSendDINRTransaction(
        fromAddress: fromAddress,
        toAddress: toAddress,
        amount: amount,
        memo: memo,
        accountNumber: accountInfo.accountNumber,
        sequence: accountInfo.sequence,
      );
      
      // Sign transaction
      final signedTx = await _signTransaction(tx, wallet, accountIndex);
      
      // Broadcast transaction
      final txHash = await _broadcastTransaction(signedTx);
      
      AppLogger.info('DINR sent successfully: $txHash');
      return txHash;
    } catch (e) {
      AppLogger.error('Error sending DINR: $e');
      throw DeshChainException('Failed to send DINR: $e');
    }
  }
  
  /// Get DINR price data
  Future<DINRPriceData> getDINRPriceData() async {
    try {
      final response = await _dio.get('/deshchain/dinr/v1/price');
      
      if (response.statusCode == 200) {
        final data = response.data;
        return DINRPriceData.fromJson(data);
      } else {
        throw DeshChainException('Failed to get DINR price data');
      }
    } catch (e) {
      AppLogger.error('Error getting DINR price data: $e');
      // Return default stable price
      return DINRPriceData(
        priceINR: 1.00,
        priceUSD: 0.012,
        timestamp: DateTime.now(),
      );
    }
  }
  
  /// Get DINR stability metrics
  Future<DINRStabilityMetrics> getDINRStabilityMetrics() async {
    try {
      final response = await _dio.get('/deshchain/dinr/v1/stability-metrics');
      
      if (response.statusCode == 200) {
        final data = response.data;
        return DINRStabilityMetrics.fromJson(data);
      } else {
        throw DeshChainException('Failed to get DINR stability metrics');
      }
    } catch (e) {
      AppLogger.error('Error getting DINR stability metrics: $e');
      return DINRStabilityMetrics.defaultMetrics();
    }
  }
  
  /// Get supported collateral assets
  Future<List<CollateralAsset>> getSupportedCollaterals() async {
    try {
      final response = await _dio.get('/deshchain/dinr/v1/supported-collaterals');
      
      if (response.statusCode == 200) {
        final data = response.data;
        final collaterals = data['collaterals'] as List<dynamic>;
        
        return collaterals.map((c) => CollateralAsset.fromJson(c)).toList();
      } else {
        throw DeshChainException('Failed to get supported collaterals');
      }
    } catch (e) {
      AppLogger.error('Error getting supported collaterals: $e');
      return [];
    }
  }
  
  /// Prepare mint DINR transaction
  Future<Map<String, dynamic>> _prepareMintDINRTransaction({
    required String userAddress,
    required String collateralDenom,
    required BigInt collateralAmount,
    required BigInt dinrToMint,
    required int accountNumber,
    required int sequence,
  }) async {
    final chainId = _isTestnet ? 'deshchain-testnet-1' : 'deshchain-1';
    
    return {
      'body': {
        'messages': [
          {
            '@type': '/deshchain.dinr.v1.MsgMintDINR',
            'minter': userAddress,
            'collateral': {
              'denom': collateralDenom,
              'amount': collateralAmount.toString(),
            },
            'dinr_to_mint': {
              'denom': 'dinr',
              'amount': dinrToMint.toString(),
            },
          }
        ],
        'memo': 'Mint DINR with collateral',
        'timeout_height': '0',
        'extension_options': [],
        'non_critical_extension_options': [],
      },
      'auth_info': {
        'signer_infos': [],
        'fee': {
          'amount': [
            {
              'denom': 'namo',
              'amount': '5000', // 0.005 NAMO gas fee
            }
          ],
          'gas_limit': '300000',
          'payer': '',
          'granter': '',
        },
      },
      'signatures': [],
    };
  }
  
  /// Prepare burn DINR transaction
  Future<Map<String, dynamic>> _prepareBurnDINRTransaction({
    required String userAddress,
    required BigInt dinrToBurn,
    required int accountNumber,
    required int sequence,
  }) async {
    final chainId = _isTestnet ? 'deshchain-testnet-1' : 'deshchain-1';
    
    return {
      'body': {
        'messages': [
          {
            '@type': '/deshchain.dinr.v1.MsgBurnDINR',
            'burner': userAddress,
            'dinr_to_burn': {
              'denom': 'dinr',
              'amount': dinrToBurn.toString(),
            },
          }
        ],
        'memo': 'Burn DINR to retrieve collateral',
        'timeout_height': '0',
        'extension_options': [],
        'non_critical_extension_options': [],
      },
      'auth_info': {
        'signer_infos': [],
        'fee': {
          'amount': [
            {
              'denom': 'namo',
              'amount': '5000', // 0.005 NAMO gas fee
            }
          ],
          'gas_limit': '250000',
          'payer': '',
          'granter': '',
        },
      },
      'signatures': [],
    };
  }
  
  /// Prepare send DINR transaction
  Future<Map<String, dynamic>> _prepareSendDINRTransaction({
    required String fromAddress,
    required String toAddress,
    required BigInt amount,
    required String memo,
    required int accountNumber,
    required int sequence,
  }) async {
    final chainId = _isTestnet ? 'deshchain-testnet-1' : 'deshchain-1';
    
    return {
      'body': {
        'messages': [
          {
            '@type': '/cosmos.bank.v1beta1.MsgSend',
            'from_address': fromAddress,
            'to_address': toAddress,
            'amount': [
              {
                'denom': 'dinr',
                'amount': amount.toString(),
              }
            ],
          }
        ],
        'memo': memo,
        'timeout_height': '0',
        'extension_options': [],
        'non_critical_extension_options': [],
      },
      'auth_info': {
        'signer_infos': [],
        'fee': {
          'amount': [
            {
              'denom': 'namo',
              'amount': '5000', // 0.005 NAMO gas fee
            }
          ],
          'gas_limit': '200000',
          'payer': '',
          'granter': '',
        },
      },
      'signatures': [],
    };
  }
  
  /// Check if address is valid
  bool isValidAddress(String address) {
    return address.startsWith('desh') && address.length == 45;
  }
  
  /// Get network info
  Future<NetworkInfo> getNetworkInfo() async {
    try {
      final response = await _dio.get('/cosmos/base/tendermint/v1beta1/node_info');
      
      if (response.statusCode == 200) {
        final data = response.data;
        return NetworkInfo.fromJson(data);
      } else {
        throw DeshChainException('Failed to get network info');
      }
    } catch (e) {
      AppLogger.error('Error getting network info: $e');
      throw DeshChainException('Failed to get network info: $e');
    }
  }
}

/// DeshChain balance model
class DeshChainBalance {
  final String address;
  final BigInt namoBalance;
  final List<TokenBalance> otherTokens;
  
  DeshChainBalance({
    required this.address,
    required this.namoBalance,
    required this.otherTokens,
  });
  
  double get namoBalanceFormatted => namoBalance.toDouble() / 1000000; // 6 decimals
}

/// Token balance model
class TokenBalance {
  final String denom;
  final BigInt amount;
  final int decimals;
  
  TokenBalance({
    required this.denom,
    required this.amount,
    required this.decimals,
  });
  
  double get formattedAmount => amount.toDouble() / (10 * decimals);
}

/// DeshChain transaction model
class DeshChainTransaction {
  final String txHash;
  final int height;
  final String timestamp;
  final bool success;
  final String? memo;
  final List<TransactionMessage> messages;
  
  DeshChainTransaction({
    required this.txHash,
    required this.height,
    required this.timestamp,
    required this.success,
    this.memo,
    required this.messages,
  });
  
  factory DeshChainTransaction.fromJson(Map<String, dynamic> json) {
    return DeshChainTransaction(
      txHash: json['txhash'],
      height: int.parse(json['height']),
      timestamp: json['timestamp'],
      success: json['code'] == 0,
      memo: json['tx']['body']['memo'],
      messages: (json['tx']['body']['messages'] as List<dynamic>)
          .map((m) => TransactionMessage.fromJson(m))
          .toList(),
    );
  }
}

/// Transaction message model
class TransactionMessage {
  final String type;
  final Map<String, dynamic> data;
  
  TransactionMessage({
    required this.type,
    required this.data,
  });
  
  factory TransactionMessage.fromJson(Map<String, dynamic> json) {
    return TransactionMessage(
      type: json['@type'],
      data: json,
    );
  }
}

/// Account info model
class AccountInfo {
  final String address;
  final int accountNumber;
  final int sequence;
  
  AccountInfo({
    required this.address,
    required this.accountNumber,
    required this.sequence,
  });
}

/// DeshPay quote model
class DeshPayQuote {
  final String text;
  final String author;
  final String category;
  final String language;
  
  DeshPayQuote({
    required this.text,
    required this.author,
    required this.category,
    required this.language,
  });
  
  factory DeshPayQuote.fromJson(Map<String, dynamic> json) {
    return DeshPayQuote(
      text: json['text'],
      author: json['author'],
      category: json['category'],
      language: json['language'],
    );
  }
}

/// Validator model
class Validator {
  final String operatorAddress;
  final String moniker;
  final String description;
  final double commission;
  final bool jailed;
  final String status;
  
  Validator({
    required this.operatorAddress,
    required this.moniker,
    required this.description,
    required this.commission,
    required this.jailed,
    required this.status,
  });
  
  factory Validator.fromJson(Map<String, dynamic> json) {
    return Validator(
      operatorAddress: json['operator_address'],
      moniker: json['description']['moniker'],
      description: json['description']['details'] ?? '',
      commission: double.parse(json['commission']['commission_rates']['rate']),
      jailed: json['jailed'],
      status: json['status'],
    );
  }
}

/// Network info model
class NetworkInfo {
  final String chainId;
  final String nodeId;
  final String version;
  final int latestBlockHeight;
  
  NetworkInfo({
    required this.chainId,
    required this.nodeId,
    required this.version,
    required this.latestBlockHeight,
  });
  
  factory NetworkInfo.fromJson(Map<String, dynamic> json) {
    return NetworkInfo(
      chainId: json['default_node_info']['network'],
      nodeId: json['default_node_info']['id'],
      version: json['default_node_info']['version'],
      latestBlockHeight: int.parse(json['sync_info']['latest_block_height']),
    );
  }
}

/// Collateral Position model for DeshChain client
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
  
  factory CollateralPosition.fromJson(Map<String, dynamic> json) {
    return CollateralPosition(
      userAddress: json['user_address'],
      collateralDenom: json['collateral_denom'],
      collateralAmount: BigInt.parse(json['collateral_amount']),
      dinrMinted: BigInt.parse(json['dinr_minted']),
      collateralRatio: double.parse(json['collateral_ratio']),
      isHealthy: json['is_healthy'],
      lastUpdated: DateTime.parse(json['last_updated']),
    );
  }
}

/// DINR Price Data model
class DINRPriceData {
  final double priceINR;
  final double priceUSD;
  final DateTime timestamp;
  
  DINRPriceData({
    required this.priceINR,
    required this.priceUSD,
    required this.timestamp,
  });
  
  factory DINRPriceData.fromJson(Map<String, dynamic> json) {
    return DINRPriceData(
      priceINR: double.parse(json['price_inr']),
      priceUSD: double.parse(json['price_usd']),
      timestamp: DateTime.parse(json['timestamp']),
    );
  }
}

/// DINR Stability Metrics model
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
  
  factory DINRStabilityMetrics.fromJson(Map<String, dynamic> json) {
    return DINRStabilityMetrics(
      currentPrice: double.parse(json['current_price']),
      totalSupply: BigInt.parse(json['total_supply']),
      collateralRatio: double.parse(json['collateral_ratio']),
      stabilityFee: double.parse(json['stability_fee']),
      yieldAPY: double.parse(json['yield_apy']),
      oracleStatus: json['oracle_status'],
      lastUpdated: DateTime.parse(json['last_updated']),
    );
  }
  
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

/// Collateral Asset model for DeshChain client
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
  
  factory CollateralAsset.fromJson(Map<String, dynamic> json) {
    return CollateralAsset(
      denom: json['denom'],
      name: json['name'],
      symbol: json['symbol'],
      decimals: json['decimals'],
      collateralRatio: json['collateral_ratio'],
      liquidationRatio: json['liquidation_ratio'],
      minCollateralAmount: BigInt.parse(json['min_collateral_amount']),
      color: Color(int.parse(json['color'], radix: 16)),
    );
  }
  
  /// Format collateral amount
  String formatAmount(BigInt amount) {
    final value = amount.toDouble() / BigInt.from(10).pow(decimals).toDouble();
    return value.toStringAsFixed(decimals <= 6 ? decimals : 6);
  }
}

/// DeshChain exception
class DeshChainException implements Exception {
  final String message;
  
  DeshChainException(this.message);
  
  @override
  String toString() => 'DeshChainException: $message';
}