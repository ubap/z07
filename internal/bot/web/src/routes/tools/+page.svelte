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

    <!-- Light Intensity -->
    <div class="p-6 space-y-4">
        <div class="flex justify-between items-center">
            <span class="text-sm font-medium text-slate-300">Light Intensity</span>

            <!-- Editable Number Input -->
            <input
                    type="number"
                    min="0"
                    max="16"
                    value={bot.lighthackLevel}
                    oninput={(e) => {
        let val = parseInt(e.target.value);

        if (val > 16) {
            val = 16;
        } else if (val < 0) {
            val = 0;
        }
        e.target.value = val;

        bot.setLighthackLevel(val);
        }}
                    disabled={!bot.lighthackEnabled}
                    class="w-12 text-center font-mono text-orange-500 bg-slate-950 px-1 py-0.5 rounded border border-slate-800 text-xs
                   focus:ring-1 focus:ring-orange-500 focus:outline-none focus:border-orange-500/50
                   disabled:opacity-30 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
            />
        </div>

        <input
                type="range"
                min="0"
                max="16"
                value={bot.lighthackLevel}
                oninput={(e) => bot.setLighthackLevel(e.target.value)}
                disabled={!bot.lighthackEnabled}
                class="w-full h-2 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-orange-500 disabled:opacity-30 disabled:cursor-not-allowed"
        />

        <div class="flex justify-between text-[10px] font-mono text-slate-600 uppercase">
            <span>Min</span>
            <span>Max</span>
        </div>
    </div>

    <!-- Light Color -->
    <div class="p-6 space-y-4">
        <div class="flex justify-between items-center">
            <span class="text-sm font-medium text-slate-300">Light Color</span>

            <!-- Editable Number Input -->
            <input
                    type="number"
                    min="0"
                    max="255"
                    value={bot.lighthackColor}
                    oninput={(e) => {
        let val = parseInt(e.target.value);

        if (val > 255) {
            val = 255;
        } else if (val < 0) {
            val = 0;
        }
        e.target.value = val;

        bot.setLighthackColor(val);
    }}
                    disabled={!bot.lighthackEnabled}
                    class="w-12 text-center font-mono text-orange-500 bg-slate-950 px-1 py-0.5 rounded border border-slate-800 text-xs
                   focus:ring-1 focus:ring-orange-500 focus:outline-none focus:border-orange-500/50
                   disabled:opacity-30 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
            />
        </div>

        <input
                type="range"
                min="0"
                max="255"
                value={bot.lighthackColor}
                oninput={(e) => bot.setLighthackColor(e.target.value)}
                disabled={!bot.lighthackEnabled}
                class="w-full h-2 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-orange-500 disabled:opacity-30 disabled:cursor-not-allowed"
        />

        <div class="flex justify-between text-[10px] font-mono text-slate-600 uppercase">
            <span>Min</span>
            <span>Max</span>
        </div>
    </div>
</div>

<style>
    /* Optional: Custom styling to make the slider thumb look more like a pro tool */
    input[type='range']::-webkit-slider-thumb {
        appearance: none;
        height: 18px;
        width: 18px;
        border-radius: 50%;
        background: #f97316; /* orange-500 */
        cursor: pointer;
        border: 2px solid #0f172a; /* slate-900 */
        box-shadow: 0 0 10px rgba(0,0,0,0.5);
    }
</style>

