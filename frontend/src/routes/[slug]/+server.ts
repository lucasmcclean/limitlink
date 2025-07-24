import type { RequestHandler } from '@sveltejs/kit';

export const GET: RequestHandler = forwardToBackend;
export const POST: RequestHandler = forwardToBackend;
export const PUT: RequestHandler = forwardToBackend;
export const DELETE: RequestHandler = forwardToBackend;

async function forwardToBackend({ params, url, request }: any) {
	const slug = params.slug;
	const backendBase = 'http://backend:8080';
	const backendUrl = `${backendBase}/${slug}`;

	const init: RequestInit = {
		method: request.method,
		headers: Object.fromEntries(request.headers)
	};

	if (request.method !== 'GET' && request.method !== 'HEAD') {
		init.body = await request.arrayBuffer();
	}

	const backendResponse = await fetch(backendUrl, init);
	const contentType = backendResponse.headers.get('content-type') || 'application/json';
	const status = backendResponse.status;
	const body = await backendResponse.arrayBuffer();

	return new Response(body, {
		status,
		headers: {
			'content-type': contentType
		}
	});
}
