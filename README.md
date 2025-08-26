<div align="center">

# 🚀 Prometheus cgroup v2 Exporter

<img src="https://raw.githubusercontent.com/prometheus/prometheus/main/documentation/images/prometheus-logo.svg" width="200" alt="Prometheus Logo"/>

### ⚡ High-Performance Monitoring for Modern Container Environments

[![Go](https://img.shields.io/badge/Go-1.21%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Prometheus](https://img.shields.io/badge/Prometheus-Compatible-E6522C?style=for-the-badge&logo=prometheus&logoColor=white)](https://prometheus.io/)
[![Docker](https://img.shields.io/badge/Docker-Available-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://hub.docker.com/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Ready-326CE5?style=for-the-badge&logo=kubernetes&logoColor=white)](https://kubernetes.io/)

[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![GitHub Issues](https://img.shields.io/github/issues/stillalive04/prometheus-cgroup-v2-exporter?style=for-the-badge)](https://github.com/stillalive04/prometheus-cgroup-v2-exporter/issues)
[![GitHub Stars](https://img.shields.io/github/stars/stillalive04/prometheus-cgroup-v2-exporter?style=for-the-badge)](https://github.com/stillalive04/prometheus-cgroup-v2-exporter/stargazers)
[![CI/CD](https://img.shields.io/github/actions/workflow/status/stillalive04/prometheus-cgroup-v2-exporter/ci.yml?style=for-the-badge&label=CI%2FCD)](https://github.com/stillalive04/prometheus-cgroup-v2-exporter/actions)

---

**🎯 A comprehensive, enterprise-grade Prometheus exporter for cgroup v2 metrics**  
*Designed for modern containerized environments and Kubernetes clusters*

[📖 Documentation](#-documentation) • [🚀 Quick Start](#-quick-start) • [📊 Metrics](#-metrics-overview) • [🐳 Docker](#-docker-deployment) • [☸️ Kubernetes](#-kubernetes-deployment)

</div>

---

## ✨ Features

<table>
<tr>
<td width="50%">

### 📊 **Comprehensive Metrics Collection**
```
🔥 CPU Metrics
   ├── Usage & Throttling
   ├── User/System Time
   └── Pressure Stall Info

💾 Memory Metrics  
   ├── Usage & Limits
   ├── Cache & RSS
   ├── OOM Events
   └── Swap Usage

💿 I/O Metrics
   ├── Read/Write Bytes
   ├── Operations Count
   └── Pressure Stalls

🔢 Process Metrics
   ├── PID Counts
   ├── Task States
   └── Process Stats
```

</td>
<td width="50%">

### 🏗️ **Enterprise Architecture**
```
⚡ High Performance
   ├── ~50MB Memory Usage
   ├── <1% CPU Overhead
   └── <100ms Scrape Time

🔧 Scalable Design
   ├── 1000+ Containers
   ├── Smart Caching
   └── Concurrent Collection

🛡️ Production Ready
   ├── Graceful Shutdown
   ├── Health Endpoints
   └── Error Recovery

☁️ Cloud Native
   ├── Kubernetes Ready
   ├── Helm Charts
   └── Service Discovery
```

</td>
</tr>
</table>

<div align="center">

### 🎯 **Why Choose cgroup v2 Exporter?**

| Feature | Traditional Exporters | **cgroup v2 Exporter** |
|---------|----------------------|------------------------|
| 🚀 **Performance** | High overhead | **Optimized & Fast** |
| 📊 **Metrics** | Limited scope | **20+ Comprehensive** |
| 🔧 **Configuration** | Complex setup | **Simple & Flexible** |
| 🐳 **Containers** | Basic support | **Native cgroup v2** |
| ☸️ **Kubernetes** | Manual setup | **Helm Charts Ready** |
| 🛡️ **Security** | Basic | **Enterprise Grade** |

</div>

---

## 📋 Prerequisites

<div align="center">

| Component | Version | Purpose |
|-----------|---------|----------|
| 🐧 **Linux Kernel** | `4.5+` | cgroup v2 support |
| 🔧 **Go** | `1.21+` | Building from source |
| 🐳 **Docker** | `20.10+` | Container deployment |
| ☸️ **Kubernetes** | `1.20+` | Orchestration |
| 📊 **Prometheus** | `2.0+` | Metrics collection |

</div>

### 🔍 **System Requirements**

```bash
# Check cgroup v2 availability
$ mount | grep cgroup2
cgroup2 on /sys/fs/cgroup type cgroup2 (rw,nosuid,nodev,noexec,relatime)

# Verify kernel version
$ uname -r
5.4.0+  # Should be 4.5 or higher
```

---

## 🚀 Quick Start

<div align="center">

### Choose Your Deployment Method

</div>

<table>
<tr>
<td width="25%" align="center">

### 📦 **Binary**
```bash
# Download & Install
wget -qO- https://github.com/stillalive04/prometheus-cgroup-v2-exporter/releases/latest/download/prometheus-cgroup-v2-exporter-linux-amd64.tar.gz | tar xz
sudo mv prometheus-cgroup-v2-exporter /usr/local/bin/

# Run
prometheus-cgroup-v2-exporter
```

</td>
<td width="25%" align="center">

### 🐳 **Docker**
```bash
# Quick Run
docker run -d \
  --name cgroup-exporter \
  --pid=host --privileged \
  -p 9753:9753 \
  -v /sys/fs/cgroup:/sys/fs/cgroup:ro \
  stillalive04/prometheus-cgroup-v2-exporter
```

</td>
<td width="25%" align="center">

### ☸️ **Kubernetes**
```bash
# Helm Install
helm repo add cgroup-exporter \
  https://stillalive04.github.io/prometheus-cgroup-v2-exporter

helm install cgroup-exporter \
  cgroup-exporter/prometheus-cgroup-v2-exporter
```

</td>
<td width="25%" align="center">

### 🔨 **Build**
```bash
# From Source
git clone https://github.com/stillalive04/prometheus-cgroup-v2-exporter.git
cd prometheus-cgroup-v2-exporter
make build
./bin/prometheus-cgroup-v2-exporter
```

</td>
</tr>
</table>

### 🎉 **Verify Installation**

```bash
# Check health
curl http://localhost:9753/health

# View metrics
curl http://localhost:9753/metrics | grep cgroup_
```

---

## 📊 Metrics Overview

<div align="center">

### 🔥 **Available Metrics**

</div>

<table>
<tr>
<td width="50%">

#### 🖥️ **CPU Metrics**
```prometheus
# Usage and timing
cgroup_cpu_usage_seconds_total{cgroup, mode}
cgroup_cpu_user_seconds_total{cgroup}
cgroup_cpu_system_seconds_total{cgroup}

# Throttling information
cgroup_cpu_throttled_seconds_total{cgroup}
cgroup_cpu_throttled_periods_total{cgroup}
cgroup_cpu_periods_total{cgroup}

# Pressure stall information
cgroup_cpu_pressure_seconds_total{cgroup, type}
```

#### 💾 **Memory Metrics**
```prometheus
# Usage and limits
cgroup_memory_usage_bytes{cgroup}
cgroup_memory_limit_bytes{cgroup}
cgroup_memory_cache_bytes{cgroup}
cgroup_memory_rss_bytes{cgroup}

# Swap information
cgroup_memory_swap_usage_bytes{cgroup}
cgroup_memory_swap_limit_bytes{cgroup}

# Events and pressure
cgroup_memory_oom_events_total{cgroup}
cgroup_memory_pressure_seconds_total{cgroup, type}
```

</td>
<td width="50%">

#### 💿 **I/O Metrics**
```prometheus
# Read/Write operations
cgroup_io_read_bytes_total{cgroup, device}
cgroup_io_write_bytes_total{cgroup, device}
cgroup_io_read_operations_total{cgroup, device}
cgroup_io_write_operations_total{cgroup, device}

# Pressure information
cgroup_io_pressure_seconds_total{cgroup, type}
```

#### 🔢 **Process Metrics**
```prometheus
# Process counts
cgroup_processes_count{cgroup}
cgroup_processes_running{cgroup}
cgroup_processes_sleeping{cgroup}
cgroup_processes_zombie{cgroup}
```

#### 📈 **Exporter Metrics**
```prometheus
# Performance metrics
prometheus_cgroup_v2_exporter_scrape_duration_seconds
prometheus_cgroup_v2_exporter_scrape_errors_total
prometheus_cgroup_v2_exporter_last_scrape_timestamp_seconds
prometheus_cgroup_v2_exporter_cgroups_scraped
```

</td>
</tr>
</table>

---

## 🐳 Docker Deployment

### 🚀 **Quick Start with Docker**

```bash
# Simple deployment
docker run -d \
  --name prometheus-cgroup-v2-exporter \
  --pid=host \
  --privileged \
  -p 9753:9753 \
  -v /sys/fs/cgroup:/sys/fs/cgroup:ro \
  -v /proc:/host/proc:ro \
  stillalive04/prometheus-cgroup-v2-exporter:latest
```

### 🔧 **Docker Compose**

```yaml
version: '3.8'
services:
  cgroup-v2-exporter:
    image: stillalive04/prometheus-cgroup-v2-exporter:latest
    container_name: prometheus-cgroup-v2-exporter
    restart: unless-stopped
    pid: host
    privileged: true
    ports:
      - "9753:9753"
    volumes:
      - /sys/fs/cgroup:/sys/fs/cgroup:ro
      - /proc:/host/proc:ro
    environment:
      - LOG_LEVEL=info
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:9753/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### 📊 **Complete Stack with Monitoring**

```bash
# Clone and run complete stack
git clone https://github.com/stillalive04/prometheus-cgroup-v2-exporter.git
cd prometheus-cgroup-v2-exporter
docker-compose up -d

# Access services
echo "🚀 Exporter:   http://localhost:9753"
echo "📊 Prometheus: http://localhost:9090" 
echo "📈 Grafana:    http://localhost:3000 (admin/admin)"
```

---

## ☸️ Kubernetes Deployment

### 🎯 **Helm Chart Installation**

```bash
# Add repository
helm repo add prometheus-cgroup-v2-exporter https://stillalive04.github.io/prometheus-cgroup-v2-exporter
helm repo update

# Install with default values
helm install cgroup-v2-exporter prometheus-cgroup-v2-exporter/prometheus-cgroup-v2-exporter \
  --namespace monitoring \
  --create-namespace

# Install with custom values
helm install cgroup-v2-exporter prometheus-cgroup-v2-exporter/prometheus-cgroup-v2-exporter \
  --namespace monitoring \
  --create-namespace \
  --set resources.requests.memory=128Mi \
  --set resources.limits.memory=256Mi
```

### 📋 **Manual Deployment**

```bash
# Apply DaemonSet
kubectl apply -f https://raw.githubusercontent.com/stillalive04/prometheus-cgroup-v2-exporter/main/deployments/kubernetes/daemonset.yaml

# Check deployment
kubectl get pods -n monitoring -l app.kubernetes.io/name=prometheus-cgroup-v2-exporter
```

### 🔍 **ServiceMonitor for Prometheus Operator**

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: prometheus-cgroup-v2-exporter
  namespace: monitoring
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: prometheus-cgroup-v2-exporter
  endpoints:
  - port: metrics
    interval: 30s
    path: /metrics
```

---

## ⚙️ Configuration

<div align="center">

### 🎛️ **Configuration Options**

</div>

<table>
<tr>
<td width="50%">

#### 📝 **Command Line Flags**
```bash
--web.listen-address=:9753
--web.telemetry-path=/metrics
--cgroup.path=/sys/fs/cgroup
--collector.enable=cpu,memory,io,pids
--collector.disable=
--log.level=info
--log.format=logfmt
--scan.interval=30s
--max.cgroups=10000
--cache.duration=60s
```

</td>
<td width="50%">

#### 📄 **Configuration File**
```yaml
web:
  listen_address: ":9753"
  telemetry_path: "/metrics"

cgroup:
  path: "/sys/fs/cgroup"
  refresh_interval: "15s"

collectors:
  cpu:
    enabled: true
    include_pressure: true
  memory:
    enabled: true
    include_pressure: true
    include_swap: true
  io:
    enabled: true
    include_pressure: true
  pids:
    enabled: true

logging:
  level: "info"
  format: "json"
```

</td>
</tr>
</table>

---

## 📊 Grafana Dashboards

<div align="center">

### 📈 **Pre-built Dashboards**

</div>

<table>
<tr>
<td width="33%" align="center">

#### 🌐 **Overview Dashboard**
- Cluster-wide metrics
- Resource utilization
- Top consumers
- Alert status

[📥 Import Dashboard](dashboards/cgroup-v2-overview.json)

</td>
<td width="33%" align="center">

#### 🔍 **Detailed Dashboard**  
- Per-cgroup metrics
- Historical trends
- Drill-down views
- Performance analysis

[📥 Import Dashboard](dashboards/cgroup-v2-detailed.json)

</td>
<td width="33%" align="center">

#### 🚨 **Troubleshooting Dashboard**
- Error tracking
- Performance issues
- Resource constraints
- Debug information

[📥 Import Dashboard](dashboards/cgroup-v2-troubleshooting.json)

</td>
</tr>
</table>

### 🎨 **Dashboard Screenshots**

<div align="center">
<img src="docs/images/dashboard-overview.png" width="800" alt="Overview Dashboard" style="border-radius: 10px; box-shadow: 0 4px 8px rgba(0,0,0,0.1);"/>
</div>

---

## 🚨 Alerting Rules

### 📊 **Prometheus Alert Rules**

```yaml
groups:
  - name: cgroup-v2-exporter
    rules:
      - alert: HighCPUThrottling
        expr: rate(cgroup_cpu_throttled_seconds_total[5m]) > 0.1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "High CPU throttling in {{ $labels.cgroup }}"
          description: "CPU throttling rate is {{ $value | humanizePercentage }}"

      - alert: HighMemoryUsage
        expr: (cgroup_memory_usage_bytes / cgroup_memory_limit_bytes) > 0.9
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High memory usage in {{ $labels.cgroup }}"
          description: "Memory usage is {{ $value | humanizePercentage }}"

      - alert: OOMKilled
        expr: increase(cgroup_memory_oom_events_total[5m]) > 0
        for: 0m
        labels:
          severity: critical
        annotations:
          summary: "OOM events detected in {{ $labels.cgroup }}"
          description: "{{ $value }} OOM events in the last 5 minutes"
```

---

## 🔧 Development

<div align="center">

### 🛠️ **Development Workflow**

</div>

```bash
# 1. Clone repository
git clone https://github.com/stillalive04/prometheus-cgroup-v2-exporter.git
cd prometheus-cgroup-v2-exporter

# 2. Install dependencies
go mod download

# 3. Run tests
make test

# 4. Build binary
make build

# 5. Run locally
./bin/prometheus-cgroup-v2-exporter --log.level=debug

# 6. Build Docker image
make docker-build

# 7. Run complete stack
docker-compose up -d
```

### 📁 **Project Structure**

```
prometheus-cgroup-v2-exporter/
├── 📁 cmd/prometheus-cgroup-v2-exporter/  # Application entry point
├── 📁 internal/
│   ├── 📁 collector/                      # Metrics collectors
│   ├── 📁 cgroup/                        # cgroup v2 parsing
│   └── 📁 config/                        # Configuration management
├── 📁 deployments/
│   ├── 📁 docker/                        # Docker configurations
│   ├── 📁 kubernetes/                    # K8s manifests
│   └── 📁 helm/                          # Helm charts
├── 📁 dashboards/                        # Grafana dashboards
├── 📁 examples/                          # Usage examples
├── 📁 docs/                              # Documentation
└── 📁 tests/                             # Test suites
```

---

## 🤝 Contributing

<div align="center">

### 💡 **We Welcome Contributions!**

[![Contributors](https://img.shields.io/github/contributors/stillalive04/prometheus-cgroup-v2-exporter?style=for-the-badge)](https://github.com/stillalive04/prometheus-cgroup-v2-exporter/graphs/contributors)
[![Pull Requests](https://img.shields.io/github/issues-pr/stillalive04/prometheus-cgroup-v2-exporter?style=for-the-badge)](https://github.com/stillalive04/prometheus-cgroup-v2-exporter/pulls)

</div>

### 🚀 **How to Contribute**

1. **🍴 Fork** the repository
2. **🌿 Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **✨ Make** your changes
4. **🧪 Add** tests for your changes
5. **✅ Run** the test suite: `make test`
6. **📝 Commit** your changes: `git commit -m 'Add amazing feature'`
7. **📤 Push** to the branch: `git push origin feature/amazing-feature`
8. **🔄 Open** a Pull Request

### 📋 **Development Guidelines**

- **Code Style**: Follow Go conventions and use `gofmt`
- **Testing**: Maintain >80% test coverage
- **Documentation**: Update relevant documentation
- **Commit Messages**: Use conventional commit format

---

## 📖 Documentation

<div align="center">

| Document | Description |
|----------|-------------|
| [🏗️ Architecture](docs/architecture.md) | System design and components |
| [📊 Metrics Reference](docs/metrics-reference.md) | Complete metrics documentation |
| [🔧 Performance Tuning](docs/performance-tuning.md) | Optimization guidelines |
| [🚨 Troubleshooting](docs/troubleshooting.md) | Common issues and solutions |
| [🔌 API Reference](docs/api-reference.md) | HTTP endpoints and responses |

</div>

---

## 🆘 Support & Community

<div align="center">

### 💬 **Get Help**

[![GitHub Issues](https://img.shields.io/badge/Issues-Report%20Bug-red?style=for-the-badge&logo=github)](https://github.com/stillalive04/prometheus-cgroup-v2-exporter/issues)
[![GitHub Discussions](https://img.shields.io/badge/Discussions-Ask%20Question-blue?style=for-the-badge&logo=github)](https://github.com/stillalive04/prometheus-cgroup-v2-exporter/discussions)
[![Documentation](https://img.shields.io/badge/Docs-Read%20More-green?style=for-the-badge&logo=gitbook)](docs/)

</div>

### 🐛 **Found a Bug?**

1. Check [existing issues](https://github.com/stillalive04/prometheus-cgroup-v2-exporter/issues)
2. Create a [new issue](https://github.com/stillalive04/prometheus-cgroup-v2-exporter/issues/new) with:
   - Clear description
   - Steps to reproduce
   - Expected vs actual behavior
   - System information

### 💡 **Have a Feature Request?**

1. Check [discussions](https://github.com/stillalive04/prometheus-cgroup-v2-exporter/discussions)
2. Open a new [feature request](https://github.com/stillalive04/prometheus-cgroup-v2-exporter/issues/new?template=feature_request.md)

---

## 📄 License

<div align="center">

[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

**This project is licensed under the MIT License**  
*See the [LICENSE](LICENSE) file for details*

</div>

---

## 🙏 Acknowledgments

<div align="center">

### 🌟 **Special Thanks**

</div>

<table>
<tr>
<td width="25%" align="center">

#### 📊 **Prometheus Community**
For the excellent monitoring ecosystem

[![Prometheus](https://img.shields.io/badge/Prometheus-Community-E6522C?style=for-the-badge&logo=prometheus)](https://prometheus.io/community/)

</td>
<td width="25%" align="center">

#### 🐧 **Linux Kernel Team**
For cgroup v2 implementation

[![Linux](https://img.shields.io/badge/Linux-Kernel-FCC624?style=for-the-badge&logo=linux)](https://kernel.org/)

</td>
<td width="25%" align="center">

#### 🔧 **Go Community**
For the amazing programming language

[![Go](https://img.shields.io/badge/Go-Community-00ADD8?style=for-the-badge&logo=go)](https://golang.org/community/)

</td>
<td width="25%" align="center">

#### 👥 **Contributors**
Everyone who helped make this better

[![Contributors](https://img.shields.io/badge/All-Contributors-purple?style=for-the-badge&logo=github)](https://github.com/stillalive04/prometheus-cgroup-v2-exporter/graphs/contributors)

</td>
</tr>
</table>

---

## 🗺️ Roadmap

<div align="center">

### 🚀 **What's Coming Next?**

</div>

- [ ] 🔌 **Extended Metrics**: Additional cgroup v2 controllers (net_cls, net_prio)
- [ ] 🏗️ **Multi-Architecture**: ARM64 and ARM32 support
- [ ] ⚡ **Performance**: Zero-allocation metric collection
- [ ] 🎯 **Advanced Filtering**: Regex-based cgroup filtering
- [ ] 🔌 **Plugin System**: Custom metrics exporters
- [ ] 🌐 **Web UI**: Built-in configuration and monitoring interface
- [ ] 📊 **Grafana Plugin**: Native Grafana data source
- [ ] 🔄 **Auto-Discovery**: Dynamic service discovery

---

<div align="center">

### 🌟 **Star History**

[![Star History Chart](https://api.star-history.com/svg?repos=stillalive04/prometheus-cgroup-v2-exporter&type=Date)](https://star-history.com/#stillalive04/prometheus-cgroup-v2-exporter&Date)

---

**Made with ❤️ by the Open Source Community**

*For more information, visit the [project repository](https://github.com/stillalive04/prometheus-cgroup-v2-exporter) on GitHub*

[![GitHub](https://img.shields.io/badge/GitHub-Repository-181717?style=for-the-badge&logo=github)](https://github.com/stillalive04/prometheus-cgroup-v2-exporter)

</div>