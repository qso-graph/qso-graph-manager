<script lang="ts">
  import { onMount } from 'svelte';
  import {
    ListPersonas,
    GetPersona,
    AddPersona,
    RemovePersona,
    SetActivePersona,
    GetActivePersona,
    SetCredential,
    GetCredential,
    DeleteCredential,
    DoctorCheck,
    ListProviders,
    EnableProvider,
    DisableProvider,
    GetConfigDir,
  } from '../wailsjs/go/main/App';

  import PersonaList from './lib/PersonaList.svelte';
  import PersonaForm from './lib/PersonaForm.svelte';
  import CredentialForm from './lib/CredentialForm.svelte';
  import HealthDashboard from './lib/HealthDashboard.svelte';
  import ProviderCard from './lib/ProviderCard.svelte';

  // State
  let personas: any[] = [];
  let providers: any[] = [];
  let selectedPersona: string = '';
  let selectedPersonaData: any = null;
  let configDir: string = '';
  let healthEntries: any[] = [];
  let credInfos: Map<string, any> = new Map();

  // UI modes
  type View = 'detail' | 'addPersona' | 'health' | 'configureProvider';
  let view: View = 'detail';
  let configuringProvider: any = null;
  let configuringCredInfo: any = null;
  let statusMessage: string = '';

  onMount(async () => {
    providers = await ListProviders();
    configDir = await GetConfigDir();
    await refresh();
  });

  async function refresh() {
    try {
      personas = await ListPersonas();
      if (personas.length > 0 && !selectedPersona) {
        const active = personas.find(p => p.isActive);
        await selectPersona(active ? active.name : personas[0].name);
      } else if (selectedPersona) {
        await selectPersona(selectedPersona);
      }
    } catch (e: any) {
      showStatus('Error loading personas: ' + e);
    }
  }

  async function selectPersona(name: string) {
    selectedPersona = name;
    view = 'detail';
    try {
      selectedPersonaData = await GetPersona(name);
      await loadCredInfos(name);
    } catch (e: any) {
      showStatus('Error: ' + e);
    }
  }

  async function loadCredInfos(personaName: string) {
    credInfos = new Map();
    for (const prov of providers) {
      try {
        const info = await GetCredential(personaName, prov.slug);
        credInfos.set(prov.slug, info);
      } catch {
        credInfos.set(prov.slug, { hasCreds: false, authType: prov.authType });
      }
    }
    credInfos = credInfos; // trigger reactivity
  }

  async function handleAddPersona(e: CustomEvent) {
    const { name, callsign, start, end } = e.detail;
    try {
      await AddPersona(name, callsign, start, end);
      showStatus(`Persona "${name}" created`);
      await refresh();
      await selectPersona(name);
    } catch (err: any) {
      showStatus('Error: ' + err);
    }
  }

  async function handleRemovePersona(e: CustomEvent) {
    const name = e.detail;
    if (!confirm(`Remove persona "${name}" and all its credentials?`)) return;
    try {
      await RemovePersona(name);
      showStatus(`Persona "${name}" removed`);
      selectedPersona = '';
      selectedPersonaData = null;
      await refresh();
    } catch (err: any) {
      showStatus('Error: ' + err);
    }
  }

  async function handleSetActive(e: CustomEvent) {
    try {
      await SetActivePersona(e.detail);
      showStatus(`Active persona set to "${e.detail}"`);
      await refresh();
    } catch (err: any) {
      showStatus('Error: ' + err);
    }
  }

  async function handleSaveCred(e: CustomEvent) {
    const { persona, provider, username, secret } = e.detail;
    try {
      if (secret) {
        await SetCredential(persona, provider, username, secret);
        showStatus(`Credentials saved for ${provider}`);
      }
      await loadCredInfos(persona);
      view = 'detail';
      configuringProvider = null;
    } catch (err: any) {
      showStatus('Error: ' + err);
    }
  }

  async function handleDeleteCred(e: CustomEvent) {
    const { persona, provider } = e.detail;
    try {
      await DeleteCredential(persona, provider);
      showStatus(`Credentials removed for ${provider}`);
      await loadCredInfos(persona);
      view = 'detail';
      configuringProvider = null;
    } catch (err: any) {
      showStatus('Error: ' + err);
    }
  }

  async function handleEnableProvider(e: CustomEvent) {
    const prov = e.detail;
    const username = selectedPersonaData?.providers?.[prov.slug]?.username || '';
    try {
      await EnableProvider(selectedPersona, prov.slug, username);
      showStatus(`${prov.name} enabled`);
      await selectPersona(selectedPersona);
      await refresh();
    } catch (err: any) {
      showStatus('Error: ' + err);
    }
  }

  async function handleDisableProvider(e: CustomEvent) {
    const prov = e.detail;
    try {
      await DisableProvider(selectedPersona, prov.slug);
      showStatus(`${prov.name} disabled`);
      await selectPersona(selectedPersona);
      await refresh();
    } catch (err: any) {
      showStatus('Error: ' + err);
    }
  }

  function handleConfigureProvider(e: CustomEvent) {
    configuringProvider = e.detail;
    configuringCredInfo = credInfos.get(e.detail.slug) || null;
    view = 'configureProvider';
  }

  async function showHealthDashboard() {
    try {
      healthEntries = await DoctorCheck();
      view = 'health';
    } catch (err: any) {
      showStatus('Error: ' + err);
    }
  }

  function showStatus(msg: string) {
    statusMessage = msg;
    setTimeout(() => { if (statusMessage === msg) statusMessage = ''; }, 4000);
  }

  function isProviderEnabled(slug: string): boolean {
    if (!selectedPersonaData?.providers) return false;
    return slug in selectedPersonaData.providers;
  }

  function providerHasCreds(slug: string): boolean {
    return credInfos.get(slug)?.hasCreds || false;
  }
</script>

<div class="flex h-screen bg-navy-900 text-slate-200">
  <!-- Sidebar -->
  <div class="w-64 flex-shrink-0 bg-slate-800/50 border-r border-slate-700 flex flex-col">
    <!-- Header -->
    <div class="px-4 py-3 border-b border-slate-700 flex items-center gap-2">
      <span class="text-amber-400 font-bold text-base">qso-graph</span>
      <span class="text-slate-400 text-sm">Manager</span>
    </div>

    <!-- Persona List -->
    <div class="flex-1 overflow-hidden">
      <PersonaList
        {personas}
        {selectedPersona}
        on:select={(e) => selectPersona(e.detail)}
        on:setActive={handleSetActive}
        on:remove={handleRemovePersona}
        on:add={() => { view = 'addPersona'; }}
      />
    </div>

    <!-- Footer -->
    <div class="px-3 py-2 border-t border-slate-700 space-y-1">
      <button
        class="w-full py-1.5 text-xs text-slate-400 hover:text-slate-200 hover:bg-slate-700/50 rounded transition-colors"
        on:click={showHealthDashboard}
      >
        Credential Health
      </button>
      <div class="text-xs text-slate-600 text-center truncate" title={configDir}>
        {configDir}
      </div>
    </div>
  </div>

  <!-- Main Content -->
  <div class="flex-1 flex flex-col overflow-hidden">
    <!-- Status bar -->
    {#if statusMessage}
      <div class="px-4 py-2 bg-slate-800 border-b border-slate-700 text-sm text-amber-400">
        {statusMessage}
      </div>
    {/if}

    <div class="flex-1 overflow-y-auto">
      {#if view === 'addPersona'}
        <PersonaForm
          on:save={handleAddPersona}
          on:cancel={() => { view = 'detail'; }}
        />

      {:else if view === 'health'}
        <HealthDashboard
          entries={healthEntries}
          personas={personas.map(p => p.name)}
          {providers}
        />

      {:else if view === 'configureProvider' && configuringProvider}
        <div class="p-6 max-w-lg">
          <button
            class="text-sm text-slate-400 hover:text-slate-200 mb-4"
            on:click={() => { view = 'detail'; configuringProvider = null; }}
          >
            &larr; Back
          </button>
          <h2 class="text-lg font-semibold text-slate-200 mb-4">
            Configure {configuringProvider.name} for {selectedPersona}
          </h2>
          <CredentialForm
            persona={selectedPersona}
            provider={configuringProvider}
            credInfo={configuringCredInfo}
            on:save={handleSaveCred}
            on:delete={handleDeleteCred}
          />
        </div>

      {:else if selectedPersonaData}
        <!-- Persona Detail View -->
        <div class="p-6">
          <!-- Persona header -->
          <div class="flex items-center justify-between mb-6">
            <div>
              <h2 class="text-xl font-semibold text-slate-200">{selectedPersonaData.name}</h2>
              <div class="text-sm text-slate-400 mt-0.5">
                {selectedPersonaData.callsign}
                {#if selectedPersonaData.start}
                  &middot; {selectedPersonaData.start}
                  {#if selectedPersonaData.end}
                    &ndash; {selectedPersonaData.end}
                  {:else}
                    &ndash; present
                  {/if}
                {/if}
              </div>
            </div>
            <div class="flex gap-2">
              {#if !personas.find(p => p.name === selectedPersona)?.isActive}
                <button
                  class="px-3 py-1.5 text-sm bg-amber-500/20 hover:bg-amber-500/30 text-amber-400 rounded transition-colors"
                  on:click={() => { SetActivePersona(selectedPersona).then(() => { showStatus('Active persona set to "' + selectedPersona + '"'); refresh(); }).catch(e => showStatus('Error: ' + e)); }}
                >
                  Set Active
                </button>
              {/if}
              <button
                class="px-3 py-1.5 text-sm bg-red-600/20 hover:bg-red-600/40 text-red-400 rounded transition-colors"
                on:click={() => { if (confirm('Remove persona "' + selectedPersona + '" and all its credentials?')) { RemovePersona(selectedPersona).then(() => { showStatus('Persona "' + selectedPersona + '" removed'); selectedPersona = ''; selectedPersonaData = null; refresh(); }).catch(e => showStatus('Error: ' + e)); } }}
              >
                Remove
              </button>
            </div>
          </div>

          <!-- Providers -->
          <h3 class="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-3">Providers</h3>
          <div class="space-y-2">
            {#each providers as prov}
              <ProviderCard
                provider={prov}
                enabled={isProviderEnabled(prov.slug)}
                hasCreds={providerHasCreds(prov.slug)}
                on:configure={handleConfigureProvider}
                on:enable={handleEnableProvider}
                on:disable={handleDisableProvider}
              />
            {/each}
          </div>
        </div>

      {:else}
        <div class="flex items-center justify-center h-full text-slate-500">
          Select a persona or create a new one
        </div>
      {/if}
    </div>
  </div>
</div>
