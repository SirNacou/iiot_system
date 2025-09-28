FROM node:24-alpine AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable


# Stage 1: Build the React Application
# Use a Node.js base image with enough capabilities to build the project.
FROM base AS build

# Set the working directory for all subsequent commands
WORKDIR /app

# Copy package.json and package-lock.json first to leverage Docker's build cache
# This means npm install only re-runs if these files change.
COPY package.json pnpm-lock.yaml ./

# Install project dependencies
# The --omit=dev flag is technically not needed here if only 'dependencies' are used, 
# but it's a good practice to only install necessary packages if you weren't building.
# For a Vite build, devDependencies are necessary, so install all.
RUN pnpm install

# Copy the rest of the application source code
COPY . .

# Build the application
# The output will be placed in the /app/dist directory
RUN pnpm run build

# ----------------------------------------------------------------------------------

# Stage 2: Serve the Static Files with a Lightweight Web Server (Nginx)
# Use a minimal web server image for the final production container.
FROM nginx:1.29-alpine AS production

# Set the environment variable to run Nginx in the foreground
ENV NODE_ENV="production"

# Remove the default Nginx index.html and configuration
RUN rm -rf /etc/nginx/conf.d/*

# Copy a custom Nginx configuration
# This configures Nginx to serve the static files and handle client-side routing.
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy the static built assets from the 'build' stage into the Nginx static directory
COPY --from=build /app/dist /usr/share/nginx/html

# Expose the port Nginx is listening on (default is 80)
EXPOSE 3000

# The default command (CMD) in nginx:alpine already starts Nginx in the foreground.
# CMD ["nginx", "-g", "daemon off;"]
