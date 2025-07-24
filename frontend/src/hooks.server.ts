import type { Handle } from '@sveltejs/kit';
import { sequence } from '@sveltejs/kit/hooks';
import { RetryAfterRateLimiter } from 'sveltekit-rate-limiter/server';
import { env } from '$env/dynamic/private';

const isDev = env.ENV === 'dev';

console.log(isDev);

const httpsRedirect: Handle = async ({ event, resolve }) => {
	const proto = event.request.headers.get('x-forwarded-proto');
	if (proto && proto !== 'https') {
		const host = event.request.headers.get('host');
		return Response.redirect(`https://${host}${event.url.pathname}${event.url.search}`, 301);
	}
	return resolve(event);
};

const limiter = new RetryAfterRateLimiter({
	IP: [
		[20, '10s'],
		[200, 'h']
	],
	IPUA: [60, '10m'],
	cookie: {
		name: 'rl_cookie',
		secret: env.COOKIE_SECRET!,
		rate: isDev ? [9999, 'm'] : [100, '10m'],
		preflight: !isDev
	}
});

const rateLimitMiddleware: Handle = async ({ event, resolve }) => {
	if (isDev) return resolve(event);

	const forwardedFor = event.request.headers.get('x-forwarded-for');
	const ip = forwardedFor?.split(',')[0].trim() ?? crypto.randomUUID();

	(event as any).getClientAddress = () => ip;

	if (event.request.method !== 'GET' && limiter.cookieLimiter?.preflight) {
		await limiter.cookieLimiter.preflight(event);
	}

	const status = await limiter.check(event);
	if (status.limited) {
		console.warn('Rate limit triggered:', {
			IP: ip,
			method: event.request.method,
			path: event.url.pathname,
			reason: status.reason
		});

		return new Response(`Too Many Requests. Retry in ${status.retryAfter}s.`, {
			status: 429,
			headers: {
				'Retry-After': status.retryAfter.toString()
			}
		});
	}

	return resolve(event);
};

const securityAndCaching: Handle = async ({ event, resolve }) => {
	const nonce = crypto.randomUUID();
	event.locals.nonce = nonce;

	const response = await resolve(event, {
		transformPageChunk: ({ html }) => html.replace(/<script/g, `<script nonce="${nonce}"`)
	});

	response.headers.set('strict-transport-security', 'max-age=63072000; includeSubDomains; preload');
	response.headers.set('x-frame-options', 'DENY');
	response.headers.set('x-content-type-options', 'nosniff');
	response.headers.set('referrer-policy', 'strict-origin-when-cross-origin');
	response.headers.set('permissions-policy', 'geolocation=(), microphone=(), camera=()');
	response.headers.set('expect-ct', 'max-age=86400, enforce');
	response.headers.set(
		'content-security-policy',
		`default-src 'self'; script-src 'self' 'nonce-${nonce}'; style-src 'self' 'nonce-${nonce}'; img-src 'self' data:;`
	);

	return response;
};

export const handle: Handle = sequence(httpsRedirect, rateLimitMiddleware, securityAndCaching);
