DOCKER_COMPOSE = docker-compose -f ./docker/docker-compose.yml
ENV_FILE = settings/development.env
CONTAINER_NAME = container-golang

GREEN = \033[0;32m
NC = \033[0m

.PHONY: vendor

up: clean create-symlink start

clean:
	@echo "$(GREEN)Limpando containers e imagens...$(NC)"
	./scripts/docker/cleanup.sh

build:
	@echo "$(GREEN)Construindo e iniciando os containers...$(NC)"
	$(DOCKER_COMPOSE) build

start:
	@echo "$(GREEN)Iniciando os containers...$(NC)"
	$(DOCKER_COMPOSE) up

run: #create-symlink
	@echo "$(GREEN)Iniciando aplicação local...$(NC)"
	go run ./cmd/webserver/main.go

stop:
	@echo "$(GREEN)Parando os containers...$(NC)"
	$(DOCKER_COMPOSE) down

create-symlink:
	@echo "$(GREEN)Criando link simbólico para o arquivo de configuração de ambiente...$(NC)"
	rm -rf .env
	ln -s $(ENV_FILE) .env