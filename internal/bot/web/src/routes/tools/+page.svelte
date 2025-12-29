<script>
    import { bot } from '$lib/botStore.svelte.js';
</script>

<header class="mb-8">
    <h2 class="text-3xl font-bold">Fishing Module</h2>
</header>

<div class="bg-slate-900 rounded-xl border border-slate-800 divide-y divide-slate-800">
    <div class="p-6 flex items-center justify-between">
        <div>
            <h3 class="font-bold text-lg text-white">Auto-Fishing</h3>
            <p class="text-sm text-slate-400">Automatically use fishing rod on water tiles nearby.</p>
        </div>

        <!-- Toggle Switch -->
        <button
                onclick={bot.toggleFishing}
                class="relative inline-flex h-7 w-12 items-center rounded-full transition-colors focus:outline-none
      {bot.fishingEnabled ? 'bg-orange-600' : 'bg-slate-700'}"
        >
      <span
              class="inline-block h-5 w-5 transform rounded-full bg-white transition-transform
        {bot.fishingEnabled ? 'translate-x-6' : 'translate-x-1'}"
      />
        </button>
    </div>

    <div class="p-4 bg-slate-900/50">
        <div class="flex items-center gap-2 text-xs font-mono">
            <span class="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
            <span class="text-slate-500 text-uppercase">Log:</span>
            <span class="text-slate-300 italic">Waiting for water tiles...</span>
        </div>
    </div>
</div>

<div class="bg-slate-900 rounded-xl border border-slate-800 divide-y divide-slate-800">
    <div class="p-6 flex items-center justify-between">
        <div>
            <h3 class="font-bold text-lg text-white">Lighthack</h3>
            <p class="text-sm text-slate-400">Full light every 100ms.</p>
        </div>

        <!-- Toggle Switch -->
        <button
                onclick={bot.toggleLighthack}
                class="relative inline-flex h-7 w-12 items-center rounded-full transition-colors focus:outline-none
      {bot.lighthackEnabled ? 'bg-orange-600' : 'bg-slate-700'}"
        >
      <span
              class="inline-block h-5 w-5 transform rounded-full bg-white transition-transform
        {bot.lighthackEnabled ? 'translate-x-6' : 'translate-x-1'}"
      />
        </button>
    </div>

</div>