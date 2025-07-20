'use client'

import { motion } from 'framer-motion'
import { Calendar, Gift, Sparkles, Star } from 'lucide-react'

interface Festival {
  id: string
  name: string
  nameHindi: string
  date: Date
  description: string
  significance: string
  colors: {
    primary: string
    secondary: string
    accent: string
  }
  symbol: string
  isActive: boolean
  bonusPercentage: number
}

interface FestivalBannerProps {
  festival: Festival
}

export function FestivalBanner({ festival }: FestivalBannerProps) {
  const daysUntil = Math.ceil((festival.date.getTime() - Date.now()) / (1000 * 60 * 60 * 24))
  const isToday = daysUntil === 0
  const isPast = daysUntil < 0

  return (
    <motion.div
      className="relative overflow-hidden rounded-2xl bg-gradient-to-r shadow-2xl"
      style={{
        background: `linear-gradient(135deg, ${festival.colors.primary}20 0%, ${festival.colors.secondary}20 50%, ${festival.colors.accent}20 100%)`
      }}
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.6 }}
    >
      {/* Background Pattern */}
      <div className="absolute inset-0 opacity-10">
        <div className="absolute top-4 left-4 text-6xl animate-bounce-soft">
          {festival.symbol}
        </div>
        <div className="absolute top-8 right-8 text-4xl animate-pulse-slow">
          {festival.symbol}
        </div>
        <div className="absolute bottom-4 left-1/3 text-5xl animate-bounce-soft" style={{ animationDelay: '1s' }}>
          {festival.symbol}
        </div>
        <div className="absolute bottom-8 right-1/4 text-3xl animate-pulse-slow" style={{ animationDelay: '2s' }}>
          {festival.symbol}
        </div>
      </div>

      {/* Content */}
      <div className="relative z-10 p-8">
        <div className="max-w-6xl mx-auto">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-center">
            {/* Left Content */}
            <div className="space-y-6">
              <motion.div
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.6, delay: 0.2 }}
              >
                <div className="flex items-center space-x-3 mb-2">
                  <motion.div
                    className="text-4xl"
                    animate={{ 
                      rotate: [0, 10, -10, 0],
                      scale: [1, 1.1, 1]
                    }}
                    transition={{ 
                      duration: 2,
                      repeat: Infinity,
                      repeatType: "reverse"
                    }}
                  >
                    {festival.symbol}
                  </motion.div>
                  <div>
                    <h2 className="text-3xl font-bold text-gray-900 dark:text-dark-text">
                      {festival.name}
                    </h2>
                    <p className="text-xl font-hindi text-gray-700 dark:text-gray-300">
                      {festival.nameHindi}
                    </p>
                  </div>
                </div>
                
                <p className="text-lg text-gray-600 dark:text-gray-400 leading-relaxed">
                  {festival.description}
                </p>
              </motion.div>

              {/* Festival Benefits */}
              <motion.div
                className="flex flex-wrap gap-4"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.6, delay: 0.4 }}
              >
                <div className="flex items-center space-x-2 bg-white/80 dark:bg-dark-surface/80 backdrop-blur-sm rounded-full px-4 py-2">
                  <Gift className="w-5 h-5 text-primary-600" />
                  <span className="font-semibold text-primary-700 dark:text-primary-400">
                    {festival.bonusPercentage}% Bonus
                  </span>
                </div>
                
                <div className="flex items-center space-x-2 bg-white/80 dark:bg-dark-surface/80 backdrop-blur-sm rounded-full px-4 py-2">
                  <Sparkles className="w-5 h-5 text-orange-600" />
                  <span className="font-medium text-gray-700 dark:text-gray-300">
                    Special Rates
                  </span>
                </div>

                <div className="flex items-center space-x-2 bg-white/80 dark:bg-dark-surface/80 backdrop-blur-sm rounded-full px-4 py-2">
                  <Star className="w-5 h-5 text-yellow-600" />
                  <span className="font-medium text-gray-700 dark:text-gray-300">
                    Cultural Rewards
                  </span>
                </div>
              </motion.div>

              {/* CTA Button */}
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.6, delay: 0.6 }}
              >
                <button
                  className="btn-cultural text-lg px-8 py-3"
                  style={{
                    background: `linear-gradient(135deg, ${festival.colors.primary}, ${festival.colors.secondary})`
                  }}
                >
                  {isToday ? '🎉 Celebrate Now' : 
                   isPast ? '📜 View Past Celebration' : 
                   '🗓️ Set Reminder'}
                </button>
              </motion.div>
            </div>

            {/* Right Content - Festival Info */}
            <motion.div
              className="space-y-6"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.6, delay: 0.3 }}
            >
              {/* Date Countdown */}
              <div className="bg-white/90 dark:bg-dark-surface/90 backdrop-blur-sm rounded-2xl p-6 border border-white/20">
                <div className="flex items-center space-x-3 mb-4">
                  <Calendar className="w-6 h-6 text-primary-600" />
                  <h3 className="font-semibold text-gray-900 dark:text-dark-text">
                    {isToday ? 'Celebrating Today!' : 
                     isPast ? 'Past Celebration' : 
                     'Upcoming Celebration'}
                  </h3>
                </div>
                
                <div className="text-center">
                  <div className="text-3xl font-bold text-primary-600 mb-2">
                    {isToday ? '🎉 TODAY!' : 
                     isPast ? '✨ CELEBRATED' : 
                     `${Math.abs(daysUntil)} ${Math.abs(daysUntil) === 1 ? 'DAY' : 'DAYS'}`}
                  </div>
                  <div className="text-sm text-gray-600 dark:text-gray-400">
                    {festival.date.toLocaleDateString('en-IN', {
                      weekday: 'long',
                      year: 'numeric',
                      month: 'long',
                      day: 'numeric'
                    })}
                  </div>
                </div>
              </div>

              {/* Significance */}
              <div className="bg-white/90 dark:bg-dark-surface/90 backdrop-blur-sm rounded-2xl p-6 border border-white/20">
                <h4 className="font-semibold text-gray-900 dark:text-dark-text mb-3">
                  Cultural Significance
                </h4>
                <p className="text-sm text-gray-600 dark:text-gray-400 leading-relaxed">
                  {festival.significance}
                </p>
              </div>

              {/* Festival Stats */}
              <div className="bg-white/90 dark:bg-dark-surface/90 backdrop-blur-sm rounded-2xl p-6 border border-white/20">
                <h4 className="font-semibold text-gray-900 dark:text-dark-text mb-4">
                  Festival Impact
                </h4>
                <div className="grid grid-cols-2 gap-4 text-center">
                  <div>
                    <div className="text-2xl font-bold text-primary-600">
                      {Math.floor(Math.random() * 1000) + 500}
                    </div>
                    <div className="text-xs text-gray-600 dark:text-gray-400">
                      Participants
                    </div>
                  </div>
                  <div>
                    <div className="text-2xl font-bold text-green-600">
                      ₹{Math.floor(Math.random() * 10) + 5}L
                    </div>
                    <div className="text-xs text-gray-600 dark:text-gray-400">
                      Bonus Distributed
                    </div>
                  </div>
                  <div>
                    <div className="text-2xl font-bold text-blue-600">
                      {Math.floor(Math.random() * 50) + 100}
                    </div>
                    <div className="text-xs text-gray-600 dark:text-gray-400">
                      Transactions
                    </div>
                  </div>
                  <div>
                    <div className="text-2xl font-bold text-purple-600">
                      {festival.bonusPercentage}%
                    </div>
                    <div className="text-xs text-gray-600 dark:text-gray-400">
                      Bonus Rate
                    </div>
                  </div>
                </div>
              </div>
            </motion.div>
          </div>
        </div>
      </div>

      {/* Animated Particles */}
      <div className="absolute inset-0 pointer-events-none">
        {Array.from({ length: 20 }).map((_, i) => (
          <motion.div
            key={i}
            className="absolute w-2 h-2 rounded-full opacity-30"
            style={{
              background: festival.colors.accent,
              left: `${Math.random() * 100}%`,
              top: `${Math.random() * 100}%`,
            }}
            animate={{
              y: [0, -20, 0],
              opacity: [0.3, 0.7, 0.3],
              scale: [1, 1.2, 1],
            }}
            transition={{
              duration: 3 + Math.random() * 2,
              repeat: Infinity,
              delay: Math.random() * 2,
            }}
          />
        ))}
      </div>
    </motion.div>
  )
}