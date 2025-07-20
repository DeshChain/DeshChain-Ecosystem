'use client'

import { useState, useRef, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { 
  Search, 
  Hash, 
  User, 
  Blocks, 
  Zap,
  ArrowRight,
  Clock,
  Sparkles
} from 'lucide-react'

interface SearchSuggestion {
  type: 'transaction' | 'block' | 'address' | 'validator' | 'lending'
  title: string
  subtitle: string
  hash: string
  icon: React.ReactNode
}

export function QuickSearch() {
  const [query, setQuery] = useState('')
  const [isOpen, setIsOpen] = useState(false)
  const [suggestions, setSuggestions] = useState<SearchSuggestion[]>([])
  const inputRef = useRef<HTMLInputElement>(null)

  // Mock suggestions based on query
  useEffect(() => {
    if (query.length > 2) {
      const mockSuggestions: SearchSuggestion[] = [
        {
          type: 'transaction',
          title: `Transaction ${query}...`,
          subtitle: 'Recent transfer: 100 NAMO',
          hash: `0x${query}${'0'.repeat(60 - query.length)}`,
          icon: <Hash className="w-4 h-4" />
        },
        {
          type: 'block',
          title: `Block #${Math.floor(Math.random() * 1000000)}`,
          subtitle: '47 transactions, 6 seconds ago',
          hash: Math.floor(Math.random() * 1000000).toString(),
          icon: <Blocks className="w-4 h-4" />
        },
        {
          type: 'address',
          title: `deshchain1${query}...`,
          subtitle: 'Balance: 1,234 NAMO',
          hash: `deshchain1${query}${'x'.repeat(35 - query.length)}`,
          icon: <User className="w-4 h-4" />
        },
        {
          type: 'lending',
          title: 'Krishi Mitra Loan',
          subtitle: 'Agricultural loan: ₹50,000',
          hash: `loan_${query}`,
          icon: <Sparkles className="w-4 h-4" />
        }
      ]
      setSuggestions(mockSuggestions)
    } else {
      setSuggestions([])
    }
  }, [query])

  const handleSearch = (searchValue?: string) => {
    const searchQuery = searchValue || query
    if (searchQuery.trim()) {
      window.location.href = `/search?q=${encodeURIComponent(searchQuery.trim())}`
    }
  }

  const recentSearches = [
    'deshchain1abc...def',
    'Block #1234567',
    'Festival transactions',
    'Lending statistics'
  ]

  const quickLinks = [
    { label: 'Latest Block', icon: <Blocks className="w-4 h-4" />, href: '/blocks' },
    { label: 'Top Validators', icon: <Zap className="w-4 h-4" />, href: '/validators' },
    { label: 'Lending Stats', icon: <Sparkles className="w-4 h-4" />, href: '/lending' },
    { label: 'Governance', icon: <User className="w-4 h-4" />, href: '/governance' }
  ]

  return (
    <div className="relative max-w-4xl mx-auto">
      <motion.div
        className={`relative transition-all duration-300 ${
          isOpen ? 'transform scale-105' : ''
        }`}
        whileFocus={{ scale: 1.02 }}
      >
        <div className="relative">
          <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
          <input
            ref={inputRef}
            type="text"
            placeholder="Search transactions, blocks, addresses, or lending data..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onFocus={() => setIsOpen(true)}
            onBlur={() => setTimeout(() => setIsOpen(false), 200)}
            onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
            className="w-full pl-12 pr-16 py-4 text-lg bg-white dark:bg-dark-surface border-2 border-gray-200 dark:border-dark-border rounded-2xl focus:outline-none focus:border-primary-500 focus:ring-4 focus:ring-primary-100 dark:focus:ring-primary-900/20 transition-all duration-200 shadow-sm hover:shadow-md"
          />
          {query && (
            <motion.button
              onClick={() => handleSearch()}
              className="absolute right-3 top-1/2 transform -translate-y-1/2 bg-primary-500 text-white px-4 py-2 rounded-lg hover:bg-primary-600 transition-colors"
              initial={{ opacity: 0, x: 10 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: 10 }}
            >
              Search
            </motion.button>
          )}
        </div>

        {/* Search Dropdown */}
        <AnimatePresence>
          {isOpen && (
            <motion.div
              className="absolute top-full mt-2 w-full bg-white dark:bg-dark-surface border border-gray-200 dark:border-dark-border rounded-xl shadow-xl z-50 overflow-hidden"
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              transition={{ duration: 0.2 }}
            >
              {/* Suggestions */}
              {suggestions.length > 0 && (
                <div className="p-4">
                  <h4 className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-3">
                    Search Results
                  </h4>
                  <div className="space-y-2">
                    {suggestions.map((suggestion, index) => (
                      <motion.button
                        key={index}
                        onClick={() => handleSearch(suggestion.hash)}
                        className="w-full flex items-center space-x-3 p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors text-left"
                        whileHover={{ x: 4 }}
                      >
                        <div className={`p-2 rounded-lg ${
                          suggestion.type === 'transaction' ? 'bg-blue-100 text-blue-600' :
                          suggestion.type === 'block' ? 'bg-green-100 text-green-600' :
                          suggestion.type === 'address' ? 'bg-purple-100 text-purple-600' :
                          'bg-orange-100 text-orange-600'
                        } dark:bg-opacity-20`}>
                          {suggestion.icon}
                        </div>
                        <div className="flex-1 min-w-0">
                          <div className="font-medium text-gray-900 dark:text-dark-text">
                            {suggestion.title}
                          </div>
                          <div className="text-sm text-gray-500 dark:text-gray-400">
                            {suggestion.subtitle}
                          </div>
                        </div>
                        <ArrowRight className="w-4 h-4 text-gray-400" />
                      </motion.button>
                    ))}
                  </div>
                </div>
              )}

              {/* Recent Searches or Quick Links */}
              {suggestions.length === 0 && (
                <div className="p-4">
                  {query.length === 0 && (
                    <>
                      {/* Quick Links */}
                      <div className="mb-6">
                        <h4 className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-3">
                          Quick Links
                        </h4>
                        <div className="grid grid-cols-2 gap-2">
                          {quickLinks.map((link) => (
                            <motion.a
                              key={link.label}
                              href={link.href}
                              className="flex items-center space-x-2 p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors text-sm"
                              whileHover={{ x: 2 }}
                            >
                              {link.icon}
                              <span className="text-gray-700 dark:text-gray-300">
                                {link.label}
                              </span>
                            </motion.a>
                          ))}
                        </div>
                      </div>

                      {/* Recent Searches */}
                      <div>
                        <h4 className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-3 flex items-center">
                          <Clock className="w-4 h-4 mr-1" />
                          Recent Searches
                        </h4>
                        <div className="space-y-1">
                          {recentSearches.map((search, index) => (
                            <motion.button
                              key={index}
                              onClick={() => handleSearch(search)}
                              className="w-full text-left p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors text-sm text-gray-600 dark:text-gray-400"
                              whileHover={{ x: 4 }}
                            >
                              {search}
                            </motion.button>
                          ))}
                        </div>
                      </div>
                    </>
                  )}

                  {query.length > 0 && query.length <= 2 && (
                    <div className="text-center py-8 text-gray-500 dark:text-gray-400">
                      <Search className="w-8 h-8 mx-auto mb-2 opacity-50" />
                      <p>Type at least 3 characters to search</p>
                    </div>
                  )}
                </div>
              )}

              {/* Search Tips */}
              <div className="bg-gray-50 dark:bg-gray-800/50 p-4 border-t border-gray-100 dark:border-gray-700">
                <div className="text-xs text-gray-500 dark:text-gray-400">
                  <strong>Search tips:</strong> Try transaction hashes, block numbers, addresses, 
                  or keywords like "lending", "festival", "governance"
                </div>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </motion.div>
    </div>
  )
}