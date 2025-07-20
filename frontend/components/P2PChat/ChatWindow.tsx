import React, { useState, useEffect, useRef, useCallback } from 'react';
import {
  Box,
  Paper,
  Typography,
  TextField,
  IconButton,
  Avatar,
  Chip,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Menu,
  MenuItem,
  Tooltip,
  Badge,
  Alert,
  LinearProgress,
  Divider,
  Card,
  CardContent,
} from '@mui/material';
import {
  Send,
  AttachFile,
  MoreVert,
  Security,
  Warning,
  CheckCircle,
  Cancel,
  Info,
  Timer,
  Verified,
  Report,
  Block,
  Phone,
  VideoCall,
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';

// Types
interface ChatMessage {
  id: string;
  sender: string;
  receiver: string;
  message: string;
  timestamp: Date;
  isSystem: boolean;
  messageType: 'text' | 'payment_request' | 'trade_update' | 'system' | 'attachment';
  metadata?: {
    tradeId?: string;
    amount?: string;
    currency?: string;
    paymentMethod?: string;
    status?: string;
    attachmentUrl?: string;
    attachmentType?: string;
  };
  isEncrypted: boolean;
  readBy?: string[];
  reactions?: { [emoji: string]: string[] };
}

interface TradeInfo {
  tradeId: string;
  status: 'matched' | 'payment_pending' | 'payment_confirmed' | 'completed' | 'disputed' | 'cancelled';
  buyer: string;
  seller: string;
  amount: string;
  fiatAmount: string;
  currency: string;
  paymentMethod: string;
  expiresAt: Date;
  createdAt: Date;
}

interface User {
  address: string;
  displayName: string;
  trustScore: number;
  isKYCVerified: boolean;
  avatar?: string;
  isOnline: boolean;
  lastSeen?: Date;
}

interface ChatWindowProps {
  tradeInfo: TradeInfo;
  currentUser: User;
  counterparty: User;
  messages: ChatMessage[];
  onSendMessage: (message: string, type?: string, metadata?: any) => void;
  onReportUser: (reason: string) => void;
  onBlockUser: () => void;
  onConfirmPayment: () => void;
  onDisputeTrade: (reason: string) => void;
  onCancelTrade: () => void;
  isLoading?: boolean;
  error?: string;
}

const ChatContainer = styled(Paper)(({ theme }) => ({
  height: '600px',
  display: 'flex',
  flexDirection: 'column',
  borderRadius: theme.spacing(2),
  overflow: 'hidden',
  boxShadow: '0 8px 32px rgba(0,0,0,0.12)',
}));

const ChatHeader = styled(Box)(({ theme }) => ({
  padding: theme.spacing(2),
  backgroundColor: theme.palette.primary.main,
  color: theme.palette.primary.contrastText,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
}));

const MessagesContainer = styled(Box)({
  flex: 1,
  overflowY: 'auto',
  padding: '16px',
  display: 'flex',
  flexDirection: 'column',
  gap: '12px',
});

const MessageBubble = styled(Box)<{ isown: boolean; isystem?: boolean }>(({ theme, isown, isystem }) => ({
  display: 'flex',
  flexDirection: isown ? 'row-reverse' : 'row',
  alignItems: 'flex-end',
  gap: theme.spacing(1),
  maxWidth: '80%',
  alignSelf: isown ? 'flex-end' : 'flex-start',
  ...(isystem && {
    alignSelf: 'center',
    maxWidth: '90%',
  }),
}));

const MessageContent = styled(Paper)<{ isown: boolean; isystem?: boolean }>(({ theme, isown, isystem }) => ({
  padding: theme.spacing(1.5),
  borderRadius: theme.spacing(2),
  maxWidth: '100%',
  wordWrap: 'break-word',
  ...(isystem ? {
    backgroundColor: theme.palette.action.hover,
    color: theme.palette.text.secondary,
    borderRadius: theme.spacing(1),
    padding: theme.spacing(1),
  } : {
    backgroundColor: isown ? theme.palette.primary.main : theme.palette.grey[100],
    color: isown ? theme.palette.primary.contrastText : theme.palette.text.primary,
  }),
}));

const InputContainer = styled(Box)(({ theme }) => ({
  padding: theme.spacing(2),
  borderTop: `1px solid ${theme.palette.divider}`,
  display: 'flex',
  alignItems: 'center',
  gap: theme.spacing(1),
}));

const TradeStatusCard = styled(Card)(({ theme }) => ({
  margin: theme.spacing(1),
  borderLeft: `4px solid ${theme.palette.primary.main}`,
}));

const TrustScoreChip = styled(Chip)<{ trustscore: number }>(({ theme, trustscore }) => ({
  fontWeight: 'bold',
  color: theme.palette.getContrastText(
    trustscore >= 90 ? '#9c27b0' : // Diamond
    trustscore >= 80 ? '#607d8b' : // Platinum
    trustscore >= 70 ? '#ff9800' : // Gold
    trustscore >= 60 ? '#9e9e9e' : // Silver
    trustscore >= 50 ? '#795548' : // Bronze
    '#f44336' // New User
  ),
  backgroundColor:
    trustscore >= 90 ? '#9c27b0' :
    trustscore >= 80 ? '#607d8b' :
    trustscore >= 70 ? '#ff9800' :
    trustscore >= 60 ? '#9e9e9e' :
    trustscore >= 50 ? '#795548' :
    '#f44336',
}));

export const P2PChatWindow: React.FC<ChatWindowProps> = ({
  tradeInfo,
  currentUser,
  counterparty,
  messages,
  onSendMessage,
  onReportUser,
  onBlockUser,
  onConfirmPayment,
  onDisputeTrade,
  onCancelTrade,
  isLoading = false,
  error,
}) => {
  const [newMessage, setNewMessage] = useState('');
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [reportDialog, setReportDialog] = useState(false);
  const [disputeDialog, setDisputeDialog] = useState(false);
  const [reportReason, setReportReason] = useState('');
  const [disputeReason, setDisputeReason] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = useCallback(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages, scrollToBottom]);

  const handleSendMessage = () => {
    if (newMessage.trim()) {
      onSendMessage(newMessage.trim());
      setNewMessage('');
    }
  };

  const handleKeyPress = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      handleSendMessage();
    }
  };

  const handleMenuClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleReport = () => {
    if (reportReason.trim()) {
      onReportUser(reportReason.trim());
      setReportDialog(false);
      setReportReason('');
    }
  };

  const handleDispute = () => {
    if (disputeReason.trim()) {
      onDisputeTrade(disputeReason.trim());
      setDisputeDialog(false);
      setDisputeReason('');
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'matched': return 'info';
      case 'payment_pending': return 'warning';
      case 'payment_confirmed': return 'success';
      case 'completed': return 'success';
      case 'disputed': return 'error';
      case 'cancelled': return 'error';
      default: return 'default';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'matched': return <Info />;
      case 'payment_pending': return <Timer />;
      case 'payment_confirmed': return <CheckCircle />;
      case 'completed': return <CheckCircle />;
      case 'disputed': return <Warning />;
      case 'cancelled': return <Cancel />;
      default: return <Info />;
    }
  };

  const getTrustBadge = (trustScore: number) => {
    if (trustScore >= 90) return 'Diamond';
    if (trustScore >= 80) return 'Platinum';
    if (trustScore >= 70) return 'Gold';
    if (trustScore >= 60) return 'Silver';
    if (trustScore >= 50) return 'Bronze';
    return 'New';
  };

  const formatTime = (date: Date) => {
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const hours = diff / (1000 * 60 * 60);
    
    if (hours < 1) {
      const minutes = Math.floor(diff / (1000 * 60));
      return `${minutes}m ago`;
    } else if (hours < 24) {
      return `${Math.floor(hours)}h ago`;
    } else {
      return date.toLocaleDateString();
    }
  };

  const getTimeRemaining = () => {
    const now = new Date();
    const timeLeft = tradeInfo.expiresAt.getTime() - now.getTime();
    
    if (timeLeft <= 0) return 'Expired';
    
    const hours = Math.floor(timeLeft / (1000 * 60 * 60));
    const minutes = Math.floor((timeLeft % (1000 * 60 * 60)) / (1000 * 60));
    
    return `${hours}h ${minutes}m left`;
  };

  const canConfirmPayment = () => {
    return tradeInfo.status === 'payment_pending' && 
           currentUser.address === tradeInfo.seller;
  };

  const canCancelTrade = () => {
    return ['matched', 'payment_pending'].includes(tradeInfo.status);
  };

  const renderMessage = (message: ChatMessage) => {
    const isOwn = message.sender === currentUser.address;
    const isSystem = message.isSystem;

    if (message.messageType === 'trade_update') {
      return (
        <Box key={message.id} sx={{ alignSelf: 'center', mb: 2, width: '100%' }}>
          <TradeStatusCard>
            <CardContent sx={{ py: 1.5 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                {getStatusIcon(message.metadata?.status || '')}
                <Typography variant="body2" fontWeight="bold">
                  Trade Status Update
                </Typography>
                <Chip
                  label={message.metadata?.status?.replace('_', ' ').toUpperCase()}
                  size="small"
                  color={getStatusColor(message.metadata?.status || '') as any}
                />
              </Box>
              <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                {message.message}
              </Typography>
            </CardContent>
          </TradeStatusCard>
        </Box>
      );
    }

    if (message.messageType === 'payment_request') {
      return (
        <Box key={message.id} sx={{ alignSelf: 'center', mb: 2, width: '100%' }}>
          <Card sx={{ mx: 2, backgroundColor: 'warning.light' }}>
            <CardContent sx={{ py: 1.5 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
                <Security color="warning" />
                <Typography variant="body2" fontWeight="bold">
                  Payment Request
                </Typography>
              </Box>
              <Typography variant="body2">
                {message.message}
              </Typography>
              <Box sx={{ mt: 1, display: 'flex', gap: 1 }}>
                <Chip label={`₹${message.metadata?.amount}`} size="small" />
                <Chip label={message.metadata?.paymentMethod} size="small" variant="outlined" />
              </Box>
            </CardContent>
          </Card>
        </Box>
      );
    }

    return (
      <MessageBubble key={message.id} isown={isOwn} isystem={isSystem}>
        {!isOwn && !isSystem && (
          <Avatar sx={{ width: 32, height: 32 }}>
            {counterparty.displayName.charAt(0)}
          </Avatar>
        )}
        
        <Box>
          <MessageContent isown={isOwn} isystem={isSystem}>
            <Typography variant="body2" sx={{ mb: 0.5 }}>
              {message.message}
            </Typography>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mt: 0.5 }}>
              <Typography variant="caption" color="text.secondary">
                {formatTime(message.timestamp)}
              </Typography>
              {message.isEncrypted && (
                <Tooltip title="End-to-end encrypted">
                  <Security sx={{ fontSize: 12 }} />
                </Tooltip>
              )}
            </Box>
          </MessageContent>
        </Box>
      </MessageBubble>
    );
  };

  return (
    <ChatContainer>
      {/* Header */}
      <ChatHeader>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <Badge
            overlap="circular"
            anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
            variant="dot"
            sx={{
              '& .MuiBadge-badge': {
                backgroundColor: counterparty.isOnline ? '#44b700' : '#grey.500',
                color: counterparty.isOnline ? '#44b700' : '#grey.500',
              },
            }}
          >
            <Avatar>{counterparty.displayName.charAt(0)}</Avatar>
          </Badge>
          
          <Box>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <Typography variant="h6">{counterparty.displayName}</Typography>
              {counterparty.isKYCVerified && (
                <Verified sx={{ fontSize: 18 }} />
              )}
            </Box>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <TrustScoreChip
                label={`${getTrustBadge(counterparty.trustScore)} ${counterparty.trustScore}`}
                size="small"
                trustscore={counterparty.trustScore}
              />
              <Typography variant="caption">
                {counterparty.isOnline ? 'Online' : `Last seen ${formatTime(counterparty.lastSeen || new Date())}`}
              </Typography>
            </Box>
          </Box>
        </Box>

        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <IconButton color="inherit" size="small">
            <Phone />
          </IconButton>
          <IconButton color="inherit" size="small">
            <VideoCall />
          </IconButton>
          <IconButton color="inherit" size="small" onClick={handleMenuClick}>
            <MoreVert />
          </IconButton>
        </Box>
      </ChatHeader>

      {/* Trade Status Bar */}
      <Box sx={{ p: 2, backgroundColor: 'background.paper', borderBottom: 1, borderColor: 'divider' }}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Typography variant="body2" fontWeight="bold">
              Trade: {tradeInfo.tradeId}
            </Typography>
            <Chip
              icon={getStatusIcon(tradeInfo.status)}
              label={tradeInfo.status.replace('_', ' ').toUpperCase()}
              size="small"
              color={getStatusColor(tradeInfo.status) as any}
            />
          </Box>
          <Typography variant="body2" color="text.secondary">
            {getTimeRemaining()}
          </Typography>
        </Box>

        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Typography variant="body2">
            {tradeInfo.amount} NAMO ↔ ₹{tradeInfo.fiatAmount} ({tradeInfo.paymentMethod})
          </Typography>
          <Box sx={{ display: 'flex', gap: 1 }}>
            {canConfirmPayment() && (
              <Button
                size="small"
                variant="contained"
                color="success"
                onClick={onConfirmPayment}
                startIcon={<CheckCircle />}
              >
                Confirm Payment
              </Button>
            )}
            {canCancelTrade() && (
              <Button
                size="small"
                variant="outlined"
                color="error"
                onClick={onCancelTrade}
                startIcon={<Cancel />}
              >
                Cancel
              </Button>
            )}
          </Box>
        </Box>

        {tradeInfo.status === 'payment_pending' && (
          <LinearProgress
            variant="determinate"
            value={75}
            sx={{ mt: 1, borderRadius: 1 }}
          />
        )}
      </Box>

      {/* Error Alert */}
      {error && (
        <Alert severity="error" sx={{ m: 1 }}>
          {error}
        </Alert>
      )}

      {/* Messages */}
      <MessagesContainer>
        {messages.map(renderMessage)}
        <div ref={messagesEndRef} />
      </MessagesContainer>

      {/* Input */}
      <InputContainer>
        <TextField
          fullWidth
          multiline
          maxRows={3}
          placeholder="Type your message..."
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          onKeyPress={handleKeyPress}
          disabled={isLoading || tradeInfo.status === 'completed' || tradeInfo.status === 'cancelled'}
          size="small"
        />
        <IconButton disabled={isLoading}>
          <AttachFile />
        </IconButton>
        <IconButton
          color="primary"
          onClick={handleSendMessage}
          disabled={!newMessage.trim() || isLoading}
        >
          <Send />
        </IconButton>
      </InputContainer>

      {/* Menu */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
      >
        <MenuItem onClick={() => { setReportDialog(true); handleMenuClose(); }}>
          <Report sx={{ mr: 1 }} />
          Report User
        </MenuItem>
        <MenuItem onClick={() => { setDisputeDialog(true); handleMenuClose(); }}>
          <Warning sx={{ mr: 1 }} />
          Dispute Trade
        </MenuItem>
        <MenuItem onClick={() => { onBlockUser(); handleMenuClose(); }}>
          <Block sx={{ mr: 1 }} />
          Block User
        </MenuItem>
      </Menu>

      {/* Report Dialog */}
      <Dialog open={reportDialog} onClose={() => setReportDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Report User</DialogTitle>
        <DialogContent>
          <TextField
            fullWidth
            multiline
            rows={4}
            placeholder="Please describe the issue..."
            value={reportReason}
            onChange={(e) => setReportReason(e.target.value)}
            sx={{ mt: 1 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setReportDialog(false)}>Cancel</Button>
          <Button onClick={handleReport} variant="contained" color="error">
            Report
          </Button>
        </DialogActions>
      </Dialog>

      {/* Dispute Dialog */}
      <Dialog open={disputeDialog} onClose={() => setDisputeDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Dispute Trade</DialogTitle>
        <DialogContent>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            Disputing a trade will freeze the escrow and involve moderators. Please provide detailed information.
          </Typography>
          <TextField
            fullWidth
            multiline
            rows={4}
            placeholder="Please describe the dispute reason..."
            value={disputeReason}
            onChange={(e) => setDisputeReason(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDisputeDialog(false)}>Cancel</Button>
          <Button onClick={handleDispute} variant="contained" color="error">
            File Dispute
          </Button>
        </DialogActions>
      </Dialog>
    </ChatContainer>
  );
};

export default P2PChatWindow;