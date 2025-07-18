import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:qr_flutter/qr_flutter.dart';
import 'package:share_plus/share_plus.dart';

import '../../widgets/cultural_gradient_text.dart';
import '../../widgets/namo_balance_widget.dart';
import '../../../core/tokens/namo_token.dart';
import '../../../utils/constants.dart';
import '../../../utils/logger.dart';

/// NAMO Receive Screen
class NAMOReceiveScreen extends ConsumerStatefulWidget {
  final String address;
  final String? displayName;
  
  const NAMOReceiveScreen({
    super.key,
    required this.address,
    this.displayName,
  });
  
  @override
  ConsumerState<NAMOReceiveScreen> createState() => _NAMOReceiveScreenState();
}

class _NAMOReceiveScreenState extends ConsumerState<NAMOReceiveScreen> 
    with SingleTickerProviderStateMixin {
  late AnimationController _animationController;
  late Animation<double> _scaleAnimation;
  late Animation<double> _rotationAnimation;
  
  final TextEditingController _amountController = TextEditingController();
  final TextEditingController _memoController = TextEditingController();
  
  bool _showCustomAmount = false;
  BigInt _requestedAmount = BigInt.zero;
  String _qrData = '';
  
  @override
  void initState() {
    super.initState();
    
    _animationController = AnimationController(
      duration: const Duration(seconds: 2),
      vsync: this,
    );
    
    _scaleAnimation = Tween<double>(
      begin: 0.8,
      end: 1.0,
    ).animate(CurvedAnimation(
      parent: _animationController,
      curve: Curves.elasticOut,
    ));
    
    _rotationAnimation = Tween<double>(
      begin: 0.0,
      end: 0.1,
    ).animate(CurvedAnimation(
      parent: _animationController,
      curve: Curves.easeInOut,
    ));
    
    _generateQRData();
    _animationController.forward();
    
    _amountController.addListener(_updateQRData);
    _memoController.addListener(_updateQRData);
  }
  
  @override
  void dispose() {
    _animationController.dispose();
    _amountController.dispose();
    _memoController.dispose();
    super.dispose();
  }
  
  void _generateQRData() {
    final data = {
      'address': widget.address,
      'chain': 'deshchain',
      'token': 'NAMO',
    };
    
    if (_requestedAmount > BigInt.zero) {
      data['amount'] = NAMOToken.formatAmount(_requestedAmount);
    }
    
    if (_memoController.text.isNotEmpty) {
      data['memo'] = _memoController.text;
    }
    
    // Simple URI format for QR code
    final uri = 'deshchain:${widget.address}?'
        '${_requestedAmount > BigInt.zero ? 'amount=${NAMOToken.formatAmount(_requestedAmount)}&' : ''}'
        '${_memoController.text.isNotEmpty ? 'memo=${Uri.encodeComponent(_memoController.text)}&' : ''}'
        'token=NAMO';
    
    setState(() {
      _qrData = uri;
    });
  }
  
  void _updateQRData() {
    final amountText = _amountController.text;
    if (amountText.isNotEmpty) {
      try {
        _requestedAmount = NAMOToken.parseAmount(amountText);
      } catch (e) {
        _requestedAmount = BigInt.zero;
      }
    } else {
      _requestedAmount = BigInt.zero;
    }
    
    _generateQRData();
  }
  
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.grey[50],
      appBar: AppBar(
        elevation: 0,
        backgroundColor: Colors.transparent,
        flexibleSpace: Container(
          decoration: const BoxDecoration(
            gradient: LinearGradient(
              colors: [AppColors.saffron, AppColors.green],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
          ),
        ),
        title: const CulturalGradientText(
          text: 'Receive NAMO',
          style: TextStyle(
            fontSize: 20,
            fontWeight: FontWeight.bold,
            color: Colors.white,
          ),
        ),
        centerTitle: true,
        iconTheme: const IconThemeData(color: Colors.white),
        actions: [
          IconButton(
            onPressed: _shareAddress,
            icon: const Icon(Icons.share),
          ),
        ],
      ),
      body: Column(
        children: [
          // Header
          Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              gradient: const LinearGradient(
                colors: [AppColors.saffron, AppColors.green],
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
              ),
              borderRadius: const BorderRadius.only(
                bottomLeft: Radius.circular(30),
                bottomRight: Radius.circular(30),
              ),
            ),
            child: Column(
              children: [
                // Balance Display
                CompactNAMOBalanceWidget(
                  address: widget.address,
                ).animate()
                    .fadeIn(duration: 600.ms)
                    .scale(begin: 0.8, end: 1.0),
                
                const SizedBox(height: 16),
                
                // Cultural Quote
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: Colors.white.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(
                      color: Colors.white.withOpacity(0.2),
                    ),
                  ),
                  child: Text(
                    NAMOToken.getCulturalQuote(),
                    style: TextStyle(
                      color: Colors.white.withOpacity(0.9),
                      fontSize: 14,
                      fontStyle: FontStyle.italic,
                    ),
                    textAlign: TextAlign.center,
                  ),
                ).animate()
                    .fadeIn(delay: 200.ms)
                    .slideY(begin: 0.3, end: 0),
              ],
            ),
          ),
          
          // Content
          Expanded(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
                  // QR Code
                  Container(
                    padding: const EdgeInsets.all(20),
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(20),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.grey.withOpacity(0.1),
                          blurRadius: 15,
                          offset: const Offset(0, 5),
                        ),
                      ],
                    ),
                    child: Column(
                      children: [
                        // QR Code Title
                        Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            const Icon(
                              Icons.qr_code_2,
                              color: AppColors.saffron,
                              size: 24,
                            ),
                            const SizedBox(width: 8),
                            const Text(
                              'Scan to Send NAMO',
                              style: TextStyle(
                                fontSize: 18,
                                fontWeight: FontWeight.bold,
                                color: Colors.black87,
                              ),
                            ),
                          ],
                        ),
                        
                        const SizedBox(height: 20),
                        
                        // QR Code
                        AnimatedBuilder(
                          animation: _scaleAnimation,
                          builder: (context, child) {
                            return Transform.scale(
                              scale: _scaleAnimation.value,
                              child: Transform.rotate(
                                angle: _rotationAnimation.value,
                                child: Container(
                                  width: 250,
                                  height: 250,
                                  decoration: BoxDecoration(
                                    color: Colors.white,
                                    borderRadius: BorderRadius.circular(20),
                                    border: Border.all(
                                      color: AppColors.saffron,
                                      width: 2,
                                    ),
                                  ),
                                  child: QrImageView(
                                    data: _qrData,
                                    version: QrVersions.auto,
                                    size: 250,
                                    backgroundColor: Colors.white,
                                    foregroundColor: Colors.black,
                                    padding: const EdgeInsets.all(20),
                                    embeddedImage: const AssetImage('assets/images/namo_logo.png'),
                                    embeddedImageStyle: QrEmbeddedImageStyle(
                                      size: const Size(40, 40),
                                    ),
                                  ),
                                ),
                              ),
                            );
                          },
                        ),
                        
                        const SizedBox(height: 20),
                        
                        // Regenerate QR Button
                        ElevatedButton.icon(
                          onPressed: () {
                            _animationController.reset();
                            _animationController.forward();
                            _generateQRData();
                          },
                          style: ElevatedButton.styleFrom(
                            backgroundColor: AppColors.saffron.withOpacity(0.1),
                            foregroundColor: AppColors.saffron,
                            elevation: 0,
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(12),
                            ),
                          ),
                          icon: const Icon(Icons.refresh),
                          label: const Text('Regenerate QR'),
                        ),
                      ],
                    ),
                  ).animate()
                      .fadeIn(delay: 400.ms)
                      .slideY(begin: 0.3, end: 0),
                  
                  const SizedBox(height: 24),
                  
                  // Address Display
                  Container(
                    padding: const EdgeInsets.all(20),
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(20),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.grey.withOpacity(0.1),
                          blurRadius: 15,
                          offset: const Offset(0, 5),
                        ),
                      ],
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Text(
                          'Your NAMO Address',
                          style: TextStyle(
                            fontSize: 18,
                            fontWeight: FontWeight.bold,
                            color: Colors.black87,
                          ),
                        ),
                        const SizedBox(height: 12),
                        
                        Container(
                          padding: const EdgeInsets.all(16),
                          decoration: BoxDecoration(
                            color: Colors.grey.shade50,
                            borderRadius: BorderRadius.circular(12),
                            border: Border.all(
                              color: Colors.grey.shade300,
                            ),
                          ),
                          child: Column(
                            children: [
                              Row(
                                children: [
                                  Expanded(
                                    child: Text(
                                      widget.address,
                                      style: const TextStyle(
                                        fontFamily: 'monospace',
                                        fontSize: 12,
                                        color: Colors.black87,
                                      ),
                                    ),
                                  ),
                                  IconButton(
                                    onPressed: _copyAddress,
                                    icon: const Icon(
                                      Icons.copy,
                                      size: 20,
                                      color: AppColors.saffron,
                                    ),
                                  ),
                                ],
                              ),
                              
                              const SizedBox(height: 8),
                              
                              Row(
                                children: [
                                  Expanded(
                                    child: ElevatedButton.icon(
                                      onPressed: _copyAddress,
                                      style: ElevatedButton.styleFrom(
                                        backgroundColor: AppColors.saffron,
                                        foregroundColor: Colors.white,
                                        shape: RoundedRectangleBorder(
                                          borderRadius: BorderRadius.circular(8),
                                        ),
                                      ),
                                      icon: const Icon(Icons.copy, size: 16),
                                      label: const Text('Copy Address'),
                                    ),
                                  ),
                                  const SizedBox(width: 12),
                                  Expanded(
                                    child: ElevatedButton.icon(
                                      onPressed: _shareAddress,
                                      style: ElevatedButton.styleFrom(
                                        backgroundColor: AppColors.green,
                                        foregroundColor: Colors.white,
                                        shape: RoundedRectangleBorder(
                                          borderRadius: BorderRadius.circular(8),
                                        ),
                                      ),
                                      icon: const Icon(Icons.share, size: 16),
                                      label: const Text('Share'),
                                    ),
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ).animate()
                      .fadeIn(delay: 600.ms)
                      .slideY(begin: 0.3, end: 0),
                  
                  const SizedBox(height: 24),
                  
                  // Custom Amount Section
                  Container(
                    padding: const EdgeInsets.all(20),
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(20),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.grey.withOpacity(0.1),
                          blurRadius: 15,
                          offset: const Offset(0, 5),
                        ),
                      ],
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            const Text(
                              'Request Specific Amount',
                              style: TextStyle(
                                fontSize: 18,
                                fontWeight: FontWeight.bold,
                                color: Colors.black87,
                              ),
                            ),
                            Switch(
                              value: _showCustomAmount,
                              onChanged: (value) {
                                setState(() {
                                  _showCustomAmount = value;
                                  if (!value) {
                                    _amountController.clear();
                                    _memoController.clear();
                                  }
                                });
                              },
                              activeColor: AppColors.saffron,
                            ),
                          ],
                        ),
                        
                        if (_showCustomAmount) ...[
                          const SizedBox(height: 16),
                          
                          // Amount Input
                          TextField(
                            controller: _amountController,
                            keyboardType: const TextInputType.numberWithOptions(decimal: true),
                            decoration: InputDecoration(
                              hintText: 'Enter amount in NAMO',
                              prefixIcon: const Icon(Icons.account_balance),
                              suffixText: 'NAMO',
                              border: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(12),
                              ),
                              focusedBorder: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(12),
                                borderSide: const BorderSide(color: AppColors.saffron),
                              ),
                            ),
                          ),
                          
                          const SizedBox(height: 16),
                          
                          // Memo Input
                          TextField(
                            controller: _memoController,
                            decoration: InputDecoration(
                              hintText: 'Add a note (optional)',
                              prefixIcon: const Icon(Icons.message),
                              border: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(12),
                              ),
                              focusedBorder: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(12),
                                borderSide: const BorderSide(color: AppColors.saffron),
                              ),
                            ),
                          ),
                          
                          const SizedBox(height: 16),
                          
                          // Quick Amount Buttons
                          Row(
                            children: [
                              _buildQuickAmountButton('100', BigInt.from(100000000)),
                              const SizedBox(width: 8),
                              _buildQuickAmountButton('500', BigInt.from(500000000)),
                              const SizedBox(width: 8),
                              _buildQuickAmountButton('1000', BigInt.from(1000000000)),
                            ],
                          ),
                        ],
                      ],
                    ),
                  ).animate()
                      .fadeIn(delay: 800.ms)
                      .slideY(begin: 0.3, end: 0),
                  
                  const SizedBox(height: 24),
                  
                  // Instructions
                  Container(
                    padding: const EdgeInsets.all(20),
                    decoration: BoxDecoration(
                      gradient: const LinearGradient(
                        colors: [AppColors.saffron, AppColors.green],
                        begin: Alignment.topLeft,
                        end: Alignment.bottomRight,
                      ),
                      borderRadius: BorderRadius.circular(20),
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Row(
                          children: [
                            Icon(
                              Icons.info,
                              color: Colors.white,
                              size: 24,
                            ),
                            SizedBox(width: 12),
                            Text(
                              'How to Receive NAMO',
                              style: TextStyle(
                                fontSize: 18,
                                fontWeight: FontWeight.bold,
                                color: Colors.white,
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 16),
                        
                        _buildInstructionStep(
                          '1',
                          'Share your address or QR code with the sender',
                        ),
                        _buildInstructionStep(
                          '2',
                          'Sender scans QR code or enters your address',
                        ),
                        _buildInstructionStep(
                          '3',
                          'Transaction will appear in your wallet once confirmed',
                        ),
                      ],
                    ),
                  ).animate()
                      .fadeIn(delay: 1000.ms)
                      .slideY(begin: 0.3, end: 0),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
  
  Widget _buildQuickAmountButton(String label, BigInt amount) {
    return Expanded(
      child: ElevatedButton(
        onPressed: () {
          _amountController.text = NAMOToken.formatAmount(amount);
        },
        style: ElevatedButton.styleFrom(
          backgroundColor: AppColors.saffron.withOpacity(0.1),
          foregroundColor: AppColors.saffron,
          elevation: 0,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(8),
          ),
        ),
        child: Text(
          '$label NAMO',
          style: const TextStyle(
            fontSize: 12,
            fontWeight: FontWeight.w600,
          ),
        ),
      ),
    );
  }
  
  Widget _buildInstructionStep(String number, String instruction) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            width: 24,
            height: 24,
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.2),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Center(
              child: Text(
                number,
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 12,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              instruction,
              style: TextStyle(
                color: Colors.white.withOpacity(0.9),
                fontSize: 14,
              ),
            ),
          ),
        ],
      ),
    );
  }
  
  void _copyAddress() {
    Clipboard.setData(ClipboardData(text: widget.address));
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Address copied to clipboard'),
        backgroundColor: AppColors.green,
      ),
    );
  }
  
  void _shareAddress() {
    final message = 'Send NAMO to my address:\n${widget.address}';
    Share.share(message);
  }
}