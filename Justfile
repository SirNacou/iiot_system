set windows-shell := ["powershell"]
set shell := ["bash", "-c"]

# This recipe will run if you just type `just` in your terminal
default:
  @echo "Running the default task!"
  @echo "You can also run specific tasks like 'just build' or 'just test'."
  @echo "To see all available recipes, run: just --list" # Suggest the explicit command
    
# Migration for backend
dbmate *ARG:
  docker run --rm -it --network=host -v "$(pwd)/backend:/db" ghcr.io/amacneil/dbmate {{ARG}}
  
  
up:
  docker compose up -d --build