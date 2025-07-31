<script lang="ts">
	let { data } = $props();
	console.log(data);
</script>

{#if data}
	<h1>LimitL.ink Admin Page</h1>

	<h2>Read Only</h2>
	<p>Slug: {data.slug}</p>
	<p>adminToken: {data.adminToken}</p>
	<p>target: {data.target}</p>
	<p>hitCount: {data.hitCount}</p>
	<p>maxHits: {data.maxHits}</p>
	<p>validFrom: {data.validFrom}</p>
	<p>createdAt: {data.createdAt}</p>
	<p>expiresAt: {data.expiresAt}</p>
	<p>adminExpiresAt: {data.adminExpiresAt}</p>
	<p>updatedAt: {data.updatedAt}</p>
{:else}
	<p>Loading or failed to fetch data.</p>
{/if}
