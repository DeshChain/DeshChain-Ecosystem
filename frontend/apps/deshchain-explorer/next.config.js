/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  images: {
    domains: ['ipfs.io', 'gateway.pinata.cloud'],
    formats: ['image/webp', 'image/avif'],
  },
  experimental: {
    appDir: true,
    serverComponentsExternalPackages: ['@cosmjs/stargate', '@cosmjs/tendermint-rpc'],
  },
  webpack: (config) => {
    config.resolve.fallback = {
      ...config.resolve.fallback,
      fs: false,
      net: false,
      tls: false,
      crypto: require.resolve('crypto-browserify'),
      stream: require.resolve('stream-browserify'),
      buffer: require.resolve('buffer'),
    };
    return config;
  },
  env: {
    NEXT_PUBLIC_CHAIN_ID: process.env.NEXT_PUBLIC_CHAIN_ID || 'deshchain-1',
    NEXT_PUBLIC_RPC_ENDPOINT: process.env.NEXT_PUBLIC_RPC_ENDPOINT || 'http://localhost:26657',
    NEXT_PUBLIC_REST_ENDPOINT: process.env.NEXT_PUBLIC_REST_ENDPOINT || 'http://localhost:1317',
    NEXT_PUBLIC_EXPLORER_NAME: 'DeshChain Explorer',
  },
  // Cultural theming support
  i18n: {
    locales: ['en', 'hi', 'bn', 'te', 'ta', 'gu', 'mr', 'kn', 'ml', 'pa'],
    defaultLocale: 'en',
    localeDetection: true,
  },
  // Performance optimizations
  compress: true,
  poweredByHeader: false,
  generateEtags: false,
  // Progressive Web App support
  pwa: {
    dest: 'public',
    register: true,
    skipWaiting: true,
    runtimeCaching: [
      {
        urlPattern: /^https?.*/,
        handler: 'NetworkFirst',
        options: {
          cacheName: 'offlineCache',
          expiration: {
            maxEntries: 200,
          },
        },
      },
    ],
  },
}

module.exports = nextConfig