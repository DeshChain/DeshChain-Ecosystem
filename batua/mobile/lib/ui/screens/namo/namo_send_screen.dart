import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:qr_code_scanner/qr_code_scanner.dart';

import '../../widgets/cultural_gradient_text.dart';
import '../../widgets/namo_balance_widget.dart';
import '../../../core/tokens/namo_token.dart';
import '../../../core/wallet/hd_wallet.dart';
import '../../../utils/constants.dart';
import '../../../utils/logger.dart';

/// NAMO Send Screen
class NAMOSendScreen extends ConsumerStatefulWidget {
  final String fromAddress;
  final HDWallet wallet;
  final int accountIndex;
  
  const NAMOSendScreen({
    super.key,
    required this.fromAddress,
    required this.wallet,
    required this.accountIndex,
  });
  
  @override
  ConsumerState<NAMOSendScreen> createState() => _NAMOSendScreenState();
}

class _NAMOSendScreenState extends ConsumerState<NAMOSendScreen> {
  final TextEditingController _addressController = TextEditingController();
  final TextEditingController _amountController = TextEditingController();
  final TextEditingController _memoController = TextEditingController();
  final FocusNode _addressFocusNode = FocusNode();
  final FocusNode _amountFocusNode = FocusNode();
  final FocusNode _memoFocusNode = FocusNode();
  
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();
  
  BigInt _balance = BigInt.zero;
  BigInt _selectedAmount = BigInt.zero;
  BigInt _transactionFee = BigInt.zero;
  bool _isLoading = false;
  bool _isValidAddress = false;
  QRViewController? _qrController;
  
  @override
  void initState() {
    super.initState();
    _loadBalance();
    _loadTransactionFee();
    _addressController.addListener(_validateAddress);
    _amountController.addListener(_validateAmount);
  }
  
  @override
  void dispose() {
    _addressController.dispose();
    _amountController.dispose();
    _memoController.dispose();
    _addressFocusNode.dispose();
    _amountFocusNode.dispose();
    _memoFocusNode.dispose();
    _qrController?.dispose();
    super.dispose();
  }
  
  Future<void> _loadBalance() async {
    try {
      final balance = await NAMOToken.getBalance(widget.fromAddress);
      setState(() {
        _balance = balance;
      });
    } catch (e) {
      AppLogger.error('Error loading balance: $e');
    }
  }
  
  Future<void> _loadTransactionFee() async {
    setState(() {
      _transactionFee = NAMOToken.getTransactionFee();
    });
  }
  
  void _validateAddress() {
    final address = _addressController.text;
    setState(() {
      _isValidAddress = NAMOToken.isValidAddress(address);
    });
  }
  
  void _validateAmount() {
    final amountText = _amountController.text;
    if (amountText.isNotEmpty) {
      try {
        final amount = NAMOToken.parseAmount(amountText);
        setState(() {
          _selectedAmount = amount;
        });
      } catch (e) {
        setState(() {
          _selectedAmount = BigInt.zero;
        });
      }
    } else {
      setState(() {
        _selectedAmount = BigInt.zero;
      });
    }
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
          text: 'Send NAMO',
          style: TextStyle(
            fontSize: 20,
            fontWeight: FontWeight.bold,
            color: Colors.white,
          ),
        ),
        centerTitle: true,
        iconTheme: const IconThemeData(color: Colors.white),
      ),
      body: Column(
        children: [
          // Header with balance
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
                  address: widget.fromAddress,
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
          
          // Send Form
          Expanded(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(20),
              child: Form(
                key: _formKey,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Recipient Address
                    _buildSectionTitle('Recipient Address'),
                    const SizedBox(height: 12),
                    
                    TextFormField(
                      controller: _addressController,
                      focusNode: _addressFocusNode,
                      decoration: InputDecoration(
                        hintText: 'Enter DeshChain address (desh...)',
                        prefixIcon: const Icon(Icons.person),
                        suffixIcon: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            if (_isValidAddress)
                              const Icon(
                                Icons.check_circle,
                                color: AppColors.green,
                              ),
                            IconButton(
                              onPressed: _scanQRCode,
                              icon: const Icon(Icons.qr_code_scanner),
                            ),
                          ],
                        ),
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                        focusedBorder: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(12),
                          borderSide: const BorderSide(color: AppColors.saffron),
                        ),
                        errorBorder: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(12),
                          borderSide: const BorderSide(color: Colors.red),
                        ),
                      ),
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'Please enter recipient address';
                        }
                        if (!NAMOToken.isValidAddress(value)) {
                          return 'Please enter a valid DeshChain address';
                        }
                        return null;
                      },
                    ).animate()
                        .fadeIn(delay: 300.ms)
                        .slideX(begin: -0.3, end: 0),
                    
                    const SizedBox(height: 24),
                    
                    // Amount
                    _buildSectionTitle('Amount'),
                    const SizedBox(height: 12),
                    
                    TextFormField(
                      controller: _amountController,
                      focusNode: _amountFocusNode,
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
                        errorBorder: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(12),
                          borderSide: const BorderSide(color: Colors.red),
                        ),
                      ),
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'Please enter amount';
                        }
                        
                        try {
                          final amount = NAMOToken.parseAmount(value);
                          if (amount <= BigInt.zero) {
                            return 'Amount must be greater than 0';
                          }
                          if (amount < NAMOToken.getMinimumTransferAmount()) {
                            return 'Amount must be at least ${NAMOToken.formatAmount(NAMOToken.getMinimumTransferAmount())} NAMO';
                          }
                          if (amount + _transactionFee > _balance) {
                            return 'Insufficient balance';
                          }
                          return null;
                        } catch (e) {
                          return 'Please enter a valid amount';
                        }
                      },
                    ).animate()
                        .fadeIn(delay: 400.ms)
                        .slideX(begin: -0.3, end: 0),
                    
                    const SizedBox(height: 12),
                    
                    // Quick Amount Buttons
                    Row(
                      children: [
                        _buildQuickAmountButton('25%', _balance * BigInt.from(25) ~/ BigInt.from(100)),
                        const SizedBox(width: 8),
                        _buildQuickAmountButton('50%', _balance * BigInt.from(50) ~/ BigInt.from(100)),
                        const SizedBox(width: 8),
                        _buildQuickAmountButton('75%', _balance * BigInt.from(75) ~/ BigInt.from(100)),
                        const SizedBox(width: 8),
                        _buildQuickAmountButton('Max', _balance - _transactionFee),
                      ],
                    ).animate()
                        .fadeIn(delay: 500.ms)
                        .slideX(begin: -0.3, end: 0),
                    
                    const SizedBox(height: 24),
                    
                    // Memo
                    _buildSectionTitle('Memo (Optional)'),
                    const SizedBox(height: 12),
                    
                    TextFormField(
                      controller: _memoController,
                      focusNode: _memoFocusNode,
                      maxLines: 3,
                      decoration: InputDecoration(
                        hintText: 'Add a message (optional)',
                        prefixIcon: const Icon(Icons.message),
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                        focusedBorder: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(12),
                          borderSide: const BorderSide(color: AppColors.saffron),
                        ),
                      ),
                    ).animate()
                        .fadeIn(delay: 600.ms)
                        .slideX(begin: -0.3, end: 0),
                    
                    const SizedBox(height: 24),
                    
                    // Transaction Summary
                    Container(
                      padding: const EdgeInsets.all(16),
                      decoration: BoxDecoration(
                        color: Colors.white,
                        borderRadius: BorderRadius.circular(12),
                        border: Border.all(
                          color: Colors.grey.shade300,
                        ),
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(
                            'Transaction Summary',
                            style: TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 12),
                          
                          _buildSummaryRow('Amount', '${NAMOToken.formatAmount(_selectedAmount)} NAMO'),
                          _buildSummaryRow('Transaction Fee', '${NAMOToken.formatAmount(_transactionFee)} NAMO'),
                          const Divider(),
                          _buildSummaryRow(
                            'Total',
                            '${NAMOToken.formatAmount(_selectedAmount + _transactionFee)} NAMO',
                            isTotal: true,
                          ),
                        ],
                      ),
                    ).animate()
                        .fadeIn(delay: 700.ms)
                        .slideY(begin: 0.3, end: 0),
                    
                    const SizedBox(height: 32),
                    
                    // Send Button
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(
                        onPressed: _isLoading ? null : _sendNAMO,
                        style: ElevatedButton.styleFrom(
                          backgroundColor: AppColors.saffron,
                          foregroundColor: Colors.white,
                          padding: const EdgeInsets.all(18),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                          ),
                        ),
                        child: _isLoading
                            ? const Row(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  SizedBox(
                                    width: 20,
                                    height: 20,
                                    child: CircularProgressIndicator(
                                      strokeWidth: 2,
                                      valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
                                    ),
                                  ),
                                  SizedBox(width: 12),
                                  Text('Sending...'),
                                ],
                              )
                            : const Text(
                                'Send NAMO',
                                style: TextStyle(
                                  fontSize: 18,
                                  fontWeight: FontWeight.bold,
                                ),
                              ),
                      ),
                    ).animate()
                        .fadeIn(delay: 800.ms)
                        .slideY(begin: 0.3, end: 0),
                  ],
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
  
  Widget _buildSectionTitle(String title) {
    return Text(
      title,
      style: const TextStyle(
        fontSize: 16,
        fontWeight: FontWeight.bold,
        color: Colors.black87,
      ),
    );
  }
  
  Widget _buildQuickAmountButton(String label, BigInt amount) {
    return Expanded(
      child: ElevatedButton(
        onPressed: () {
          if (amount > BigInt.zero && amount <= _balance) {
            _amountController.text = NAMOToken.formatAmount(amount);
          }
        },
        style: ElevatedButton.styleFrom(
          backgroundColor: Colors.white,
          foregroundColor: AppColors.saffron,
          padding: const EdgeInsets.symmetric(vertical: 8),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(8),
            side: const BorderSide(color: AppColors.saffron),
          ),
        ),
        child: Text(
          label,
          style: const TextStyle(
            fontSize: 12,
            fontWeight: FontWeight.w600,
          ),
        ),
      ),
    );
  }
  
  Widget _buildSummaryRow(String label, String value, {bool isTotal = false}) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: TextStyle(
              fontSize: isTotal ? 16 : 14,
              fontWeight: isTotal ? FontWeight.bold : FontWeight.normal,
              color: Colors.black87,
            ),
          ),
          Text(
            value,
            style: TextStyle(
              fontSize: isTotal ? 16 : 14,
              fontWeight: isTotal ? FontWeight.bold : FontWeight.normal,
              color: isTotal ? AppColors.saffron : Colors.black87,
            ),
          ),
        ],
      ),
    );
  }
  
  Future<void> _scanQRCode() async {
    // TODO: Implement QR code scanning
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('QR Code Scanner'),
        content: const Text('QR code scanning will be implemented soon.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }
  
  Future<void> _sendNAMO() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }
    
    setState(() {
      _isLoading = true;
    });
    
    try {
      // Show confirmation dialog
      final confirmed = await _showConfirmationDialog();
      if (!confirmed) {
        setState(() {
          _isLoading = false;
        });
        return;
      }
      
      // Send transaction
      final txHash = await NAMOToken.sendTokens(
        fromAddress: widget.fromAddress,
        toAddress: _addressController.text,
        amount: _selectedAmount,
        memo: _memoController.text.isEmpty ? NAMOToken.getCulturalQuote() : _memoController.text,
        wallet: widget.wallet,
        accountIndex: widget.accountIndex,
      );
      
      // Show success dialog
      await _showSuccessDialog(txHash);
      
      // Navigate back
      Navigator.pop(context);
    } catch (e) {
      AppLogger.error('Error sending NAMO: $e');
      
      // Show error dialog
      showDialog(
        context: context,
        builder: (context) => AlertDialog(
          title: const Text('Transaction Failed'),
          content: Text('Failed to send NAMO: ${e.toString()}'),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('OK'),
            ),
          ],
        ),
      );
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }
  
  Future<bool> _showConfirmationDialog() async {
    return await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Confirm Transaction'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('To: ${_addressController.text}'),
            const SizedBox(height: 8),
            Text('Amount: ${NAMOToken.formatAmount(_selectedAmount)} NAMO'),
            const SizedBox(height: 8),
            Text('Fee: ${NAMOToken.formatAmount(_transactionFee)} NAMO'),
            const SizedBox(height: 8),
            Text('Total: ${NAMOToken.formatAmount(_selectedAmount + _transactionFee)} NAMO'),
            if (_memoController.text.isNotEmpty) ...[
              const SizedBox(height: 8),
              Text('Memo: ${_memoController.text}'),
            ],
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.pop(context, true),
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.saffron,
              foregroundColor: Colors.white,
            ),
            child: const Text('Confirm'),
          ),
        ],
      ),
    ) ?? false;
  }
  
  Future<void> _showSuccessDialog(String txHash) async {
    await showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Row(
          children: [
            Icon(Icons.check_circle, color: AppColors.green),
            SizedBox(width: 12),
            Text('Transaction Sent'),
          ],
        ),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text('Your NAMO transaction has been sent successfully!'),
            const SizedBox(height: 12),
            Text('Transaction Hash:'),
            const SizedBox(height: 4),
            GestureDetector(
              onTap: () {
                Clipboard.setData(ClipboardData(text: txHash));
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('Transaction hash copied to clipboard')),
                );
              },
              child: Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: Colors.grey.shade100,
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  txHash,
                  style: const TextStyle(
                    fontFamily: 'monospace',
                    fontSize: 12,
                  ),
                ),
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }
}