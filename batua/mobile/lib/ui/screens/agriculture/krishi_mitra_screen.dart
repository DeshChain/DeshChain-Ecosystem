import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:url_launcher/url_launcher.dart';

import '../../widgets/cultural_gradient_text.dart';
import '../../../utils/constants.dart';
import '../../../utils/logger.dart';

/// Krishi Mitra Coming Soon Screen - Agriculture Loans
class KrishiMitraScreen extends ConsumerStatefulWidget {
  const KrishiMitraScreen({super.key});
  
  @override
  ConsumerState<KrishiMitraScreen> createState() => _KrishiMitraScreenState();
}

class _KrishiMitraScreenState extends ConsumerState<KrishiMitraScreen> 
    with SingleTickerProviderStateMixin {
  late AnimationController _animationController;
  late Animation<double> _fadeAnimation;
  late Animation<Offset> _slideAnimation;
  
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _phoneController = TextEditingController();
  bool _isNotifyMePressed = false;
  
  // Agriculture-related features
  final List<KrishiFeature> _features = [
    KrishiFeature(
      icon: Icons.agriculture,
      title: 'Crop-Specific Loans',
      description: 'Customized loans for different crops and seasons',
      color: AppColors.green,
    ),
    KrishiFeature(
      icon: Icons.percent,
      title: 'Low Interest Rates',
      description: '6-9% interest rates vs 12-18% from banks',
      color: AppColors.saffron,
    ),
    KrishiFeature(
      icon: Icons.shield,
      title: 'Fraud Protection',
      description: 'Triple-layer verification system',
      color: AppColors.green,
    ),
    KrishiFeature(
      icon: Icons.groups,
      title: 'Community Support',
      description: 'Village panchayat verification and peer monitoring',
      color: AppColors.saffron,
    ),
    KrishiFeature(
      icon: Icons.phone,
      title: 'Easy Application',
      description: 'Apply from home with minimal documentation',
      color: AppColors.green,
    ),
    KrishiFeature(
      icon: Icons.trending_up,
      title: 'Quick Approval',
      description: 'Get approved within 24-48 hours',
      color: AppColors.saffron,
    ),
  ];
  
  // Cultural quotes related to agriculture
  final List<String> _agriculturalQuotes = [
    'जय जवान, जय किसान - वो धरती का पुत्र है',
    'Agriculture is the backbone of India',
    'अन्नदाता हमारे देश का गौरव है',
    'Krishi Mitra - किसान का सच्चा साथी',
    'भारत की मिट्टी में सोना है',
    'Farm to Fork - हमारी जिम्मेदारी',
  ];
  
  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      duration: const Duration(seconds: 2),
      vsync: this,
    );
    
    _fadeAnimation = Tween<double>(
      begin: 0.0,
      end: 1.0,
    ).animate(CurvedAnimation(
      parent: _animationController,
      curve: Curves.easeInOut,
    ));
    
    _slideAnimation = Tween<Offset>(
      begin: const Offset(0, 0.5),
      end: Offset.zero,
    ).animate(CurvedAnimation(
      parent: _animationController,
      curve: Curves.easeOutCubic,
    ));
    
    _animationController.forward();
  }
  
  @override
  void dispose() {
    _animationController.dispose();
    _emailController.dispose();
    _phoneController.dispose();
    super.dispose();
  }
  
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.grey[50],
      body: CustomScrollView(
        slivers: [
          // App Bar
          SliverAppBar(
            expandedHeight: 300,
            floating: false,
            pinned: true,
            backgroundColor: AppColors.green,
            flexibleSpace: FlexibleSpaceBar(
              background: Container(
                decoration: const BoxDecoration(
                  gradient: LinearGradient(
                    begin: Alignment.topLeft,
                    end: Alignment.bottomRight,
                    colors: [
                      AppColors.green,
                      AppColors.saffron,
                    ],
                  ),
                ),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const SizedBox(height: 60),
                    
                    // Coming Soon Badge
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 20,
                        vertical: 8,
                      ),
                      decoration: BoxDecoration(
                        color: Colors.white.withOpacity(0.2),
                        borderRadius: BorderRadius.circular(20),
                        border: Border.all(
                          color: Colors.white.withOpacity(0.3),
                          width: 1,
                        ),
                      ),
                      child: const Text(
                        'Coming Soon',
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 14,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ).animate()
                        .scale(duration: 800.ms)
                        .shimmer(duration: 2000.ms),
                    
                    const SizedBox(height: 20),
                    
                    // Logo and Title
                    Container(
                      width: 100,
                      height: 100,
                      decoration: BoxDecoration(
                        color: Colors.white.withOpacity(0.2),
                        borderRadius: BorderRadius.circular(25),
                      ),
                      child: const Icon(
                        Icons.agriculture,
                        size: 50,
                        color: Colors.white,
                      ),
                    ).animate()
                        .scale(duration: 600.ms)
                        .then(delay: 200.ms)
                        .shimmer(duration: 1500.ms),
                    
                    const SizedBox(height: 20),
                    
                    // Title
                    const CulturalGradientText(
                      text: 'Krishi Mitra',
                      style: TextStyle(
                        fontSize: 36,
                        fontWeight: FontWeight.bold,
                        color: Colors.white,
                      ),
                      gradientType: GradientType.flag,
                    ).animate()
                        .fadeIn(delay: 400.ms)
                        .slideY(begin: 0.3, end: 0),
                    
                    const SizedBox(height: 8),
                    
                    // Subtitle
                    const Text(
                      'किसान का सच्चा साथी',
                      style: TextStyle(
                        fontSize: 18,
                        color: Colors.white,
                        fontWeight: FontWeight.w500,
                      ),
                    ).animate()
                        .fadeIn(delay: 600.ms)
                        .slideY(begin: 0.3, end: 0),
                    
                    const SizedBox(height: 4),
                    
                    const Text(
                      'Community-Backed Agricultural Lending Platform',
                      style: TextStyle(
                        fontSize: 14,
                        color: Colors.white,
                        fontWeight: FontWeight.w400,
                      ),
                    ).animate()
                        .fadeIn(delay: 800.ms)
                        .slideY(begin: 0.3, end: 0),
                  ],
                ),
              ),
            ),
          ),
          
          // Content
          SliverToBoxAdapter(
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Launch Info Card
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
                          'Launch Timeline',
                          style: TextStyle(
                            fontSize: 20,
                            fontWeight: FontWeight.bold,
                            color: Colors.black87,
                          ),
                        ),
                        const SizedBox(height: 16),
                        
                        _buildTimelineItem(
                          'Q2 2024',
                          'Beta Testing with Select Villages',
                          Icons.science,
                          true,
                        ),
                        _buildTimelineItem(
                          'Q3 2024',
                          'Pilot Launch in 5 States',
                          Icons.rocket_launch,
                          true,
                        ),
                        _buildTimelineItem(
                          'Q4 2024',
                          'Full Launch Across India',
                          Icons.public,
                          false,
                        ),
                        _buildTimelineItem(
                          'Q1 2025',
                          'Advanced Features & AI Integration',
                          Icons.smart_toy,
                          false,
                        ),
                      ],
                    ),
                  ).animate()
                      .fadeIn(delay: 200.ms)
                      .slideY(begin: 0.3, end: 0),
                  
                  const SizedBox(height: 30),
                  
                  // Features Section
                  const Text(
                    'What\'s Coming?',
                    style: TextStyle(
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                      color: Colors.black87,
                    ),
                  ).animate()
                      .fadeIn(delay: 400.ms)
                      .slideX(begin: -0.3, end: 0),
                  
                  const SizedBox(height: 20),
                  
                  // Features Grid
                  GridView.builder(
                    shrinkWrap: true,
                    physics: const NeverScrollableScrollPhysics(),
                    gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                      crossAxisCount: 2,
                      childAspectRatio: 1.0,
                      crossAxisSpacing: 16,
                      mainAxisSpacing: 16,
                    ),
                    itemCount: _features.length,
                    itemBuilder: (context, index) {
                      final feature = _features[index];
                      return _buildFeatureCard(feature, index);
                    },
                  ),
                  
                  const SizedBox(height: 30),
                  
                  // Impact Stats
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
                        const Text(
                          'Expected Impact',
                          style: TextStyle(
                            fontSize: 20,
                            fontWeight: FontWeight.bold,
                            color: Colors.white,
                          ),
                        ),
                        const SizedBox(height: 20),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                          children: [
                            _buildImpactStat('1,00,000+', 'Farmers'),
                            _buildImpactStat('₹500 Cr', 'Loan Amount'),
                            _buildImpactStat('95%', 'Success Rate'),
                          ],
                        ),
                      ],
                    ),
                  ).animate()
                      .fadeIn(delay: 800.ms)
                      .slideY(begin: 0.3, end: 0),
                  
                  const SizedBox(height: 30),
                  
                  // Cultural Quote
                  Container(
                    padding: const EdgeInsets.all(20),
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(20),
                      border: Border.all(
                        color: AppColors.green.withOpacity(0.3),
                        width: 2,
                      ),
                    ),
                    child: Column(
                      children: [
                        const Icon(
                          Icons.format_quote,
                          size: 40,
                          color: AppColors.green,
                        ),
                        const SizedBox(height: 12),
                        Text(
                          _agriculturalQuotes[DateTime.now().millisecondsSinceEpoch % _agriculturalQuotes.length],
                          style: const TextStyle(
                            fontSize: 18,
                            fontStyle: FontStyle.italic,
                            color: Colors.black87,
                          ),
                          textAlign: TextAlign.center,
                        ),
                      ],
                    ),
                  ).animate()
                      .fadeIn(delay: 1000.ms)
                      .scale(begin: 0.8, end: 1.0),
                  
                  const SizedBox(height: 30),
                  
                  // Notify Me Section
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
                          'Get Notified When We Launch',
                          style: TextStyle(
                            fontSize: 20,
                            fontWeight: FontWeight.bold,
                            color: Colors.black87,
                          ),
                        ),
                        const SizedBox(height: 16),
                        
                        // Email Input
                        TextField(
                          controller: _emailController,
                          decoration: InputDecoration(
                            hintText: 'Enter your email address',
                            prefixIcon: const Icon(Icons.email),
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                            ),
                            focusedBorder: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: const BorderSide(color: AppColors.green),
                            ),
                          ),
                          keyboardType: TextInputType.emailAddress,
                        ),
                        
                        const SizedBox(height: 12),
                        
                        // Phone Input
                        TextField(
                          controller: _phoneController,
                          decoration: InputDecoration(
                            hintText: 'Enter your phone number',
                            prefixIcon: const Icon(Icons.phone),
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                            ),
                            focusedBorder: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: const BorderSide(color: AppColors.green),
                            ),
                          ),
                          keyboardType: TextInputType.phone,
                        ),
                        
                        const SizedBox(height: 20),
                        
                        // Notify Me Button
                        SizedBox(
                          width: double.infinity,
                          child: ElevatedButton(
                            onPressed: _isNotifyMePressed ? null : _registerForNotifications,
                            style: ElevatedButton.styleFrom(
                              backgroundColor: AppColors.green,
                              foregroundColor: Colors.white,
                              padding: const EdgeInsets.all(16),
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(12),
                              ),
                            ),
                            child: _isNotifyMePressed
                                ? const Row(
                                    mainAxisAlignment: MainAxisAlignment.center,
                                    children: [
                                      Icon(Icons.check_circle, size: 20),
                                      SizedBox(width: 8),
                                      Text('Registered Successfully!'),
                                    ],
                                  )
                                : const Text(
                                    'Notify Me When Available',
                                    style: TextStyle(
                                      fontSize: 16,
                                      fontWeight: FontWeight.bold,
                                    ),
                                  ),
                          ),
                        ),
                        
                        const SizedBox(height: 12),
                        
                        Text(
                          'We\'ll send you an email and SMS when Krishi Mitra is available in your area.',
                          style: TextStyle(
                            fontSize: 12,
                            color: Colors.grey[600],
                          ),
                          textAlign: TextAlign.center,
                        ),
                      ],
                    ),
                  ).animate()
                      .fadeIn(delay: 1200.ms)
                      .slideY(begin: 0.3, end: 0),
                  
                  const SizedBox(height: 30),
                  
                  // Social Media & Contact
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
                          'Stay Connected',
                          style: TextStyle(
                            fontSize: 20,
                            fontWeight: FontWeight.bold,
                            color: Colors.black87,
                          ),
                        ),
                        const SizedBox(height: 16),
                        
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                          children: [
                            _buildSocialButton(
                              'Twitter',
                              Icons.alternate_email,
                              () => _launchURL('https://twitter.com/deshchain'),
                            ),
                            _buildSocialButton(
                              'Telegram',
                              Icons.send,
                              () => _launchURL('https://t.me/deshchain'),
                            ),
                            _buildSocialButton(
                              'WhatsApp',
                              Icons.chat,
                              () => _launchURL('https://wa.me/+91XXXXXXXXXX'),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ).animate()
                      .fadeIn(delay: 1400.ms)
                      .slideY(begin: 0.3, end: 0),
                  
                  const SizedBox(height: 20),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
  
  Widget _buildTimelineItem(String date, String title, IconData icon, bool isCompleted) {
    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      child: Row(
        children: [
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: isCompleted ? AppColors.green : Colors.grey.shade300,
              borderRadius: BorderRadius.circular(20),
            ),
            child: Icon(
              isCompleted ? Icons.check : icon,
              color: isCompleted ? Colors.white : Colors.grey.shade600,
              size: 20,
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  date,
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.bold,
                    color: isCompleted ? AppColors.green : Colors.grey.shade600,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  title,
                  style: const TextStyle(
                    fontSize: 16,
                    color: Colors.black87,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
  
  Widget _buildFeatureCard(KrishiFeature feature, int index) {
    return Container(
      padding: const EdgeInsets.all(16),
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
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            width: 60,
            height: 60,
            decoration: BoxDecoration(
              color: feature.color.withOpacity(0.1),
              borderRadius: BorderRadius.circular(15),
            ),
            child: Icon(
              feature.icon,
              size: 30,
              color: feature.color,
            ),
          ),
          const SizedBox(height: 12),
          Text(
            feature.title,
            style: const TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.bold,
              color: Colors.black87,
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 8),
          Text(
            feature.description,
            style: TextStyle(
              fontSize: 12,
              color: Colors.grey[600],
            ),
            textAlign: TextAlign.center,
          ),
        ],
      ),
    ).animate(delay: Duration(milliseconds: index * 200))
        .fadeIn(duration: 600.ms)
        .slideY(begin: 0.3, end: 0);
  }
  
  Widget _buildImpactStat(String value, String label) {
    return Column(
      children: [
        Text(
          value,
          style: const TextStyle(
            color: Colors.white,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: TextStyle(
            color: Colors.white.withOpacity(0.8),
            fontSize: 14,
          ),
        ),
      ],
    );
  }
  
  Widget _buildSocialButton(String label, IconData icon, VoidCallback onTap) {
    return Column(
      children: [
        Container(
          width: 60,
          height: 60,
          decoration: BoxDecoration(
            color: AppColors.green.withOpacity(0.1),
            borderRadius: BorderRadius.circular(15),
          ),
          child: Material(
            color: Colors.transparent,
            child: InkWell(
              onTap: onTap,
              borderRadius: BorderRadius.circular(15),
              child: Icon(
                icon,
                size: 30,
                color: AppColors.green,
              ),
            ),
          ),
        ),
        const SizedBox(height: 8),
        Text(
          label,
          style: const TextStyle(
            fontSize: 12,
            color: Colors.black87,
          ),
        ),
      ],
    );
  }
  
  Future<void> _registerForNotifications() async {
    if (_emailController.text.isEmpty || _phoneController.text.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Please enter both email and phone number'),
          backgroundColor: Colors.red,
        ),
      );
      return;
    }
    
    setState(() {
      _isNotifyMePressed = true;
    });
    
    // TODO: Implement actual notification registration
    AppLogger.info('Registered for notifications: ${_emailController.text}, ${_phoneController.text}');
    
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Successfully registered! We\'ll notify you when Krishi Mitra launches.'),
        backgroundColor: AppColors.green,
      ),
    );
  }
  
  Future<void> _launchURL(String url) async {
    final Uri uri = Uri.parse(url);
    if (await canLaunchUrl(uri)) {
      await launchUrl(uri);
    } else {
      AppLogger.error('Could not launch URL: $url');
    }
  }
}

/// Krishi Feature Model
class KrishiFeature {
  final IconData icon;
  final String title;
  final String description;
  final Color color;
  
  KrishiFeature({
    required this.icon,
    required this.title,
    required this.description,
    required this.color,
  });
}