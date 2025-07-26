const express = require('express');
const cors = require('cors');
const path = require('path');
const fs = require('fs');
const crypto = require('crypto');
const { Pool } = require('pg');
const Redis = require('redis');

const app = express();
const port = process.env.PORT || 4000;

// Configuration
const FAUCET_CONFIG = {
    tokensPerRequest: 1000,
    rateLimitHours: 24,
    maxDailyRequests: 100,
    minAddressLength: 43,
    maxAddressLength: 43,
    addressPrefix: 'desh'
};

// Database connection
const pool = new Pool({
    user: process.env.DB_USER || 'deshchain',
    host: process.env.DB_HOST || 'postgres',
    database: process.env.DB_NAME || 'deshchain_explorer',
    password: process.env.DB_PASSWORD || 'deshchain123',
    port: process.env.DB_PORT || 5432,
});

// Redis connection for rate limiting
const redisClient = Redis.createClient({
    url: process.env.REDIS_URL || 'redis://redis:6379'
});

// Middleware
app.use(cors());
app.use(express.json());

// Rate limiting middleware
const rateLimit = require('express-rate-limit');
const limiter = rateLimit({
    windowMs: 15 * 60 * 1000, // 15 minutes
    max: 10, // Limit each IP to 10 requests per windowMs
    message: { success: false, message: 'Too many requests, please try again later.' }
});

app.use('/api/request', limiter);

// Connect to Redis
redisClient.connect().catch(console.error);

// Initialize database schema
async function initDatabase() {
    try {
        await pool.query(`
            CREATE TABLE IF NOT EXISTS faucet_requests (
                id SERIAL PRIMARY KEY,
                address VARCHAR(100) NOT NULL,
                amount DECIMAL(20, 8) NOT NULL,
                tx_hash VARCHAR(64) UNIQUE,
                ip_address INET,
                user_agent TEXT,
                status VARCHAR(20) DEFAULT 'pending',
                created_at TIMESTAMP DEFAULT NOW(),
                processed_at TIMESTAMP,
                error_message TEXT
            );
        `);

        await pool.query(`
            CREATE INDEX IF NOT EXISTS idx_faucet_requests_address ON faucet_requests(address);
        `);

        await pool.query(`
            CREATE INDEX IF NOT EXISTS idx_faucet_requests_created_at ON faucet_requests(created_at);
        `);

        await pool.query(`
            CREATE INDEX IF NOT EXISTS idx_faucet_requests_ip ON faucet_requests(ip_address);
        `);

        await pool.query(`
            CREATE TABLE IF NOT EXISTS faucet_stats (
                id SERIAL PRIMARY KEY,
                total_requests BIGINT DEFAULT 0,
                total_distributed DECIMAL(20, 8) DEFAULT 0,
                daily_requests INTEGER DEFAULT 0,
                last_reset_date DATE DEFAULT CURRENT_DATE,
                updated_at TIMESTAMP DEFAULT NOW()
            );
        `);

        // Insert initial stats if not exists
        await pool.query(`
            INSERT INTO faucet_stats (total_requests, total_distributed)
            SELECT 0, 0
            WHERE NOT EXISTS (SELECT 1 FROM faucet_stats)
        `);

        console.log('Faucet database initialized successfully');
    } catch (error) {
        console.error('Database initialization error:', error);
    }
}

// Validate DeshChain address
function validateAddress(address) {
    if (!address || typeof address !== 'string') {
        return { valid: false, error: 'Address is required' };
    }

    if (address.length !== FAUCET_CONFIG.maxAddressLength) {
        return { valid: false, error: `Address must be ${FAUCET_CONFIG.maxAddressLength} characters long` };
    }

    if (!address.startsWith(FAUCET_CONFIG.addressPrefix)) {
        return { valid: false, error: `Address must start with "${FAUCET_CONFIG.addressPrefix}"` };
    }

    // Additional validation for character set
    const validChars = /^[a-zA-Z0-9]+$/;
    if (!validChars.test(address.slice(FAUCET_CONFIG.addressPrefix.length))) {
        return { valid: false, error: 'Address contains invalid characters' };
    }

    return { valid: true };
}

// Check rate limiting
async function checkRateLimit(address, ipAddress) {
    const rateLimitMs = FAUCET_CONFIG.rateLimitHours * 60 * 60 * 1000;
    const now = new Date();
    const cutoffTime = new Date(now.getTime() - rateLimitMs);

    try {
        // Check address-based rate limit
        const addressResult = await pool.query(`
            SELECT COUNT(*) as count, MAX(created_at) as last_request
            FROM faucet_requests 
            WHERE address = $1 AND created_at > $2 AND status = 'completed'
        `, [address, cutoffTime]);

        if (parseInt(addressResult.rows[0].count) > 0) {
            const lastRequest = new Date(addressResult.rows[0].last_request);
            const timeRemaining = rateLimitMs - (now.getTime() - lastRequest.getTime());
            const hoursRemaining = Math.ceil(timeRemaining / (60 * 60 * 1000));
            
            return {
                allowed: false,
                error: `Rate limit exceeded. You can request tokens again in ${hoursRemaining} hours.`,
                timeRemaining: timeRemaining
            };
        }

        // Check IP-based rate limit (more lenient)
        const ipResult = await pool.query(`
            SELECT COUNT(*) as count
            FROM faucet_requests 
            WHERE ip_address = $1 AND created_at > $2 AND status = 'completed'
        `, [ipAddress, cutoffTime]);

        if (parseInt(ipResult.rows[0].count) >= 5) { // Max 5 requests per IP per day
            return {
                allowed: false,
                error: 'IP address rate limit exceeded. Too many requests from this IP.',
                timeRemaining: rateLimitMs
            };
        }

        return { allowed: true };
    } catch (error) {
        console.error('Rate limit check error:', error);
        return { allowed: false, error: 'Rate limit check failed' };
    }
}

// Generate transaction hash (mock)
function generateTxHash() {
    return '0x' + crypto.randomBytes(32).toString('hex');
}

// Simulate token distribution
async function distributeTokens(address, amount, requestId) {
    // In a real implementation, this would interact with the blockchain
    // For now, we'll simulate the process
    
    return new Promise((resolve) => {
        setTimeout(async () => {
            try {
                const txHash = generateTxHash();
                
                // Update request status
                await pool.query(`
                    UPDATE faucet_requests 
                    SET status = 'completed', tx_hash = $1, processed_at = NOW()
                    WHERE id = $2
                `, [txHash, requestId]);

                // Update stats
                await pool.query(`
                    UPDATE faucet_stats 
                    SET total_requests = total_requests + 1,
                        total_distributed = total_distributed + $1,
                        daily_requests = daily_requests + 1,
                        updated_at = NOW()
                `, [amount]);

                // Insert mock transaction into transactions table
                await pool.query(`
                    INSERT INTO transactions (hash, from_address, to_address, amount, status, timestamp, tx_type)
                    VALUES ($1, $2, $3, $4, 'success', NOW(), 'faucet')
                    ON CONFLICT (hash) DO NOTHING
                `, [txHash, 'faucet', address, amount]);

                resolve({
                    success: true,
                    txHash: txHash,
                    amount: amount
                });
            } catch (error) {
                console.error('Token distribution error:', error);
                
                // Update request status to failed
                await pool.query(`
                    UPDATE faucet_requests 
                    SET status = 'failed', error_message = $1, processed_at = NOW()
                    WHERE id = $2
                `, [error.message, requestId]);

                resolve({
                    success: false,
                    error: 'Token distribution failed'
                });
            }
        }, Math.random() * 3000 + 2000); // Simulate 2-5 second processing time
    });
}

// API Routes

// Get faucet statistics
app.get('/api/stats', async (req, res) => {
    try {
        const cacheKey = 'faucet:stats';
        const cached = await redisClient.get(cacheKey);
        
        if (cached) {
            return res.json(JSON.parse(cached));
        }

        const result = await pool.query('SELECT * FROM faucet_stats ORDER BY id DESC LIMIT 1');
        const stats = result.rows[0] || {
            total_requests: 0,
            total_distributed: 0,
            daily_requests: 0
        };

        const responseData = {
            tokensPerRequest: FAUCET_CONFIG.tokensPerRequest,
            rateLimitHours: FAUCET_CONFIG.rateLimitHours,
            totalRequests: parseInt(stats.total_requests),
            totalDistributed: parseFloat(stats.total_distributed),
            dailyRequests: parseInt(stats.daily_requests),
            maxDailyRequests: FAUCET_CONFIG.maxDailyRequests
        };

        await redisClient.setEx(cacheKey, 60, JSON.stringify(responseData)); // Cache for 1 minute
        res.json(responseData);
    } catch (error) {
        console.error('Error fetching faucet stats:', error);
        res.status(500).json({ success: false, message: 'Failed to fetch statistics' });
    }
});

// Get recent faucet transactions
app.get('/api/recent', async (req, res) => {
    try {
        const limit = Math.min(parseInt(req.query.limit) || 10, 50);
        
        const result = await pool.query(`
            SELECT address, amount, tx_hash as hash, created_at as timestamp, status
            FROM faucet_requests 
            WHERE status = 'completed' AND tx_hash IS NOT NULL
            ORDER BY created_at DESC 
            LIMIT $1
        `, [limit]);

        const transactions = result.rows.map(tx => ({
            ...tx,
            amount: parseFloat(tx.amount),
            address: tx.address.substring(0, 8) + '...' + tx.address.substring(tx.address.length - 8) // Anonymize
        }));

        res.json(transactions);
    } catch (error) {
        console.error('Error fetching recent transactions:', error);
        res.status(500).json({ success: false, message: 'Failed to fetch recent transactions' });
    }
});

// Request tokens
app.post('/api/request', async (req, res) => {
    const { address } = req.body;
    const ipAddress = req.ip || req.connection.remoteAddress;
    const userAgent = req.get('User-Agent') || '';

    try {
        // Validate address
        const validation = validateAddress(address);
        if (!validation.valid) {
            return res.status(400).json({
                success: false,
                message: validation.error
            });
        }

        // Check rate limiting
        const rateLimitCheck = await checkRateLimit(address, ipAddress);
        if (!rateLimitCheck.allowed) {
            return res.status(429).json({
                success: false,
                message: rateLimitCheck.error,
                timeRemaining: rateLimitCheck.timeRemaining
            });
        }

        // Check daily limit
        const dailyStats = await pool.query(`
            SELECT daily_requests, last_reset_date FROM faucet_stats ORDER BY id DESC LIMIT 1
        `);

        if (dailyStats.rows.length > 0) {
            const stats = dailyStats.rows[0];
            const today = new Date().toISOString().split('T')[0];
            
            if (stats.last_reset_date !== today) {
                // Reset daily counter
                await pool.query(`
                    UPDATE faucet_stats 
                    SET daily_requests = 0, last_reset_date = CURRENT_DATE
                `);
            } else if (parseInt(stats.daily_requests) >= FAUCET_CONFIG.maxDailyRequests) {
                return res.status(429).json({
                    success: false,
                    message: 'Daily request limit reached. Please try again tomorrow.'
                });
            }
        }

        // Create faucet request record
        const insertResult = await pool.query(`
            INSERT INTO faucet_requests (address, amount, ip_address, user_agent, status)
            VALUES ($1, $2, $3, $4, 'pending')
            RETURNING id
        `, [address, FAUCET_CONFIG.tokensPerRequest, ipAddress, userAgent]);

        const requestId = insertResult.rows[0].id;

        // Process token distribution asynchronously
        const distributionResult = await distributeTokens(address, FAUCET_CONFIG.tokensPerRequest, requestId);

        if (distributionResult.success) {
            res.json({
                success: true,
                message: `Successfully sent ${FAUCET_CONFIG.tokensPerRequest} NAMO tokens to ${address}`,
                amount: `${FAUCET_CONFIG.tokensPerRequest} NAMO`,
                txHash: distributionResult.txHash,
                explorerUrl: `https://explorer.deshchain.com/tx/${distributionResult.txHash}`
            });
        } else {
            res.status(500).json({
                success: false,
                message: distributionResult.error || 'Failed to distribute tokens'
            });
        }

    } catch (error) {
        console.error('Faucet request error:', error);
        res.status(500).json({
            success: false,
            message: 'Internal server error. Please try again later.'
        });
    }
});

// Get request status
app.get('/api/request/:id', async (req, res) => {
    try {
        const requestId = parseInt(req.params.id);
        
        const result = await pool.query(`
            SELECT id, address, amount, tx_hash, status, created_at, processed_at, error_message
            FROM faucet_requests 
            WHERE id = $1
        `, [requestId]);

        if (result.rows.length === 0) {
            return res.status(404).json({
                success: false,
                message: 'Request not found'
            });
        }

        const request = result.rows[0];
        res.json({
            success: true,
            request: {
                ...request,
                amount: parseFloat(request.amount)
            }
        });
    } catch (error) {
        console.error('Error fetching request status:', error);
        res.status(500).json({
            success: false,
            message: 'Failed to fetch request status'
        });
    }
});

// Health check
app.get('/health', (req, res) => {
    res.json({ 
        status: 'healthy', 
        timestamp: new Date().toISOString(),
        version: '1.0.0'
    });
});

// Serve static files
app.use(express.static('/app/frontend'));

// Serve the frontend for all non-API routes
app.get('*', (req, res) => {
    // Only serve HTML for non-API routes
    if (!req.path.startsWith('/api/')) {
        try {
            const htmlContent = fs.readFileSync('/app/frontend/index.html', 'utf8');
            res.setHeader('Content-Type', 'text/html');
            res.send(htmlContent);
        } catch (error) {
            console.error('Error serving HTML:', error);
            res.status(500).send('Error loading page');
        }
    } else {
        res.status(404).json({ success: false, message: 'API endpoint not found' });
    }
});

// Error handling middleware
app.use((error, req, res, next) => {
    console.error('Unhandled error:', error);
    res.status(500).json({ 
        success: false, 
        message: 'Internal server error' 
    });
});

// Background tasks
async function backgroundTasks() {
    try {
        // Clean up old failed requests (older than 7 days)
        await pool.query(`
            DELETE FROM faucet_requests 
            WHERE status = 'failed' AND created_at < NOW() - INTERVAL '7 days'
        `);

        // Update daily stats if needed
        const today = new Date().toISOString().split('T')[0];
        await pool.query(`
            UPDATE faucet_stats 
            SET daily_requests = 0, last_reset_date = CURRENT_DATE
            WHERE last_reset_date < CURRENT_DATE
        `);

    } catch (error) {
        console.error('Background task error:', error);
    }
}

// Run background tasks every hour
setInterval(backgroundTasks, 60 * 60 * 1000);

// Start server
app.listen(port, async () => {
    console.log(`DeshChain Faucet running on port ${port}`);
    await initDatabase();
    
    // Run background tasks on startup
    setTimeout(backgroundTasks, 5000);
});

// Graceful shutdown
process.on('SIGINT', async () => {
    console.log('Shutting down gracefully...');
    await pool.end();
    await redisClient.quit();
    process.exit(0);
});