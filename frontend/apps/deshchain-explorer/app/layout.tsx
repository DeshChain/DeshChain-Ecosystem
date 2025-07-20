import './globals.css'
import { Inter, Noto_Sans_Devanagari } from 'next/font/google'
import { Providers } from './providers'
import { Header } from '@/components/layout/Header'
import { Footer } from '@/components/layout/Footer'
import { Sidebar } from '@/components/layout/Sidebar'
import { CulturalProvider } from '@/providers/CulturalProvider'
import { ExplorerProvider } from '@/providers/ExplorerProvider'
import { Toaster } from 'react-hot-toast'

const inter = Inter({ 
  subsets: ['latin'],
  variable: '--font-inter',
})

const notoSansDevanagari = Noto_Sans_Devanagari({
  subsets: ['devanagari'],
  variable: '--font-hindi',
})

export const metadata = {
  title: 'DeshChain Explorer | Blockchain for Cultural Heritage',
  description: 'Explore the DeshChain blockchain - where cultural heritage meets cutting-edge DeFi. Track transactions, lending modules, festivals, and more.',
  keywords: ['blockchain', 'explorer', 'deshchain', 'cultural', 'heritage', 'defi', 'lending'],
  authors: [{ name: 'DeshChain Team' }],
  creator: 'DeshChain',
  publisher: 'DeshChain',
  formatDetection: {
    email: false,
    address: false,
    telephone: false,
  },
  viewport: 'width=device-width, initial-scale=1',
  themeColor: '#f97316',
  colorScheme: 'light dark',
  openGraph: {
    type: 'website',
    locale: 'en_US',
    url: 'https://explorer.deshchain.network',
    title: 'DeshChain Explorer',
    description: 'Explore the DeshChain blockchain ecosystem',
    siteName: 'DeshChain Explorer',
    images: [
      {
        url: '/og-image.png',
        width: 1200,
        height: 630,
        alt: 'DeshChain Explorer',
      },
    ],
  },
  twitter: {
    card: 'summary_large_image',
    title: 'DeshChain Explorer',
    description: 'Explore the DeshChain blockchain ecosystem',
    images: ['/og-image.png'],
    creator: '@deshchain',
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      'max-video-preview': -1,
      'max-image-preview': 'large',
      'max-snippet': -1,
    },
  },
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`${inter.variable} ${notoSansDevanagari.variable} font-sans`}>
        <Providers>
          <CulturalProvider>
            <ExplorerProvider>
              <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-dark-bg dark:to-dark-surface">
                <Header />
                <div className="flex">
                  <Sidebar />
                  <main className="flex-1 ml-64 transition-all duration-300">
                    <div className="container mx-auto px-6 py-8">
                      {children}
                    </div>
                  </main>
                </div>
                <Footer />
                <Toaster
                  position="top-right"
                  toastOptions={{
                    duration: 4000,
                    style: {
                      background: '#363636',
                      color: '#fff',
                      borderRadius: '10px',
                    },
                    success: {
                      style: {
                        background: '#22c55e',
                      },
                    },
                    error: {
                      style: {
                        background: '#ef4444',
                      },
                    },
                  }}
                />
              </div>
            </ExplorerProvider>
          </CulturalProvider>
        </Providers>
      </body>
    </html>
  )
}