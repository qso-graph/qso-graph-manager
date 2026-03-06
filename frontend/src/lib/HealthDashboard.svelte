<script lang="ts">
  export let entries: any[] = [];
  export let personas: string[] = [];
  export let providers: any[] = [];

  function getEntry(persona: string, providerSlug: string): any | undefined {
    return entries.find(e => e.persona === persona && e.provider === providerSlug);
  }

  $: stored = entries.filter(e => e.hasCreds).length;
  $: missing = entries.filter(e => !e.hasCreds).length;
</script>

<div class="p-6">
  <div class="flex items-center justify-between mb-4">
    <h2 class="text-lg font-semibold text-slate-200">Credential Health</h2>
    <div class="text-sm text-slate-400">
      <span class="text-emerald-400">{stored} stored</span>
      {#if missing > 0}
        <span class="mx-1">&middot;</span>
        <span class="text-red-400">{missing} missing</span>
      {/if}
    </div>
  </div>

  {#if entries.length === 0}
    <div class="text-slate-500 text-center py-8">
      No personas configured yet. Add a persona to see credential health.
    </div>
  {:else}
    <div class="overflow-x-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-slate-700">
            <th class="text-left py-2 px-3 text-slate-400 font-medium">Persona</th>
            {#each providers as prov}
              <th class="text-center py-2 px-3 text-slate-400 font-medium">{prov.name}</th>
            {/each}
          </tr>
        </thead>
        <tbody>
          {#each personas as persona}
            <tr class="border-b border-slate-700/50 hover:bg-slate-800/50">
              <td class="py-2 px-3 text-slate-200 font-medium">{persona}</td>
              {#each providers as prov}
                {@const entry = getEntry(persona, prov.slug)}
                <td class="text-center py-2 px-3">
                  {#if entry?.hasCreds}
                    <span class="text-emerald-400" title="Credentials stored">&#10003;</span>
                  {:else}
                    <span class="text-red-400" title="Credentials missing">&#10007;</span>
                  {/if}
                </td>
              {/each}
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>
