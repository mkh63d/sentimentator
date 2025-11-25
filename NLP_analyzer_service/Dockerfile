# Base image
FROM python:3.10-slim

# Set working directory
WORKDIR /app

# Copy requirements first (for caching)
COPY resources/requirements_20251028_191858.txt requirements.txt

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy source code
COPY src/ . 

# Run the main.py when container starts
CMD ["python", "main.py"]
