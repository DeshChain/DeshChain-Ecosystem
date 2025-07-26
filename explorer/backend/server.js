const express = require('express');
const cors = require('cors');
const path = require('path');
const { Pool } = require('pg');
const Redis = require('redis');

const app = express();
const port = process.env.PORT || 3001;

// Database connections
const pool = new Pool({
    user: process.env.DB_USER || 'deshchain',
    host: process.env.DB_HOST || 'postgres',
    database: process.env.DB_NAME || 'deshchain_explorer',
    password: process.env.DB_PASSWORD || 'deshchain123',
    port: process.env.DB_PORT || 5432,
});

const redisClient = Redis.createClient({
    url: process.env.REDIS_URL || 'redis://redis:6379'
});

// Middleware
app.use(cors({
    origin: ['https://deshchain.com', 'https://www.deshchain.com', 'https://explorer.deshchain.com', 'http://localhost:3000', 'http://localhost'],
    credentials: true,
    methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
    allowedHeaders: ['Content-Type', 'Authorization', 'X-Requested-With']
}));
app.use(express.json());
app.use(express.static(path.join(__dirname, '../frontend')));

// Connect to Redis
redisClient.connect().catch(console.error);

// Initialize database schema
async function initDatabase() {
    try {
        await pool.query(`
            CREATE TABLE IF NOT EXISTS blocks (
                id SERIAL PRIMARY KEY,
                height BIGINT UNIQUE NOT NULL,
                hash VARCHAR(64) NOT NULL,
                prev_hash VARCHAR(64),
                timestamp TIMESTAMP NOT NULL,
                validator VARCHAR(100),
                tx_count INTEGER DEFAULT 0,
                size_bytes INTEGER DEFAULT 0,
                gas_used BIGINT DEFAULT 0,
                gas_limit BIGINT DEFAULT 0,
                created_at TIMESTAMP DEFAULT NOW()
            );
        `);

        await pool.query(`
            CREATE TABLE IF NOT EXISTS transactions (
                id SERIAL PRIMARY KEY,
                hash VARCHAR(64) UNIQUE NOT NULL,
                block_height BIGINT REFERENCES blocks(height),
                tx_index INTEGER,
                from_address VARCHAR(100),
                to_address VARCHAR(100),
                amount DECIMAL(20, 8) DEFAULT 0,
                fee DECIMAL(20, 8) DEFAULT 0,
                gas_used INTEGER DEFAULT 0,
                gas_price DECIMAL(20, 8) DEFAULT 0,
                status VARCHAR(20) DEFAULT 'success',
                timestamp TIMESTAMP NOT NULL,
                memo TEXT,
                tx_type VARCHAR(50),
                created_at TIMESTAMP DEFAULT NOW()
            );
        `);

        await pool.query(`
            CREATE TABLE IF NOT EXISTS validators (
                id SERIAL PRIMARY KEY,
                address VARCHAR(100) UNIQUE NOT NULL,
                name VARCHAR(200) NOT NULL,
                website VARCHAR(500),
                description TEXT,
                voting_power DECIMAL(10, 4) DEFAULT 0,
                commission DECIMAL(5, 2) DEFAULT 0,
                uptime DECIMAL(5, 2) DEFAULT 100,
                delegators INTEGER DEFAULT 0,
                self_bonded DECIMAL(20, 8) DEFAULT 0,
                total_bonded DECIMAL(20, 8) DEFAULT 0,
                jailed BOOLEAN DEFAULT FALSE,
                active BOOLEAN DEFAULT TRUE,
                last_updated TIMESTAMP DEFAULT NOW()
            );
        `);

        await pool.query(`
            CREATE TABLE IF NOT EXISTS chain_stats (
                id SERIAL PRIMARY KEY,
                latest_block BIGINT DEFAULT 0,
                total_transactions BIGINT DEFAULT 0,
                active_validators INTEGER DEFAULT 21,
                total_supply DECIMAL(20, 8) DEFAULT 100000000,
                market_cap DECIMAL(20, 2) DEFAULT 0,
                namo_price DECIMAL(10, 4) DEFAULT 2.50,
                updated_at TIMESTAMP DEFAULT NOW()
            );
        `);

        // Insert initial data
        await insertSampleData();
        console.log('Database initialized successfully');
    } catch (error) {
        console.error('Database initialization error:', error);
    }
}

// Insert sample data for demonstration
async function insertSampleData() {
    try {
        // Check if we already have data
        const blockCount = await pool.query('SELECT COUNT(*) FROM blocks');
        if (parseInt(blockCount.rows[0].count) > 0) return;

        // Insert sample validators
        const validators = [
            { name: 'Bharatmata Validator', address: 'deshval1bharatmata123456789', commission: 5.0, voting_power: 15.5 },
            { name: 'Akhand Bharat Node', address: 'deshval1akhandbharat123456', commission: 3.0, voting_power: 12.8 },
            { name: 'Sanatan Dharma Validator', address: 'deshval1sanatandharma1234', commission: 4.5, voting_power: 11.2 },
            { name: 'Vande Mataram Node', address: 'deshval1vandemataram12345', commission: 2.5, voting_power: 10.8 },
            { name: 'Jai Hind Validator', address: 'deshval1jaihind1234567890', commission: 3.5, voting_power: 9.5 },
            { name: 'Bharat Mata Ki Jai', address: 'deshval1bharatmatakijai123', commission: 4.0, voting_power: 8.7 },
            { name: 'Swaraj Node', address: 'deshval1swaraj12345678901', commission: 3.8, voting_power: 7.9 },
            { name: 'Azadi Validator', address: 'deshval1azadi123456789012', commission: 4.2, voting_power: 6.8 },
            { name: 'Swadeshi Node', address: 'deshval1swadeshi1234567890', commission: 3.2, voting_power: 5.9 },
            { name: 'Rashtra Validator', address: 'deshval1rashtra123456789', commission: 3.7, voting_power: 5.2 },
            { name: 'Dharti Node', address: 'deshval1dharti12345678901', commission: 4.1, voting_power: 4.8 },
            { name: 'Ganga Validator', address: 'deshval1ganga123456789012', commission: 3.9, voting_power: 4.3 },
            { name: 'Himalaya Node', address: 'deshval1himalaya1234567890', commission: 3.4, voting_power: 3.9 },
            { name: 'Saraswati Validator', address: 'deshval1saraswati123456789', commission: 3.6, voting_power: 3.5 },
            { name: 'Yamuna Node', address: 'deshval1yamuna12345678901', commission: 4.3, voting_power: 3.2 },
            { name: 'Krishna Validator', address: 'deshval1krishna123456789012', commission: 3.1, voting_power: 2.8 },
            { name: 'Rama Node', address: 'deshval1rama1234567890123', commission: 3.8, voting_power: 2.5 },
            { name: 'Hanuman Validator', address: 'deshval1hanuman12345678901', commission: 4.0, voting_power: 2.2 },
            { name: 'Shiva Node', address: 'deshval1shiva123456789012', commission: 3.3, voting_power: 1.9 },
            { name: 'Durga Validator', address: 'deshval1durga1234567890123', commission: 3.7, voting_power: 1.6 },
            { name: 'Lakshmi Node', address: 'deshval1lakshmi12345678901', commission: 4.4, voting_power: 1.3 }
        ];

        for (const val of validators) {
            await pool.query(`
                INSERT INTO validators (name, address, commission, voting_power, delegators, uptime, active)
                VALUES ($1, $2, $3, $4, $5, $6, $7)
                ON CONFLICT (address) DO NOTHING
            `, [val.name, val.address.substring(0, 50), val.commission, val.voting_power, 
                Math.floor(Math.random() * 1000) + 50, 
                95 + Math.random() * 5, true]);
        }

        // Insert sample blocks and transactions
        const currentTime = new Date();
        for (let i = 1; i <= 50; i++) {
            const blockTime = new Date(currentTime.getTime() - (50 - i) * 6000); // 6 seconds per block
            const txCount = Math.floor(Math.random() * 20) + 1;
            
            await pool.query(`
                INSERT INTO blocks (height, hash, prev_hash, timestamp, validator, tx_count, size_bytes, gas_used, gas_limit)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
                ON CONFLICT (height) DO NOTHING
            `, [
                i,
                generateHash(),
                i > 1 ? generateHash() : '0x0000000000000000000000000000000000000000000000000000000000000000',
                blockTime,
                validators[i % validators.length].name,
                txCount,
                1024 + Math.floor(Math.random() * 2048),
                Math.floor(Math.random() * 1000000),
                2000000
            ]);

            // Insert transactions for this block
            for (let j = 0; j < txCount; j++) {
                await pool.query(`
                    INSERT INTO transactions (hash, block_height, tx_index, from_address, to_address, amount, fee, gas_used, status, timestamp, tx_type)
                    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
                    ON CONFLICT (hash) DO NOTHING
                `, [
                    generateHash(),
                    i,
                    j,
                    generateAddress(),
                    generateAddress(),
                    (Math.random() * 1000).toFixed(4),
                    (Math.random() * 0.1).toFixed(6),
                    Math.floor(Math.random() * 50000),
                    Math.random() > 0.05 ? 'success' : 'failed',
                    blockTime,
                    ['transfer', 'delegate', 'vote', 'mint'][Math.floor(Math.random() * 4)]
                ]);
            }
        }

        // Insert initial chain stats
        await pool.query(`
            INSERT INTO chain_stats (latest_block, total_transactions, active_validators, total_supply, namo_price)
            VALUES (50, 500, 21, 100000000, 2.50)
            ON CONFLICT (id) DO UPDATE SET
                latest_block = EXCLUDED.latest_block,
                total_transactions = EXCLUDED.total_transactions,
                updated_at = NOW()
        `);

        console.log('Sample data inserted successfully');
    } catch (error) {
        console.error('Error inserting sample data:', error);
    }
}

// Helper functions
function generateHash() {
    return '0x' + Array.from({length: 64}, () => Math.floor(Math.random() * 16).toString(16)).join('');
}

function generateAddress() {
    return 'desh' + Array.from({length: 39}, () => 
        'abcdefghijklmnopqrstuvwxyz0123456789'[Math.floor(Math.random() * 36)]
    ).join('');
}

// API Routes

// Chain statistics
app.get('/api/stats', async (req, res) => {
    try {
        const cacheKey = 'chain:stats';
        const cached = await redisClient.get(cacheKey);
        
        if (cached) {
            return res.json(JSON.parse(cached));
        }

        const result = await pool.query('SELECT * FROM chain_stats ORDER BY id DESC LIMIT 1');
        const stats = result.rows[0] || {
            latest_block: 0,
            total_transactions: 0,
            active_validators: 21,
            namo_price: 2.50
        };

        // Add real-time updates
        const latestBlock = await pool.query('SELECT MAX(height) as max_height FROM blocks');
        const totalTx = await pool.query('SELECT COUNT(*) as count FROM transactions');
        const activeVals = await pool.query('SELECT COUNT(*) as count FROM validators WHERE active = true');

        const responseData = {
            latestBlock: parseInt(latestBlock.rows[0]?.max_height) || stats.latest_block,
            totalTransactions: parseInt(totalTx.rows[0]?.count) || stats.total_transactions,
            activeValidators: parseInt(activeVals.rows[0]?.count) || stats.active_validators,
            namoPrice: parseFloat(stats.namo_price),
            marketCap: parseFloat(stats.namo_price) * 100000000, // Total supply * price
            totalSupply: 100000000
        };

        await redisClient.setEx(cacheKey, 30, JSON.stringify(responseData)); // Cache for 30 seconds
        res.json(responseData);
    } catch (error) {
        console.error('Error fetching stats:', error);
        res.status(500).json({ error: 'Failed to fetch stats' });
    }
});

// Latest blocks
app.get('/api/blocks', async (req, res) => {
    try {
        const limit = Math.min(parseInt(req.query.limit) || 10, 100);
        const offset = parseInt(req.query.offset) || 0;

        const result = await pool.query(`
            SELECT height, hash, timestamp, validator, tx_count as txs, size_bytes
            FROM blocks 
            ORDER BY height DESC 
            LIMIT $1 OFFSET $2
        `, [limit, offset]);

        const blocks = result.rows.map(block => ({
            ...block,
            time: block.timestamp,
            hash: block.hash,
            txs: block.txs || 0
        }));

        res.json(blocks);
    } catch (error) {
        console.error('Error fetching blocks:', error);
        res.status(500).json({ error: 'Failed to fetch blocks' });
    }
});

// Latest transactions
app.get('/api/transactions', async (req, res) => {
    try {
        const limit = Math.min(parseInt(req.query.limit) || 10, 100);
        const offset = parseInt(req.query.offset) || 0;

        const result = await pool.query(`
            SELECT hash, block_height, from_address as "from", to_address as "to", 
                   amount, fee, status, timestamp as time, tx_type
            FROM transactions 
            ORDER BY block_height DESC, tx_index DESC 
            LIMIT $1 OFFSET $2
        `, [limit, offset]);

        const transactions = result.rows.map(tx => ({
            ...tx,
            amount: parseFloat(tx.amount),
            fee: parseFloat(tx.fee)
        }));

        res.json(transactions);
    } catch (error) {
        console.error('Error fetching transactions:', error);
        res.status(500).json({ error: 'Failed to fetch transactions' });
    }
});

// Validators
app.get('/api/validators', async (req, res) => {
    try {
        const result = await pool.query(`
            SELECT name, address, voting_power, commission, uptime, 
                   delegators, active, total_bonded, jailed
            FROM validators 
            ORDER BY voting_power DESC
        `);

        const validators = result.rows.map(val => ({
            ...val,
            votingPower: parseFloat(val.voting_power),
            commission: parseFloat(val.commission),
            uptime: parseFloat(val.uptime),
            active: val.active && !val.jailed
        }));

        res.json(validators);
    } catch (error) {
        console.error('Error fetching validators:', error);
        res.status(500).json({ error: 'Failed to fetch validators' });
    }
});

// Block by height
app.get('/api/block/:height', async (req, res) => {
    try {
        const height = parseInt(req.params.height);
        
        const blockResult = await pool.query('SELECT * FROM blocks WHERE height = $1', [height]);
        if (blockResult.rows.length === 0) {
            return res.status(404).json({ error: 'Block not found' });
        }

        const block = blockResult.rows[0];
        
        // Get transactions for this block
        const txResult = await pool.query(`
            SELECT * FROM transactions WHERE block_height = $1 ORDER BY tx_index
        `, [height]);

        res.json({
            ...block,
            transactions: txResult.rows
        });
    } catch (error) {
        console.error('Error fetching block:', error);
        res.status(500).json({ error: 'Failed to fetch block' });
    }
});

// Transaction by hash
app.get('/api/transaction/:hash', async (req, res) => {
    try {
        const hash = req.params.hash;
        
        const result = await pool.query('SELECT * FROM transactions WHERE hash = $1', [hash]);
        if (result.rows.length === 0) {
            return res.status(404).json({ error: 'Transaction not found' });
        }

        res.json(result.rows[0]);
    } catch (error) {
        console.error('Error fetching transaction:', error);
        res.status(500).json({ error: 'Failed to fetch transaction' });
    }
});

// Address information
app.get('/api/address/:address', async (req, res) => {
    try {
        const address = req.params.address;
        
        // Get transactions for this address
        const txResult = await pool.query(`
            SELECT * FROM transactions 
            WHERE from_address = $1 OR to_address = $1 
            ORDER BY timestamp DESC 
            LIMIT 100
        `, [address]);

        // Calculate balance (simplified)
        const sent = await pool.query(`
            SELECT COALESCE(SUM(amount + fee), 0) as total 
            FROM transactions 
            WHERE from_address = $1 AND status = 'success'
        `, [address]);

        const received = await pool.query(`
            SELECT COALESCE(SUM(amount), 0) as total 
            FROM transactions 
            WHERE to_address = $1 AND status = 'success'
        `, [address]);

        const balance = parseFloat(received.rows[0].total) - parseFloat(sent.rows[0].total);

        res.json({
            address,
            balance: Math.max(0, balance),
            transactionCount: txResult.rows.length,
            transactions: txResult.rows
        });
    } catch (error) {
        console.error('Error fetching address:', error);
        res.status(500).json({ error: 'Failed to fetch address information' });
    }
});

// Health check
app.get('/health', (req, res) => {
    res.json({ status: 'healthy', timestamp: new Date().toISOString() });
});

// Serve the frontend
app.get('*', (req, res) => {
    res.sendFile(path.join(__dirname, '../frontend/index.html'));
});

// Error handling middleware
app.use((error, req, res, next) => {
    console.error('Unhandled error:', error);
    res.status(500).json({ error: 'Internal server error' });
});

// Start server
app.listen(port, async () => {
    console.log(`DeshChain Explorer Backend running on port ${port}`);
    await initDatabase();
});

// Graceful shutdown
process.on('SIGINT', async () => {
    console.log('Shutting down gracefully...');
    await pool.end();
    await redisClient.quit();
    process.exit(0);
});