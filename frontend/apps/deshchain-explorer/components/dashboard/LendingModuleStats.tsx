'use client'

import { motion } from 'framer-motion'
import { 
  Sprout, 
  Building2, 
  GraduationCap, 
  TrendingUp,
  Users,
  DollarSign,
  Percent,
  AlertTriangle
} from 'lucide-react'
import { StatCard } from './StatCard'

interface LendingStatsProps {
  stats: {
    krishiMitra: {
      totalLoans: number
      totalDisbursed: string
      averageRate: number
      activeLoans: number
      defaultRate: number
    }
    vyavasayaMitra: {
      totalLoans: number
      totalDisbursed: string
      averageRate: number
      activeLoans: number
      defaultRate: number
    }
    shikshamitra: {
      totalLoans: number
      totalDisbursed: string
      averageRate: number
      activeLoans: number
      defaultRate: number
    }
    combined: {
      totalDisbursed: string
      avgInterestRate: number
      totalBorrowers: number
      totalLenders: number
    }
  } | null | undefined
  loading?: boolean
}

export function LendingModuleStats({ stats, loading = false }: LendingStatsProps) {
  if (loading || !stats) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div className="w-48 h-8 bg-gray-200 dark:bg-gray-700 rounded animate-pulse" />
          <div className="w-24 h-6 bg-gray-200 dark:bg-gray-700 rounded animate-pulse" />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="card p-6 animate-pulse">
              <div className="w-12 h-12 bg-gray-200 dark:bg-gray-700 rounded-xl mb-4" />
              <div className="space-y-2">
                <div className="w-24 h-4 bg-gray-200 dark:bg-gray-700 rounded" />
                <div className="w-32 h-8 bg-gray-200 dark:bg-gray-700 rounded" />
              </div>
            </div>
          ))}
        </div>
      </div>
    )
  }

  const modules = [
    {
      name: 'Krishi Mitra',
      nameHindi: 'कृषि मित्र',
      description: 'Agricultural Lending',
      icon: <Sprout className="w-6 h-6" />,
      color: 'green' as const,
      data: stats.krishiMitra,
      rateRange: '6-9%'
    },
    {
      name: 'Vyavasaya Mitra',
      nameHindi: 'व्यवसाय मित्र',
      description: 'Business Lending',
      icon: <Building2 className="w-6 h-6" />,
      color: 'blue' as const,
      data: stats.vyavasayaMitra,
      rateRange: '8-12%'
    },
    {
      name: 'Shiksha Mitra',
      nameHindi: 'शिक्षा मित्र',
      description: 'Education Loans',
      icon: <GraduationCap className="w-6 h-6" />,
      color: 'purple' as const,
      data: stats.shikshamitra,
      rateRange: '4-7%'
    }
  ]

  return (
    <motion.div
      className="space-y-6"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
    >
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-gray-900 dark:text-dark-text">
            Lending Modules
          </h2>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Empowering India through accessible finance
          </p>
        </div>
        <motion.button
          className="btn-lending"
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
        >
          View All
        </motion.button>
      </div>

      {/* Combined Stats */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatCard
          title="Total Disbursed"
          value={stats.combined.totalDisbursed}
          icon={<DollarSign className="w-6 h-6" />}
          color="green"
          subtitle="Across all modules"
        />
        <StatCard
          title="Average Rate"
          value={`${stats.combined.avgInterestRate.toFixed(1)}%`}
          icon={<Percent className="w-6 h-6" />}
          color="blue"
          subtitle="Weighted average"
        />
        <StatCard
          title="Total Borrowers"
          value={stats.combined.totalBorrowers}
          icon={<Users className="w-6 h-6" />}
          color="purple"
          subtitle="Active users"
        />
        <StatCard
          title="Growth Rate"
          value="23.4%"
          icon={<TrendingUp className="w-6 h-6" />}
          color="orange"
          trend={23.4}
          subtitle="Month over month"
        />
      </div>

      {/* Module-wise Stats */}
      <div className="space-y-6">
        <h3 className="text-xl font-semibold text-gray-900 dark:text-dark-text">
          Module Performance
        </h3>
        
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {modules.map((module, index) => (
            <motion.div
              key={module.name}
              className="card-lending p-6"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: index * 0.1 }}
              whileHover={{ y: -4 }}
            >
              {/* Module Header */}
              <div className="flex items-center space-x-3 mb-4">
                <div className={`p-2 rounded-lg ${
                  module.color === 'green' ? 'bg-green-100 text-green-600 dark:bg-green-900/20 dark:text-green-400' :
                  module.color === 'blue' ? 'bg-blue-100 text-blue-600 dark:bg-blue-900/20 dark:text-blue-400' :
                  'bg-purple-100 text-purple-600 dark:bg-purple-900/20 dark:text-purple-400'
                }`}>
                  {module.icon}
                </div>
                <div>
                  <h4 className="font-semibold text-gray-900 dark:text-dark-text">
                    {module.name}
                  </h4>
                  <p className="text-sm text-gray-600 dark:text-gray-400 font-hindi">
                    {module.nameHindi}
                  </p>
                </div>
              </div>

              {/* Key Metrics */}
              <div className="space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600 dark:text-gray-400">Total Loans</span>
                  <span className="font-semibold text-gray-900 dark:text-dark-text">
                    {module.data.totalLoans.toLocaleString()}
                  </span>
                </div>

                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600 dark:text-gray-400">Disbursed</span>
                  <span className="font-semibold text-gray-900 dark:text-dark-text">
                    {module.data.totalDisbursed}
                  </span>
                </div>

                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600 dark:text-gray-400">Avg Rate</span>
                  <span className="font-semibold text-gray-900 dark:text-dark-text">
                    {module.data.averageRate.toFixed(1)}%
                  </span>
                </div>

                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600 dark:text-gray-400">Active</span>
                  <span className="font-semibold text-gray-900 dark:text-dark-text">
                    {module.data.activeLoans.toLocaleString()}
                  </span>
                </div>

                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600 dark:text-gray-400">Default Rate</span>
                  <div className="flex items-center space-x-1">
                    {module.data.defaultRate > 5 && (
                      <AlertTriangle className="w-4 h-4 text-yellow-500" />
                    )}
                    <span className={`font-semibold ${
                      module.data.defaultRate > 5 ? 'text-yellow-600' : 
                      module.data.defaultRate > 3 ? 'text-orange-600' : 'text-green-600'
                    }`}>
                      {module.data.defaultRate.toFixed(1)}%
                    </span>
                  </div>
                </div>

                {/* Rate Range */}
                <div className="pt-3 border-t border-gray-200 dark:border-gray-700">
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-gray-600 dark:text-gray-400">Rate Range</span>
                    <span className="text-sm font-medium text-gray-900 dark:text-dark-text">
                      {module.rateRange}
                    </span>
                  </div>
                </div>

                {/* Progress Bar */}
                <div className="space-y-2">
                  <div className="flex justify-between text-xs text-gray-600 dark:text-gray-400">
                    <span>Utilization</span>
                    <span>{((module.data.activeLoans / module.data.totalLoans) * 100).toFixed(1)}%</span>
                  </div>
                  <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                    <motion.div
                      className={`h-2 rounded-full ${
                        module.color === 'green' ? 'bg-green-500' :
                        module.color === 'blue' ? 'bg-blue-500' : 'bg-purple-500'
                      }`}
                      initial={{ width: 0 }}
                      animate={{ width: `${(module.data.activeLoans / module.data.totalLoans) * 100}%` }}
                      transition={{ duration: 1, delay: 0.5 }}
                    />
                  </div>
                </div>
              </div>

              {/* Action Button */}
              <motion.button
                className={`w-full mt-4 py-2 px-4 rounded-lg text-sm font-medium transition-colors ${
                  module.color === 'green' ? 'bg-green-50 text-green-700 hover:bg-green-100 dark:bg-green-900/20 dark:text-green-400' :
                  module.color === 'blue' ? 'bg-blue-50 text-blue-700 hover:bg-blue-100 dark:bg-blue-900/20 dark:text-blue-400' :
                  'bg-purple-50 text-purple-700 hover:bg-purple-100 dark:bg-purple-900/20 dark:text-purple-400'
                }`}
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                View Details
              </motion.button>
            </motion.div>
          ))}
        </div>
      </div>

      {/* Quick Insights */}
      <div className="card p-6">
        <h4 className="font-semibold text-gray-900 dark:text-dark-text mb-4">
          Quick Insights
        </h4>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
          <div className="flex items-center space-x-2">
            <div className="w-2 h-2 bg-green-500 rounded-full" />
            <span className="text-gray-600 dark:text-gray-400">
              Education loans have the lowest default rate at {stats.shikshamitra.defaultRate}%
            </span>
          </div>
          <div className="flex items-center space-x-2">
            <div className="w-2 h-2 bg-blue-500 rounded-full" />
            <span className="text-gray-600 dark:text-gray-400">
              Business lending shows highest growth with 34% increase
            </span>
          </div>
          <div className="flex items-center space-x-2">
            <div className="w-2 h-2 bg-purple-500 rounded-full" />
            <span className="text-gray-600 dark:text-gray-400">
              Agricultural loans serve 67% rural population
            </span>
          </div>
        </div>
      </div>
    </motion.div>
  )
}