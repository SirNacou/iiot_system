set windows-shell := ["powershell"]
set shell := ["bash", "-c"]

# This recipe will run if you just type `just` in your terminal
default:
  @echo "Running the default task!"
  @echo "You can also run specific tasks like 'just build' or 'just test'."
  @echo "To see all available recipes, run: just --list" # Suggest the explicit command
    
# Migration for backend
dbmate *ARG:
  docker run --rm -it --network=host -v "$(pwd)/db:/db" ghcr.io/amacneil/dbmate {{ARG}} -d "./backend/migrations/"
  
# Gen code from db
bob-gen:
  # With configuration file
  go run github.com/stephenafamo/bob/gen/bobgen-psql@latest -c ./backend/bob.yaml
  
up:
  docker compose up -d --build