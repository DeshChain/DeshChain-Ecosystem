FROM node:18-alpine

WORKDIR /app

# Create a simple React app placeholder
RUN echo '<!DOCTYPE html>\
<html>\
<head>\
<title>DeshChain Explorer</title>\
<style>\
body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }\
.container { max-width: 1200px; margin: 0 auto; }\
.header { background: white; padding: 20px; border-radius: 8px; margin-bottom: 20px; }\
.stats { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; }\
.stat-card { background: white; padding: 20px; border-radius: 8px; }\
.stat-value { font-size: 24px; font-weight: bold; color: #ff6b00; }\
</style>\
</head>\
<body>\
<div class="container">\
<div class="header">\
<h1>DeshChain Explorer</h1>\
<p>Explore the cultural blockchain of India</p>\
</div>\
<div class="stats">\
<div class="stat-card">\
<h3>Latest Block</h3>\
<div class="stat-value">1,000</div>\
</div>\
<div class="stat-card">\
<h3>Total Validators</h3>\
<div class="stat-value">21</div>\
</div>\
<div class="stat-card">\
<h3>NAMO Price</h3>\
<div class="stat-value">₹2.50</div>\
</div>\
<div class="stat-card">\
<h3>Total Transactions</h3>\
<div class="stat-value">10,000+</div>\
</div>\
</div>\
</div>\
</body>\
</html>' > index.html

# Simple HTTP server
RUN npm install -g http-server

EXPOSE 3000

CMD ["http-server", "-p", "3000", "-c-1"]