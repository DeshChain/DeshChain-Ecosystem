import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/tokens/dinr_token.dart';
import '../../utils/constants.dart';
import 'cultural_gradient_text.dart';

/// DINR Balance Display Widget
class DINRBalanceWidget extends ConsumerStatefulWidget {
  final String address;
  final bool showDetails;
  final VoidCallback? onTap;
  final bool animated;
  
  const DINRBalanceWidget({
    super.key,
    required this.address,
    this.showDetails = true,
    this.onTap,
    this.animated = false,
  });
  
  @override
  ConsumerState<DINRBalanceWidget> createState() => _DINRBalanceWidgetState();
}

class _DINRBalanceWidgetState extends ConsumerState<DINRBalanceWidget> {
  BigInt _balance = BigInt.zero;
  bool _isLoading = false;
  DINRStabilityMetrics? _stabilityMetrics;
  List<CollateralPosition> _positions = [];
  
  @override
  void initState() {
    super.initState();
    _loadBalance();
    _loadStabilityMetrics();
    _loadCollateralPositions();
  }
  
  Future<void> _loadBalance() async {
    setState(() {
      _isLoading = true;
    });
    
    try {
      final balance = await DINRToken.getBalance(widget.address);
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
  
  Future<void> _loadStabilityMetrics() async {
    try {
      final metrics = await DINRToken.getStabilityMetrics();
      setState(() {
        _stabilityMetrics = metrics;
      });
    } catch (e) {
      // Handle error silently
    }
  }
  
  Future<void> _loadCollateralPositions() async {
    try {
      final positions = await DINRToken.getCollateralPositions(widget.address);
      setState(() {
        _positions = positions;
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
            AppColors.navy,
            AppColors.white,
            AppColors.saffron,
          ],
        ),
        borderRadius: BorderRadius.circular(20),
        boxShadow: [
          BoxShadow(
            color: AppColors.navy.withOpacity(0.3),
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
                  Icons.currency_rupee,
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
                      'DINR Stablecoin',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    Text(
                      'Digital INR • 1:1 Peg',
                      style: TextStyle(
                        color: Colors.white.withOpacity(0.8),
                        fontSize: 12,
                      ),
                    ),
                  ],
                ),
              ),
              if (widget.showDetails)
                Row(
                  children: [
                    // Stability Indicator
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                      decoration: BoxDecoration(
                        color: Colors.green.withOpacity(0.2),
                        borderRadius: BorderRadius.circular(12),
                        border: Border.all(
                          color: Colors.green.withOpacity(0.3),
                          width: 1,
                        ),
                      ),
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Icon(
                            Icons.shield,
                            color: Colors.green.shade300,
                            size: 12,
                          ),
                          const SizedBox(width: 4),
                          Text(
                            'STABLE',
                            style: TextStyle(
                              color: Colors.green.shade300,
                              fontSize: 10,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(width: 8),
                    IconButton(
                      onPressed: () {
                        _loadBalance();
                        _loadStabilityMetrics();
                        _loadCollateralPositions();
                      },
                      icon: Icon(
                        Icons.refresh,
                        color: Colors.white.withOpacity(0.8),
                        size: 20,
                      ),
                    ),
                  ],
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
                  text: '₹${DINRToken.formatAmount(_balance)}',
                  style: const TextStyle(
                    fontSize: 32,
                    fontWeight: FontWeight.bold,
                  ),
                  gradientType: GradientType.flag,
                ),
                
                const SizedBox(height: 4),
                
                // DINR Amount
                Text(
                  '${DINRToken.formatAmount(_balance)} DINR',
                  style: TextStyle(
                    color: Colors.white.withOpacity(0.9),
                    fontSize: 16,
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
                    DINRToken.getCulturalQuote(),
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
          
          if (widget.showDetails && _stabilityMetrics != null) ...[
            const SizedBox(height: 20),
            
            // Stability Metrics
            Row(
              children: [
                Expanded(
                  child: _buildStatCard(
                    'Price',
                    '₹${_stabilityMetrics!.currentPrice.toStringAsFixed(4)}',
                    Icons.currency_rupee,
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _buildStatCard(
                    'APY',
                    '${_stabilityMetrics!.yieldAPY.toStringAsFixed(1)}%',
                    Icons.trending_up,
                    color: Colors.green,
                  ),
                ),
              ],
            ),
            
            const SizedBox(height: 12),
            
            Row(
              children: [
                Expanded(
                  child: _buildStatCard(
                    'Collateral',
                    '${_stabilityMetrics!.collateralRatio.toStringAsFixed(0)}%',
                    Icons.shield,
                    color: _stabilityMetrics!.collateralRatio >= 150 ? Colors.green : Colors.orange,
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _buildStatCard(
                    'Supply',
                    _stabilityMetrics!.formattedTotalSupply,
                    Icons.account_balance_wallet,
                  ),
                ),
              ],
            ),
          ],
          
          // Collateral Positions Summary
          if (_positions.isNotEmpty) ...[
            const SizedBox(width: 16),
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
              child: Row(
                children: [
                  Icon(
                    Icons.account_balance,
                    color: Colors.white.withOpacity(0.8),
                    size: 16,
                  ),
                  const SizedBox(width: 8),
                  Text(
                    '${_positions.length} Collateral Position${_positions.length > 1 ? 's' : ''}',
                    style: TextStyle(
                      color: Colors.white.withOpacity(0.9),
                      fontSize: 14,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                  const Spacer(),
                  Text(
                    '${_positions.where((p) => p.isHealthy).length}/${_positions.length} Healthy',
                    style: TextStyle(
                      color: Colors.white.withOpacity(0.7),
                      fontSize: 12,
                    ),
                  ),
                ],
              ),
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
                    // TODO: Navigate to DINR send screen
                  },
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildActionButton(
                  'Receive',
                  Icons.qr_code,
                  () {
                    // TODO: Navigate to DINR receive screen
                  },
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildActionButton(
                  'Mint',
                  Icons.add_circle,
                  () {
                    // TODO: Navigate to DINR mint screen
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

/// Compact DINR Balance Widget for headers
class CompactDINRBalanceWidget extends ConsumerStatefulWidget {
  final String address;
  final VoidCallback? onTap;
  
  const CompactDINRBalanceWidget({
    super.key,
    required this.address,
    this.onTap,
  });
  
  @override
  ConsumerState<CompactDINRBalanceWidget> createState() => _CompactDINRBalanceWidgetState();
}

class _CompactDINRBalanceWidgetState extends ConsumerState<CompactDINRBalanceWidget> {
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
      final balance = await DINRToken.getBalance(widget.address);
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
            colors: [AppColors.navy, AppColors.saffron],
          ),
          borderRadius: BorderRadius.circular(25),
          boxShadow: [
            BoxShadow(
              color: AppColors.navy.withOpacity(0.3),
              blurRadius: 8,
              offset: const Offset(0, 4),
            ),
          ],
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(
              Icons.currency_rupee,
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
                '₹${DINRToken.formatAmount(_balance)}',
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                ),
              ),
              const SizedBox(width: 4),
              const Text(
                'DINR',
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

/// DINR Balance Card for lists
class DINRBalanceCard extends StatelessWidget {
  final BigInt balance;
  final VoidCallback? onTap;
  final bool showActions;
  final List<CollateralPosition> positions;
  
  const DINRBalanceCard({
    super.key,
    required this.balance,
    this.onTap,
    this.showActions = true,
    this.positions = const [],
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
              AppColors.navy,
              AppColors.white,
              AppColors.saffron,
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
                        Icons.currency_rupee,
                        color: Colors.white,
                        size: 24,
                      ),
                      const SizedBox(width: 12),
                      const Expanded(
                        child: Text(
                          'DINR Stablecoin',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 18,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                        decoration: BoxDecoration(
                          color: Colors.green.withOpacity(0.2),
                          borderRadius: BorderRadius.circular(12),
                          border: Border.all(
                            color: Colors.green.withOpacity(0.3),
                            width: 1,
                          ),
                        ),
                        child: Text(
                          '₹1.00',
                          style: TextStyle(
                            color: Colors.green.shade300,
                            fontSize: 12,
                            fontWeight: FontWeight.w600,
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
                    text: '₹${DINRToken.formatAmount(balance)}',
                    style: const TextStyle(
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                    ),
                    gradientType: GradientType.flag,
                  ),
                  
                  const SizedBox(height: 4),
                  
                  Text(
                    '${DINRToken.formatAmount(balance)} DINR',
                    style: TextStyle(
                      color: Colors.white.withOpacity(0.8),
                      fontSize: 14,
                    ),
                  ),
                  
                  if (positions.isNotEmpty) ...[
                    const SizedBox(height: 12),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                      decoration: BoxDecoration(
                        color: Colors.white.withOpacity(0.1),
                        borderRadius: BorderRadius.circular(16),
                        border: Border.all(
                          color: Colors.white.withOpacity(0.2),
                          width: 1,
                        ),
                      ),
                      child: Text(
                        '${positions.length} Position${positions.length > 1 ? 's' : ''}',
                        style: TextStyle(
                          color: Colors.white.withOpacity(0.9),
                          fontSize: 12,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                    ),
                  ],
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}