# Contributing to DhanSetu Mobile

Thank you for your interest in contributing to DhanSetu! We welcome contributions from the community to help build India's revolutionary DeFi super app.

## 🤝 Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md) to ensure a welcoming environment for all contributors.

## 🚀 Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/dhansetu-mobile.git
   cd dhansetu-mobile
   ```
3. Install dependencies:
   ```bash
   npm install
   ```
4. Create a new branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## 📝 Development Guidelines

### Code Style
- Follow the existing code style and patterns
- Use TypeScript for all new code
- Run `npm run lint` before committing
- Ensure `npm run type-check` passes without errors

### Component Guidelines
- Create reusable components in `src/components/common`
- Cultural components go in `src/components/cultural`
- Follow the atomic design pattern
- Add proper TypeScript interfaces for props

### State Management
- Use Redux Toolkit for global state
- Use React Context for cross-cutting concerns
- Keep component state local when possible

### Blockchain Integration
- Use the existing DeshChainClient for blockchain interactions
- Follow the established patterns in `src/services/blockchain`
- Add proper error handling for all blockchain calls

### Cultural Sensitivity
- Respect Indian cultural values in all contributions
- Use appropriate language in comments and documentation
- Test UI with different Indian language settings
- Consider festival themes when adding new features

## 🧪 Testing

### Unit Tests
- Write tests for all new functions and components
- Aim for >80% code coverage
- Run tests with `npm test`

### Integration Tests
- Test blockchain interactions with mock data
- Test navigation flows end-to-end
- Ensure wallet operations are properly tested

### Manual Testing
- Test on both Android and iOS devices
- Test with different screen sizes
- Test with poor network conditions
- Test with different language settings

## 🔄 Pull Request Process

1. **Before submitting:**
   - Update documentation for any API changes
   - Add tests for new functionality
   - Ensure all tests pass
   - Run linting and type checking
   - Update README.md if needed

2. **PR Description:**
   - Clearly describe the changes
   - Link to any related issues
   - Include screenshots for UI changes
   - List any breaking changes

3. **Review Process:**
   - At least one maintainer approval required
   - All CI checks must pass
   - No merge conflicts
   - Follow up on review comments promptly

## 🏷️ Commit Convention

We follow conventional commits:

```
feat: Add new feature
fix: Fix bug
docs: Update documentation
style: Format code
refactor: Refactor code
test: Add tests
chore: Update dependencies
```

Examples:
- `feat: Add biometric authentication for wallet`
- `fix: Resolve transaction signing issue`
- `docs: Update README with new features`

## 🐛 Reporting Issues

### Bug Reports
Include:
- Device and OS version
- Steps to reproduce
- Expected vs actual behavior
- Screenshots if applicable
- Error logs from console

### Feature Requests
Include:
- Clear description of the feature
- Use cases and benefits
- Mockups or examples if possible
- Impact on existing features

## 💡 Feature Ideas

We especially welcome contributions in these areas:
- Improved cultural UI components
- Additional language support
- Performance optimizations
- Accessibility improvements
- Security enhancements
- New DeFi features

## 📱 Platform-Specific Guidelines

### Android
- Test on Android 7.0+
- Consider low-end device performance
- Follow Material Design guidelines where appropriate

### iOS
- Test on iOS 13.0+
- Follow iOS Human Interface Guidelines
- Ensure proper safe area handling

## 🌐 Localization

When adding new strings:
1. Add to `src/locales/en.json` first
2. Use the translation key in code
3. Create a separate PR for translations
4. Ensure all UI text is translatable

## 📞 Getting Help

- Join our [Discord](https://discord.gg/deshchain)
- Check existing issues and discussions
- Ask in the #dev-help channel
- Tag @maintainers for urgent issues

## 🙏 Recognition

All contributors will be recognized in:
- The project README
- Release notes
- Our website contributors page

Thank you for helping build the future of Indian DeFi! 🇮🇳

---

Jai Hind! 🇮🇳