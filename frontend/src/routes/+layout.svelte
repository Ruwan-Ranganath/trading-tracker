<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  
  let currentPath = '/';
  
  onMount(() => {
    const unsubscribe = page.subscribe(($page) => {
      currentPath = $page.url.pathname;
    });
    
    return () => {
      unsubscribe();
    };
  });
  
  const isActive = (path: string) => currentPath === path ? 'bg-gray-900 text-white' : 'text-gray-300 hover:bg-gray-700 hover:text-white';
</script>

<div class="min-h-screen bg-gray-100">
  <!-- Navigation -->
  <nav class="bg-gray-800">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <span class="text-white font-bold">Trading Tracker</span>
          </div>
          <div class="hidden md:block">
            <div class="ml-10 flex items-baseline space-x-4">
              <a 
                href="/" 
                class={`px-3 py-2 rounded-md text-sm font-medium ${isActive('/')}`}
              >
                Dashboard
              </a>
              <a 
                href="/trades" 
                class={`px-3 py-2 rounded-md text-sm font-medium ${isActive('/trades')}`}
              >
                Trades
              </a>
              <a 
                href="/pnl" 
                class={`px-3 py-2 rounded-md text-sm font-medium ${isActive('/pnl')}`}
              >
                P&L Analytics
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  </nav>

  <!-- Page Content -->
  <main class="py-6">
    <slot />
  </main>
  
  <!-- Footer -->
  <footer class="bg-white border-t border-gray-200 mt-8">
    <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <p class="text-center text-gray-500 text-sm">
        &copy; {new Date().getFullYear()} Trading Tracker. All rights reserved.
      </p>
    </div>
  </footer>
</div>
