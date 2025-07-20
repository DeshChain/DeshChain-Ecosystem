/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import React, { useState, useEffect } from 'react';
import {
  Box,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  ListItemAvatar,
  Avatar,
  Typography,
  TextField,
  InputAdornment,
  IconButton,
  Chip,
  Paper,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions
} from '@mui/material';
import {
  Search as SearchIcon,
  Person as PersonIcon,
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Star as StarIcon,
  StarBorder as StarBorderIcon
} from '@mui/icons-material';
import { motion, AnimatePresence } from 'framer-motion';

interface Contact {
  id: string;
  name: string;
  address: string;
  tags?: string[];
  isFavorite?: boolean;
  lastUsed?: Date;
}

interface AddressBookProps {
  onSelect: (address: string) => void;
  showAddButton?: boolean;
}

export const AddressBook: React.FC<AddressBookProps> = ({ 
  onSelect, 
  showAddButton = true 
}) => {
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [showAddDialog, setShowAddDialog] = useState(false);
  const [editingContact, setEditingContact] = useState<Contact | null>(null);
  const [newContact, setNewContact] = useState<Partial<Contact>>({
    name: '',
    address: '',
    tags: []
  });

  // Load contacts from localStorage
  useEffect(() => {
    const savedContacts = localStorage.getItem('deshchain-address-book');
    if (savedContacts) {
      setContacts(JSON.parse(savedContacts));
    } else {
      // Default contacts for demo
      setContacts([
        {
          id: '1',
          name: 'Family Fund',
          address: 'desh1abcdef1234567890abcdef1234567890abcdef',
          tags: ['family', 'monthly'],
          isFavorite: true,
          lastUsed: new Date('2024-01-10')
        },
        {
          id: '2',
          name: 'Business Partner',
          address: 'desh1fedcba0987654321fedcba0987654321fedcba',
          tags: ['business'],
          isFavorite: false,
          lastUsed: new Date('2024-01-15')
        },
        {
          id: '3',
          name: 'Village Pool',
          address: 'desh1village123456789012345678901234567890',
          tags: ['community', 'village'],
          isFavorite: true,
          lastUsed: new Date('2024-01-18')
        }
      ]);
    }
  }, []);

  // Save contacts to localStorage whenever they change
  useEffect(() => {
    if (contacts.length > 0) {
      localStorage.setItem('deshchain-address-book', JSON.stringify(contacts));
    }
  }, [contacts]);

  const filteredContacts = contacts.filter(contact =>
    contact.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    contact.address.toLowerCase().includes(searchTerm.toLowerCase()) ||
    contact.tags?.some(tag => tag.toLowerCase().includes(searchTerm.toLowerCase()))
  );

  const sortedContacts = [...filteredContacts].sort((a, b) => {
    // Favorites first
    if (a.isFavorite && !b.isFavorite) return -1;
    if (!a.isFavorite && b.isFavorite) return 1;
    
    // Then by last used
    const aTime = a.lastUsed?.getTime() || 0;
    const bTime = b.lastUsed?.getTime() || 0;
    return bTime - aTime;
  });

  const handleAddContact = () => {
    if (newContact.name && newContact.address) {
      const contact: Contact = {
        id: Date.now().toString(),
        name: newContact.name,
        address: newContact.address,
        tags: newContact.tags || [],
        isFavorite: false,
        lastUsed: new Date()
      };

      if (editingContact) {
        setContacts(contacts.map(c => c.id === editingContact.id ? { ...contact, id: editingContact.id } : c));
        setEditingContact(null);
      } else {
        setContacts([...contacts, contact]);
      }

      setNewContact({ name: '', address: '', tags: [] });
      setShowAddDialog(false);
    }
  };

  const handleDeleteContact = (id: string) => {
    setContacts(contacts.filter(c => c.id !== id));
  };

  const toggleFavorite = (id: string) => {
    setContacts(contacts.map(c => 
      c.id === id ? { ...c, isFavorite: !c.isFavorite } : c
    ));
  };

  const handleSelectContact = (contact: Contact) => {
    // Update last used
    setContacts(contacts.map(c => 
      c.id === contact.id ? { ...c, lastUsed: new Date() } : c
    ));
    onSelect(contact.address);
  };

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map(word => word[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  return (
    <Paper elevation={1} sx={{ p: 2, borderRadius: 2 }}>
      <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
        <Typography variant="subtitle1" fontWeight="bold">
          Address Book
        </Typography>
        {showAddButton && (
          <Button
            size="small"
            startIcon={<AddIcon />}
            onClick={() => {
              setEditingContact(null);
              setNewContact({ name: '', address: '', tags: [] });
              setShowAddDialog(true);
            }}
          >
            Add Contact
          </Button>
        )}
      </Box>

      <TextField
        fullWidth
        size="small"
        placeholder="Search contacts..."
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <SearchIcon fontSize="small" />
            </InputAdornment>
          )
        }}
        sx={{ mb: 2 }}
      />

      <List sx={{ maxHeight: 300, overflow: 'auto' }}>
        <AnimatePresence>
          {sortedContacts.map((contact, index) => (
            <motion.div
              key={contact.id}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              transition={{ duration: 0.2, delay: index * 0.05 }}
            >
              <ListItem
                disablePadding
                secondaryAction={
                  <Box>
                    <IconButton
                      size="small"
                      onClick={() => toggleFavorite(contact.id)}
                    >
                      {contact.isFavorite ? (
                        <StarIcon color="warning" fontSize="small" />
                      ) : (
                        <StarBorderIcon fontSize="small" />
                      )}
                    </IconButton>
                    <IconButton
                      size="small"
                      onClick={() => {
                        setEditingContact(contact);
                        setNewContact(contact);
                        setShowAddDialog(true);
                      }}
                    >
                      <EditIcon fontSize="small" />
                    </IconButton>
                    <IconButton
                      size="small"
                      onClick={() => handleDeleteContact(contact.id)}
                    >
                      <DeleteIcon fontSize="small" />
                    </IconButton>
                  </Box>
                }
              >
                <ListItemButton onClick={() => handleSelectContact(contact)}>
                  <ListItemAvatar>
                    <Avatar sx={{ bgcolor: 'primary.main', width: 36, height: 36 }}>
                      {getInitials(contact.name)}
                    </Avatar>
                  </ListItemAvatar>
                  <ListItemText
                    primary={
                      <Box display="flex" alignItems="center" gap={1}>
                        <Typography variant="body2" fontWeight="medium">
                          {contact.name}
                        </Typography>
                        {contact.tags?.map(tag => (
                          <Chip
                            key={tag}
                            label={tag}
                            size="small"
                            variant="outlined"
                          />
                        ))}
                      </Box>
                    }
                    secondary={
                      <Typography variant="caption" sx={{ fontFamily: 'monospace' }}>
                        {contact.address.slice(0, 20)}...{contact.address.slice(-10)}
                      </Typography>
                    }
                  />
                </ListItemButton>
              </ListItem>
            </motion.div>
          ))}
        </AnimatePresence>

        {sortedContacts.length === 0 && (
          <ListItem>
            <ListItemText
              primary={
                <Typography variant="body2" color="text.secondary" textAlign="center">
                  {searchTerm ? 'No contacts found' : 'No contacts yet'}
                </Typography>
              }
            />
          </ListItem>
        )}
      </List>

      {/* Add/Edit Contact Dialog */}
      <Dialog open={showAddDialog} onClose={() => setShowAddDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>
          {editingContact ? 'Edit Contact' : 'Add New Contact'}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              fullWidth
              label="Name"
              value={newContact.name || ''}
              onChange={(e) => setNewContact({ ...newContact, name: e.target.value })}
              autoFocus
            />
            <TextField
              fullWidth
              label="Address"
              value={newContact.address || ''}
              onChange={(e) => setNewContact({ ...newContact, address: e.target.value })}
              placeholder="desh1..."
              sx={{ fontFamily: 'monospace' }}
            />
            <TextField
              fullWidth
              label="Tags (comma separated)"
              value={newContact.tags?.join(', ') || ''}
              onChange={(e) => setNewContact({ 
                ...newContact, 
                tags: e.target.value.split(',').map(tag => tag.trim()).filter(Boolean)
              })}
              placeholder="family, business, friend"
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowAddDialog(false)}>Cancel</Button>
          <Button 
            onClick={handleAddContact} 
            variant="contained"
            disabled={!newContact.name || !newContact.address}
          >
            {editingContact ? 'Save' : 'Add'}
          </Button>
        </DialogActions>
      </Dialog>
    </Paper>
  );
};