import 'dart:convert';
import 'dart:typed_data';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:crypto/crypto.dart';
import 'package:pointycastle/export.dart';

import '../../utils/logger.dart';

/// Secure storage wrapper with encryption and cultural context
class SecureStorage {
  static const FlutterSecureStorage _storage = FlutterSecureStorage(
    aOptions: AndroidOptions(
      encryptedSharedPreferences: true,
      sharedPreferencesName: 'batua_secure_prefs',
      preferencesKeyPrefix: 'batua_',
    ),
    iOptions: IOSOptions(
      groupId: 'group.com.deshchain.batua',
      accountName: 'Batua Wallet',
      accessibility: IOSAccessibility.first_unlock_this_device,
    ),
  );
  
  static const String _masterKeyKey = 'master_key';
  static const String _saltKey = 'salt';
  static const String _culturalPrefixKey = 'cultural_prefix';
  
  static Uint8List? _masterKey;
  static String? _culturalPrefix;
  
  /// Initialize secure storage
  static Future<void> init() async {
    try {
      await _initializeMasterKey();
      await _initializeCulturalPrefix();
      AppLogger.info('Secure storage initialized');
    } catch (e) {
      AppLogger.error('Failed to initialize secure storage: $e');
      rethrow;
    }
  }
  
  /// Initialize master key for encryption
  static Future<void> _initializeMasterKey() async {
    String? masterKeyHex = await _storage.read(key: _masterKeyKey);
    String? saltHex = await _storage.read(key: _saltKey);
    
    if (masterKeyHex == null || saltHex == null) {
      // Generate new master key and salt
      final salt = _generateRandomBytes(32);
      final masterKey = _generateRandomBytes(32);
      
      await _storage.write(key: _masterKeyKey, value: hex.encode(masterKey));
      await _storage.write(key: _saltKey, value: hex.encode(salt));
      
      _masterKey = masterKey;
    } else {
      _masterKey = Uint8List.fromList(hex.decode(masterKeyHex));
    }
  }
  
  /// Initialize cultural prefix for keys
  static Future<void> _initializeCulturalPrefix() async {
    _culturalPrefix = await _storage.read(key: _culturalPrefixKey);
    
    if (_culturalPrefix == null) {
      // Generate cultural prefix based on Sanskrit/Hindi
      final culturalWords = [
        'सुरक्षा', 'गुप्त', 'रहस्य', 'सुरक्षित', 'संरक्षित',
        'गोपनीय', 'निजी', 'व्यक्तिगत', 'सुरक्षा_कवच'
      ];
      
      final random = Random.secure();
      _culturalPrefix = culturalWords[random.nextInt(culturalWords.length)];
      
      await _storage.write(key: _culturalPrefixKey, value: _culturalPrefix!);
    }
  }
  
  /// Generate random bytes
  static Uint8List _generateRandomBytes(int length) {
    final random = Random.secure();
    return Uint8List.fromList(
      List<int>.generate(length, (i) => random.nextInt(256)),
    );
  }
  
  /// Get cultural key with prefix
  static String _getCulturalKey(String key) {
    return '${_culturalPrefix}_$key';
  }
  
  /// Encrypt data using AES-256-GCM
  static Uint8List _encrypt(String data) {
    if (_masterKey == null) {
      throw SecurityException('Master key not initialized');
    }
    
    final plainBytes = utf8.encode(data);
    final iv = _generateRandomBytes(12); // 96-bit IV for GCM
    
    final cipher = GCMBlockCipher(AESFastEngine());
    final params = AEADParameters(
      KeyParameter(_masterKey!),
      128, // 128-bit authentication tag
      iv,
    );
    
    cipher.init(true, params);
    
    final cipherBytes = cipher.process(plainBytes);
    
    // Combine IV and ciphertext
    final result = Uint8List(iv.length + cipherBytes.length);
    result.setRange(0, iv.length, iv);
    result.setRange(iv.length, result.length, cipherBytes);
    
    return result;
  }
  
  /// Decrypt data using AES-256-GCM
  static String _decrypt(Uint8List encryptedData) {
    if (_masterKey == null) {
      throw SecurityException('Master key not initialized');
    }
    
    if (encryptedData.length < 12) {
      throw SecurityException('Invalid encrypted data');
    }
    
    final iv = encryptedData.sublist(0, 12);
    final cipherBytes = encryptedData.sublist(12);
    
    final cipher = GCMBlockCipher(AESFastEngine());
    final params = AEADParameters(
      KeyParameter(_masterKey!),
      128,
      iv,
    );
    
    cipher.init(false, params);
    
    try {
      final plainBytes = cipher.process(cipherBytes);
      return utf8.decode(plainBytes);
    } catch (e) {
      throw SecurityException('Failed to decrypt data: $e');
    }
  }
  
  /// Write encrypted data to secure storage
  static Future<void> write(String key, String value) async {
    try {
      final culturalKey = _getCulturalKey(key);
      final encryptedData = _encrypt(value);
      final encodedData = base64Encode(encryptedData);
      
      await _storage.write(key: culturalKey, value: encodedData);
      AppLogger.debug('Secure data written for key: $key');
    } catch (e) {
      AppLogger.error('Failed to write secure data: $e');
      throw SecurityException('Failed to write secure data');
    }
  }
  
  /// Read encrypted data from secure storage
  static Future<String?> read(String key) async {
    try {
      final culturalKey = _getCulturalKey(key);
      final encodedData = await _storage.read(key: culturalKey);
      
      if (encodedData == null) {
        return null;
      }
      
      final encryptedData = base64Decode(encodedData);
      final decryptedValue = _decrypt(encryptedData);
      
      AppLogger.debug('Secure data read for key: $key');
      return decryptedValue;
    } catch (e) {
      AppLogger.error('Failed to read secure data: $e');
      return null;
    }
  }
  
  /// Delete data from secure storage
  static Future<void> delete(String key) async {
    try {
      final culturalKey = _getCulturalKey(key);
      await _storage.delete(key: culturalKey);
      AppLogger.debug('Secure data deleted for key: $key');
    } catch (e) {
      AppLogger.error('Failed to delete secure data: $e');
      throw SecurityException('Failed to delete secure data');
    }
  }
  
  /// Check if key exists in secure storage
  static Future<bool> containsKey(String key) async {
    try {
      final culturalKey = _getCulturalKey(key);
      return await _storage.containsKey(key: culturalKey);
    } catch (e) {
      AppLogger.error('Failed to check key existence: $e');
      return false;
    }
  }
  
  /// Get all keys from secure storage
  static Future<Map<String, String>> readAll() async {
    try {
      final allData = await _storage.readAll();
      final result = <String, String>{};
      
      for (final entry in allData.entries) {
        if (entry.key.startsWith('${_culturalPrefix}_')) {
          final originalKey = entry.key.substring('${_culturalPrefix}_'.length);
          try {
            final encryptedData = base64Decode(entry.value);
            final decryptedValue = _decrypt(encryptedData);
            result[originalKey] = decryptedValue;
          } catch (e) {
            AppLogger.error('Failed to decrypt key ${entry.key}: $e');
          }
        }
      }
      
      return result;
    } catch (e) {
      AppLogger.error('Failed to read all secure data: $e');
      return {};
    }
  }
  
  /// Clear all secure storage data
  static Future<void> deleteAll() async {
    try {
      await _storage.deleteAll();
      AppLogger.info('All secure data cleared');
    } catch (e) {
      AppLogger.error('Failed to clear secure data: $e');
      throw SecurityException('Failed to clear secure data');
    }
  }
  
  /// Write JSON data to secure storage
  static Future<void> writeJson(String key, Map<String, dynamic> json) async {
    final jsonString = jsonEncode(json);
    await write(key, jsonString);
  }
  
  /// Read JSON data from secure storage
  static Future<Map<String, dynamic>?> readJson(String key) async {
    final jsonString = await read(key);
    if (jsonString == null) return null;
    
    try {
      return jsonDecode(jsonString) as Map<String, dynamic>;
    } catch (e) {
      AppLogger.error('Failed to decode JSON for key $key: $e');
      return null;
    }
  }
  
  /// Store cultural preference
  static Future<void> writeCulturalPreference(String key, String value) async {
    final culturalKey = 'cultural_$key';
    await write(culturalKey, value);
  }
  
  /// Read cultural preference
  static Future<String?> readCulturalPreference(String key) async {
    final culturalKey = 'cultural_$key';
    return await read(culturalKey);
  }
  
  /// Store user authentication data
  static Future<void> writeAuthData(String key, String value) async {
    final authKey = 'auth_$key';
    await write(authKey, value);
  }
  
  /// Read user authentication data
  static Future<String?> readAuthData(String key) async {
    final authKey = 'auth_$key';
    return await read(authKey);
  }
  
  /// Backup secure storage to encrypted string
  static Future<String> createBackup() async {
    try {
      final allData = await readAll();
      final backup = {
        'timestamp': DateTime.now().toIso8601String(),
        'version': '1.0.0',
        'data': allData,
        'cultural_prefix': _culturalPrefix,
      };
      
      final backupString = jsonEncode(backup);
      final encryptedBackup = _encrypt(backupString);
      
      return base64Encode(encryptedBackup);
    } catch (e) {
      AppLogger.error('Failed to create backup: $e');
      throw SecurityException('Failed to create backup');
    }
  }
  
  /// Restore secure storage from encrypted backup
  static Future<void> restoreFromBackup(String backupString) async {
    try {
      final encryptedBackup = base64Decode(backupString);
      final decryptedString = _decrypt(encryptedBackup);
      final backup = jsonDecode(decryptedString);
      
      // Validate backup structure
      if (backup['version'] != '1.0.0') {
        throw SecurityException('Incompatible backup version');
      }
      
      await deleteAll();
      
      final data = backup['data'] as Map<String, dynamic>;
      for (final entry in data.entries) {
        await write(entry.key, entry.value);
      }
      
      AppLogger.info('Backup restored successfully');
    } catch (e) {
      AppLogger.error('Failed to restore backup: $e');
      throw SecurityException('Failed to restore backup');
    }
  }
}

/// Security exception class
class SecurityException implements Exception {
  final String message;
  
  SecurityException(this.message);
  
  @override
  String toString() => 'SecurityException: $message';
}