import 'dart:convert';
import 'dart:math';
import 'dart:typed_data';
import 'package:crypto/crypto.dart';
import 'package:bip39/bip39.dart' as bip39;
import 'package:bip32/bip32.dart' as bip32;
import 'package:ed25519_hd_key/ed25519_hd_key.dart';
import 'package:pointycastle/export.dart';
import 'package:web3dart/web3dart.dart';

import '../storage/secure_storage.dart';
import '../../utils/constants.dart';
import '../../utils/logger.dart';

/// HD Wallet implementation following BIP32/BIP39/BIP44 standards
/// Supporting DeshChain, Ethereum, and Bitcoin
class HDWallet {
  static const String _mnemonicKey = 'wallet_mnemonic';
  static const String _seedKey = 'wallet_seed';
  static const String _walletIdKey = 'wallet_id';
  
  // BIP44 coin types
  static const int _ethereumCoinType = 60;
  static const int _bitcoinCoinType = 0;
  static const int _deshchainCoinType = 118; // Cosmos standard
  
  String? _mnemonic;
  Uint8List? _seed;
  String? _walletId;
  
  // Cached keys for performance
  Map<String, EthereumAddress> _ethereumAddresses = {};
  Map<String, ECPrivateKey> _privateKeys = {};
  
  /// Generate new wallet with mnemonic
  static Future<HDWallet> generateWallet({
    int strength = 256, // 24 words
    String? passphrase,
  }) async {
    final mnemonic = bip39.generateMnemonic(strength: strength);
    return await _createWalletFromMnemonic(mnemonic, passphrase);
  }
  
  /// Import wallet from mnemonic
  static Future<HDWallet> importFromMnemonic(
    String mnemonic, {
    String? passphrase,
  }) async {
    if (!bip39.validateMnemonic(mnemonic)) {
      throw WalletException('Invalid mnemonic phrase');
    }
    return await _createWalletFromMnemonic(mnemonic, passphrase);
  }
  
  /// Load existing wallet from secure storage
  static Future<HDWallet?> loadWallet() async {
    try {
      final mnemonic = await SecureStorage.read(_mnemonicKey);
      final seedHex = await SecureStorage.read(_seedKey);
      final walletId = await SecureStorage.read(_walletIdKey);
      
      if (mnemonic == null || seedHex == null || walletId == null) {
        return null;
      }
      
      final wallet = HDWallet._();
      wallet._mnemonic = mnemonic;
      wallet._seed = Uint8List.fromList(hex.decode(seedHex));
      wallet._walletId = walletId;
      
      AppLogger.info('Wallet loaded successfully');
      return wallet;
    } catch (e) {
      AppLogger.error('Failed to load wallet: $e');
      return null;
    }
  }
  
  /// Create wallet from mnemonic
  static Future<HDWallet> _createWalletFromMnemonic(
    String mnemonic,
    String? passphrase,
  ) async {
    final seed = bip39.mnemonicToSeed(mnemonic, passphrase: passphrase ?? '');
    final walletId = _generateWalletId(seed);
    
    final wallet = HDWallet._();
    wallet._mnemonic = mnemonic;
    wallet._seed = seed;
    wallet._walletId = walletId;
    
    await wallet._saveToSecureStorage();
    
    AppLogger.info('Wallet created successfully');
    return wallet;
  }
  
  HDWallet._();
  
  /// Save wallet to secure storage
  Future<void> _saveToSecureStorage() async {
    if (_mnemonic == null || _seed == null || _walletId == null) {
      throw WalletException('Wallet not initialized');
    }
    
    await SecureStorage.write(_mnemonicKey, _mnemonic!);
    await SecureStorage.write(_seedKey, hex.encode(_seed!));
    await SecureStorage.write(_walletIdKey, _walletId!);
  }
  
  /// Generate unique wallet ID
  static String _generateWalletId(Uint8List seed) {
    final hash = sha256.convert(seed);
    return 'batua_${hash.toString().substring(0, 16)}';
  }
  
  /// Get wallet mnemonic (requires authentication)
  String? get mnemonic => _mnemonic;
  
  /// Get wallet ID
  String? get walletId => _walletId;
  
  /// Check if wallet is initialized
  bool get isInitialized => _mnemonic != null && _seed != null;
  
  /// Get Ethereum address for account
  EthereumAddress getEthereumAddress(int accountIndex) {
    final cacheKey = 'eth_$accountIndex';
    if (_ethereumAddresses.containsKey(cacheKey)) {
      return _ethereumAddresses[cacheKey]!;
    }
    
    final privateKey = _deriveEthereumPrivateKey(accountIndex);
    final address = EthereumAddress.fromPublicKey(privateKey.publicKey.getEncoded());
    
    _ethereumAddresses[cacheKey] = address;
    return address;
  }
  
  /// Get DeshChain address for account
  String getDeshChainAddress(int accountIndex) {
    final privateKey = _deriveDeshChainPrivateKey(accountIndex);
    return _generateDeshChainAddress(privateKey);
  }
  
  /// Get Bitcoin address for account
  String getBitcoinAddress(int accountIndex) {
    final privateKey = _deriveBitcoinPrivateKey(accountIndex);
    return _generateBitcoinAddress(privateKey);
  }
  
  /// Sign Ethereum transaction
  Future<Uint8List> signEthereumTransaction(
    Transaction transaction,
    int accountIndex,
  ) async {
    final privateKey = _deriveEthereumPrivateKey(accountIndex);
    final credentials = EthPrivateKey(privateKey.d!);
    
    return await credentials.signTransaction(transaction);
  }
  
  /// Sign DeshChain transaction
  Future<Uint8List> signDeshChainTransaction(
    Uint8List transactionBytes,
    int accountIndex,
  ) async {
    final privateKey = _deriveDeshChainPrivateKey(accountIndex);
    return _signWithEd25519(transactionBytes, privateKey);
  }
  
  /// Sign message with Ethereum key
  Future<Uint8List> signMessage(
    Uint8List message,
    int accountIndex,
  ) async {
    final privateKey = _deriveEthereumPrivateKey(accountIndex);
    final credentials = EthPrivateKey(privateKey.d!);
    
    return await credentials.signPersonalMessage(message);
  }
  
  /// Derive Ethereum private key using BIP44
  ECPrivateKey _deriveEthereumPrivateKey(int accountIndex) {
    final cacheKey = 'eth_priv_$accountIndex';
    if (_privateKeys.containsKey(cacheKey)) {
      return _privateKeys[cacheKey]!;
    }
    
    final path = "m/44'/$_ethereumCoinType'/$accountIndex'/0/0";
    final node = bip32.BIP32.fromSeed(_seed!);
    final derived = node.derivePath(path);
    
    final privateKey = ECPrivateKey(
      BigInt.parse(hex.encode(derived.privateKey!), radix: 16),
      ECDomainParameters('secp256k1'),
    );
    
    _privateKeys[cacheKey] = privateKey;
    return privateKey;
  }
  
  /// Derive DeshChain private key using BIP44
  Uint8List _deriveDeshChainPrivateKey(int accountIndex) {
    final path = "m/44'/$_deshchainCoinType'/$accountIndex'/0/0";
    final masterKey = ED25519_HD_KEY.getMasterKeyFromSeed(_seed!);
    final derived = ED25519_HD_KEY.derivePath(path, masterKey);
    
    return derived.key;
  }
  
  /// Derive Bitcoin private key using BIP44
  ECPrivateKey _deriveBitcoinPrivateKey(int accountIndex) {
    final path = "m/44'/$_bitcoinCoinType'/$accountIndex'/0/0";
    final node = bip32.BIP32.fromSeed(_seed!);
    final derived = node.derivePath(path);
    
    return ECPrivateKey(
      BigInt.parse(hex.encode(derived.privateKey!), radix: 16),
      ECDomainParameters('secp256k1'),
    );
  }
  
  /// Generate DeshChain address from private key
  String _generateDeshChainAddress(Uint8List privateKey) {
    // Implementation depends on DeshChain address format
    // This is a placeholder implementation
    final publicKey = ED25519_HD_KEY.getPublicKey(privateKey, false);
    final hash = sha256.convert(publicKey);
    return 'desh${_bech32Encode(hash.bytes.take(20).toList())}';
  }
  
  /// Generate Bitcoin address from private key
  String _generateBitcoinAddress(ECPrivateKey privateKey) {
    // Simplified Bitcoin address generation
    // In production, use a proper Bitcoin library
    final publicKey = privateKey.publicKey.getEncoded();
    final hash = sha256.convert(publicKey);
    return '1${_base58Encode(hash.bytes.take(20).toList())}';
  }
  
  /// Sign with Ed25519 for DeshChain
  Uint8List _signWithEd25519(Uint8List message, Uint8List privateKey) {
    final signer = Ed25519Signer();
    final keyPair = AsymmetricKeyPair(
      Ed25519PublicKey(ED25519_HD_KEY.getPublicKey(privateKey, false)),
      Ed25519PrivateKey(privateKey),
    );
    
    signer.init(true, PrivateKeyParameter(keyPair.privateKey));
    return signer.generateSignature(message).bytes;
  }
  
  /// Simplified Bech32 encoding for DeshChain addresses
  String _bech32Encode(List<int> data) {
    // This is a simplified implementation
    // In production, use a proper Bech32 library
    return base64Encode(data).replaceAll('/', '').replaceAll('+', '').toLowerCase();
  }
  
  /// Simplified Base58 encoding for Bitcoin addresses
  String _base58Encode(List<int> data) {
    // This is a simplified implementation
    // In production, use a proper Base58 library
    return base64Encode(data).replaceAll('/', '').replaceAll('+', '');
  }
  
  /// Delete wallet from secure storage
  Future<void> deleteWallet() async {
    await SecureStorage.delete(_mnemonicKey);
    await SecureStorage.delete(_seedKey);
    await SecureStorage.delete(_walletIdKey);
    
    _mnemonic = null;
    _seed = null;
    _walletId = null;
    _ethereumAddresses.clear();
    _privateKeys.clear();
    
    AppLogger.info('Wallet deleted successfully');
  }
  
  /// Export wallet as JSON (requires authentication)
  Map<String, dynamic> exportWallet() {
    if (!isInitialized) {
      throw WalletException('Wallet not initialized');
    }
    
    return {
      'walletId': _walletId,
      'mnemonic': _mnemonic,
      'created': DateTime.now().toIso8601String(),
      'version': '1.0.0',
      'type': 'HD',
    };
  }
  
  /// Get wallet summary
  Map<String, dynamic> getWalletSummary() {
    if (!isInitialized) {
      throw WalletException('Wallet not initialized');
    }
    
    return {
      'walletId': _walletId,
      'type': 'HD Wallet',
      'supportedChains': ['DeshChain', 'Ethereum', 'Bitcoin'],
      'accountCount': 1, // Can be extended for multiple accounts
      'created': DateTime.now().toIso8601String(),
    };
  }
}

/// Wallet exception class
class WalletException implements Exception {
  final String message;
  
  WalletException(this.message);
  
  @override
  String toString() => 'WalletException: $message';
}

/// Account information
class Account {
  final int index;
  final String name;
  final String address;
  final String? balance;
  final bool isActive;
  
  Account({
    required this.index,
    required this.name,
    required this.address,
    this.balance,
    this.isActive = true,
  });
  
  factory Account.fromJson(Map<String, dynamic> json) {
    return Account(
      index: json['index'],
      name: json['name'],
      address: json['address'],
      balance: json['balance'],
      isActive: json['isActive'] ?? true,
    );
  }
  
  Map<String, dynamic> toJson() {
    return {
      'index': index,
      'name': name,
      'address': address,
      'balance': balance,
      'isActive': isActive,
    };
  }
}