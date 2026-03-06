<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let personas: any[] = [];
  export let selectedPersona: string = '';

  const dispatch = createEventDispatcher();

  function healthColor(p: any): string {
    if (!p.providers || p.providers.length === 0) return 'bg-gray-500';
    // We'll get health data from parent
    return 'bg-emerald-500';
  }

  function select(name: string) {
    dispatch('select', name);
  }

  function setActive(name: string) {
    dispatch('setActive', name);
  }

  function remove(name: string) {
    dispatch('remove', name);
  }
</script>

<div class="flex flex-col h-full">
  <div class="px-4 py-3 border-b border-slate-700">
    <h2 class="text-sm font-semibold text-slate-300 uppercase tracking-wider">Personas</h2>
  </div>

  <div class="flex-1 overflow-y-auto">
    {#each personas as p}
      <button
        class="w-full text-left px-4 py-3 border-b border-slate-700/50 hover:bg-slate-700/50 transition-colors
               {selectedPersona === p.name ? 'bg-slate-700/70 border-l-2 border-l-amber-400' : ''}"
        on:click={() => select(p.name)}
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2 min-w-0">
            <span class="w-2 h-2 rounded-full flex-shrink-0 {healthColor(p)}"></span>
            <span class="font-medium text-slate-200 truncate">{p.name}</span>
          </div>
          {#if p.isActive}
            <span class="text-xs bg-amber-500/20 text-amber-400 px-1.5 py-0.5 rounded font-medium flex-shrink-0">
              ACTIVE
            </span>
          {/if}
        </div>
        <div class="ml-4 mt-0.5 text-xs text-slate-400">{p.callsign}</div>
      </button>
    {/each}

    {#if personas.length === 0}
      <div class="px-4 py-8 text-center text-slate-500 text-sm">
        No personas yet.<br/>Click "Add" to create one.
      </div>
    {/if}
  </div>

  <div class="px-3 py-2 border-t border-slate-700">
    <button
      class="w-full py-2 text-sm bg-amber-500 hover:bg-amber-400 text-slate-900 font-medium rounded transition-colors"
      on:click={() => dispatch('add')}
    >
      + Add Persona
    </button>
  </div>
</div>
