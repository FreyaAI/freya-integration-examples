# Use the official PHP 7.4 image from Docker Hub
FROM php:7.4-cli

# Install dependencies
RUN apt-get update && apt-get upgrade -y \
    && apt-get install -y openssl zip unzip git libonig-dev libxml2-dev

# Clear cache
RUN apt-get clean && rm -rf /var/lib/apt/lists/*

# Install PHP extensions
RUN docker-php-ext-install pdo_mysql mbstring exif pcntl bcmath opcache

# Install Composer
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer

# Set working directory to /var/www/html
WORKDIR /var/www/html

# Apache is already configured to use /var/www/html as DocumentRoot, so we do not need to remove it
# Instead we will just make sure it is empty
RUN rm -rf /var/www/html/*

# Copy application source to the current working directory
COPY . /var/www/html

# Install PHP dependencies
RUN composer install --no-interaction --prefer-dist --optimize-autoloader

# Change ownership of our applications to the www-data user
RUN chown -R www-data:www-data /var/www/html

# Expose port 80
EXPOSE 80

# Start Apache in the foreground
CMD ["php", "-S", "0.0.0.0:80", "/var/www/html/index.php"]