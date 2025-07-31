import type { RequestHandler } from '@sveltejs/kit';
import { json } from '@sveltejs/kit';
import { error } from '@sveltejs/kit';

export const POST: RequestHandler = async ({ request }) => {
	try {
		const body = await request.json();

		const backendResponse = await fetch('http://backend:8080/links', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(body)
		});

		if (!backendResponse.ok) {
			const message = await backendResponse.text();
			throw error(backendResponse.status, message || 'Backend error');
		}

		const responseData = await backendResponse.json();
		return json(responseData);
	} catch (err: any) {
		console.error('Error handling /links POST:', err);
		throw error(500, 'Internal server error');
	}
};
