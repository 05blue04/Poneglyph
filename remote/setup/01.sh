#!/bin/bash
set -eu

# ==================================================================================== #
# VARIABLES
# ==================================================================================== #

# Set the timezone for the server. A full list of available timezones can be found by 
# running timedatectl list-timezones.
TIMEZONE=America/New_York

# Set the name of the new user to create.
USERNAME=NicoRobin

# Prompt to enter a password for the PostgreSQL user 
read -p "Enter password for poneglyph DB user: " DB_PASSWORD

# Force all output to be presented in en_US for the duration of this script. This avoids  
# any "setting locale failed" errors while this script is running, before we have 
# installed support for all locales. Do not change this setting!
export LC_ALL=en_US.UTF-8 

# ==================================================================================== #
# SCRIPT LOGIC
# ==================================================================================== #

# Enable the "universe" repository.
add-apt-repository --yes universe

# Update all software packages.
apt update

# Set the system timezone and install all locales.
timedatectl set-timezone ${TIMEZONE}
apt --yes install locales-all

# Add the new user (and give them sudo privileges).
useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"

# Force a password to be set for the new user the first time they log in.
passwd --delete "${USERNAME}"
chage --lastday 0 "${USERNAME}"

# Copy the SSH keys from the root user to the new user.
rsync --archive --chown=${USERNAME}:${USERNAME} /root/.ssh /home/${USERNAME}

# Configure the firewall to allow SSH, HTTP and HTTPS traffic.
ufw allow 22
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

# Install fail2ban.
apt --yes install fail2ban

# Install the migrate CLI tool.
curl -L https://github.com/pressly/goose/releases/download/v3.25.0/goose_linux_x86_64 -o goose
chmod +x goose
mv goose /usr/local/bin/goose

# Install PostgreSQL.
apt --yes install postgresql

# Set up the DB and create a user account with the password entered earlier.
sudo -i -u postgres psql -c "CREATE DATABASE poneglyph"
sudo -i -u postgres psql -d poneglyph -c "CREATE EXTENSION IF NOT EXISTS citext"
sudo -i -u postgres psql -d poneglyph -c "CREATE ROLE archaeologist WITH LOGIN PASSWORD '${DB_PASSWORD}'"

# Add a DSN for connecting to the database to the system-wide environment 
# variables in the /etc/environment file.
echo "PONEGLYPH_DB_DSN='postgres://archaeologist:${DB_PASSWORD}@localhost/poneglyph'" >> /etc/environment

# Install Caddy (see https://caddyserver.com/docs/install#debian-ubuntu-raspbian).
apt --yes install debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
apt update
apt --yes install caddy

# Upgrade all packages. Using the --force-confnew flag means that configuration 
# files will be replaced if newer ones are available.
apt --yes -o Dpkg::Options::="--force-confnew" upgrade

echo "Script complete! Rebooting..."
reboot
