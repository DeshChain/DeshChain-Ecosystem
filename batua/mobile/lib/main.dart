import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:intl/date_symbol_data_local.dart';

import 'core/storage/secure_storage.dart';
import 'ui/themes/app_theme.dart';
import 'ui/screens/splash_screen.dart';
import 'utils/constants.dart';
import 'providers/app_providers.dart';
import 'services/localization_service.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Initialize date formatting
  await initializeDateFormatting('hi_IN', null);
  
  // Set preferred orientations
  await SystemChrome.setPreferredOrientations([
    DeviceOrientation.portraitUp,
    DeviceOrientation.portraitDown,
  ]);
  
  // Initialize Hive
  await Hive.initFlutter();
  
  // Initialize secure storage
  await SecureStorage.init();
  
  // Set system UI overlay style
  SystemChrome.setSystemUIOverlayStyle(
    const SystemUiOverlayStyle(
      statusBarColor: Colors.transparent,
      statusBarIconBrightness: Brightness.dark,
      systemNavigationBarColor: AppColors.white,
      systemNavigationBarIconBrightness: Brightness.dark,
    ),
  );
  
  runApp(
    const ProviderScope(
      child: BatuaApp(),
    ),
  );
}

class BatuaApp extends ConsumerWidget {
  const BatuaApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final locale = ref.watch(localeProvider);
    final themeMode = ref.watch(themeModeProvider);
    
    return MaterialApp(
      title: 'Batua - बटुआ',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.lightTheme,
      darkTheme: AppTheme.darkTheme,
      themeMode: themeMode,
      locale: locale,
      supportedLocales: LocalizationService.supportedLocales,
      localizationsDelegates: const [
        AppLocalizations.delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      home: const SplashScreen(),
    );
  }
}

// App Colors with Indian Flag Theme
class AppColors {
  static const Color saffron = Color(0xFFFF9933);
  static const Color white = Color(0xFFFFFFFF);
  static const Color green = Color(0xFF138808);
  static const Color blue = Color(0xFF000080);
  
  static const Color primaryDark = Color(0xFF1A1A2E);
  static const Color primaryLight = Color(0xFFF5F5F5);
  static const Color accent = Color(0xFFFFD700);
  
  static const Color success = Color(0xFF4CAF50);
  static const Color error = Color(0xFFF44336);
  static const Color warning = Color(0xFFFF9800);
  static const Color info = Color(0xFF2196F3);
  
  // Gradient for Batua logo
  static const LinearGradient batuaGradient = LinearGradient(
    begin: Alignment.topCenter,
    end: Alignment.bottomCenter,
    colors: [saffron, white, green],
    stops: [0.0, 0.5, 1.0],
  );
  
  // Festival theme colors
  static const Map<String, Color> festivalColors = {
    'diwali': Color(0xFFFFD700),
    'holi': Color(0xFFFF1493),
    'dussehra': Color(0xFFFF6347),
    'ganesh': Color(0xFFFF8C00),
    'navratri': Color(0xFF8B008B),
  };
}