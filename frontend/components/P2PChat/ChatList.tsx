import React, { useState, useEffect } from 'react';
import {
  Box,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  Avatar,
  Typography,
  Badge,
  Chip,
  TextField,
  InputAdornment,
  Tabs,
  Tab,
  Paper,
  IconButton,
  Menu,
  MenuItem,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Divider,
  Tooltip,
} from '@mui/material';
import {
  Search,
  MoreVert,
  Circle,
  Security,
  Warning,
  CheckCircle,
  Timer,
  Block,
  Delete,
  Archive,
  Star,
  StarBorder,
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';

// Types
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

interface ChatListProps {
  conversations: ChatConversation[];
  selectedConversationId?: string;
  onSelectConversation: (conversation: ChatConversation) => void;
  onArchiveConversation: (conversationId: string) => void;
  onBlockUser: (conversationId: string) => void;
  onDeleteConversation: (conversationId: string) => void;
  onToggleFavorite: (conversationId: string) => void;
  currentUserAddress: string;
}

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

const StyledPaper = styled(Paper)(({ theme }) => ({
  height: '600px',
  borderRadius: theme.spacing(2),
  overflow: 'hidden',
  display: 'flex',
  flexDirection: 'column',
}));

const SearchContainer = styled(Box)(({ theme }) => ({
  padding: theme.spacing(2),
  borderBottom: `1px solid ${theme.palette.divider}`,
}));

const ConversationItem = styled(ListItem)<{ selected?: boolean }>(({ theme, selected }) => ({
  cursor: 'pointer',
  borderRadius: theme.spacing(1),
  margin: theme.spacing(0.5, 1),
  backgroundColor: selected ? theme.palette.action.selected : 'transparent',
  '&:hover': {
    backgroundColor: theme.palette.action.hover,
  },
  position: 'relative',
}));

const StatusIndicator = styled(Box)<{ status: string }>(({ theme, status }) => ({
  position: 'absolute',
  top: 8,
  right: 8,
  width: 12,
  height: 12,
  borderRadius: '50%',
  backgroundColor:
    status === 'completed' ? theme.palette.success.main :
    status === 'disputed' ? theme.palette.error.main :
    status === 'cancelled' ? theme.palette.error.main :
    status === 'payment_confirmed' ? theme.palette.success.main :
    status === 'payment_pending' ? theme.palette.warning.main :
    theme.palette.info.main,
}));

const TrustScoreChip = styled(Chip)<{ trustscore: number }>(({ theme, trustscore }) => ({
  height: 20,
  fontSize: '0.7rem',
  backgroundColor:
    trustscore >= 90 ? '#9c27b0' :
    trustscore >= 80 ? '#607d8b' :
    trustscore >= 70 ? '#ff9800' :
    trustscore >= 60 ? '#9e9e9e' :
    trustscore >= 50 ? '#795548' :
    '#f44336',
  color: 'white',
}));

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`chat-tabpanel-${index}`}
      aria-labelledby={`chat-tab-${index}`}
      {...other}
      style={{ height: '100%', display: value === index ? 'flex' : 'none', flexDirection: 'column' }}
    >
      {children}
    </div>
  );
}

export const P2PChatList: React.FC<ChatListProps> = ({
  conversations,
  selectedConversationId,
  onSelectConversation,
  onArchiveConversation,
  onBlockUser,
  onDeleteConversation,
  onToggleFavorite,
  currentUserAddress,
}) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [activeTab, setActiveTab] = useState(0);
  const [anchorEl, setAnchorEl] = useState<{ element: HTMLElement; conversationId: string } | null>(null);
  const [confirmDialog, setConfirmDialog] = useState<{ open: boolean; action: string; conversationId: string }>({
    open: false,
    action: '',
    conversationId: '',
  });

  const handleMenuClick = (event: React.MouseEvent<HTMLElement>, conversationId: string) => {
    event.stopPropagation();
    setAnchorEl({ element: event.currentTarget, conversationId });
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleMenuAction = (action: string, conversationId: string) => {
    setAnchorEl(null);
    
    if (action === 'delete' || action === 'block') {
      setConfirmDialog({ open: true, action, conversationId });
    } else if (action === 'archive') {
      onArchiveConversation(conversationId);
    } else if (action === 'favorite') {
      onToggleFavorite(conversationId);
    }
  };

  const handleConfirmAction = () => {
    const { action, conversationId } = confirmDialog;
    
    if (action === 'delete') {
      onDeleteConversation(conversationId);
    } else if (action === 'block') {
      onBlockUser(conversationId);
    }
    
    setConfirmDialog({ open: false, action: '', conversationId: '' });
  };

  const filterConversations = (conversations: ChatConversation[], filter: string) => {
    let filtered = conversations;

    // Apply search filter
    if (searchTerm) {
      filtered = filtered.filter(
        conv =>
          conv.counterparty.displayName.toLowerCase().includes(searchTerm.toLowerCase()) ||
          conv.tradeId.toLowerCase().includes(searchTerm.toLowerCase()) ||
          conv.lastMessage.content.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    // Apply tab filter
    switch (filter) {
      case 'active':
        return filtered.filter(conv => 
          ['matched', 'payment_pending', 'payment_confirmed'].includes(conv.tradeStatus) &&
          !conv.isArchived && !conv.isBlocked
        );
      case 'completed':
        return filtered.filter(conv => 
          conv.tradeStatus === 'completed' && !conv.isArchived
        );
      case 'disputed':
        return filtered.filter(conv => 
          conv.tradeStatus === 'disputed' && !conv.isArchived
        );
      case 'archived':
        return filtered.filter(conv => conv.isArchived);
      default:
        return filtered.filter(conv => !conv.isArchived && !conv.isBlocked);
    }
  };

  const getFilteredConversations = () => {
    const filters = ['all', 'active', 'completed', 'disputed', 'archived'];
    return filterConversations(conversations, filters[activeTab]);
  };

  const getTabCounts = () => {
    return {
      all: conversations.filter(c => !c.isArchived && !c.isBlocked).length,
      active: conversations.filter(c => 
        ['matched', 'payment_pending', 'payment_confirmed'].includes(c.tradeStatus) &&
        !c.isArchived && !c.isBlocked
      ).length,
      completed: conversations.filter(c => c.tradeStatus === 'completed' && !c.isArchived).length,
      disputed: conversations.filter(c => c.tradeStatus === 'disputed' && !c.isArchived).length,
      archived: conversations.filter(c => c.isArchived).length,
    };
  };

  const formatTimestamp = (date: Date) => {
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const hours = diff / (1000 * 60 * 60);
    
    if (hours < 1) {
      const minutes = Math.floor(diff / (1000 * 60));
      return `${minutes}m`;
    } else if (hours < 24) {
      return `${Math.floor(hours)}h`;
    } else {
      const days = Math.floor(hours / 24);
      return days === 1 ? '1d' : `${days}d`;
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed': return <CheckCircle sx={{ fontSize: 16, color: 'success.main' }} />;
      case 'disputed': return <Warning sx={{ fontSize: 16, color: 'error.main' }} />;
      case 'cancelled': return <Block sx={{ fontSize: 16, color: 'error.main' }} />;
      case 'payment_confirmed': return <CheckCircle sx={{ fontSize: 16, color: 'success.main' }} />;
      case 'payment_pending': return <Timer sx={{ fontSize: 16, color: 'warning.main' }} />;
      default: return <Circle sx={{ fontSize: 16, color: 'info.main' }} />;
    }
  };

  const getTrustBadge = (trustScore: number) => {
    if (trustScore >= 90) return 'D'; // Diamond
    if (trustScore >= 80) return 'P'; // Platinum
    if (trustScore >= 70) return 'G'; // Gold
    if (trustScore >= 60) return 'S'; // Silver
    if (trustScore >= 50) return 'B'; // Bronze
    return 'N'; // New
  };

  const tabCounts = getTabCounts();
  const filteredConversations = getFilteredConversations();

  // Sort conversations: pinned first, then by last message timestamp
  const sortedConversations = [...filteredConversations].sort((a, b) => {
    if (a.isPinned && !b.isPinned) return -1;
    if (!a.isPinned && b.isPinned) return 1;
    return new Date(b.lastMessage.timestamp).getTime() - new Date(a.lastMessage.timestamp).getTime();
  });

  return (
    <StyledPaper>
      {/* Search */}
      <SearchContainer>
        <TextField
          fullWidth
          size="small"
          placeholder="Search conversations..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <Search />
              </InputAdornment>
            ),
          }}
        />
      </SearchContainer>

      {/* Tabs */}
      <Tabs
        value={activeTab}
        onChange={(_, newValue) => setActiveTab(newValue)}
        variant="scrollable"
        scrollButtons="auto"
        sx={{ borderBottom: 1, borderColor: 'divider', minHeight: 48 }}
      >
        <Tab label={`All (${tabCounts.all})`} />
        <Tab label={`Active (${tabCounts.active})`} />
        <Tab label={`Completed (${tabCounts.completed})`} />
        <Tab label={`Disputed (${tabCounts.disputed})`} />
        <Tab label={`Archived (${tabCounts.archived})`} />
      </Tabs>

      {/* Conversation List */}
      <TabPanel value={activeTab} index={activeTab}>
        <Box sx={{ flex: 1, overflow: 'auto' }}>
          {sortedConversations.length === 0 ? (
            <Box sx={{ p: 3, textAlign: 'center' }}>
              <Typography variant="body2" color="text.secondary">
                {searchTerm ? 'No conversations found' : 'No conversations yet'}
              </Typography>
            </Box>
          ) : (
            <List sx={{ p: 0 }}>
              {sortedConversations.map((conversation) => (
                <ConversationItem
                  key={conversation.id}
                  selected={selectedConversationId === conversation.id}
                  onClick={() => onSelectConversation(conversation)}
                >
                  <ListItemAvatar>
                    <Badge
                      overlap="circular"
                      anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
                      variant="dot"
                      sx={{
                        '& .MuiBadge-badge': {
                          backgroundColor: conversation.counterparty.isOnline ? '#44b700' : 'transparent',
                        },
                      }}
                    >
                      <Avatar>
                        {conversation.counterparty.displayName.charAt(0)}
                      </Avatar>
                    </Badge>
                  </ListItemAvatar>

                  <ListItemText
                    primary={
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5 }}>
                        <Typography variant="subtitle2" sx={{ fontWeight: conversation.unreadCount > 0 ? 'bold' : 'normal' }}>
                          {conversation.counterparty.displayName}
                        </Typography>
                        <TrustScoreChip
                          label={getTrustBadge(conversation.counterparty.trustScore)}
                          size="small"
                          trustscore={conversation.counterparty.trustScore}
                        />
                        {conversation.counterparty.isKYCVerified && (
                          <Security sx={{ fontSize: 14, color: 'success.main' }} />
                        )}
                        {conversation.isFavorite && (
                          <Star sx={{ fontSize: 14, color: 'warning.main' }} />
                        )}
                      </Box>
                    }
                    secondary={
                      <Box>
                        <Typography
                          variant="body2"
                          color="text.secondary"
                          sx={{
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            whiteSpace: 'nowrap',
                            fontWeight: conversation.unreadCount > 0 ? 'bold' : 'normal',
                          }}
                        >
                          {conversation.lastMessage.isSystem ? '🤖 ' : ''}
                          {conversation.lastMessage.content}
                        </Typography>
                        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mt: 0.5 }}>
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            {getStatusIcon(conversation.tradeStatus)}
                            <Typography variant="caption" color="text.secondary">
                              {conversation.tradeAmount} NAMO • {conversation.tradeId.slice(-6)}
                            </Typography>
                          </Box>
                          <Typography variant="caption" color="text.secondary">
                            {formatTimestamp(conversation.lastMessage.timestamp)}
                          </Typography>
                        </Box>
                      </Box>
                    }
                  />

                  {/* Unread count badge */}
                  {conversation.unreadCount > 0 && (
                    <Badge
                      badgeContent={conversation.unreadCount}
                      color="primary"
                      sx={{ position: 'absolute', top: 8, right: 32 }}
                    />
                  )}

                  {/* Status indicator */}
                  <StatusIndicator status={conversation.tradeStatus} />

                  {/* Menu button */}
                  <IconButton
                    size="small"
                    onClick={(e) => handleMenuClick(e, conversation.id)}
                    sx={{ position: 'absolute', top: 8, right: 8, opacity: 0.7 }}
                  >
                    <MoreVert sx={{ fontSize: 16 }} />
                  </IconButton>
                </ConversationItem>
              ))}
            </List>
          )}
        </Box>
      </TabPanel>

      {/* Context Menu */}
      <Menu
        anchorEl={anchorEl?.element}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
      >
        <MenuItem onClick={() => handleMenuAction('favorite', anchorEl?.conversationId || '')}>
          {conversations.find(c => c.id === anchorEl?.conversationId)?.isFavorite ? (
            <><Star sx={{ mr: 1 }} /> Remove from Favorites</>
          ) : (
            <><StarBorder sx={{ mr: 1 }} /> Add to Favorites</>
          )}
        </MenuItem>
        <MenuItem onClick={() => handleMenuAction('archive', anchorEl?.conversationId || '')}>
          <Archive sx={{ mr: 1 }} />
          Archive
        </MenuItem>
        <Divider />
        <MenuItem onClick={() => handleMenuAction('block', anchorEl?.conversationId || '')}>
          <Block sx={{ mr: 1 }} />
          Block User
        </MenuItem>
        <MenuItem onClick={() => handleMenuAction('delete', anchorEl?.conversationId || '')}>
          <Delete sx={{ mr: 1 }} />
          Delete Conversation
        </MenuItem>
      </Menu>

      {/* Confirmation Dialog */}
      <Dialog open={confirmDialog.open} onClose={() => setConfirmDialog({ ...confirmDialog, open: false })}>
        <DialogTitle>
          {confirmDialog.action === 'delete' ? 'Delete Conversation' : 'Block User'}
        </DialogTitle>
        <DialogContent>
          <Typography>
            {confirmDialog.action === 'delete'
              ? 'Are you sure you want to delete this conversation? This action cannot be undone.'
              : 'Are you sure you want to block this user? You will no longer receive messages from them.'}
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setConfirmDialog({ ...confirmDialog, open: false })}>
            Cancel
          </Button>
          <Button onClick={handleConfirmAction} color="error" variant="contained">
            {confirmDialog.action === 'delete' ? 'Delete' : 'Block'}
          </Button>
        </DialogActions>
      </Dialog>
    </StyledPaper>
  );
};

export default P2PChatList;