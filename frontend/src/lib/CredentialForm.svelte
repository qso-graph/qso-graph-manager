<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let persona: string;
  export let provider: any;
  export let credInfo: any = null;

  const dispatch = createEventDispatcher();

  let username = credInfo?.username || '';
  let secret = '';
  let showSecret = false;
  let error = '';
  let saving = false;

  $: secretLabel = provider?.authType === 'api_key' ? 'API Key' : 'Password';
  $: secretPlaceholder = provider?.authType === 'api_key' ? 'Enter API key' : 'Enter password';

  async function submit() {
    if (!username.trim()) { error = 'Username is required'; return; }
    if (!secret.trim() && !credInfo?.hasCreds) {
      error = `${secretLabel} is required`;
      return;
    }
    error = '';
    saving = true;
    dispatch('save', {
      persona,
      provider: provider.slug,
      username: username.trim(),
      secret: secret.trim()
    });
    saving = false;
  }

  function remove() {
    dispatch('delete', { persona, provider: provider.slug });
  }
</script>

<div class="bg-slate-800 border border-slate-700 rounded-lg p-4">
  <div class="flex items-center justify-between mb-3">
    <h3 class="font-medium text-slate-200">{provider.name}</h3>
    <span class="text-xs px-2 py-0.5 rounded {credInfo?.hasCreds ? 'bg-emerald-500/20 text-emerald-400' : 'bg-red-500/20 text-red-400'}">
      {credInfo?.hasCreds ? 'Stored' : 'Missing'}
    </span>
  </div>

  <div class="space-y-3">
    <div>
      <label class="block text-xs text-slate-400 mb-1">Username</label>
      <input
        type="text"
        bind:value={username}
        placeholder="Enter username"
        class="w-full px-3 py-1.5 bg-slate-900 border border-slate-600 rounded text-sm text-slate-200
               focus:border-amber-500 focus:outline-none"
      />
    </div>

    <div>
      <label class="block text-xs text-slate-400 mb-1">{secretLabel}</label>
      <div class="relative">
        {#if showSecret}
          <input
            type="text"
            bind:value={secret}
            placeholder={credInfo?.hasCreds ? '(unchanged)' : secretPlaceholder}
            class="w-full px-3 py-1.5 bg-slate-900 border border-slate-600 rounded text-sm text-slate-200
                   focus:border-amber-500 focus:outline-none pr-16"
          />
        {:else}
          <input
            type="password"
            bind:value={secret}
            placeholder={credInfo?.hasCreds ? '(unchanged)' : secretPlaceholder}
            class="w-full px-3 py-1.5 bg-slate-900 border border-slate-600 rounded text-sm text-slate-200
                   focus:border-amber-500 focus:outline-none pr-16"
          />
        {/if}
        <button
          type="button"
          on:click={() => showSecret = !showSecret}
          class="absolute right-2 top-1/2 -translate-y-1/2 text-xs text-slate-400 hover:text-slate-300"
        >
          {showSecret ? 'Hide' : 'Show'}
        </button>
      </div>
    </div>

    {#if error}
      <p class="text-red-400 text-xs">{error}</p>
    {/if}

    <div class="flex gap-2 pt-1">
      <button
        on:click={submit}
        disabled={saving}
        class="px-3 py-1.5 text-sm bg-amber-500 hover:bg-amber-400 text-slate-900 font-medium rounded
               transition-colors disabled:opacity-50"
      >
        Save
      </button>
      {#if credInfo?.hasCreds}
        <button
          on:click={remove}
          class="px-3 py-1.5 text-sm bg-red-600/20 hover:bg-red-600/40 text-red-400 rounded transition-colors"
        >
          Delete
        </button>
      {/if}
    </div>
  </div>
</div>
