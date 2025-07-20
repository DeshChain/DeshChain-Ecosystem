'use client'

import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { 
  Activity, 
  TrendingUp, 
  Users, 
  Coins, 
  Clock,
  ArrowRight,
  Sparkles,
  Heart,
  Zap
} from 'lucide-react'

import { StatCard } from '@/components/dashboard/StatCard'
import { ChainOverview } from '@/components/dashboard/ChainOverview'
import { RecentTransactions } from '@/components/dashboard/RecentTransactions'
import { LendingModuleStats } from '@/components/dashboard/LendingModuleStats'
import { CulturalHighlights } from '@/components/dashboard/CulturalHighlights'
import { FestivalBanner } from '@/components/cultural/FestivalBanner'
import { QuickSearch } from '@/components/search/QuickSearch'
import { NetworkStatus } from '@/components/network/NetworkStatus'
import { useChainData } from '@/hooks/useChainData'
import { useCultural } from '@/hooks/useCultural'

export default function HomePage() {
  const { 
    chainStats, 
    recentTransactions, 
    lendingStats, 
    loading, 
    error 
  } = useChainData()
  
  const { 
    currentFestival, 
    dailyQuote, 
    culturalEvents 
  } = useCultural()

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1
      }
    }
  }

  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    visible: {
      y: 0,
      opacity: 1,
      transition: {
        duration: 0.5
      }
    }
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-96">
        <div className="text-center">
          <div className="text-red-500 text-6xl mb-4">⚠️</div>
          <h2 className="text-2xl font-bold text-gray-900 dark:text-dark-text mb-2">
            Connection Error
          </h2>
          <p className="text-gray-600 dark:text-gray-400 mb-4">
            Unable to connect to DeshChain network
          </p>
          <button 
            onClick={() => window.location.reload()}
            className="btn-primary"
          >
            Retry Connection
          </button>
        </div>
      </div>
    )
  }

  return (
    <motion.div
      variants={containerVariants}
      initial="hidden"
      animate="visible"
      className="space-y-8"
    >
      {/* Header Section */}
      <motion.div variants={itemVariants} className="text-center space-y-4">
        <div className="flex items-center justify-center space-x-2">
          <Sparkles className="w-8 h-8 text-primary-500" />
          <h1 className="text-4xl font-bold text-cultural">
            DeshChain Explorer
          </h1>
          <Sparkles className="w-8 h-8 text-primary-500" />
        </div>
        <p className="text-xl text-gray-600 dark:text-gray-400 max-w-3xl mx-auto">
          Explore the blockchain that bridges cultural heritage with cutting-edge DeFi. 
          Discover lending modules, festival celebrations, and transparent transactions.
        </p>
        <NetworkStatus />
      </motion.div>

      {/* Festival Banner */}
      {currentFestival && (
        <motion.div variants={itemVariants}>
          <FestivalBanner festival={currentFestival} />
        </motion.div>
      )}

      {/* Quick Search */}
      <motion.div variants={itemVariants}>
        <QuickSearch />
      </motion.div>

      {/* Main Stats Grid */}
      <motion.div 
        variants={itemVariants}
        className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6"
      >
        <StatCard
          title="Block Height"
          value={chainStats?.blockHeight || 0}
          icon={<Activity className="w-6 h-6" />}
          trend={chainStats?.blockTrend}
          loading={loading}
          color="blue"
        />
        <StatCard
          title="Total Transactions"
          value={chainStats?.totalTransactions || 0}
          icon={<TrendingUp className="w-6 h-6" />}
          trend={chainStats?.txTrend}
          loading={loading}
          color="green"
        />
        <StatCard
          title="Active Validators"
          value={chainStats?.activeValidators || 0}
          icon={<Users className="w-6 h-6" />}
          loading={loading}
          color="purple"
        />
        <StatCard
          title="NAMO Supply"
          value={chainStats?.namoSupply || 0}
          icon={<Coins className="w-6 h-6" />}
          trend={chainStats?.supplyTrend}
          loading={loading}
          color="orange"
          format="currency"
        />
      </motion.div>

      {/* Chain Overview and Cultural Highlights */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <motion.div variants={itemVariants} className="lg:col-span-2">
          <ChainOverview 
            chainStats={chainStats}
            loading={loading}
          />
        </motion.div>
        <motion.div variants={itemVariants}>
          <CulturalHighlights 
            dailyQuote={dailyQuote}
            culturalEvents={culturalEvents}
            loading={loading}
          />
        </motion.div>
      </div>

      {/* Lending Module Statistics */}
      <motion.div variants={itemVariants}>
        <LendingModuleStats 
          stats={lendingStats}
          loading={loading}
        />
      </motion.div>

      {/* Recent Transactions */}
      <motion.div variants={itemVariants}>
        <RecentTransactions 
          transactions={recentTransactions}
          loading={loading}
        />
      </motion.div>

      {/* Quick Links Grid */}
      <motion.div variants={itemVariants} className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <QuickLinkCard
          title="Lending Modules"
          description="Explore agricultural, business, and education loans"
          icon={<Heart className="w-8 h-8" />}
          href="/lending"
          color="from-lending-krishi to-lending-shiksha"
        />
        <QuickLinkCard
          title="Sikkebaaz Launchpad"
          description="Discover and launch community tokens"
          icon={<Zap className="w-8 h-8" />}
          href="/sikkebaaz"
          color="from-purple-500 to-pink-500"
        />
        <QuickLinkCard
          title="Money Order DEX"
          description="Traditional money transfers on blockchain"
          icon={<TrendingUp className="w-8 h-8" />}
          href="/money-order"
          color="from-blue-500 to-cyan-500"
        />
        <QuickLinkCard
          title="Cultural Heritage"
          description="Explore festivals, quotes, and traditions"
          icon={<Sparkles className="w-8 h-8" />}
          href="/cultural"
          color="from-cultural-saffron to-cultural-green"
        />
        <QuickLinkCard
          title="Governance"
          description="Participate in community decisions"
          icon={<Users className="w-8 h-8" />}
          href="/governance"
          color="from-gray-600 to-gray-800"
        />
        <QuickLinkCard
          title="Validators"
          description="Network security and validation"
          icon={<Activity className="w-8 h-8" />}
          href="/validators"
          color="from-green-500 to-emerald-600"
        />
      </motion.div>

      {/* Footer CTA */}
      <motion.div 
        variants={itemVariants}
        className="text-center py-12 bg-gradient-to-r from-primary-50 to-orange-50 dark:from-primary-900/20 dark:to-orange-900/20 rounded-2xl"
      >
        <h2 className="text-2xl font-bold text-gray-900 dark:text-dark-text mb-4">
          Ready to Explore DeshChain?
        </h2>
        <p className="text-gray-600 dark:text-gray-400 mb-6 max-w-2xl mx-auto">
          Dive deeper into the blockchain that's revolutionizing finance while preserving cultural heritage. 
          Join millions of users building the future of decentralized India.
        </p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <button className="btn-cultural">
            Start Exploring
          </button>
          <button className="btn-secondary">
            Learn More
          </button>
        </div>
      </motion.div>
    </motion.div>
  )
}

interface QuickLinkCardProps {
  title: string
  description: string
  icon: React.ReactNode
  href: string
  color: string
}

function QuickLinkCard({ title, description, icon, href, color }: QuickLinkCardProps) {
  return (
    <motion.a
      href={href}
      className={`block p-6 rounded-xl bg-gradient-to-br ${color} text-white hover:scale-105 transition-all duration-300 shadow-lg hover:shadow-xl group`}
      whileHover={{ scale: 1.02 }}
      whileTap={{ scale: 0.98 }}
    >
      <div className="flex items-start space-x-4">
        <div className="flex-shrink-0">
          {icon}
        </div>
        <div className="flex-1 min-w-0">
          <h3 className="text-lg font-semibold mb-2">{title}</h3>
          <p className="text-sm opacity-90 mb-3">{description}</p>
          <div className="flex items-center text-sm font-medium group-hover:translate-x-1 transition-transform duration-200">
            Explore <ArrowRight className="w-4 h-4 ml-1" />
          </div>
        </div>
      </div>
    </motion.a>
  )
}