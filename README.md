# 🤖 AI Bot Constructor

This is a microservice system that allows users to quickly create and run their own neural network-based Telegram bots without the need to manage infrastructure

## 🚀 Quick Start

1. **Clone repo**
   ```bash
   git clone https://github.com/BraunKc/ai-bot-constructor.git
   cd ai-bot-constructor
   ```

2. **Initialize Environment Variables**
   
   Generate `.env` files:
   ```bash
   make dotenvs
   ```
   
   > ⚠️ **Important**. This instruction duplicated at MakeFile
   > 1. Get your api key from [OpenRouter](https://openrouter.ai/workspaces/default/keys)
   > 2. Open `./executor-service/.env`
   > 3. Paste api key into `OPEN_ROUTER_TOKEN`

3. **Build and Run**
   
   Build all Docker containers:
   ```bash
   make docker-builds
   ```
   
   Start the services:
   ```bash
   docker-compose up -d
   ```

## 🛠 Tech Stack

| Component | Technology |
| :--- | :--- |
| Communication | gRPC |
| Database | PostgreSQL |
| Message Broker | Kafka |
| Containerization | Docker SDK |
| Bot Interface | Telegram Bot API |
| AI Provider | OpenRouter |

## 🔌 Ports Used

| Service | Port |
|---------|------|
| Auth Gateway | 8080 |
| Orchestrator Gateway | 8081 |
| Auth Service (gRPC) | 50050 |
| Orchestrator Service (gRPC) | 50051 |
| PostgreSQL | 5432 |
| Kafka | 9092, 9093 |

## 🏗 Architecture

The system follows a microservices architecture

![Architecture Diagram](assets/arch.png)
![Flow Visualization](assets/visualization.png)
