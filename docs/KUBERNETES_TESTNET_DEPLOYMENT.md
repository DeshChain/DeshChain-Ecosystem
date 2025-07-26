# DeshChain Testnet Kubernetes Deployment Guide

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Prerequisites](#prerequisites)
3. [Infrastructure Setup](#infrastructure-setup)
4. [Kubernetes Cluster Configuration](#kubernetes-cluster-configuration)
5. [DeshChain Node Deployment](#deshchain-node-deployment)
6. [Validator Setup](#validator-setup)
7. [Frontend Deployment](#frontend-deployment)
8. [Monitoring and Observability](#monitoring-and-observability)
9. [Security Configuration](#security-configuration)
10. [Maintenance and Scaling](#maintenance-and-scaling)
11. [Troubleshooting](#troubleshooting)

## Architecture Overview

### DeshChain Testnet Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Load Balancer (Ingress)                  │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │   Frontend  │  │   Explorer  │  │   Faucet    │              │
│  │   (React)   │  │   (React)   │  │   (API)     │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │    API      │  │   gRPC      │  │   WebSocket │              │
│  │  Gateway    │  │   Service   │  │   Service   │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
├─────────────────────────────────────────────────────────────────┤
│                    DeshChain Validator Nodes                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ Validator-1 │  │ Validator-2 │  │ Validator-3 │              │
│  │  (Genesis)  │  │  (Genesis)  │  │  (Genesis)  │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ Validator-4 │  │  Sentry-1   │  │  Sentry-2   │              │
│  │  (Genesis)  │  │   (RPC)     │  │   (RPC)     │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
├─────────────────────────────────────────────────────────────────┤
│                     Storage & Database                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ PostgreSQL  │  │    Redis    │  │    IPFS     │              │
│  │  (Primary)  │  │   (Cache)   │  │ (Storage)   │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
├─────────────────────────────────────────────────────────────────┤
│                   Monitoring & Logging                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ Prometheus  │  │   Grafana   │  │ ELK Stack   │              │
│  │ (Metrics)   │  │(Dashboard)  │  │  (Logs)     │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└─────────────────────────────────────────────────────────────────┘
```

### Component Breakdown

| **Component** | **Replicas** | **Purpose** | **Resources** |
|---------------|--------------|-------------|---------------|
| **Genesis Validators** | 4 | Initial consensus nodes | 4 CPU, 16GB RAM, 1TB SSD |
| **Sentry Nodes** | 2 | Public RPC endpoints | 2 CPU, 8GB RAM, 500GB SSD |
| **API Gateway** | 3 | REST API proxy | 1 CPU, 4GB RAM |
| **Frontend** | 3 | React applications | 0.5 CPU, 2GB RAM |
| **Explorer** | 2 | Blockchain explorer | 1 CPU, 4GB RAM |
| **Faucet** | 2 | Testnet token distribution | 0.5 CPU, 2GB RAM |
| **PostgreSQL** | 1 | Primary database | 2 CPU, 8GB RAM, 200GB SSD |
| **Redis** | 1 | Caching layer | 1 CPU, 4GB RAM |
| **IPFS** | 3 | Distributed storage | 1 CPU, 4GB RAM, 100GB SSD |

## Prerequisites

### Required Tools

```bash
# Install kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Install Helm
curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg > /dev/null
sudo apt-get update && sudo apt-get install helm

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install kind (for local testing)
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
sudo install -o root -g root -m 0755 kind /usr/local/bin/kind
```

### System Requirements

| **Environment** | **Nodes** | **CPU** | **Memory** | **Storage** | **Network** |
|-----------------|-----------|---------|------------|-------------|-------------|
| **Development** | 3 nodes | 8 cores | 32GB RAM | 500GB SSD | 1Gbps |
| **Staging** | 5 nodes | 16 cores | 64GB RAM | 1TB SSD | 10Gbps |
| **Production** | 10+ nodes | 32 cores | 128GB RAM | 2TB NVMe | 25Gbps |

## Infrastructure Setup

### 1. Kubernetes Cluster Setup

#### Option A: Cloud Provider (Recommended)

**AWS EKS Setup:**
```bash
# Create cluster using eksctl
cat << EOF > cluster-config.yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: deshchain-testnet
  region: ap-south-1
  version: "1.28"

nodeGroups:
  - name: validator-nodes
    instanceType: c5.2xlarge
    minSize: 4
    maxSize: 8
    desiredCapacity: 4
    volumeSize: 1000
    volumeType: gp3
    ssh:
      allow: true
    labels:
      node-type: validator
    taints:
      - key: validator
        value: "true"
        effect: NoSchedule

  - name: service-nodes
    instanceType: c5.xlarge
    minSize: 2
    maxSize: 6
    desiredCapacity: 3
    volumeSize: 500
    volumeType: gp3
    ssh:
      allow: true
    labels:
      node-type: service

addons:
  - name: aws-ebs-csi-driver
  - name: aws-efs-csi-driver
  - name: aws-load-balancer-controller
EOF

eksctl create cluster -f cluster-config.yaml
```

**Google GKE Setup:**
```bash
# Create GKE cluster
gcloud container clusters create deshchain-testnet \
  --zone=asia-south1-a \
  --machine-type=c2-standard-8 \
  --num-nodes=4 \
  --disk-size=1000GB \
  --disk-type=pd-ssd \
  --enable-network-policy \
  --enable-ip-alias \
  --enable-autoscaling \
  --min-nodes=4 \
  --max-nodes=10
```

#### Option B: Local Development (Kind)

```bash
# Create local cluster for development
cat << EOF > kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: deshchain-local
nodes:
  - role: control-plane
    kubeadmConfigPatches:
    - |
      kind: InitConfiguration
      nodeRegistration:
        kubeletExtraArgs:
          node-labels: "ingress-ready=true"
    extraPortMappings:
    - containerPort: 80
      hostPort: 80
      protocol: TCP
    - containerPort: 443
      hostPort: 443
      protocol: TCP
  - role: worker
    labels:
      node-type: validator
  - role: worker
    labels:
      node-type: validator
  - role: worker
    labels:
      node-type: service
EOF

kind create cluster --config=kind-config.yaml
```

### 2. Container Images

Create multi-stage Dockerfile for DeshChain:

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o deshchaind ./cmd/deshchaind

# Runtime stage
FROM alpine:3.18

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/deshchaind .
COPY --from=builder /app/config ./config

EXPOSE 26656 26657 26658 9090 1317

CMD ["./deshchaind", "start"]
```

Build and push images:
```bash
# Build DeshChain node image
docker build -t deshchain/node:testnet-v1.0.0 .
docker push deshchain/node:testnet-v1.0.0

# Build frontend image
cd frontend
docker build -t deshchain/frontend:testnet-v1.0.0 .
docker push deshchain/frontend:testnet-v1.0.0
```

## Kubernetes Cluster Configuration

### 1. Namespace Setup

```yaml
# namespaces.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: deshchain-testnet
  labels:
    name: deshchain-testnet
    environment: testnet
---
apiVersion: v1
kind: Namespace
metadata:
  name: deshchain-monitoring
  labels:
    name: deshchain-monitoring
    environment: testnet
---
apiVersion: v1
kind: Namespace
metadata:
  name: deshchain-storage
  labels:
    name: deshchain-storage
    environment: testnet
```

### 2. Storage Classes

```yaml
# storage-classes.yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: validator-storage
  namespace: deshchain-testnet
provisioner: kubernetes.io/aws-ebs # Change for your provider
parameters:
  type: gp3
  iops: "3000"
  throughput: "125"
volumeBindingMode: WaitForFirstConsumer
reclaimPolicy: Retain
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: database-storage
  namespace: deshchain-storage
provisioner: kubernetes.io/aws-ebs
parameters:
  type: io2
  iops: "5000"
volumeBindingMode: WaitForFirstConsumer
reclaimPolicy: Retain
```

### 3. Network Policies

```yaml
# network-policies.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: validator-network-policy
  namespace: deshchain-testnet
spec:
  podSelector:
    matchLabels:
      app: deshchain-validator
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: deshchain-validator
    ports:
    - protocol: TCP
      port: 26656
  - from:
    - podSelector:
        matchLabels:
          app: deshchain-sentry
    ports:
    - protocol: TCP
      port: 26657
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: deshchain-validator
    ports:
    - protocol: TCP
      port: 26656
  - to: []
    ports:
    - protocol: TCP
      port: 53
    - protocol: UDP
      port: 53
```

## DeshChain Node Deployment

### 1. ConfigMaps

```yaml
# configmaps.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: deshchain-genesis
  namespace: deshchain-testnet
data:
  genesis.json: |
    {
      "genesis_time": "2024-01-01T00:00:00.000000000Z",
      "chain_id": "deshchain-testnet-1",
      "initial_height": "1",
      "consensus_params": {
        "block": {
          "max_bytes": "22020096",
          "max_gas": "10000000",
          "time_iota_ms": "1000"
        },
        "evidence": {
          "max_age_num_blocks": "100000",
          "max_age_duration": "172800000000000",
          "max_bytes": "1048576"
        },
        "validator": {
          "pub_key_types": ["ed25519"]
        },
        "version": {}
      },
      "app_hash": "",
      "app_state": {
        "namo": {
          "params": {
            "total_supply": "10000000000000000",
            "initial_price_usd": "0.01"
          },
          "balances": [
            {
              "address": "deshchain1genesis1...",
              "coins": [{"denom": "namo", "amount": "1000000000000"}]
            }
          ]
        },
        "validator": {
          "genesis_nfts": [
            {
              "token_id": "1",
              "rank": 1,
              "english_name": "Narendra Modi",
              "hindi_name": "नरेंद्र मोदी",
              "current_owner": "deshchain1genesis1...",
              "image_uri": "/nfts/1.png",
              "minted_at": "2024-01-01T00:00:00.000000000Z"
            }
          ]
        }
      }
    }
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: deshchain-config
  namespace: deshchain-testnet
data:
  app.toml: |
    minimum-gas-prices = "0.025namo"
    halt-height = 0
    halt-time = 0
    inter-block-cache = true
    index-events = []
    
    [telemetry]
    service-name = "deshchain"
    enabled = true
    enable-hostname = true
    enable-hostname-label = true
    enable-service-label = true
    prometheus-retention-time = 600
    
    [api]
    enable = true
    swagger = true
    address = "tcp://0.0.0.0:1317"
    
    [grpc]
    enable = true
    address = "0.0.0.0:9090"
    
    [grpc-web]
    enable = true
    address = "0.0.0.0:9091"
    
  config.toml: |
    proxy_app = "tcp://127.0.0.1:26658"
    moniker = "deshchain-validator"
    
    [rpc]
    laddr = "tcp://0.0.0.0:26657"
    cors_allowed_origins = ["*"]
    cors_allowed_methods = ["HEAD", "GET", "POST"]
    cors_allowed_headers = ["Origin", "Accept", "Content-Type", "X-Requested-With", "X-Server-Time"]
    
    [p2p]
    laddr = "tcp://0.0.0.0:26656"
    external_address = ""
    persistent_peers = ""
    private_peer_ids = ""
    
    [consensus]
    timeout_propose = "3s"
    timeout_propose_delta = "500ms"
    timeout_prevote = "1s"
    timeout_prevote_delta = "500ms"
    timeout_precommit = "1s"
    timeout_precommit_delta = "500ms"
    timeout_commit = "5s"
```

### 2. Secrets

```yaml
# secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: validator-keys
  namespace: deshchain-testnet
type: Opaque
data:
  # Base64 encoded validator private keys
  priv_validator_key.json: |
    <base64-encoded-validator-key>
  node_key.json: |
    <base64-encoded-node-key>
---
apiVersion: v1
kind: Secret
metadata:
  name: database-credentials
  namespace: deshchain-storage
type: Opaque
data:
  username: <base64-encoded-username>
  password: <base64-encoded-password>
```

### 3. Genesis Validator StatefulSet

```yaml
# validator-statefulset.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: deshchain-validator
  namespace: deshchain-testnet
  labels:
    app: deshchain-validator
    component: validator
spec:
  serviceName: deshchain-validator
  replicas: 4
  selector:
    matchLabels:
      app: deshchain-validator
  template:
    metadata:
      labels:
        app: deshchain-validator
        component: validator
    spec:
      nodeSelector:
        node-type: validator
      tolerations:
      - key: validator
        operator: Equal
        value: "true"
        effect: NoSchedule
      initContainers:
      - name: init-validator
        image: deshchain/node:testnet-v1.0.0
        command:
        - sh
        - -c
        - |
          if [ ! -f /deshchain/config/genesis.json ]; then
            deshchaind init $VALIDATOR_NAME --chain-id deshchain-testnet-1 --home /deshchain
            cp /config/genesis.json /deshchain/config/genesis.json
            cp /config/app.toml /deshchain/config/app.toml
            cp /config/config.toml /deshchain/config/config.toml
          fi
        env:
        - name: VALIDATOR_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        volumeMounts:
        - name: validator-data
          mountPath: /deshchain
        - name: config-volume
          mountPath: /config
      containers:
      - name: deshchain
        image: deshchain/node:testnet-v1.0.0
        command:
        - deshchaind
        - start
        - --home
        - /deshchain
        - --log_level
        - info
        ports:
        - containerPort: 26656
          name: p2p
        - containerPort: 26657
          name: rpc
        - containerPort: 26658
          name: abci
        - containerPort: 9090
          name: grpc
        - containerPort: 1317
          name: api
        env:
        - name: VALIDATOR_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        resources:
          requests:
            memory: "16Gi"
            cpu: "4"
          limits:
            memory: "32Gi"
            cpu: "8"
        volumeMounts:
        - name: validator-data
          mountPath: /deshchain
        - name: validator-keys
          mountPath: /deshchain/config/priv_validator_key.json
          subPath: priv_validator_key.json
        - name: validator-keys
          mountPath: /deshchain/config/node_key.json
          subPath: node_key.json
        livenessProbe:
          httpGet:
            path: /health
            port: 26657
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /status
            port: 26657
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: config-volume
        configMap:
          name: deshchain-config
      - name: validator-keys
        secret:
          secretName: validator-keys
  volumeClaimTemplates:
  - metadata:
      name: validator-data
    spec:
      accessModes: ["ReadWriteOnce"]
      storageClassName: validator-storage
      resources:
        requests:
          storage: 1Ti
```

### 4. Sentry Node Deployment

```yaml
# sentry-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deshchain-sentry
  namespace: deshchain-testnet
  labels:
    app: deshchain-sentry
    component: sentry
spec:
  replicas: 2
  selector:
    matchLabels:
      app: deshchain-sentry
  template:
    metadata:
      labels:
        app: deshchain-sentry
        component: sentry
    spec:
      nodeSelector:
        node-type: service
      initContainers:
      - name: init-sentry
        image: deshchain/node:testnet-v1.0.0
        command:
        - sh
        - -c
        - |
          deshchaind init sentry --chain-id deshchain-testnet-1 --home /deshchain
          cp /config/genesis.json /deshchain/config/genesis.json
          cp /config/app.toml /deshchain/config/app.toml
          cp /config/config.toml /deshchain/config/config.toml
        volumeMounts:
        - name: sentry-data
          mountPath: /deshchain
        - name: config-volume
          mountPath: /config
      containers:
      - name: deshchain
        image: deshchain/node:testnet-v1.0.0
        command:
        - deshchaind
        - start
        - --home
        - /deshchain
        - --log_level
        - info
        ports:
        - containerPort: 26656
          name: p2p
        - containerPort: 26657
          name: rpc
        - containerPort: 9090
          name: grpc
        - containerPort: 1317
          name: api
        resources:
          requests:
            memory: "8Gi"
            cpu: "2"
          limits:
            memory: "16Gi"
            cpu: "4"
        volumeMounts:
        - name: sentry-data
          mountPath: /deshchain
        livenessProbe:
          httpGet:
            path: /health
            port: 26657
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /status
            port: 26657
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: sentry-data
        emptyDir: {}
      - name: config-volume
        configMap:
          name: deshchain-config
```

### 5. Services

```yaml
# services.yaml
apiVersion: v1
kind: Service
metadata:
  name: deshchain-validator-svc
  namespace: deshchain-testnet
  labels:
    app: deshchain-validator
spec:
  selector:
    app: deshchain-validator
  ports:
  - name: p2p
    port: 26656
    targetPort: 26656
  - name: rpc
    port: 26657
    targetPort: 26657
  - name: grpc
    port: 9090
    targetPort: 9090
  - name: api
    port: 1317
    targetPort: 1317
  type: ClusterIP
---
apiVersion: v1
kind: Service
metadata:
  name: deshchain-sentry-svc
  namespace: deshchain-testnet
  labels:
    app: deshchain-sentry
spec:
  selector:
    app: deshchain-sentry
  ports:
  - name: rpc
    port: 26657
    targetPort: 26657
  - name: grpc
    port: 9090
    targetPort: 9090
  - name: api
    port: 1317
    targetPort: 1317
  type: LoadBalancer
---
apiVersion: v1
kind: Service
metadata:
  name: deshchain-p2p-svc
  namespace: deshchain-testnet
  labels:
    app: deshchain-validator
spec:
  selector:
    app: deshchain-validator
  ports:
  - name: p2p
    port: 26656
    targetPort: 26656
    protocol: TCP
  type: LoadBalancer
```

## Validator Setup

### 1. Create Genesis Validators

```bash
#!/bin/bash
# create-genesis-validators.sh

CHAIN_ID="deshchain-testnet-1"
VALIDATORS=("validator-0" "validator-1" "validator-2" "validator-3")
GENESIS_ACCOUNTS=(
  "deshchain1genesis1..." 
  "deshchain1genesis2..." 
  "deshchain1genesis3..." 
  "deshchain1genesis4..."
)

# Create genesis file
deshchaind init genesis --chain-id $CHAIN_ID

# Add genesis accounts
for i in "${!GENESIS_ACCOUNTS[@]}"; do
  deshchaind add-genesis-account ${GENESIS_ACCOUNTS[$i]} 1000000000000namo
done

# Create validator keys
for validator in "${VALIDATORS[@]}"; do
  mkdir -p keys/$validator
  deshchaind init $validator --chain-id $CHAIN_ID --home keys/$validator
  
  # Generate validator transaction
  deshchaind gentx $validator 100000000000namo \
    --chain-id $CHAIN_ID \
    --home keys/$validator \
    --commission-rate 0.05 \
    --commission-max-rate 0.10 \
    --commission-max-change-rate 0.01 \
    --min-self-delegation 100000000000
done

# Collect genesis transactions
deshchaind collect-gentxs

# Validate genesis
deshchaind validate-genesis

echo "Genesis validators created successfully!"
```

### 2. Validator Key Management

```yaml
# validator-key-job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: generate-validator-keys
  namespace: deshchain-testnet
spec:
  template:
    spec:
      containers:
      - name: keygen
        image: deshchain/node:testnet-v1.0.0
        command:
        - sh
        - -c
        - |
          for i in {0..3}; do
            mkdir -p /keys/validator-$i
            deshchaind init validator-$i --chain-id deshchain-testnet-1 --home /keys/validator-$i
            echo "Generated keys for validator-$i"
          done
        volumeMounts:
        - name: key-storage
          mountPath: /keys
      volumes:
      - name: key-storage
        persistentVolumeClaim:
          claimName: validator-keys-pvc
      restartPolicy: Never
```

### 3. Validator Monitoring

```yaml
# validator-monitor.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: validator-monitor
  namespace: deshchain-testnet
spec:
  replicas: 1
  selector:
    matchLabels:
      app: validator-monitor
  template:
    metadata:
      labels:
        app: validator-monitor
    spec:
      containers:
      - name: monitor
        image: prom/node-exporter:latest
        ports:
        - containerPort: 9100
        args:
        - --path.procfs=/host/proc
        - --path.sysfs=/host/sys
        - --collector.filesystem.ignored-mount-points
        - ^/(sys|proc|dev|host|etc|rootfs/var/lib/docker/containers|rootfs/var/lib/docker/overlay2|rootfs/run/docker/netns|rootfs/var/lib/docker/aufs)($$|/)
        volumeMounts:
        - name: proc
          mountPath: /host/proc
          readOnly: true
        - name: sys
          mountPath: /host/sys
          readOnly: true
        - name: rootfs
          mountPath: /rootfs
          readOnly: true
      volumes:
      - name: proc
        hostPath:
          path: /proc
      - name: sys
        hostPath:
          path: /sys
      - name: rootfs
        hostPath:
          path: /
```

## Frontend Deployment

### 1. React Frontend

```yaml
# frontend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deshchain-frontend
  namespace: deshchain-testnet
  labels:
    app: deshchain-frontend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: deshchain-frontend
  template:
    metadata:
      labels:
        app: deshchain-frontend
    spec:
      containers:
      - name: frontend
        image: deshchain/frontend:testnet-v1.0.0
        ports:
        - containerPort: 3000
        env:
        - name: REACT_APP_API_URL
          value: "https://testnet-api.deshchain.com"
        - name: REACT_APP_RPC_URL
          value: "https://testnet-rpc.deshchain.com"
        - name: REACT_APP_CHAIN_ID
          value: "deshchain-testnet-1"
        resources:
          requests:
            memory: "2Gi"
            cpu: "0.5"
          limits:
            memory: "4Gi"
            cpu: "1"
        livenessProbe:
          httpGet:
            path: /
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: deshchain-frontend-svc
  namespace: deshchain-testnet
spec:
  selector:
    app: deshchain-frontend
  ports:
  - port: 80
    targetPort: 3000
  type: ClusterIP
```

### 2. Blockchain Explorer

```yaml
# explorer-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deshchain-explorer
  namespace: deshchain-testnet
  labels:
    app: deshchain-explorer
spec:
  replicas: 2
  selector:
    matchLabels:
      app: deshchain-explorer
  template:
    metadata:
      labels:
        app: deshchain-explorer
    spec:
      containers:
      - name: explorer
        image: deshchain/explorer:testnet-v1.0.0
        ports:
        - containerPort: 3001
        env:
        - name: RPC_ENDPOINT
          value: "http://deshchain-sentry-svc:26657"
        - name: API_ENDPOINT
          value: "http://deshchain-sentry-svc:1317"
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: database-credentials
              key: connection-string
        resources:
          requests:
            memory: "4Gi"
            cpu: "1"
          limits:
            memory: "8Gi"
            cpu: "2"
---
apiVersion: v1
kind: Service
metadata:
  name: deshchain-explorer-svc
  namespace: deshchain-testnet
spec:
  selector:
    app: deshchain-explorer
  ports:
  - port: 80
    targetPort: 3001
  type: ClusterIP
```

### 3. Testnet Faucet

```yaml
# faucet-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deshchain-faucet
  namespace: deshchain-testnet
  labels:
    app: deshchain-faucet
spec:
  replicas: 2
  selector:
    matchLabels:
      app: deshchain-faucet
  template:
    metadata:
      labels:
        app: deshchain-faucet
    spec:
      containers:
      - name: faucet
        image: deshchain/faucet:testnet-v1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: RPC_ENDPOINT
          value: "http://deshchain-sentry-svc:26657"
        - name: FAUCET_MNEMONIC
          valueFrom:
            secretKeyRef:
              name: faucet-keys
              key: mnemonic
        - name: RATE_LIMIT
          value: "1000000namo"
        - name: RATE_LIMIT_WINDOW
          value: "24h"
        resources:
          requests:
            memory: "2Gi"
            cpu: "0.5"
          limits:
            memory: "4Gi"
            cpu: "1"
---
apiVersion: v1
kind: Service
metadata:
  name: deshchain-faucet-svc
  namespace: deshchain-testnet
spec:
  selector:
    app: deshchain-faucet
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

## Monitoring and Observability

### 1. Prometheus Setup

```yaml
# prometheus-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-config
  namespace: deshchain-monitoring
data:
  prometheus.yml: |
    global:
      scrape_interval: 15s
      evaluation_interval: 15s
    
    rule_files:
      - "deshchain_rules.yml"
    
    scrape_configs:
      - job_name: 'deshchain-validators'
        static_configs:
          - targets: ['deshchain-validator-svc.deshchain-testnet:26657']
        metrics_path: /metrics
        scrape_interval: 10s
      
      - job_name: 'deshchain-sentries'
        static_configs:
          - targets: ['deshchain-sentry-svc.deshchain-testnet:26657']
        metrics_path: /metrics
        scrape_interval: 10s
      
      - job_name: 'node-exporter'
        static_configs:
          - targets: ['validator-monitor.deshchain-testnet:9100']
    
    alerting:
      alertmanagers:
        - static_configs:
            - targets: ['alertmanager:9093']

  deshchain_rules.yml: |
    groups:
      - name: deshchain_alerts
        rules:
          - alert: ValidatorDown
            expr: up{job="deshchain-validators"} == 0
            for: 1m
            labels:
              severity: critical
            annotations:
              summary: "DeshChain validator is down"
              description: "Validator {{ $labels.instance }} has been down for more than 1 minute."
          
          - alert: HighBlockTime
            expr: tendermint_consensus_block_interval_seconds > 10
            for: 2m
            labels:
              severity: warning
            annotations:
              summary: "High block time detected"
              description: "Block time is {{ $value }}s on {{ $labels.instance }}"
          
          - alert: LowPeerCount
            expr: tendermint_p2p_peers < 2
            for: 5m
            labels:
              severity: warning
            annotations:
              summary: "Low peer count"
              description: "Node {{ $labels.instance }} has only {{ $value }} peers"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prometheus
  namespace: deshchain-monitoring
spec:
  replicas: 1
  selector:
    matchLabels:
      app: prometheus
  template:
    metadata:
      labels:
        app: prometheus
    spec:
      containers:
      - name: prometheus
        image: prom/prometheus:latest
        ports:
        - containerPort: 9090
        args:
        - --config.file=/etc/prometheus/prometheus.yml
        - --storage.tsdb.path=/prometheus/
        - --web.console.libraries=/etc/prometheus/console_libraries
        - --web.console.templates=/etc/prometheus/consoles
        - --storage.tsdb.retention.time=200h
        - --web.enable-lifecycle
        volumeMounts:
        - name: prometheus-config
          mountPath: /etc/prometheus
        - name: prometheus-storage
          mountPath: /prometheus
        resources:
          requests:
            memory: "4Gi"
            cpu: "1"
          limits:
            memory: "8Gi"
            cpu: "2"
      volumes:
      - name: prometheus-config
        configMap:
          name: prometheus-config
      - name: prometheus-storage
        persistentVolumeClaim:
          claimName: prometheus-storage-pvc
```

### 2. Grafana Dashboard

```yaml
# grafana-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: grafana
  namespace: deshchain-monitoring
spec:
  replicas: 1
  selector:
    matchLabels:
      app: grafana
  template:
    metadata:
      labels:
        app: grafana
    spec:
      containers:
      - name: grafana
        image: grafana/grafana:latest
        ports:
        - containerPort: 3000
        env:
        - name: GF_SECURITY_ADMIN_PASSWORD
          valueFrom:
            secretKeyRef:
              name: grafana-credentials
              key: admin-password
        volumeMounts:
        - name: grafana-storage
          mountPath: /var/lib/grafana
        - name: grafana-config
          mountPath: /etc/grafana/provisioning
        resources:
          requests:
            memory: "2Gi"
            cpu: "0.5"
          limits:
            memory: "4Gi"
            cpu: "1"
      volumes:
      - name: grafana-storage
        persistentVolumeClaim:
          claimName: grafana-storage-pvc
      - name: grafana-config
        configMap:
          name: grafana-config
```

### 3. Logging with ELK Stack

```yaml
# elasticsearch-deployment.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: elasticsearch
  namespace: deshchain-monitoring
spec:
  serviceName: elasticsearch
  replicas: 3
  selector:
    matchLabels:
      app: elasticsearch
  template:
    metadata:
      labels:
        app: elasticsearch
    spec:
      containers:
      - name: elasticsearch
        image: docker.elastic.co/elasticsearch/elasticsearch:8.11.0
        ports:
        - containerPort: 9200
        - containerPort: 9300
        env:
        - name: cluster.name
          value: deshchain-logs
        - name: node.name
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: discovery.seed_hosts
          value: "elasticsearch-0.elasticsearch,elasticsearch-1.elasticsearch,elasticsearch-2.elasticsearch"
        - name: cluster.initial_master_nodes
          value: "elasticsearch-0,elasticsearch-1,elasticsearch-2"
        - name: ES_JAVA_OPTS
          value: "-Xms2g -Xmx2g"
        - name: xpack.security.enabled
          value: "false"
        volumeMounts:
        - name: elasticsearch-data
          mountPath: /usr/share/elasticsearch/data
        resources:
          requests:
            memory: "4Gi"
            cpu: "1"
          limits:
            memory: "8Gi"
            cpu: "2"
  volumeClaimTemplates:
  - metadata:
      name: elasticsearch-data
    spec:
      accessModes: ["ReadWriteOnce"]
      storageClassName: database-storage
      resources:
        requests:
          storage: 100Gi
```

## Security Configuration

### 1. RBAC Configuration

```yaml
# rbac.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: deshchain-service-account
  namespace: deshchain-testnet
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: deshchain-cluster-role
rules:
- apiGroups: [""]
  resources: ["pods", "services", "endpoints"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["networking.k8s.io"]
  resources: ["networkpolicies"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: deshchain-cluster-role-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: deshchain-cluster-role
subjects:
- kind: ServiceAccount
  name: deshchain-service-account
  namespace: deshchain-testnet
```

### 2. Pod Security Policies

```yaml
# pod-security-policy.yaml
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata:
  name: deshchain-psp
spec:
  privileged: false
  allowPrivilegeEscalation: false
  requiredDropCapabilities:
    - ALL
  volumes:
    - 'configMap'
    - 'emptyDir'
    - 'projected'
    - 'secret'
    - 'downwardAPI'
    - 'persistentVolumeClaim'
  runAsUser:
    rule: 'MustRunAsNonRoot'
  seLinux:
    rule: 'RunAsAny'
  fsGroup:
    rule: 'RunAsAny'
```

### 3. Network Security

```yaml
# ingress-nginx.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: deshchain-ingress
  namespace: deshchain-testnet
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/force-ssl-redirect: "true"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/rate-limit: "100"
    nginx.ingress.kubernetes.io/rate-limit-window: "1m"
spec:
  tls:
  - hosts:
    - testnet.deshchain.com
    - testnet-api.deshchain.com
    - testnet-rpc.deshchain.com
    - explorer.testnet.deshchain.com
    - faucet.testnet.deshchain.com
    secretName: deshchain-tls
  rules:
  - host: testnet.deshchain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: deshchain-frontend-svc
            port:
              number: 80
  - host: testnet-api.deshchain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: deshchain-sentry-svc
            port:
              number: 1317
  - host: testnet-rpc.deshchain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: deshchain-sentry-svc
            port:
              number: 26657
  - host: explorer.testnet.deshchain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: deshchain-explorer-svc
            port:
              number: 80
  - host: faucet.testnet.deshchain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: deshchain-faucet-svc
            port:
              number: 80
```

## Maintenance and Scaling

### 1. Backup Strategy

```yaml
# backup-cronjob.yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: deshchain-backup
  namespace: deshchain-testnet
spec:
  schedule: "0 2 * * *"  # Daily at 2 AM
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: deshchain/backup:latest
            command:
            - /bin/bash
            - -c
            - |
              # Backup blockchain data
              tar -czf /backup/deshchain-$(date +%Y%m%d).tar.gz /deshchain/data
              
              # Upload to cloud storage
              aws s3 cp /backup/deshchain-$(date +%Y%m%d).tar.gz s3://deshchain-backups/testnet/
              
              # Clean old backups (keep 30 days)
              find /backup -name "deshchain-*.tar.gz" -mtime +30 -delete
            env:
            - name: AWS_ACCESS_KEY_ID
              valueFrom:
                secretKeyRef:
                  name: aws-credentials
                  key: access-key-id
            - name: AWS_SECRET_ACCESS_KEY
              valueFrom:
                secretKeyRef:
                  name: aws-credentials
                  key: secret-access-key
            volumeMounts:
            - name: backup-storage
              mountPath: /backup
            - name: validator-data
              mountPath: /deshchain
              readOnly: true
          volumes:
          - name: backup-storage
            persistentVolumeClaim:
              claimName: backup-storage-pvc
          - name: validator-data
            persistentVolumeClaim:
              claimName: validator-data-validator-0
          restartPolicy: OnFailure
```

### 2. Horizontal Pod Autoscaler

```yaml
# hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: deshchain-frontend-hpa
  namespace: deshchain-testnet
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: deshchain-frontend
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: deshchain-sentry-hpa
  namespace: deshchain-testnet
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: deshchain-sentry
  minReplicas: 2
  maxReplicas: 6
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### 3. Upgrade Strategy

```bash
#!/bin/bash
# upgrade-deshchain.sh

VERSION=$1
if [ -z "$VERSION" ]; then
  echo "Usage: $0 <version>"
  exit 1
fi

echo "Upgrading DeshChain testnet to version $VERSION"

# Update validator images
kubectl set image statefulset/deshchain-validator \
  deshchain=deshchain/node:$VERSION \
  -n deshchain-testnet

# Update sentry images
kubectl set image deployment/deshchain-sentry \
  deshchain=deshchain/node:$VERSION \
  -n deshchain-testnet

# Update frontend images
kubectl set image deployment/deshchain-frontend \
  frontend=deshchain/frontend:$VERSION \
  -n deshchain-testnet

# Wait for rollout
kubectl rollout status statefulset/deshchain-validator -n deshchain-testnet
kubectl rollout status deployment/deshchain-sentry -n deshchain-testnet
kubectl rollout status deployment/deshchain-frontend -n deshchain-testnet

echo "Upgrade completed successfully!"
```

## Deployment Commands

### 1. Initial Deployment

```bash
#!/bin/bash
# deploy-testnet.sh

# Create namespaces
kubectl apply -f namespaces.yaml

# Apply storage classes
kubectl apply -f storage-classes.yaml

# Create secrets (ensure these are properly generated first)
kubectl apply -f secrets.yaml

# Deploy ConfigMaps
kubectl apply -f configmaps.yaml

# Deploy database
kubectl apply -f database/

# Deploy validators
kubectl apply -f validator-statefulset.yaml

# Deploy sentry nodes
kubectl apply -f sentry-deployment.yaml

# Deploy services
kubectl apply -f services.yaml

# Deploy frontend applications
kubectl apply -f frontend-deployment.yaml
kubectl apply -f explorer-deployment.yaml
kubectl apply -f faucet-deployment.yaml

# Deploy monitoring
kubectl apply -f monitoring/

# Deploy ingress
kubectl apply -f ingress-nginx.yaml

echo "DeshChain testnet deployed successfully!"
echo "Access points:"
echo "- Frontend: https://testnet.deshchain.com"
echo "- Explorer: https://explorer.testnet.deshchain.com"
echo "- API: https://testnet-api.deshchain.com"
echo "- RPC: https://testnet-rpc.deshchain.com"
echo "- Faucet: https://faucet.testnet.deshchain.com"
```

### 2. Health Check Script

```bash
#!/bin/bash
# health-check.sh

echo "DeshChain Testnet Health Check"
echo "==============================="

# Check node status
echo "Checking validator nodes..."
kubectl get pods -l app=deshchain-validator -n deshchain-testnet

echo "Checking sentry nodes..."
kubectl get pods -l app=deshchain-sentry -n deshchain-testnet

# Check services
echo "Checking services..."
kubectl get svc -n deshchain-testnet

# Check ingress
echo "Checking ingress..."
kubectl get ingress -n deshchain-testnet

# Test RPC endpoints
echo "Testing RPC endpoints..."
curl -s https://testnet-rpc.deshchain.com/status | jq '.result.node_info.network'

# Test API endpoints
echo "Testing API endpoints..."
curl -s https://testnet-api.deshchain.com/cosmos/base/tendermint/v1beta1/node_info | jq '.default_node_info.network'

echo "Health check completed!"
```

## Troubleshooting

### Common Issues and Solutions

#### 1. Pod Startup Issues

```bash
# Check pod logs
kubectl logs -f deshchain-validator-0 -n deshchain-testnet

# Check pod events
kubectl describe pod deshchain-validator-0 -n deshchain-testnet

# Check resource usage
kubectl top pod deshchain-validator-0 -n deshchain-testnet
```

#### 2. Network Connectivity Issues

```bash
# Test internal connectivity
kubectl exec -it deshchain-validator-0 -n deshchain-testnet -- \
  nc -zv deshchain-sentry-svc 26657

# Check DNS resolution
kubectl exec -it deshchain-validator-0 -n deshchain-testnet -- \
  nslookup deshchain-sentry-svc.deshchain-testnet.svc.cluster.local
```

#### 3. Storage Issues

```bash
# Check PVC status
kubectl get pvc -n deshchain-testnet

# Check storage usage
kubectl exec -it deshchain-validator-0 -n deshchain-testnet -- \
  df -h /deshchain
```

#### 4. Performance Tuning

```bash
# Monitor resource usage
kubectl top nodes
kubectl top pods -n deshchain-testnet

# Check cluster events
kubectl get events --sort-by=.metadata.creationTimestamp -n deshchain-testnet
```

### Emergency Procedures

#### 1. Scale Down for Maintenance

```bash
# Scale down non-essential services
kubectl scale deployment deshchain-frontend --replicas=0 -n deshchain-testnet
kubectl scale deployment deshchain-explorer --replicas=0 -n deshchain-testnet

# Backup critical data before maintenance
kubectl exec -it deshchain-validator-0 -n deshchain-testnet -- \
  tar -czf /backup/emergency-backup.tar.gz /deshchain/data
```

#### 2. Disaster Recovery

```bash
# Restore from backup
kubectl create -f disaster-recovery-job.yaml

# Verify data integrity
kubectl exec -it deshchain-validator-0 -n deshchain-testnet -- \
  deshchaind validate-genesis --home /deshchain
```

This comprehensive guide provides everything needed to deploy and manage a DeshChain testnet on Kubernetes, from initial setup through production monitoring and maintenance.