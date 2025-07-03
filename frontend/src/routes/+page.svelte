<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>LimitL.ink</title>
    <link rel="icon" href="/img/favicon.ico" sizes="any" />
    <link rel="icon" href="/img/favicon.svg" type="image/svg+xml" />
  </head>

  <body>
    <header class="h-screen w-9/10 mx-auto flex items-center justify-evenly">
      <div class="w-fit flex flex-col justify-center gap-4">
        <h1 class="text-6xl sm:text-7xl md:text-8xl font-display font-extrabold tracking-normal">
          LimitL<span class="opacity-60">.</span>ink
        </h1>
        <h2 class="text-xl sm:text-2xl">Create expiring links in seconds!</h2>
        <div class="flex gap-4 text-lg sm:text-xl font-bold">
          <a href="/about" class="relative group w-fit px-4 py-2 rounded-xl border border-surface overflow-hidden transition-colors duration-300 hover:border-accent">
            <div class="absolute inset-0 bg-surface group-hover:opacity-0 transition-opacity duration-300"></div>
            <span class="relative z-10 text-text">Learn More</span>
          </a>
          <a href="#new-link" class="relative group w-fit px-4 py-2 rounded-xl border border-surface overflow-hidden transition-colors duration-300 hover:border-accent">
            <div class="absolute inset-0 bg-gradient-to-r from-accent to-secondary group-hover:opacity-0 transition-opacity duration-300"></div>
            <span class="relative z-10 text-text">Get Started</span>
          </a>
        </div>
      </div>
      <div class="hidden sm:block h-fit w-1/4">
        <!-- SVG content preserved as-is -->
        <!-- [SVG remains unchanged] -->
      </div>
    </header>

    <main>
      <section id="new-link" class="w-full mb-32 flex flex-col justify-center items-center">
        <form action="/links" method="POST" class="w-full max-w-2xl p-8 flex flex-col gap-4" aria-labelledby="form-heading">
          <h3 id="form-heading" class="text-3xl font-bold text-text">Create a New Link</h3>

          <!-- Error Placeholder -->
          <div id="form-error" tabindex="-1" class="hidden text-error" role="alert" aria-live="polite"></div>

          <!-- Target URL -->
          <div>
            <label for="target" class="block text-lg mb-1">Target URL:</label>
            <input
              type="url"
              id="target"
              name="target"
              placeholder="https://example.com"
              class="w-full px-4 py-2 border rounded-lg border-text"
              required
            />
          </div>

          <!-- Expiry Selection -->
          <div>
            <label for="expires-in" class="block text-lg mb-1">Expires In:</label>
            <select id="expires-in" name="expires-in" class="w-full px-4 py-2 border border-text rounded-lg">
              <option value="1">1 day</option>
              <option value="7" selected>7 days</option>
              <option value="30">30 days</option>
              <option value="custom">Custom date/time</option>
            </select>
          </div>

          <!-- Custom Expiry -->
          <div id="expires-custom-container" hidden>
            <label for="expires-at-display" class="block text-lg mb-1">Custom Expiry:</label>
            <input
              type="datetime-local"
              id="expires-at-display"
              name="expires-at-display"
              class="w-full px-4 py-2 border rounded-lg border-text"
            />
            <input type="hidden" id="expires-at" name="expires-at" disabled />
          </div>

          <!-- Advanced Options Toggle -->
          <div>
            <button type="button" id="toggle-advanced" class="text-accent hover:underline" aria-expanded="false" aria-controls="advanced-options">
              Show Advanced Options
            </button>
          </div>

          <!-- Advanced Options -->
          <div id="advanced-options" hidden class="flex flex-col gap-4">
            <div>
              <label for="max-hits" class="block text-lg mb-1">
                Max Hits <span class="font-semibold opacity-60">(optional)</span>:
              </label>
              <input
                type="number"
                id="max-hits"
                name="max-hits"
                min="1"
                placeholder="Unlimited"
                class="w-full px-4 py-2 border rounded-lg border-text"
              />
            </div>
          </div>

          <!-- Submit -->
          <button type="submit" class="mt-4 px-6 py-3 bg-accent text-white font-bold rounded-xl">
            Create Link
          </button>
        </form>
      </section>
    </main>
  </body>
</html>
