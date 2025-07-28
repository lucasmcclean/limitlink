<script lang="ts">
	import { goto } from '$app/navigation';
	import Logo from '$lib/components/Logo.svelte';
	import { redirect } from '@sveltejs/kit';

	let { target, expiresIn, maxHits, validFrom, password, slugLength, slugCharset, errors } = $state(
		{
			target: '',
			expiresIn: 7,
			maxHits: '',
			validFrom: '',
			password: '',
			slugLength: '7',
			slugCharset: 'alphanumeric',
			errors: {} as Record<string, string>
		}
	);

	$effect(() => {
		errors.target = !/^https?:\/\//.test(target.trim())
			? 'Must be a valid URL starting with http:// or https://'
			: '';

		errors.expiresIn = expiresIn < 1 || expiresIn > 30 ? 'Must be between 1 and 30' : '';

		errors.slugLength =
			slugLength !== '' && (+slugLength < 6 || +slugLength > 12) ? 'Must be between 6 and 12' : '';

		errors.maxHits = maxHits !== '' && +maxHits <= 0 ? 'Must be greater than 0' : '';

		errors.validFrom =
			validFrom && (isNaN(Date.parse(validFrom)) || new Date(validFrom) < new Date())
				? 'Must be a valid future date'
				: '';

		errors.password = password && password.length < 8 ? 'Must be at least 8 characters' : '';

		errors.slugCharset = !['alphanumeric', 'letters', 'numbers'].includes(slugCharset)
			? 'Invalid character set'
			: '';
	});

	let isInvalid: boolean = $state(true);
	$effect(() => {
		isInvalid = Object.values(errors).some(Boolean);
	});

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();

		if (Object.values(errors).some(Boolean)) return;

		const now = new Date();
		const expiresAt = new Date(now);
		expiresAt.setUTCDate(now.getUTCDate() + Number(expiresIn));

		const payload = {
			target: target.trim(),
			expiresAt: expiresAt.toISOString(),
			maxHits: maxHits !== '' ? +maxHits : null,
			validFrom: validFrom || null,
			password: password || null,
			slugLength: slugLength !== '' ? +slugLength : 7,
			slugCharset
		};

		try {
			const res = await fetch('/api/create', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(payload)
			});

			if (!res.ok) {
				const error = await res.text();
				console.error('Submission error:', error);
			} else {
				const data = await res.json();
				console.log('Success:', data);
				goto(`/admin/${data.adminToken}`);
			}
		} catch (err) {
			console.error('Request failed:', err);
		}
	}
</script>

<svelte:head>
	<title>LimitL.ink</title>

	<link rel="icon" href="/favicon.ico" sizes="any" />
	<link rel="icon" href="/favicon.svg" type="image/svg+xml" />

	<meta
		name="description"
		content="Create and manage short links with a privacy-first, minimal URL shortener. No tracking. No bloat. Just links."
	/>
	<link rel="canonical" href="https://limitl.ink/" />

	<meta property="og:title" content="LimitL.ink" />
	<meta
		property="og:description"
		content="A no-nonsense tool to create and manage short links. Secure. Lightweight. No tracking."
	/>

	<meta property="og:type" content="website" />
	<meta property="og:url" content="https://limitl.ink/" />
	<meta property="og:image" content="https://limitl.ink/img/limitlink-logo.svg" />
	<meta property="og:site_name" content="LimitL.ink" />

	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:title" content="LimitL.ink" />
	<meta
		name="twitter:description"
		content="Create simple, protected short links with zero tracking or analytics. Just what you need."
	/>
	<meta name="twitter:image" content="https://limitl.ink/img/limitlink-logo.svg" />

	<script type="application/ld+json">
		{
			"@context": "https://schema.org",
			"@type": "WebSite",
			"name": "LimitL.ink",
			"url": "https://limitl.ink",
			"description": "Privacy-respecting URL shortener with zero tracking. Minimal and fast."
		}
	</script>
</svelte:head>

<div class="mx-auto max-w-4xl px-4">
	<header class="flex min-h-screen items-center justify-around">
		<div>
			<h1 class="text-6xl font-normal sm:text-7xl md:text-8xl">
				LimitL<span class="opacity-60">.</span>ink
			</h1>
			<h2 class="mt-2 text-xl sm:text-2xl md:text-3xl">Create expiring links in seconds!</h2>
			<div class="mt-4 flex gap-4 text-lg sm:text-xl md:text-2xl">
				<a href="/about">Learn More</a>
				<a href="#new-link">Get Started</a>
			</div>
		</div>
		<div class="hidden w-1/4 sm:block"><Logo /></div>
	</header>

	<main>
		<section id="new-link" class="mx-auto max-w-2xl px-4">
			<form
				novalidate
				onsubmit={handleSubmit}
				class="
        [&_input]:bg-surface [&_select]:bg-surface
        [&_select]:accent-primary [&_input]:accent-primary
        space-y-6
        [&_input]:rounded [&_input]:border-none
        [&_select]:rounded [&_select]:border-none"
			>
				<div>
					<label for="target">Target URL:</label>
					<input
						id="target"
						name="target"
						type="url"
						required
						placeholder="https://example.com"
						bind:value={target}
						aria-invalid={errors.target ? 'true' : 'false'}
						aria-describedby={errors.target ? 'target-error' : undefined}
					/>
					{#if errors.target}
						<p id="target-error" class="error" aria-live="polite">{errors.target}</p>
					{/if}
				</div>

				<div>
					<label for="expires-in"
						>Expires In: <span id="expires-in-days"
							>{expiresIn} {expiresIn === 1 ? 'day' : 'days'}</span
						></label
					>
					<input
						id="expires-in"
						name="expires-in"
						class="block pt-8"
						type="range"
						min="1"
						max="30"
						aria-valuemin="0"
						aria-valuemax="30"
						aria-valuenow={expiresIn}
						bind:value={expiresIn}
						aria-label="Volume"
						aria-labelledby="expires-in-days"
						aria-invalid={errors.expiresIn ? 'true' : 'false'}
						aria-describedby={errors.expiresIn ? 'expires-in-error' : undefined}
					/>
					{#if errors.expiresIn}
						<p id="expires-in-error" class="error" aria-live="polite">{errors.expiresIn}</p>
					{/if}
				</div>

				<div>
					<label for="max-hits">
						Maximum number of visits<span class="opacity-80">(optional)</span>:
					</label>
					<input
						id="max-hits"
						name="max-hits"
						type="number"
						min="1"
						max="999999999"
						placeholder="Unlimited"
						bind:value={maxHits}
						aria-describedby={errors.maxHits ? 'max-hits-error' : undefined}
					/>
					{#if errors.maxHits}
						<p id="max-hits-error" class="error" aria-live="polite">{errors.maxHits}</p>
					{/if}
				</div>

				<div>
					<label for="valid-from">
						Becomes valid on <span class="opacity-80">(optional)</span>:
					</label>
					<input
						id="valid-from"
						name="valid-from"
						type="date"
						bind:value={validFrom}
						aria-describedby={errors.validFrom ? 'valid-from-error' : undefined}
					/>
					{#if errors.validFrom}
						<p id="valid-from-error" class="error" aria-live="polite">{errors.validFrom}</p>
					{/if}
				</div>

				<div>
					<label for="password">
						Password <span class="opacity-80">(optional)</span>:
					</label>
					<input
						id="password"
						name="password"
						type="password"
						placeholder="Leave blank for none"
						bind:value={password}
						aria-describedby={errors.password ? 'password-error' : undefined}
					/>
					{#if errors.password}
						<p id="password-error" class="error" aria-live="polite">{errors.password}</p>
					{/if}
				</div>

				<div>
					<label for="slug-length">
						Length of generated URL <span class="opacity-80">(default: 7)</span>:
					</label>
					<input
						id="slug-length"
						name="slug-length"
						type="number"
						min="6"
						max="12"
						bind:value={slugLength}
						aria-describedby={errors.slugLength ? 'slug-length-error' : undefined}
					/>
					{#if errors.slugLength}
						<p id="slug-length-error" class="error" aria-live="polite">{errors.slugLength}</p>
					{/if}
				</div>

				<div>
					<label for="slug-charset" class="mb-1 block">
						Character set for generated URL <span class="opacity-80">(default: alphanumeric)</span>:
					</label>
					<select
						id="slug-charset"
						name="slug-charset"
						bind:value={slugCharset}
						aria-describedby={errors.slugCharset ? 'slug-charset-error' : undefined}
					>
						<option value="alphanumeric">Alphanumeric</option>
						<option value="letters">Letters</option>
						<option value="numbers">Numbers</option>
					</select>

					{#if errors.slugCharset}
						<p id="slug-charset-error" class="error" aria-live="polite">{errors.slugCharset}</p>
					{/if}
				</div>

				<button type="submit" disabled={isInvalid}>Create Link</button>
			</form>
		</section>
	</main>

	<footer class="mt-auto py-8 opacity-80">&copy; 2025 Lucas McClean — LimitLink™</footer>
</div>
