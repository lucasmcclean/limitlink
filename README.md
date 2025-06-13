# LimitL.ink

**LimitL.ink** is an open-source URL shortener built with Go, Tailwind CSS, and
HTMX. It’s more than just a link shortener; it gives you control over your
links with:

- Usage limits
- Expiration dates
- Click tracking

## Running Locally

To run the app locally, clone the repo and launch it with Docker:

```bash
git clone https://github.com/lucasmcclean/limitlink.git
cd limitlink
docker compose up
```

Once it’s up, visit [http://localhost:8080](http://localhost:8080)

To use the Tailwind CLI (standalone), run:

```bash
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v4.1.8/tailwindcss-{os}-{arch}
chmod +x tailwindcss-{os}-{arch}
mv tailwindcss-{os}-{arch} tailwindcss
```

replacing {os} and {arch} with the correct values for your system (you can find
the available releases on
[Tailwind's GitHub](https://github.com/tailwindlabs/tailwindcss/releases/)).
Note that the version this project is using is **4.1.8**. I will keep this
version updated should it change in the future. More details can be found on
[Tailwind's guide](https://tailwindcss.com/blog/standalone-cli).

Once you have Tailwind installed, you can run it with:

```bash
./tailwindcss -i ./tailwind.css -o ./assets/static/css/tailwind.css --watch \
  --content "./assets/templates/**/*.html,./assets/static/html/**/*.html"
```

or, for production, run:

```bash
./tailwindcss -i ./tailwind.css -o ./assets/static/css/tailwind.css --minify \
  --content "./assets/templates/**/*.html,./assets/static/html/**/*.html"
```

## License

This project is licensed under the MIT License. You’re free to use,
modify, and distribute it for personal or commercial purposes.

> Please note: while the code is open-source, the LimitL.ink name and
> branding are not available for reuse or redistribution without permission.
