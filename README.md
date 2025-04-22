# LimitL.ink

**LimitL.ink** is an open-source URL shortener built with Go and Svelte.
It’s more than just a link shortener—it gives you control over your links with:

- Usage limits
- Expiration dates
- Click tracking


## Running Locally

To run the app locally, clone the repo and launch it with Docker:

```bash
git clone https://github.com/lucasmcclean/limitlink.git
cd limitlink
docker compose -f docker-compose.local.yml up --build
```

Once it’s up, visit: [http://localhost:8080](http://localhost:8080)


## Project Structure

- **`main.go`** – Entrypoint for the Go backend
- **`frontend/`** – Contains the Svelte frontend
- **`static/`** – Created at runtime to store the built frontend assets
- **`Dockerfile` / `docker-compose.*.yml`** – Build configs for local
  and production environments


## Purpose

I decided to build this tool because a URL shortener is a classic backend
project, but I wanted to put a spin on it. This project provided me the
perfect opportunity to learn more about one of my favorite programming
languages, Go, while building something genuinely useful for me and
hopefully others too.


## License

This project is licensed under the MIT License. You’re free to use,
modify, and distribute it for personal or commercial purposes.

Please note: while the code is open-source, the LimitL.ink name and
branding are not available for reuse or redistribution without permission.
