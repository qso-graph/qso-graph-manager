<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let provider: any;
  export let enabled: boolean = false;
  export let hasCreds: boolean = false;

  const dispatch = createEventDispatcher();

  $: authLabel = provider.authType === 'api_key' ? 'API Key' : 'Password';
</script>

<div class="bg-slate-800 border border-slate-700 rounded-lg p-4 flex items-center justify-between">
  <div class="flex items-center gap-3">
    <div class="w-2 h-2 rounded-full {hasCreds ? 'bg-emerald-500' : enabled ? 'bg-amber-500' : 'bg-slate-600'}"></div>
    <div>
      <div class="font-medium text-slate-200">{provider.name}</div>
      <div class="text-xs text-slate-400">{authLabel} auth &middot; {provider.pypi}</div>
    </div>
  </div>

  <div class="flex items-center gap-2">
    {#if hasCreds}
      <span class="text-xs text-emerald-400">Stored</span>
    {:else if enabled}
      <span class="text-xs text-amber-400">No creds</span>
    {:else}
      <span class="text-xs text-slate-500">Disabled</span>
    {/if}

    <button
      on:click={() => dispatch('configure', provider)}
      class="px-3 py-1 text-xs bg-slate-700 hover:bg-slate-600 text-slate-300 rounded transition-colors"
    >
      Configure
    </button>

    <button
      on:click={() => dispatch(enabled ? 'disable' : 'enable', provider)}
      class="px-3 py-1 text-xs rounded transition-colors
             {enabled ? 'bg-slate-600 hover:bg-slate-500 text-slate-300' : 'bg-amber-500/20 hover:bg-amber-500/30 text-amber-400'}"
    >
      {enabled ? 'Disable' : 'Enable'}
    </button>
  </div>
</div>
