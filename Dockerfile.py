FROM python:3.11-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    g++ \
    && rm -rf /var/lib/apt/lists/*

COPY cmd/scoring-service/requirements.txt .
RUN pip3 install --no-cache-dir -r requirements.txt
RUN pip3 install grpcio grpcio-tools

RUN pip3 install spacy
RUN python3 -m spacy download ru_core_news_sm

COPY proto /app/proto

RUN python3 -m grpc_tools.protoc -I /app/proto --python_out=/app --grpc_python_out=/app /app/proto/nlp.proto

COPY cmd/scoring-service/main.py /app/

CMD ["python3", "main.py"]