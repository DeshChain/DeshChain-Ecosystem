import React, { useState, useMemo } from 'react';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  Badge,
  Button,
  Input,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  Textarea,
  Alert,
  AlertDescription,
} from '@/components/ui';
import {
  MapPin,
  Clock,
  Phone,
  MessageCircle,
  CheckCircle,
  XCircle,
  Eye,
  Filter,
  Search,
  Star,
  IndianRupee,
  Calendar,
  User,
} from 'lucide-react';
import { format } from 'date-fns';

interface ServiceRequestsTableProps {
  serviceRequests: ServiceRequest[];
  onAcceptRequest: (requestId: string) => Promise<void>;
  onCompleteRequest: (requestId: string, note: string) => Promise<void>;
  onRejectRequest: (requestId: string, reason: string) => Promise<void>;
  onRefresh: () => void;
}

interface ServiceRequest {
  id: string;
  customerName: string;
  customerAddress: string;
  serviceType: string;
  amount: string;
  status: string;
  priority: string;
  createdAt: string;
  acceptedAt?: string;
  estimatedCompletionTime?: string;
  location: {
    address: string;
    city: string;
    state: string;
    pincode: string;
  };
  customerRating?: number;
  customerFeedback?: string;
  earningsAmount?: string;
}

export const ServiceRequestsTable: React.FC<ServiceRequestsTableProps> = ({
  serviceRequests,
  onAcceptRequest,
  onCompleteRequest,
  onRejectRequest,
  onRefresh,
}) => {
  const [selectedRequest, setSelectedRequest] = useState<ServiceRequest | null>(null);
  const [showDetailsDialog, setShowDetailsDialog] = useState(false);
  const [showCompleteDialog, setShowCompleteDialog] = useState(false);
  const [showRejectDialog, setShowRejectDialog] = useState(false);
  const [completionNote, setCompletionNote] = useState('');
  const [rejectionReason, setRejectionReason] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [priorityFilter, setPriorityFilter] = useState('all');
  const [serviceTypeFilter, setServiceTypeFilter] = useState('all');
  const [loading, setLoading] = useState<string | null>(null);

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending': return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      case 'accepted': return 'bg-blue-100 text-blue-800 border-blue-200';
      case 'in_progress': return 'bg-orange-100 text-orange-800 border-orange-200';
      case 'completed': return 'bg-green-100 text-green-800 border-green-200';
      case 'cancelled': return 'bg-red-100 text-red-800 border-red-200';
      case 'disputed': return 'bg-purple-100 text-purple-800 border-purple-200';
      default: return 'bg-gray-100 text-gray-800 border-gray-200';
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

  const getServiceTypeIcon = (serviceType: string) => {
    switch (serviceType.toLowerCase()) {
      case 'cash_in': return '💳';
      case 'cash_out': return '💰';
      case 'money_transfer': return '💸';
      case 'bill_payment': return '🧾';
      default: return '🔄';
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

  const filteredRequests = useMemo(() => {
    return serviceRequests.filter((request) => {
      const matchesSearch = 
        request.customerName.toLowerCase().includes(searchQuery.toLowerCase()) ||
        request.serviceType.toLowerCase().includes(searchQuery.toLowerCase()) ||
        request.location.city.toLowerCase().includes(searchQuery.toLowerCase()) ||
        request.location.pincode.includes(searchQuery);
      
      const matchesStatus = statusFilter === 'all' || request.status === statusFilter;
      const matchesPriority = priorityFilter === 'all' || request.priority === priorityFilter;
      const matchesServiceType = serviceTypeFilter === 'all' || request.serviceType === serviceTypeFilter;
      
      return matchesSearch && matchesStatus && matchesPriority && matchesServiceType;
    });
  }, [serviceRequests, searchQuery, statusFilter, priorityFilter, serviceTypeFilter]);

  const handleAcceptRequest = async (requestId: string) => {
    setLoading(requestId);
    try {
      await onAcceptRequest(requestId);
      onRefresh();
    } catch (error) {
      console.error('Error accepting request:', error);
    } finally {
      setLoading(null);
    }
  };

  const handleCompleteRequest = async () => {
    if (!selectedRequest) return;
    
    setLoading(selectedRequest.id);
    try {
      await onCompleteRequest(selectedRequest.id, completionNote);
      setShowCompleteDialog(false);
      setCompletionNote('');
      setSelectedRequest(null);
      onRefresh();
    } catch (error) {
      console.error('Error completing request:', error);
    } finally {
      setLoading(null);
    }
  };

  const handleRejectRequest = async () => {
    if (!selectedRequest) return;
    
    setLoading(selectedRequest.id);
    try {
      await onRejectRequest(selectedRequest.id, rejectionReason);
      setShowRejectDialog(false);
      setRejectionReason('');
      setSelectedRequest(null);
      onRefresh();
    } catch (error) {
      console.error('Error rejecting request:', error);
    } finally {
      setLoading(null);
    }
  };

  const getActionButtons = (request: ServiceRequest) => {
    switch (request.status.toLowerCase()) {
      case 'pending':
        return (
          <div className="flex space-x-2">
            <Button
              size="sm"
              onClick={() => handleAcceptRequest(request.id)}
              disabled={loading === request.id}
              className="bg-green-600 hover:bg-green-700"
            >
              {loading === request.id ? 'Processing...' : 'Accept'}
            </Button>
            <Button
              size="sm"
              variant="outline"
              onClick={() => {
                setSelectedRequest(request);
                setShowRejectDialog(true);
              }}
              className="border-red-200 text-red-600 hover:bg-red-50"
            >
              Reject
            </Button>
          </div>
        );
      case 'accepted':
      case 'in_progress':
        return (
          <div className="flex space-x-2">
            <Button
              size="sm"
              onClick={() => {
                setSelectedRequest(request);
                setShowCompleteDialog(true);
              }}
              disabled={loading === request.id}
              className="bg-blue-600 hover:bg-blue-700"
            >
              Complete
            </Button>
            <Button
              size="sm"
              variant="outline"
              onClick={() => {
                setSelectedRequest(request);
                setShowDetailsDialog(true);
              }}
            >
              <MessageCircle className="h-4 w-4" />
            </Button>
          </div>
        );
      case 'completed':
        return (
          <Badge className="bg-green-100 text-green-800">
            ✅ Completed
          </Badge>
        );
      default:
        return (
          <Button
            size="sm"
            variant="outline"
            onClick={() => {
              setSelectedRequest(request);
              setShowDetailsDialog(true);
            }}
          >
            <Eye className="h-4 w-4" />
          </Button>
        );
    }
  };

  const uniqueServiceTypes = [...new Set(serviceRequests.map(r => r.serviceType))];

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            <span>Service Requests Management</span>
            <div className="flex items-center space-x-2">
              <Badge variant="outline">
                {filteredRequests.length} of {serviceRequests.length}
              </Badge>
              <Button onClick={onRefresh} variant="outline" size="sm">
                Refresh
              </Button>
            </div>
          </CardTitle>
        </CardHeader>
        <CardContent>
          {/* Filters */}
          <div className="flex flex-wrap gap-4 mb-6 p-4 bg-gray-50 rounded-lg">
            <div className="flex-1 min-w-[200px]">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                <Input
                  placeholder="Search by customer, service, location..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
            
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-[140px]">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="pending">Pending</SelectItem>
                <SelectItem value="accepted">Accepted</SelectItem>
                <SelectItem value="in_progress">In Progress</SelectItem>
                <SelectItem value="completed">Completed</SelectItem>
                <SelectItem value="cancelled">Cancelled</SelectItem>
              </SelectContent>
            </Select>

            <Select value={priorityFilter} onValueChange={setPriorityFilter}>
              <SelectTrigger className="w-[140px]">
                <SelectValue placeholder="Priority" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Priority</SelectItem>
                <SelectItem value="urgent">Urgent</SelectItem>
                <SelectItem value="high">High</SelectItem>
                <SelectItem value="normal">Normal</SelectItem>
                <SelectItem value="low">Low</SelectItem>
              </SelectContent>
            </Select>

            <Select value={serviceTypeFilter} onValueChange={setServiceTypeFilter}>
              <SelectTrigger className="w-[160px]">
                <SelectValue placeholder="Service Type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Services</SelectItem>
                {uniqueServiceTypes.map((type) => (
                  <SelectItem key={type} value={type}>
                    {getServiceTypeIcon(type)} {type.replace('_', ' ').toUpperCase()}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {/* Table */}
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Customer</TableHead>
                  <TableHead>Service</TableHead>
                  <TableHead>Amount</TableHead>
                  <TableHead>Location</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Priority</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredRequests.map((request) => (
                  <TableRow
                    key={request.id}
                    className="hover:bg-gray-50 cursor-pointer"
                    onClick={() => {
                      setSelectedRequest(request);
                      setShowDetailsDialog(true);
                    }}
                  >
                    <TableCell>
                      <div className="flex items-center space-x-3">
                        <div className="w-8 h-8 bg-gradient-to-r from-orange-400 to-orange-600 rounded-full flex items-center justify-center text-white text-sm font-bold">
                          {request.customerName.charAt(0)}
                        </div>
                        <div>
                          <p className="font-medium">{request.customerName}</p>
                          <p className="text-sm text-gray-500">#{request.id.slice(-6)}</p>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center space-x-2">
                        <span className="text-lg">{getServiceTypeIcon(request.serviceType)}</span>
                        <span className="font-medium">
                          {request.serviceType.replace('_', ' ').toUpperCase()}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="font-semibold text-green-600">
                        {formatCurrency(request.amount)}
                      </div>
                      {request.earningsAmount && (
                        <div className="text-sm text-gray-500">
                          Earn: {formatCurrency(request.earningsAmount)}
                        </div>
                      )}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center space-x-1">
                        <MapPin className="h-4 w-4 text-gray-400" />
                        <div>
                          <p className="text-sm">{request.location.city}</p>
                          <p className="text-xs text-gray-500">{request.location.pincode}</p>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge className={getStatusColor(request.status)}>
                        {request.status.replace('_', ' ').toUpperCase()}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center space-x-2">
                        <div
                          className={`w-3 h-3 rounded-full ${getPriorityColor(request.priority)}`}
                        ></div>
                        <span className="text-sm capitalize">{request.priority}</span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center space-x-1">
                        <Calendar className="h-4 w-4 text-gray-400" />
                        <div>
                          <p className="text-sm">
                            {format(new Date(request.createdAt), 'MMM dd')}
                          </p>
                          <p className="text-xs text-gray-500">
                            {format(new Date(request.createdAt), 'HH:mm')}
                          </p>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell onClick={(e) => e.stopPropagation()}>
                      {getActionButtons(request)}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>

            {filteredRequests.length === 0 && (
              <div className="text-center py-12">
                <Filter className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <p className="text-gray-500">No service requests match your filters</p>
                <Button
                  variant="outline"
                  className="mt-4"
                  onClick={() => {
                    setSearchQuery('');
                    setStatusFilter('all');
                    setPriorityFilter('all');
                    setServiceTypeFilter('all');
                  }}
                >
                  Clear Filters
                </Button>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Request Details Dialog */}
      <Dialog open={showDetailsDialog} onOpenChange={setShowDetailsDialog}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Service Request Details</DialogTitle>
          </DialogHeader>
          {selectedRequest && (
            <div className="space-y-6">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <h4 className="font-semibold mb-2">Customer Information</h4>
                  <div className="space-y-2">
                    <p><User className="inline h-4 w-4 mr-2" />{selectedRequest.customerName}</p>
                    <p><MapPin className="inline h-4 w-4 mr-2" />{selectedRequest.location.address}</p>
                    <p className="text-sm text-gray-600">{selectedRequest.location.city}, {selectedRequest.location.state} - {selectedRequest.location.pincode}</p>
                  </div>
                </div>
                <div>
                  <h4 className="font-semibold mb-2">Service Details</h4>
                  <div className="space-y-2">
                    <p>{getServiceTypeIcon(selectedRequest.serviceType)} {selectedRequest.serviceType.replace('_', ' ').toUpperCase()}</p>
                    <p><IndianRupee className="inline h-4 w-4 mr-2" />{formatCurrency(selectedRequest.amount)}</p>
                    <Badge className={getStatusColor(selectedRequest.status)}>
                      {selectedRequest.status.replace('_', ' ').toUpperCase()}
                    </Badge>
                  </div>
                </div>
              </div>

              <div>
                <h4 className="font-semibold mb-2">Timeline</h4>
                <div className="space-y-2 text-sm">
                  <p><Clock className="inline h-4 w-4 mr-2" />Created: {format(new Date(selectedRequest.createdAt), 'PPpp')}</p>
                  {selectedRequest.acceptedAt && (
                    <p><CheckCircle className="inline h-4 w-4 mr-2" />Accepted: {format(new Date(selectedRequest.acceptedAt), 'PPpp')}</p>
                  )}
                  {selectedRequest.estimatedCompletionTime && (
                    <p><Clock className="inline h-4 w-4 mr-2" />Est. Completion: {format(new Date(selectedRequest.estimatedCompletionTime), 'PPpp')}</p>
                  )}
                </div>
              </div>

              {selectedRequest.customerRating && (
                <div>
                  <h4 className="font-semibold mb-2">Customer Feedback</h4>
                  <div className="flex items-center space-x-2 mb-2">
                    <Star className="h-4 w-4 text-yellow-500 fill-current" />
                    <span>{selectedRequest.customerRating}/5</span>
                  </div>
                  {selectedRequest.customerFeedback && (
                    <p className="text-sm text-gray-600 bg-gray-50 p-3 rounded">
                      "{selectedRequest.customerFeedback}"
                    </p>
                  )}
                </div>
              )}
            </div>
          )}
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowDetailsDialog(false)}>
              Close
            </Button>
            <Button onClick={() => window.open(`tel:${selectedRequest?.customerAddress}`)}>
              <Phone className="h-4 w-4 mr-2" />
              Call Customer
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Complete Request Dialog */}
      <Dialog open={showCompleteDialog} onOpenChange={setShowCompleteDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Complete Service Request</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <Alert>
              <CheckCircle className="h-4 w-4" />
              <AlertDescription>
                You are about to mark this service as completed. Please provide any completion notes.
              </AlertDescription>
            </Alert>
            <Textarea
              placeholder="Enter completion notes (optional)"
              value={completionNote}
              onChange={(e) => setCompletionNote(e.target.value)}
              rows={3}
            />
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowCompleteDialog(false)}>
              Cancel
            </Button>
            <Button
              onClick={handleCompleteRequest}
              disabled={loading === selectedRequest?.id}
              className="bg-green-600 hover:bg-green-700"
            >
              {loading === selectedRequest?.id ? 'Processing...' : 'Complete Service'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Reject Request Dialog */}
      <Dialog open={showRejectDialog} onOpenChange={setShowRejectDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Reject Service Request</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <Alert>
              <XCircle className="h-4 w-4" />
              <AlertDescription>
                You are about to reject this service request. Please provide a reason.
              </AlertDescription>
            </Alert>
            <Textarea
              placeholder="Enter rejection reason"
              value={rejectionReason}
              onChange={(e) => setRejectionReason(e.target.value)}
              rows={3}
              required
            />
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setShowRejectDialog(false)}>
              Cancel
            </Button>
            <Button
              onClick={handleRejectRequest}
              disabled={!rejectionReason.trim() || loading === selectedRequest?.id}
              variant="destructive"
            >
              {loading === selectedRequest?.id ? 'Processing...' : 'Reject Request'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};