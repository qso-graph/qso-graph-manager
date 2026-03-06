<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let editing: boolean = false;
  export let initial: any = null;

  const dispatch = createEventDispatcher();

  let name = initial?.name || '';
  let callsign = initial?.callsign || '';
  let start = initial?.start || '';
  let end = initial?.end || '';
  let error = '';

  $: callsign = callsign.toUpperCase();

  function validate(): boolean {
    if (!name.trim()) { error = 'Name is required'; return false; }
    if (!callsign.trim()) { error = 'Callsign is required'; return false; }
    if (!/^[A-Z0-9]{2,10}(\/[A-Z0-9]+)?$/.test(callsign.trim())) {
      error = 'Invalid callsign format';
      return false;
    }
    if (start && !/^\d{4}-\d{2}-\d{2}$/.test(start)) {
      error = 'Start date must be YYYY-MM-DD';
      return false;
    }
    if (end && !/^\d{4}-\d{2}-\d{2}$/.test(end)) {
      error = 'End date must be YYYY-MM-DD';
      return false;
    }
    error = '';
    return true;
  }

  function submit() {
    if (!validate()) return;
    dispatch('save', { name: name.trim(), callsign: callsign.trim(), start, end });
  }
</script>

<div class="p-6">
  <h2 class="text-lg font-semibold text-slate-200 mb-4">
    {editing ? 'Edit Persona' : 'Add Persona'}
  </h2>

  <div class="space-y-4 max-w-md">
    <div>
      <label class="block text-sm text-slate-400 mb-1" for="name">Name</label>
      <input
        id="name"
        type="text"
        bind:value={name}
        disabled={editing}
        placeholder="e.g., Watson"
        class="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded text-slate-200
               focus:border-amber-500 focus:outline-none disabled:opacity-50"
      />
    </div>

    <div>
      <label class="block text-sm text-slate-400 mb-1" for="callsign">Callsign</label>
      <input
        id="callsign"
        type="text"
        bind:value={callsign}
        placeholder="e.g., KI7MT"
        class="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded text-slate-200
               focus:border-amber-500 focus:outline-none uppercase"
      />
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div>
        <label class="block text-sm text-slate-400 mb-1" for="start">Start Date</label>
        <input
          id="start"
          type="date"
          bind:value={start}
          class="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded text-slate-200
                 focus:border-amber-500 focus:outline-none"
        />
      </div>
      <div>
        <label class="block text-sm text-slate-400 mb-1" for="end">End Date</label>
        <input
          id="end"
          type="date"
          bind:value={end}
          class="w-full px-3 py-2 bg-slate-800 border border-slate-600 rounded text-slate-200
                 focus:border-amber-500 focus:outline-none"
        />
      </div>
    </div>

    {#if error}
      <p class="text-red-400 text-sm">{error}</p>
    {/if}

    <div class="flex gap-3 pt-2">
      <button
        on:click={submit}
        class="px-4 py-2 bg-amber-500 hover:bg-amber-400 text-slate-900 font-medium rounded transition-colors"
      >
        {editing ? 'Update' : 'Create'}
      </button>
      <button
        on:click={() => dispatch('cancel')}
        class="px-4 py-2 bg-slate-700 hover:bg-slate-600 text-slate-300 rounded transition-colors"
      >
        Cancel
      </button>
    </div>
  </div>
</div>
