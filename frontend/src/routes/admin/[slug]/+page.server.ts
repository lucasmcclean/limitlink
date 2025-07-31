import { error } from '@sveltejs/kit';

export async function load({ params }: any) {
	const { slug } = params;

	const res = await fetch(`http://backend:8080/links/${slug}`);

	if (!res.ok) {
		console.log(res);
		throw error(res.status, `Failed to fetch data for slug: ${slug}`);
	}

	const data = await res.json();

	return data;
}
