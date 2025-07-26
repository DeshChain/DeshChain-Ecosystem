import express from 'express';
import cors from 'cors';
import helmet from 'helmet';
import { createServer } from 'http';
import { Server } from 'socket.io';
import { Pool } from 'pg';
import Redis from 'redis';
import { StargateClient } from '@cosmjs/stargate';
import * as cron from 'node-cron';
import winston from 'winston';

// Initialize logger
const logger = winston.createLogger({
  level: 'info',
  format: winston.format.json(),
  transports: [
    new winston.transports.Console({
      format: winston.format.simple(),
    }),
  ],
});

// Initialize Express app
const app = express();
const httpServer = createServer(app);
const io = new Server(httpServer, {
  cors: {
    origin: '*',
    methods: ['GET', 'POST'],
  },
});

// Middleware
app.use(helmet());
app.use(cors());
app.use(express.json());

// Database connection
const pool = new Pool({
  connectionString: process.env.DATABASE_URL,
});

// Redis connection
const redis = Redis.createClient({
  url: process.env.REDIS_URL,
});

// Connect to blockchain
let client: StargateClient;

async function connectToBlockchain() {
  try {
    client = await StargateClient.connect(process.env.RPC_URL || 'http://localhost:26657');
    logger.info('Connected to DeshChain');
  } catch (error) {
    logger.error('Failed to connect to blockchain:', error);
    setTimeout(connectToBlockchain, 5000);
  }
}

// API Routes
app.get('/api/v1/status', async (req, res) => {
  try {
    const chainId = await client.getChainId();
    const height = await client.getHeight();
    const status = await client.getTx(chainId);
    
    res.json({
      chainId,
      height,
      status: 'online',
      timestamp: new Date(),
    });
  } catch (error) {
    res.status(500).json({ error: 'Failed to get status' });
  }
});

app.get('/api/v1/blocks', async (req, res) => {
  try {
    const limit = parseInt(req.query.limit as string) || 20;
    const offset = parseInt(req.query.offset as string) || 0;
    
    const query = 'SELECT * FROM blocks ORDER BY height DESC LIMIT $1 OFFSET $2';
    const result = await pool.query(query, [limit, offset]);
    
    res.json({
      blocks: result.rows,
      total: result.rowCount,
    });
  } catch (error) {
    res.status(500).json({ error: 'Failed to get blocks' });
  }
});

app.get('/api/v1/transactions', async (req, res) => {
  try {
    const limit = parseInt(req.query.limit as string) || 20;
    const offset = parseInt(req.query.offset as string) || 0;
    
    const query = 'SELECT * FROM transactions ORDER BY timestamp DESC LIMIT $1 OFFSET $2';
    const result = await pool.query(query, [limit, offset]);
    
    res.json({
      transactions: result.rows,
      total: result.rowCount,
    });
  } catch (error) {
    res.status(500).json({ error: 'Failed to get transactions' });
  }
});

app.get('/api/v1/validators', async (req, res) => {
  try {
    const query = 'SELECT * FROM validators WHERE active = true ORDER BY voting_power DESC';
    const result = await pool.query(query);
    
    res.json({
      validators: result.rows,
      total: result.rowCount,
    });
  } catch (error) {
    res.status(500).json({ error: 'Failed to get validators' });
  }
});

app.get('/api/v1/address/:address', async (req, res) => {
  try {
    const { address } = req.params;
    
    // Get from cache first
    const cached = await redis.get(`address:${address}`);
    if (cached) {
      return res.json(JSON.parse(cached));
    }
    
    // Get balance from blockchain
    const account = await client.getAccount(address);
    const balance = await client.getAllBalances(address);
    
    const data = {
      address,
      account,
      balance,
      timestamp: new Date(),
    };
    
    // Cache for 30 seconds
    await redis.setex(`address:${address}`, 30, JSON.stringify(data));
    
    res.json(data);
  } catch (error) {
    res.status(500).json({ error: 'Failed to get address info' });
  }
});

// WebSocket for real-time updates
io.on('connection', (socket) => {
  logger.info('New WebSocket connection');
  
  socket.on('subscribe', (channel) => {
    socket.join(channel);
    logger.info(`Socket subscribed to ${channel}`);
  });
  
  socket.on('disconnect', () => {
    logger.info('Socket disconnected');
  });
});

// Block indexer
async function indexBlocks() {
  try {
    const height = await client.getHeight();
    const lastIndexed = await getLastIndexedBlock();
    
    for (let i = lastIndexed + 1; i <= height; i++) {
      const block = await client.getBlock(i);
      await saveBlock(block);
      
      // Emit to WebSocket
      io.to('blocks').emit('new_block', block);
    }
  } catch (error) {
    logger.error('Block indexing error:', error);
  }
}

async function getLastIndexedBlock(): Promise<number> {
  const result = await pool.query('SELECT MAX(height) as max_height FROM blocks');
  return result.rows[0].max_height || 0;
}

async function saveBlock(block: any) {
  const query = `
    INSERT INTO blocks (height, hash, time, proposer, num_txs, total_gas)
    VALUES ($1, $2, $3, $4, $5, $6)
    ON CONFLICT (height) DO NOTHING
  `;
  
  await pool.query(query, [
    block.header.height,
    block.id,
    block.header.time,
    block.header.proposerAddress,
    block.txs.length,
    0, // Calculate total gas
  ]);
}

// Cron jobs
cron.schedule('*/5 * * * * *', indexBlocks); // Every 5 seconds

// Health check endpoint
app.get('/health', (req, res) => {
  res.json({ status: 'ok' });
});

// Start server
const PORT = process.env.PORT || 3001;

async function start() {
  await redis.connect();
  await connectToBlockchain();
  
  httpServer.listen(PORT, () => {
    logger.info(`Explorer backend running on port ${PORT}`);
  });
}

start().catch((error) => {
  logger.error('Failed to start:', error);
  process.exit(1);
});