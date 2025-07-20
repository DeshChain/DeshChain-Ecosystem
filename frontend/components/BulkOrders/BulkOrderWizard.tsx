import React, { useState, useCallback, useMemo } from 'react';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  Button,
  Input,
  Label,
  Textarea,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Badge,
  Progress,
  Alert,
  AlertDescription,
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui';
import {
  Upload,
  Download,
  FileText,
  CheckCircle,
  XCircle,
  AlertTriangle,
  Users,
  IndianRupee,
  Clock,
  Send,
  Eye,
  Trash2,
  Plus,
  FileSpreadsheet,
} from 'lucide-react';
import { useDropzone } from 'react-dropzone';
import Papa from 'papaparse';

interface BulkOrderWizardProps {
  businessAccount: BusinessAccount;
  onCreateBulkOrder: (orders: BulkOrderData) => Promise<void>;
  onValidateTemplate: (template: BulkOrderTemplate) => Promise<ValidationResult>;
}

interface BusinessAccount {
  address: string;
  businessName: string;
  dailyLimit: string;
  monthlyLimit: string;
  maxBulkOrderSize: number;
  isActive: boolean;
  bulkOrdersEnabled: boolean;
}

interface BulkOrderData {
  orders: OrderItem[];
  metadata: OrderMetadata;
  settings: ProcessingSettings;
}

interface OrderItem {
  recipientAddress: string;
  amount: string;
  memo?: string;
  priority: 'LOW' | 'NORMAL' | 'HIGH' | 'URGENT';
  customerRef?: string;
}

interface OrderMetadata {
  description: string;
  reference: string;
  department?: string;
  projectCode?: string;
  notifyEmail?: string;
  scheduledTime?: Date;
}

interface ProcessingSettings {
  batchSize: number;
  maxRetries: number;
  stopOnFirstFailure: boolean;
  validateRecipients: boolean;
  notifyOnCompletion: boolean;
}

interface ValidationResult {
  isValid: boolean;
  totalOrders: number;
  validOrders: number;
  invalidOrders: number;
  totalAmount: string;
  warnings: string[];
  errors: string[];
}

const STEP_TITLES = [
  'Upload Orders',
  'Review & Validate',
  'Configure Settings',
  'Confirm & Submit'
];

export const BulkOrderWizard: React.FC<BulkOrderWizardProps> = ({
  businessAccount,
  onCreateBulkOrder,
  onValidateTemplate,
}) => {
  const [currentStep, setCurrentStep] = useState(0);
  const [orders, setOrders] = useState<OrderItem[]>([]);
  const [metadata, setMetadata] = useState<OrderMetadata>({
    description: '',
    reference: '',
  });
  const [settings, setSettings] = useState<ProcessingSettings>({
    batchSize: 50,
    maxRetries: 3,
    stopOnFirstFailure: false,
    validateRecipients: true,
    notifyOnCompletion: true,
  });
  const [validation, setValidation] = useState<ValidationResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [showPreview, setShowPreview] = useState(false);
  const [uploadMethod, setUploadMethod] = useState<'csv' | 'manual'>('csv');
  const [csvData, setCsvData] = useState<string>('');
  const [manualEntry, setManualEntry] = useState<OrderItem>({
    recipientAddress: '',
    amount: '',
    memo: '',
    priority: 'NORMAL',
  });

  const onDrop = useCallback((acceptedFiles: File[]) => {
    const file = acceptedFiles[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (e) => {
        const csv = e.target?.result as string;
        setCsvData(csv);
        parseCSV(csv);
      };
      reader.readAsText(file);
    }
  }, []);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      'text/csv': ['.csv'],
      'application/vnd.ms-excel': ['.xls'],
      'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet': ['.xlsx'],
    },
    multiple: false,
  });

  const parseCSV = (csvString: string) => {
    Papa.parse(csvString, {
      header: true,
      skipEmptyLines: true,
      complete: (results) => {
        const parsedOrders: OrderItem[] = results.data.map((row: any, index) => ({
          recipientAddress: row.recipient_address || row.address || '',
          amount: row.amount || '',
          memo: row.memo || row.description || '',
          priority: (row.priority || 'NORMAL').toUpperCase() as 'LOW' | 'NORMAL' | 'HIGH' | 'URGENT',
          customerRef: row.customer_ref || row.reference || `${Date.now()}-${index}`,
        })).filter(order => order.recipientAddress && order.amount);
        
        setOrders(parsedOrders);
      },
      error: (error) => {
        console.error('CSV parsing error:', error);
      },
    });
  };

  const addManualOrder = () => {
    if (manualEntry.recipientAddress && manualEntry.amount) {
      setOrders([...orders, { ...manualEntry, customerRef: `${Date.now()}` }]);
      setManualEntry({
        recipientAddress: '',
        amount: '',
        memo: '',
        priority: 'NORMAL',
      });
    }
  };

  const removeOrder = (index: number) => {
    setOrders(orders.filter((_, i) => i !== index));
  };

  const validateOrders = async () => {
    setLoading(true);
    try {
      const template = {
        orders: orders.map(order => ({
          ...order,
          amount: parseFloat(order.amount),
        })),
        metadata,
        settings,
      };
      
      const result = await onValidateTemplate(template);
      setValidation(result);
    } catch (error) {
      console.error('Validation error:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async () => {
    setLoading(true);
    try {
      await onCreateBulkOrder({
        orders,
        metadata,
        settings,
      });
    } catch (error) {
      console.error('Submission error:', error);
    } finally {
      setLoading(false);
    }
  };

  const totalAmount = useMemo(() => {
    return orders.reduce((sum, order) => sum + parseFloat(order.amount || '0'), 0);
  }, [orders]);

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-IN', {
      style: 'currency',
      currency: 'INR',
      maximumFractionDigits: 0,
    }).format(amount);
  };

  const downloadTemplate = () => {
    const csvContent = [
      'recipient_address,amount,memo,priority,customer_ref',
      'desh1abc123...,1000,Salary payment,NORMAL,EMP001',
      'desh1def456...,2500,Vendor payment,HIGH,VND002',
      'desh1ghi789...,500,Refund,LOW,REF003',
    ].join('\n');

    const blob = new Blob([csvContent], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'bulk_order_template.csv';
    a.click();
    window.URL.revokeObjectURL(url);
  };

  const canProceed = () => {
    switch (currentStep) {
      case 0: return orders.length > 0;
      case 1: return validation?.isValid || false;
      case 2: return true;
      case 3: return true;
      default: return false;
    }
  };

  const nextStep = () => {
    if (currentStep === 1 && !validation) {
      validateOrders();
      return;
    }
    if (currentStep < STEP_TITLES.length - 1) {
      setCurrentStep(currentStep + 1);
    }
  };

  const prevStep = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
    }
  };

  const renderStepContent = () => {
    switch (currentStep) {
      case 0:
        return (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center justify-between">
                <span>Upload Money Orders</span>
                <Button onClick={downloadTemplate} variant="outline" size="sm">
                  <Download className="h-4 w-4 mr-2" />
                  Download Template
                </Button>
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-6">
              <Tabs value={uploadMethod} onValueChange={(value) => setUploadMethod(value as 'csv' | 'manual')}>
                <TabsList className="grid w-full grid-cols-2">
                  <TabsTrigger value="csv">CSV Upload</TabsTrigger>
                  <TabsTrigger value="manual">Manual Entry</TabsTrigger>
                </TabsList>
                
                <TabsContent value="csv" className="space-y-4">
                  <div
                    {...getRootProps()}
                    className={`border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors ${
                      isDragActive ? 'border-orange-500 bg-orange-50' : 'border-gray-300 hover:border-gray-400'
                    }`}
                  >
                    <input {...getInputProps()} />
                    <Upload className="h-12 w-12 mx-auto mb-4 text-gray-400" />
                    {isDragActive ? (
                      <p className="text-orange-600">Drop the CSV file here...</p>
                    ) : (
                      <div>
                        <p className="text-lg font-medium mb-2">Drop CSV file here or click to browse</p>
                        <p className="text-gray-500">Supports CSV, XLS, XLSX files</p>
                      </div>
                    )}
                  </div>
                  
                  {csvData && (
                    <div className="bg-green-50 p-4 rounded-lg">
                      <div className="flex items-center space-x-2">
                        <CheckCircle className="h-5 w-5 text-green-600" />
                        <span className="text-green-800">CSV file uploaded successfully</span>
                      </div>
                      <p className="text-sm text-green-700 mt-1">
                        Parsed {orders.length} orders from CSV
                      </p>
                    </div>
                  )}
                </TabsContent>
                
                <TabsContent value="manual" className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <Label htmlFor="recipient">Recipient Address</Label>
                      <Input
                        id="recipient"
                        placeholder="desh1..."
                        value={manualEntry.recipientAddress}
                        onChange={(e) => setManualEntry({...manualEntry, recipientAddress: e.target.value})}
                      />
                    </div>
                    <div>
                      <Label htmlFor="amount">Amount (INR)</Label>
                      <Input
                        id="amount"
                        type="number"
                        placeholder="1000"
                        value={manualEntry.amount}
                        onChange={(e) => setManualEntry({...manualEntry, amount: e.target.value})}
                      />
                    </div>
                    <div>
                      <Label htmlFor="memo">Memo (Optional)</Label>
                      <Input
                        id="memo"
                        placeholder="Payment description"
                        value={manualEntry.memo}
                        onChange={(e) => setManualEntry({...manualEntry, memo: e.target.value})}
                      />
                    </div>
                    <div>
                      <Label htmlFor="priority">Priority</Label>
                      <Select value={manualEntry.priority} onValueChange={(value) => setManualEntry({...manualEntry, priority: value as any})}>
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="LOW">Low</SelectItem>
                          <SelectItem value="NORMAL">Normal</SelectItem>
                          <SelectItem value="HIGH">High</SelectItem>
                          <SelectItem value="URGENT">Urgent</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                  </div>
                  <Button onClick={addManualOrder} className="w-full">
                    <Plus className="h-4 w-4 mr-2" />
                    Add Order
                  </Button>
                </TabsContent>
              </Tabs>

              {orders.length > 0 && (
                <div>
                  <div className="flex items-center justify-between mb-4">
                    <h3 className="text-lg font-semibold">Orders Preview</h3>
                    <div className="flex items-center space-x-4">
                      <Badge variant="outline">
                        {orders.length} orders
                      </Badge>
                      <Badge variant="outline">
                        Total: {formatCurrency(totalAmount)}
                      </Badge>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setShowPreview(true)}
                      >
                        <Eye className="h-4 w-4 mr-2" />
                        View All
                      </Button>
                    </div>
                  </div>
                  
                  <div className="bg-gray-50 p-4 rounded-lg">
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
                      <div>
                        <p className="text-2xl font-bold text-orange-600">{orders.length}</p>
                        <p className="text-sm text-gray-600">Total Orders</p>
                      </div>
                      <div>
                        <p className="text-2xl font-bold text-green-600">{formatCurrency(totalAmount)}</p>
                        <p className="text-sm text-gray-600">Total Amount</p>
                      </div>
                      <div>
                        <p className="text-2xl font-bold text-blue-600">
                          {orders.filter(o => o.priority === 'HIGH' || o.priority === 'URGENT').length}
                        </p>
                        <p className="text-sm text-gray-600">High Priority</p>
                      </div>
                      <div>
                        <p className="text-2xl font-bold text-purple-600">
                          {formatCurrency(totalAmount / orders.length)}
                        </p>
                        <p className="text-sm text-gray-600">Avg Amount</p>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        );

      case 1:
        return (
          <Card>
            <CardHeader>
              <CardTitle>Review & Validate Orders</CardTitle>
            </CardHeader>
            <CardContent className="space-y-6">
              {!validation ? (
                <div className="text-center py-8">
                  <Button onClick={validateOrders} disabled={loading} size="lg">
                    {loading ? 'Validating...' : 'Validate Orders'}
                  </Button>
                  <p className="text-sm text-gray-600 mt-2">
                    This will check all addresses and amounts for validity
                  </p>
                </div>
              ) : (
                <div className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <Card className="border-green-200">
                      <CardContent className="p-4 text-center">
                        <CheckCircle className="h-8 w-8 text-green-600 mx-auto mb-2" />
                        <p className="text-2xl font-bold text-green-600">{validation.validOrders}</p>
                        <p className="text-sm text-gray-600">Valid Orders</p>
                      </CardContent>
                    </Card>
                    
                    <Card className="border-red-200">
                      <CardContent className="p-4 text-center">
                        <XCircle className="h-8 w-8 text-red-600 mx-auto mb-2" />
                        <p className="text-2xl font-bold text-red-600">{validation.invalidOrders}</p>
                        <p className="text-sm text-gray-600">Invalid Orders</p>
                      </CardContent>
                    </Card>
                    
                    <Card className="border-orange-200">
                      <CardContent className="p-4 text-center">
                        <IndianRupee className="h-8 w-8 text-orange-600 mx-auto mb-2" />
                        <p className="text-2xl font-bold text-orange-600">
                          {formatCurrency(parseFloat(validation.totalAmount))}
                        </p>
                        <p className="text-sm text-gray-600">Total Amount</p>
                      </CardContent>
                    </Card>
                  </div>

                  {validation.errors.length > 0 && (
                    <Alert>
                      <XCircle className="h-4 w-4" />
                      <AlertDescription>
                        <div className="space-y-1">
                          <p className="font-medium">Validation Errors:</p>
                          {validation.errors.map((error, index) => (
                            <p key={index} className="text-sm">• {error}</p>
                          ))}
                        </div>
                      </AlertDescription>
                    </Alert>
                  )}

                  {validation.warnings.length > 0 && (
                    <Alert>
                      <AlertTriangle className="h-4 w-4" />
                      <AlertDescription>
                        <div className="space-y-1">
                          <p className="font-medium">Warnings:</p>
                          {validation.warnings.map((warning, index) => (
                            <p key={index} className="text-sm">• {warning}</p>
                          ))}
                        </div>
                      </AlertDescription>
                    </Alert>
                  )}
                  
                  {validation.isValid && (
                    <div className="bg-green-50 p-4 rounded-lg">
                      <div className="flex items-center space-x-2">
                        <CheckCircle className="h-5 w-5 text-green-600" />
                        <span className="text-green-800 font-medium">All orders are valid and ready for processing!</span>
                      </div>
                    </div>
                  )}
                </div>
              )}
            </CardContent>
          </Card>
        );

      case 2:
        return (
          <Card>
            <CardHeader>
              <CardTitle>Configure Processing Settings</CardTitle>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-4">
                  <div>
                    <Label htmlFor="description">Description</Label>
                    <Input
                      id="description"
                      placeholder="Monthly salary payments"
                      value={metadata.description}
                      onChange={(e) => setMetadata({...metadata, description: e.target.value})}
                    />
                  </div>
                  
                  <div>
                    <Label htmlFor="reference">Reference Number</Label>
                    <Input
                      id="reference"
                      placeholder="PAY-2024-001"
                      value={metadata.reference}
                      onChange={(e) => setMetadata({...metadata, reference: e.target.value})}
                    />
                  </div>
                  
                  <div>
                    <Label htmlFor="department">Department (Optional)</Label>
                    <Input
                      id="department"
                      placeholder="HR Department"
                      value={metadata.department || ''}
                      onChange={(e) => setMetadata({...metadata, department: e.target.value})}
                    />
                  </div>
                  
                  <div>
                    <Label htmlFor="notifyEmail">Notification Email</Label>
                    <Input
                      id="notifyEmail"
                      type="email"
                      placeholder="admin@company.com"
                      value={metadata.notifyEmail || ''}
                      onChange={(e) => setMetadata({...metadata, notifyEmail: e.target.value})}
                    />
                  </div>
                </div>
                
                <div className="space-y-4">
                  <div>
                    <Label htmlFor="batchSize">Batch Size</Label>
                    <Select 
                      value={settings.batchSize.toString()} 
                      onValueChange={(value) => setSettings({...settings, batchSize: parseInt(value)})}
                    >
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="10">10 orders per batch</SelectItem>
                        <SelectItem value="25">25 orders per batch</SelectItem>
                        <SelectItem value="50">50 orders per batch</SelectItem>
                        <SelectItem value="100">100 orders per batch</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  
                  <div>
                    <Label htmlFor="maxRetries">Max Retries</Label>
                    <Select 
                      value={settings.maxRetries.toString()} 
                      onValueChange={(value) => setSettings({...settings, maxRetries: parseInt(value)})}
                    >
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="0">No retries</SelectItem>
                        <SelectItem value="1">1 retry</SelectItem>
                        <SelectItem value="3">3 retries</SelectItem>
                        <SelectItem value="5">5 retries</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  
                  <div className="space-y-3">
                    <div className="flex items-center space-x-2">
                      <input
                        type="checkbox"
                        id="stopOnFirstFailure"
                        checked={settings.stopOnFirstFailure}
                        onChange={(e) => setSettings({...settings, stopOnFirstFailure: e.target.checked})}
                      />
                      <Label htmlFor="stopOnFirstFailure">Stop on first failure</Label>
                    </div>
                    
                    <div className="flex items-center space-x-2">
                      <input
                        type="checkbox"
                        id="validateRecipients"
                        checked={settings.validateRecipients}
                        onChange={(e) => setSettings({...settings, validateRecipients: e.target.checked})}
                      />
                      <Label htmlFor="validateRecipients">Validate recipient addresses</Label>
                    </div>
                    
                    <div className="flex items-center space-x-2">
                      <input
                        type="checkbox"
                        id="notifyOnCompletion"
                        checked={settings.notifyOnCompletion}
                        onChange={(e) => setSettings({...settings, notifyOnCompletion: e.target.checked})}
                      />
                      <Label htmlFor="notifyOnCompletion">Send completion notification</Label>
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        );

      case 3:
        return (
          <Card>
            <CardHeader>
              <CardTitle>Confirm & Submit Bulk Order</CardTitle>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="bg-gray-50 p-6 rounded-lg">
                <h3 className="text-lg font-semibold mb-4">Order Summary</h3>
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                  <div>
                    <p className="text-sm text-gray-600">Total Orders</p>
                    <p className="text-2xl font-bold">{orders.length}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Total Amount</p>
                    <p className="text-2xl font-bold text-green-600">{formatCurrency(totalAmount)}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Estimated Fees</p>
                    <p className="text-2xl font-bold text-orange-600">{formatCurrency(totalAmount * 0.002)}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Processing Time</p>
                    <p className="text-2xl font-bold text-blue-600">~{Math.ceil(orders.length / settings.batchSize)}min</p>
                  </div>
                </div>
              </div>
              
              <div className="space-y-4">
                <div>
                  <h4 className="font-medium">Metadata</h4>
                  <div className="text-sm text-gray-600 space-y-1">
                    <p><strong>Description:</strong> {metadata.description || 'None'}</p>
                    <p><strong>Reference:</strong> {metadata.reference || 'None'}</p>
                    <p><strong>Department:</strong> {metadata.department || 'None'}</p>
                  </div>
                </div>
                
                <div>
                  <h4 className="font-medium">Processing Settings</h4>
                  <div className="text-sm text-gray-600 space-y-1">
                    <p><strong>Batch Size:</strong> {settings.batchSize} orders</p>
                    <p><strong>Max Retries:</strong> {settings.maxRetries}</p>
                    <p><strong>Stop on Failure:</strong> {settings.stopOnFirstFailure ? 'Yes' : 'No'}</p>
                    <p><strong>Validate Recipients:</strong> {settings.validateRecipients ? 'Yes' : 'No'}</p>
                  </div>
                </div>
              </div>
              
              <Alert>
                <AlertTriangle className="h-4 w-4" />
                <AlertDescription>
                  Please review all details carefully. Once submitted, this bulk order cannot be cancelled.
                </AlertDescription>
              </Alert>
              
              <Button 
                onClick={handleSubmit} 
                disabled={loading} 
                size="lg" 
                className="w-full bg-green-600 hover:bg-green-700"
              >
                {loading ? 'Submitting...' : 'Submit Bulk Order'}
                <Send className="h-4 w-4 ml-2" />
              </Button>
            </CardContent>
          </Card>
        );

      default:
        return null;
    }
  };

  return (
    <div className="max-w-6xl mx-auto p-6">
      {/* Progress Indicator */}
      <div className="mb-8">
        <div className="flex items-center justify-between mb-4">
          {STEP_TITLES.map((title, index) => (
            <div
              key={index}
              className={`flex items-center ${index < STEP_TITLES.length - 1 ? 'flex-1' : ''}`}
            >
              <div
                className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium ${
                  index <= currentStep
                    ? 'bg-orange-600 text-white'
                    : 'bg-gray-200 text-gray-600'
                }`}
              >
                {index < currentStep ? <CheckCircle className="h-5 w-5" /> : index + 1}
              </div>
              <span className={`ml-2 text-sm ${index <= currentStep ? 'text-orange-600 font-medium' : 'text-gray-500'}`}>
                {title}
              </span>
              {index < STEP_TITLES.length - 1 && (
                <div className={`flex-1 h-0.5 mx-4 ${index < currentStep ? 'bg-orange-600' : 'bg-gray-200'}`} />
              )}
            </div>
          ))}
        </div>
        <Progress value={(currentStep / (STEP_TITLES.length - 1)) * 100} className="h-2" />
      </div>

      {/* Step Content */}
      {renderStepContent()}

      {/* Navigation Buttons */}
      <div className="flex justify-between mt-8">
        <Button
          onClick={prevStep}
          disabled={currentStep === 0}
          variant="outline"
        >
          Previous
        </Button>
        
        <Button
          onClick={nextStep}
          disabled={!canProceed() || currentStep === STEP_TITLES.length - 1}
        >
          {currentStep === STEP_TITLES.length - 1 ? 'Submit' : 'Next'}
        </Button>
      </div>

      {/* Orders Preview Dialog */}
      <Dialog open={showPreview} onOpenChange={setShowPreview}>
        <DialogContent className="max-w-4xl max-h-[80vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Orders Preview ({orders.length} orders)</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>#</TableHead>
                  <TableHead>Recipient</TableHead>
                  <TableHead>Amount</TableHead>
                  <TableHead>Priority</TableHead>
                  <TableHead>Memo</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {orders.map((order, index) => (
                  <TableRow key={index}>
                    <TableCell>{index + 1}</TableCell>
                    <TableCell className="font-mono text-sm">
                      {order.recipientAddress.slice(0, 12)}...
                    </TableCell>
                    <TableCell>{formatCurrency(parseFloat(order.amount))}</TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          order.priority === 'URGENT' ? 'destructive' :
                          order.priority === 'HIGH' ? 'default' :
                          'secondary'
                        }
                      >
                        {order.priority}
                      </Badge>
                    </TableCell>
                    <TableCell className="max-w-[200px] truncate">
                      {order.memo || '-'}
                    </TableCell>
                    <TableCell>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => removeOrder(index)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
          <DialogFooter>
            <Button onClick={() => setShowPreview(false)}>Close</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};