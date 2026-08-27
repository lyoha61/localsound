FROM nvidia/cuda:12.6.3-runtime-ubuntu24.04

RUN apt-get update && \
    apt-get install -y python3 python3-pip python3-venv ffmpeg && \
    rm -rf /var/lib/apt/lists/*

RUN python3 -m venv /opt/venv

RUN /opt/venv/bin/pip install --upgrade pip && \
    /opt/venv/bin/pip install whisperx fastapi uvicorn python-multipart

COPY app.py /app/app.py

WORKDIR /app

CMD ["/opt/venv/bin/uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000"]