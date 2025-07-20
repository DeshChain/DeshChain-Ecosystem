import React, { useState, useEffect, useCallback } from 'react';
import {
  Box,
  Container,
  Grid,
  Typography,
  Paper,
  Alert,
  Snackbar,
  useTheme,
  useMediaQuery,
  Drawer,
  IconButton,
  AppBar,
  Toolbar,
} from '@mui/material';
import { Menu, Close } from '@mui/icons-material';
import { styled } from '@mui/material/styles';

import P2PChatList from './ChatList';
import P2PChatWindow from './ChatWindow';

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

interface ChatConversation {
  id: string;
  tradeId: string;
  counterparty: {
    address: string;
    displayName: string;
    avatar?: string;
    trustScore: number;
    isKYCVerified: boolean;
    isOnline: boolean;
    lastSeen: Date;
  };
  lastMessage: {
    content: string;
    timestamp: Date;
    isSystem: boolean;
    sender: string;
  };
  unreadCount: number;
  tradeStatus: 'matched' | 'payment_pending' | 'payment_confirmed' | 'completed' | 'disputed' | 'cancelled';
  tradeAmount: string;
  tradeCurrency: string;
  isArchived: boolean;
  isBlocked: boolean;
  isFavorite: boolean;
  isPinned: boolean;
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

const StyledContainer = styled(Container)(({ theme }) => ({
  paddingTop: theme.spacing(3),
  paddingBottom: theme.spacing(3),
  height: '100vh',
  display: 'flex',
  flexDirection: 'column',
}));

const ChatGrid = styled(Grid)(({ theme }) => ({
  height: '100%',
  '& .MuiGrid-item': {
    display: 'flex',
    flexDirection: 'column',
  },
}));

const MobileChatHeader = styled(AppBar)(({ theme }) => ({
  [theme.breakpoints.up('md')]: {
    display: 'none',
  },
}));

const EmptyState = styled(Paper)(({ theme }) => ({
  height: '600px',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  borderRadius: theme.spacing(2),
  backgroundColor: theme.palette.background.paper,
}));

export const P2PChatSystem: React.FC = () => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  
  const [conversations, setConversations] = useState<ChatConversation[]>([]);
  const [selectedConversation, setSelectedConversation] = useState<ChatConversation | null>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>('');
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [notification, setNotification] = useState<{ message: string; severity: 'success' | 'error' | 'info' }>({ message: '', severity: 'info' });

  // Initialize with mock data
  useEffect(() => {
    initializeMockData();
  }, []);

  const initializeMockData = () => {
    // Mock current user
    const mockCurrentUser: User = {
      address: 'desh1current123456',
      displayName: 'You',
      trustScore: 85,
      isKYCVerified: true,
      isOnline: true,
    };
    setCurrentUser(mockCurrentUser);

    // Mock conversations
    const mockConversations: ChatConversation[] = [
      {
        id: 'conv-1',
        tradeId: 'TRADE-2025011901-001',
        counterparty: {
          address: 'desh1user789',
          displayName: 'राम शर्मा',
          trustScore: 92,
          isKYCVerified: true,
          isOnline: true,
          lastSeen: new Date(),
        },
        lastMessage: {
          content: 'Payment received. Confirming transaction now.',
          timestamp: new Date(Date.now() - 1000 * 60 * 5), // 5 minutes ago
          isSystem: false,
          sender: 'desh1user789',
        },
        unreadCount: 2,
        tradeStatus: 'payment_confirmed',
        tradeAmount: '10000',
        tradeCurrency: 'NAMO',
        isArchived: false,
        isBlocked: false,
        isFavorite: true,
        isPinned: true,
      },
      {
        id: 'conv-2',
        tradeId: 'TRADE-2025011901-002',
        counterparty: {
          address: 'desh1trader456',
          displayName: 'Priya Patel',
          trustScore: 78,
          isKYCVerified: true,
          isOnline: false,
          lastSeen: new Date(Date.now() - 1000 * 60 * 30), // 30 minutes ago
        },
        lastMessage: {
          content: 'Trade matched successfully. Please proceed with payment.',
          timestamp: new Date(Date.now() - 1000 * 60 * 15), // 15 minutes ago
          isSystem: true,
          sender: 'system',
        },
        unreadCount: 0,
        tradeStatus: 'payment_pending',
        tradeAmount: '5000',
        tradeCurrency: 'NAMO',
        isArchived: false,
        isBlocked: false,
        isFavorite: false,
        isPinned: false,
      },
      {
        id: 'conv-3',
        tradeId: 'TRADE-2025011801-003',
        counterparty: {
          address: 'desh1merchant123',
          displayName: 'Ahmed Khan',
          trustScore: 88,
          isKYCVerified: true,
          isOnline: true,
          lastSeen: new Date(),
        },
        lastMessage: {
          content: 'Thanks for the smooth transaction! ⭐⭐⭐⭐⭐',
          timestamp: new Date(Date.now() - 1000 * 60 * 60 * 2), // 2 hours ago
          isSystem: false,
          sender: 'desh1merchant123',
        },
        unreadCount: 0,
        tradeStatus: 'completed',
        tradeAmount: '25000',
        tradeCurrency: 'NAMO',
        isArchived: false,
        isBlocked: false,
        isFavorite: false,
        isPinned: false,
      },
    ];
    setConversations(mockConversations);
  };

  const loadMessages = useCallback(async (conversationId: string) => {
    setLoading(true);
    try {
      // Mock messages for the conversation
      const mockMessages: ChatMessage[] = [
        {
          id: 'msg-1',
          sender: 'system',
          receiver: '',
          message: 'Trade matched successfully. Escrow has been created.',
          timestamp: new Date(Date.now() - 1000 * 60 * 60), // 1 hour ago
          isSystem: true,
          messageType: 'trade_update',
          metadata: {
            tradeId: selectedConversation?.tradeId,
            status: 'matched',
          },
          isEncrypted: false,
        },
        {
          id: 'msg-2',
          sender: selectedConversation?.counterparty.address || '',
          receiver: currentUser?.address || '',
          message: 'Hi! I see we have been matched for the NAMO trade. I am ready to proceed.',
          timestamp: new Date(Date.now() - 1000 * 60 * 55), // 55 minutes ago
          isSystem: false,
          messageType: 'text',
          isEncrypted: true,
        },
        {
          id: 'msg-3',
          sender: currentUser?.address || '',
          receiver: selectedConversation?.counterparty.address || '',
          message: 'Great! I will send the fiat payment now via UPI. Please confirm once received.',
          timestamp: new Date(Date.now() - 1000 * 60 * 50), // 50 minutes ago
          isSystem: false,
          messageType: 'text',
          isEncrypted: true,
        },
        {
          id: 'msg-4',
          sender: currentUser?.address || '',
          receiver: selectedConversation?.counterparty.address || '',
          message: 'Payment details request',
          timestamp: new Date(Date.now() - 1000 * 60 * 45), // 45 minutes ago
          isSystem: false,
          messageType: 'payment_request',
          metadata: {
            amount: selectedConversation?.tradeAmount,
            currency: 'INR',
            paymentMethod: 'UPI',
          },
          isEncrypted: true,
        },
        {
          id: 'msg-5',
          sender: selectedConversation?.counterparty.address || '',
          receiver: currentUser?.address || '',
          message: 'Payment received successfully! Releasing NAMO from escrow now.',
          timestamp: new Date(Date.now() - 1000 * 60 * 10), // 10 minutes ago
          isSystem: false,
          messageType: 'text',
          isEncrypted: true,
        },
        {
          id: 'msg-6',
          sender: 'system',
          receiver: '',
          message: 'Payment confirmed by seller. NAMO released to buyer.',
          timestamp: new Date(Date.now() - 1000 * 60 * 5), // 5 minutes ago
          isSystem: true,
          messageType: 'trade_update',
          metadata: {
            tradeId: selectedConversation?.tradeId,
            status: 'payment_confirmed',
          },
          isEncrypted: false,
        },
      ];

      setMessages(mockMessages);
    } catch (error) {
      setError('Failed to load messages');
    } finally {
      setLoading(false);
    }
  }, [selectedConversation?.tradeId, selectedConversation?.counterparty.address, currentUser?.address]);

  useEffect(() => {
    if (selectedConversation) {
      loadMessages(selectedConversation.id);
    }
  }, [selectedConversation, loadMessages]);

  const handleSelectConversation = (conversation: ChatConversation) => {
    setSelectedConversation(conversation);
    
    // Mark messages as read
    setConversations(prev => 
      prev.map(conv => 
        conv.id === conversation.id ? { ...conv, unreadCount: 0 } : conv
      )
    );

    // Close mobile drawer
    if (isMobile) {
      setDrawerOpen(false);
    }
  };

  const handleSendMessage = async (message: string, type = 'text', metadata?: any) => {
    if (!selectedConversation || !currentUser) return;

    const newMessage: ChatMessage = {
      id: `msg-${Date.now()}`,
      sender: currentUser.address,
      receiver: selectedConversation.counterparty.address,
      message,
      timestamp: new Date(),
      isSystem: false,
      messageType: type as any,
      metadata,
      isEncrypted: true,
    };

    setMessages(prev => [...prev, newMessage]);

    // Update conversation's last message
    setConversations(prev =>
      prev.map(conv =>
        conv.id === selectedConversation.id
          ? {
              ...conv,
              lastMessage: {
                content: message,
                timestamp: new Date(),
                isSystem: false,
                sender: currentUser.address,
              },
            }
          : conv
      )
    );

    try {
      // Here you would send the message to the backend
      // await chatAPI.sendMessage(newMessage);
    } catch (error) {
      setError('Failed to send message');
    }
  };

  const handleArchiveConversation = (conversationId: string) => {
    setConversations(prev =>
      prev.map(conv =>
        conv.id === conversationId ? { ...conv, isArchived: true } : conv
      )
    );
    setNotification({ message: 'Conversation archived', severity: 'success' });
  };

  const handleBlockUser = (conversationId: string) => {
    setConversations(prev =>
      prev.map(conv =>
        conv.id === conversationId ? { ...conv, isBlocked: true } : conv
      )
    );
    setNotification({ message: 'User blocked', severity: 'success' });
  };

  const handleDeleteConversation = (conversationId: string) => {
    setConversations(prev => prev.filter(conv => conv.id !== conversationId));
    if (selectedConversation?.id === conversationId) {
      setSelectedConversation(null);
    }
    setNotification({ message: 'Conversation deleted', severity: 'success' });
  };

  const handleToggleFavorite = (conversationId: string) => {
    setConversations(prev =>
      prev.map(conv =>
        conv.id === conversationId ? { ...conv, isFavorite: !conv.isFavorite } : conv
      )
    );
  };

  const handleReportUser = (reason: string) => {
    // Handle user reporting
    setNotification({ message: 'User reported', severity: 'success' });
  };

  const handleBlockUserFromChat = () => {
    if (selectedConversation) {
      handleBlockUser(selectedConversation.id);
    }
  };

  const handleConfirmPayment = () => {
    // Handle payment confirmation
    if (selectedConversation) {
      const systemMessage: ChatMessage = {
        id: `msg-${Date.now()}`,
        sender: 'system',
        receiver: '',
        message: 'Payment confirmed. Trade completed successfully.',
        timestamp: new Date(),
        isSystem: true,
        messageType: 'trade_update',
        metadata: {
          tradeId: selectedConversation.tradeId,
          status: 'completed',
        },
        isEncrypted: false,
      };

      setMessages(prev => [...prev, systemMessage]);
      setConversations(prev =>
        prev.map(conv =>
          conv.id === selectedConversation.id
            ? { ...conv, tradeStatus: 'completed' }
            : conv
        )
      );
      setNotification({ message: 'Payment confirmed successfully', severity: 'success' });
    }
  };

  const handleDisputeTrade = (reason: string) => {
    // Handle trade dispute
    if (selectedConversation) {
      setConversations(prev =>
        prev.map(conv =>
          conv.id === selectedConversation.id
            ? { ...conv, tradeStatus: 'disputed' }
            : conv
        )
      );
      setNotification({ message: 'Dispute filed successfully', severity: 'info' });
    }
  };

  const handleCancelTrade = () => {
    // Handle trade cancellation
    if (selectedConversation) {
      setConversations(prev =>
        prev.map(conv =>
          conv.id === selectedConversation.id
            ? { ...conv, tradeStatus: 'cancelled' }
            : conv
        )
      );
      setNotification({ message: 'Trade cancelled', severity: 'info' });
    }
  };

  const getCurrentTradeInfo = (): TradeInfo | null => {
    if (!selectedConversation) return null;

    return {
      tradeId: selectedConversation.tradeId,
      status: selectedConversation.tradeStatus,
      buyer: currentUser?.address === selectedConversation.counterparty.address 
        ? selectedConversation.counterparty.address 
        : currentUser?.address || '',
      seller: currentUser?.address === selectedConversation.counterparty.address 
        ? currentUser?.address || ''
        : selectedConversation.counterparty.address,
      amount: selectedConversation.tradeAmount,
      fiatAmount: '50000', // Mock fiat amount
      currency: 'INR',
      paymentMethod: 'UPI',
      expiresAt: new Date(Date.now() + 1000 * 60 * 60 * 2), // 2 hours from now
      createdAt: new Date(Date.now() - 1000 * 60 * 60), // 1 hour ago
    };
  };

  const chatListComponent = (
    <P2PChatList
      conversations={conversations}
      selectedConversationId={selectedConversation?.id}
      onSelectConversation={handleSelectConversation}
      onArchiveConversation={handleArchiveConversation}
      onBlockUser={handleBlockUser}
      onDeleteConversation={handleDeleteConversation}
      onToggleFavorite={handleToggleFavorite}
      currentUserAddress={currentUser?.address || ''}
    />
  );

  return (
    <StyledContainer maxWidth="xl">
      {/* Mobile Header */}
      {isMobile && selectedConversation && (
        <MobileChatHeader position="static" color="default" elevation={0}>
          <Toolbar>
            <IconButton edge="start" onClick={() => setDrawerOpen(true)}>
              <Menu />
            </IconButton>
            <Typography variant="h6" sx={{ flexGrow: 1 }}>
              {selectedConversation.counterparty.displayName}
            </Typography>
          </Toolbar>
        </MobileChatHeader>
      )}

      {/* Header */}
      <Box sx={{ mb: 3 }}>
        <Typography variant="h4" component="h1" gutterBottom sx={{ fontWeight: 'bold' }}>
          P2P Trade Chat
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Secure, encrypted communication for your peer-to-peer trades
        </Typography>
      </Box>

      {/* Main Content */}
      <ChatGrid container spacing={3} sx={{ flex: 1 }}>
        {/* Desktop: Chat List */}
        {!isMobile && (
          <Grid item md={4}>
            {chatListComponent}
          </Grid>
        )}

        {/* Chat Window */}
        <Grid item xs={12} md={8}>
          {selectedConversation && currentUser ? (
            <P2PChatWindow
              tradeInfo={getCurrentTradeInfo()!}
              currentUser={currentUser}
              counterparty={selectedConversation.counterparty}
              messages={messages}
              onSendMessage={handleSendMessage}
              onReportUser={handleReportUser}
              onBlockUser={handleBlockUserFromChat}
              onConfirmPayment={handleConfirmPayment}
              onDisputeTrade={handleDisputeTrade}
              onCancelTrade={handleCancelTrade}
              isLoading={loading}
              error={error}
            />
          ) : (
            <EmptyState>
              <Box sx={{ textAlign: 'center' }}>
                <Typography variant="h6" color="text.secondary" gutterBottom>
                  Select a conversation to start chatting
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Your P2P trade conversations will appear here
                </Typography>
                {isMobile && (
                  <IconButton
                    onClick={() => setDrawerOpen(true)}
                    sx={{ mt: 2 }}
                    color="primary"
                  >
                    <Menu />
                  </IconButton>
                )}
              </Box>
            </EmptyState>
          )}
        </Grid>
      </ChatGrid>

      {/* Mobile Drawer */}
      {isMobile && (
        <Drawer
          anchor="left"
          open={drawerOpen}
          onClose={() => setDrawerOpen(false)}
          sx={{
            '& .MuiDrawer-paper': {
              width: '80%',
              maxWidth: 400,
            },
          }}
        >
          <Box sx={{ p: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Typography variant="h6">Conversations</Typography>
            <IconButton onClick={() => setDrawerOpen(false)}>
              <Close />
            </IconButton>
          </Box>
          {chatListComponent}
        </Drawer>
      )}

      {/* Notifications */}
      <Snackbar
        open={!!notification.message}
        autoHideDuration={4000}
        onClose={() => setNotification({ ...notification, message: '' })}
      >
        <Alert
          onClose={() => setNotification({ ...notification, message: '' })}
          severity={notification.severity}
        >
          {notification.message}
        </Alert>
      </Snackbar>
    </StyledContainer>
  );
};

export default P2PChatSystem;