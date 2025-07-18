import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:lottie/lottie.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../widgets/cultural_gradient_text.dart';
import '../widgets/diya_animation.dart';
import '../../providers/app_providers.dart';
import '../../core/wallet/hd_wallet.dart';
import '../../utils/constants.dart';
import 'onboarding/onboarding_screen.dart';
import 'home/home_screen.dart';

class SplashScreen extends ConsumerStatefulWidget {
  const SplashScreen({super.key});

  @override
  ConsumerState<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends ConsumerState<SplashScreen>
    with TickerProviderStateMixin {
  late AnimationController _logoController;
  late AnimationController _diyaController;
  late AnimationController _quoteController;
  
  bool _showQuote = false;
  String _currentQuote = '';
  
  final List<String> _splashQuotes = [
    'सुरक्षा आपकी, सुविधा हमारी',
    'Your Digital Wallet with Indian Soul',
    'भारत का पहला सांस्कृतिक बटुआ',
    'Powered by DeshChain Technology',
    'डिजिटल सुरक्षा, भारतीय संस्कार',
  ];
  
  @override
  void initState() {
    super.initState();
    
    _logoController = AnimationController(
      duration: const Duration(seconds: 2),
      vsync: this,
    );
    
    _diyaController = AnimationController(
      duration: const Duration(seconds: 3),
      vsync: this,
    );
    
    _quoteController = AnimationController(
      duration: const Duration(milliseconds: 800),
      vsync: this,
    );
    
    _initializeApp();
  }
  
  @override
  void dispose() {
    _logoController.dispose();
    _diyaController.dispose();
    _quoteController.dispose();
    super.dispose();
  }
  
  Future<void> _initializeApp() async {
    // Start animations
    _logoController.forward();
    _diyaController.forward();
    
    // Show quote after logo animation
    await Future.delayed(const Duration(milliseconds: 1000));
    setState(() {
      _showQuote = true;
      _currentQuote = _splashQuotes[0];
    });
    _quoteController.forward();
    
    // Cycle through quotes
    for (int i = 1; i < _splashQuotes.length; i++) {
      await Future.delayed(const Duration(milliseconds: 1500));
      await _quoteController.reverse();
      setState(() {
        _currentQuote = _splashQuotes[i];
      });
      await _quoteController.forward();
    }
    
    // Check if wallet exists
    final hasWallet = await _checkWalletExists();
    
    // Navigate to appropriate screen
    await Future.delayed(const Duration(milliseconds: 1000));
    
    if (mounted) {
      Navigator.of(context).pushReplacement(
        MaterialPageRoute(
          builder: (context) => hasWallet 
              ? const HomeScreen()
              : const OnboardingScreen(),
        ),
      );
    }
  }
  
  Future<bool> _checkWalletExists() async {
    try {
      final wallet = await HDWallet.loadWallet();
      return wallet != null && wallet.isInitialized;
    } catch (e) {
      return false;
    }
  }
  
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      body: SafeArea(
        child: Column(
          children: [
            // Top cultural pattern
            Container(
              height: 100,
              decoration: const BoxDecoration(
                gradient: LinearGradient(
                  colors: [AppColors.saffron, AppColors.white],
                  begin: Alignment.topCenter,
                  end: Alignment.bottomCenter,
                ),
              ),
              child: const Center(
                child: CulturalPattern(
                  type: CulturalPatternType.rangoli,
                  size: 60,
                ),
              ),
            ),
            
            Expanded(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  // Logo animation
                  AnimatedBuilder(
                    animation: _logoController,
                    builder: (context, child) {
                      return Transform.scale(
                        scale: _logoController.value,
                        child: const BatuaLogo(size: 120),
                      );
                    },
                  ),
                  
                  const SizedBox(height: 40),
                  
                  // Diya animation
                  SizedBox(
                    height: 100,
                    child: DiyaAnimation(
                      controller: _diyaController,
                      count: 5,
                    ),
                  ),
                  
                  const SizedBox(height: 40),
                  
                  // Quote animation
                  AnimatedBuilder(
                    animation: _quoteController,
                    builder: (context, child) {
                      return Opacity(
                        opacity: _quoteController.value,
                        child: Transform.translate(
                          offset: Offset(0, 20 * (1 - _quoteController.value)),
                          child: _showQuote
                              ? CulturalGradientText(
                                  text: _currentQuote,
                                  style: const TextStyle(
                                    fontSize: 18,
                                    fontWeight: FontWeight.w500,
                                  ),
                                  textAlign: TextAlign.center,
                                )
                              : const SizedBox(),
                        ),
                      );
                    },
                  ),
                  
                  const SizedBox(height: 60),
                  
                  // Loading indicator
                  const CulturalLoadingIndicator(),
                ],
              ),
            ),
            
            // Bottom cultural pattern
            Container(
              height: 100,
              decoration: const BoxDecoration(
                gradient: LinearGradient(
                  colors: [AppColors.white, AppColors.green],
                  begin: Alignment.topCenter,
                  end: Alignment.bottomCenter,
                ),
              ),
              child: const Center(
                child: CulturalPattern(
                  type: CulturalPatternType.lotus,
                  size: 60,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

/// Batua logo with gradient
class BatuaLogo extends StatelessWidget {
  final double size;
  
  const BatuaLogo({
    super.key,
    required this.size,
  });
  
  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        // Logo gradient text
        ShaderMask(
          shaderCallback: (bounds) => AppColors.batuaGradient.createShader(bounds),
          child: Text(
            'बटुआ',
            style: TextStyle(
              fontSize: size,
              fontWeight: FontWeight.bold,
              color: Colors.white,
              fontFamily: 'Samarkan',
            ),
          ),
        ),
        
        const SizedBox(height: 8),
        
        // English text
        Text(
          'BATUA',
          style: TextStyle(
            fontSize: size * 0.3,
            fontWeight: FontWeight.w600,
            color: AppColors.primaryDark,
            letterSpacing: 4,
          ),
        ),
        
        const SizedBox(height: 4),
        
        // Tagline
        Text(
          'Your Digital Wallet with Indian Soul',
          style: TextStyle(
            fontSize: size * 0.12,
            color: AppColors.primaryDark.withOpacity(0.7),
            fontStyle: FontStyle.italic,
          ),
        ),
      ],
    );
  }
}

/// Cultural pattern widget
class CulturalPattern extends StatelessWidget {
  final CulturalPatternType type;
  final double size;
  final Color? color;
  
  const CulturalPattern({
    super.key,
    required this.type,
    required this.size,
    this.color,
  });
  
  @override
  Widget build(BuildContext context) {
    return CustomPaint(
      size: Size(size, size),
      painter: CulturalPatternPainter(
        type: type,
        color: color ?? AppColors.accent,
      ),
    );
  }
}

/// Cultural pattern types
enum CulturalPatternType { rangoli, lotus, diya, mandala }

/// Cultural pattern painter
class CulturalPatternPainter extends CustomPainter {
  final CulturalPatternType type;
  final Color color;
  
  CulturalPatternPainter({
    required this.type,
    required this.color,
  });
  
  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = color
      ..strokeWidth = 2
      ..style = PaintingStyle.stroke;
    
    final center = Offset(size.width / 2, size.height / 2);
    final radius = size.width / 2;
    
    switch (type) {
      case CulturalPatternType.rangoli:
        _drawRangoli(canvas, center, radius, paint);
        break;
      case CulturalPatternType.lotus:
        _drawLotus(canvas, center, radius, paint);
        break;
      case CulturalPatternType.diya:
        _drawDiya(canvas, center, radius, paint);
        break;
      case CulturalPatternType.mandala:
        _drawMandala(canvas, center, radius, paint);
        break;
    }
  }
  
  void _drawRangoli(Canvas canvas, Offset center, double radius, Paint paint) {
    // Draw concentric squares rotated
    for (int i = 0; i < 3; i++) {
      final size = radius * (1 - i * 0.3);
      final rect = Rect.fromCenter(
        center: center,
        width: size * 2,
        height: size * 2,
      );
      
      canvas.save();
      canvas.translate(center.dx, center.dy);
      canvas.rotate(i * 0.785398); // 45 degrees
      canvas.translate(-center.dx, -center.dy);
      canvas.drawRect(rect, paint);
      canvas.restore();
    }
  }
  
  void _drawLotus(Canvas canvas, Offset center, double radius, Paint paint) {
    // Draw lotus petals
    for (int i = 0; i < 8; i++) {
      final angle = i * 0.785398; // 45 degrees
      final petalPath = Path();
      
      final startX = center.dx + radius * 0.3 * cos(angle);
      final startY = center.dy + radius * 0.3 * sin(angle);
      final endX = center.dx + radius * cos(angle);
      final endY = center.dy + radius * sin(angle);
      
      petalPath.moveTo(startX, startY);
      petalPath.quadraticBezierTo(
        center.dx + radius * 0.8 * cos(angle + 0.3),
        center.dy + radius * 0.8 * sin(angle + 0.3),
        endX,
        endY,
      );
      petalPath.quadraticBezierTo(
        center.dx + radius * 0.8 * cos(angle - 0.3),
        center.dy + radius * 0.8 * sin(angle - 0.3),
        startX,
        startY,
      );
      
      canvas.drawPath(petalPath, paint);
    }
  }
  
  void _drawDiya(Canvas canvas, Offset center, double radius, Paint paint) {
    // Draw diya shape
    final diyaPath = Path();
    
    diyaPath.moveTo(center.dx - radius, center.dy);
    diyaPath.quadraticBezierTo(
      center.dx - radius * 0.5,
      center.dy - radius * 0.5,
      center.dx,
      center.dy,
    );
    diyaPath.quadraticBezierTo(
      center.dx + radius * 0.5,
      center.dy - radius * 0.5,
      center.dx + radius,
      center.dy,
    );
    diyaPath.quadraticBezierTo(
      center.dx + radius * 0.8,
      center.dy + radius * 0.3,
      center.dx,
      center.dy + radius * 0.3,
    );
    diyaPath.quadraticBezierTo(
      center.dx - radius * 0.8,
      center.dy + radius * 0.3,
      center.dx - radius,
      center.dy,
    );
    
    canvas.drawPath(diyaPath, paint);
  }
  
  void _drawMandala(Canvas canvas, Offset center, double radius, Paint paint) {
    // Draw mandala pattern
    canvas.drawCircle(center, radius, paint);
    canvas.drawCircle(center, radius * 0.7, paint);
    canvas.drawCircle(center, radius * 0.4, paint);
    
    // Draw radial lines
    for (int i = 0; i < 12; i++) {
      final angle = i * 0.5236; // 30 degrees
      final startX = center.dx + radius * 0.4 * cos(angle);
      final startY = center.dy + radius * 0.4 * sin(angle);
      final endX = center.dx + radius * cos(angle);
      final endY = center.dy + radius * sin(angle);
      
      canvas.drawLine(Offset(startX, startY), Offset(endX, endY), paint);
    }
  }
  
  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}

/// Cultural loading indicator
class CulturalLoadingIndicator extends StatefulWidget {
  const CulturalLoadingIndicator({super.key});
  
  @override
  State<CulturalLoadingIndicator> createState() => _CulturalLoadingIndicatorState();
}

class _CulturalLoadingIndicatorState extends State<CulturalLoadingIndicator>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  
  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      duration: const Duration(seconds: 2),
      vsync: this,
    )..repeat();
  }
  
  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }
  
  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: _controller,
      builder: (context, child) {
        return Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: List.generate(3, (index) {
            final delay = index * 0.2;
            final animValue = (_controller.value + delay) % 1.0;
            
            return Container(
              margin: const EdgeInsets.symmetric(horizontal: 4),
              width: 8,
              height: 8,
              decoration: BoxDecoration(
                color: AppColors.accent.withOpacity(animValue),
                shape: BoxShape.circle,
              ),
            );
          }),
        );
      },
    );
  }
}

// Import required math functions
import 'dart:math' show cos, sin;