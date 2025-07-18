import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';

/// Cultural gradient text widget with Indian flag colors
class CulturalGradientText extends StatelessWidget {
  final String text;
  final TextStyle? style;
  final TextAlign? textAlign;
  final GradientType gradientType;
  final bool animated;
  
  const CulturalGradientText({
    super.key,
    required this.text,
    this.style,
    this.textAlign,
    this.gradientType = GradientType.flag,
    this.animated = false,
  });
  
  @override
  Widget build(BuildContext context) {
    Widget textWidget = ShaderMask(
      shaderCallback: (bounds) => _getGradient().createShader(bounds),
      child: Text(
        text,
        style: (style ?? const TextStyle()).copyWith(
          color: Colors.white,
          fontWeight: FontWeight.w600,
        ),
        textAlign: textAlign,
      ),
    );
    
    if (animated) {
      textWidget = textWidget
          .animate()
          .fadeIn(duration: 600.ms)
          .slideY(begin: 0.2, end: 0, duration: 600.ms);
    }
    
    return textWidget;
  }
  
  LinearGradient _getGradient() {
    switch (gradientType) {
      case GradientType.flag:
        return const LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [
            Color(0xFFFF9933), // Saffron
            Color(0xFFFFFFFF), // White
            Color(0xFF138808), // Green
          ],
          stops: [0.0, 0.5, 1.0],
        );
      case GradientType.sunset:
        return const LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            Color(0xFFFF6B35), // Orange
            Color(0xFFFFD700), // Gold
            Color(0xFFFF1493), // Pink
          ],
        );
      case GradientType.ocean:
        return const LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [
            Color(0xFF00CED1), // Dark Turquoise
            Color(0xFF4169E1), // Royal Blue
            Color(0xFF000080), // Navy
          ],
        );
      case GradientType.forest:
        return const LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            Color(0xFF32CD32), // Lime Green
            Color(0xFF228B22), // Forest Green
            Color(0xFF006400), // Dark Green
          ],
        );
      case GradientType.festival:
        return const LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            Color(0xFFFFD700), // Gold
            Color(0xFFFF69B4), // Hot Pink
            Color(0xFF9370DB), // Medium Purple
            Color(0xFF00CED1), // Dark Turquoise
          ],
        );
    }
  }
}

/// Gradient types for different occasions
enum GradientType {
  flag,
  sunset,
  ocean,
  forest,
  festival,
}

/// Animated cultural text with typing effect
class AnimatedCulturalText extends StatefulWidget {
  final String text;
  final TextStyle? style;
  final Duration duration;
  final GradientType gradientType;
  final VoidCallback? onComplete;
  
  const AnimatedCulturalText({
    super.key,
    required this.text,
    this.style,
    this.duration = const Duration(milliseconds: 2000),
    this.gradientType = GradientType.flag,
    this.onComplete,
  });
  
  @override
  State<AnimatedCulturalText> createState() => _AnimatedCulturalTextState();
}

class _AnimatedCulturalTextState extends State<AnimatedCulturalText>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  late Animation<int> _characterCount;
  
  @override
  void initState() {
    super.initState();
    
    _controller = AnimationController(
      duration: widget.duration,
      vsync: this,
    );
    
    _characterCount = IntTween(
      begin: 0,
      end: widget.text.length,
    ).animate(CurvedAnimation(
      parent: _controller,
      curve: Curves.easeInOut,
    ));
    
    _controller.addStatusListener((status) {
      if (status == AnimationStatus.completed) {
        widget.onComplete?.call();
      }
    });
    
    _controller.forward();
  }
  
  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }
  
  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: _characterCount,
      builder: (context, child) {
        final displayText = widget.text.substring(0, _characterCount.value);
        return CulturalGradientText(
          text: displayText,
          style: widget.style,
          gradientType: widget.gradientType,
        );
      },
    );
  }
}

/// Cultural text with glow effect
class GlowingCulturalText extends StatelessWidget {
  final String text;
  final TextStyle? style;
  final Color glowColor;
  final double glowRadius;
  final GradientType gradientType;
  
  const GlowingCulturalText({
    super.key,
    required this.text,
    this.style,
    this.glowColor = const Color(0xFFFFD700),
    this.glowRadius = 10.0,
    this.gradientType = GradientType.flag,
  });
  
  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        // Glow effect
        Text(
          text,
          style: (style ?? const TextStyle()).copyWith(
            color: glowColor,
            shadows: [
              Shadow(
                color: glowColor.withOpacity(0.6),
                blurRadius: glowRadius,
              ),
              Shadow(
                color: glowColor.withOpacity(0.4),
                blurRadius: glowRadius * 2,
              ),
            ],
          ),
        ),
        // Main text
        CulturalGradientText(
          text: text,
          style: style,
          gradientType: gradientType,
        ),
      ],
    );
  }
}

/// Cultural quote card with gradient background
class CulturalQuoteCard extends StatelessWidget {
  final String quote;
  final String author;
  final String? translation;
  final GradientType gradientType;
  final EdgeInsets padding;
  
  const CulturalQuoteCard({
    super.key,
    required this.quote,
    required this.author,
    this.translation,
    this.gradientType = GradientType.flag,
    this.padding = const EdgeInsets.all(16),
  });
  
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: padding,
      decoration: BoxDecoration(
        gradient: _getBackgroundGradient(),
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.1),
            blurRadius: 10,
            offset: const Offset(0, 5),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Quote text
          Text(
            '"$quote"',
            style: const TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.w500,
              color: Colors.white,
              height: 1.4,
            ),
          ),
          
          if (translation != null) ...[
            const SizedBox(height: 12),
            Text(
              translation!,
              style: TextStyle(
                fontSize: 14,
                color: Colors.white.withOpacity(0.8),
                fontStyle: FontStyle.italic,
              ),
            ),
          ],
          
          const SizedBox(height: 16),
          
          // Author
          Align(
            alignment: Alignment.centerRight,
            child: Text(
              '- $author',
              style: const TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w600,
                color: Colors.white,
              ),
            ),
          ),
        ],
      ),
    );
  }
  
  LinearGradient _getBackgroundGradient() {
    switch (gradientType) {
      case GradientType.flag:
        return const LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            Color(0xFFFF9933),
            Color(0xFFFFD700),
            Color(0xFF138808),
          ],
        );
      case GradientType.sunset:
        return const LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            Color(0xFFFF6B35),
            Color(0xFFFFD700),
            Color(0xFFFF1493),
          ],
        );
      case GradientType.ocean:
        return const LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            Color(0xFF00CED1),
            Color(0xFF4169E1),
            Color(0xFF000080),
          ],
        );
      case GradientType.forest:
        return const LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            Color(0xFF32CD32),
            Color(0xFF228B22),
            Color(0xFF006400),
          ],
        );
      case GradientType.festival:
        return const LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            Color(0xFFFFD700),
            Color(0xFFFF69B4),
            Color(0xFF9370DB),
          ],
        );
    }
  }
}

/// Cultural gradient button
class CulturalGradientButton extends StatelessWidget {
  final String text;
  final VoidCallback onPressed;
  final GradientType gradientType;
  final double height;
  final double borderRadius;
  final TextStyle? textStyle;
  
  const CulturalGradientButton({
    super.key,
    required this.text,
    required this.onPressed,
    this.gradientType = GradientType.flag,
    this.height = 56,
    this.borderRadius = 12,
    this.textStyle,
  });
  
  @override
  Widget build(BuildContext context) {
    return Container(
      height: height,
      decoration: BoxDecoration(
        gradient: _getGradient(),
        borderRadius: BorderRadius.circular(borderRadius),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.2),
            blurRadius: 8,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: onPressed,
          borderRadius: BorderRadius.circular(borderRadius),
          child: Center(
            child: Text(
              text,
              style: textStyle ?? const TextStyle(
                color: Colors.white,
                fontSize: 16,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
        ),
      ),
    );
  }
  
  LinearGradient _getGradient() {
    switch (gradientType) {
      case GradientType.flag:
        return const LinearGradient(
          colors: [
            Color(0xFFFF9933),
            Color(0xFFFFD700),
            Color(0xFF138808),
          ],
        );
      case GradientType.sunset:
        return const LinearGradient(
          colors: [
            Color(0xFFFF6B35),
            Color(0xFFFFD700),
          ],
        );
      case GradientType.ocean:
        return const LinearGradient(
          colors: [
            Color(0xFF00CED1),
            Color(0xFF4169E1),
          ],
        );
      case GradientType.forest:
        return const LinearGradient(
          colors: [
            Color(0xFF32CD32),
            Color(0xFF228B22),
          ],
        );
      case GradientType.festival:
        return const LinearGradient(
          colors: [
            Color(0xFFFFD700),
            Color(0xFFFF69B4),
            Color(0xFF9370DB),
          ],
        );
    }
  }
}