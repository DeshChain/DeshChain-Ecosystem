import React, { useEffect, useState } from 'react';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  Badge,
  Button,
  Progress,
  Separator,
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui';
import {
  Activity,
  Users,
  IndianRupee,
  Star,
  Clock,
  TrendingUp,
  Bell,
  MapPin,
  Phone,
  MessageCircle,
} from 'lucide-react';
import { LineChart, Line, AreaChart, Area, PieChart, Pie, Cell, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

interface DashboardOverviewProps {
  dashboardData: SevaMitraDashboardData;
  onRefresh: () => void;
}

interface SevaMitraDashboardData {
  mitraInfo: SevaMitra;
  summary: DashboardSummary;
  earningsData: EarningsData;
  serviceRequests: ServiceRequest[];
  performanceStats: PerformanceStats;
  analytics: AnalyticsData;
  notifications: Notification[];
  rankings: RankingData;
}

interface DashboardSummary {
  totalServices: number;
  todayServices: number;
  monthlyServices: number;
  totalEarnings: string;
  todayEarnings: string;
  monthlyEarnings: string;
  pendingRequests: number;
  trustScore: number;
  responseTime: string;
  customerRating: number;
  onlineStatus: boolean;
  lastActiveTime: string;
}

interface ServiceRequest {
  id: string;
  customerName: string;
  serviceType: string;
  amount: string;
  status: string;
  priority: string;
  createdAt: string;
  location: {
    address: string;
    city: string;
    pincode: string;
  };
}

interface Notification {
  id: string;
  type: string;
  title: string;
  message: string;
  priority: string;
  createdAt: string;
  data?: any;
}

const COLORS = ['#FF9933', '#138808', '#000080', '#FFD700', '#FF6B6B'];

export const DashboardOverview: React.FC<DashboardOverviewProps> = ({
  dashboardData,
  onRefresh,
}) => {
  const { mitraInfo, summary, earningsData, serviceRequests, notifications, rankings, analytics } = dashboardData;

  const getTrustScoreBadge = (score: number) => {
    if (score >= 90) return { label: 'Diamond', color: 'bg-purple-500', textColor: 'text-white' };
    if (score >= 80) return { label: 'Platinum', color: 'bg-gray-500', textColor: 'text-white' };
    if (score >= 70) return { label: 'Gold', color: 'bg-yellow-500', textColor: 'text-white' };
    if (score >= 60) return { label: 'Silver', color: 'bg-gray-400', textColor: 'text-white' };
    if (score >= 50) return { label: 'Bronze', color: 'bg-amber-600', textColor: 'text-white' };
    return { label: 'New', color: 'bg-red-500', textColor: 'text-white' };
  };

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending': return 'bg-yellow-100 text-yellow-800';
      case 'accepted': return 'bg-blue-100 text-blue-800';
      case 'in_progress': return 'bg-orange-100 text-orange-800';
      case 'completed': return 'bg-green-100 text-green-800';
      case 'cancelled': return 'bg-red-100 text-red-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  const getPriorityColor = (priority: string) => {
    switch (priority.toLowerCase()) {
      case 'urgent': return 'bg-red-500';
      case 'high': return 'bg-orange-500';
      case 'normal': return 'bg-blue-500';
      case 'low': return 'bg-gray-500';
      default: return 'bg-gray-500';
    }
  };

  const formatCurrency = (amount: string) => {
    const num = parseFloat(amount);
    return new Intl.NumberFormat('en-IN', {
      style: 'currency',
      currency: 'INR',
      maximumFractionDigits: 0,
    }).format(num);
  };

  const trustBadge = getTrustScoreBadge(summary.trustScore);

  return (
    <div className="min-h-screen bg-gradient-to-br from-orange-50 to-green-50 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <div className="relative">
              <div className="w-16 h-16 bg-gradient-to-r from-orange-400 to-orange-600 rounded-full flex items-center justify-center text-white text-xl font-bold">
                {mitraInfo.name.charAt(0)}
              </div>
              {summary.onlineStatus && (
                <div className="absolute -bottom-1 -right-1 w-5 h-5 bg-green-500 border-2 border-white rounded-full"></div>
              )}
            </div>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">
                नमस्ते, {mitraInfo.name}! 🙏
              </h1>
              <p className="text-gray-600">Seva Mitra Dashboard</p>
              <div className="flex items-center space-x-2 mt-1">
                <Badge className={`${trustBadge.color} ${trustBadge.textColor}`}>
                  {trustBadge.label} {summary.trustScore}
                </Badge>
                <Badge variant="outline">
                  Rank #{rankings.localRank} Local
                </Badge>
              </div>
            </div>
          </div>
          <div className="flex items-center space-x-3">
            <Button
              variant="outline"
              size="sm"
              onClick={onRefresh}
              className="border-orange-200 hover:bg-orange-50"
            >
              Refresh
            </Button>
            <Badge variant={summary.onlineStatus ? 'default' : 'secondary'}>
              {summary.onlineStatus ? 'Online' : 'Offline'}
            </Badge>
          </div>
        </div>

        {/* Key Metrics Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <Card className="border-l-4 border-l-orange-500">
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">Today's Services</p>
                  <p className="text-2xl font-bold text-gray-900">{summary.todayServices}</p>
                  <p className="text-xs text-gray-500">Total: {summary.totalServices}</p>
                </div>
                <Activity className="h-8 w-8 text-orange-600" />
              </div>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-green-500">
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">Today's Earnings</p>
                  <p className="text-2xl font-bold text-gray-900">
                    {formatCurrency(summary.todayEarnings)}
                  </p>
                  <p className="text-xs text-gray-500">
                    Total: {formatCurrency(summary.totalEarnings)}
                  </p>
                </div>
                <IndianRupee className="h-8 w-8 text-green-600" />
              </div>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-blue-500">
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">Customer Rating</p>
                  <div className="flex items-center space-x-1">
                    <p className="text-2xl font-bold text-gray-900">{summary.customerRating}</p>
                    <Star className="h-5 w-5 text-yellow-500 fill-current" />
                  </div>
                  <p className="text-xs text-gray-500">Trust Score: {summary.trustScore}</p>
                </div>
                <Users className="h-8 w-8 text-blue-600" />
              </div>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-l-purple-500">
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">Pending Requests</p>
                  <p className="text-2xl font-bold text-gray-900">{summary.pendingRequests}</p>
                  <p className="text-xs text-gray-500">Avg Response: {summary.responseTime}</p>
                </div>
                <Clock className="h-8 w-8 text-purple-600" />
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Main Dashboard Content */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left Column - Service Requests & Notifications */}
          <div className="lg:col-span-2 space-y-6">
            {/* Recent Service Requests */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <span>Recent Service Requests</span>
                  <Badge variant="outline">{serviceRequests.length} requests</Badge>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {serviceRequests.slice(0, 5).map((request) => (
                    <div
                      key={request.id}
                      className="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
                    >
                      <div className="flex items-center space-x-4">
                        <div className={`w-3 h-3 rounded-full ${getPriorityColor(request.priority)}`}></div>
                        <div>
                          <p className="font-medium text-gray-900">{request.customerName}</p>
                          <p className="text-sm text-gray-600">{request.serviceType}</p>
                          <div className="flex items-center space-x-2 mt-1">
                            <MapPin className="h-3 w-3 text-gray-400" />
                            <span className="text-xs text-gray-500">
                              {request.location.city}, {request.location.pincode}
                            </span>
                          </div>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="font-semibold text-gray-900">
                          {formatCurrency(request.amount)}
                        </p>
                        <Badge className={getStatusColor(request.status)}>
                          {request.status.replace('_', ' ')}
                        </Badge>
                        <p className="text-xs text-gray-500 mt-1">
                          {new Date(request.createdAt).toLocaleDateString()}
                        </p>
                      </div>
                    </div>
                  ))}
                  {serviceRequests.length === 0 && (
                    <div className="text-center py-8 text-gray-500">
                      <Activity className="h-12 w-12 mx-auto mb-4 opacity-50" />
                      <p>No recent service requests</p>
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>

            {/* Service Type Analytics */}
            <Card>
              <CardHeader>
                <CardTitle>Service Type Distribution</CardTitle>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                  <PieChart>
                    <Pie
                      data={analytics.serviceTypesChart}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      label={({ name, percentage }) => `${name} ${percentage}%`}
                      outerRadius={80}
                      fill="#8884d8"
                      dataKey="count"
                    >
                      {analytics.serviceTypesChart.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                      ))}
                    </Pie>
                    <Tooltip />
                  </PieChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>
          </div>

          {/* Right Column - Notifications & Quick Actions */}
          <div className="space-y-6">
            {/* Quick Actions */}
            <Card>
              <CardHeader>
                <CardTitle>Quick Actions</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <Button 
                  className="w-full justify-start bg-orange-600 hover:bg-orange-700" 
                  size="lg"
                >
                  <Activity className="h-4 w-4 mr-2" />
                  Accept New Request
                </Button>
                <Button 
                  variant="outline" 
                  className="w-full justify-start border-green-200 hover:bg-green-50" 
                  size="lg"
                >
                  <Phone className="h-4 w-4 mr-2" />
                  Update Availability
                </Button>
                <Button 
                  variant="outline" 
                  className="w-full justify-start border-blue-200 hover:bg-blue-50" 
                  size="lg"
                >
                  <MessageCircle className="h-4 w-4 mr-2" />
                  Customer Support
                </Button>
              </CardContent>
            </Card>

            {/* Notifications */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <span>Notifications</span>
                  <Badge variant="outline">{notifications.length}</Badge>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {notifications.slice(0, 4).map((notification) => (
                    <div
                      key={notification.id}
                      className="flex items-start space-x-3 p-3 bg-gray-50 rounded-lg"
                    >
                      <Bell className="h-4 w-4 text-orange-600 mt-1 flex-shrink-0" />
                      <div className="flex-1 min-w-0">
                        <p className="font-medium text-gray-900 text-sm">
                          {notification.title}
                        </p>
                        <p className="text-xs text-gray-600">
                          {notification.message}
                        </p>
                        <p className="text-xs text-gray-500 mt-1">
                          {new Date(notification.createdAt).toLocaleTimeString()}
                        </p>
                      </div>
                    </div>
                  ))}
                  {notifications.length === 0 && (
                    <div className="text-center py-4 text-gray-500">
                      <Bell className="h-8 w-8 mx-auto mb-2 opacity-50" />
                      <p className="text-sm">No notifications</p>
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>

            {/* Performance Summary */}
            <Card>
              <CardHeader>
                <CardTitle>Performance Summary</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <div className="flex justify-between text-sm mb-1">
                    <span>Completion Rate</span>
                    <span>{dashboardData.performanceStats.completionRate}%</span>
                  </div>
                  <Progress value={dashboardData.performanceStats.completionRate} className="h-2" />
                </div>
                <div>
                  <div className="flex justify-between text-sm mb-1">
                    <span>Customer Satisfaction</span>
                    <span>{dashboardData.performanceStats.customerSatisfaction}/5</span>
                  </div>
                  <Progress value={dashboardData.performanceStats.customerSatisfaction * 20} className="h-2" />
                </div>
                <div>
                  <div className="flex justify-between text-sm mb-1">
                    <span>On-time Delivery</span>
                    <span>{dashboardData.performanceStats.onTimeDelivery}%</span>
                  </div>
                  <Progress value={dashboardData.performanceStats.onTimeDelivery} className="h-2" />
                </div>
                <Separator />
                <div className="text-center">
                  <p className="text-2xl font-bold text-green-600">
                    +{dashboardData.performanceStats.monthlyGrowth}%
                  </p>
                  <p className="text-sm text-gray-600">Monthly Growth</p>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Earnings Chart */}
        <Card>
          <CardHeader>
            <CardTitle>Weekly Earnings Trend</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <AreaChart data={earningsData.weeklyEarnings}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis 
                  dataKey="date" 
                  tickFormatter={(value) => new Date(value).toLocaleDateString('en-IN', { day: 'numeric', month: 'short' })}
                />
                <YAxis tickFormatter={(value) => `₹${value}`} />
                <Tooltip 
                  labelFormatter={(value) => new Date(value).toLocaleDateString('en-IN')}
                  formatter={(value) => [`₹${value}`, 'Earnings']}
                />
                <Area type="monotone" dataKey="amount" stroke="#FF9933" fill="#FF9933" fillOpacity={0.6} />
              </AreaChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};