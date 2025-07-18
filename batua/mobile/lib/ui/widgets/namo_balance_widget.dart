import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/tokens/namo_token.dart';
import '../../utils/constants.dart';
import 'cultural_gradient_text.dart';

/// NAMO Balance Display Widget
class NAMOBalanceWidget extends ConsumerStatefulWidget {
  final String address;
  final bool showDetails;
  final VoidCallback? onTap;
  final bool animated;
  
  const NAMOBalanceWidget({
    super.key,
    required this.address,
    this.showDetails = true,
    this.onTap,
    this.animated = false,
  });
  
  @override
  ConsumerState<NAMOBalanceWidget> createState() => _NAMOBalanceWidgetState();
}

class _NAMOBalanceWidgetState extends ConsumerState<NAMOBalanceWidget> {
  BigInt _balance = BigInt.zero;
  bool _isLoading = false;
  NAMOMarketStats? _marketStats;
  
  @override
  void initState() {
    super.initState();
    _loadBalance();
    _loadMarketStats();
  }
  
  Future<void> _loadBalance() async {
    setState(() {
      _isLoading = true;
    });
    
    try {
      final balance = await NAMOToken.getBalance(widget.address);
      setState(() {
        _balance = balance;
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _isLoading = false;
      });
    }
  }
  
  Future<void> _loadMarketStats() async {
    try {
      final stats = await NAMOToken.getMarketStats();
      setState(() {
        _marketStats = stats;
      });
    } catch (e) {
      // Handle error silently
    }
  }
  
  @override
  Widget build(BuildContext context) {
    Widget content = Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            AppColors.saffron,
            AppColors.white,
            AppColors.green,
          ],
        ),
        borderRadius: BorderRadius.circular(20),
        boxShadow: [
          BoxShadow(
            color: AppColors.saffron.withOpacity(0.3),
            blurRadius: 15,
            offset: const Offset(0, 8),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Token Header
          Row(
            children: [
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: Colors.white.withOpacity(0.2),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: const Icon(
                  Icons.account_balance,
                  color: Colors.white,
                  size: 24,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'NAMO Token',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    Text(
                      'Native Token of DeshChain',
                      style: TextStyle(
                        color: Colors.white.withOpacity(0.8),
                        fontSize: 12,
                      ),
                    ),
                  ],
                ),
              ),
              if (widget.showDetails)
                IconButton(
                  onPressed: _loadBalance,
                  icon: Icon(
                    Icons.refresh,
                    color: Colors.white.withOpacity(0.8),
                  ),
                ),
            ],
          ),
          
          const SizedBox(height: 20),
          
          // Balance Display
          if (_isLoading) ...[
            const Center(
              child: CircularProgressIndicator(
                valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
              ),
            ),
          ] else ...[
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Main Balance
                CulturalGradientText(
                  text: '${NAMOToken.formatAmount(_balance)} NAMO',
                  style: const TextStyle(
                    fontSize: 32,
                    fontWeight: FontWeight.bold,
                  ),
                  gradientType: GradientType.flag,
                ),
                
                const SizedBox(height: 8),
                
                // INR Value
                if (_marketStats != null)
                  Text(
                    '≈ ₹${(_balance.toDouble() / 1000000 * _marketStats!.priceINR).toStringAsFixed(2)}',
                    style: TextStyle(
                      color: Colors.white.withOpacity(0.9),
                      fontSize: 18,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                
                const SizedBox(height: 16),
                
                // Cultural Quote
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: Colors.white.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(
                      color: Colors.white.withOpacity(0.2),
                      width: 1,
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
                ),
              ],
            ),
          ],
          
          if (widget.showDetails && _marketStats != null) ...[
            const SizedBox(height: 20),
            
            // Market Stats
            Row(
              children: [
                Expanded(
                  child: _buildStatCard(
                    'Price',
                    '₹${_marketStats!.priceINR.toStringAsFixed(4)}',
                    Icons.currency_rupee,
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _buildStatCard(
                    '24h Change',
                    '${_marketStats!.change24h.toStringAsFixed(2)}%',
                    _marketStats!.changeIcon,
                    color: _marketStats!.changeColor,
                  ),
                ),
              ],
            ),
            
            const SizedBox(height: 12),
            
            Row(
              children: [
                Expanded(
                  child: _buildStatCard(
                    'Market Cap',
                    _marketStats!.formattedMarketCap,
                    Icons.pie_chart,
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _buildStatCard(
                    'Volume 24h',
                    _marketStats!.formattedVolume,
                    Icons.bar_chart,
                  ),
                ),
              ],
            ),
          ],
          
          const SizedBox(height: 20),
          
          // Action Buttons
          Row(
            children: [
              Expanded(
                child: _buildActionButton(
                  'Send',
                  Icons.send,
                  () {
                    // TODO: Navigate to send screen
                  },
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildActionButton(
                  'Receive',
                  Icons.qr_code,
                  () {
                    // TODO: Navigate to receive screen
                  },
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildActionButton(
                  'Stake',
                  Icons.trending_up,
                  () {
                    // TODO: Navigate to staking screen
                  },
                ),
              ),
            ],
          ),
        ],
      ),
    );
    
    if (widget.animated) {
      content = content
          .animate()
          .fadeIn(duration: 600.ms)
          .slideY(begin: 0.3, end: 0, duration: 600.ms);
    }
    
    return GestureDetector(
      onTap: widget.onTap,
      child: content,
    );
  }
  
  Widget _buildStatCard(String label, String value, IconData icon, {Color? color}) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: Colors.white.withOpacity(0.2),
          width: 1,
        ),
      ),
      child: Column(
        children: [
          Icon(
            icon,
            color: color ?? Colors.white.withOpacity(0.8),
            size: 18,
          ),
          const SizedBox(height: 4),
          Text(
            label,
            style: TextStyle(
              color: Colors.white.withOpacity(0.7),
              fontSize: 12,
            ),
          ),
          const SizedBox(height: 2),
          Text(
            value,
            style: TextStyle(
              color: color ?? Colors.white,
              fontSize: 14,
              fontWeight: FontWeight.w600,
            ),
          ),
        ],
      ),
    );
  }
  
  Widget _buildActionButton(String label, IconData icon, VoidCallback onTap) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.15),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: Colors.white.withOpacity(0.3),
          width: 1,
        ),
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(12),
          child: Padding(
            padding: const EdgeInsets.symmetric(vertical: 12),
            child: Column(
              children: [
                Icon(
                  icon,
                  color: Colors.white,
                  size: 20,
                ),
                const SizedBox(height: 4),
                Text(
                  label,
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 12,
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

/// Compact NAMO Balance Widget for headers
class CompactNAMOBalanceWidget extends ConsumerStatefulWidget {
  final String address;
  final VoidCallback? onTap;
  
  const CompactNAMOBalanceWidget({
    super.key,
    required this.address,
    this.onTap,
  });
  
  @override
  ConsumerState<CompactNAMOBalanceWidget> createState() => _CompactNAMOBalanceWidgetState();
}

class _CompactNAMOBalanceWidgetState extends ConsumerState<CompactNAMOBalanceWidget> {
  BigInt _balance = BigInt.zero;
  bool _isLoading = false;
  
  @override
  void initState() {
    super.initState();
    _loadBalance();
  }
  
  Future<void> _loadBalance() async {
    setState(() {
      _isLoading = true;
    });
    
    try {
      final balance = await NAMOToken.getBalance(widget.address);
      setState(() {
        _balance = balance;
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _isLoading = false;
      });
    }
  }
  
  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: widget.onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        decoration: BoxDecoration(
          gradient: const LinearGradient(
            colors: [AppColors.saffron, AppColors.green],
          ),
          borderRadius: BorderRadius.circular(25),
          boxShadow: [
            BoxShadow(
              color: AppColors.saffron.withOpacity(0.3),
              blurRadius: 8,
              offset: const Offset(0, 4),
            ),
          ],
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(
              Icons.account_balance,
              color: Colors.white,
              size: 18,
            ),
            const SizedBox(width: 8),
            if (_isLoading) ...[
              const SizedBox(
                width: 16,
                height: 16,
                child: CircularProgressIndicator(
                  strokeWidth: 2,
                  valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
                ),
              ),
            ] else ...[
              Text(
                NAMOToken.formatAmount(_balance),
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                ),
              ),
              const SizedBox(width: 4),
              const Text(
                'NAMO',
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 12,
                  fontWeight: FontWeight.w500,
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }
}

/// NAMO Balance Card for lists
class NAMOBalanceCard extends StatelessWidget {
  final BigInt balance;
  final VoidCallback? onTap;
  final bool showActions;
  
  const NAMOBalanceCard({
    super.key,
    required this.balance,
    this.onTap,
    this.showActions = true,
  });
  
  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 4,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
      ),
      child: Container(
        decoration: BoxDecoration(
          gradient: const LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [
              AppColors.saffron,
              AppColors.white,
              AppColors.green,
            ],
          ),
          borderRadius: BorderRadius.circular(16),
        ),
        child: Material(
          color: Colors.transparent,
          child: InkWell(
            onTap: onTap,
            borderRadius: BorderRadius.circular(16),
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      const Icon(
                        Icons.account_balance,
                        color: Colors.white,
                        size: 24,
                      ),
                      const SizedBox(width: 12),
                      const Expanded(
                        child: Text(
                          'NAMO Token',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 18,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                      if (showActions)
                        IconButton(
                          onPressed: () {
                            // TODO: Show more options
                          },
                          icon: const Icon(
                            Icons.more_vert,
                            color: Colors.white,
                          ),
                        ),
                    ],
                  ),
                  
                  const SizedBox(height: 16),
                  
                  CulturalGradientText(
                    text: '${NAMOToken.formatAmount(balance)} NAMO',
                    style: const TextStyle(
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                    ),
                    gradientType: GradientType.flag,
                  ),
                  
                  const SizedBox(height: 8),
                  
                  Text(
                    'Native token of DeshChain',
                    style: TextStyle(
                      color: Colors.white.withOpacity(0.8),
                      fontSize: 14,
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}