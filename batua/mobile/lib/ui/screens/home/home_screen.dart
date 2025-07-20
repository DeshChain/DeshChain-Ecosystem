import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../widgets/cultural_gradient_text.dart';
import '../../widgets/namo_balance_widget.dart';
import '../namo/namo_send_screen.dart';
import '../namo/namo_receive_screen.dart';
import '../suraksha/suraksha_scheme_screen.dart';
import '../agriculture/krishi_mitra_screen.dart';
import '../../../core/wallet/hd_wallet.dart';
import '../../../core/tokens/namo_token.dart';
import '../../../utils/constants.dart';
import '../../../utils/logger.dart';

/// Main Home Screen - Batua Wallet Dashboard
class HomeScreen extends ConsumerStatefulWidget {
  const HomeScreen({super.key});
  
  @override
  ConsumerState<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends ConsumerState<HomeScreen> 
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  HDWallet? _wallet;
  String _currentAddress = '';
  bool _isLoading = false;
  
  // Cultural quotes for home screen
  final List<String> _homeQuotes = [
    'आपका डिजिटल बटुआ, भारतीय संस्कारों के साथ',
    'Your Digital Wallet with Indian Soul',
    'सुरक्षा और सुविधा का संगम',
    'Where Technology Meets Tradition',
    'भविष्य का बटुआ, भारतीय जड़ों के साथ',
  ];
  
  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
    _loadWallet();
  }
  
  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }
  
  Future<void> _loadWallet() async {
    setState(() {
      _isLoading = true;
    });
    
    try {
      final wallet = await HDWallet.loadWallet();
      if (wallet != null) {
        final address = wallet.getDeshChainAddress(0);
        setState(() {
          _wallet = wallet;
          _currentAddress = address;
          _isLoading = false;
        });
      } else {
        // Navigate to onboarding if no wallet exists
        Navigator.pushReplacementNamed(context, '/onboarding');
      }
    } catch (e) {
      AppLogger.error('Error loading wallet: $e');
      setState(() {
        _isLoading = false;
      });
    }
  }
  
  @override
  Widget build(BuildContext context) {
    if (_isLoading || _wallet == null) {
      return const Scaffold(
        body: Center(
          child: CircularProgressIndicator(),
        ),
      );
    }
    
    return Scaffold(
      backgroundColor: Colors.grey[50],
      body: Column(
        children: [
          // Header
          Container(
            padding: const EdgeInsets.fromLTRB(20, 60, 20, 20),
            decoration: const BoxDecoration(
              gradient: LinearGradient(
                colors: [AppColors.saffron, AppColors.green],
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
              ),
              borderRadius: BorderRadius.only(
                bottomLeft: Radius.circular(30),
                bottomRight: Radius.circular(30),
              ),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Header Row
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Text(
                          'नमस्ते',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 24,
                            fontWeight: FontWeight.bold,
                          ),
                        ).animate()
                            .fadeIn(duration: 600.ms)
                            .slideX(begin: -0.3, end: 0),
                        const SizedBox(height: 4),
                        Text(
                          'Welcome to Batua',
                          style: TextStyle(
                            color: Colors.white.withOpacity(0.9),
                            fontSize: 16,
                          ),
                        ).animate()
                            .fadeIn(delay: 200.ms)
                            .slideX(begin: -0.3, end: 0),
                      ],
                    ),
                    Row(
                      children: [
                        // Notifications
                        IconButton(
                          onPressed: () {
                            // TODO: Show notifications
                          },
                          icon: const Icon(
                            Icons.notifications,
                            color: Colors.white,
                          ),
                        ),
                        // Settings
                        IconButton(
                          onPressed: () {
                            // TODO: Show settings
                          },
                          icon: const Icon(
                            Icons.settings,
                            color: Colors.white,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
                
                const SizedBox(height: 20),
                
                // Cultural Quote of the Day
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: Colors.white.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(
                      color: Colors.white.withOpacity(0.2),
                    ),
                  ),
                  child: Text(
                    _homeQuotes[DateTime.now().day % _homeQuotes.length],
                    style: TextStyle(
                      color: Colors.white.withOpacity(0.9),
                      fontSize: 16,
                      fontStyle: FontStyle.italic,
                    ),
                    textAlign: TextAlign.center,
                  ),
                ).animate()
                    .fadeIn(delay: 400.ms)
                    .slideY(begin: 0.3, end: 0),
              ],
            ),
          ),
          
          // Tab Bar
          Container(
            margin: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(15),
              boxShadow: [
                BoxShadow(
                  color: Colors.grey.withOpacity(0.1),
                  blurRadius: 10,
                  offset: const Offset(0, 5),
                ),
              ],
            ),
            child: TabBar(
              controller: _tabController,
              indicator: BoxDecoration(
                gradient: const LinearGradient(
                  colors: [AppColors.saffron, AppColors.green],
                ),
                borderRadius: BorderRadius.circular(15),
              ),
              labelColor: Colors.white,
              unselectedLabelColor: Colors.grey[600],
              tabs: const [
                Tab(text: 'Wallet'),
                Tab(text: 'Pension'),
                Tab(text: 'Krishi'),
                Tab(text: 'More'),
              ],
            ),
          ).animate()
              .fadeIn(delay: 600.ms)
              .slideY(begin: 0.3, end: 0),
          
          // Tab Content
          Expanded(
            child: TabBarView(
              controller: _tabController,
              children: [
                _buildWalletTab(),
                _buildPensionTab(),
                _buildKrishiTab(),
                _buildMoreTab(),
              ],
            ),
          ),
        ],
      ),
    );
  }
  
  Widget _buildWalletTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // NAMO Balance Widget
          NAMOBalanceWidget(
            address: _currentAddress,
            animated: true,
            onTap: () {
              // TODO: Navigate to NAMO details
            },
          ).animate()
              .fadeIn(delay: 200.ms)
              .slideY(begin: 0.3, end: 0),
          
          const SizedBox(height: 24),
          
          // Quick Actions
          const Text(
            'Quick Actions',
            style: TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: Colors.black87,
            ),
          ).animate()
              .fadeIn(delay: 400.ms)
              .slideX(begin: -0.3, end: 0),
          
          const SizedBox(height: 16),
          
          Row(
            children: [
              Expanded(
                child: _buildQuickActionCard(
                  'Send',
                  Icons.send,
                  AppColors.saffron,
                  () => _navigateToSend(),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildQuickActionCard(
                  'Receive',
                  Icons.qr_code,
                  AppColors.green,
                  () => _navigateToReceive(),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildQuickActionCard(
                  'Stake',
                  Icons.trending_up,
                  AppColors.saffron,
                  () {
                    // TODO: Navigate to staking
                  },
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _buildQuickActionCard(
                  'History',
                  Icons.history,
                  AppColors.green,
                  () {
                    // TODO: Navigate to transaction history
                  },
                ),
              ),
            ],
          ).animate()
              .fadeIn(delay: 600.ms)
              .slideY(begin: 0.3, end: 0),
          
          const SizedBox(height: 24),
          
          // Recent Transactions
          const Text(
            'Recent Transactions',
            style: TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: Colors.black87,
            ),
          ).animate()
              .fadeIn(delay: 800.ms)
              .slideX(begin: -0.3, end: 0),
          
          const SizedBox(height: 16),
          
          // Transaction List
          FutureBuilder<List<NAMOTransaction>>(
            future: NAMOToken.getTransactionHistory(_currentAddress, limit: 5),
            builder: (context, snapshot) {
              if (snapshot.connectionState == ConnectionState.waiting) {
                return const Center(
                  child: CircularProgressIndicator(),
                );
              }
              
              if (snapshot.hasError || !snapshot.hasData || snapshot.data!.isEmpty) {
                return Container(
                  padding: const EdgeInsets.all(40),
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(16),
                    boxShadow: [
                      BoxShadow(
                        color: Colors.grey.withOpacity(0.1),
                        blurRadius: 10,
                        offset: const Offset(0, 5),
                      ),
                    ],
                  ),
                  child: Column(
                    children: [
                      Icon(
                        Icons.history,
                        size: 60,
                        color: Colors.grey[400],
                      ),
                      const SizedBox(height: 16),
                      Text(
                        'No transactions yet',
                        style: TextStyle(
                          fontSize: 16,
                          color: Colors.grey[600],
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        'Your transaction history will appear here',
                        style: TextStyle(
                          fontSize: 14,
                          color: Colors.grey[500],
                        ),
                        textAlign: TextAlign.center,
                      ),
                    ],
                  ),
                );
              }
              
              return Column(
                children: snapshot.data!.map((tx) => _buildTransactionCard(tx)).toList(),
              );
            },
          ).animate()
              .fadeIn(delay: 1000.ms)
              .slideY(begin: 0.3, end: 0),
        ],
      ),
    );
  }
  
  Widget _buildPensionTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Pension Overview Card
          Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              gradient: const LinearGradient(
                colors: [AppColors.saffron, AppColors.green],
              ),
              borderRadius: BorderRadius.circular(20),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Row(
                  children: [
                    Icon(
                      Icons.savings,
                      color: Colors.white,
                      size: 30,
                    ),
                    SizedBox(width: 12),
                    Text(
                      'Gram Suraksha Pool',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 22,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                Text(
                  'Secure your future with guaranteed 50% returns',
                  style: TextStyle(
                    color: Colors.white.withOpacity(0.9),
                    fontSize: 16,
                  ),
                ),
                const SizedBox(height: 20),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                  children: [
                    _buildPensionStat('50%', 'Guaranteed Returns'),
                    _buildPensionStat('₹500Cr', 'Total Invested'),
                    _buildPensionStat('1L+', 'Investors'),
                  ],
                ),
              ],
            ),
          ).animate()
              .fadeIn(delay: 200.ms)
              .slideY(begin: 0.3, end: 0),
          
          const SizedBox(height: 24),
          
          // Pension Actions
          Row(
            children: [
              Expanded(
                child: ElevatedButton.icon(
                  onPressed: () {
                    Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (context) => const SurakshaSchemeScreen(),
                      ),
                    );
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: AppColors.saffron,
                    foregroundColor: Colors.white,
                    padding: const EdgeInsets.all(16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  icon: const Icon(Icons.launch),
                  label: const Text('Open Pension'),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: ElevatedButton.icon(
                  onPressed: () {
                    // TODO: Navigate to suraksha calculator
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: AppColors.green,
                    foregroundColor: Colors.white,
                    padding: const EdgeInsets.all(16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  icon: const Icon(Icons.calculate),
                  label: const Text('Calculate'),
                ),
              ),
            ],
          ).animate()
              .fadeIn(delay: 400.ms)
              .slideY(begin: 0.3, end: 0),
          
          const SizedBox(height: 24),
          
          // Pension Features
          _buildFeatureList([
            'Guaranteed 50% returns',
            'KYC verified platform',
            'Cultural bonus rewards',
            'Monthly payout options',
            'Patriotism scoring benefits',
          ]).animate()
              .fadeIn(delay: 600.ms)
              .slideY(begin: 0.3, end: 0),
        ],
      ),
    );
  }
  
  Widget _buildKrishiTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Krishi Mitra Coming Soon Card
          Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              gradient: const LinearGradient(
                colors: [AppColors.green, AppColors.saffron],
              ),
              borderRadius: BorderRadius.circular(20),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    const Icon(
                      Icons.agriculture,
                      color: Colors.white,
                      size: 30,
                    ),
                    const SizedBox(width: 12),
                    const Expanded(
                      child: Text(
                        'Krishi Mitra',
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 22,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 12,
                        vertical: 6,
                      ),
                      decoration: BoxDecoration(
                        color: Colors.white.withOpacity(0.2),
                        borderRadius: BorderRadius.circular(20),
                      ),
                      child: const Text(
                        'Coming Soon',
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 12,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                Text(
                  'Community-backed agricultural lending platform',
                  style: TextStyle(
                    color: Colors.white.withOpacity(0.9),
                    fontSize: 16,
                  ),
                ),
                const SizedBox(height: 20),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                  children: [
                    _buildKrishiStat('6-9%', 'Interest Rate'),
                    _buildKrishiStat('₹500Cr', 'Target Amount'),
                    _buildKrishiStat('1L+', 'Farmers'),
                  ],
                ),
              ],
            ),
          ).animate()
              .fadeIn(delay: 200.ms)
              .slideY(begin: 0.3, end: 0),
          
          const SizedBox(height: 24),
          
          // Krishi Actions
          Row(
            children: [
              Expanded(
                child: ElevatedButton.icon(
                  onPressed: () {
                    Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (context) => const KrishiMitraScreen(),
                      ),
                    );
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: AppColors.green,
                    foregroundColor: Colors.white,
                    padding: const EdgeInsets.all(16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  icon: const Icon(Icons.launch),
                  label: const Text('Learn More'),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: ElevatedButton.icon(
                  onPressed: () {
                    // TODO: Navigate to notification signup
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: AppColors.saffron,
                    foregroundColor: Colors.white,
                    padding: const EdgeInsets.all(16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  icon: const Icon(Icons.notifications),
                  label: const Text('Notify Me'),
                ),
              ),
            ],
          ).animate()
              .fadeIn(delay: 400.ms)
              .slideY(begin: 0.3, end: 0),
          
          const SizedBox(height: 24),
          
          // Krishi Features
          _buildFeatureList([
            'Low interest rates (6-9%)',
            'Community verification',
            'Crop-specific loans',
            'Quick approval process',
            'Triple-layer fraud protection',
          ]).animate()
              .fadeIn(delay: 600.ms)
              .slideY(begin: 0.3, end: 0),
        ],
      ),
    );
  }
  
  Widget _buildMoreTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // DeshChain Ecosystem
          const Text(
            'DeshChain Ecosystem',
            style: TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: Colors.black87,
            ),
          ).animate()
              .fadeIn(delay: 200.ms)
              .slideX(begin: -0.3, end: 0),
          
          const SizedBox(height: 16),
          
          // Ecosystem Cards
          _buildEcosystemCard(
            'Sikkebaaz',
            'Desi Memecoin Launchpad',
            Icons.rocket_launch,
            AppColors.saffron,
            'Coming Soon',
          ).animate()
              .fadeIn(delay: 400.ms)
              .slideY(begin: 0.3, end: 0),
          
          const SizedBox(height: 12),
          
          _buildEcosystemCard(
            'Money Order DEX',
            'Culturally-rooted Exchange',
            Icons.currency_exchange,
            AppColors.green,
            'Coming Soon',
          ).animate()
              .fadeIn(delay: 600.ms)
              .slideY(begin: 0.3, end: 0),
          
          const SizedBox(height: 12),
          
          _buildEcosystemCard(
            'Cultural NFTs',
            'Indian Heritage Collection',
            Icons.collections,
            AppColors.saffron,
            'Coming Soon',
          ).animate()
              .fadeIn(delay: 800.ms)
              .slideY(begin: 0.3, end: 0),
          
          const SizedBox(height: 24),
          
          // Wallet Settings
          const Text(
            'Wallet Settings',
            style: TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: Colors.black87,
            ),
          ).animate()
              .fadeIn(delay: 1000.ms)
              .slideX(begin: -0.3, end: 0),
          
          const SizedBox(height: 16),
          
          // Settings Options
          _buildSettingsCard(
            'Security',
            'Biometric & PIN settings',
            Icons.security,
            () {
              // TODO: Navigate to security settings
            },
          ).animate()
              .fadeIn(delay: 1200.ms)
              .slideY(begin: 0.3, end: 0),
          
          const SizedBox(height: 12),
          
          _buildSettingsCard(
            'Backup & Recovery',
            'Manage seed phrase',
            Icons.backup,
            () {
              // TODO: Navigate to backup settings
            },
          ).animate()
              .fadeIn(delay: 1400.ms)
              .slideY(begin: 0.3, end: 0),
          
          const SizedBox(height: 12),
          
          _buildSettingsCard(
            'About',
            'App version & support',
            Icons.info,
            () {
              // TODO: Navigate to about page
            },
          ).animate()
              .fadeIn(delay: 1600.ms)
              .slideY(begin: 0.3, end: 0),
        ],
      ),
    );
  }
  
  Widget _buildQuickActionCard(String title, IconData icon, Color color, VoidCallback onTap) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.1),
            blurRadius: 10,
            offset: const Offset(0, 5),
          ),
        ],
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(16),
          child: Padding(
            padding: const EdgeInsets.all(20),
            child: Column(
              children: [
                Container(
                  width: 50,
                  height: 50,
                  decoration: BoxDecoration(
                    color: color.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(15),
                  ),
                  child: Icon(
                    icon,
                    color: color,
                    size: 24,
                  ),
                ),
                const SizedBox(height: 12),
                Text(
                  title,
                  style: const TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                    color: Colors.black87,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
  
  Widget _buildTransactionCard(NAMOTransaction transaction) {
    final direction = transaction.getDirection(_currentAddress);
    final isReceived = direction == TransactionDirection.received;
    
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.1),
            blurRadius: 10,
            offset: const Offset(0, 5),
          ),
        ],
      ),
      child: Row(
        children: [
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: (isReceived ? AppColors.green : AppColors.saffron).withOpacity(0.1),
              borderRadius: BorderRadius.circular(10),
            ),
            child: Icon(
              isReceived ? Icons.arrow_downward : Icons.arrow_upward,
              color: isReceived ? AppColors.green : AppColors.saffron,
              size: 20,
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  isReceived ? 'Received' : 'Sent',
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  transaction.getDisplayAddress(_currentAddress),
                  style: const TextStyle(
                    fontSize: 12,
                    color: Colors.grey,
                    fontFamily: 'monospace',
                  ),
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
          ),
          Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                '${isReceived ? '+' : '-'}${transaction.formattedAmount} NAMO',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  color: isReceived ? AppColors.green : AppColors.saffron,
                ),
              ),
              const SizedBox(height: 4),
              Text(
                '${transaction.timestamp.day}/${transaction.timestamp.month}',
                style: const TextStyle(
                  fontSize: 12,
                  color: Colors.grey,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
  
  Widget _buildPensionStat(String value, String label) {
    return Column(
      children: [
        Text(
          value,
          style: const TextStyle(
            color: Colors.white,
            fontSize: 18,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: TextStyle(
            color: Colors.white.withOpacity(0.8),
            fontSize: 12,
          ),
        ),
      ],
    );
  }
  
  Widget _buildKrishiStat(String value, String label) {
    return Column(
      children: [
        Text(
          value,
          style: const TextStyle(
            color: Colors.white,
            fontSize: 18,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: TextStyle(
            color: Colors.white.withOpacity(0.8),
            fontSize: 12,
          ),
        ),
      ],
    );
  }
  
  Widget _buildFeatureList(List<String> features) {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.1),
            blurRadius: 10,
            offset: const Offset(0, 5),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            'Key Features',
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.bold,
              color: Colors.black87,
            ),
          ),
          const SizedBox(height: 16),
          ...features.map((feature) => Padding(
            padding: const EdgeInsets.only(bottom: 12),
            child: Row(
              children: [
                const Icon(
                  Icons.check_circle,
                  color: AppColors.green,
                  size: 20,
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Text(
                    feature,
                    style: const TextStyle(
                      fontSize: 14,
                      color: Colors.black87,
                    ),
                  ),
                ),
              ],
            ),
          )),
        ],
      ),
    );
  }
  
  Widget _buildEcosystemCard(String title, String description, IconData icon, Color color, String status) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.1),
            blurRadius: 10,
            offset: const Offset(0, 5),
          ),
        ],
      ),
      child: Row(
        children: [
          Container(
            width: 50,
            height: 50,
            decoration: BoxDecoration(
              color: color.withOpacity(0.1),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Icon(
              icon,
              color: color,
              size: 24,
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                    color: Colors.black87,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  description,
                  style: const TextStyle(
                    fontSize: 14,
                    color: Colors.grey,
                  ),
                ),
              ],
            ),
          ),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
            decoration: BoxDecoration(
              color: color.withOpacity(0.1),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Text(
              status,
              style: TextStyle(
                fontSize: 12,
                color: color,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
        ],
      ),
    );
  }
  
  Widget _buildSettingsCard(String title, String description, IconData icon, VoidCallback onTap) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.1),
            blurRadius: 10,
            offset: const Offset(0, 5),
          ),
        ],
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(12),
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                Container(
                  width: 40,
                  height: 40,
                  decoration: BoxDecoration(
                    color: AppColors.saffron.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(10),
                  ),
                  child: Icon(
                    icon,
                    color: AppColors.saffron,
                    size: 20,
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        title,
                        style: const TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: Colors.black87,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        description,
                        style: const TextStyle(
                          fontSize: 14,
                          color: Colors.grey,
                        ),
                      ),
                    ],
                  ),
                ),
                const Icon(
                  Icons.chevron_right,
                  color: Colors.grey,
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
  
  void _navigateToSend() {
    if (_wallet != null) {
      Navigator.push(
        context,
        MaterialPageRoute(
          builder: (context) => NAMOSendScreen(
            fromAddress: _currentAddress,
            wallet: _wallet!,
            accountIndex: 0,
          ),
        ),
      );
    }
  }
  
  void _navigateToReceive() {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => NAMOReceiveScreen(
          address: _currentAddress,
          displayName: 'My Wallet',
        ),
      ),
    );
  }
}