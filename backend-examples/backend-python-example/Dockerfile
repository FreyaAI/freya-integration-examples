# Use the official Python image from the Docker Hub
FROM python:3.12-slim

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

ENV KEYS_DIRECTORY=/code/keys

# Set the working directory in the container
WORKDIR /code

# Upgrade pip
RUN pip install --upgrade pip

# Copy the requirements file into the container at /code
COPY ./requirements.txt /code/requirements.txt

# Install any dependencies
RUN pip install -r requirements.txt

# Copy the current directory contents into the container at /code
COPY . /code/

# Command to run the application
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "80"]
