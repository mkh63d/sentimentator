# ---- Builder stage ----
FROM python:3.10-slim AS builder

# Install minimal build tools for compiling wheels
RUN apt-get update && apt-get install -y --no-install-recommends \
        python3-dev build-essential \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

# Copy only requirements to leverage cache
COPY resources/requirements.txt .

# Upgrade pip, wheel, setuptools
# Build wheels for all dependencies
RUN python -m pip install --upgrade pip wheel setuptools \
 && python -m pip wheel --wheel-dir=/wheels -r requirements.txt \
    --extra-index-url https://download.pytorch.org/whl/cpu

# ---- Runtime stage ----
FROM python:3.10-slim AS runtime

# Create non-root user
RUN useradd --create-home appuser
WORKDIR /app

# Install only runtime libs needed by wheels
RUN apt-get update && apt-get install -y --no-install-recommends \
        libpq5 ca-certificates libstdc++6 \
    && rm -rf /var/lib/apt/lists/*

# Copy pre-built wheels and install (no compilation)
COPY --from=builder /wheels /wheels
RUN python -m pip install --no-cache-dir --no-deps /wheels/* \
 && rm -rf /wheels

# Install SpaCy model in runtime
RUN python -m pip install --no-cache-dir \
    https://github.com/explosion/spacy-models/releases/download/en_core_web_sm-3.8.0/en_core_web_sm-3.8.0-py3-none-any.whl

# Copy application source code
COPY --chown=appuser:appuser src/ ./ 

# Minimal environment
ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1

USER appuser

# Run the application
CMD ["python", "main.py"]
