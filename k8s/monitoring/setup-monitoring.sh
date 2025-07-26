#!/bin/bash

# DeshChain Monitoring Stack Setup Script
# This script sets up comprehensive monitoring for DeshChain testnet

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
NAMESPACE="monitoring"
PROMETHEUS_CHART_VERSION="51.2.0"
LOKI_CHART_VERSION="5.36.1"
ELASTICSEARCH_CHART_VERSION="8.5.1"

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if kubectl is installed
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl is not installed. Please install kubectl first."
        exit 1
    fi
    
    # Check if helm is installed
    if ! command -v helm &> /dev/null; then
        log_error "helm is not installed. Please install Helm first."
        exit 1
    fi
    
    # Check if cluster is accessible
    if ! kubectl cluster-info &> /dev/null; then
        log_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
        exit 1
    fi
    
    log_success "Prerequisites check passed"
}

create_namespace() {
    log_info "Creating monitoring namespace..."
    
    kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -
    
    # Label namespace for network policies
    kubectl label namespace $NAMESPACE name=$NAMESPACE --overwrite
    
    log_success "Namespace '$NAMESPACE' created/updated"
}

add_helm_repositories() {
    log_info "Adding Helm repositories..."
    
    # Add Prometheus community repository
    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    
    # Add Grafana repository
    helm repo add grafana https://grafana.github.io/helm-charts
    
    # Add Elastic repository
    helm repo add elastic https://helm.elastic.co
    
    # Update repositories
    helm repo update
    
    log_success "Helm repositories added and updated"
}

create_secrets() {
    log_info "Creating monitoring secrets..."
    
    # Create basic auth secret for Loki
    kubectl create secret generic loki-basic-auth \
        --from-literal=auth="$(echo -n "admin:$(openssl passwd -apr1 'deshchain-logs-2024!')" | base64 -w 0)" \
        --namespace=$NAMESPACE \
        --dry-run=client -o yaml | kubectl apply -f -
    
    # Create Grafana admin password secret
    kubectl create secret generic grafana-admin-secret \
        --from-literal=admin-password="deshchain-grafana-2024!" \
        --namespace=$NAMESPACE \
        --dry-run=client -o yaml | kubectl apply -f -
    
    log_success "Monitoring secrets created"
}

install_prometheus_stack() {
    log_info "Installing Prometheus monitoring stack..."
    
    # Install kube-prometheus-stack
    helm upgrade --install prometheus-stack prometheus-community/kube-prometheus-stack \
        --namespace=$NAMESPACE \
        --version=$PROMETHEUS_CHART_VERSION \
        --values=prometheus-values.yaml \
        --timeout=600s \
        --wait
    
    log_success "Prometheus stack installed successfully"
}

install_loki_stack() {
    log_info "Installing Loki logging stack..."
    
    # Install Loki
    helm upgrade --install loki grafana/loki \
        --namespace=$NAMESPACE \
        --version=$LOKI_CHART_VERSION \
        --values=loki-values.yaml \
        --timeout=600s \
        --wait
    
    log_success "Loki stack installed successfully"
}

configure_grafana_dashboards() {
    log_info "Configuring Grafana dashboards..."
    
    # Wait for Grafana to be ready
    kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=grafana -n $NAMESPACE --timeout=300s
    
    # Create ConfigMap with custom dashboards
    kubectl create configmap deshchain-dashboards \
        --from-file=grafana-dashboards/ \
        --namespace=$NAMESPACE \
        --dry-run=client -o yaml | kubectl apply -f -
    
    # Restart Grafana to load new dashboards
    kubectl rollout restart deployment/prometheus-stack-grafana -n $NAMESPACE
    
    log_success "Grafana dashboards configured"
}

setup_alerts() {
    log_info "Setting up monitoring alerts..."
    
    # Create custom PrometheusRule for DeshChain alerts
    cat <<EOF | kubectl apply -f -
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: deshchain-custom-alerts
  namespace: $NAMESPACE
  labels:
    prometheus: kube-prometheus
spec:
  groups:
    - name: deshchain.custom
      rules:
        - alert: DeshChainValidatorStakeBelow50k
          expr: deshchain_validator_stake_usd < 50000
          for: 5m
          labels:
            severity: warning
          annotations:
            summary: "Validator stake below minimum"
            description: "Validator {{ \$labels.validator }} has stake of {{ \$value }} USD, below 50k minimum."
        
        - alert: DeshChainReferralCommissionHigh
          expr: deshchain_referral_commission_total > 1000000
          for: 1h
          labels:
            severity: info
          annotations:
            summary: "High referral commission earned"
            description: "Validator {{ \$labels.validator }} has earned {{ \$value }} NAMO in referral commissions."
        
        - alert: DeshChainTokenLaunchEligible
          expr: deshchain_validator_token_launch_eligible == 1
          labels:
            severity: info
          annotations:
            summary: "Validator eligible for token launch"
            description: "Validator {{ \$labels.validator }} is now eligible to launch their validator token."
EOF
    
    log_success "Custom alerts configured"
}

setup_network_policies() {
    log_info "Setting up network policies for monitoring..."
    
    cat <<EOF | kubectl apply -f -
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: monitoring-network-policy
  namespace: $NAMESPACE
spec:
  podSelector: {}
  policyTypes:
    - Ingress
    - Egress
  ingress:
    # Allow ingress from ingress-nginx for Grafana web UI
    - from:
        - namespaceSelector:
            matchLabels:
              name: ingress-nginx
      ports:
        - protocol: TCP
          port: 3000
    
    # Allow Prometheus to scrape metrics from default namespace
    - from:
        - namespaceSelector:
            matchLabels:
              name: default
      ports:
        - protocol: TCP
          port: 26660
    
    # Allow internal communication within monitoring namespace
    - from:
        - podSelector: {}
  
  egress:
    # Allow egress to default namespace for scraping
    - to:
        - namespaceSelector:
            matchLabels:
              name: default
      ports:
        - protocol: TCP
          port: 26660
    
    # Allow DNS
    - to: []
      ports:
        - protocol: UDP
          port: 53
        - protocol: TCP
          port: 53
    
    # Allow HTTPS for external integrations
    - to: []
      ports:
        - protocol: TCP
          port: 443
    
    # Allow internal communication
    - to:
        - podSelector: {}
EOF
    
    log_success "Network policies configured"
}

verify_installation() {
    log_info "Verifying monitoring stack installation..."
    
    # Wait for all deployments to be ready
    kubectl wait --for=condition=available deployment --all -n $NAMESPACE --timeout=600s
    
    # Check Prometheus
    if kubectl get pods -n $NAMESPACE -l app.kubernetes.io/name=prometheus | grep -q Running; then
        log_success "Prometheus is running"
    else
        log_error "Prometheus is not running properly"
        return 1
    fi
    
    # Check Grafana
    if kubectl get pods -n $NAMESPACE -l app.kubernetes.io/name=grafana | grep -q Running; then
        log_success "Grafana is running"
    else
        log_error "Grafana is not running properly"
        return 1
    fi
    
    # Check Loki
    if kubectl get pods -n $NAMESPACE -l app.kubernetes.io/name=loki | grep -q Running; then
        log_success "Loki is running"
    else
        log_error "Loki is not running properly"
        return 1
    fi
    
    # Check AlertManager
    if kubectl get pods -n $NAMESPACE -l app.kubernetes.io/name=alertmanager | grep -q Running; then
        log_success "AlertManager is running"
    else
        log_error "AlertManager is not running properly"
        return 1
    fi
    
    log_success "All monitoring components are running successfully"
}

get_access_info() {
    log_info "Getting access information..."
    
    echo ""
    echo "=== DeshChain Monitoring Stack Access Information ==="
    echo ""
    
    # Grafana access
    echo "Grafana Dashboard:"
    echo "  URL: https://monitoring.testnet.deshchain.com"
    echo "  Username: admin"
    echo "  Password: deshchain-grafana-2024!"
    echo ""
    
    # Prometheus access
    echo "Prometheus (internal):"
    echo "  URL: http://prometheus-stack-prometheus-server.$NAMESPACE.svc.cluster.local:80"
    echo ""
    
    # Loki access
    echo "Loki Logs:"
    echo "  URL: https://logs.testnet.deshchain.com"
    echo "  Username: admin"
    echo "  Password: deshchain-logs-2024!"
    echo ""
    
    # AlertManager access
    echo "AlertManager (internal):"
    echo "  URL: http://prometheus-stack-alertmanager.$NAMESPACE.svc.cluster.local:9093"
    echo ""
    
    echo "=== Port Forward Commands (for local access) ==="
    echo "Grafana:      kubectl port-forward -n $NAMESPACE svc/prometheus-stack-grafana 3000:80"
    echo "Prometheus:   kubectl port-forward -n $NAMESPACE svc/prometheus-stack-prometheus 9090:9090"
    echo "AlertManager: kubectl port-forward -n $NAMESPACE svc/prometheus-stack-alertmanager 9093:9093"
    echo ""
}

main() {
    log_info "Starting DeshChain monitoring stack setup..."
    
    check_prerequisites
    create_namespace
    add_helm_repositories
    create_secrets
    install_prometheus_stack
    install_loki_stack
    configure_grafana_dashboards
    setup_alerts
    setup_network_policies
    verify_installation
    get_access_info
    
    log_success "DeshChain monitoring stack setup completed successfully!"
    echo ""
    echo "You can now access your monitoring dashboards and start monitoring your DeshChain testnet."
    echo "Remember to update the Slack webhook URLs in AlertManager configuration for proper alerting."
}

# Handle script interruption
trap 'log_error "Script interrupted"; exit 1' INT TERM

# Change to script directory
cd "$(dirname "$0")"

# Run main function
main "$@"