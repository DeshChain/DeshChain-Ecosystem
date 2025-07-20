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

import { openDB, DBSchema, IDBPDatabase } from 'idb';
import { MoneyOrderFormData, ReceiptData, PoolInfo } from '../types';

// IndexedDB schema
interface MoneyOrderDB extends DBSchema {
  'pending-orders': {
    key: string;
    value: PendingOrder;
    indexes: { 'by-timestamp': number; 'by-status': string };
  };
  'receipts': {
    key: string;
    value: ReceiptData;
    indexes: { 'by-timestamp': string; 'by-sender': string };
  };
  'pools': {
    key: string;
    value: PoolInfo;
    indexes: { 'by-type': string; 'by-village': string };
  };
  'sync-queue': {
    key: string;
    value: SyncQueueItem;
    indexes: { 'by-priority': number; 'by-created': number };
  };
}

interface PendingOrder {
  id: string;
  data: MoneyOrderFormData;
  timestamp: number;
  status: 'pending' | 'syncing' | 'synced' | 'failed';
  retryCount: number;
  error?: string;
  receiptId?: string;
}

interface SyncQueueItem {
  id: string;
  type: 'order' | 'pool-update' | 'receipt-fetch';
  data: any;
  priority: number;
  created: number;
  attempts: number;
  lastAttempt?: number;
  error?: string;
}

export class OfflineService {
  private db: IDBPDatabase<MoneyOrderDB> | null = null;
  private syncInterval: NodeJS.Timer | null = null;
  private isOnline: boolean = navigator.onLine;
  private syncInProgress: boolean = false;

  constructor() {
    this.initializeDB();
    this.setupEventListeners();
    this.startSyncProcess();
  }

  // Initialize IndexedDB
  private async initializeDB() {
    try {
      this.db = await openDB<MoneyOrderDB>('money-order-offline', 1, {
        upgrade(db) {
          // Pending orders store
          if (!db.objectStoreNames.contains('pending-orders')) {
            const orderStore = db.createObjectStore('pending-orders', { keyPath: 'id' });
            orderStore.createIndex('by-timestamp', 'timestamp');
            orderStore.createIndex('by-status', 'status');
          }

          // Receipts store
          if (!db.objectStoreNames.contains('receipts')) {
            const receiptStore = db.createObjectStore('receipts', { keyPath: 'receiptId' });
            receiptStore.createIndex('by-timestamp', 'timestamp');
            receiptStore.createIndex('by-sender', 'sender');
          }

          // Pools store
          if (!db.objectStoreNames.contains('pools')) {
            const poolStore = db.createObjectStore('pools', { keyPath: 'poolId' });
            poolStore.createIndex('by-type', 'type');
            poolStore.createIndex('by-village', 'postalCode');
          }

          // Sync queue store
          if (!db.objectStoreNames.contains('sync-queue')) {
            const syncStore = db.createObjectStore('sync-queue', { keyPath: 'id' });
            syncStore.createIndex('by-priority', 'priority');
            syncStore.createIndex('by-created', 'created');
          }
        }
      });
    } catch (error) {
      console.error('Failed to initialize offline database:', error);
    }
  }

  // Setup online/offline event listeners
  private setupEventListeners() {
    window.addEventListener('online', () => {
      this.isOnline = true;
      console.log('Connection restored - starting sync');
      this.syncPendingOrders();
    });

    window.addEventListener('offline', () => {
      this.isOnline = false;
      console.log('Connection lost - offline mode activated');
    });

    // Service Worker message listener
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.addEventListener('message', (event) => {
        if (event.data.type === 'sync-complete') {
          this.handleSyncComplete(event.data.orderId);
        }
      });
    }
  }

  // Start periodic sync process
  private startSyncProcess() {
    // Sync every 30 seconds when online
    this.syncInterval = setInterval(() => {
      if (this.isOnline && !this.syncInProgress) {
        this.syncPendingOrders();
      }
    }, 30000);
  }

  // Create offline money order
  public async createOfflineOrder(orderData: MoneyOrderFormData): Promise<string> {
    if (!this.db) {
      throw new Error('Offline database not initialized');
    }

    const orderId = this.generateOfflineOrderId();
    const pendingOrder: PendingOrder = {
      id: orderId,
      data: orderData,
      timestamp: Date.now(),
      status: 'pending',
      retryCount: 0
    };

    // Store in IndexedDB
    await this.db.put('pending-orders', pendingOrder);

    // Add to sync queue
    await this.addToSyncQueue({
      id: `sync-${orderId}`,
      type: 'order',
      data: pendingOrder,
      priority: 1,
      created: Date.now(),
      attempts: 0
    });

    // Try immediate sync if online
    if (this.isOnline) {
      this.syncOrder(orderId);
    }

    return orderId;
  }

  // Sync pending orders
  private async syncPendingOrders() {
    if (!this.db || this.syncInProgress) return;

    this.syncInProgress = true;

    try {
      // Get all pending orders
      const pendingOrders = await this.db.getAllFromIndex(
        'pending-orders',
        'by-status',
        'pending'
      );

      console.log(`Found ${pendingOrders.length} pending orders to sync`);

      // Sync each order
      for (const order of pendingOrders) {
        await this.syncOrder(order.id);
      }

      // Process sync queue
      await this.processSyncQueue();

    } catch (error) {
      console.error('Sync process failed:', error);
    } finally {
      this.syncInProgress = false;
    }
  }

  // Sync individual order
  private async syncOrder(orderId: string): Promise<void> {
    if (!this.db) return;

    try {
      const order = await this.db.get('pending-orders', orderId);
      if (!order || order.status !== 'pending') return;

      // Update status
      order.status = 'syncing';
      await this.db.put('pending-orders', order);

      // Send to server
      const response = await fetch('/api/money-orders', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Offline-Sync': 'true',
          'X-Offline-Order-Id': orderId
        },
        body: JSON.stringify(order.data)
      });

      if (!response.ok) {
        throw new Error(`Sync failed: ${response.statusText}`);
      }

      const result = await response.json();

      // Update order with receipt
      order.status = 'synced';
      order.receiptId = result.receipt.receiptId;
      await this.db.put('pending-orders', order);

      // Store receipt
      await this.db.put('receipts', result.receipt);

      // Notify UI
      this.notifyOrderSynced(orderId, result.receipt);

    } catch (error) {
      console.error(`Failed to sync order ${orderId}:`, error);
      
      if (!this.db) return;
      
      const order = await this.db.get('pending-orders', orderId);
      if (order) {
        order.status = 'pending';
        order.retryCount++;
        order.error = error instanceof Error ? error.message : 'Unknown error';
        await this.db.put('pending-orders', order);
      }
    }
  }

  // Process sync queue
  private async processSyncQueue() {
    if (!this.db) return;

    const items = await this.db.getAllFromIndex(
      'sync-queue',
      'by-priority'
    );

    for (const item of items) {
      if (item.attempts > 3) {
        // Remove after too many attempts
        await this.db.delete('sync-queue', item.id);
        continue;
      }

      try {
        await this.processSyncItem(item);
        await this.db.delete('sync-queue', item.id);
      } catch (error) {
        item.attempts++;
        item.lastAttempt = Date.now();
        item.error = error instanceof Error ? error.message : 'Unknown error';
        await this.db.put('sync-queue', item);
      }
    }
  }

  // Process individual sync item
  private async processSyncItem(item: SyncQueueItem): Promise<void> {
    switch (item.type) {
      case 'order':
        await this.syncOrder(item.data.id);
        break;
      case 'pool-update':
        await this.syncPoolUpdate(item.data);
        break;
      case 'receipt-fetch':
        await this.fetchReceipt(item.data.receiptId);
        break;
    }
  }

  // Cache pool data for offline use
  public async cachePoolData(pools: PoolInfo[]): Promise<void> {
    if (!this.db) return;

    const tx = this.db.transaction('pools', 'readwrite');
    await Promise.all(pools.map(pool => tx.store.put(pool)));
    await tx.done;
  }

  // Get cached pools
  public async getCachedPools(type?: string): Promise<PoolInfo[]> {
    if (!this.db) return [];

    if (type) {
      return await this.db.getAllFromIndex('pools', 'by-type', type);
    }
    return await this.db.getAll('pools');
  }

  // Get pending orders
  public async getPendingOrders(): Promise<PendingOrder[]> {
    if (!this.db) return [];
    return await this.db.getAllFromIndex('pending-orders', 'by-timestamp');
  }

  // Get order status
  public async getOrderStatus(orderId: string): Promise<PendingOrder | undefined> {
    if (!this.db) return undefined;
    return await this.db.get('pending-orders', orderId);
  }

  // Get cached receipts
  public async getCachedReceipts(sender?: string): Promise<ReceiptData[]> {
    if (!this.db) return [];

    if (sender) {
      return await this.db.getAllFromIndex('receipts', 'by-sender', sender);
    }
    return await this.db.getAll('receipts');
  }

  // Clear old data
  public async clearOldData(daysToKeep: number = 30): Promise<void> {
    if (!this.db) return;

    const cutoffTime = Date.now() - (daysToKeep * 24 * 60 * 60 * 1000);

    // Clear old pending orders
    const oldOrders = await this.db.getAllFromIndex('pending-orders', 'by-timestamp');
    for (const order of oldOrders) {
      if (order.timestamp < cutoffTime && order.status === 'synced') {
        await this.db.delete('pending-orders', order.id);
      }
    }

    // Clear old receipts
    const receipts = await this.db.getAll('receipts');
    for (const receipt of receipts) {
      const timestamp = new Date(receipt.timestamp).getTime();
      if (timestamp < cutoffTime) {
        await this.db.delete('receipts', receipt.receiptId);
      }
    }
  }

  // Helper methods

  private generateOfflineOrderId(): string {
    return `offline-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
  }

  private async addToSyncQueue(item: SyncQueueItem): Promise<void> {
    if (!this.db) return;
    await this.db.put('sync-queue', item);
  }

  private async syncPoolUpdate(poolData: any): Promise<void> {
    // Implement pool update sync
    console.log('Syncing pool update:', poolData);
  }

  private async fetchReceipt(receiptId: string): Promise<void> {
    // Implement receipt fetching
    console.log('Fetching receipt:', receiptId);
  }

  private notifyOrderSynced(orderId: string, receipt: ReceiptData) {
    // Notify UI components
    window.dispatchEvent(new CustomEvent('order-synced', {
      detail: { orderId, receipt }
    }));

    // Show notification if supported
    if ('Notification' in window && Notification.permission === 'granted') {
      new Notification('Money Order Synced', {
        body: `Order ${orderId} has been successfully synced`,
        icon: '/icon-192x192.png'
      });
    }
  }

  private handleSyncComplete(orderId: string) {
    console.log(`Sync complete for order: ${orderId}`);
    // Update UI or trigger callbacks
  }

  // Get sync status
  public getSyncStatus(): {
    isOnline: boolean;
    syncInProgress: boolean;
    pendingCount: number;
  } {
    return {
      isOnline: this.isOnline,
      syncInProgress: this.syncInProgress,
      pendingCount: 0 // Would need to track this
    };
  }

  // Force sync
  public async forceSync(): Promise<void> {
    if (this.isOnline && !this.syncInProgress) {
      await this.syncPendingOrders();
    }
  }

  // Export data for backup
  public async exportData(): Promise<string> {
    if (!this.db) throw new Error('Database not initialized');

    const data = {
      pendingOrders: await this.db.getAll('pending-orders'),
      receipts: await this.db.getAll('receipts'),
      pools: await this.db.getAll('pools'),
      exportDate: new Date().toISOString()
    };

    return JSON.stringify(data, null, 2);
  }

  // Import backup data
  public async importData(jsonData: string): Promise<void> {
    if (!this.db) throw new Error('Database not initialized');

    try {
      const data = JSON.parse(jsonData);

      // Import pending orders
      if (data.pendingOrders) {
        for (const order of data.pendingOrders) {
          await this.db.put('pending-orders', order);
        }
      }

      // Import receipts
      if (data.receipts) {
        for (const receipt of data.receipts) {
          await this.db.put('receipts', receipt);
        }
      }

      // Import pools
      if (data.pools) {
        for (const pool of data.pools) {
          await this.db.put('pools', pool);
        }
      }

    } catch (error) {
      console.error('Failed to import data:', error);
      throw error;
    }
  }

  // Cleanup
  public destroy() {
    if (this.syncInterval) {
      clearInterval(this.syncInterval);
    }
    if (this.db) {
      this.db.close();
    }
  }
}