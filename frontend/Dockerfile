# Usar uma imagem do Python como base
FROM python:3.9-slim

# Definir o diretório de trabalho dentro do container
WORKDIR /app

# Copiar os arquivos do frontend para o container
COPY . /app

# Expor a porta 80 para servir o frontend
EXPOSE 80

# Comando para iniciar o servidor HTTP do Python
CMD ["python3", "-m", "http.server", "80"]

