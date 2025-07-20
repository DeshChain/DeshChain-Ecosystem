'use client'

import { motion } from 'framer-motion'
import { TrendingUp, TrendingDown, Minus } from 'lucide-react'

interface StatCardProps {
  title: string
  value: number | string
  icon: React.ReactNode
  trend?: number
  loading?: boolean
  color?: 'blue' | 'green' | 'purple' | 'orange' | 'red'
  format?: 'number' | 'currency' | 'percentage'
  subtitle?: string
}

const colorClasses = {
  blue: {
    bg: 'from-blue-500 to-blue-600',
    icon: 'bg-blue-100 text-blue-600 dark:bg-blue-900/20 dark:text-blue-400',
    trend: 'text-blue-600 dark:text-blue-400'
  },
  green: {
    bg: 'from-green-500 to-green-600',
    icon: 'bg-green-100 text-green-600 dark:bg-green-900/20 dark:text-green-400',
    trend: 'text-green-600 dark:text-green-400'
  },
  purple: {
    bg: 'from-purple-500 to-purple-600',
    icon: 'bg-purple-100 text-purple-600 dark:bg-purple-900/20 dark:text-purple-400',
    trend: 'text-purple-600 dark:text-purple-400'
  },
  orange: {
    bg: 'from-orange-500 to-orange-600',
    icon: 'bg-orange-100 text-orange-600 dark:bg-orange-900/20 dark:text-orange-400',
    trend: 'text-orange-600 dark:text-orange-400'
  },
  red: {
    bg: 'from-red-500 to-red-600',
    icon: 'bg-red-100 text-red-600 dark:bg-red-900/20 dark:text-red-400',
    trend: 'text-red-600 dark:text-red-400'
  }
}

export function StatCard({ 
  title, 
  value, 
  icon, 
  trend, 
  loading = false, 
  color = 'blue',
  format = 'number',
  subtitle 
}: StatCardProps) {
  const formatValue = (val: number | string) => {
    if (typeof val === 'string') return val
    
    switch (format) {
      case 'currency':
        return new Intl.NumberFormat('en-IN', {
          style: 'currency',
          currency: 'INR',
          notation: 'compact',
          maximumFractionDigits: 1
        }).format(val)
      case 'percentage':
        return `${val.toFixed(1)}%`
      default:
        return new Intl.NumberFormat('en-IN', {
          notation: 'compact',
          maximumFractionDigits: 1
        }).format(val)
    }
  }

  const getTrendIcon = () => {
    if (!trend || trend === 0) return <Minus className="w-4 h-4" />
    return trend > 0 ? <TrendingUp className="w-4 h-4" /> : <TrendingDown className="w-4 h-4" />
  }

  const getTrendColor = () => {
    if (!trend || trend === 0) return 'text-gray-500'
    return trend > 0 ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'
  }

  if (loading) {
    return (
      <div className="card p-6 animate-pulse">
        <div className="flex items-center justify-between mb-4">
          <div className="w-12 h-12 bg-gray-200 dark:bg-gray-700 rounded-xl" />
          <div className="w-16 h-6 bg-gray-200 dark:bg-gray-700 rounded" />
        </div>
        <div className="space-y-2">
          <div className="w-24 h-4 bg-gray-200 dark:bg-gray-700 rounded" />
          <div className="w-32 h-8 bg-gray-200 dark:bg-gray-700 rounded" />
        </div>
      </div>
    )
  }

  return (
    <motion.div
      className="card p-6 hover:shadow-lg transition-all duration-300"
      whileHover={{ y: -2 }}
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
    >
      <div className="flex items-center justify-between mb-4">
        <div className={`p-3 rounded-xl ${colorClasses[color].icon}`}>
          {icon}
        </div>
        
        {trend !== undefined && (
          <div className={`flex items-center space-x-1 text-sm font-medium ${getTrendColor()}`}>
            {getTrendIcon()}
            <span>{Math.abs(trend).toFixed(1)}%</span>
          </div>
        )}
      </div>

      <div className="space-y-1">
        <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400">
          {title}
        </h3>
        <motion.div
          className="text-2xl font-bold text-gray-900 dark:text-dark-text"
          key={value.toString()}
          initial={{ scale: 1.1 }}
          animate={{ scale: 1 }}
          transition={{ duration: 0.2 }}
        >
          {formatValue(value)}
        </motion.div>
        {subtitle && (
          <p className="text-xs text-gray-500 dark:text-gray-400">
            {subtitle}
          </p>
        )}
      </div>

      {/* Optional gradient accent */}
      <div className={`absolute bottom-0 left-0 right-0 h-1 bg-gradient-to-r ${colorClasses[color].bg} rounded-b-lg`} />
    </motion.div>
  )
}